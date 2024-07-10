package routers

import (
	"context"
	"crm-glonass/api/controllers"
	"crm-glonass/config"
	"crm-glonass/pkg/logging"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

var log = logging.NewLogger(config.GetConfig())

func Vehicles(r *gin.RouterGroup, collection *mongo.Collection) {

	ctx := context.Background()
	h := controllers.NewVehiclesController(collection, ctx, log)
	//
	r.POST("/", h.Create)
	//r.PATCH("/:id", h.Update)
	//r.DELETE("/:id", h.Delete)
	//r.GET("/:id", h.GetById)
	//r.POST("/get-by-filter", h.GetByFilter)
}
