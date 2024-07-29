package services

import (
	"context"
	"drivers-service/api/components"
	"drivers-service/api/dto"
	"drivers-service/config"
	"drivers-service/constants"
	"drivers-service/data/models"
	"drivers-service/pkg/logging"
	"drivers-service/pkg/service_errors"
	"drivers-service/pkg/tools"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

type MemberService struct {
	Mongo        *mongo.Database
	Collection   *mongo.Collection
	ctx          context.Context
	logger       logging.Logger
	tokenService *TokenService
	config       *config.Config
	totp         *TotpService
	mailService  *EmailService
}

// NewMemberService creates a new instance of the MemberService struct.
//
// Parameters:
// - db: a pointer to a mongo.Database object representing the MongoDB database connection.
// - cfg: a pointer to a config.Config object representing the application configuration.
// - ctx: a context.Context object representing the context of the function call.
// - collectionName: a string representing the name of the MongoDB collection.
//
// Returns:
// - a pointer to a MemberService struct.
func NewMemberService(db *mongo.Database, cfg *config.Config, ctx context.Context, collectionName string) *MemberService {
	return &MemberService{
		Mongo:        db,
		Collection:   db.Collection(collectionName),
		ctx:          ctx,
		logger:       logging.NewLogger(cfg),
		config:       cfg,
		mailService:  NewEmailService(cfg),
		totp:         NewTotpService(db, cfg, ctx),
		tokenService: NewTokenService(cfg),
	}
}

// Register registers a new member in the system.
//
// The function takes a pointer to a MemberRegistration struct as input, which contains the details of the member to be registered.
// It returns an error if there was any issue during the registration process.
func (m *MemberService) Register(memberCreate *dto.MemberRegistration) error {

	memberCreate.ID = tools.GenerateUUID()
	memberCreate.CreatedAt = time.Now()
	memberCreate.UpdatedAt = memberCreate.CreatedAt
	memberCreate.Birthday = memberCreate.CreatedAt
	memberCreate.Verification = tools.GenerateUUID()

	rolesCollection := m.Mongo.Collection("roles")

	role, err := findRoleByName(m.ctx, rolesCollection, "member")
	if err != nil {
		panic(err)
	}
	memberCreate.Role = []models.MemberRole{*role}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := m.Collection.Indexes().CreateOne(m.ctx, index); err != nil {
		m.logger.Error(logging.MongoDB, logging.CreateIndex, err.Error(), nil)
		return err
	}

	bp, err := bcrypt.GenerateFromPassword([]byte(memberCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		m.logger.Error(logging.MongoDB, logging.HashPassword, err.Error(), nil)
		return errors.New("Пароль не может быть хеширован")
	}

	memberCreate.Password = string(bp)

	if m.config.SMTP.Auth {
		code := memberCreate.Verification
		firstName := memberCreate.FirstName
		emailData := models.EmailData{
			URL:       fmt.Sprintf("http://%s:%d/api/v1/member/verifyemail/%s", m.config.Server.Domain, m.config.Server.IPort, code),
			FirstName: firstName,
			Subject:   "Your account verification code",
		}

		err = m.mailService.SendEmail(memberCreate.Email, &emailData, "verificationCode.html")
		if err != nil {
			m.logger.Error(logging.Email, logging.SendEmail, err.Error(), nil)
			return err
		}
	}

	res, err := m.Collection.InsertOne(m.ctx, memberCreate)
	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			m.logger.Error(logging.MongoDB, logging.Insert, err.Error(), nil)
			return &service_errors.ServiceError{EndUserMessage: service_errors.EmailExists}
		}
		return err
	}
	var member *dto.MemberResponse
	query := bson.M{"_id": res.InsertedID}
	if err = m.Collection.FindOne(m.ctx, query).Decode(&member); err != nil {
		return err
	}

	return err
}

