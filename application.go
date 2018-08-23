package wstf

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type Application struct {
	// Main router.
	RootRouter *Router
	// Router that will be handled on device is connected.
	OnConnectionRouter *Router
	// Router that will be handled on device is disconnected.
	OnDisconnectionRouter *Router

	// On read message from client failed, which literally means connection is closed.
	OnReadMessageFailed func(conn *Connection, err error)
	// On received message from client.
	OnReceiveMessage func(conn *Connection, messageType int, message []byte)
	// Call when received message is not recognized(and hence unhandled) as requests or other commands.
	OnReceiveUnhandledMessage func(conn *Connection, messageType int, message []byte)
	// Call when received message starts with '{', but is not a valid request.
	OnReceiveInvalidRequest func(conn *Connection, messageType int, message []byte)
	// On received a valid request from client.
	OnReceiveRequest func(req *Request, res *Response, message []byte)

	// On the response is going to be encoded to JSON string.
	OnEncodingResponse func(req *Request, res *Response)
	// On failed to encode the JsonResponse.
	OnFailedToEncodeResponse func(err error, req *Request, res *Response)
	// On the response is going to be sent to client.
	OnSendingResponse func(req *Request, res *Response, message []byte)
	// On failed to send the response to client.
	OnFailedToSendResponse func(err error, req *Request, res *Response, message []byte)
}

// Create a new Application with root router.
// Use the RootRouter as param in constructor because it is required.
func NewApplication(rootRouter *Router) *Application {
	return &Application{RootRouter: rootRouter}
}

// Get the WebSocket handler for http routing.
//
// The #onNotAbleToUpgradeToWebSocket() will be called if the request not be able to upgrade( by checking request
// metadata, like headers, @see upgrader#AbleToUpgrade()).
//
// The #getWebSocketResponseHeader() will be called if the request can be upgraded. It is expected to
// 1. return `*, false` to *reject* WebSocket connection request.
//     - You should handle http.ResponseWrite yourself, like response a status code with some message.
// 2. return `*http.Header, true` to *accept* the WebSocket connection request.
//     - In this case, you should ignore the http.ResponseWrite.
//     - If the #responseHeader is not nil, it will be included in the response to the client's upgrade request and may
//        be used to specify cookies (Set-Cookie) and the application negotiated subprotocol (Sec-WebSocket-Protocol),
//        you may checkout github.com/gorilla/websocket/server.go#Upgrade() for more info.
//
// The #onUpgradeToWebSocketFailed() will be called if the request can be upgraded to WebSocket but failed to. The callback
// may be rarely called.
func (m *Application) GetWebsocketHandlerFunc(
	upgrader *websocket.Upgrader,
	onNotAbleToUpgradeToWebSocket func(w http.ResponseWriter, r *http.Request),
	getWebSocketResponseHeader func(conn *Connection, w http.ResponseWriter) (*http.Header, bool),
	onUpgradingToWebSocketFailed func(err error, conn *Connection, w http.ResponseWriter),
) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if !AbleToUpgrade(upgrader, w, r) {
			onNotAbleToUpgradeToWebSocket(w, r)
			return
		}
		conn := m.NewConnection(nil, r)
		responseHeader, accept := getWebSocketResponseHeader(conn, w)
		if !accept {
			return
		}
		connection, err := upgrader.Upgrade(w, r, *responseHeader)
		if err != nil {
			onUpgradingToWebSocketFailed(err, conn, w)
			return
		}
		conn.WebSocketConn = connection
		conn.OnConnect()
	}
}

// The RootRouter is required and hence we don't provide function to set or modify root router.
// An application can have only one router as the home router.
//func (m *Application) SetRootRouter(router *Router) {
//	m.RootRouter = router
//}

func (m *Application) OnConnected(router *Router) {
	m.OnConnectionRouter = router
}

func (m *Application) OnDisconnected(router *Router) {
	m.OnDisconnectionRouter = router
}

// Create a connection instance once a WebSocket connection is established and connected.
func (m *Application) NewConnection(conn *websocket.Conn, req *http.Request) *Connection {
	return &Connection{
		m,
		req,
		conn,
		make(map[string]interface{}),
	}
}
