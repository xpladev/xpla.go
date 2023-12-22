package staking_test

import (
	"fmt"
	"math/rand"
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
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	stakingcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"
)

var validatorNumber = 2

type IntegrationTestSuite struct {
	suite.Suite

	xplac    provider.XplaClient
	apis     []string
	accounts []simtypes.Account

	cfg     network.Config
	network network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	src := rand.NewSource(1)
	r := rand.New(src)
	s.accounts = testutil.RandomAccounts(r, 2)

	s.network = network.New(s.T(), s.cfg)
	s.Require().NoError(s.network.WaitForNextBlock())

	val := s.network.Validators[0]
	val2 := s.network.Validators[1]

	del, err := sdk.ParseCoinNormalized("1000axpla")
	s.Require().NoError(err)

	unbond, err := sdk.ParseCoinNormalized("10axpla")
	s.Require().NoError(err)

	_, err = MsgDelegateExec(
		val.ClientCtx,
		val.Address,
		val2.ValAddress,
		del,
	)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// redelegate
	_, err = MsgRedelegateExec(
		val.ClientCtx,
		val.Address,
		val.ValAddress,
		val2.ValAddress,
		unbond,
		fmt.Sprintf("--%s=%d", flags.FlagGas, 300000),
	)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// unbonding
	_, err = MsgUnbondExec(val.ClientCtx, val.Address, val.ValAddress, unbond)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	s.xplac = client.NewXplaClient(testutil.TestChainId)
	s.apis = []string{
		s.network.Validators[0].APIAddress,
		s.network.Validators[0].AppConfig.GRPC.Address,
	}
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestQueryValidators() {
	val1 := s.network.Validators[0]

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		// query a validator
		queryValidatorMsg := types.QueryValidatorMsg{
			ValidatorAddr: val1.ValAddress.String(),
		}

		res1, err := s.xplac.QueryValidators(queryValidatorMsg).Query()
		s.Require().NoError(err)

		var queryValidatorResponse stakingtypes.QueryValidatorResponse
		jsonpb.Unmarshal(strings.NewReader(res1), &queryValidatorResponse)

		s.Require().Equal(val1.ValAddress.String(), queryValidatorResponse.Validator.OperatorAddress)
		s.Require().Equal(false, queryValidatorResponse.Validator.Jailed)
		s.Require().Equal("BOND_STATUS_BONDED", queryValidatorResponse.Validator.Status.String())
		s.Require().Equal("node0", queryValidatorResponse.Validator.Description.Moniker)

		// query validators
		res2, err := s.xplac.QueryValidators().Query()
		s.Require().NoError(err)

		var queryValidatorsResponse stakingtypes.QueryValidatorsResponse
		jsonpb.Unmarshal(strings.NewReader(res2), &queryValidatorsResponse)

		s.Require().Equal(validatorNumber, len(queryValidatorsResponse.Validators))

	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestDelegation() {
	val1 := s.network.Validators[0]
	val2 := s.network.Validators[1]

	// staking delegation
	// del addr != "" / val adr != ""
	queryDelegationMsg := types.QueryDelegationMsg{
		DelegatorAddr: val1.Address.String(),
		ValidatorAddr: val2.ValAddress.String(),
	}

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)

			res, err := s.xplac.QueryDelegation(queryDelegationMsg).Query()
			s.Require().NoError(err)

			var queryValidatorResponse stakingtypes.QueryValidatorResponse
			jsonpb.Unmarshal(strings.NewReader(res), &queryValidatorResponse)

			s.Require().Equal(val2.ValAddress.String(), queryValidatorResponse.Validator.OperatorAddress)

		} else {
			s.xplac.WithGrpc(api)

			res, err := s.xplac.QueryDelegation(queryDelegationMsg).Query()
			s.Require().NoError(err)

			var queryDelegationResponse stakingtypes.QueryDelegationResponse
			jsonpb.Unmarshal(strings.NewReader(res), &queryDelegationResponse)

			s.Require().Equal(val1.Address.String(), queryDelegationResponse.DelegationResponse.Delegation.DelegatorAddress)
			s.Require().Equal("1010.000000000000000000", queryDelegationResponse.DelegationResponse.Delegation.Shares.String())
			s.Require().Equal("1010", queryDelegationResponse.DelegationResponse.Balance.Amount.String())
		}
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestDelegationDelegator() {
	val1 := s.network.Validators[0]

	// staking delegation
	// del addr != "" / val adr == ""
	queryDelegationMsg := types.QueryDelegationMsg{
		DelegatorAddr: val1.Address.String(),
	}

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)

			res, err := s.xplac.QueryDelegation(queryDelegationMsg).Query()
			s.Require().NoError(err)

			var queryValidatorsResponse stakingtypes.QueryValidatorsResponse
			jsonpb.Unmarshal(strings.NewReader(res), &queryValidatorsResponse)

			s.Require().Equal(2, len(queryValidatorsResponse.Validators))

		} else {
			s.xplac.WithGrpc(api)

			res, err := s.xplac.QueryDelegation(queryDelegationMsg).Query()
			s.Require().NoError(err)

			var queryDelegatorDelegationsResponse stakingtypes.QueryDelegatorDelegationsResponse
			jsonpb.Unmarshal(strings.NewReader(res), &queryDelegatorDelegationsResponse)

			s.Require().Equal(2, len(queryDelegatorDelegationsResponse.DelegationResponses))
		}
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestDelegationValidator() {
	val2 := s.network.Validators[1]

	// staking delegation
	// del addr == "" / val adr != ""
	queryDelegationMsg := types.QueryDelegationMsg{
		ValidatorAddr: val2.ValAddress.String(),
	}

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.QueryDelegation(queryDelegationMsg).Query()
		s.Require().NoError(err)

		var queryValidatorDelegationsResponse stakingtypes.QueryValidatorDelegationsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryValidatorDelegationsResponse)

		s.Require().Equal(2, len(queryValidatorDelegationsResponse.DelegationResponses))
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestUnbonding() {
	val1 := s.network.Validators[0]

	queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{
		DelegatorAddr: val1.Address.String(),
		ValidatorAddr: val1.ValAddress.String(),
	}

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.QueryUnbondingDelegation(queryUnbondingDelegationMsg).Query()
		s.Require().NoError(err)

		var queryUnbondingDelegationResponse stakingtypes.QueryUnbondingDelegationResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryUnbondingDelegationResponse)

		s.Require().Equal(val1.Address.String(), queryUnbondingDelegationResponse.Unbond.DelegatorAddress)
		s.Require().Equal(val1.ValAddress.String(), queryUnbondingDelegationResponse.Unbond.ValidatorAddress)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestUnbondingDelegator() {
	val1 := s.network.Validators[0]

	queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{
		DelegatorAddr: val1.Address.String(),
	}

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.QueryUnbondingDelegation(queryUnbondingDelegationMsg).Query()
		s.Require().NoError(err)

		var queryDelegatorUnbondingDelegationsResponse stakingtypes.QueryDelegatorUnbondingDelegationsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryDelegatorUnbondingDelegationsResponse)

		s.Require().Equal(val1.Address.String(), queryDelegatorUnbondingDelegationsResponse.UnbondingResponses[0].DelegatorAddress)
		s.Require().Equal(val1.ValAddress.String(), queryDelegatorUnbondingDelegationsResponse.UnbondingResponses[0].ValidatorAddress)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestUnbondingValidator() {
	val1 := s.network.Validators[0]

	queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{
		ValidatorAddr: val1.ValAddress.String(),
	}

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.QueryUnbondingDelegation(queryUnbondingDelegationMsg).Query()
		s.Require().NoError(err)

		var queryValidatorUnbondingDelegationsResponse stakingtypes.QueryValidatorUnbondingDelegationsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryValidatorUnbondingDelegationsResponse)

		s.Require().Equal(val1.Address.String(), queryValidatorUnbondingDelegationsResponse.UnbondingResponses[0].DelegatorAddress)
		s.Require().Equal(val1.ValAddress.String(), queryValidatorUnbondingDelegationsResponse.UnbondingResponses[0].ValidatorAddress)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestRedelegation() {
	val1 := s.network.Validators[0]
	val2 := s.network.Validators[1]

	queryRedelegationMsg := types.QueryRedelegationMsg{
		DelegatorAddr:    val1.Address.String(),
		SrcValidatorAddr: val1.ValAddress.String(),
		DstValidatorAddr: val2.ValAddress.String(),
	}

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)

			// LCD unsupported
			_, err := s.xplac.QueryRedelegation(queryRedelegationMsg).Query()
			s.Require().Error(err)

		} else {
			s.xplac.WithGrpc(api)

			res, err := s.xplac.QueryRedelegation(queryRedelegationMsg).Query()
			s.Require().NoError(err)

			var queryRedelegationsResponse stakingtypes.QueryRedelegationsResponse
			jsonpb.Unmarshal(strings.NewReader(res), &queryRedelegationsResponse)

			s.Require().Equal(
				val1.Address.String(),
				queryRedelegationsResponse.RedelegationResponses[0].Redelegation.DelegatorAddress,
			)
			s.Require().Equal(
				val1.ValAddress.String(),
				queryRedelegationsResponse.RedelegationResponses[0].Redelegation.ValidatorSrcAddress,
			)
			s.Require().Equal(
				val2.ValAddress.String(),
				queryRedelegationsResponse.RedelegationResponses[0].Redelegation.ValidatorDstAddress,
			)
			s.Require().Equal(
				int64(4),
				queryRedelegationsResponse.RedelegationResponses[0].Entries[0].RedelegationEntry.CreationHeight,
			)
		}
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestRedelegationDelegator() {
	val1 := s.network.Validators[0]
	val2 := s.network.Validators[1]

	queryRedelegationMsg := types.QueryRedelegationMsg{
		DelegatorAddr: val1.Address.String(),
	}

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.QueryRedelegation(queryRedelegationMsg).Query()
		s.Require().NoError(err)

		var queryRedelegationsResponse stakingtypes.QueryRedelegationsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryRedelegationsResponse)

		s.Require().Equal(
			val1.Address.String(),
			queryRedelegationsResponse.RedelegationResponses[0].Redelegation.DelegatorAddress,
		)
		s.Require().Equal(
			val1.ValAddress.String(),
			queryRedelegationsResponse.RedelegationResponses[0].Redelegation.ValidatorSrcAddress,
		)
		s.Require().Equal(
			val2.ValAddress.String(),
			queryRedelegationsResponse.RedelegationResponses[0].Redelegation.ValidatorDstAddress,
		)
		s.Require().Equal(
			int64(4),
			queryRedelegationsResponse.RedelegationResponses[0].Entries[0].RedelegationEntry.CreationHeight,
		)
	}

	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestRedelegationSrcValidator() {
	val1 := s.network.Validators[0]
	val2 := s.network.Validators[1]

	queryRedelegationMsg := types.QueryRedelegationMsg{
		SrcValidatorAddr: val1.ValAddress.String(),
	}

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)

			// LCD unsupported
			_, err := s.xplac.QueryRedelegation(queryRedelegationMsg).Query()
			s.Require().Error(err)

		} else {
			s.xplac.WithGrpc(api)

			res, err := s.xplac.QueryRedelegation(queryRedelegationMsg).Query()
			s.Require().NoError(err)

			var queryRedelegationsResponse stakingtypes.QueryRedelegationsResponse
			jsonpb.Unmarshal(strings.NewReader(res), &queryRedelegationsResponse)

			s.Require().Equal(
				val1.Address.String(),
				queryRedelegationsResponse.RedelegationResponses[0].Redelegation.DelegatorAddress,
			)
			s.Require().Equal(
				val1.ValAddress.String(),
				queryRedelegationsResponse.RedelegationResponses[0].Redelegation.ValidatorSrcAddress,
			)
			s.Require().Equal(
				val2.ValAddress.String(),
				queryRedelegationsResponse.RedelegationResponses[0].Redelegation.ValidatorDstAddress,
			)
			s.Require().Equal(
				int64(4),
				queryRedelegationsResponse.RedelegationResponses[0].Entries[0].RedelegationEntry.CreationHeight,
			)
		}
	}

	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestHistoricalInfo() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		historicalInfoMsg := types.HistoricalInfoMsg{
			Height: "2",
		}
		res, err := s.xplac.HistoricalInfo(historicalInfoMsg).Query()
		s.Require().NoError(err)

		var queryHistoricalInfoResponse stakingtypes.QueryHistoricalInfoResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryHistoricalInfoResponse)

		s.Require().Equal(validatorNumber, len(queryHistoricalInfoResponse.Hist.Valset))
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestStakingPool() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.StakingPool().Query()
		s.Require().NoError(err)

		var queryPoolResponse stakingtypes.QueryPoolResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryPoolResponse)

		s.Require().Equal("200000000000000000990", queryPoolResponse.Pool.BondedTokens.String())
		s.Require().Equal("10", queryPoolResponse.Pool.NotBondedTokens.String())
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestParams() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.StakingParams().Query()
		s.Require().NoError(err)

		var queryParamsResponse stakingtypes.QueryParamsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryParamsResponse)

		s.Require().Equal("axpla", queryParamsResponse.Params.BondDenom)
		s.Require().Equal(uint32(100), queryParamsResponse.Params.MaxValidators)
		s.Require().Equal(uint32(7), queryParamsResponse.Params.MaxEntries)
		s.Require().Equal(uint32(10000), queryParamsResponse.Params.HistoricalEntries)

	}
	s.xplac = provider.ResetXplac(s.xplac)
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(types.XplaDenom, sdk.NewInt(10))).String()),
}

