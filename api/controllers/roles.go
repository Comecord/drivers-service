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

type RolesController struct {
	service *services.RoleService
}

func NewRoleController(db *mongo.Database, ctx context.Context, conf *config.Config) *RolesController {
	service, ok := services.NewRoleService(db, conf, ctx, "roles").(*services.RoleService)
	if !ok {
		return nil
	}
	return &RolesController{
		service: service,
	}
}

// CreateRole Создание роли
//
//	@Summary		Создание роли
//	@Description	Создание роли
//	@Tags			Roles
//	@Accept			json
//	@produces		json
//	@Param			Request	body		dto.Role					true	"role"
//	@Success		201		{object}	components.BaseHttpResponse	"Success"
//	@Failure		400		{object}	components.BaseHttpResponse	"Failed"
//	@Failure		409		{object}	components.BaseHttpResponse	"Failed"
//	@Router			/api/v1/roles/create [post]
func (r *RolesController) CreateRole(ctx *gin.Context) {
	req := new(dto.Role)
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			components.GenerateBaseResponseWithValidationError(nil, false, components.ValidationError, err))
		return
	}
	role := r.service.CreateRole(req)

	ctx.JSON(http.StatusCreated, components.GenerateBaseResponse(role, true, components.Success))
}

// ListRoles Вывод всех ролей
//
//	@Summary		Вывод всех ролей
//	@Description	Вывод всех ролей
//	@Tags			Roles
//	@Accept			json
//	@produces		json
//	@Success		200	{array}		dto.RoleList				"Success"
//	@Failure		400	{object}	components.BaseHttpResponse	"Failed"
//	@Failure		409	{object}	components.BaseHttpResponse	"Failed"
//	@Router			/api/v1/roles/list [get]
func (r *RolesController) ListRoles(ctx *gin.Context) {
	roles, err := r.service.ListRoles()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest,
			components.GenerateBaseResponseWithValidationError(nil, false, components.NotFoundError, err))
	}
	ctx.JSON(http.StatusOK, roles)
}
