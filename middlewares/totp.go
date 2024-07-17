package middlewares

import (
	"bytes"
	"crm-glonass/data/models"
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"fmt"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"image/png"
	"os"
)

type TOTPGenerator struct {
	Key *otp.Key
}

func GeneratePassCode(payload *models.AuthTOTP) (*otp.Key, error) {
	passcode, _ := totp.Generate(totp.GenerateOpts{
		Issuer:      payload.Issuer,
		AccountName: payload.AccountName,
	})

	// Конвертация TOTP ключа в PNG
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

	fmt.Println("QR code saved to uploads/" + filename)

	return passcode, err
}

func ValidatePassCode(passcode string, utf8string string) bool {
	secret := base32.StdEncoding.EncodeToString([]byte(utf8string))
	valid := totp.Validate(passcode, secret)
	return valid
}
