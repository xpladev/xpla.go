package gov

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type GovExternal struct {
	Xplac provider.XplaClient
}

func NewGovExternal(xplac provider.XplaClient) (e GovExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Submit a proposal along with an initial deposit.
func (e GovExternal) SubmitProposal(submitProposalMsg types.SubmitProposalMsg) provider.XplaClient {
	msg, err := MakeSubmitProposalMsg(submitProposalMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(GovModule).
		WithMsgType(GovSubmitProposalMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Deposit tokens for an active proposal.
func (e GovExternal) GovDeposit(govDepositMsg types.GovDepositMsg) provider.XplaClient {
	msg, err := MakeGovDepositMsg(govDepositMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(GovModule).
		WithMsgType(GovDepositMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Vote for an active proposal, options: yes/no/no_with_veto/abstain.
func (e GovExternal) Vote(voteMsg types.VoteMsg) provider.XplaClient {
	msg, err := MakeVoteMsg(voteMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(GovModule).
		WithMsgType(GovVoteMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Vote for an active proposal, options: yes/no/no_with_veto/abstain.
func (e GovExternal) WeightedVote(weightedVoteMsg types.WeightedVoteMsg) provider.XplaClient {
	msg, err := MakeWeightedVoteMsg(weightedVoteMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(GovModule).
		WithMsgType(GovWeightedVoteMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Query details of a singla proposal.
func (e GovExternal) QueryProposal(queryProposal types.QueryProposalMsg) provider.XplaClient {
	msg, err := MakeQueryProposalMsg(queryProposal)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(GovModule).
		WithMsgType(GovQueryProposalMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query proposals with optional filters.
func (e GovExternal) QueryProposals(queryProposals types.QueryProposalsMsg) provider.XplaClient {
	msg, err := MakeQueryProposalsMsg(queryProposals)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(GovModule).
		WithMsgType(GovQueryProposalsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query details of a deposit or deposits on a proposal.
func (e GovExternal) QueryDeposit(queryDepositMsg types.QueryDepositMsg) provider.XplaClient {
	var queryType int
	if e.Xplac.GetGrpcUrl() != "" {
		queryType = types.QueryGrpc
	} else {
		queryType = types.QueryLcd
	}

	if queryDepositMsg.Depositor != "" {
		msg, argsType, err := MakeQueryDepositMsg(queryDepositMsg, e.Xplac.GetHttpMutex(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext(), e.Xplac.GetLcdURL(), queryType)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		if argsType == "params" {
			e.Xplac.WithModule(GovModule).
				WithMsgType(GovQueryDepositParamsMsgType).
				WithMsg(msg)
		} else {
			e.Xplac.WithModule(GovModule).
				WithMsgType(GovQueryDepositRequestMsgType).
				WithMsg(msg)
		}
	} else {
		msg, argsType, err := MakeQueryDepositsMsg(queryDepositMsg, e.Xplac.GetHttpMutex(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext(), e.Xplac.GetLcdURL(), queryType)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		if argsType == "params" {
			e.Xplac.WithModule(GovModule).
				WithMsgType(GovQueryDepositsParamsMsgType).
				WithMsg(msg)
		} else {
			e.Xplac.WithModule(GovModule).
				WithMsgType(GovQueryDepositsRequestMsgType).
				WithMsg(msg)
		}
	}
	return e.Xplac
}

// Query details of a single vote or votes on a proposal.
func (e GovExternal) QueryVote(queryVoteMsg types.QueryVoteMsg) provider.XplaClient {
	var queryType int
	if e.Xplac.GetGrpcUrl() != "" {
		queryType = types.QueryGrpc
	} else {
		queryType = types.QueryLcd
	}

	if queryVoteMsg.VoterAddr != "" {
		msg, err := MakeQueryVoteMsg(queryVoteMsg, e.Xplac.GetHttpMutex(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext(), e.Xplac.GetLcdURL(), queryType)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(GovModule).
			WithMsgType(GovQueryVoteMsgType).
			WithMsg(msg)

	} else {
		msg, status, err := MakeQueryVotesMsg(queryVoteMsg, e.Xplac.GetHttpMutex(), e.Xplac.GetGrpcClient(), e.Xplac.GetContext(), e.Xplac.GetLcdURL(), queryType)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		if status == "notPassed" {
			e.Xplac.WithModule(GovModule).
				WithMsgType(GovQueryVotesNotPassedMsgType).
				WithMsg(msg)
		} else {
			e.Xplac.WithModule(GovModule).
				WithMsgType(GovQueryVotesPassedMsgType).
				WithMsg(msg)
		}
	}
	return e.Xplac
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
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(GovModule).
		WithMsgType(GovTallyMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query parameters of the governance process or the parameters (voting|tallying|deposit) of the governance process.
func (e GovExternal) GovParams(govParamsMsg ...types.GovParamsMsg) provider.XplaClient {
	if len(govParamsMsg) == 0 {
		e.Xplac.WithModule(GovModule).
			WithMsgType(GovQueryGovParamsMsgType).
			WithMsg(nil)
	} else if len(govParamsMsg) == 1 {
		msg, err := MakeGovParamsMsg(govParamsMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(GovModule)
		switch govParamsMsg[0].ParamType {
		case "voting":
			e.Xplac.WithMsgType(GovQueryGovParamVotingMsgType)
		case "tallying":
			e.Xplac.WithMsgType(GovQueryGovParamTallyingMsgType)
		case "deposit":
			e.Xplac.WithMsgType(GovQueryGovParamDepositMsgType)
		}
		e.Xplac.WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}

// Query the proposer of a governance proposal.
func (e GovExternal) Proposer(proposerMsg types.ProposerMsg) provider.XplaClient {
	e.Xplac.WithModule(GovModule).
		WithMsgType(GovQueryProposerMsgType).
		WithMsg(proposerMsg.ProposalID)
	return e.Xplac
}
