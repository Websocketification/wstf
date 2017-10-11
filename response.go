package wstf

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"net/http"
	"fmt"
)

type Response struct {
	Connection   *websocket.Conn
	JsonResponse JsonResponse
}

type JsonResponse struct {
	// Unique http request identifier.
	ID string `json:"id"`
	// Headers to be send to client.
	Headers map[string](string) `json:"headers"`
	// Http status.
	Status int `json:"status"`
	// Response body.
	Body interface{} `json:"body"`
}

func (m JsonResponse) ToJson() string {
	str, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(str)
}

func NewResponse(conn *websocket.Conn, ID string) Response {
	res := Response{}
	res.Connection = conn
	res.JsonResponse = JsonResponse{ID: ID}
	return res
}

// Response the request with status code 200.
func (m Response) Done(body interface{}) {
	m.JsonResponse.Status = http.StatusOK
	m.JsonResponse.Body = body
	m.End()
}

// Response the request with specific status code.
func (m Response) Error(httpStatusCode int, error interface{}) {
	m.JsonResponse.Status = httpStatusCode
	m.End()
}

// Set header.
func (m Response) SetHeader(key, value string) {
	m.JsonResponse.Headers[key] = value
}

// Finish the request.
func (m Response) End() error {
	return m.Write(websocket.TextMessage, []byte(m.JsonResponse.ToJson()))
}

// Response to client.
func (m Response) Write(mt int, message []byte) error {
	if m.Connection == nil {
		fmt.Println("DEBUGGING MODE: Sending Message: ", string(message))
		return nil
	}
	err := m.Connection.WriteMessage(mt, message)
	return err
}
