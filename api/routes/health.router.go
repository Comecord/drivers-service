package routers

import (
	handlers "crm-glonass/api/controllers"
	"github.com/gin-gonic/gin"
)

func Health(r *gin.RouterGroup) {
	handler := handlers.NewHealthHandler()

	r.GET("/", handler.HandlerGet)
}
