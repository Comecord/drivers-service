package controllers

import (
	"crm-glonass/api/components"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthHandler struct {
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// HandlerGet godoc
//
//	@Summary		Health Check
//	@Description	Health Check
//	@Tags			Health Check
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	components.BaseHttpResponse	"Success"
//	@Failure		400	{object}	components.BaseHttpResponse	"Failed"
//	@Router			/api/v1/health/ [get]
func (h *HealthHandler) HandlerGet(c *gin.Context) {
	c.JSON(http.StatusOK, components.GenerateBaseResponse("ok", true, components.Success))
}
