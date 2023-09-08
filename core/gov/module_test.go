package gov_test

import (
	"math/rand"

	"github.com/xpladev/xpla.go/core/gov"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := gov.NewCoreModule()

	// test get name
	s.Require().Equal(gov.GovModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// submit proposal
	submitProposalMsg := types.SubmitProposalMsg{
		Title:       "Test proposal",
		Description: "Proposal description",
		Type:        "text",
		Deposit:     "1000",
	}

	makeSubmitProposalMsg, err := gov.MakeSubmitProposalMsg(submitProposalMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeSubmitProposalMsg
	txBuilder, err = c.NewTxRouter(txBuilder, gov.GovSubmitProposalMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeSubmitProposalMsg, txBuilder.GetTx().GetMsgs()[0])

	// deposit
	govDepositMsg := types.GovDepositMsg{
		ProposalID: "1",
		Deposit:    "1000",
	}

	makeGovDepositMsg, err := gov.MakeGovDepositMsg(govDepositMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeGovDepositMsg
	txBuilder, err = c.NewTxRouter(txBuilder, gov.GovDepositMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeGovDepositMsg, txBuilder.GetTx().GetMsgs()[0])

	// vote
	voteMsg := types.VoteMsg{
		ProposalID: "1",
		Option:     "yes",
	}

	makeVoteMsg, err := gov.MakeVoteMsg(voteMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeVoteMsg
	txBuilder, err = c.NewTxRouter(txBuilder, gov.GovVoteMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeVoteMsg, txBuilder.GetTx().GetMsgs()[0])

	// weighted vote
	weightedVoteMsg := types.WeightedVoteMsg{
		ProposalID: "1",
		Yes:        "0.6",
		No:         "0.3",
		Abstain:    "0.05",
		NoWithVeto: "0.05",
	}

	makeWeightedVoteMsg, err := gov.MakeWeightedVoteMsg(weightedVoteMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeWeightedVoteMsg
	txBuilder, err = c.NewTxRouter(txBuilder, gov.GovWeightedVoteMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeWeightedVoteMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = provider.ResetXplac(s.xplac)
}
