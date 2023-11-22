package distribution

import (
	"context"

	"github.com/xpladev/xpla.go/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - fund community pool
func MakeFundCommunityPoolMsg(fundCommunityPoolMsg types.FundCommunityPoolMsg, depositorAddr sdk.AccAddress) (disttypes.MsgFundCommunityPool, error) {
	return parseFundCommunityPoolArgs(fundCommunityPoolMsg, depositorAddr)
}

// (Tx) make msg - proposal community pool
func MakeProposalCommunityPoolSpendMsg(communityPoolSpendMsg types.CommunityPoolSpendMsg, from sdk.AccAddress, encodingConfig params.EncodingConfig) (govtypes.MsgSubmitProposal, error) {
	return parseProposalCommunityPoolSpendArgs(communityPoolSpendMsg, from, encodingConfig)
}

// (Tx) make msg - withdraw rewards
func MakeWithdrawRewardsMsg(withdrawRewardsMsg types.WithdrawRewardsMsg, delAddr sdk.AccAddress) ([]sdk.Msg, error) {
	return parseWithdrawRewardsArgs(withdrawRewardsMsg, delAddr)
}

// (Tx) make msg - withdraw all rewards
func MakeWithdrawAllRewardsMsg(delAddr sdk.AccAddress, grpcConn grpc.ClientConn, ctx context.Context) ([]sdk.Msg, error) {
	return parseWithdrawAllRewardsArgs(delAddr, grpcConn, ctx)
}

// (Tx) make msg - withdraw address
func MakeSetWithdrawAddrMsg(setWithdrawAddrMsg types.SetWithdrawAddrMsg, delAddr sdk.AccAddress) (disttypes.MsgSetWithdrawAddress, error) {
	return parseSetWithdrawAddrArgs(setWithdrawAddrMsg, delAddr)
}

// (Query) make msg - distribution params
func MakeQueryDistributionParamsMsg() (disttypes.QueryParamsRequest, error) {
	return disttypes.QueryParamsRequest{}, nil
}

// (Query) make msg - validator outstanding rewards
func MakeValidatorOutstandingRewardsMsg(validatorOutstandingRewardsMsg types.ValidatorOutstandingRewardsMsg) (disttypes.QueryValidatorOutstandingRewardsRequest, error) {
	return parseValidatorOutstandingRewardsArgs(validatorOutstandingRewardsMsg)
}

// (Query) make msg - commission
func MakeQueryDistCommissionMsg(queryDistCommissionMsg types.QueryDistCommissionMsg) (disttypes.QueryValidatorCommissionRequest, error) {
	return parseQueryDistCommissionArgs(queryDistCommissionMsg)
}

// (Query) make msg - distribution slashes
func MakeQueryDistSlashesMsg(queryDistSlashesMsg types.QueryDistSlashesMsg) (disttypes.QueryValidatorSlashesRequest, error) {
	return parseDistSlashesArgs(queryDistSlashesMsg)
}

// (Query) make msg - distribution rewards
func MakeQueryDistRewardsMsg(queryDistRewardsMsg types.QueryDistRewardsMsg) (disttypes.QueryDelegationRewardsRequest, error) {
	return parseQueryDistRewardsArgs(queryDistRewardsMsg)
}

// (Query) make msg - distribution all rewards
func MakeQueryDistTotalRewardsMsg(queryDistRewardsMsg types.QueryDistRewardsMsg) (disttypes.QueryDelegationTotalRewardsRequest, error) {
	return disttypes.QueryDelegationTotalRewardsRequest{
		DelegatorAddress: queryDistRewardsMsg.DelegatorAddr,
	}, nil
}

// (Query) make msg - community pool
func MakeQueryCommunityPoolMsg() (disttypes.QueryCommunityPoolRequest, error) {
	return disttypes.QueryCommunityPoolRequest{}, nil
}
