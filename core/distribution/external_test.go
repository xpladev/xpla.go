package distribution_test

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mdist "github.com/xpladev/xpla.go/core/distribution"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestDistributionTx() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)

	s.xplac.WithPrivateKey(accounts[0].PrivKey)
	// fund community pool
	fundCommunityPoolMsg := types.FundCommunityPoolMsg{
		Amount: "1000",
	}
	s.xplac.FundCommunityPool(fundCommunityPoolMsg)

	makeFundCommunityPoolMsg, err := mdist.MakeFundCommunityPoolMsg(fundCommunityPoolMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeFundCommunityPoolMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionFundCommunityPoolMsgType, s.xplac.GetMsgType())

	distFundCommunityPoolTxbytes, err := s.xplac.FundCommunityPool(fundCommunityPoolMsg).CreateAndSignTx()
	s.Require().NoError(err)

	distFundCommunityPoolJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(distFundCommunityPoolTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.DistFundCommunityPoolTxTemplates, string(distFundCommunityPoolJsonTxbytes))

	// community pool spend
	communityPoolSpendMsg := types.CommunityPoolSpendMsg{
		Title:       "community pool spend",
		Description: "pay me",
		Recipient:   accounts[0].Address.String(),
		Amount:      "1000",
		Deposit:     "1000",
	}
	s.xplac.CommunityPoolSpend(communityPoolSpendMsg)

	makeProposalCommunityPoolSpendMsg, err := mdist.MakeProposalCommunityPoolSpendMsg(communityPoolSpendMsg, s.xplac.GetPrivateKey(), s.xplac.GetEncoding())
	s.Require().NoError(err)

	s.Require().Equal(makeProposalCommunityPoolSpendMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionProposalCommunityPoolSpendMsgType, s.xplac.GetMsgType())

	distCommunityPoolSpendTxbytes, err := s.xplac.CommunityPoolSpend(communityPoolSpendMsg).CreateAndSignTx()
	s.Require().NoError(err)

	distCommunityPoolSpendJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(distCommunityPoolSpendTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.DistCommunityPoolSpendTxTemplates, string(distCommunityPoolSpendJsonTxbytes))

	// withdraw rewards
	withdrawRewardsMsg := types.WithdrawRewardsMsg{
		DelegatorAddr: accounts[0].Address.String(),
		ValidatorAddr: sdk.ValAddress(accounts[0].Address).String(),
		Commission:    true,
	}
	s.xplac.WithdrawRewards(withdrawRewardsMsg)

	makeWithdrawRewardsMsg, err := mdist.MakeWithdrawRewardsMsg(withdrawRewardsMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeWithdrawRewardsMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionWithdrawRewardsMsgType, s.xplac.GetMsgType())

	distWithdrawRewardsTxbytes, err := s.xplac.WithdrawRewards(withdrawRewardsMsg).CreateAndSignTx()
	s.Require().NoError(err)

	distWithdrawRewardsJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(distWithdrawRewardsTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.DistWithdrawRewardsTxTemplates, string(distWithdrawRewardsJsonTxbytes))

	// set withdraw address
	setWithdrawAddrMsg := types.SetWithdrawAddrMsg{
		WithdrawAddr: accounts[0].Address.String(),
	}
	s.xplac.SetWithdrawAddr(setWithdrawAddrMsg)

	makeSetWithdrawAddrMsg, err := mdist.MakeSetWithdrawAddrMsg(setWithdrawAddrMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeSetWithdrawAddrMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionSetWithdrawAddrMsgType, s.xplac.GetMsgType())

	distSetWithdrawAddrTxbytes, err := s.xplac.SetWithdrawAddr(setWithdrawAddrMsg).CreateAndSignTx()
	s.Require().NoError(err)

	distSetWithdrawAddrJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(distSetWithdrawAddrTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.DistSetWithdrawAddrTxTemplates, string(distSetWithdrawAddrJsonTxbytes))
}

func (s *IntegrationTestSuite) TestDistribution() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)

	val := s.network.Validators[0].ValAddress.String()
	// query dist params
	s.xplac.DistributionParams()

	makeQueryDistributionParamsMsg, err := mdist.MakeQueryDistributionParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDistributionParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQueryDistributionParamsMsgType, s.xplac.GetMsgType())

	// validator outstanding rewards
	validatorOutstandingRewardsMsg := types.ValidatorOutstandingRewardsMsg{
		ValidatorAddr: val,
	}
	s.xplac.ValidatorOutstandingRewards(validatorOutstandingRewardsMsg)

	makeValidatorOutstandingRewardsMsg, err := mdist.MakeValidatorOutstandingRewardsMsg(validatorOutstandingRewardsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeValidatorOutstandingRewardsMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionValidatorOutstandingRewardsMsgType, s.xplac.GetMsgType())

	// dist commission
	queryDistCommissionMsg := types.QueryDistCommissionMsg{
		ValidatorAddr: val,
	}
	s.xplac.DistCommission(queryDistCommissionMsg)

	makeQueryDistCommissionMsg, err := mdist.MakeQueryDistCommissionMsg(queryDistCommissionMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDistCommissionMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQueryDistCommissionMsgType, s.xplac.GetMsgType())

	// dist slashes
	queryDistSlashesMsg := types.QueryDistSlashesMsg{
		ValidatorAddr: val,
	}
	s.xplac.DistSlashes(queryDistSlashesMsg)

	makeQueryDistSlashesMsg, err := mdist.MakeQueryDistSlashesMsg(queryDistSlashesMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDistSlashesMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQuerySlashesMsgType, s.xplac.GetMsgType())

	// dist rewards
	queryDistRewardsMsg := types.QueryDistRewardsMsg{
		DelegatorAddr: accounts[0].Address.String(),
		ValidatorAddr: val,
	}
	s.xplac.DistRewards(queryDistRewardsMsg)

	makeQueryDistRewwradsMsg, err := mdist.MakeQueryDistRewardsMsg(queryDistRewardsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDistRewwradsMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQueryRewardsMsgType, s.xplac.GetMsgType())

	// total rewards
	queryDistRewardsMsg = types.QueryDistRewardsMsg{
		DelegatorAddr: accounts[0].Address.String(),
	}
	s.xplac.DistRewards(queryDistRewardsMsg)

	makeQueryDistTotalRewardsMsg, err := mdist.MakeQueryDistTotalRewardsMsg(queryDistRewardsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDistTotalRewardsMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQueryTotalRewardsMsgType, s.xplac.GetMsgType())

	// community pool
	s.xplac.CommunityPool()

	makeQueryCommunityPoolMsg, err := mdist.MakeQueryCommunityPoolMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryCommunityPoolMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQueryCommunityPoolMsgType, s.xplac.GetMsgType())
}
