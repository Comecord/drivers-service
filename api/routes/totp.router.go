package routers

import (
	"context"
	"drivers-service/api/controllers"
	"drivers-service/config"
	"drivers-service/middlewares"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func AuthTotp(r *gin.RouterGroup, db *mongo.Database) {
	cfg := config.GetConfig()
	ctx := context.Background()
	h := controllers.NewTotpController(db, ctx, cfg)

	r.POST("/generate", h.GenerateAuthentication)
	r.GET("/active/:code", middlewares.Authentication(cfg), middlewares.Authorization([]string{"member"}), h.ActiveAuthentication)
}
