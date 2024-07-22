package harvester

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{} // use default options

func ServerStart(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("Upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Read:", err)
			break
		}
		log.Printf("Recv: %s", message)
		err = c.WriteMessage(mt, message)
		log.Printf("Write:%d - %s", mt, message)
		if err != nil {
			log.Println("Write:", err)
			break
		}
	}
}