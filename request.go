package wstf

import (
	"encoding/json"
	"errors"
)

type Request struct {
	// The request unique ID.
	ID string `json:"id"`
	// HTTP Method.
	Method string `json:"method"`
	// Non-empty path.
	Path string `json:"path"`
	// The query part of request which is a JSON object and usually typed as map[string]string.
	Query interface{} `json:"query"`
	// Headers of the request.
	Headers map[string](string) `json:"headers"`
	// The request body of request, a JSON object.
	Body interface{} `json:"body"`
}

func NewRequest(jsonBytes []byte) (*Request, error) {
	var req *Request
	err := json.Unmarshal(jsonBytes, &req)
	if !req.IsRequestValid() {
		return nil, errors.New("expected fields is empty")
	}
	return req, err
}

// Is the request valid.
func (m Request) IsRequestValid() bool {
	return !(m.ID == "" || m.Method == "" || m.Path == "")
}
