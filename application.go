package wstf

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	// Enable cross domain.
	return true
}} // use default options

type Application struct {
	// Main router.
	RootRouter *Router
	// Router that will be handled on device is connected.
	OnConnectionRouter *Router
	// Router that will be handled on device is disconnected.
	OnDisconnectionRouter *Router
}

// Create a new Application with root router.
// Use the RootRouter as param in constructor because it is required.
func NewApplication(rootRouter *Router) *Application {
	return &Application{RootRouter: rootRouter}
}

// Get the handler func for websocket.
func (m *Application) GetWebsocketHandlerFunc() func(w http.ResponseWriter, r *http.Request) {
	mHandler := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("Upgrade Error:", err)
			return
		}
		NewConnection(m, c, r)
	}
	return mHandler
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
