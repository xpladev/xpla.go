package distribution

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type DistributionExternal struct {
	Xplac provider.XplaClient
}

func NewDistributionExternal(xplac provider.XplaClient) (e DistributionExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Funds the community pool with the specified amount.
func (e DistributionExternal) FundCommunityPool(fundCommunityPoolMsg types.FundCommunityPoolMsg) provider.XplaClient {
	msg, err := MakeFundCommunityPoolMsg(fundCommunityPoolMsg, e.Xplac.GetPrivateKey())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(DistributionModule).
		WithMsgType(DistributionFundCommunityPoolMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Submit a community pool spend proposal.
func (e DistributionExternal) CommunityPoolSpend(communityPoolSpendMsg types.CommunityPoolSpendMsg) provider.XplaClient {
	msg, err := MakeProposalCommunityPoolSpendMsg(communityPoolSpendMsg, e.Xplac.GetPrivateKey(), e.Xplac.GetEncoding())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(DistributionModule).
		WithMsgType(DistributionProposalCommunityPoolSpendMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator.
func (e DistributionExternal) WithdrawRewards(withdrawRewardsMsg types.WithdrawRewardsMsg) provider.XplaClient {
	msg, err := MakeWithdrawRewardsMsg(withdrawRewardsMsg, e.Xplac.GetPrivateKey())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(DistributionModule).
		WithMsgType(DistributionWithdrawRewardsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Withdraw all delegations rewards for a delegator.
func (e DistributionExternal) WithdrawAllRewards() provider.XplaClient {
	msg, err := MakeWithdrawAllRewardsMsg(e.Xplac.GetPrivateKey(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(DistributionModule).
		WithMsgType(DistributionWithdrawAllRewardsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Change the default withdraw address for rewards associated with an address.
func (e DistributionExternal) SetWithdrawAddr(setWithdrawAddrMsg types.SetWithdrawAddrMsg) provider.XplaClient {
	msg, err := MakeSetWithdrawAddrMsg(setWithdrawAddrMsg, e.Xplac.GetPrivateKey())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(DistributionModule).
		WithMsgType(DistributionSetWithdrawAddrMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Query distribution parameters.
func (e DistributionExternal) DistributionParams() provider.XplaClient {
	msg, err := MakeQueryDistributionParamsMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(DistributionModule).
		WithMsgType(DistributionQueryDistributionParamsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations.
func (e DistributionExternal) ValidatorOutstandingRewards(validatorOutstandingRewardsMsg types.ValidatorOutstandingRewardsMsg) provider.XplaClient {
	msg, err := MakeValidatorOutstandingRewardsMsg(validatorOutstandingRewardsMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(DistributionModule).
		WithMsgType(DistributionValidatorOutstandingRewardsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query distribution validator commission.
func (e DistributionExternal) DistCommission(queryDistCommissionMsg types.QueryDistCommissionMsg) provider.XplaClient {
	msg, err := MakeQueryDistCommissionMsg(queryDistCommissionMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(DistributionModule).
		WithMsgType(DistributionQueryDistCommissionMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query distribution validator slashes.
func (e DistributionExternal) DistSlashes(queryDistSlashesMsg types.QueryDistSlashesMsg) provider.XplaClient {
	msg, err := MakeQueryDistSlashesMsg(queryDistSlashesMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(DistributionModule).
		WithMsgType(DistributionQuerySlashesMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query all ditribution delegator rewards or rewards from a particular validator.
func (e DistributionExternal) DistRewards(queryDistRewardsMsg types.QueryDistRewardsMsg) provider.XplaClient {
	if queryDistRewardsMsg.DelegatorAddr == "" {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInsufficientParams, "must set a delegator address"))
	}

	if queryDistRewardsMsg.ValidatorAddr != "" {
		msg, err := MakeQueryDistRewardsMsg(queryDistRewardsMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(DistributionModule).
			WithMsgType(DistributionQueryRewardsMsgType).
			WithMsg(msg)
	} else {
		msg, err := MakeQueryDistTotalRewardsMsg(queryDistRewardsMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(DistributionModule).
			WithMsgType(DistributionQueryTotalRewardsMsgType).
			WithMsg(msg)
	}
	return e.Xplac
}

// Query the amount of coins in the community pool.
func (e DistributionExternal) CommunityPool() provider.XplaClient {
	msg, err := MakeQueryCommunityPoolMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(DistributionModule).
		WithMsgType(DistributionQueryCommunityPoolMsgType).
		WithMsg(msg)
	return e.Xplac
}
