package wstf

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Defined commands.
const PrefixCmd = '$'
const CmdPingString = "$PING"
const CmdPongString = "$PONG"
const PrefixJsonObject = '{'

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
	httpRequest := m.HttpRequest
	// Fake request and response for connection handlers.
	mConnectionRequest := &Request{Connection: m, HttpRequest: m.HttpRequest}
	mConnectionResponse := NewResponse(m, m.Locals, mConnectionRequest)

	if app.OnConnectionRouter != nil {
		// CALLBACK for OnConnect.
		app.OnConnectionRouter.Handle(httpRequest.URL.Path, mConnectionRequest, mConnectionResponse, func() {
			fmt.Println("A device is connected:", httpRequest.URL.Path)
		})
	}
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			if app.OnReadMessageFailed != nil {
				// CALLBACK for read message failed.
				app.OnReadMessageFailed(m, err)
			}
			break
		}
		if app.OnReceiveMessage != nil {
			// CALLBACK for message.
			app.OnReceiveMessage(m, mt, message)
		}
		if mt != websocket.TextMessage || len(message) == 0 || message[0] != PrefixJsonObject {
			m.HandleMessage(mt, message)
			continue
		}
		req, err := NewRequest(message, m)
		if err != nil {
			if app.OnReceiveInvalidRequest != nil {
				// CALLBACK for invalid request.
				app.OnReceiveInvalidRequest(m, mt, message)
			}
			continue
		}
		res := NewResponse(m, m.Locals, req)
		if app.OnReceiveRequest != nil {
			// CALLBACK for request.
			app.OnReceiveRequest(req, res, message)
		}
		app.RootRouter.Handle(req.Path, req, res, func() {
			fmt.Println("Unhandled request!")
			res.Error(http.StatusNotFound, "Unhandled request!")
		})
	}
	if app.OnDisconnectionRouter != nil {
		// CALLBACK for onDisconnect.
		app.OnDisconnectionRouter.Handle(httpRequest.URL.Path, mConnectionRequest, mConnectionResponse, func() {
			fmt.Println("A device is connected:", httpRequest.URL.Path)
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
		// CALLBACK for unhandled message.
		m.Application.OnReceiveUnhandledMessage(m, messageType, message)
	}
}
