package client

import (
	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"
)

var addr = flag.String("addr", "localhost:5900", "http service address")

func ConnectWebsocket() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("Connecting to %s", u.String())

	// Пытаемся подключиться к WebSocket серверу
	for {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			log.Printf("Не удалось подключиться к WebSocket серверу: %v\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		done := make(chan struct{})

		go func() {
			defer close(done)
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
					return
				}
				log.Printf("recv: %s", message)
			}
		}()

		// Периодически отправляем сообщение на сервер
		timer := time.NewTimer(10 * time.Second)

		for {
			select {
			case <-done:
				return
			case <-timer.C:
				err := c.WriteMessage(websocket.TextMessage, []byte("ping"))
				if err != nil {
					log.Println("write:", err)
					c.Close()
					break
				}
				timer.Reset(10 * time.Second) // Сброс таймера для следующего периода
			case <-interrupt:
				log.Println("interrupt")

				// Cleanly close the connection by sending a close message and then
				// waiting (with timeout) for the server to close the connection.
				err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				if err != nil {
					log.Println("write close:", err)
					return
				}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
				return
			}
		}
	}
}
