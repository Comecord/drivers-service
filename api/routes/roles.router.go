package routers

import (
	"context"
	"crm-glonass/api/controllers"
	"crm-glonass/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Roles(r *gin.RouterGroup, db *mongo.Database) {
	cfg := config.GetConfig()
	ctx := context.Background()
	h := controllers.NewRoleController(db, ctx, cfg)

	r.POST("/create", h.Create)
	r.GET("/list", h.List)
}
