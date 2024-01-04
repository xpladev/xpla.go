package gov

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &GovExternal{}

type GovExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e GovExternal) {
	e.Xplac = xplac
	e.Name = GovModule
	return e
}

func (e GovExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e GovExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Submit a proposal along with an initial deposit.
func (e GovExternal) SubmitProposal(submitProposalMsg types.SubmitProposalMsg) provider.XplaClient {
	msg, err := MakeSubmitProposalMsg(submitProposalMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(GovSubmitProposalMsgType, err)
	}

	return e.ToExternal(GovSubmitProposalMsgType, msg)
}

// Deposit tokens for an active proposal.
func (e GovExternal) GovDeposit(govDepositMsg types.GovDepositMsg) provider.XplaClient {
	msg, err := MakeGovDepositMsg(govDepositMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(GovDepositMsgType, err)
	}

	return e.ToExternal(GovDepositMsgType, msg)
}

// Vote for an active proposal, options: yes/no/no_with_veto/abstain.
func (e GovExternal) Vote(voteMsg types.VoteMsg) provider.XplaClient {
	msg, err := MakeVoteMsg(voteMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(GovVoteMsgType, err)
	}

	return e.ToExternal(GovVoteMsgType, msg)
}

// Vote for an active proposal, options: yes/no/no_with_veto/abstain.
func (e GovExternal) WeightedVote(weightedVoteMsg types.WeightedVoteMsg) provider.XplaClient {
	msg, err := MakeWeightedVoteMsg(weightedVoteMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(GovWeightedVoteMsgType, err)
	}

	return e.ToExternal(GovWeightedVoteMsgType, msg)
}

// Query

// Query details of a singla proposal.
func (e GovExternal) QueryProposal(queryProposal types.QueryProposalMsg) provider.XplaClient {
	msg, err := MakeQueryProposalMsg(queryProposal)
	if err != nil {
		return e.Err(GovQueryProposalMsgType, err)
	}

	return e.ToExternal(GovQueryProposalMsgType, msg)
}

// Query proposals with optional filters.
func (e GovExternal) QueryProposals(queryProposals types.QueryProposalsMsg) provider.XplaClient {
	msg, err := MakeQueryProposalsMsg(queryProposals)
	if err != nil {
		return e.Err(GovQueryProposalsMsgType, err)
	}

	return e.ToExternal(GovQueryProposalsMsgType, msg)
}

// Query details of a deposit or deposits on a proposal.
func (e GovExternal) QueryDeposit(queryDepositMsg types.QueryDepositMsg) provider.XplaClient {
	var queryType int
	if e.Xplac.GetGrpcUrl() != "" {
		queryType = types.QueryGrpc
	} else {
		queryType = types.QueryLcd
	}

	switch {
	case queryDepositMsg.Depositor != "":
		msg, argsType, err := MakeQueryDepositMsg(queryDepositMsg, e.Xplac.GetHttpMutex(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext(), e.Xplac.GetLcdURL(), queryType)
		if err != nil {
			return e.Err(GovQueryDepositRequestMsgType, err)
		}
		if argsType == "params" {
			return e.ToExternal(GovQueryDepositParamsMsgType, msg)

		} else {
			return e.ToExternal(GovQueryDepositRequestMsgType, msg)
		}

	default:
		msg, argsType, err := MakeQueryDepositsMsg(queryDepositMsg, e.Xplac.GetHttpMutex(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext(), e.Xplac.GetLcdURL(), queryType)
		if err != nil {
			return e.Err(GovQueryDepositsRequestMsgType, err)
		}
		if argsType == "params" {
			return e.ToExternal(GovQueryDepositsParamsMsgType, msg)

		} else {
			return e.ToExternal(GovQueryDepositsRequestMsgType, msg)
		}
	}
}

// Query details of a single vote or votes on a proposal.
func (e GovExternal) QueryVote(queryVoteMsg types.QueryVoteMsg) provider.XplaClient {
	var queryType int
	if e.Xplac.GetGrpcUrl() != "" {
		queryType = types.QueryGrpc
	} else {
		queryType = types.QueryLcd
	}

	switch {
	case queryVoteMsg.VoterAddr != "":
		msg, err := MakeQueryVoteMsg(queryVoteMsg, e.Xplac.GetHttpMutex(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext(), e.Xplac.GetLcdURL(), queryType)
		if err != nil {
			return e.Err(GovQueryVoteMsgType, err)
		}

		return e.ToExternal(GovQueryVoteMsgType, msg)

	default:
		msg, status, err := MakeQueryVotesMsg(queryVoteMsg, e.Xplac.GetHttpMutex(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext(), e.Xplac.GetLcdURL(), queryType)
		if err != nil {
			return e.Err(GovQueryVotesPassedMsgType, err)
		}
		if status == "notPassed" {
			return e.ToExternal(GovQueryVotesNotPassedMsgType, msg)

		} else {
			return e.ToExternal(GovQueryVotesPassedMsgType, msg)
		}
	}
}

// Query the tally of a proposal vote.
func (e GovExternal) Tally(tallyMsg types.TallyMsg) provider.XplaClient {
	var queryType int
	if e.Xplac.GetGrpcUrl() != "" {
		queryType = types.QueryGrpc
	} else {
		queryType = types.QueryLcd
	}

	msg, err := MakeGovTallyMsg(tallyMsg, e.Xplac.GetHttpMutex(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext(), e.Xplac.GetLcdURL(), queryType)
	if err != nil {
		return e.Err(GovTallyMsgType, err)
	}

	return e.ToExternal(GovTallyMsgType, msg)
}

// Query parameters of the governance process or the parameters (voting|tallying|deposit) of the governance process.
func (e GovExternal) GovParams(govParamsMsg ...types.GovParamsMsg) provider.XplaClient {
	switch {
	case len(govParamsMsg) == 0:
		return e.ToExternal(GovQueryGovParamsMsgType, nil)

	case len(govParamsMsg) == 1:
		msg, err := MakeGovParamsMsg(govParamsMsg[0])
		if err != nil {
			return e.Err(GovQueryGovParamsMsgType, err)
		}

		switch govParamsMsg[0].ParamType {
		case "voting":
			return e.ToExternal(GovQueryGovParamVotingMsgType, msg)
		case "tallying":
			return e.ToExternal(GovQueryGovParamTallyingMsgType, msg)
		case "deposit":
			return e.ToExternal(GovQueryGovParamDepositMsgType, msg)
		default:
			return e.Err(GovQueryGovParamsMsgType, types.ErrWrap(types.ErrInvalidRequest, "invalid param type"))
		}

	default:
		return e.Err(GovQueryGovParamsMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}

// Query the proposer of a governance proposal.
func (e GovExternal) Proposer(proposerMsg types.ProposerMsg) provider.XplaClient {
	return e.ToExternal(GovQueryProposerMsgType, proposerMsg.ProposalID)
}
