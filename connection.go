package wstf

import (
	"log"
	"github.com/gorilla/websocket"
	"fmt"
	"net/http"
)

type Connection struct {
}

func NewConnection(app Application, conn *websocket.Conn, request *http.Request) Connection {
	connection := Connection{}
	connectionLocals := map[string]interface{}{}
	res := NewResponse(conn, connectionLocals, request, "")
	req := Request{}
	app.OnConnectionRouter.Handle(request.URL.Path, req, res, func() {
		fmt.Println("A device is connected:", request.URL.Path)
	})
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
		res := NewResponse(conn, connectionLocals, request, req.ID)
		app.RootRouter.Handle(req.Path, *req, res, func() {
			fmt.Println("Unhandled request!")
			res.Error(404, "Unhandled request!")
		})
	}
	app.OnDisconnectionRouter.Handle(request.URL.Path, req, res, func() {
		fmt.Println("A device is connected:", request.URL.Path)
	})
	conn.Close()
	return connection
}
