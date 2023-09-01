package distribution

import (
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"golang.org/x/net/context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - fund community pool
func MakeFundCommunityPoolMsg(fundCommunityPoolMsg types.FundCommunityPoolMsg, privKey key.PrivateKey) (disttypes.MsgFundCommunityPool, error) {
	return parseFundCommunityPoolArgs(fundCommunityPoolMsg, privKey)
}

// (Tx) make msg - proposal community pool
func MakeProposalCommunityPoolSpendMsg(communityPoolSpendMsg types.CommunityPoolSpendMsg, privKey key.PrivateKey, encodingConfig params.EncodingConfig) (govtypes.MsgSubmitProposal, error) {
	return parseProposalCommunityPoolSpendArgs(communityPoolSpendMsg, privKey, encodingConfig)
}

// (Tx) make msg - withdraw rewards
func MakeWithdrawRewardsMsg(withdrawRewardsMsg types.WithdrawRewardsMsg, privKey key.PrivateKey) ([]sdk.Msg, error) {
	return parseWithdrawRewardsArgs(withdrawRewardsMsg, privKey)
}

// (Tx) make msg - withdraw all rewards
func MakeWithdrawAllRewardsMsg(privKey key.PrivateKey, grpcConn grpc.ClientConn, ctx context.Context) ([]sdk.Msg, error) {
	return parseWithdrawAllRewardsArgs(privKey, grpcConn, ctx)
}

// (Tx) make msg - withdraw address
func MakeSetWithdrawAddrMsg(setWithdrawAddrMsg types.SetWithdrawAddrMsg, privKey key.PrivateKey) (disttypes.MsgSetWithdrawAddress, error) {
	return parseSetWithdrawAddrArgs(setWithdrawAddrMsg, privKey)
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
