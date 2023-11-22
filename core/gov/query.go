package gov

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	govv1beta1 "cosmossdk.io/api/cosmos/gov/v1beta1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govutils "github.com/cosmos/cosmos-sdk/x/gov/client/utils"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var out []byte
var res proto.Message
var err error

// Query client for gov module.
func QueryGov(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcGov(i)
	} else {
		return queryByLcdGov(i)
	}
}

func queryByGrpcGov(i core.QueryClient) (string, error) {
	queryClient := govtypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Gov proposal
	case i.Ixplac.GetMsgType() == GovQueryProposalMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryProposalRequest)
		res, err = queryClient.Proposal(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Gov proposals
	case i.Ixplac.GetMsgType() == GovQueryProposalsMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryProposalsRequest)
		res, err = queryClient.Proposals(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Gov deposit parameter
	case i.Ixplac.GetMsgType() == GovQueryDepositParamsMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryDepositParams)

		var deposit govtypes.Deposit

		clientCtx, err := core.ClientForQuery(i)
		if err != nil {
			return "", err
		}

		resByTxQuery, err := govutils.QueryDepositByTxQuery(clientCtx, convertMsg)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}
		clientCtx.Codec.MustUnmarshalJSON(resByTxQuery, &deposit)
		res = &deposit

	// Gov deposit
	case i.Ixplac.GetMsgType() == GovQueryDepositRequestMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryDepositRequest)
		res, err = queryClient.Deposit(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Gov deposits parameter
	case i.Ixplac.GetMsgType() == GovQueryDepositsParamsMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryProposalParams)

		var deposit govtypes.Deposits
		clientCtx, err := core.ClientForQuery(i)
		if err != nil {
			return "", err
		}

		resByTxQuery, err := govutils.QueryDepositsByTxQuery(clientCtx, convertMsg)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		clientCtx.LegacyAmino.MustUnmarshalJSON(resByTxQuery, &deposit)
		out, err := core.PrintObjectLegacy(i, deposit)
		if err != nil {
			return "", util.LogErr(errors.ErrParse, err)
		}
		return string(out), nil

	// Gov deposits
	case i.Ixplac.GetMsgType() == GovQueryDepositsRequestMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryDepositsRequest)
		res, err = queryClient.Deposits(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Gov tally
	case i.Ixplac.GetMsgType() == GovTallyMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryTallyResultRequest)
		res, err = queryClient.TallyResult(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Gov params
	case i.Ixplac.GetMsgType() == GovQueryGovParamsMsgType:
		votingRes, err := queryClient.Params(
			i.Ixplac.GetContext(),
			&govtypes.QueryParamsRequest{ParamsType: "voting"},
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		tallyRes, err := queryClient.Params(
			i.Ixplac.GetContext(),
			&govtypes.QueryParamsRequest{ParamsType: "tallying"},
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		depositRes, err := queryClient.Params(
			i.Ixplac.GetContext(),
			&govtypes.QueryParamsRequest{ParamsType: "deposit"},
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		govAllParams := govtypes.NewParams(
			votingRes.GetVotingParams(),
			tallyRes.GetTallyParams(),
			depositRes.GetDepositParams(),
		)

		bytes, err := util.JsonMarshalData(govAllParams)
		if err != nil {
			return "", util.LogErr(errors.ErrFailedToMarshal, err)
		}
		return string(bytes), nil

	// Gov params of voting
	case i.Ixplac.GetMsgType() == GovQueryGovParamVotingMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)
		resParams, err := queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		bytes, err := util.JsonMarshalData(resParams.GetVotingParams())
		if err != nil {
			return "", util.LogErr(errors.ErrFailedToMarshal, err)
		}
		return string(bytes), nil

	// Gov params of tally
	case i.Ixplac.GetMsgType() == GovQueryGovParamTallyingMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)
		resParams, err := queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		bytes, err := util.JsonMarshalData(resParams.GetTallyParams())
		if err != nil {
			return "", util.LogErr(errors.ErrFailedToMarshal, err)
		}
		return string(bytes), nil

	// Gov params of deposit
	case i.Ixplac.GetMsgType() == GovQueryGovParamDepositMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)
		resParams, err := queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		bytes, err := util.JsonMarshalData(resParams.GetDepositParams())
		if err != nil {
			return "", util.LogErr(errors.ErrFailedToMarshal, err)
		}
		return string(bytes), nil

	// Gov proposer
	case i.Ixplac.GetMsgType() == GovQueryProposerMsgType:
		convertMsg := i.Ixplac.GetMsg().(string)
		proposalId, err := util.FromStringToUint64(convertMsg)
		if err != nil {
			return "", err
		}

		clientCtx, err := core.ClientForQuery(i)
		if err != nil {
			return "", err
		}

		prop, err := govutils.QueryProposerByTxQuery(clientCtx, proposalId)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		bytes, err := util.JsonMarshalData(prop)
		if err != nil {
			return "", util.LogErr(errors.ErrFailedToMarshal, err)
		}
		return string(bytes), nil

	// Gov vote
	case i.Ixplac.GetMsgType() == GovQueryVoteMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryVoteRequest)
		resVote, err := queryClient.Vote(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		clientCtx, err := core.ClientForQuery(i)
		if err != nil {
			return "", err
		}

		voterAddr, err := sdk.AccAddressFromBech32(convertMsg.Voter)
		if err != nil {
			return "", util.LogErr(errors.ErrParse, err)
		}

		vote := resVote.GetVote()
		if vote.Empty() {
			params := govtypes.NewQueryVoteParams(convertMsg.ProposalId, voterAddr)
			resByTxQuery, err := govutils.QueryVoteByTxQuery(clientCtx, params)
			if err != nil {
				return "", util.LogErr(errors.ErrGrpcRequest, err)
			}

			if err := clientCtx.Codec.UnmarshalJSON(resByTxQuery, &vote); err != nil {
				return "", util.LogErr(errors.ErrFailedToUnmarshal, err)
			}
		}

		res = &resVote.Vote

	// Gov votes not passed
	case i.Ixplac.GetMsgType() == GovQueryVotesNotPassedMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryProposalVotesParams)
		clientCtx, err := core.ClientForQuery(i)
		if err != nil {
			return "", err
		}
		resByTxQuery, err := govutils.QueryVotesByTxQuery(clientCtx, convertMsg)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		var votes govtypes.Votes

		clientCtx.LegacyAmino.MustUnmarshalJSON(resByTxQuery, &votes)
		out, err := core.PrintObjectLegacy(i, votes)
		if err != nil {
			return "", util.LogErr(errors.ErrParse, err)
		}
		return string(out), nil

	// Gov votes passed
	case i.Ixplac.GetMsgType() == GovQueryVotesPassedMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryVotesRequest)
		res, err = queryClient.Votes(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = core.PrintProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	govProposalsLabel = "proposals"
	govDeposistsLabel = "deposits"
	govTallyLabel     = "tally"
	govParamsLabel    = "params"
	govVotesLabel     = "votes"
)

