package wstf

import (
	"net/http"
	"log"
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

// Create a new Application.
func NewApplication() *Application {
	return &Application{}
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

// Set root router.
// An application can have only one router as the home router.
func (m *Application) SetRootRouter(router *Router) {
	m.RootRouter = router
}

func (m *Application) OnConnected(router *Router) {
	m.OnConnectionRouter = router
}

func (m *Application) OnDisconnected(router *Router) {
	m.OnDisconnectionRouter = router
}
