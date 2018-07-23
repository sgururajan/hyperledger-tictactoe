package websocketHandler

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
)

type SocketMessage struct {
	MsgType string `json:"type,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type WebSocketHandler struct {
	hub *Hub
}

var logger = utils.NewAppLogger("wsHandler","")

func NewWebSocketHanler() *WebSocketHandler {
	return &WebSocketHandler{
		hub: newHub(),
	}
}

func (m *WebSocketHandler) Initialize() {
	go m.hub.run()
}

func (m *WebSocketHandler) SocketConnectionHandlerFunc(w http.ResponseWriter, r *http.Request) {
	logger.Infof("new socket connection received")
	serveWebSocketConn(m.hub, w, r)
}

func (m *WebSocketHandler) BroadcastMessage(msgType string, msg interface{}) {
	msgPayload:= SocketMessage{
		MsgType: msgType,
		Payload: msg,
	}

	msgJson,err:= json.Marshal(msgPayload)
	if err != nil {
		msgJson, err = getErrorMessage(err)
		if err==nil {
			m.hub.broadcast <- msgJson
		}
	}

	m.hub.broadcast <- msgJson
}

func getErrorMessage(err error) ([]byte, error) {
	msgPayload:= SocketMessage{
		MsgType: "error",
		Payload: []byte(fmt.Sprintf("$#v", err)),
	}

	result,err:= json.Marshal(msgPayload)
	return result, err
}
