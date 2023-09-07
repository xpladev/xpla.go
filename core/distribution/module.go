package distribution

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type coreModule struct{}

func NewCoreModule() core.CoreModule {
	return &coreModule{}
}

func (c *coreModule) Name() string {
	return DistributionModule
}

func (c *coreModule) NewTxRouter(builder cmclient.TxBuilder, msgType string, msg interface{}) (cmclient.TxBuilder, error) {
	switch {
	case msgType == DistributionFundCommunityPoolMsgType:
		convertMsg := msg.(disttypes.MsgFundCommunityPool)
		builder.SetMsgs(&convertMsg)

	case msgType == DistributionProposalCommunityPoolSpendMsgType:
		convertMsg := msg.(govtypes.MsgSubmitProposal)
		builder.SetMsgs(&convertMsg)

	case msgType == DistributionWithdrawRewardsMsgType:
		convertMsg := msg.([]sdk.Msg)
		builder.SetMsgs(convertMsg...)

	case msgType == DistributionWithdrawAllRewardsMsgType:
		convertMsg := msg.([]sdk.Msg)
		builder.SetMsgs(convertMsg...)

	case msgType == DistributionSetWithdrawAddrMsgType:
		convertMsg := msg.(disttypes.MsgSetWithdrawAddress)
		builder.SetMsgs(&convertMsg)

	default:
		return nil, util.LogErr(errors.ErrInvalidMsgType, msgType)
	}

	return builder, nil
}

func (c *coreModule) NewQueryRouter(q core.QueryClient) (string, error) {
	return QueryDistribution(q)
}
