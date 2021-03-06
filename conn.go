package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
	"time"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	sendBinary chan []byte
}

func checkOrigin(r *http.Request) bool {

	actualOrigin := r.Header.Get("origin")
	clientIp := strings.Split(r.RemoteAddr, ":")[0]
	ipAllowed := true

	log.Print("Checking origin ", actualOrigin, "...")

	if len(config.Allow) > 0 {
		ipAllowed = false
		for _, ip := range config.Allow {
			if ip == clientIp {
				ipAllowed = true
				break
			}
		}
	}

	if !ipAllowed {
		log.Print("Ip ", clientIp, " is not allowed")
	}

	if len(config.Origins) == 0 && ipAllowed {
		return true
	}

	for _, origin := range config.Origins {
		if origin == actualOrigin && ipAllowed {
			return true
		}
	}

	return false
}

// readPump pumps messages from the websocket connection to the hub.
func (c *connection) readPump() {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		mtype, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		if mtype == 2 {

			packet, err := ReadPacket(message)

			if err != nil {
				log.Print("Fail to read binary packet")
				break
			}

			Execute(c, packet)
		}
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (c *connection) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case bytes, ok := <-c.sendBinary:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.BinaryMessage, bytes); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serverWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws, sendBinary: make(chan []byte, 256)}
	h.register <- c

	go c.writePump()
	c.readPump()
}
