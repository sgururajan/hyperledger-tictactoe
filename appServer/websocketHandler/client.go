package websocketHandler

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/segmentio/ksuid"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
	"net/http"
	"time"
	"github.com/sgururajan/hyperledger-tictactoe/domainModel"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	hub        *Hub
	conn       *websocket.Conn
	send       chan []byte
	blockEvent chan domainModel.BlockInfo
	orgName    string
	beSubsFn   BlockEventSubscribe
	beUnSubsFn BlockEvnetUnSubscribe
}

type subscribeToBlockEventMsg struct {
	MsgType string `json:"type,omitempty"`
	Payload struct {
		OrgName string `json:"orgName"`
	} `json:"payload,omitempty"`
}

type BlockEventSubscribe func(string, chan domainModel.BlockInfo)
type BlockEvnetUnSubscribe func(string, chan domainModel.BlockInfo)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if r.Header.Get("origin") == "http://localhost:4200" {
			return true
		}
		return false
	},
}

func (c *Client) readPump() {
	log := utils.NewAppLogger("WebsocketReadPump", "")
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("error reading socket: %#v", err)
			}
			break
		}
		//nBytes:= bytes.IndexByte(msg,0)
		socketMessage := SocketMessage{}

		err = json.Unmarshal(msg, &socketMessage)
		log.Infof("mesage received from websocket: %#v", socketMessage)
		log.Infof("mesage received from websocket: %s", string(msg))
		//message:= bytes.TrimSpace(bytes.Replace(msg, newLine, space,-1))

		switch socketMessage.MsgType {
		case "subscribeToBlockEvent":
			subsMsg := subscribeToBlockEventMsg{}
			err = json.Unmarshal(msg, &subsMsg)
			log.Infof("refined value: %#v", subsMsg)
			c.orgName = subsMsg.Payload.OrgName
			c.beSubsFn(c.orgName, c.blockEvent)
			break

		}

	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				break
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				break
			}

			w.Write(message)
			w.Close()
			break
		case blockMsg, ok := <-c.blockEvent:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.beUnSubsFn(c.orgName, c.blockEvent)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			socketMessage := SocketMessage{
				MsgType: "block",
				Payload: blockMsg,
			}
			data, err := json.Marshal(socketMessage)
			if err != nil {
				break
			}
			w.Write(data)
			w.Close()
			break
		}
	}
}

func serveWebSocketConn(hub *Hub, blockEventSubscribeFn BlockEventSubscribe, blockEventUnsubscribeFn BlockEvnetUnSubscribe, w http.ResponseWriter, r *http.Request) {
	logger := utils.NewAppLogger("ServeWebsocketConn", "")

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Errorf("error while creating web socket connection. Err: %#v", err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), beSubsFn: blockEventSubscribeFn, beUnSubsFn: blockEventUnsubscribeFn, blockEvent: make(chan domainModel.BlockInfo)}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()

	clientId := ksuid.New().String()
	msgBytes, err := json.Marshal(SocketMessage{MsgType: "clientid", Payload: clientId})
	if err == nil {
		client.send <- msgBytes
	}
	//client.send <- []byte(clientId)
}
