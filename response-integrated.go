package wstf

import (
	"fmt"
	"net/http"
)

// Frequently used methods that call raw methods.

// Response the request with status code 200.
func (m *Response) Done(body interface{}) {
	m.SetStatusCode(http.StatusOK).SetBody(body).End()
}

// Response the request with specific status code.
func (m *Response) Error(statusCode int, error interface{}) {
	m.SetStatusCode(statusCode).End()
	fmt.Println("Responsing Error:", error)
}
