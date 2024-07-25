package routers

import (
	"context"
	"drivers-service/api/handlers"
	"drivers-service/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Roles(r *gin.RouterGroup, db *mongo.Database) {
	cfg := config.GetConfig()
	ctx := context.Background()
	h := handlers.NewRoleController(db, ctx, cfg)

	r.POST("/create", h.Create)
	r.GET("/list", h.List)
}
