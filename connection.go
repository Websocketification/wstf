package wstf

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Defined commands.
const CmdPrefix = '$'
const CmdPingString = "$PING"
const CmdPongString = "$PONG"
const JsonObjectPrefix = '{'

var CmdPingBytes = []byte(CmdPingString)
var CmdPongBytes = []byte(CmdPongString)

type Connection struct {
	Application *Application `json:"application"`
	// The original http request.
	HttpRequest *http.Request `json:"httpRequest"`
	// The websocket connection.
	WebSocketConn *websocket.Conn `json:"webSocketConn"`
	// The local variables exist in the scope/lifecycle of connection.
	Locals map[string]interface{} `json:"locals"`
}

func (m *Connection) OnConnect() {
	app := m.Application
	conn := m.WebSocketConn
	request := m.HttpRequest
	req := &Request{Connection: m, HttpRequest: m.HttpRequest}
	res := NewResponse(m, m.Locals, req)

	if app.OnConnectionRouter != nil {
		app.OnConnectionRouter.Handle(request.URL.Path, req, res, func() {
			fmt.Println("A device is connected:", request.URL.Path)
		})
	}
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			if app.OnReadMessageFailed != nil {
				app.OnReadMessageFailed(m, err)
			}
			break
		}
		if app.OnReceiveMessage != nil {
			app.OnReceiveMessage(m, mt, message)
		}
		if mt != websocket.TextMessage || len(message) == 0 || message[0] != JsonObjectPrefix {
			m.HandleMessage(mt, message)
			continue
		}
		req, err := NewRequest(message, m)
		if err != nil {
			if app.OnReceiveInvalidRequest != nil {
				app.OnReceiveInvalidRequest(m, mt, message)
			}
			if app.OnReceiveUnhandledMessage != nil {
				m.Application.OnReceiveUnhandledMessage(m, mt, message)
			}
			continue
		}
		res := NewResponse(m, m.Locals, req)
		app.RootRouter.Handle(req.Path, req, res, func() {
			fmt.Println("Unhandled request!")
			res.Error(http.StatusNotFound, "Unhandled request!")
		})
	}
	if app.OnDisconnectionRouter != nil {
		app.OnDisconnectionRouter.Handle(request.URL.Path, req, res, func() {
			fmt.Println("A device is connected:", request.URL.Path)
		})
	}
	conn.Close()
}

// Handle known and unknown messages.
func (m *Connection) HandleMessage(messageType int, message []byte) {
	if messageType == websocket.TextMessage {
		msg := string(message)
		if msg == CmdPingString {
			m.WebSocketConn.WriteMessage(websocket.TextMessage, CmdPongBytes)
			return
		}
		if msg == CmdPongString {
			return
		}
	}
	if m.Application.OnReceiveUnhandledMessage != nil {
		m.Application.OnReceiveUnhandledMessage(m, messageType, message)
	}
}
