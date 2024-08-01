package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type Message struct {
	UserID string `json:"userID"`
	Type   string `json:"type"`
}

func ConnectSocket() {
	url := "ws://localhost:8000/ws"

	header := http.Header{}
	header.Add("Authorization", "Bearer your_auth_token") // Здесь добавьте ваш авторизационный токен

	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
	}
	defer conn.Close()

	fmt.Println("Connected to server:", url)

	// Отправка сообщения "messageData"
	msg := Message{
		UserID: "29318209382",
		Type:   "vehicles",
	}
	err = conn.WriteJSON(msg)
	if err != nil {
		log.Println("Error sending JSON:", err)
		return
	}

	// Ожидание ответа от сервера
	go func() {
		for {
			response := make(map[string]string)
			err = conn.ReadJSON(&response)
			if err != nil {
				log.Println("Error reading JSON:", err)
				break
			}
			if response["status"] == "success" {
				fmt.Printf("[responseData] %s\n", response["message"])
			} else {
				fmt.Printf("[errorsMessage] %s\n", response["message"])
			}
		}
	}()

	// Отправка сообщения "userJoin"
	joinMsg := Message{
		UserID: "29318209382",
		Type:   "join",
	}
	err = conn.WriteJSON(joinMsg)
	if err != nil {
		log.Println("Error sending JSON:", err)
		return
	}

	// Задержка для ожидания ответов от сервера
	time.Sleep(5 * time.Second)
}
