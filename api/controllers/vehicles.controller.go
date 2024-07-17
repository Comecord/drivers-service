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

type VehiclesController struct {
	service *services.VehicleService
}

func NewVehiclesController(db *mongo.Database, ctx context.Context, conf *config.Config) *VehiclesController {
	service, ok := services.NewVehicleService(db, conf, ctx, "vehicles").(*services.VehicleService)
	if !ok {
		// handle the error case where the type assertion fails
		return nil
	}
	return &VehiclesController{
		service: service,
	}
}

// Create Vehicle godoc
//
//		@Summary		Create a vehicle
//		@Description	Create a vehicle
//		@Tags			Vehicles
//		@Accept			json
//		@produces		json
//		@Param			Request	body		dto.CreateVehicleRequest								true	"Create a vehicle model"
//		@Success		201		{object}	components.BaseHttpResponse{result=dto.DBVehicleDTO}	"Created response"
//		@Failure		400		{object}	components.BaseHttpResponse								"Bad request"
//		@Router			/api/v1/vehicles/ [post]
//	 @Security       BearerAuth
func (vc *VehiclesController) Create(ctx *gin.Context) {
	var vehicle *dto.CreateVehicleRequest

	if err := ctx.ShouldBindJSON(&vehicle); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newPost, err := vc.service.Create(vehicle)

	if err != nil {
		ctx.AbortWithStatusJSON(components.TranslateErrorToStatusCode(err),
			components.GenerateBaseResponseWithError(nil, false, components.InternalError, err))
		return
	}
	ctx.JSON(http.StatusCreated, components.GenerateBaseResponse(newPost, true, components.Success))

}

func (vc *VehiclesController) Update(ctx *gin.Context) {

}

func (vc *VehiclesController) Delete(ctx *gin.Context) {

}

func (vc *VehiclesController) GetById(ctx *gin.Context) {

}

func (vc *VehiclesController) GetByFilter(ctx *gin.Context) {

}
