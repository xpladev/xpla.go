package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
)

// (Tx) make msg - create validator
func MakeCreateValidatorMsg(createValidatorMsg types.CreateValidatorMsg, privKey key.PrivateKey, output string) (sdk.Msg, error) {
	return parseCreateValidatorArgs(createValidatorMsg, privKey, output)
}

// (Tx) make msg - edit validator
func MakeEditValidatorMsg(editValidatorMsg types.EditValidatorMsg, privKey key.PrivateKey) (stakingtypes.MsgEditValidator, error) {
	return parseEditValidatorArgs(editValidatorMsg, privKey)
}

// (Tx) make msg - delegate
func MakeDelegateMsg(delegateMsg types.DelegateMsg, privKey key.PrivateKey) (stakingtypes.MsgDelegate, error) {
	return parseDelegateArgs(delegateMsg, privKey)
}

// (Tx) make msg - unbond
func MakeUnbondMsg(unbondMsg types.UnbondMsg, privKey key.PrivateKey) (stakingtypes.MsgUndelegate, error) {
	return parseUnbondArgs(unbondMsg, privKey)
}

// (Tx) make msg - redelegate
func MakeRedelegateMsg(redelegateMsg types.RedelegateMsg, privKey key.PrivateKey) (stakingtypes.MsgBeginRedelegate, error) {
	return parseRedelegateArgs(redelegateMsg, privKey)
}

// (Query) make msg - validator
func MakeQueryValidatorMsg(queryValidatorMsg types.QueryValidatorMsg) (stakingtypes.QueryValidatorRequest, error) {
	return stakingtypes.QueryValidatorRequest{
		ValidatorAddr: queryValidatorMsg.ValidatorAddr,
	}, nil
}

// (Query) make msg - validators
func MakeQueryValidatorsMsg() (stakingtypes.QueryValidatorsRequest, error) {
	return stakingtypes.QueryValidatorsRequest{
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - query delegation
func MakeQueryDelegationMsg(queryDelegationMsg types.QueryDelegationMsg) (stakingtypes.QueryDelegationRequest, error) {
	return stakingtypes.QueryDelegationRequest{
		DelegatorAddr: queryDelegationMsg.DelegatorAddr,
		ValidatorAddr: queryDelegationMsg.ValidatorAddr,
	}, nil
}

// (Query) make msg - query delegations
func MakeQueryDelegationsMsg(queryDelegationMsg types.QueryDelegationMsg) (stakingtypes.QueryDelegatorDelegationsRequest, error) {
	return stakingtypes.QueryDelegatorDelegationsRequest{
		DelegatorAddr: queryDelegationMsg.DelegatorAddr,
		Pagination:    core.PageRequest,
	}, nil
}

// (Query) make msg - query delegations to
func MakeQueryDelegationsToMsg(queryDelegationMsg types.QueryDelegationMsg) (stakingtypes.QueryValidatorDelegationsRequest, error) {
	return stakingtypes.QueryValidatorDelegationsRequest{
		ValidatorAddr: queryDelegationMsg.ValidatorAddr,
		Pagination:    core.PageRequest,
	}, nil
}

// (Query) make msg - query unbonding delegation
func MakeQueryUnbondingDelegationMsg(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) (stakingtypes.QueryUnbondingDelegationRequest, error) {
	return stakingtypes.QueryUnbondingDelegationRequest{
		DelegatorAddr: queryUnbondingDelegationMsg.DelegatorAddr,
		ValidatorAddr: queryUnbondingDelegationMsg.ValidatorAddr,
	}, nil
}

// (Query) make msg - query unbonding delegations
func MakeQueryUnbondingDelegationsMsg(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) (stakingtypes.QueryDelegatorUnbondingDelegationsRequest, error) {
	return stakingtypes.QueryDelegatorUnbondingDelegationsRequest{
		DelegatorAddr: queryUnbondingDelegationMsg.DelegatorAddr,
		Pagination:    core.PageRequest,
	}, nil
}

// (Query) make msg - query unbonding delegations from
func MakeQueryUnbondingDelegationsFromMsg(queryUnbondingDelegationMsg types.QueryUnbondingDelegationMsg) (stakingtypes.QueryValidatorUnbondingDelegationsRequest, error) {
	return stakingtypes.QueryValidatorUnbondingDelegationsRequest{
		ValidatorAddr: queryUnbondingDelegationMsg.ValidatorAddr,
		Pagination:    core.PageRequest,
	}, nil
}

// (Query) make msg - query redelegation
func MakeQueryRedelegationMsg(queryRedelegationMsg types.QueryRedelegationMsg) (stakingtypes.QueryRedelegationsRequest, error) {
	delAddr := queryRedelegationMsg.DelegatorAddr
	valSrcAddr := queryRedelegationMsg.SrcValidatorAddr
	valDstAddr := queryRedelegationMsg.DstValidatorAddr

	return stakingtypes.QueryRedelegationsRequest{
		DelegatorAddr:    delAddr,
		DstValidatorAddr: valDstAddr,
		SrcValidatorAddr: valSrcAddr,
	}, nil
}

// (Query) make msg - query redelegations
func MakeQueryRedelegationsMsg(queryRedelegationMsg types.QueryRedelegationMsg) (stakingtypes.QueryRedelegationsRequest, error) {
	return stakingtypes.QueryRedelegationsRequest{
		DelegatorAddr: queryRedelegationMsg.DelegatorAddr,
		Pagination:    core.PageRequest,
	}, nil
}

// (Query) make msg - query redelegations from
func MakeQueryRedelegationsFromMsg(queryRedelegationMsg types.QueryRedelegationMsg) (stakingtypes.QueryRedelegationsRequest, error) {
	return stakingtypes.QueryRedelegationsRequest{
		SrcValidatorAddr: queryRedelegationMsg.SrcValidatorAddr,
		Pagination:       core.PageRequest,
	}, nil
}

// (Query) make msg - historical
func MakeHistoricalInfoMsg(historicalInfoMsg types.HistoricalInfoMsg) (stakingtypes.QueryHistoricalInfoRequest, error) {
	return parseHistoricalInfoArgs(historicalInfoMsg)
}

// (Query) make msg - staking pool
func MakeQueryStakingPoolMsg() (stakingtypes.QueryPoolRequest, error) {
	return stakingtypes.QueryPoolRequest{}, nil
}

// (Query) make msg - staking params
func MakeQueryStakingParamsMsg() (stakingtypes.QueryParamsRequest, error) {
	return stakingtypes.QueryParamsRequest{}, nil
}
