package authz

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

type coreModule struct{}

func NewCoreModule() core.CoreModule {
	return &coreModule{}
}

func (c *coreModule) Name() string {
	return AuthzModule
}

func (c *coreModule) NewTxRouter(logger types.Logger, builder cmclient.TxBuilder, msgType string, msg interface{}) (cmclient.TxBuilder, error) {
	switch {
	case msgType == AuthzGrantMsgType:
		convertMsg := msg.(authz.MsgGrant)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == AuthzRevokeMsgType:
		convertMsg := msg.(authz.MsgRevoke)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == AuthzExecMsgType:
		convertMsg := msg.(authz.MsgExec)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	default:
		return nil, logger.Err(types.ErrWrap(types.ErrInvalidMsgType, msgType))
	}

	return builder, nil
}

func (c *coreModule) NewQueryRouter(q core.QueryClient) (string, error) {
	return QueryAuthz(q)
}
