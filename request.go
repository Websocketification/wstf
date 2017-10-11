package wstf

import (
	"encoding/json"
	"errors"
)

type Request struct {
	ID      string              `json:"id"`
	Method  string              `json:"method"`
	Path    string              `json:"path"`
	Query   string              `json:"query"`
	Headers map[string](string) `json:"headers"`
	Body    string              `json:"body"`
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
