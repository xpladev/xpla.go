package distribution

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &DistributionExternal{}

type DistributionExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e DistributionExternal) {
	e.Xplac = xplac
	e.Name = DistributionModule
	return e
}

func (e DistributionExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e DistributionExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Funds the community pool with the specified amount.
func (e DistributionExternal) FundCommunityPool(fundCommunityPoolMsg types.FundCommunityPoolMsg) provider.XplaClient {
	msg, err := MakeFundCommunityPoolMsg(fundCommunityPoolMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(DistributionFundCommunityPoolMsgType, err)
	}

	return e.ToExternal(DistributionFundCommunityPoolMsgType, msg)
}

// Submit a community pool spend proposal.
func (e DistributionExternal) CommunityPoolSpend(communityPoolSpendMsg types.CommunityPoolSpendMsg) provider.XplaClient {
	msg, err := MakeProposalCommunityPoolSpendMsg(communityPoolSpendMsg, e.Xplac.GetFromAddress(), e.Xplac.GetEncoding())
	if err != nil {
		return e.Err(DistributionProposalCommunityPoolSpendMsgType, err)
	}

	return e.ToExternal(DistributionProposalCommunityPoolSpendMsgType, msg)
}

// Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator.
func (e DistributionExternal) WithdrawRewards(withdrawRewardsMsg types.WithdrawRewardsMsg) provider.XplaClient {
	msg, err := MakeWithdrawRewardsMsg(withdrawRewardsMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(DistributionWithdrawRewardsMsgType, err)
	}

	return e.ToExternal(DistributionWithdrawRewardsMsgType, msg)
}

// Withdraw all delegations rewards for a delegator.
func (e DistributionExternal) WithdrawAllRewards() provider.XplaClient {
	msg, err := MakeWithdrawAllRewardsMsg(e.Xplac.GetFromAddress(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext())
	if err != nil {
		return e.Err(DistributionWithdrawAllRewardsMsgType, err)
	}

	return e.ToExternal(DistributionWithdrawAllRewardsMsgType, msg)
}

// Change the default withdraw address for rewards associated with an address.
func (e DistributionExternal) SetWithdrawAddr(setWithdrawAddrMsg types.SetWithdrawAddrMsg) provider.XplaClient {
	msg, err := MakeSetWithdrawAddrMsg(setWithdrawAddrMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(DistributionSetWithdrawAddrMsgType, err)
	}

	return e.ToExternal(DistributionSetWithdrawAddrMsgType, msg)
}

// Query

// Query distribution parameters.
func (e DistributionExternal) DistributionParams() provider.XplaClient {
	msg, err := MakeQueryDistributionParamsMsg()
	if err != nil {
		return e.Err(DistributionQueryDistributionParamsMsgType, err)
	}

	return e.ToExternal(DistributionQueryDistributionParamsMsgType, msg)
}

// Query distribution outstanding (un-withdrawn) rewards for a validator and all their delegations.
func (e DistributionExternal) ValidatorOutstandingRewards(validatorOutstandingRewardsMsg types.ValidatorOutstandingRewardsMsg) provider.XplaClient {
	msg, err := MakeValidatorOutstandingRewardsMsg(validatorOutstandingRewardsMsg)
	if err != nil {
		return e.Err(DistributionValidatorOutstandingRewardsMsgType, err)
	}

	return e.ToExternal(DistributionValidatorOutstandingRewardsMsgType, msg)
}

// Query distribution validator commission.
func (e DistributionExternal) DistCommission(queryDistCommissionMsg types.QueryDistCommissionMsg) provider.XplaClient {
	msg, err := MakeQueryDistCommissionMsg(queryDistCommissionMsg)
	if err != nil {
		return e.Err(DistributionQueryDistCommissionMsgType, err)
	}

	return e.ToExternal(DistributionQueryDistCommissionMsgType, msg)
}

// Query distribution validator slashes.
func (e DistributionExternal) DistSlashes(queryDistSlashesMsg types.QueryDistSlashesMsg) provider.XplaClient {
	msg, err := MakeQueryDistSlashesMsg(queryDistSlashesMsg)
	if err != nil {
		return e.Err(DistributionQuerySlashesMsgType, err)
	}

	return e.ToExternal(DistributionQuerySlashesMsgType, msg)
}

// Query all ditribution delegator rewards or rewards from a particular validator.
func (e DistributionExternal) DistRewards(queryDistRewardsMsg types.QueryDistRewardsMsg) provider.XplaClient {
	if queryDistRewardsMsg.DelegatorAddr == "" {
		return e.Err(DistributionQueryRewardsMsgType, types.ErrWrap(types.ErrInsufficientParams, "must set a delegator address"))
	}

	switch {
	case queryDistRewardsMsg.ValidatorAddr != "":
		msg, err := MakeQueryDistRewardsMsg(queryDistRewardsMsg)
		if err != nil {
			return e.Err(DistributionQueryRewardsMsgType, err)
		}

		return e.ToExternal(DistributionQueryRewardsMsgType, msg)

	default:
		msg, err := MakeQueryDistTotalRewardsMsg(queryDistRewardsMsg)
		if err != nil {
			return e.Err(DistributionQueryTotalRewardsMsgType, err)
		}

		return e.ToExternal(DistributionQueryTotalRewardsMsgType, msg)
	}
}

// Query the amount of coins in the community pool.
func (e DistributionExternal) CommunityPool() provider.XplaClient {
	msg, err := MakeQueryCommunityPoolMsg()
	if err != nil {
		return e.Err(DistributionQueryCommunityPoolMsgType, err)
	}

	return e.ToExternal(DistributionQueryCommunityPoolMsgType, msg)
}
