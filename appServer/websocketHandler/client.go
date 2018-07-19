package websocketHandler

import (
	"github.com/gorilla/websocket"
	"net/http"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
	"time"
	"github.com/segmentio/ksuid"
	"encoding/json"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan []byte
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if r.Header.Get("origin")=="http://localhost:4200" {
			return true
		}
		return false
	},
}

func (c *Client) writePump() {
	ticker:= time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok:= <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err:= c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)
			if err:= w.Close(); err!=nil {
				return
			}
		}
	}
}

func serveWebSocketConn(hub *Hub, w http.ResponseWriter, r *http.Request) {
	logger:= utils.NewAppLogger("ServeWebsocketConn","")

	conn, err:= upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("error while creating web socket connection. Err: %#v", err)
		return
	}

	client:= &Client{hub:hub, conn:conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	go client.writePump()

	clientId:= ksuid.New().String()
	msgBytes,err:= json.Marshal(SocketMessage{MsgType:"clientid", Payload:clientId})
	if err ==nil {
		client.send <- msgBytes
	}
	//client.send <- []byte(clientId)
}
