package wstf

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

const MethodPush = "$PUSH"
const MethodBroadcast = "$BROADCAST"
const MethodNotification = "$NOTIFICATION"

// Get the registered(or newly created and registered) Pusher instance.
// The lifecycle of the original response should end here.
func (m *Response) GetRegisteredPusher() *Pusher {
	return m.Connection.GetRegisteredPusher(m)
}

// Remove the registered Pusher instance.
// The lifecycle of the original response should end here.
func (m *Response) RemoveRegisteredPusher() *Pusher {
	return m.Connection.RemoveRegisteredPusher(m)
}

func (m *Response) GetPusherId() string {
	// Use the request path as the pusher ID.
	return m.Request.Path
}

// Get a corresponding Pusher instance, and the lifecycle of the original response should end here.
func (m *Response) NewPusher() *Pusher {
	return &Pusher{m.Connection, m.ConnectionLocals, m.Request, m.Locals, m.Request.Path, MethodPush}
}

// The pusher used to proactively push messages to the corresponding devices.
type Pusher struct {
	// The reference of wstf Connection.
	Connection *Connection `json:"-"`
	// Shortcut for Connection.Locals
	ConnectionLocals map[string]interface{}
	// The reference of corresponding wstf Request.
	Request *Request
	// A map that contains response local variables scoped to the request.
	// This property is useful for exposing request-level information such as the request path name, authenticated user, user settings, and so on.
	Locals map[string]interface{}
	// The path used in messages to be pushed.
	Path string `json:"path"`
	// The method used in messages to be pushed, which may be '$PUSH', '$BROADCAST', '$NOTIFICATION', '$SMS', '$AD'.
	Method string `json:"method"`
}

// The message/response used for the pushing services.
type PushMessage struct {
	// Unique http request identifier.
	Id string `json:"id"`
	// The method used to push the message.
	Method string `json:"method"`
	// The path of the pushing response is used as an identifier.
	Path string `json:"path"`
	// Http status.
	Status int `json:"status"`
	// The http status text.
	//StatusText string `json:"statusText"`
	// Headers to be send to client.
	Headers map[string]string `json:"headers"`
	// Response body.
	Body interface{} `json:"body"`
}

func NewPusher(method, path string, status int, body interface{}) *PushMessage {
	return &PushMessage{path, method, path, status, nil, body}
}

// Push a message with m.Path and the status code http.StatusOK(200).
func (m *Pusher) Push(message interface{}) error {
	return m.PushAny(m.Path, http.StatusOK, message)
}

// Push a message with a custom path and the status code http.StatusOK(200).
func (m *Pusher) PushWithPath(path string, message interface{}) error {
	return m.PushAny(path, http.StatusOK, message)
}

// Push a message with a custom path and the status code http.StatusOK(200).
func (m *Pusher) PushWithCode(status int, message interface{}) error {
	return m.PushAny(m.Path, status, message)
}

// Push a message with a custom path and a custom status code.
func (m *Pusher) PushAny(path string, status int, body interface{}) error {
	app := m.Connection.Application
	message := NewPusher(m.Method, path, status, body)
	if app.OnEncodingPushMessage != nil {
		// CALLBACK for all push messages.
		app.OnEncodingPushMessage(m, message)
	}
	bytes, err := json.Marshal(message)
	if err != nil {
		if app.OnFailedToEncodePushMessage != nil {
			app.OnFailedToEncodePushMessage(err, m, message)
		} else {
			fmt.Println("[WSTF] Failed to encode the PUSH MESSAGE:", message)
		}
		return err
	}
	return m.Write(websocket.TextMessage, message, bytes)
}

// Write the push message to the client.
// The instance of the #PushMessage is used to call back.
func (m *Pusher) Write(mt int, message *PushMessage, bytes []byte) error {
	app := m.Connection.Application
	if app.OnPushingMessage != nil {
		// CALLBACK for all push messages.
		app.OnPushingMessage(m, message, bytes)
	}
	err := m.Connection.WebSocketConn.WriteMessage(mt, bytes)
	if err != nil {
		if app.OnFailedToPushMessage != nil {
			app.OnFailedToPushMessage(err, m, message, bytes)
		} else {
			fmt.Println("[WSTF] Failed to push the message to client:", err, bytes)
		}
	}
	return err
}
