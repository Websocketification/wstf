package wstf

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Response struct {
	// The reference of wstf Connection.
	Connection *Connection `json:"-"`
	// Shortcut for Connection.Locals
	ConnectionLocals map[string]interface{}
	// The reference of corresponding wstf Request.
	Request *Request
	// JSON Response that will be sent as response to corresponding request.
	JsonResponse *JsonResponse
	// A map that contains response local variables scoped to the request.
	// This property is useful for exposing request-level information such as the request path name, authenticated user, user settings, and so on.
	Locals map[string]interface{}
}

func NewResponse(connection *Connection, connectionLocals map[string]interface{}, request *Request) *Response {
	return &Response{
		connection,
		connectionLocals,
		request,
		&JsonResponse{Id: request.Id},
		make(map[string]interface{}),
	}
}

// The raw methods to access Response.JsonResponse.

// Set status code.
func (m *Response) SetStatusCode(statusCode int) *Response {
	m.JsonResponse.Status = statusCode
	return m
}

// Set header.
func (m *Response) SetHeader(key, value string) {
	m.JsonResponse.Headers[key] = value
}

// Set json-response body.
func (m *Response) SetBody(body interface{}) *Response {
	m.JsonResponse.Body = body
	return m
}

// Finish the request.
func (m *Response) End() error {
	app := m.Connection.Application
	if app.OnEncodingResponse != nil {
		// CALLBACK for all responses.
		app.OnEncodingResponse(m.Request, m)
	}
	bytes, err := json.Marshal(m.JsonResponse)
	if err != nil {
		if app.OnFailedToEncodeResponse != nil {
			app.OnFailedToEncodeResponse(err, m.Request, m)
		} else {
			fmt.Println("[WSTF] Failed to encode the JSON Response:", m.JsonResponse)
		}
		return err
	}
	return m.Write(websocket.TextMessage, bytes)
}

// Send the response to the client.
func (m *Response) Write(mt int, message []byte) error {
	// FIX-ME Support debugging mode( , which can be used for unit test).
	if m.Connection == nil {
		fmt.Println("DEBUGGING MODE: Sending Message: ", string(message))
		return nil
	}
	app := m.Connection.Application
	if app.OnSendingResponse != nil {
		// CALLBACK for all responses.
		app.OnSendingResponse(m.Request, m, message)
	}
	err := m.Connection.WebSocketConn.WriteMessage(mt, message)
	if err != nil {
		if app.OnFailedToSendResponse != nil {
			app.OnFailedToSendResponse(err, m.Request, m, message)
		} else {
			fmt.Println("[WSTF] Failed to send the response to client:", err, message)
		}
	}
	return err
}
