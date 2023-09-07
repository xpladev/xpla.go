package core

import cmclient "github.com/cosmos/cosmos-sdk/client"

// The standard form for a module in the core package.
type CoreModule interface {
	Name() string
	NewTxRouter(cmclient.TxBuilder, string, interface{}) (cmclient.TxBuilder, error)
	NewQueryRouter(QueryClient) (string, error)
}
