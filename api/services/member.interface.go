package services

import (
	"drivers-service/api/dto"
	"github.com/gin-gonic/gin"
)

type MemberInterface interface {
	Register(createDTO *dto.MemberRegistration) error
	Login(*dto.MemberAuth) (*dto.TokenDetail, error)
	Update(updateDTO *dto.MemberUpdate, ctx *gin.Context) (*dto.MemberResponse, error)
}
