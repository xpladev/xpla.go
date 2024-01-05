package gov

import (
	"context"
	"sync"

	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	govv1beta1 "cosmossdk.io/api/cosmos/gov/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govutils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
)

// Parsing - submit proposal
func parseSubmitProposalArgs(submitProposalMsg types.SubmitProposalMsg, proposer sdk.AccAddress) (govtypes.MsgSubmitProposal, error) {
	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(submitProposalMsg.Deposit))
	if err != nil {
		return govtypes.MsgSubmitProposal{}, types.ErrWrap(types.ErrParse, err)
	}

	content := govtypes.ContentFromProposalType(
		submitProposalMsg.Title,
		submitProposalMsg.Description,
		govutils.NormalizeProposalType(submitProposalMsg.Type),
	)

	msg, err := govtypes.NewMsgSubmitProposal(content, amount, proposer)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, types.ErrWrap(types.ErrInvalidRequest, err)
	}

	return *msg, nil
}

// Parsing - deposit
func parseGovDepositArgs(govDepositMsg types.GovDepositMsg, from sdk.AccAddress) (govtypes.MsgDeposit, error) {
	proposalId, err := util.FromStringToUint64(govDepositMsg.ProposalID)
	if err != nil {
		return govtypes.MsgDeposit{}, types.ErrWrap(types.ErrConvert, err)
	}

	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(govDepositMsg.Deposit))
	if err != nil {
		return govtypes.MsgDeposit{}, types.ErrWrap(types.ErrParse, err)
	}

	msg := govtypes.NewMsgDeposit(from, proposalId, amount)

	return *msg, nil
}

// Parsing - vote
func parseVoteArgs(voteMsg types.VoteMsg, from sdk.AccAddress) (govtypes.MsgVote, error) {
	proposalId, err := util.FromStringToUint64(voteMsg.ProposalID)
	if err != nil {
		return govtypes.MsgVote{}, types.ErrWrap(types.ErrConvert, err)
	}

	byteVoteOption, err := govtypes.VoteOptionFromString(govutils.NormalizeVoteOption(voteMsg.Option))
	if err != nil {
		return govtypes.MsgVote{}, types.ErrWrap(types.ErrParse, err)
	}

	msg := govtypes.NewMsgVote(from, proposalId, byteVoteOption)
	return *msg, nil
}

// Parsing - weighted vote
func parseWeightedVoteArgs(weightedVoteMsg types.WeightedVoteMsg, from sdk.AccAddress) (govtypes.MsgVoteWeighted, error) {
	proposalId, err := util.FromStringToUint64(weightedVoteMsg.ProposalID)
	if err != nil {
		return govtypes.MsgVoteWeighted{}, types.ErrWrap(types.ErrConvert, err)
	}

	options := weightedVoteOptionConverting(weightedVoteMsg)

	msg := govtypes.NewMsgVoteWeighted(from, proposalId, options)
	err = msg.ValidateBasic()
	if err != nil {
		return govtypes.MsgVoteWeighted{}, types.ErrWrap(types.ErrInvalidRequest, err)
	}

	return *msg, nil
}

// Parsing - proposals
func parseQueryProposalsArgs(queryProposalsMsg types.QueryProposalsMsg) (govtypes.QueryProposalsRequest, error) {
	depositorAddr := queryProposalsMsg.Depositor
	voterAddr := queryProposalsMsg.Voter
	strProposalStatus := queryProposalsMsg.Status

	var proposalStatus govtypes.ProposalStatus

	if len(depositorAddr) != 0 {
		_, err := sdk.AccAddressFromBech32(depositorAddr)
		if err != nil {
			return govtypes.QueryProposalsRequest{}, types.ErrWrap(types.ErrParse, err)
		}
	}

	if len(voterAddr) != 0 {
		_, err := sdk.AccAddressFromBech32(voterAddr)
		if err != nil {
			return govtypes.QueryProposalsRequest{}, types.ErrWrap(types.ErrParse, err)
		}
	}

	if len(strProposalStatus) != 0 {
		proposalStatus1, err := govtypes.ProposalStatusFromString(govutils.NormalizeProposalStatus(strProposalStatus))
		proposalStatus = proposalStatus1
		if err != nil {
			return govtypes.QueryProposalsRequest{}, types.ErrWrap(types.ErrParse, err)
		}
	}

	return govtypes.QueryProposalsRequest{
		ProposalStatus: proposalStatus,
		Voter:          voterAddr,
		Depositor:      depositorAddr,
		Pagination:     core.PageRequest,
	}, nil
}

