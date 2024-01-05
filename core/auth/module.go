package auth

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
)

type coreModule struct{}

func NewCoreModule() core.CoreModule {
	return &coreModule{}
}

func (c *coreModule) Name() string {
	return AuthModule
}

func (c *coreModule) NewTxRouter(logger types.Logger, _ cmclient.TxBuilder, _ string, _ interface{}) (cmclient.TxBuilder, error) {
	return nil, logger.Err(types.ErrWrap(types.ErrInvalidRequest, c.Name(), "module has not tx"))
}

func (c *coreModule) NewQueryRouter(q core.QueryClient) (string, error) {
	return QueryAuth(q)
}
