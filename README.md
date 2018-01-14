# Websocketification in Golang

<!-- > 2:48 PM, 6/23/17. -->

Websocketification in Golang.

## Examples

```go
package main

import (
	"fmt"
	"github.com/Websocketification/wstf"
	"net/http"
	"log"
)

func main() {
	app := wstf.NewApplication()

	// Set up main router
	mRouter := wstf.NewRouter()
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

	// Set up sub router.
	mSubRouter := wstf.NewRouter()
	type User struct {
		ID   string
		Name string
	}
	mSubRouter.Use("/{userName}").Get(func(req *wstf.Request, res *wstf.Response, next func()) {
		res.Done(User{ID: "BeFisher", Name: "Berton Fisher"})
	})
	// All else requests.
	mSubRouter.Use(".*").All(func(req *wstf.Request, res *wstf.Response, next func()) {
		fmt.Println("The request is not processed!")
	})
	// Adding sub router uses empty string to match the child routers.
	mRouter.Push("/users", mSubRouter)

	// Redirect requests to wstf handler func.
	http.HandleFunc("/WebsocketServer", app.GetWebsocketHandlerFunc())
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
