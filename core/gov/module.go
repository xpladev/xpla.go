package gov

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
	return GovModule
}

func (c *coreModule) NewTxRouter(builder cmclient.TxBuilder, msgType string, msg interface{}) (cmclient.TxBuilder, error) {
	switch {
	case msgType == GovSubmitProposalMsgType:
		convertMsg := msg.(govtypes.MsgSubmitProposal)
		builder.SetMsgs(&convertMsg)

	case msgType == GovDepositMsgType:
		convertMsg := msg.(govtypes.MsgDeposit)
		builder.SetMsgs(&convertMsg)

	case msgType == GovVoteMsgType:
		convertMsg := msg.(govtypes.MsgVote)
		builder.SetMsgs(&convertMsg)

	case msgType == GovWeightedVoteMsgType:
		convertMsg := msg.(govtypes.MsgVoteWeighted)
		builder.SetMsgs(&convertMsg)

	default:
		return nil, util.LogErr(errors.ErrInvalidMsgType, msgType)
	}

	return builder, nil
}

func (c *coreModule) NewQueryRouter(q core.QueryClient) (string, error) {
	return QueryGov(q)
}
