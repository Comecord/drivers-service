package services

import (
	"context"
	"crm-glonass/api/dto"
	"crm-glonass/config"
	"crm-glonass/data/models"
	"crm-glonass/pkg/logging"
	"crm-glonass/pkg/service_errors"
	"crm-glonass/pkg/tools"
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
}

func NewMemberService(db *mongo.Database, cfg *config.Config, ctx context.Context, collectionName string) MemberInterface {
	return &MemberService{
		Mongo:      db,
		Collection: db.Collection(collectionName),
		ctx:        ctx,
		logger:     logging.NewLogger(cfg),
		tokenService: &TokenService{
			logger: logging.NewLogger(cfg),
			cfg:    cfg,
		},
	}
}

func (m *MemberService) Register(memberCreate *dto.MemberCreate) error {

	memberCreate.ID = tools.GenerateUUID()
	memberCreate.CreateAt = time.Now()
	memberCreate.UpdatedAt = memberCreate.CreateAt

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

	m.logger.Infof("memberCreate: %v", memberCreate)

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

	exists, err := m.existEmail(req.Email)
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

	tdto := tokenDto{Id: member.ID, MobileNumber: member.Phone, Email: member.Email}

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

func (m *MemberService) existEmail(email string) (bool, error) {
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
