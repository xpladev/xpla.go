package reward

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &RewardExternal{}

type RewardExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e RewardExternal) {
	e.Xplac = xplac
	e.Name = RewardModule
	return e
}

func (e RewardExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e RewardExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Funds the fee collector with the specified amount
func (e RewardExternal) FundFeeCollector(fundFeeCollectorMsg types.FundFeeCollectorMsg) provider.XplaClient {
	msg, err := MakeFundFeeCollectorMsg(fundFeeCollectorMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(RewardFundFeeCollectorMsgType, err)
	}

	return e.ToExternal(RewardFundFeeCollectorMsgType, msg)
}

// Query

// Query reward params
func (e RewardExternal) RewardParams() provider.XplaClient {
	msg, err := MakeQueryRewardParamsMsg()
	if err != nil {
		return e.Err(RewardQueryRewardParamsMsgType, err)
	}

	return e.ToExternal(RewardQueryRewardParamsMsgType, msg)
}

// Query reward pool
func (e RewardExternal) RewardPool() provider.XplaClient {
	msg, err := MakeQueryRewardPoolMsg()
	if err != nil {
		return e.Err(RewardQueryRewardPoolMsgType, err)
	}

	return e.ToExternal(RewardQueryRewardPoolMsgType, msg)
}
