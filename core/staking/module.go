package staking

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

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

func (c *coreModule) NewTxRouter(builder cmclient.TxBuilder, msgType string, msg interface{}) (cmclient.TxBuilder, error) {
	switch {
	case msgType == StakingCreateValidatorMsgType:
		convertMsg := msg.(sdk.Msg)
		builder.SetMsgs(convertMsg)

	case msgType == StakingEditValidatorMsgType:
		convertMsg := msg.(stakingtypes.MsgEditValidator)
		builder.SetMsgs(&convertMsg)

	case msgType == StakingDelegateMsgType:
		convertMsg := msg.(stakingtypes.MsgDelegate)
		builder.SetMsgs(&convertMsg)

	case msgType == StakingUnbondMsgType:
		convertMsg := msg.(stakingtypes.MsgUndelegate)
		builder.SetMsgs(&convertMsg)

	case msgType == StakingRedelegateMsgType:
		convertMsg := msg.(stakingtypes.MsgBeginRedelegate)
		builder.SetMsgs(&convertMsg)

	default:
		return nil, util.LogErr(errors.ErrInvalidMsgType, msgType)
	}

	return builder, nil
}

func (c *coreModule) NewQueryRouter(q core.QueryClient) (string, error) {
	return QueryStaking(q)
}