func queryByLcdGov(i core.QueryClient) (string, error) {
	url := util.MakeQueryLcdUrl(govv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Gov proposal
	case i.Ixplac.GetMsgType() == GovQueryProposalMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryProposalRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId))

	// Gov proposals
	case i.Ixplac.GetMsgType() == GovQueryProposalsMsgType:
		url = url + govProposalsLabel

	// Gov deposit parameter
	case i.Ixplac.GetMsgType() == GovQueryDepositParamsMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryDepositParams)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalID), govDeposistsLabel, convertMsg.Depositor.String())

	// Gov deposit
	case i.Ixplac.GetMsgType() == GovQueryDepositRequestMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryDepositRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId), govDeposistsLabel, convertMsg.Depositor)

	// Gov deposits parameter
	case i.Ixplac.GetMsgType() == GovQueryDepositsParamsMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryProposalParams)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalID), govDeposistsLabel)

	// Gov deposits
	case i.Ixplac.GetMsgType() == GovQueryDepositsRequestMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryDepositsRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId), govDeposistsLabel)

	// Gov tally
	case i.Ixplac.GetMsgType() == GovTallyMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryTallyResultRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId), govTallyLabel)

	// Gov params
	case i.Ixplac.GetMsgType() == GovQueryGovParamsMsgType:
		return "", util.LogErr(errors.ErrNotSupport, "unsupported querying all gov params by using LCD. query each parameter(voting|tallying|deposit)")

	// Gov params of voting
	case i.Ixplac.GetMsgType() == GovQueryGovParamVotingMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)

		url = url + util.MakeQueryLabels(govParamsLabel, convertMsg.ParamsType)

	// Gov params of tally
	case i.Ixplac.GetMsgType() == GovQueryGovParamTallyingMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)

		url = url + util.MakeQueryLabels(govParamsLabel, convertMsg.ParamsType)

	// Gov params of deposit
	case i.Ixplac.GetMsgType() == GovQueryGovParamDepositMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryParamsRequest)

		url = url + util.MakeQueryLabels(govParamsLabel, convertMsg.ParamsType)

	// Gov proposer
	case i.Ixplac.GetMsgType() == GovQueryProposerMsgType:
		return "", util.LogErr(errors.ErrNotSupport, "unsupported querying proposer by using LCD")

	// Gov vote
	case i.Ixplac.GetMsgType() == GovQueryVoteMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryVoteRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId), govVotesLabel, convertMsg.Voter)

	// Gov votes not passed
	case i.Ixplac.GetMsgType() == GovQueryVotesNotPassedMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryProposalVotesParams)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalID), govVotesLabel)

	// Gov votes passed
	case i.Ixplac.GetMsgType() == GovQueryVotesPassedMsgType:
		convertMsg := i.Ixplac.GetMsg().(govtypes.QueryVotesRequest)

		url = url + util.MakeQueryLabels(govProposalsLabel, util.FromUint64ToString(convertMsg.ProposalId), govVotesLabel)

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	i.Ixplac.GetHttpMutex().Lock()
	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		i.Ixplac.GetHttpMutex().Unlock()
		return "", err
	}
	i.Ixplac.GetHttpMutex().Unlock()

	return string(out), nil

}
