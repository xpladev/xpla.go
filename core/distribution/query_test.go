package distribution_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/xpladev/xpla.go/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/stretchr/testify/suite"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"
)

var validatorNumber = 2

type IntegrationTestSuite struct {
	suite.Suite

	xplac provider.XplaClient
	apis  []string

	cfg     network.Config
	network network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), s.cfg)
	s.Require().NoError(s.network.WaitForNextBlock())

	val := s.network.Validators[0]
	val2 := s.network.Validators[1]

	del := sdk.NewCoin(types.XplaDenom, sdk.NewInt(1000))

	_, err := msgDelegateExec(
		val.ClientCtx,
		val.Address,
		val2.ValAddress,
		del,
	)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	s.xplac = client.NewXplaClient(testutil.TestChainId).WithVerbose(1)
	s.apis = []string{
		s.network.Validators[0].APIAddress,
		s.network.Validators[0].AppConfig.GRPC.Address,
	}
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestDistributionParams() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.DistributionParams().Query()
		s.Require().NoError(err)

		var queryParamsResponse disttypes.QueryParamsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryParamsResponse)

		s.Require().Equal("0.020000000000000000", queryParamsResponse.Params.CommunityTax.String())
		s.Require().Equal("0.010000000000000000", queryParamsResponse.Params.BaseProposerReward.String())
		s.Require().Equal("0.040000000000000000", queryParamsResponse.Params.BonusProposerReward.String())
		s.Require().True(queryParamsResponse.Params.WithdrawAddrEnabled)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestValidatorOutstandingRewards() {
	valAddr := s.network.Validators[0].ValAddress.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		msg := types.ValidatorOutstandingRewardsMsg{
			ValidatorAddr: valAddr,
		}

		res, err := s.xplac.ValidatorOutstandingRewards(msg).Query()
		s.Require().NoError(err)

		var queryValidatorOutstandingRewardsResponse disttypes.QueryValidatorOutstandingRewardsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryValidatorOutstandingRewardsResponse)

		s.Require().Equal(types.XplaDenom, queryValidatorOutstandingRewardsResponse.Rewards.Rewards[0].Denom)
		s.Require().NotEqual("0", queryValidatorOutstandingRewardsResponse.Rewards.Rewards[0].Amount.String())
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestDistCommission() {
	valAddr := s.network.Validators[0].ValAddress.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryDistCommissionMsg := types.QueryDistCommissionMsg{
			ValidatorAddr: valAddr,
		}
		res, err := s.xplac.DistCommission(queryDistCommissionMsg).Query()
		s.Require().NoError(err)

		var queryValidatorCommissionResponse disttypes.QueryValidatorCommissionResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryValidatorCommissionResponse)

		s.Require().Equal(types.XplaDenom, queryValidatorCommissionResponse.Commission.Commission[0].Denom)
		s.Require().NotEqual("0", queryValidatorCommissionResponse.Commission.Commission[0].Amount.String())

	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestDistSlashes() {
	valAddr := s.network.Validators[0].ValAddress.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryDistSlashesMsg := types.QueryDistSlashesMsg{
			ValidatorAddr: valAddr,
		}

		res, err := s.xplac.DistSlashes(queryDistSlashesMsg).Query()
		s.Require().NoError(err)

		var queryValidatorSlashesResponse disttypes.QueryValidatorSlashesResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryValidatorSlashesResponse)

		s.Require().Equal(0, len(queryValidatorSlashesResponse.Slashes))
	}
	s.xplac = provider.ResetXplac(s.xplac)
}
func (s *IntegrationTestSuite) TestDistRewards() {

	s.Require().NoError(s.network.WaitForNextBlock())

	delegator := s.network.Validators[0].Address.String()
	validator := s.network.Validators[1].ValAddress.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryDistRewardsMsg := types.QueryDistRewardsMsg{
			DelegatorAddr: delegator,
			ValidatorAddr: validator,
		}

		res, err := s.xplac.DistRewards(queryDistRewardsMsg).Query()
		s.Require().NoError(err)

		var queryDelegationRewardsResponse disttypes.QueryDelegationRewardsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryDelegationRewardsResponse)

		s.Require().NotEqual("0", queryDelegationRewardsResponse.Rewards[0].Amount.String())

		queryDistRewardsMsgTotal := types.QueryDistRewardsMsg{
			DelegatorAddr: delegator,
		}

		res1, err := s.xplac.DistRewards(queryDistRewardsMsgTotal).Query()
		s.Require().NoError(err)

		var queryDelegatorTotalRewardsResponse disttypes.QueryDelegatorTotalRewardsResponse
		json.Unmarshal([]byte(res1), &queryDelegatorTotalRewardsResponse)

		s.Require().NotEqual("0", queryDelegatorTotalRewardsResponse.Rewards[0].Reward[0].Amount.String())
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestCommunityPool() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.CommunityPool().Query()
		s.Require().NoError(err)

		var queryCommunityPoolResponse disttypes.QueryCommunityPoolResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryCommunityPoolResponse)

		s.Require().Equal(types.XplaDenom, queryCommunityPoolResponse.Pool[0].Denom)
		s.Require().NotEqual("0", queryCommunityPoolResponse.Pool[0].Amount.String())
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func msgDelegateExec(clientCtx cmclient.Context, delegator, validator, amount fmt.Stringer, extraArgs ...string) (sdktestutil.BufferWriter, error) {
	args := []string{
		validator.String(),
		amount.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, delegator.String()),
		fmt.Sprintf("--%s=%d", flags.FlagGas, 300000),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(types.XplaDenom, sdk.NewInt(10))).String()),
	}

	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, stakingcli.NewDelegateCmd(), args)
}
func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