// MsgRedelegateExec creates a redelegate message.
func MsgRedelegateExec(clientCtx cmclient.Context, from, src, dst, amount fmt.Stringer,
	extraArgs ...string,
) (sdktestutil.BufferWriter, error) {
	args := []string{
		src.String(),
		dst.String(),
		amount.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from.String()),
		fmt.Sprintf("--%s=%d", flags.FlagGas, 300000),
	}
	args = append(args, extraArgs...)

	args = append(args, commonArgs...)
	return clitestutil.ExecTestCLICmd(clientCtx, stakingcli.NewRedelegateCmd(), args)
}

// MsgUnbondExec creates a unbond message.
func MsgUnbondExec(clientCtx cmclient.Context, from fmt.Stringer, valAddress,
	amount fmt.Stringer, extraArgs ...string,
) (sdktestutil.BufferWriter, error) {
	args := []string{
		valAddress.String(),
		amount.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from.String()),
		fmt.Sprintf("--%s=%d", flags.FlagGas, 300000),
	}

	args = append(args, commonArgs...)
	args = append(args, extraArgs...)
	return clitestutil.ExecTestCLICmd(clientCtx, stakingcli.NewUnbondCmd(), args)
}

func MsgDelegateExec(clientCtx cmclient.Context, delegator, validator, amount fmt.Stringer, extraArgs ...string) (sdktestutil.BufferWriter, error) {
	args := []string{
		validator.String(),
		amount.String(),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, delegator.String()),
		fmt.Sprintf("--%s=%d", flags.FlagGas, 300000),
	}

	args = append(args, extraArgs...)
	args = append(args, commonArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, stakingcli.NewDelegateCmd(), args)
}
func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
