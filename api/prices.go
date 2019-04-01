package api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func SubscribePrice(productsIDs []string, c *websocket.Conn) bool {
	payload := struct {
		Type       string   `json:"type"`
		ProductIds []string `json:"product_ids"`
		Channels   []string `json:"channels"`
		Key        string   `json:"key"`
		Signature  string   `json:"signature"`
		Timestamp  string   `json:"timestamp"`
		Passphrase string   `json:"passphrase"`
	}{
		Type:       "subscribe",
		ProductIds: productsIDs,
		Channels: []string{
			"level2",
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if err := c.WriteMessage(websocket.TextMessage, data); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func InitWebsocketConnection(messages chan []byte) (chan struct{}, *websocket.Conn) {
	c, _, err := websocket.DefaultDialer.Dial(os.Getenv("API_WEBSOCKET_URL"), nil)
	if err != nil {
		log.Fatal("dial:", err)
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
			messages <- message
		}
	}()
	return done, c
}
