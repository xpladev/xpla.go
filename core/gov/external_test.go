package gov_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mgov "github.com/xpladev/xpla.go/core/gov"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestGovTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// submit proposal
	submitProposalMsg := types.SubmitProposalMsg{
		Title:       "Test proposal",
		Description: "Proposal description",
		Type:        "text",
		Deposit:     "1000",
	}
	s.xplac.SubmitProposal(submitProposalMsg)

	makeSubmitProposalMsg, err := mgov.MakeSubmitProposalMsg(submitProposalMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeSubmitProposalMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovSubmitProposalMsgType, s.xplac.GetMsgType())

	govSubmitProposalTxbytes, err := s.xplac.SubmitProposal(submitProposalMsg).CreateAndSignTx()
	s.Require().NoError(err)

	govSubmitProposalJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(govSubmitProposalTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.GovSubmitProposalTxTemplates, string(govSubmitProposalJsonTxbytes))

	// deposit
	govDepositMsg := types.GovDepositMsg{
		ProposalID: "1",
		Deposit:    "1000",
	}
	s.xplac.GovDeposit(govDepositMsg)

	makeGovDepositMsg, err := mgov.MakeGovDepositMsg(govDepositMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeGovDepositMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovDepositMsgType, s.xplac.GetMsgType())

	govDepositTxbytes, err := s.xplac.GovDeposit(govDepositMsg).CreateAndSignTx()
	s.Require().NoError(err)

	govDepositJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(govDepositTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.GovDepositTxTemplates, string(govDepositJsonTxbytes))

	// vote
	voteMsg := types.VoteMsg{
		ProposalID: "1",
		Option:     "yes",
	}
	s.xplac.Vote(voteMsg)

	makeVoteMsg, err := mgov.MakeVoteMsg(voteMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeVoteMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovVoteMsgType, s.xplac.GetMsgType())

	govVoteTxbytes, err := s.xplac.Vote(voteMsg).CreateAndSignTx()
	s.Require().NoError(err)

	govVoteJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(govVoteTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.GovVoteTxTemplates, string(govVoteJsonTxbytes))

	// weighted vote
	weightedVoteMsg := types.WeightedVoteMsg{
		ProposalID: "1",
		Yes:        "0.6",
		No:         "0.3",
		Abstain:    "0.05",
		NoWithVeto: "0.05",
	}
	s.xplac.WeightedVote(weightedVoteMsg)

	makeWeightedVoteMsg, err := mgov.MakeWeightedVoteMsg(weightedVoteMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeWeightedVoteMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovWeightedVoteMsgType, s.xplac.GetMsgType())

	govWeightedVoteTxbytes, err := s.xplac.WeightedVote(weightedVoteMsg).CreateAndSignTx()
	s.Require().NoError(err)

	govWeightedVoteJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(govWeightedVoteTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.GovWeightedVoteTxTemplates, string(govWeightedVoteJsonTxbytes))
}

