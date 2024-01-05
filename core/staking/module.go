package staking

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type coreModule struct{}

func NewCoreModule() core.CoreModule {
	return &coreModule{}
}

func (c *coreModule) Name() string {
	return StakingModule
}

func (c *coreModule) NewTxRouter(logger types.Logger, builder cmclient.TxBuilder, msgType string, msg interface{}) (cmclient.TxBuilder, error) {
	switch {
	case msgType == StakingCreateValidatorMsgType:
		convertMsg := msg.(sdk.Msg)
		builder.SetMsgs(convertMsg)

	case msgType == StakingEditValidatorMsgType:
		convertMsg := msg.(stakingtypes.MsgEditValidator)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == StakingDelegateMsgType:
		convertMsg := msg.(stakingtypes.MsgDelegate)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == StakingUnbondMsgType:
		convertMsg := msg.(stakingtypes.MsgUndelegate)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == StakingRedelegateMsgType:
		convertMsg := msg.(stakingtypes.MsgBeginRedelegate)
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
	return QueryStaking(q)
}