// Login authenticates a member using their email and password, and returns a token if successful.
//
// Parameters:
// - req: A pointer to a MemberAuth DTO containing the email and password of the member.
//
// Returns:
// - tokenDetail: A pointer to a TokenDetail DTO containing the token, mobile number, and email of the member.
// - error: An error if the authentication fails.
func (m *MemberService) Login(req *dto.MemberAuth) (*dto.TokenDetail, error) {

	exists, err := m.ExistEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, &service_errors.ServiceError{EndUserMessage: service_errors.EmailNotExists}
	}
	var member models.Member
	query := bson.M{"email": req.Email}
	err = m.Collection.FindOne(m.ctx, query).Decode(&member)

	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(member.Password), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	if member.IsTotp == true {
		isVerify := m.totp.codeValidate(req.Code, member.SecretQrCode)
		if !isVerify {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.TotpNotValid}
		}
	}

	tdto := tokenDto{Id: member.ID, MobileNumber: member.Phone, Email: member.Email}

	for _, role := range member.Role {
		tdto.Roles = append(tdto.Roles, role.Name)
	}

	fmt.Printf("tdto: %v", tdto)

	token, err := m.tokenService.GenerateToken(&tdto)
	if err != nil {
		return nil, err
	}
	return token, nil

}

// FindAll возвращает отфильтрованный список объектов MemberResponse из сервиса MemberService.
//
// Параметры:
//
// - page: номер страницы результатов для извлечения (по умолчанию: 1)
//
// - limit: максимальное количество результатов на странице (по умолчанию: 10)
//
// Возвращает:
//
// - []*dto.MemberResponse: срез объектов MemberResponse, представляющих извлеченных членов
//
// - error: ошибка, если операция извлечения не удалась
func (ms *MemberService) FindAll(page int, limit int) ([]*dto.MemberResponse, error) {

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	skip := (page - 1) * limit

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))
	opt.SetSort(bson.M{"created_at": -1})

	query := bson.M{}

	cursor, err := ms.Collection.Find(ms.ctx, query, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ms.ctx)

	var members []*dto.MemberResponse

	for cursor.Next(ms.ctx) {
		var member dto.MemberResponse
		err := cursor.Decode(&member)
		if err != nil {
			return nil, err
		}
		members = append(members, &member)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return members, err
}

// GetMemberById возвращает информацию о члене по его идентификатору.
//
// id - идентификатор члена.
// Возвращает объект dto.MemberResponse с информацией о члене или ошибку, если не удалось найти члена.
func (m *MemberService) GetMemberById(id string) (*dto.MemberResponse, error) {
	var memberRes *dto.MemberResponse
	query := bson.M{"_id": id}
	err := m.Collection.FindOne(m.ctx, query).Decode(&memberRes)
	if err != nil {
		return nil, err
	}
	return memberRes, nil

}

// Update updates a member's information in the database.
//
// Parameters:
// - res: A pointer to a MemberUpdate DTO containing the updated member information.
// - ctx: A pointer to a gin.Context object for handling the HTTP request.
//
// Returns:
// - memberRes: A pointer to a MemberResponse DTO containing the updated member information.
// - error: An error if the update operation fails.
func (m *MemberService) Update(ctx *gin.Context, res *dto.MemberUpdate) (*dto.MemberResponse, error) {

	member, err := m.getEmailUser(ctx)
	if err != nil {
		return nil, err
	}

	updateData, _ := m.updateDataValidation(res)

	query := bson.M{"email": member.Email}

	_, err = m.Collection.UpdateOne(m.ctx, query, bson.M{"$set": updateData})
	if err != nil {
		return nil, err
	}

	var memberRes *dto.MemberResponse
	query = bson.M{"_id": member.ID}
	err = m.Collection.FindOne(m.ctx, query).Decode(&memberRes)
	if err != nil {
		return nil, err
	}

	return memberRes, nil

}