func (s *IntegrationTestSuite) TestGov() {
	val := s.network.Validators[0]

	_, err := MsgSubmitProposal(val.ClientCtx, val.Address.String(),
		"Text Proposal 1", "Where is the title!?", govtypes.ProposalTypeText,
		fmt.Sprintf("--%s=%s", govcli.FlagDeposit, sdk.NewCoin(s.cfg.BondDenom, govtypes.DefaultMinDepositTokens).String()))
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	_, err = MsgVote(val.ClientCtx, val.Address.String(), "1", "yes=0.6,no=0.3,abstain=0.05,no_with_veto=0.05")
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// query proposal
	queryProposalMsg := types.QueryProposalMsg{
		ProposalID: "1",
	}
	s.xplac.QueryProposal(queryProposalMsg)

	makeQueryProposalMsg, err := mgov.MakeQueryProposalMsg(queryProposalMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryProposalMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryProposalMsgType, s.xplac.GetMsgType())

	// query proposals
	queryProposalsMsg := types.QueryProposalsMsg{
		Status:    "DepositPeriod",
		Voter:     s.accounts[0].Address.String(),
		Depositor: s.accounts[1].Address.String(),
	}
	s.xplac.QueryProposals(queryProposalsMsg)

	makeQueryProposalsMsg, err := mgov.MakeQueryProposalsMsg(queryProposalsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryProposalsMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryProposalsMsgType, s.xplac.GetMsgType())

	var queryType int
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
			queryType = types.QueryLcd
		} else {
			s.xplac.WithGrpc(api)
			queryType = types.QueryGrpc
		}

		// query deposit
		queryDepositMsg := types.QueryDepositMsg{
			ProposalID: "1",
			Depositor:  s.accounts[0].Address.String(),
		}
		s.xplac.QueryDeposit(queryDepositMsg)

		makeQueryDepositMsg, _, err := mgov.MakeQueryDepositMsg(queryDepositMsg, s.xplac.GetGrpcClient(), s.xplac.GetContext(), s.xplac.GetLcdURL(), queryType)
		s.Require().NoError(err)

		s.Require().Equal(makeQueryDepositMsg, s.xplac.GetMsg())
		s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
		s.Require().Equal(mgov.GovQueryDepositRequestMsgType, s.xplac.GetMsgType())

		// query deposits
		queryDepositMsg = types.QueryDepositMsg{
			ProposalID: "1",
		}
		s.xplac.QueryDeposit(queryDepositMsg)

		makeQueryDepositsMsg, _, err := mgov.MakeQueryDepositsMsg(queryDepositMsg, s.xplac.GetGrpcClient(), s.xplac.GetContext(), s.xplac.GetLcdURL(), queryType)
		s.Require().NoError(err)

		s.Require().Equal(makeQueryDepositsMsg, s.xplac.GetMsg())
		s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
		s.Require().Equal(mgov.GovQueryDepositsRequestMsgType, s.xplac.GetMsgType())

		// query vote
		queryVoteMsg := types.QueryVoteMsg{
			ProposalID: "1",
			VoterAddr:  val.Address.String(),
		}
		s.xplac.QueryVote(queryVoteMsg)

		makeQueryVoteMsg, err := mgov.MakeQueryVoteMsg(queryVoteMsg, s.xplac.GetGrpcClient(), s.xplac.GetContext(), s.xplac.GetLcdURL(), queryType)
		s.Require().NoError(err)

		s.Require().Equal(makeQueryVoteMsg, s.xplac.GetMsg())
		s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
		s.Require().Equal(mgov.GovQueryVoteMsgType, s.xplac.GetMsgType())

		// query votes
		queryVoteMsg = types.QueryVoteMsg{
			ProposalID: "1",
		}
		s.xplac.QueryVote(queryVoteMsg)

		makeQueryVotesMsg, _, err := mgov.MakeQueryVotesMsg(queryVoteMsg, s.xplac.GetGrpcClient(), s.xplac.GetContext(), s.xplac.GetLcdURL(), queryType)
		s.Require().NoError(err)

		s.Require().Equal(makeQueryVotesMsg, s.xplac.GetMsg())
		s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
		s.Require().Equal(mgov.GovQueryVotesPassedMsgType, s.xplac.GetMsgType())

		// tally
		tallyMsg := types.TallyMsg{
			ProposalID: "1",
		}
		s.xplac.Tally(tallyMsg)

		makeGovTallyMsg, err := mgov.MakeGovTallyMsg(tallyMsg, s.xplac.GetGrpcClient(), s.xplac.GetContext(), s.xplac.GetLcdURL(), queryType)
		s.Require().NoError(err)

		s.Require().Equal(makeGovTallyMsg, s.xplac.GetMsg())
		s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
		s.Require().Equal(mgov.GovTallyMsgType, s.xplac.GetMsgType())
	}
	s.xplac = provider.ResetXplac(s.xplac)

	// gov params
	s.xplac.GovParams()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryGovParamsMsgType, s.xplac.GetMsgType())

	// gov params, paramtype voting
	govParamsMsg := types.GovParamsMsg{
		ParamType: "voting",
	}
	s.xplac.GovParams(govParamsMsg)

	makeGovParamsMsg, err := mgov.MakeGovParamsMsg(govParamsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGovParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryGovParamVotingMsgType, s.xplac.GetMsgType())

	// gov params, paramtype tallying
	govParamsMsg = types.GovParamsMsg{
		ParamType: "tallying",
	}
	s.xplac.GovParams(govParamsMsg)

	makeGovParamsMsg, err = mgov.MakeGovParamsMsg(govParamsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGovParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryGovParamTallyingMsgType, s.xplac.GetMsgType())

	// gov params, paramtype deposit
	govParamsMsg = types.GovParamsMsg{
		ParamType: "deposit",
	}
	s.xplac.GovParams(govParamsMsg)

	makeGovParamsMsg, err = mgov.MakeGovParamsMsg(govParamsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGovParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryGovParamDepositMsgType, s.xplac.GetMsgType())

	// proposer
	proposerMsg := types.ProposerMsg{
		ProposalID: "1",
	}
	s.xplac.Proposer(proposerMsg)

	s.Require().Equal(proposerMsg.ProposalID, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryProposerMsgType, s.xplac.GetMsgType())
}
