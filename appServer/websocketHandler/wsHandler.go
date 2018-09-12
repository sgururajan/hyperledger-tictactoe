package websocketHandler

import (
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/sgururajan/hyperledger-tictactoe/utils"
	"github.com/sgururajan/hyperledger-tictactoe/appServer/networkHandlers"
	"github.com/sgururajan/hyperledger-tictactoe/fabnetwork"
	"github.com/sgururajan/hyperledger-tictactoe/domainModel"
)

type SocketMessage struct {
	MsgType string `json:"type,omitempty"`
	Payload interface{} `json:"payload,omitempty"`
}

type WebSocketHandler struct {
	hub *Hub
	networkHandlers *networkHandlers.NetworkHandler
}

var logger = utils.NewAppLogger("wsHandler","")

func NewWebSocketHanler(handler *networkHandlers.NetworkHandler) *WebSocketHandler {
	return &WebSocketHandler{
		hub: newHub(),
		networkHandlers:handler,
	}
}

func (m *WebSocketHandler) Initialize() {
	go m.hub.run()
}

func (m *WebSocketHandler) SubscribeToBlockEvent(orgName string, receiver chan domainModel.BlockInfo) {
	log:= utils.NewAppLogger("SubscriveToBlockEvent","")
	network, err:= m.networkHandlers.GetNetwork("testnetwork")
	if err != nil {
		log.Errorf("error while getting network. err: %#v", err)
		return
	}

	blockEventListener:= fabnetwork.BlockEventListener{
		Receiver: receiver,
	}

	network.RegisterBlockEventListener("tictactoechannel", orgName, blockEventListener)
}

func (m *WebSocketHandler) UnsubscriveBlockEvent(orgName string, receiver chan domainModel.BlockInfo) {
	log:= utils.NewAppLogger("SubscriveToBlockEvent","")
	network, err:= m.networkHandlers.GetNetwork("testnetwork")
	if err != nil {
		log.Errorf("error while getting network. err: %#v", err)
		return
	}

	blockEventListener:= fabnetwork.BlockEventListener{
		Receiver: receiver,
	}

	network.UnRegisterBlockEventListener("tictactoechannel", orgName, blockEventListener)
}

func (m *WebSocketHandler) SocketConnectionHandlerFunc(w http.ResponseWriter, r *http.Request) {
	logger.Infof("new socket connection received")
	serveWebSocketConn(m.hub, m.SubscribeToBlockEvent, m.UnsubscriveBlockEvent, w, r)
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
