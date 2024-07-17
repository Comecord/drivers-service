package controllers

import (
	"context"
	"crm-glonass/api/components"
	"crm-glonass/api/dto"
	"crm-glonass/api/services"
	"crm-glonass/config"
	"crm-glonass/constants"
	"crm-glonass/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"strings"
)

type TotpController struct {
	service *services.TotpService
	token   *services.TokenService
}

func NewTotpController(db *mongo.Database, ctx context.Context, conf *config.Config) *TotpController {
	service, ok := services.NewTotpService(db, conf, ctx).(*services.TotpService)
	if !ok {
		return nil
	}
	return &TotpController{
		service: service,
		token: &services.TokenService{
			Logger: logging.NewLogger(conf),
			Cfg:    conf,
		},
	}
}

// GenerateAuthentication Auth godoc
//
//	@Summary		Generate TOTP
//	@Description	Generate TOTP authentication for member
//	@Tags			Auth
//	@Accept			json
//	@produces		json
//	@Param			Request	body		dto.TotpRequest				true	"payload"
//	@Success		201		{object}	components.BaseHttpResponse	"Success"
//	@Failure		400		{object}	components.BaseHttpResponse	"Failed"
//	@Failure		409		{object}	components.BaseHttpResponse	"Failed"
//	@Router			/api/v1/members/totp/generate [post]
func (to *TotpController) GenerateAuthentication(ctx *gin.Context) {
	var request dto.TotpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	response, err := to.service.GenerateTotp(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(components.TranslateErrorToStatusCode(err),
			components.GenerateBaseResponseWithError(nil, false, components.NotFoundError, err))
		return
	}
	ctx.JSON(http.StatusOK, components.GenerateBaseResponse(response, true, components.Success))
}

// ActiveAuthentication Auth godoc
//
//	@Summary		Active TOTP authentication
//	@Description	Active TOTP authentication for member
//	@Tags			Auth
//	@Accept			json
//	@produces		json
//	@Param			code	path		string							true	"Code"
//	@Success		201		{object}	components.BaseHttpResponse	"Success"
//	@Failure		400		{object}	components.BaseHttpResponse	"Failed"
//	@Failure		409		{object}	components.BaseHttpResponse	"Failed"
//	@Router			/api/v1/members/totp/active/{code} [get]
//	@Security		AuthBearer
func (toc *TotpController) ActiveAuthentication(ctx *gin.Context) {
	var request dto.TotpCodeVerify
	request.Code = ctx.Param("code")
	fmt.Printf("RESPONSE -------- %+v\n", request)
	var err error
	claimMap := map[string]interface{}{}
	auth := ctx.GetHeader(constants.AuthorizationHeaderKey)
	token := strings.Split(auth, " ")
	claimMap, err = toc.token.GetClaims(token[0])
	fmt.Printf("CLAIM -------- %+v\n", token)
	request.Email = claimMap["Email"].(string)
	fmt.Printf("RESPONSE -------- %+v\n", request)
	response, err := toc.service.Active(&request)
	fmt.Printf("RESPONSE -------- %+v\n", response)
	if err != nil {
		ctx.AbortWithStatusJSON(components.TranslateErrorToStatusCode(err),
			components.GenerateBaseResponseWithError(nil, false, components.NotFoundError, err))
		return
	}
	ctx.JSON(http.StatusOK, components.GenerateBaseResponse(response, true, components.Success))
}