// updateDataValidation validates and prepares the update data for a member.
//
// Parameters:
// - updateData: A pointer to a MemberUpdate DTO containing the updated member information.
//
// Returns:
// - update: A bson.M object containing the update data.
// - error: An error if the update data is invalid or empty.
func (m *MemberService) updateDataValidation(updateData *dto.MemberUpdate) (bson.M, error) {
	update := bson.M{}

	// Добавляем поля для обновления только если они не пустые
	if updateData.FirstName != "" {
		update["firstName"] = updateData.FirstName
	}
	if updateData.LastName != "" {
		update["lastName"] = updateData.LastName
	}
	if updateData.MiddleName != "" {
		update["middleName"] = updateData.MiddleName
	}
	if !updateData.Birthday.IsZero() {
		update["birthday"] = updateData.Birthday
	}
	if (dto.MemberLocationResponse{}) != updateData.Location { // Проверяем, не является ли Location пустой структурой
		update["location"] = updateData.Location
	}

	locationUpdate := bson.M{}
	if updateData.Location.Address != "" {
		locationUpdate["address"] = updateData.Location.Address
	}
	if updateData.Location.City != "" {
		locationUpdate["city"] = updateData.Location.City
	}
	if updateData.Location.Postcode != "" {
		locationUpdate["postcode"] = updateData.Location.Postcode
	}
	if updateData.Location.Country != "" {
		locationUpdate["country"] = updateData.Location.Country
	}

	// Если есть обновления по Location, добавляем в общий update
	if len(locationUpdate) > 0 {
		update["location"] = locationUpdate
	}

	update["updatedAt"] = time.Now() // Обновляем время

	if len(update) == 0 {
		return nil, errors.New("нет полей для обновления")
	}

	return update, nil
}

// getEmailUser retrieves a Member from the database based on the email contained in the Authorization header of the given gin.Context.
//
// Parameters:
// - ctx: A pointer to a gin.Context object representing the HTTP request and response.
//
// Returns:
// - models.Member: The Member object retrieved from the database.
// - error: An error if the retrieval fails or if the Authorization header is missing or invalid.
func (m *MemberService) getEmailUser(ctx *gin.Context) (models.Member, error) {
	claimMap := map[string]interface{}{}
	auth := ctx.GetHeader(constants.AuthorizationHeaderKey)
	fmt.Errorf("claimMap: %v", claimMap)
	tokenParts := strings.Split(auth, " ")
	if tokenParts[0] != "Bearer" {
		err := &service_errors.ServiceError{EndUserMessage: service_errors.TokenBearer}
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			components.GenerateBaseResponseWithError(nil, false, components.AuthError, err))
		return models.Member{}, err
	}

	claimMap, err := m.tokenService.GetClaims(tokenParts[1])
	fmt.Errorf("claimMap: %v", claimMap)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			components.GenerateBaseResponseWithError(nil, false, components.AuthError, err))
		return models.Member{}, err
	}

	var member models.Member
	query := bson.M{"email": claimMap[constants.EmailKey]}
	err = m.Collection.FindOne(m.ctx, query).Decode(&member)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized,
			components.GenerateBaseResponseWithError(nil, false, components.AuthError, err))
		return models.Member{}, err
	}
	return member, nil

}

// ExistEmail checks if an email exists in the database.
//
// Parameters:
// - email: The email to check.
//
// Returns:
// - bool: True if the email exists, false otherwise.
// - error: An error if there was an issue checking the email.
func (m *MemberService) ExistEmail(email string) (bool, error) {
	query := bson.M{"email": email}
	var member *models.Member
	err := m.Collection.FindOne(m.ctx, query).Decode(&member)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// chooseRole finds and returns a MemberRole from the given member's list of roles based on the provided role name.
//
// Parameters:
// - member: A pointer to a Member struct representing the member whose roles are being searched.
// - roleName: A string representing the name of the role to search for.
//
// Returns:
// - MemberRole: The MemberRole struct that matches the provided role name, or an empty MemberRole struct and an error if no match is found.
// - error: An error if the provided role name is not found in the member's roles.
func chooseRole(member *models.Member, roleName string) (models.MemberRole, error) {
	for _, role := range member.Role {
		if role.Name == roleName {
			return role, nil
		}
	}
	return models.MemberRole{}, fmt.Errorf("role '%s' not found for member", roleName)
}

// findRoleByName retrieves a MemberRole from the given collection based on the provided name.
//
// Parameters:
// - ctx: The context.Context to control the execution of the function.
// - collection: The *mongo.Collection to search for the MemberRole.
// - name: The name of the MemberRole to search for.
//
// Returns:
// - *models.MemberRole: The retrieved MemberRole, or nil if not found.
// - error: An error if there was a problem retrieving the MemberRole.
func findRoleByName(ctx context.Context, collection *mongo.Collection, name string) (*models.MemberRole, error) {
	var role *models.MemberRole
	query := bson.M{"name": name}
	err := collection.FindOne(ctx, query).Decode(&role)
	if err != nil {
		return nil, err
	}
	return role, nil
}
