package api

import (
	routers "crm-glonass/api/routes"
	"crm-glonass/api/validations"
	"crm-glonass/config"
	"crm-glonass/docs"
	_ "crm-glonass/docs"
	"crm-glonass/middlewares"
	"crm-glonass/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/mongo"
)

var logger = logging.NewLogger(config.GetConfig())

func InitialServer(cfg *config.Config, database *mongo.Database, logger logging.Logger) {
	gin.SetMode(cfg.Server.RunMode)
	r := gin.New()

	RegisterValidator()

	r.Use(middlewares.DefaultStructuredLogger(cfg))
	r.Use(middlewares.Cors(cfg))
	r.Use(gin.Logger(), gin.CustomRecovery(middlewares.ErrorHandler), middlewares.LimitByRequest())

	RegisterRouter(r, cfg, database)
	RegisterSwagger(r, cfg)

	logger.Info(logging.API, logging.Connection, fmt.Sprintf("Starting server on port:%d", cfg.Server.IPort), nil)
	err := r.Run(fmt.Sprintf(":%d", cfg.Server.IPort))
	if err != nil {
		logger.Fatal(logging.API, logging.StartUp, err.Error(), nil)
	}
}
func RegisterValidator() {
	val, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {

		err := val.RegisterValidation("password", validations.PasswordValidator, true)
		if err != nil {
			logger.Error(logging.Validation, logging.StartUp, err.Error(), nil)
		}
	}
}

func RegisterSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = "COMECORD"
	docs.SwaggerInfo.Description = "Система управление и мониторинга транспортных средств с системой GLONASS"
	docs.SwaggerInfo.Version = cfg.Version
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%d", cfg.Server.EPort)
	docs.SwaggerInfo.Schemes = []string{"http"}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and then your token.

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func RegisterRouter(r *gin.Engine, conf *config.Config, db *mongo.Database) {
	api := r.Group("/api")
	v1 := api.Group("/v1")
	{

		membersRouterGroup := v1.Group("/members")
		routers.Members(membersRouterGroup, db)

		authTotp := membersRouterGroup.Group("/totp")
		routers.AuthTotp(authTotp, db)

		vehicles := v1.Group("/vehicles", middlewares.Authentication(conf), middlewares.Authorization([]string{"member"}))
		routers.Vehicles(vehicles, db)

		health := v1.Group("/health")
		routers.Health(health)

		roles := v1.Group("/roles")
		routers.Roles(roles, db)
	}

	r.Static("/uploads", "./uploads")

}
