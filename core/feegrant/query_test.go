package feegrant_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/xpladev/xpla.go/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/x/feegrant/client/cli"
	"github.com/stretchr/testify/suite"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"
)

var (
	oneYear         = 365 * 24 * 60 * 60
	validatorNumber = 2
)

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

	granter := s.network.Validators[0].Address
	grantee := s.network.Validators[1].Address

	s.createGrant(granter, grantee)

	_, err := feegrant.NewGrant(granter, grantee, &feegrant.BasicAllowance{
		SpendLimit: sdk.NewCoins(sdk.NewCoin(types.XplaDenom, sdk.NewInt(100))),
	})
	s.Require().NoError(err)

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

func (s *IntegrationTestSuite) TestFeegrants() {
	granter := s.network.Validators[0].Address
	grantee := s.network.Validators[1].Address

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryFeeGrantMsg := types.QueryFeeGrantMsg{
			Granter: granter.String(),
			Grantee: grantee.String(),
		}
		res1, err := s.xplac.QueryFeeGrants(queryFeeGrantMsg).Query()
		s.Require().NoError(err)

		var queryAllowanceResponse feegrant.QueryAllowanceResponse
		jsonpb.Unmarshal(strings.NewReader(res1), &queryAllowanceResponse)

		s.Require().Equal(granter.String(), queryAllowanceResponse.Allowance.Granter)
		s.Require().Equal(grantee.String(), queryAllowanceResponse.Allowance.Grantee)

		// LCD not supported
		if i != 0 {
			queryFeeGrantMsgGranter := types.QueryFeeGrantMsg{
				Granter: granter.String(),
			}
			res2, err := s.xplac.QueryFeeGrants(queryFeeGrantMsgGranter).Query()
			s.Require().NoError(err)

			var queryAllowancesByGranterResponse feegrant.QueryAllowancesByGranterResponse
			jsonpb.Unmarshal(strings.NewReader(res2), &queryAllowancesByGranterResponse)

			s.Require().Equal(granter.String(), queryAllowancesByGranterResponse.Allowances[0].Granter)
			s.Require().Equal(grantee.String(), queryAllowancesByGranterResponse.Allowances[0].Grantee)
		}

		queryFeeGrantMsgGrantee := types.QueryFeeGrantMsg{
			Grantee: grantee.String(),
		}
		res3, err := s.xplac.QueryFeeGrants(queryFeeGrantMsgGrantee).Query()
		s.Require().NoError(err)

		var queryAllowancesResponse feegrant.QueryAllowancesResponse
		jsonpb.Unmarshal(strings.NewReader(res3), &queryAllowancesResponse)

		s.Require().Equal(granter.String(), queryAllowancesResponse.Allowances[0].Granter)
		s.Require().Equal(grantee.String(), queryAllowancesResponse.Allowances[0].Grantee)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) createGrant(granter, grantee sdk.Address) {
	val := s.network.Validators[0]

	clientCtx := val.ClientCtx
	commonFlags := []string{
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	fee := sdk.NewCoin(types.XplaDenom, sdk.NewInt(100))

	args := append(
		[]string{
			granter.String(),
			grantee.String(),
			fmt.Sprintf("--%s=%s", cli.FlagSpendLimit, fee.String()),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, granter),
			fmt.Sprintf("--%s=%s", cli.FlagExpiration, getFormattedExpiration(int64(oneYear))),
		},
		commonFlags...,
	)

	cmd := cli.NewCmdFeeGrant()

	_, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
	s.Require().NoError(err)
	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func getFormattedExpiration(duration int64) string {
	return time.Now().Add(time.Duration(duration) * time.Second).Format(time.RFC3339)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
