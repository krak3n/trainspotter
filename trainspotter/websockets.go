// Websocket Server

package main

import (
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

//
// Message Hub
//

type Hub struct {
	// Stores all connected clients
	connections map[*connection]bool

	// Messages to broadcast to the connected clients
	Broadcast chan []byte

	// Register a new client
	register chan *connection

	// Unregister clients, removing them from the pool
	unregister chan *connection
}

// Listens on the channels for messages and performs the
// relivant functionality
func (h *Hub) Run() {
	for {
		select {
		// On register messages, store the connection
		case c := <-h.register:
			h.connections[c] = true
			d, _ := json.Marshal(Berths)
			c.send <- d
		// On unregister messsages, delete the connection from the pool
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		// On incoming messages, loop over connected clients and send
		// the message
		case m := <-h.Broadcast:
			for c := range h.connections {
				select {
				// Put the message on the connections send channel
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

// Create a new message hub
func NewHub() Hub {
	return Hub{
		Broadcast:   make(chan []byte),
		register:    make(chan *connection),
		unregister:  make(chan *connection),
		connections: make(map[*connection]bool),
	}
}

//
// Writer
//

// Represents a WS connecton to Perceptor
type connection struct {
	// The actual WS connection
	ws *websocket.Conn

	// A channel to send messages to the connection
	send chan []byte
}

// Returns the remote address of the WS Connection
func (c *connection) addr() net.Addr {
	return c.ws.RemoteAddr()
}

// Writes a given message type and payload to the WS connection
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// Writes messages to the WS Connection
func (c *connection) writer() {
	// Create a ticker that will ping the client
	ticker := time.NewTicker(pingPeriod)

	// Ensure we stop the ticker and close the connection when we exit
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		// Get a message from the connections send channel
		case m, ok := <-c.send:
			// Not OK, send a close message
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return // break out of the loop
			}
			// Attempt to write the message to the connection, catching errors
			if err := c.write(websocket.TextMessage, m); err != nil {
				return // break out of the loop
			}
		case <-ticker.C:
			// Ping the client to keep the connection open
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return // break out of the loop
			}
		}
	}
}

//
// Service
//

// Upgrade instance to upgrade the connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSService struct {
	Hub Hub
}

// Connection handler, upgrades the connection and registers the
// connection with the hub
func (s WSService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only support GET requests
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	// Upgrade the connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// Create Connection instance
	c := &connection{send: make(chan []byte, 256), ws: ws}
	// Register the connection
	s.Hub.register <- c
	// Start the writer for the conneciton
	go c.writer()
}

// Create a new WS Service
func NewWSService(h Hub) WSService {
	return WSService{
		Hub: h,
	}
}
