package routers

import (
	handlers "drivers-service/api/controllers"
	"github.com/gin-gonic/gin"
)

func Health(r *gin.RouterGroup) {
	handler := handlers.NewHealthHandler()

	r.GET("/", handler.HandlerGet)
}
