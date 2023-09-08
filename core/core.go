package core

import cmclient "github.com/cosmos/cosmos-sdk/client"

// The standard form for a module in the core package.
// Every modules are enrolled to the controller by using this interface.
type CoreModule interface {
	// Name of a core module must not be duplicated previous names.
	Name() string

	// Routed transaction messages are built in the TxBuilder of Cosmos-SDK.
	NewTxRouter(cmclient.TxBuilder, string, interface{}) (cmclient.TxBuilder, error)

	// Route query requests by gRPC or HTTP.
	// Queries are returned with string type regardless of communication protocol.
	NewQueryRouter(QueryClient) (string, error)
}
