package staking

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type StakingExternal struct {
	Xplac provider.XplaClient
}

func NewStakingExternal(xplac provider.XplaClient) (e StakingExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Create new validator initialized with a self-delegation to it.
func (e StakingExternal) CreateValidator(createValidatorMsg types.CreateValidatorMsg) provider.XplaClient {
	msg, err := MakeCreateValidatorMsg(createValidatorMsg, e.Xplac.GetPrivateKey(), e.Xplac.GetOutputDocument())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(StakingModule).
		WithMsgType(StakingCreateValidatorMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Edit an existing validator account.
func (e StakingExternal) EditValidator(editValidatorMsg types.EditValidatorMsg) provider.XplaClient {
	msg, err := MakeEditValidatorMsg(editValidatorMsg, e.Xplac.GetPrivateKey())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(StakingModule).
		WithMsgType(StakingEditValidatorMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Delegate liquid tokens to a validator.
func (e StakingExternal) Delegate(delegateMsg types.DelegateMsg) provider.XplaClient {
	msg, err := MakeDelegateMsg(delegateMsg, e.Xplac.GetPrivateKey())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(StakingModule).
		WithMsgType(StakingDelegateMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Unbond shares from a validator.
func (e StakingExternal) Unbond(unbondMsg types.UnbondMsg) provider.XplaClient {
	msg, err := MakeUnbondMsg(unbondMsg, e.Xplac.GetPrivateKey())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(StakingModule).
		WithMsgType(StakingUnbondMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Redelegate illiquid tokens from one validator to another.
func (e StakingExternal) Redelegate(redelegateMsg types.RedelegateMsg) provider.XplaClient {
	msg, err := MakeRedelegateMsg(redelegateMsg, e.Xplac.GetPrivateKey())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(StakingModule).
		WithMsgType(StakingRedelegateMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Query a validator or for all validators.
func (e StakingExternal) QueryValidators(queryValidatorMsg ...types.QueryValidatorMsg) provider.XplaClient {
	if len(queryValidatorMsg) == 0 {
		msg, err := MakeQueryValidatorsMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryValidatorsMsgType).
			WithMsg(msg)
	} else if len(queryValidatorMsg) == 1 {
		msg, err := MakeQueryValidatorMsg(queryValidatorMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryValidatorMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}

// Query a delegation based on address and validator address, all out going redelegations from a validator or all delegations made by on delegator.
func (e StakingExternal) QueryDelegation(queryDelegationMsg types.QueryDelegationMsg) provider.XplaClient {
	if queryDelegationMsg.DelegatorAddr != "" && queryDelegationMsg.ValidatorAddr != "" {
		msg, err := MakeQueryDelegationMsg(queryDelegationMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryDelegationMsgType).
			WithMsg(msg)
	} else if queryDelegationMsg.DelegatorAddr != "" {
		msg, err := MakeQueryDelegationsMsg(queryDelegationMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryDelegationsMsgType).
			WithMsg(msg)
	} else if queryDelegationMsg.ValidatorAddr != "" {
		msg, err := MakeQueryDelegationsToMsg(queryDelegationMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryDelegationsToMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "wrong delegation message"))
	}
	return e.Xplac
}

// Query all unbonding delegatations from a validator, an unbonding-delegation record based on delegator and validator address or all unbonding-delegations records for one delegator.
func (e StakingExternal) QueryUnbondingDelegation(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) provider.XplaClient {
	if queryUnbondingDelegationMsg.DelegatorAddr != "" && queryUnbondingDelegationMsg.ValidatorAddr != "" {
		msg, err := MakeQueryUnbondingDelegationMsg(queryUnbondingDelegationMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryUnbondingDelegationMsgType).
			WithMsg(msg)
	} else if queryUnbondingDelegationMsg.DelegatorAddr != "" {
		msg, err := MakeQueryUnbondingDelegationsMsg(queryUnbondingDelegationMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryUnbondingDelegationsMsgType).
			WithMsg(msg)
	} else if queryUnbondingDelegationMsg.ValidatorAddr != "" {
		msg, err := MakeQueryUnbondingDelegationsFromMsg(queryUnbondingDelegationMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryUnbondingDelegationsFromMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "wrong unbonding delegation message"))
	}
	return e.Xplac
}

// Query a redelegation record based on delegator and a source and destination validator.
// Also, query all outgoing redelegatations from a validator or all redelegations records for one delegator.
func (e StakingExternal) QueryRedelegation(queryRedelegationMsg types.QueryRedelegationMsg) provider.XplaClient {
	if queryRedelegationMsg.DelegatorAddr != "" &&
		queryRedelegationMsg.SrcValidatorAddr != "" &&
		queryRedelegationMsg.DstValidatorAddr != "" {
		msg, err := MakeQueryRedelegationMsg(queryRedelegationMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryRedelegationMsgType).
			WithMsg(msg)
	} else if queryRedelegationMsg.DelegatorAddr != "" {
		msg, err := MakeQueryRedelegationsMsg(queryRedelegationMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryRedelegationsMsgType).
			WithMsg(msg)
	} else if queryRedelegationMsg.SrcValidatorAddr != "" {
		msg, err := MakeQueryRedelegationsFromMsg(queryRedelegationMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(StakingModule).
			WithMsgType(StakingQueryRedelegationsFromMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "wrong redelegation message"))
	}
	return e.Xplac
}

// Query historical info at given height.
func (e StakingExternal) HistoricalInfo(historicalInfoMsg types.HistoricalInfoMsg) provider.XplaClient {
	msg, err := MakeHistoricalInfoMsg(historicalInfoMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(StakingModule).
		WithMsgType(StakingHistoricalInfoMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query the current staking pool values.
func (e StakingExternal) StakingPool() provider.XplaClient {
	msg, err := MakeQueryStakingPoolMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(StakingModule).
		WithMsgType(StakingQueryStakingPoolMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query the current staking parameters information.
func (e StakingExternal) StakingParams() provider.XplaClient {
	msg, err := MakeQueryStakingParamsMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(StakingModule).
		WithMsgType(StakingQueryStakingParamsMsgType).
		WithMsg(msg)
	return e.Xplac
}
