package reward

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

type RewardExternal struct {
	Xplac provider.XplaClient
}

func NewRewardExternal(xplac provider.XplaClient) (e RewardExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Funds the fee collector with the specified amount
func (e RewardExternal) FundFeeCollector(fundFeeCollectorMsg types.FundFeeCollectorMsg) provider.XplaClient {
	msg, err := MakeFundFeeCollectorMsg(fundFeeCollectorMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(RewardModule).
		WithMsgType(RewardFundFeeCollectorMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Query reward params
func (e RewardExternal) RewardParams() provider.XplaClient {
	msg, err := MakeQueryRewardParamsMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(RewardModule).
		WithMsgType(RewardQueryRewardParamsMsgType).
		WithMsg(msg)

	return e.Xplac
}

// Query reward pool
func (e RewardExternal) RewardPool() provider.XplaClient {
	msg, err := MakeQueryRewardPoolMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(RewardModule).
		WithMsgType(RewardQueryRewardPoolMsgType).
		WithMsg(msg)

	return e.Xplac
}
