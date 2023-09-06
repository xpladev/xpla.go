package authz

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

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

func (c *coreModule) NewTxRouter(builder cmclient.TxBuilder, msgType string, msg interface{}) (cmclient.TxBuilder, error) {
	switch {
	case msgType == AuthzGrantMsgType:
		convertMsg := msg.(authz.MsgGrant)
		builder.SetMsgs(&convertMsg)

	case msgType == AuthzRevokeMsgType:
		convertMsg := msg.(authz.MsgRevoke)
		builder.SetMsgs(&convertMsg)

	case msgType == AuthzExecMsgType:
		convertMsg := msg.(authz.MsgExec)
		builder.SetMsgs(&convertMsg)

	default:
		return nil, util.LogErr(errors.ErrInvalidMsgType, msgType)
	}

	return builder, nil
}

func (c *coreModule) NewQueryRouter(q core.QueryClient) (string, error) {
	return QueryAuthz(q)
}
