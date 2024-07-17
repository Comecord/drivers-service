package services

import (
	"context"
	"drivers-service/api/dto"
	"drivers-service/config"
	"drivers-service/data/models"
	"drivers-service/pkg/logging"
	"drivers-service/pkg/service_errors"
	"drivers-service/pkg/tools"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
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

func NewMemberService(db *mongo.Database, cfg *config.Config, ctx context.Context, collectionName string) MemberInterface {
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

func (m *MemberService) Register(memberCreate *dto.MemberRegistration) error {

	memberCreate.ID = tools.GenerateUUID()
	memberCreate.CreateAt = time.Now()
	memberCreate.UpdatedAt = memberCreate.CreateAt
	memberCreate.Birthday = memberCreate.CreateAt
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

func (m *MemberService) Update(req *dto.MemberUpdate) (*dto.MemberResponse, error) {
	return nil, nil

}

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

func chooseRole(member *models.Member, roleName string) (models.MemberRole, error) {
	for _, role := range member.Role {
		if role.Name == roleName {
			return role, nil
		}
	}
	return models.MemberRole{}, fmt.Errorf("role '%s' not found for member", roleName)
}

func findRoleByName(ctx context.Context, collection *mongo.Collection, name string) (*models.MemberRole, error) {
	var role *models.MemberRole
	query := bson.M{"name": name}
	err := collection.FindOne(ctx, query).Decode(&role)
	if err != nil {
		return nil, err
	}
	return role, nil
}
