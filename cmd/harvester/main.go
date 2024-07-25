package main

import (
	"context"
	"crm-glonass/config"
	"crm-glonass/data/mongox"
	"crm-glonass/harvester"
	"crm-glonass/pkg/logging"
	"fmt"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sync"
)

type Message struct {
	UserID string `json:"userID"`
	Type   string `json:"type"`
}

var (
	conf    = config.GetConfig()
	logger  = logging.NewLogger(config.GetConfig())
	clients = make(map[*Client]bool) // Порог блокировки
)
var mu sync.RWMutex

// Функция для обработки сообщений "vehicles"
func VehicleService(msg Message) map[string]string {
	// Обработка логики для "vehicles"
	return map[string]string{"status": "success", "message": "Vehicle data processed", "userID": msg.UserID}
}

// Функция для обработки сообщений "join"
func UserJoinService(msg Message) map[string]string {
	return map[string]string{"status": "success", "message": "User joined", "userID": msg.UserID}
}

func LoginService(msg Message) map[string]string {
	data := harvester.Login()
	authData := fmt.Sprintf("AuthId: %v, UserId: %v", data.AuthId, data.UserId)
	return map[string]string{"status": "success", "message": authData, "userID": msg.UserID}
}

// Структура для маршрутов сервера
type ServerRoutes map[string]func(Message) map[string]string

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешить все источники для простоты
	},
}

type Client struct {
	Conn *websocket.Conn
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
	if err == mongo.ErrNoDocuments {
		conn.WriteMessage(websocket.TextMessage, []byte("Email не найден"))
		return
	} else if err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Ошибка базы данных"))
		return
	}

	// Определение маршрутов
	routes := ServerRoutes{
		"vehicles": VehicleService,
		"join":     UserJoinService,
		"login":    LoginService,
	}

	for {
		var msg Message
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
