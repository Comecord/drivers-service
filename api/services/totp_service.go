package services

import (
	"bytes"
	"context"
	"crypto/sha256"
	"drivers-service/api/dto"
	"drivers-service/config"
	"drivers-service/data/models"
	"drivers-service/pkg/logging"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"image/png"
	"os"
	"time"
)

type TotpService struct {
	Mongo  *mongo.Database
	ctx    context.Context
	logger logging.Logger
	config *config.Config
}

func NewTotpService(db *mongo.Database, cfg *config.Config, ctx context.Context) *TotpService {
	return &TotpService{
		Mongo:  db,
		ctx:    ctx,
		logger: logging.NewLogger(cfg),
		config: cfg,
	}
}

func (t *TotpService) GenerateTotp(payload *dto.TotpRequest) (*dto.TotpResponse, error) {
	collection := t.Mongo.Collection("members")
	var member *models.Member
	err := collection.FindOne(t.ctx, bson.M{"email": payload.AccountName}).Decode(&member)
	if err != nil {
		t.logger.Error(logging.MongoDB, logging.Find, err.Error(), nil)
		return nil, err
	}

	secret, _ := t.generateSecretKey(payload)
	t.logger.Infof("Secret Key: %s", secret.SecretKey)

	totp := &dto.TotpResponse{
		SecretKey: secret.SecretKey,
		TotpURL:   secret.TotpURL,
		QrCode:    fmt.Sprintf("http://localhost:%v/uploads/%s", t.config.Server.IPort, secret.FileName),
	}

	query := bson.D{{Key: "_id", Value: member.ID}}

	member.SecretQrCode = secret.SecretKey
	member.FileQRCode = secret.FileName
	member.UpdatedAt = time.Now()

	doc := bson.D{{Key: "$set", Value: member}}

	collection.FindOneAndUpdate(t.ctx, query, doc, options.FindOneAndUpdate().SetReturnDocument(1))

	return totp, nil
}

func (t *TotpService) Active(code *dto.TotpCodeVerify) (string, error) {
	collection := t.Mongo.Collection("members")
	var member *models.Member
	err := collection.FindOne(t.ctx, bson.M{"email": code.Email, "isTotp": false}).Decode(&member)
	if err != nil {
		t.logger.Error(logging.MongoDB, logging.Find, err.Error(), nil)
		return "", err
	}
	fmt.Printf("RESPONSE -------- %+v\n", member)
	isValid := t.codeValidate(code.Code, member.SecretQrCode)
	if !isValid {
		return "", errors.New("invalid code")
	} else {

		query := bson.D{{Key: "_id", Value: member.ID}}
		doc := bson.D{{Key: "$set", Value: member}}
		member.IsTotp = true
		collection.FindOneAndUpdate(t.ctx, query, doc, options.FindOneAndUpdate().SetReturnDocument(1))
		t.logger.Info(logging.MongoDB, logging.Update, "Totp activated", nil)
	}

	return "Двухфакторная аутентификация активирована", nil
}

func (t *TotpService) codeValidate(passcode string, secret string) bool {
	valid := totp.Validate(passcode, secret)
	return valid
}

func (t *TotpService) generateSecretKey(payload *dto.TotpRequest) (dto.TOTP, error) {
	passcode, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      payload.Issuer,
		AccountName: payload.AccountName,
	})

	var buf bytes.Buffer
	img, err := passcode.Image(200, 200)
	if err != nil {
		panic(err)
	}
	png.Encode(&buf, img)

	// Генерация хэша для имени файла
	hash := sha256.Sum256([]byte(passcode.Secret()))
	filename := hex.EncodeToString(hash[:]) + ".png"

	// Сохранение QR-кода в папку uploads
	file, err := os.Create("./uploads/" + filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}

	result := dto.TOTP{
		Issuer:      payload.Issuer,
		AccountName: payload.AccountName,
		SecretKey:   passcode.Secret(),
		FileName:    filename,
		TotpURL:     passcode.URL(),
	}
	fmt.Println("QR code saved to uploads/" + filename)

	return result, nil
}
