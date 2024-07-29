package services

import (
	"drivers-service/api/dto"
	"github.com/golang-jwt/jwt"
)

type TokenInterface interface {
	GenerateToken(token *tokenDto) (*dto.TokenDetail, error)
	VerifyToken(token string) (*jwt.Token, error)
	GetClaims(token string) (claimMap map[string]interface{}, err error)
}
