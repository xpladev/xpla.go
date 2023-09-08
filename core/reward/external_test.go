package reward_test

import (
	"math/rand"

	mreward "github.com/xpladev/xpla.go/core/reward"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestRewardTx() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)

	s.xplac.WithPrivateKey(accounts[0].PrivKey)
	// fund fee collector
	fundFeeCollectorMsg := types.FundFeeCollectorMsg{
		DepositorAddr: accounts[0].Address.String(),
		Amount:        "1000",
	}
	s.xplac.FundFeeCollector(fundFeeCollectorMsg)

	makeFundFeeCollectorMsg, err := mreward.MakeFundFeeCollectorMsg(fundFeeCollectorMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeFundFeeCollectorMsg, s.xplac.GetMsg())
	s.Require().Equal(mreward.RewardModule, s.xplac.GetModule())
	s.Require().Equal(mreward.RewardFundFeeCollectorMsgType, s.xplac.GetMsgType())

	rewardFundFeeCollectorTxbytes, err := s.xplac.FundFeeCollector(fundFeeCollectorMsg).CreateAndSignTx()
	s.Require().NoError(err)

	rewardFundFeeCollectorJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(rewardFundFeeCollectorTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.RewardFundFeeCollectorTxTemplates, string(rewardFundFeeCollectorJsonTxbytes))
}

func (s *IntegrationTestSuite) TestReward() {
	// reward params
	s.xplac.RewardParams()

	makeQueryRewardParamsMsg, err := mreward.MakeQueryRewardParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryRewardParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mreward.RewardModule, s.xplac.GetModule())
	s.Require().Equal(mreward.RewardQueryRewardParamsMsgType, s.xplac.GetMsgType())

	// reward pool
	s.xplac.RewardPool()

	makeQueryRewardPoolMsg, err := mreward.MakeQueryRewardPoolMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryRewardPoolMsg, s.xplac.GetMsg())
	s.Require().Equal(mreward.RewardModule, s.xplac.GetModule())
	s.Require().Equal(mreward.RewardQueryRewardPoolMsgType, s.xplac.GetMsgType())
}
