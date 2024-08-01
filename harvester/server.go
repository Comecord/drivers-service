package harvester

import (
	"fmt"
)

type Message struct {
	UserID string `json:"userID"`
	Type   string `json:"type"`
}

// Функция для обработки сообщений "vehicles"
func VehicleListService(msg Message) map[string]string {
	data := GetVehicleList()
	fmt.Printf("VEHICLES: %v", data)
	vehicleData := fmt.Sprintf("%v", data)
	return map[string]string{"status": "success", "message": vehicleData, "userID": msg.UserID}
}

func LoginService(msg Message) map[string]string {
	data := Login()
	authData := fmt.Sprintf("AuthId: %v, UserId: %v", data.AuthId, data.UserId)
	return map[string]string{"status": "success", "message": authData, "userID": msg.UserID}
}

// Структура для маршрутов сервера
type ServerRoutes map[string]func(Message) map[string]string
