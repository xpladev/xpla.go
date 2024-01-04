package staking

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &StakingExternal{}

type StakingExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e StakingExternal) {
	e.Xplac = xplac
	e.Name = StakingModule
	return e
}

func (e StakingExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e StakingExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Create new validator initialized with a self-delegation to it.
func (e StakingExternal) CreateValidator(createValidatorMsg types.CreateValidatorMsg) provider.XplaClient {
	msg, err := MakeCreateValidatorMsg(createValidatorMsg, e.Xplac.GetFromAddress(), e.Xplac.GetOutputDocument())
	if err != nil {
		return e.Err(StakingCreateValidatorMsgType, err)
	}

	return e.ToExternal(StakingCreateValidatorMsgType, msg)
}

// Edit an existing validator account.
func (e StakingExternal) EditValidator(editValidatorMsg types.EditValidatorMsg) provider.XplaClient {
	msg, err := MakeEditValidatorMsg(editValidatorMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(StakingEditValidatorMsgType, err)
	}

	return e.ToExternal(StakingEditValidatorMsgType, msg)
}

// Delegate liquid tokens to a validator.
func (e StakingExternal) Delegate(delegateMsg types.DelegateMsg) provider.XplaClient {
	msg, err := MakeDelegateMsg(delegateMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(StakingDelegateMsgType, err)
	}

	return e.ToExternal(StakingDelegateMsgType, msg)
}

// Unbond shares from a validator.
func (e StakingExternal) Unbond(unbondMsg types.UnbondMsg) provider.XplaClient {
	msg, err := MakeUnbondMsg(unbondMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(StakingUnbondMsgType, err)
	}

	return e.ToExternal(StakingUnbondMsgType, msg)
}

// Redelegate illiquid tokens from one validator to another.
func (e StakingExternal) Redelegate(redelegateMsg types.RedelegateMsg) provider.XplaClient {
	msg, err := MakeRedelegateMsg(redelegateMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(StakingRedelegateMsgType, err)
	}

	return e.ToExternal(StakingRedelegateMsgType, msg)
}

// Query

// Query a validator or for all validators.
func (e StakingExternal) QueryValidators(queryValidatorMsg ...types.QueryValidatorMsg) provider.XplaClient {
	switch {
	case len(queryValidatorMsg) == 0:
		msg, err := MakeQueryValidatorsMsg()
		if err != nil {
			return e.Err(StakingQueryValidatorsMsgType, err)
		}

		return e.ToExternal(StakingQueryValidatorsMsgType, msg)

	case len(queryValidatorMsg) == 1:
		msg, err := MakeQueryValidatorMsg(queryValidatorMsg[0])
		if err != nil {
			return e.Err(StakingQueryValidatorMsgType, err)
		}

		return e.ToExternal(StakingQueryValidatorMsgType, msg)

	default:
		return e.Err(StakingQueryValidatorMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}

// Query a delegation based on address and validator address, all out going redelegations from a validator or all delegations made by on delegator.
func (e StakingExternal) QueryDelegation(queryDelegationMsg types.QueryDelegationMsg) provider.XplaClient {
	switch {
	case queryDelegationMsg.DelegatorAddr != "" && queryDelegationMsg.ValidatorAddr != "":
		msg, err := MakeQueryDelegationMsg(queryDelegationMsg)
		if err != nil {
			return e.Err(StakingQueryDelegationMsgType, err)
		}

		return e.ToExternal(StakingQueryDelegationMsgType, msg)

	case queryDelegationMsg.DelegatorAddr != "":
		msg, err := MakeQueryDelegationsMsg(queryDelegationMsg)
		if err != nil {
			return e.Err(StakingQueryDelegationsMsgType, err)
		}

		return e.ToExternal(StakingQueryDelegationsMsgType, msg)

	case queryDelegationMsg.ValidatorAddr != "":
		msg, err := MakeQueryDelegationsToMsg(queryDelegationMsg)
		if err != nil {
			return e.Err(StakingQueryDelegationsToMsgType, err)
		}

		return e.ToExternal(StakingQueryDelegationsToMsgType, msg)

	default:
		return e.Err(StakingQueryDelegationsToMsgType, types.ErrWrap(types.ErrInvalidRequest, "wrong delegation message"))
	}
}

// Query all unbonding delegatations from a validator, an unbonding-delegation record based on delegator and validator address or all unbonding-delegations records for one delegator.
func (e StakingExternal) QueryUnbondingDelegation(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) provider.XplaClient {
	switch {
	case queryUnbondingDelegationMsg.DelegatorAddr != "" && queryUnbondingDelegationMsg.ValidatorAddr != "":
		msg, err := MakeQueryUnbondingDelegationMsg(queryUnbondingDelegationMsg)
		if err != nil {
			return e.Err(StakingQueryUnbondingDelegationMsgType, err)
		}

		return e.ToExternal(StakingQueryUnbondingDelegationMsgType, msg)

	case queryUnbondingDelegationMsg.DelegatorAddr != "":
		msg, err := MakeQueryUnbondingDelegationsMsg(queryUnbondingDelegationMsg)
		if err != nil {
			return e.Err(StakingQueryUnbondingDelegationsMsgType, err)
		}

		return e.ToExternal(StakingQueryUnbondingDelegationsMsgType, msg)

	case queryUnbondingDelegationMsg.ValidatorAddr != "":
		msg, err := MakeQueryUnbondingDelegationsFromMsg(queryUnbondingDelegationMsg)
		if err != nil {
			return e.Err(StakingQueryUnbondingDelegationsFromMsgType, err)
		}

		return e.ToExternal(StakingQueryUnbondingDelegationsFromMsgType, msg)

	default:
		return e.Err(StakingQueryUnbondingDelegationsFromMsgType, types.ErrWrap(types.ErrInvalidRequest, "wrong unbonding delegation message"))
	}
}

// Query a redelegation record based on delegator and a source and destination validator.
// Also, query all outgoing redelegatations from a validator or all redelegations records for one delegator.
func (e StakingExternal) QueryRedelegation(queryRedelegationMsg types.QueryRedelegationMsg) provider.XplaClient {
	switch {
	case queryRedelegationMsg.DelegatorAddr != "" &&
		queryRedelegationMsg.SrcValidatorAddr != "" &&
		queryRedelegationMsg.DstValidatorAddr != "":

		msg, err := MakeQueryRedelegationMsg(queryRedelegationMsg)
		if err != nil {
			return e.Err(StakingQueryRedelegationMsgType, err)
		}

		return e.ToExternal(StakingQueryRedelegationMsgType, msg)

	case queryRedelegationMsg.DelegatorAddr != "":
		msg, err := MakeQueryRedelegationsMsg(queryRedelegationMsg)
		if err != nil {
			return e.Err(StakingQueryRedelegationsMsgType, err)
		}

		return e.ToExternal(StakingQueryRedelegationsMsgType, msg)

	case queryRedelegationMsg.SrcValidatorAddr != "":
		msg, err := MakeQueryRedelegationsFromMsg(queryRedelegationMsg)
		if err != nil {
			return e.Err(StakingQueryRedelegationsFromMsgType, err)
		}

		return e.ToExternal(StakingQueryRedelegationsFromMsgType, msg)

	default:
		return e.Err(StakingQueryRedelegationsFromMsgType, types.ErrWrap(types.ErrInvalidRequest, "wrong redelegation message"))
	}
}

// Query historical info at given height.
func (e StakingExternal) HistoricalInfo(historicalInfoMsg types.HistoricalInfoMsg) provider.XplaClient {
	msg, err := MakeHistoricalInfoMsg(historicalInfoMsg)
	if err != nil {
		return e.Err(StakingHistoricalInfoMsgType, err)
	}

	return e.ToExternal(StakingHistoricalInfoMsgType, msg)
}

// Query the current staking pool values.
func (e StakingExternal) StakingPool() provider.XplaClient {
	msg, err := MakeQueryStakingPoolMsg()
	if err != nil {
		return e.Err(StakingQueryStakingPoolMsgType, err)
	}

	return e.ToExternal(StakingQueryStakingPoolMsgType, msg)
}

// Query the current staking parameters information.
func (e StakingExternal) StakingParams() provider.XplaClient {
	msg, err := MakeQueryStakingParamsMsg()
	if err != nil {
		return e.Err(StakingQueryStakingParamsMsgType, err)
	}

	return e.ToExternal(StakingQueryStakingParamsMsgType, msg)
}
