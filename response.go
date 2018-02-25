package wstf

import (
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

func NewResponse(conn *websocket.Conn, connectionLocals map[string]interface{}, request *http.Request, requestID string) *Response {
	res := &Response{}
	res.Connection = conn
	res.ConnectionLocals = connectionLocals
	res.HttpRequest = request
	res.JsonResponse = &JsonResponse{ID: requestID}
	res.Locals = map[string]interface{}{}
	return res
}

func (m *Response) SetStatusCode(statusCode int) *Response {
	m.JsonResponse.Status = statusCode
	return m
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

// Response the request with status code 200.
func (m *Response) Done(body interface{}) {
	m.SetStatusCode(http.StatusOK)
	m.JsonResponse.Body = body
	m.End()
}

// Response the request with specific status code.
func (m *Response) Error(statusCode int, errors ... interface{}) {
	m.SetStatusCode(statusCode)
	fmt.Println("Responsing Error: ", errors)
	m.End()
}
