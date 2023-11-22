package gov

import (
	"context"
	"sync"

	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
)

// (Tx) make msg - submit proposal
func MakeSubmitProposalMsg(submitProposalMsg types.SubmitProposalMsg, proposer sdk.AccAddress) (govtypes.MsgSubmitProposal, error) {
	return parseSubmitProposalArgs(submitProposalMsg, proposer)
}

// (Tx) make msg - deposit
func MakeGovDepositMsg(govDepositMsg types.GovDepositMsg, from sdk.AccAddress) (govtypes.MsgDeposit, error) {
	return parseGovDepositArgs(govDepositMsg, from)
}

// (Tx) make msg - vote
func MakeVoteMsg(voteMsg types.VoteMsg, from sdk.AccAddress) (govtypes.MsgVote, error) {
	return parseVoteArgs(voteMsg, from)
}

// (Tx) make msg - weighted vote
func MakeWeightedVoteMsg(weightedVoteMsg types.WeightedVoteMsg, from sdk.AccAddress) (govtypes.MsgVoteWeighted, error) {
	return parseWeightedVoteArgs(weightedVoteMsg, from)
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
func MakeQueryDepositMsg(queryDepositMsg types.QueryDepositMsg, httpMutex *sync.Mutex, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	return parseQueryDepositArgs(queryDepositMsg, httpMutex, grpcConn, ctx, lcdUrl, queryType)
}

// (Query) make msg - query deposits
func MakeQueryDepositsMsg(queryDepositMsg types.QueryDepositMsg, httpMutex *sync.Mutex, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	return parseQueryDepositsArgs(queryDepositMsg, httpMutex, grpcConn, ctx, lcdUrl, queryType)
}

// (Query) make msg - tally
func MakeGovTallyMsg(tallyMsg types.TallyMsg, httpMutex *sync.Mutex, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, error) {
	return parseGovTallyArgs(tallyMsg, httpMutex, grpcConn, ctx, lcdUrl, queryType)
}

// (Query) make msg - gov params
func MakeGovParamsMsg(govParamsMsg types.GovParamsMsg) (govtypes.QueryParamsRequest, error) {
	return parseGovParamArgs(govParamsMsg)
}

// (Query) make msg - query vote
func MakeQueryVoteMsg(queryVoteMsg types.QueryVoteMsg, httpMutex *sync.Mutex, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (govtypes.QueryVoteRequest, error) {
	return parseQueryVoteArgs(queryVoteMsg, httpMutex, grpcConn, ctx, lcdUrl, queryType)
}

// (Query) make msg - query votes
func MakeQueryVotesMsg(queryVoteMsg types.QueryVoteMsg, httpMutex *sync.Mutex, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	return parseQueryVotesArgs(queryVoteMsg, httpMutex, grpcConn, ctx, lcdUrl, queryType)
}
