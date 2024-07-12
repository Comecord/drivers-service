package controllers

import (
	"context"
	"crm-glonass/api/components"
	"crm-glonass/api/dto"
	"crm-glonass/api/services"
	"crm-glonass/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type MemberController struct {
	service *services.MemberService
}

func NewMemberController(db *mongo.Database, ctx context.Context, conf *config.Config) *MemberController {
	service, ok := services.NewMemberService(db, conf, ctx, "members").(*services.MemberService)
	if !ok {
		return nil
	}
	return &MemberController{
		service: service,
	}
}

// Register Member godoc
// @Summary Registration a member
// @Description Registration a member
// @Tags Members
// @Accept json
// @produces json
// @Param Request body dto.MemberCreate true "member"
// @Success 201 {object} components.BaseHttpResponse "Success"
// @Failure 400 {object} components.BaseHttpResponse "Failed"
// @Failure 409 {object} components.BaseHttpResponse "Failed"
// @Router /v1/members/ [post]
func (mc *MemberController) Register(ctx *gin.Context) {
	member := new(dto.MemberCreate)
	err := ctx.ShouldBindJSON(member)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			components.GenerateBaseResponseWithValidationError(nil, false, components.ValidationError, err))
		return
	}

	err = mc.service.Register(member)
	if err != nil {
		ctx.AbortWithStatusJSON(components.TranslateErrorToStatusCode(err),
			components.GenerateBaseResponseWithError(nil, false, components.InternalError, err))
		return
	}
	ctx.JSON(http.StatusCreated, components.GenerateBaseResponse(nil, true, components.Success))

}