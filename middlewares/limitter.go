package middlewares

import (
	"drivers-service/api/components"
	"github.com/didip/tollbooth"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LimitByRequest() gin.HandlerFunc {
	lmt := tollbooth.NewLimiter(1, nil)
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusTooManyRequests,
				components.GenerateBaseResponseWithError(nil, false, components.LimiterError, err))
			return
		} else {
			c.Next()
		}
	}
}
