package wstf

// Requests handler.
type Handler interface {
	// Handle a request.
	Handle(remainingPath string, req *Request, res *Response, next func())
}

// Call routes recursively in order unless the next is called.
func HandleRoutes(remainingPath string, routes []*Route, req *Request, res *Response, next func()) {
	if routes == nil || len(routes) == 0 {
		next()
		return
	}
	i := 0
	var goOn func()
	goOn = func() {
		if i < len(routes) {
			route := routes[i]
			i++
			route.Handle(remainingPath, req, res, goOn)
		} else {
			next()
		}
	}
	goOn()
}

// Call routers recursively in order unless the next is called.
func HandleRouters(remainingPath string, routers []*Router, req *Request, res *Response, next func()) {
	if routers == nil || len(routers) == 0 {
		next()
		return
	}
	i := 0
	var goOn func()
	goOn = func() {
		if i < len(routers) {
			router := routers[i]
			i++
			router.Handle(remainingPath, req, res, goOn)
		} else {
			next()
		}
	}
	goOn()
}

// Call processors recursively in order unless the next is called.
func HandleMethods(processors []func(req *Request, res *Response, next func()), req *Request, res *Response, next func()) {
	if processors == nil || len(processors) == 0 {
		next()
		return
	}
	i := 0
	var goOn func()
	goOn = func() {
		if i < len(processors) {
			processor := processors[i]
			i++
			processor(req, res, goOn)
		} else {
			next()
		}
	}
	goOn()
}