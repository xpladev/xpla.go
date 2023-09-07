package distribution_test

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/core/distribution"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := distribution.NewCoreModule()

	// test get name
	s.Require().Equal(distribution.DistributionModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// fund community pool
	fundCommunityPoolMsg := types.FundCommunityPoolMsg{
		Amount: "1000",
	}

	makeFundCommunityPoolMsg, err := distribution.MakeFundCommunityPoolMsg(fundCommunityPoolMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeFundCommunityPoolMsg
	txBuilder, err = c.NewTxRouter(txBuilder, distribution.DistributionFundCommunityPoolMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeFundCommunityPoolMsg, txBuilder.GetTx().GetMsgs()[0])

	// community pool spend
	communityPoolSpendMsg := types.CommunityPoolSpendMsg{
		Title:       "community pool spend",
		Description: "pay me",
		Recipient:   accounts[0].Address.String(),
		Amount:      "1000",
		Deposit:     "1000",
	}

	makeProposalCommunityPoolSpendMsg, err := distribution.MakeProposalCommunityPoolSpendMsg(communityPoolSpendMsg, s.xplac.GetPrivateKey(), s.xplac.GetEncoding())
	s.Require().NoError(err)

	testMsg = makeProposalCommunityPoolSpendMsg
	txBuilder, err = c.NewTxRouter(txBuilder, distribution.DistributionProposalCommunityPoolSpendMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeProposalCommunityPoolSpendMsg, txBuilder.GetTx().GetMsgs()[0])

	// withdraw rewards
	withdrawRewardsMsg := types.WithdrawRewardsMsg{
		DelegatorAddr: accounts[0].Address.String(),
		ValidatorAddr: sdk.ValAddress(accounts[0].Address).String(),
		Commission:    true,
	}

	makeWithdrawRewardsMsg, err := distribution.MakeWithdrawRewardsMsg(withdrawRewardsMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeWithdrawRewardsMsg
	txBuilder, err = c.NewTxRouter(txBuilder, distribution.DistributionWithdrawRewardsMsgType, testMsg)
	s.Require().NoError(err)

	// set withdraw address
	setWithdrawAddrMsg := types.SetWithdrawAddrMsg{
		WithdrawAddr: accounts[0].Address.String(),
	}

	makeSetWithdrawAddrMsg, err := distribution.MakeSetWithdrawAddrMsg(setWithdrawAddrMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeSetWithdrawAddrMsg
	txBuilder, err = c.NewTxRouter(txBuilder, distribution.DistributionSetWithdrawAddrMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeSetWithdrawAddrMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = client.ResetXplac(s.xplac)
}
