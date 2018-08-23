package wstf

// The JSON response used to response to corresponding wstf request.
type JsonResponse struct {
	// Unique http request identifier.
	Id string `json:"id"`
	// Http status.
	Status int `json:"status"`
	// The http status text.
	//StatusText string `json:"statusText"`
	// Headers to be send to client.
	Headers map[string]string `json:"headers"`
	// Response body.
	Body interface{} `json:"body"`
}
