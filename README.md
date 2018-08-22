# Websocketification in Golang

<!-- > 2:48 PM, 6/23/17. -->

Websocketification in Golang.

## Supports

- [x] Check the request headers before upgrading, @see [`upgrader-helper.go#AbleToUpgrade()`](https://godoc.org/github.com/Websocketification/wstf#AbleToUpgrade).
	- [x] Check the request `Origin` for [Cross-Origin Resource Sharing (CORS)](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS) before upgrading.
- [x] Support middleware to **accept** or **reject** the WebSocket connection request, @see [`application.go#GetWebsocketHandlerFunc()`](https://godoc.org/github.com/Websocketification/wstf#Application.GetWebsocketHandlerFunc).
- [x] Set cookies *only* if the request can be upgraded.

## Example of Usage

The example can be found [here](example/main.go).

```go
# ./example/main.go
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
```


## Route

> /user/hello

"/user/{userID}"

## Relations

```
Router          1:n        Routes
Route           1:1        Router
Route           1:n        HTTP_METHODS
HTTP_METHOD     1:n        Processors[func(req Request, res Response, next func())]
```


```
.
├── Router_1
│   ├── Route_1
│   │   └── Router_1a
│   │       ├── Route_1a
│   │       └── Route_1b
│   ├── Route_2
│   │   ├── All()
│   │   ├── Get()
│   │   └── Post()
│   └── Route_3
├── Router_2
└── Router_3
```

```
/users
    /users/{userName}
        /users/{userName}/profile
/weathers
    /weathers/{location}
```

A router has multiple connected routes and is responsible to pass all requests to routes one by one.

A route has the path pattern to match desired paths, but a router does not have.

- A route has only one pattern.
- A route can have one router.
- A route can have multiple processors if without router.

## Handler

Handler is responsible for handle the requests and decide whether to handle
the request or just pass it down.

The handler should be an `Handler(interface)` with a method
`Handle func(remainingPath, req Request, res Response)`
like the `Route` or `Router`.

```Golang
type Handler interface{
    Handle func(remainingPath, req Request, res Response, next func())
}
```

## Processor

Processor does nearly the same thing Handler does.
A processor directly process the request and response.
It usually get resources, process the request, or simply response.

```Golang
func(req Request, res Response, next func())
```

## Slices in methods like `Push(pattern string, router Router)`, `Get(processor func(...))`, ...

Here we use as less slices as we can.

It is recommended to push one by one.

**Handlers/Processors slice is not frequently used.**
It is recommended to use single handler/processor as you can.
If you need to use slices, please prove yourself at the Github/Issues.

**To Keep the code clean and simple.**
Using handlers/processors slice can make the codes complex and reduce the readability.

## References

- ExperssJs (a web framework for Node.js) API Docs: https://expressjs.com/en/4x/api.html
