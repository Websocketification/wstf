package wstf

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var WebSocketHeaderRequirements = map[string]string{
	"Connection": "upgrade",

	"Upgrade": "websocket",

	"Sec-Websocket-Version": "13",

	//"Pragma": "no-cache",
	//"Cache-Control": "no-cache",
}

// The helper function for the upgrader to check whether a request can be upgraded.
// @see https://tools.ietf.org/html/rfc6455#page-60
// @see https://github.com/gorilla/websocket/blob/master/server.go#L112
func AbleToUpgrade(upgrader *websocket.Upgrader, w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "GET" {
		return false
	}
	// Check headers.
	for k, v := range WebSocketHeaderRequirements {
		if strings.ToLower(r.Header.Get(k)) != v {
			return false
		}
	}
	if r.Header.Get("Sec-Websocket-Key") == "" {
		return false
	}
	// Check origin.
	if upgrader.CheckOrigin == nil || !upgrader.CheckOrigin(r) {
		return false
	}

	_, ok := w.(http.Hijacker)
	return ok
}
