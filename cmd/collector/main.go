package main

import (
	"drivers-service/config"
	"drivers-service/pkg/logging"
)

func main() {

	// Start server
	logger := logging.NewLogger(config.GetConfig())
	logger.Info(logging.General, logging.StartUp, "Started server...", map[logging.ExtraKey]interface{}{"Version": config.GetConfig().Version})

}
