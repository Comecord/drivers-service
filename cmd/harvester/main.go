package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

type Message struct {
	UserID string `json:"userID"`
	Type   string `json:"type"`
}

// Функция для обработки сообщений "vehicles"
func VehicleService(msg Message) map[string]string {
	// Обработка логики для "vehicles"
	return map[string]string{"status": "success", "message": "Vehicle data processed", "userID": msg.UserID}
}

// Функция для обработки сообщений "join"
func UserJoinService(msg Message) map[string]string {
	return map[string]string{"status": "success", "message": "User joined", "userID": msg.UserID}
}

// Структура для маршрутов сервера
type ServerRoutes map[string]func(Message) map[string]string

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
	//authHeader := r.Header.Get("Authorization")
	//if authHeader == "" {
	//	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	//	return
	//}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected:", conn.RemoteAddr())

	// Определение маршрутов
	routes := ServerRoutes{
		"vehicles": VehicleService,
		"join":     UserJoinService,
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
