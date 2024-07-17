package routers

import (
	"context"
	"crm-glonass/api/controllers"
	"crm-glonass/config"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Members(r *gin.RouterGroup, db *mongo.Database) {
	cfg := config.GetConfig()
	ctx := context.Background()
	h := controllers.NewMemberController(db, ctx, cfg)

	r.POST("/register", h.Register)
	r.POST("/login", h.Login)
}
