package core

import cmclient "github.com/cosmos/cosmos-sdk/client"

type CoreModule interface {
	Name() string
	NewTxRouter(cmclient.TxBuilder, string) (cmclient.TxBuilder, error)
	NewQueryRouter(QueryClient) (string, error)
}
