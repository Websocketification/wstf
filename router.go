package wstf

// Router
//
// TODO A method(`OnError()`) is needed to catch the 404, 500 error before send it to client.
type Router struct {
	Routes []*Route
	//ParamsRoute []Route
}

func NewRouter() *Router {
	m := &Router{}
	return m
}

// Push a router to the end of sequence.
func (m *Router) Push(pattern string, router *Router) {
	route := NewRoute(pattern, router)
	m.Routes = append(m.Routes, route)
}

// Add a func to specific pattern. Either "" or "*" will match every request.
//
// To keep the code clean and simple, here we do not use pattern slices `Use(pattern ...string)`.
func (m *Router) Use(pattern string) *Route {
	route := NewRoute(pattern, nil)
	m.Routes = append(m.Routes, route)
	return route
}

// Handle request.
func (m *Router) Handle(remainingPath string, req *Request, res *Response, next func()) {
	HandleRoutes(remainingPath, m.Routes, req, res, next)
}

