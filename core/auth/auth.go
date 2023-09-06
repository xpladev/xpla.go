package auth

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
)

type coreModule struct{}

func NewCoreModule() core.CoreModule {
	return &coreModule{}
}

func (c *coreModule) Name() string {
	return AuthModule
}

func (c *coreModule) NewTxRouter(txBuiler cmclient.TxBuilder, msgType string) (cmclient.TxBuilder, error) {
	return nil, util.LogErr(errors.ErrInvalidRequest, c.Name(), "module has not tx")
}

func (c *coreModule) NewQueryRouter(q core.QueryClient) (string, error) {
	return QueryAuth(q)
}
