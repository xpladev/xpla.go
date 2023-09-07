package params

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type coreModule struct{}

func NewCoreModule() core.CoreModule {
	return &coreModule{}
}

func (c *coreModule) Name() string {
	return ParamsModule
}

func (c *coreModule) NewTxRouter(builder cmclient.TxBuilder, msgType string, msg interface{}) (cmclient.TxBuilder, error) {
	switch {
	case msgType == ParamsProposalParamChangeMsgType:
		convertMsg := msg.(govtypes.MsgSubmitProposal)
		builder.SetMsgs(&convertMsg)

	default:
		return nil, util.LogErr(errors.ErrInvalidMsgType, msgType)
	}

	return builder, nil
}

func (c *coreModule) NewQueryRouter(q core.QueryClient) (string, error) {
	return QueryParams(q)
}
