// The server implementation of the WSTF protocol, packaged as "wstf" in Golang.
package wstf

// Versions of the WSTF protocol.
// @see http.Request#Proto
const (
	// The WSTF protocol starting from the version of 0.8 which is under development.
	Protocol08 = "WSTF/0.8"
	// The major version of the protocol.
	// The 0 means it is under development and may not ready for production.
	ProtocolMajor0 = 0
	// The minor version of the protocol.
	ProtocolMinor8 = 8
)

// The protocol supported in the current package.
const ProtocolCurrent = Protocol08