// Parsing - query deposit
func parseQueryDepositArgs(queryDepositMsg types.QueryDepositMsg, httpMutex *sync.Mutex, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	var propStatus govtypes.ProposalStatus

	proposalId, err := util.FromStringToUint64(queryDepositMsg.ProposalID)
	if err != nil {
		return nil, "", types.ErrWrap(types.ErrConvert, err)
	}
	depositorAddr, err := sdk.AccAddressFromBech32(queryDepositMsg.Depositor)
	if err != nil {
		return nil, "", types.ErrWrap(types.ErrParse, err)
	}

	if queryType == types.QueryGrpc {
		queryClient := govtypes.NewQueryClient(grpcConn)

		proposalRes, err := queryClient.Proposal(
			ctx,
			&govtypes.QueryProposalRequest{ProposalId: proposalId},
		)
		if err != nil {
			return nil, "", types.ErrWrap(types.ErrGrpcRequest, err)
		}

		propStatus = proposalRes.Proposal.Status

	} else {
		url := util.MakeQueryLcdUrl(govv1beta1.Query_ServiceDesc.Metadata.(string))
		url = url + util.MakeQueryLabels("proposals", queryDepositMsg.ProposalID)

		httpMutex.Lock()
		out, err := util.CtxHttpClient("GET", lcdUrl+url, nil, ctx)
		if err != nil {
			httpMutex.Unlock()
			return nil, "", err
		}
		httpMutex.Unlock()

		var response govtypes.QueryProposalResponse
		responseData, err := util.JsonUnmarshalData(response, out)
		if err != nil {
			return nil, "", types.ErrWrap(types.ErrFailedToUnmarshal, err)
		}
		propStatusString := responseData.(map[string]interface{})["proposal"].(map[string]interface{})["status"].(string)

		propStatus = govtypes.ProposalStatus(govtypes.ProposalStatus_value[propStatusString])
	}

	if !(propStatus == govtypes.StatusVotingPeriod || propStatus == govtypes.StatusDepositPeriod) {
		params := govtypes.NewQueryDepositParams(proposalId, depositorAddr)
		return params, "params", nil
	}

	return govtypes.QueryDepositRequest{
		ProposalId: proposalId,
		Depositor:  queryDepositMsg.Depositor,
	}, "request", nil
}

