package wstf

// The JSON response used to response to corresponding wstf request.
type JsonResponse struct {
	// Unique http request identifier.
	Id string `json:"id"`
	// Headers to be send to client.
	Headers map[string]string `json:"headers"`
	// Http status.
	Status int `json:"status"`
	// Response body.
	Body interface{} `json:"body"`
}

func (m *JsonResponse) ToJson() string {
	return Stringify(m)
}
