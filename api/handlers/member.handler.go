package handlers

import (
	"context"
	"drivers-service/api/components"
	"drivers-service/api/dto"
	"drivers-service/api/services"
	"drivers-service/config"
	"fmt"
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
//
//	@Summary		Registration a member
//	@Description	Registration a member
//	@Tags			Auth
//	@Accept			json
//	@produces		json
//	@Param			Request	body		dto.MemberRegistration		true	"member"
//	@Success		201		{object}	components.BaseHttpResponse	"Success"
//	@Failure		400		{object}	components.BaseHttpResponse	"Failed"
//	@Failure		409		{object}	components.BaseHttpResponse	"Failed"
//	@Router			/api/v1/members/register [post]
func (mc *MemberController) Register(ctx *gin.Context) {
	member := new(dto.MemberRegistration)
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
	ctx.JSON(http.StatusCreated, components.GenerateBaseResponse("Member created", true, components.Success))

}

// Login Member godoc
//
//	@Summary		Login a member
//	@Description	Login a member
//	@Tags			Auth
//	@Accept			json
//	@produces		json
//	@Param			Request	body		dto.MemberAuth				true	"member"
//	@Success		200		{object}	components.BaseHttpResponse	"Success"
//	@Failure		400		{object}	components.BaseHttpResponse	"Failed"
//	@Failure		409		{object}	components.BaseHttpResponse	"Failed"
//	@Router			/api/v1/members/login [post]
func (mc *MemberController) Login(ctx *gin.Context) {
	req := new(dto.MemberAuth)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			components.GenerateBaseResponseWithValidationError(nil, false, components.ValidationError, err))
		return
	}
	fmt.Println(req)
	token, err := mc.service.Login(req)
	if err != nil {
		ctx.AbortWithStatusJSON(components.TranslateErrorToStatusCode(err), components.GenerateBaseResponseWithError(nil, false, components.ValidationError, err))
		return
	}

	ctx.JSON(http.StatusCreated, components.GenerateBaseResponse(token, true, components.Success))
}

// Update Member godoc
//
//	@Summary		Update a member
//	@Description	Update a member
//	@Tags			Members
//	@Accept			json
//	@produces		json
//	@Param			Request	body		dto.MemberUpdate		true	"member"
//	@Success		201		{object}	components.BaseHttpResponse	"Success"
//	@Failure		400		{object}	components.BaseHttpResponse	"Failed"
//	@Failure		409		{object}	components.BaseHttpResponse	"Failed"
//	@Router			/api/v1/members/update [post]
//	@Security		AuthBearer
func (mc *MemberController) Update(ctx *gin.Context) {
	req := new(dto.MemberUpdate)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			components.GenerateBaseResponseWithValidationError(nil, false, components.ValidationError, err))
		return
	}
	res, err := mc.service.Update(req)
	if err != nil {
		ctx.AbortWithStatusJSON(components.TranslateErrorToStatusCode(err), components.GenerateBaseResponseWithError(nil, false, components.ValidationError, err))
		return
	}

	ctx.JSON(http.StatusCreated, components.GenerateBaseResponse(res, true, components.Success))
}
