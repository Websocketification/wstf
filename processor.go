package wstf

// Alias type for function handler, and here we name it as Processor, and Processors.
// A processor, works as a *handler function* or a *middleware*, directly processes the request, response, and next.
// Usually a processor gets some info from request, fetch some expected resources, then:
// 		- response the expected data to client as a *handler function*.
// 		- pass it to downstream processors/handlers as a *middleware*.
type Processor = func(req *Request, res *Response, next func())
type Processors = []Processor
