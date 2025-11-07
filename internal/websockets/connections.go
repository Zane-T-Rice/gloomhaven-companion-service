package websockets

import "github.com/gorilla/websocket"

// All the clients that have connected to a local instance of the websocket.
var Connections = map[string]*websocket.Conn{}