// Parsing - query deposits
func parseQueryDepositsArgs(queryDepositMsg types.QueryDepositMsg, httpMutex *sync.Mutex, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	var propStatus govtypes.ProposalStatus
	proposalId, err := util.FromStringToUint64(queryDepositMsg.ProposalID)
	if err != nil {
		return nil, "", types.ErrWrap(types.ErrConvert, err)
	}

	if queryType == types.QueryGrpc {
		queryClient := govtypes.NewQueryClient(grpcConn)

		proposalRes, err := queryClient.Proposal(
			ctx,
			&govtypes.QueryProposalRequest{ProposalId: proposalId},
		)
		if err != nil {
			return nil, "", types.ErrWrap(types.ErrGrpcRequest, err)
		}

		propStatus = proposalRes.Proposal.Status

	} else {
		url := util.MakeQueryLcdUrl(govv1beta1.Query_ServiceDesc.Metadata.(string))
		url = url + util.MakeQueryLabels("proposals", queryDepositMsg.ProposalID)

		httpMutex.Lock()
		out, err := util.CtxHttpClient("GET", lcdUrl+url, nil, ctx)
		if err != nil {
			httpMutex.Unlock()
			return nil, "", err
		}
		httpMutex.Unlock()

		var response govtypes.QueryProposalResponse
		responseData, err := util.JsonUnmarshalData(response, out)
		if err != nil {
			return nil, "", types.ErrWrap(types.ErrFailedToUnmarshal, err)
		}

		propStatusString := responseData.(map[string]interface{})["proposal"].(map[string]interface{})["status"].(string)

		propStatus = govtypes.ProposalStatus(govtypes.ProposalStatus_value[propStatusString])
	}

	if !(propStatus == govtypes.StatusVotingPeriod || propStatus == govtypes.StatusDepositPeriod) {
		params := govtypes.NewQueryProposalParams(proposalId)
		return params, "params", nil
	}

	return govtypes.QueryDepositsRequest{
		ProposalId: proposalId,
		Pagination: core.PageRequest,
	}, "request", nil
}

