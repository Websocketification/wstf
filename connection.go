package wstf

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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
	res := NewResponse(conn, m.Locals, request, "")
	req := &Request{}

	if app.OnConnectionRouter != nil {
		app.OnConnectionRouter.Handle(request.URL.Path, req, res, func() {
			fmt.Println("A device is connected:", request.URL.Path)
		})
	}
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read: ", err, mt)
			break
		}
		fmt.Println("Received: ", mt, string(message))
		req, err := NewRequest(message)
		if err != nil {
			log.Fatal("Failed to parse request json string.", err)
		}
		res := NewResponse(conn, m.Locals, request, req.Id)
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
