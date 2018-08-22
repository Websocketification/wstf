package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Websocketification/wstf"
	"github.com/gorilla/websocket"
)

// Set up main router
var mRouter = wstf.NewRouter()

func main() {
	app := wstf.NewApplication(mRouter)

	// The middleware use "*" to match everything rather than match the child routers.
	mRouter.Use(".*").All(func(req *wstf.Request, res *wstf.Response, next func()) {
		fmt.Println("Process cookies here!")
		next()
	}).All(func(req *wstf.Request, response *wstf.Response, next func()) {
		fmt.Println("Request database!")
		next()
	}).Get(func(req *wstf.Request, response *wstf.Response, next func()) {
		fmt.Println("All GET requests!")
		next()
	}).Post(func(req *wstf.Request, res *wstf.Response, next func()) {
		fmt.Println("Log post requests here!")
		next()
	})

	mRouter.Use("/").Get(func(req *wstf.Request, res *wstf.Response, next func()) {
		res.Done("Hello, this is a WebSocket server.")
	})

	// Set up sub router.
	mSubRouter := wstf.NewRouter()
	type User struct {
		ID   string
		Name string
	}
	mSubRouter.Use("/{userName}").Get(func(req *wstf.Request, res *wstf.Response, next func()) {
		res.Done(User{ID: "Fisher", Name: "Awesome Fisher"})
	})
	// All else requests.
	mSubRouter.Use(".*").All(func(req *wstf.Request, res *wstf.Response, next func()) {
		fmt.Println("The request is not processed!")
	})

	// Adding sub router uses empty string to match the child routers.
	mRouter.Push("/users", mSubRouter)

	// Redirect requests to wstf handler func.
	http.HandleFunc("/WebSocketServer", app.GetWebsocketHandlerFunc(
		&websocket.Upgrader{},
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Println("The http request is not upgradable:", r.URL.Path)
		},
		func(conn *wstf.Connection, w http.ResponseWriter) (*http.Header, bool) {
			// Allow all upgradable request to be upgraded.
			return nil, true
		},
		func(err error, w http.ResponseWriter, r *http.Request) {
			fmt.Println("Failed to upgrade to WebSocket:", r.URL.Path)
		},
	))

	fmt.Println("Server will be running at 127.0.0.1:3333")
	err := http.ListenAndServe("127.0.0.1:3333", nil)
	if err != nil {
		log.Fatal(err)
	}
}
