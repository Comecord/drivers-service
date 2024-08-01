package main

import (
	"context"
	"drivers-service/config"
	"drivers-service/data/mongox"
	"drivers-service/harvester"
	"drivers-service/pkg/logging"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

var (
	conf   = config.GetConfig()
	logger = logging.NewLogger(config.GetConfig())
)

// Структура для маршрутов сервера
type ServerRoutes map[string]func(message harvester.Message) map[string]string

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешить все источники для простоты
	},
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	fmt.Println("Starting server on :8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	db, err := mongox.Connection(conf, ctx, logger)
	if err != nil {
		fmt.Println("Error getting MongoDB client:", err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	var resultEmail struct {
		Email string `bson:"email"`
	}

	email := r.Header.Get("Auth-Email")

	err = db.Collection("members").FindOne(ctx, bson.M{"email": email}).Decode(&resultEmail)
	if errors.Is(err, mongo.ErrNoDocuments) {
		conn.WriteMessage(websocket.TextMessage, []byte("Email не найден"))
		return
	} else if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Ошибка базы данных"))
		return
	}

	// Определение маршрутов
	routes := ServerRoutes{
		"vehicles": harvester.VehicleListService,
	}

	for {
		var msg harvester.Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading JSON:", err)
			break
		}

		// Обработка сообщения на основе маршрутов
		if handler, exists := routes[msg.Type]; exists {
			response := handler(msg)
			conn.WriteJSON(response)
		} else {
			conn.WriteJSON(map[string]string{"status": "error", "message": "Unknown type"})
		}
	}
}
