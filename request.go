package wstf

import (
	"encoding/json"
	"errors"
)

// @see http://eagain.net/articles/go-dynamic-json/
// @see json.RawMessage
type Request struct {
	// The request unique ID.
	Id string `json:"id"`
	// HTTP Method.
	Method string `json:"method"`
	// Non-empty path.
	Path string `json:"path"`
	// Params matched from path.
	Params map[string]string `json:"params"`
	// The query part of request which is a JSON object and usually typed as map[string]string.
	Query json.RawMessage `json:"query"`
	// Headers of the request.
	Headers map[string]string `json:"headers"`
	// The request body of request, a JSON object.
	Body json.RawMessage `json:"body"`
}

// The json.RawMessage type is used for some fields(Query & Body) of request
// to store the original partial json bytes and delay JSON decoding.
func NewRequest(jsonBytes []byte) (*Request, error) {
	var req *Request
	err := json.Unmarshal(jsonBytes, &req)
	req.Params = make(map[string]string)
	if !req.IsRequestValid() {
		return nil, errors.New("expected fields is empty")
	}
	return req, err
}

// Is the request valid.
func (m *Request) IsRequestValid() bool {
	return !(m.Id == "" || m.Method == "" || m.Path == "")
}

// Unmarshal the request.Query into a given struct.
// The obj should be pointer of some struct.
func (m *Request) UnmarshalQuery(obj interface{}) error {
	return json.Unmarshal(m.Query, obj)
}

// Unmarshal the request.Body into a given struct.
// The obj should be pointer of some struct.
func (m *Request) UnmarshalBody(obj interface{}) error {
	return json.Unmarshal(m.Body, obj)
}

// Get the query unmarshaller.
func (m *Request) GetQueryUnmarshaller() (*StructUnmarshaller, error) {
	unmarshaller := &StructUnmarshaller{}
	return unmarshaller, json.Unmarshal(m.Query, &unmarshaller.DataMap)
}

// Get the body unmarshaller.
func (m *Request) GetBodyUnmarshaller() (*StructUnmarshaller, error) {
	unmarshaller := &StructUnmarshaller{
		RawMessage: m.Body,
	}
	return unmarshaller, json.Unmarshal(m.Body, &unmarshaller.DataMap)
}
