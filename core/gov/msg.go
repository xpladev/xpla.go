package gov

import (
	"context"

	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
)

// (Tx) make msg - submit proposal
func MakeSubmitProposalMsg(submitProposalMsg types.SubmitProposalMsg, privKey key.PrivateKey) (govtypes.MsgSubmitProposal, error) {
	return parseSubmitProposalArgs(submitProposalMsg, privKey)
}

// (Tx) make msg - deposit
func MakeGovDepositMsg(govDepositMsg types.GovDepositMsg, privKey key.PrivateKey) (govtypes.MsgDeposit, error) {
	return parseGovDepositArgs(govDepositMsg, privKey)
}

// (Tx) make msg - vote
func MakeVoteMsg(voteMsg types.VoteMsg, privKey key.PrivateKey) (govtypes.MsgVote, error) {
	return parseVoteArgs(voteMsg, privKey)
}

// (Tx) make msg - weighted vote
func MakeWeightedVoteMsg(weightedVoteMsg types.WeightedVoteMsg, privKey key.PrivateKey) (govtypes.MsgVoteWeighted, error) {
	return parseWeightedVoteArgs(weightedVoteMsg, privKey)
}

// (Query) make msg - proposal
func MakeQueryProposalMsg(queryProposalMsg types.QueryProposalMsg) (govtypes.QueryProposalRequest, error) {
	proposalId, err := util.FromStringToUint64(queryProposalMsg.ProposalID)
	if err != nil {
		return govtypes.QueryProposalRequest{}, util.LogErr(errors.ErrParse, err)
	}
	return govtypes.QueryProposalRequest{
		ProposalId: proposalId,
	}, nil
}

// (Query) make msg - proposals
func MakeQueryProposalsMsg(queryProposalsMsg types.QueryProposalsMsg) (govtypes.QueryProposalsRequest, error) {
	return parseQueryProposalsArgs(queryProposalsMsg)
}

// (Query) make msg - query deposit
func MakeQueryDepositMsg(queryDepositMsg types.QueryDepositMsg, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	return parseQueryDepositArgs(queryDepositMsg, grpcConn, ctx, lcdUrl, queryType)
}

// (Query) make msg - query deposits
func MakeQueryDepositsMsg(queryDepositMsg types.QueryDepositMsg, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	return parseQueryDepositsArgs(queryDepositMsg, grpcConn, ctx, lcdUrl, queryType)
}

// (Query) make msg - tally
func MakeGovTallyMsg(tallyMsg types.TallyMsg, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, error) {
	return parseGovTallyArgs(tallyMsg, grpcConn, ctx, lcdUrl, queryType)
}

// (Query) make msg - gov params
func MakeGovParamsMsg(govParamsMsg types.GovParamsMsg) (govtypes.QueryParamsRequest, error) {
	return parseGovParamArgs(govParamsMsg)
}

// (Query) make msg - query vote
func MakeQueryVoteMsg(queryVoteMsg types.QueryVoteMsg, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (govtypes.QueryVoteRequest, error) {
	return parseQueryVoteArgs(queryVoteMsg, grpcConn, ctx, lcdUrl, queryType)
}

// (Query) make msg - query votes
func MakeQueryVotesMsg(queryVoteMsg types.QueryVoteMsg, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	return parseQueryVotesArgs(queryVoteMsg, grpcConn, ctx, lcdUrl, queryType)
}
