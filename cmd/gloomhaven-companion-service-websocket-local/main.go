package main

import (
	"context"
	"fmt"
	"gloomhaven-companion-service/internal/websockets"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Upgrader is used to upgrade HTTP connections to WebSocket connections.
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket connection
	connectionUUID := uuid.New().String()
	connectResponse, err := websockets.Connect(context.TODO(), websockets.ConnectRequest{
		Headers: map[string]string{
			"Sec-WebSocket-Protocol": r.Header.Get("Sec-WebSocket-Protocol"),
		},
		QueryStringParameters: map[string]string{
			"campaignId": r.URL.Query().Get("campaignId"),
			"scenarioId": r.URL.Query().Get("scenarioId"),
		},
		RequestContext: websockets.RequestContext{
			ConnectionID: connectionUUID,
		},
	})

	if err != nil {
		log.Printf("Connection failed. %v", err)
	}

	header := http.Header{}
	header.Add("Sec-WebSocket-Protocol", connectResponse.Headers["Sec-WebSocket-Protocol"])
	conn, err := upgrader.Upgrade(w, r, header)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	websockets.Connections[connectionUUID] = conn
	defer conn.Close()
	// Listen for incoming messages
	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s\\n", message)
		websockets.Default(context.TODO(), websockets.DefaultRequest{
			Body: string(message),
			RequestContext: websockets.RequestContext{
				ConnectionID: connectionUUID,
			},
		})
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	fmt.Println("WebSocket server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
