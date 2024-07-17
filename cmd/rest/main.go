package main

import (
	"context"
	"crm-glonass/api"
	"crm-glonass/api/client"
	"crm-glonass/config"
	"crm-glonass/data/cache"
	mongox "crm-glonass/data/mongox"
	"crm-glonass/data/seeds"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	var wg sync.WaitGroup

	// Websocket client connection
	wg.Add(1)
	go func() {
		defer wg.Done()
		client.ConnectWebsocket()
		logger.Infof("WebSocket client connected")
	}()

	// Database connection
	wg.Add(1)
	go func() {
		defer wg.Done()
		database, _ := mongox.Connection(conf, ctx, logger)
		seeds.SeedRoles(database, ctx)
		err := cache.InitRedis(conf, ctx)
		if err != nil {
			logger.Error(logging.Redis, logging.Connection, err.Error(), map[logging.ExtraKey]interface{}{"Version": conf.Version})
		}
		logger.Infof("Listening on Swagger http://localhost:%d/swagger/index.html", conf.Server.IPort)
		api.InitialServer(conf, database, logger)
	}()

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Infof("Shutting down server...")

	// Wait for all goroutines to finish
	wg.Wait()
}
