package wstf

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Response struct {
	// The websocket connection.
	Connection *websocket.Conn
	// Locals scoped among the connection.
	ConnectionLocals map[string]interface{}
	// The original http request.
	HttpRequest *http.Request
	// JSON Response that will be sent as response to corresponding request.
	JsonResponse *JsonResponse
	// A map that contains response local variables scoped to the request.
	// This property is useful for exposing request-level information such as the request path name, authenticated user, user settings, and so on.
	Locals map[string]interface{}
}

type JsonResponse struct {
	// Unique http request identifier.
	ID string `json:"id"`
	// Headers to be send to client.
	Headers map[string]string `json:"headers"`
	// Http status.
	Status int `json:"status"`
	// Response body.
	Body interface{} `json:"body"`
}

func (m *JsonResponse) ToJson() string {
	str, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(str)
}

func NewResponse(conn *websocket.Conn, connectionLocals map[string]interface{}, request *http.Request, requestID string) *Response {
	res := &Response{}
	res.Connection = conn
	res.ConnectionLocals = connectionLocals
	res.HttpRequest = request
	res.JsonResponse = &JsonResponse{ID: requestID}
	res.Locals = map[string]interface{}{}
	return res
}

// Response the request with status code 200.
func (m *Response) Done(body interface{}) {
	m.JsonResponse.Status = http.StatusOK
	m.JsonResponse.Body = body
	m.End()
}

// Response the request with specific status code.
func (m *Response) Error(httpStatusCode int, error interface{}) {
	m.JsonResponse.Status = httpStatusCode
	fmt.Println("Responsing Error: ", error)
	m.End()
}

// Set header.
func (m *Response) SetHeader(key, value string) {
	m.JsonResponse.Headers[key] = value
}

// Finish the request.
func (m *Response) End() error {
	return m.Write(websocket.TextMessage, []byte(m.JsonResponse.ToJson()))
}

// Response to client.
func (m *Response) Write(mt int, message []byte) error {
	if m.Connection == nil {
		fmt.Println("DEBUGGING MODE: Sending Message: ", string(message))
		return nil
	}
	err := m.Connection.WriteMessage(mt, message)
	return err
}