// Parsing - tally
func parseGovTallyArgs(tallyMsg types.TallyMsg, httpMutex *sync.Mutex, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (govtypes.QueryTallyResultRequest, error) {
	proposalId, err := util.FromStringToUint64(tallyMsg.ProposalID)
	if err != nil {
		return govtypes.QueryTallyResultRequest{}, types.ErrWrap(types.ErrConvert, err)
	}

	if queryType == types.QueryGrpc {
		queryClient := govtypes.NewQueryClient(grpcConn)

		_, err := queryClient.Proposal(
			ctx,
			&govtypes.QueryProposalRequest{ProposalId: proposalId},
		)
		if err != nil {
			return govtypes.QueryTallyResultRequest{}, types.ErrWrap(types.ErrInvalidRequest, "failed to fetch proposal-id", proposalId, " : ", err)
		}

	} else {
		url := util.MakeQueryLcdUrl(govv1beta1.Query_ServiceDesc.Metadata.(string))
		url = url + util.MakeQueryLabels("proposals", tallyMsg.ProposalID)

		httpMutex.Lock()
		_, err := util.CtxHttpClient("GET", lcdUrl+url, nil, ctx)
		if err != nil {
			httpMutex.Unlock()
			return govtypes.QueryTallyResultRequest{}, err
		}
		httpMutex.Unlock()
	}

	return govtypes.QueryTallyResultRequest{
		ProposalId: proposalId,
	}, nil
}

// Parsing - gov params
func parseGovParamArgs(govParamsMsg types.GovParamsMsg) (govtypes.QueryParamsRequest, error) {
	if govParamsMsg.ParamType == "voting" ||
		govParamsMsg.ParamType == "tallying" ||
		govParamsMsg.ParamType == "deposit" {
		return govtypes.QueryParamsRequest{
			ParamsType: govParamsMsg.ParamType,
		}, nil
	} else {
		return govtypes.QueryParamsRequest{}, types.ErrWrap(types.ErrInvalidRequest, "argument must be one of (voting|tallying|deposit), was ", govParamsMsg.ParamType)
	}
}

// Parsing - query vote
func parseQueryVoteArgs(queryVoteMsg types.QueryVoteMsg, httpMutex *sync.Mutex, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (govtypes.QueryVoteRequest, error) {
	proposalId, err := util.FromStringToUint64(queryVoteMsg.ProposalID)
	if err != nil {
		return govtypes.QueryVoteRequest{}, types.ErrWrap(types.ErrConvert, err)
	}

	if queryType == types.QueryGrpc {
		queryClient := govtypes.NewQueryClient(grpcConn)

		_, err := queryClient.Proposal(
			ctx,
			&govtypes.QueryProposalRequest{ProposalId: proposalId},
		)
		if err != nil {
			return govtypes.QueryVoteRequest{}, types.ErrWrap(types.ErrGrpcRequest, err)
		}

	} else {
		url := util.MakeQueryLcdUrl(govv1beta1.Query_ServiceDesc.Metadata.(string))
		url = url + util.MakeQueryLabels("proposals", queryVoteMsg.ProposalID)

		httpMutex.Lock()
		_, err := util.CtxHttpClient("GET", lcdUrl+url, nil, ctx)
		if err != nil {
			httpMutex.Unlock()
			return govtypes.QueryVoteRequest{}, err
		}
		httpMutex.Unlock()
	}

	return govtypes.QueryVoteRequest{
		ProposalId: proposalId,
		Voter:      queryVoteMsg.VoterAddr,
	}, nil
}

// Parsing - query votes
func parseQueryVotesArgs(queryVoteMsg types.QueryVoteMsg, httpMutex *sync.Mutex, grpcConn grpc.ClientConn, ctx context.Context, lcdUrl string, queryType int) (interface{}, string, error) {
	var propStatus govtypes.ProposalStatus
	proposalId, err := util.FromStringToUint64(queryVoteMsg.ProposalID)
	if err != nil {
		return nil, "", types.ErrWrap(types.ErrConvert, err)
	}

	if queryType == types.QueryGrpc {
		queryClient := govtypes.NewQueryClient(grpcConn)

		proposalRes, err := queryClient.Proposal(
			ctx,
			&govtypes.QueryProposalRequest{ProposalId: proposalId},
		)
		if err != nil {
			return nil, "", types.ErrWrap(types.ErrGrpcRequest, err)
		}

		propStatus = proposalRes.Proposal.Status

	} else {
		url := util.MakeQueryLcdUrl(govv1beta1.Query_ServiceDesc.Metadata.(string))
		url = url + util.MakeQueryLabels("proposals", queryVoteMsg.ProposalID)

		httpMutex.Lock()
		out, err := util.CtxHttpClient("GET", lcdUrl+url, nil, ctx)
		if err != nil {
			httpMutex.Unlock()
			return nil, "", err
		}
		httpMutex.Unlock()

		var response govtypes.QueryProposalResponse
		responseData, err := util.JsonUnmarshalData(response, out)
		if err != nil {
			return nil, "", types.ErrWrap(types.ErrFailedToUnmarshal, err)
		}

		propStatusString := responseData.(map[string]interface{})["proposal"].(map[string]interface{})["status"].(string)

		propStatus = govtypes.ProposalStatus(govtypes.ProposalStatus_value[propStatusString])
	}

	if !(propStatus == govtypes.StatusVotingPeriod || propStatus == govtypes.StatusDepositPeriod) {
		params := govtypes.NewQueryProposalVotesParams(proposalId, 0, 0)
		return params, "notPassed", nil
	}

	return govtypes.QueryVotesRequest{
		ProposalId: proposalId,
		Pagination: core.PageRequest,
	}, "passed", nil

}

func weightedVoteOptionConverting(weightedVoteMsg types.WeightedVoteMsg) govtypes.WeightedVoteOptions {
	weightedVoteOptions := govtypes.WeightedVoteOptions{}

	if weightedVoteMsg.Yes != "" {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionYes,
			Weight: sdk.MustNewDecFromStr(weightedVoteMsg.Yes),
		})
	}
	if weightedVoteMsg.Abstain != "" {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionAbstain,
			Weight: sdk.MustNewDecFromStr(weightedVoteMsg.Abstain),
		})
	}
	if weightedVoteMsg.No != "" {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionNo,
			Weight: sdk.MustNewDecFromStr(weightedVoteMsg.No),
		})
	}
	if weightedVoteMsg.NoWithVeto != "" {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionNoWithVeto,
			Weight: sdk.MustNewDecFromStr(weightedVoteMsg.NoWithVeto),
		})
	}

	return weightedVoteOptions
}
