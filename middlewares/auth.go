package middlewares

import (
	"crm-glonass/api/components"
	"crm-glonass/api/services"
	"crm-glonass/config"
	"crm-glonass/constants"
	"crm-glonass/pkg/service_errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
)

func Authentication(cfg *config.Config) gin.HandlerFunc {
	var tokenService = services.NewTokenService(cfg)

	return func(c *gin.Context) {
		var err error
		claimMap := map[string]interface{}{}
		auth := c.GetHeader(constants.AuthorizationHeaderKey)
		token := strings.Split(auth, " ")
		if auth == "" {
			err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenRequired}
		} else {
			claimMap, err = tokenService.GetClaims(token[1])
			if err != nil {
				switch err.(*jwt.ValidationError).Errors {
				case jwt.ValidationErrorExpired:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
				default:
					err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
				}
			}
		}
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				components.GenerateBaseResponseWithError(nil, false, components.AuthError, err))
			return
		}

		c.Set(constants.UserIdKey, claimMap[constants.UserIdKey])
		c.Set(constants.FirstNameKey, claimMap[constants.FirstNameKey])
		c.Set(constants.LastNameKey, claimMap[constants.LastNameKey])
		c.Set(constants.UserNameKey, claimMap[constants.UserNameKey])
		c.Set(constants.EmailKey, claimMap[constants.EmailKey])
		c.Set(constants.MobileNumberKey, claimMap[constants.MobileNumberKey])
		c.Set(constants.RolesKey, claimMap[constants.RolesKey])
		c.Set(constants.ExpireTimeKey, claimMap[constants.ExpireTimeKey])

		c.Next()
	}
}

func Authorization(validRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(c.Keys) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden,
				components.GenerateBaseResponse(nil, false, components.ForbiddenError))
			return
		}
		rolesVal := c.Keys[constants.RolesKey]
		fmt.Println(rolesVal)
		if rolesVal == nil {
			c.AbortWithStatusJSON(http.StatusForbidden,
				components.GenerateBaseResponse(nil, false, components.ForbiddenError))
			return
		}
		roles := rolesVal.([]interface{})
		val := map[string]int{}
		for _, item := range roles {
			val[item.(string)] = 0
		}

		for _, item := range validRoles {
			if _, ok := val[item]; ok {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, components.GenerateBaseResponse(nil, false, components.ForbiddenError))
	}
}
