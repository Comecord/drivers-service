package routers

import (
	"context"
	"drivers-service/api/handlers"
	"drivers-service/config"
	"drivers-service/middlewares"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Members(r *gin.RouterGroup, db *mongo.Database) {
	cfg := config.GetConfig()
	ctx := context.Background()
	h := handlers.NewMemberController(db, ctx, cfg)

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
	r.POST("/update", middlewares.Authentication(cfg), middlewares.Authorization([]string{"member", "admin"}), h.Update)
}
