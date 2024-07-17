package main

import (
	"context"
	"crm-glonass/api"
	"crm-glonass/config"
	"crm-glonass/data/cache"
	mongox "crm-glonass/data/mongox"
	"os"

	"crm-glonass/pkg/logging"
)

var logger = logging.NewLogger(config.GetConfig())

func main() {
	conf := config.GetConfig()
	ctx := context.TODO()

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		err := os.Mkdir("uploads", os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	// Logger info
	logger.Info(logging.General, logging.StartUp, "Started server...", map[logging.ExtraKey]interface{}{"Version": conf.Version})

	// Database connection
	database, _ := mongox.Connection(conf, ctx, logger)
	cache.InitRedis(conf, ctx)
	logger.Infof("Listening on Swagger http://localhost:%d/swagger/index.html", conf.Server.IPort)
	api.InitialServer(conf, database, logger)

}
