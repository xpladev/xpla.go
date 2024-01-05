package core

import "github.com/xpladev/xpla.go/provider"

// The standard form for generating tx or query message.
type External interface {
	// If generating tx or query message is success, return message.
	ToExternal(string, interface{}) provider.XplaClient

	// If generating message is failed, return error.
	Err(string, error) provider.XplaClient
}
