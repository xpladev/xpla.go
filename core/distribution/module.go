package distribution

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"

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

func (c *coreModule) NewTxRouter(logger types.Logger, builder cmclient.TxBuilder, msgType string, msg interface{}) (cmclient.TxBuilder, error) {
	switch {
	case msgType == DistributionFundCommunityPoolMsgType:
		convertMsg := msg.(disttypes.MsgFundCommunityPool)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == DistributionProposalCommunityPoolSpendMsgType:
		convertMsg := msg.(govtypes.MsgSubmitProposal)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == DistributionWithdrawRewardsMsgType:
		convertMsg := msg.([]sdk.Msg)
		err := builder.SetMsgs(convertMsg...)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == DistributionWithdrawAllRewardsMsgType:
		convertMsg := msg.([]sdk.Msg)
		err := builder.SetMsgs(convertMsg...)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == DistributionSetWithdrawAddrMsgType:
		convertMsg := msg.(disttypes.MsgSetWithdrawAddress)
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
	return QueryDistribution(q)
}
