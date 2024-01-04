package authz_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/authz/client/cli"
	"github.com/stretchr/testify/suite"
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

	val1 := s.network.Validators[0]
	val2 := s.network.Validators[1]

	_, err := ExecGrant(
		val1,
		[]string{
			val2.Address.String(),
			"send",
			fmt.Sprintf("--%s=100axpla", cli.FlagSpendLimit),
			fmt.Sprintf("--%s=%s", flags.FlagFrom, val1.Address.String()),
			fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
			fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
			fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		},
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

func (s *IntegrationTestSuite) TestAuthzGrant() {
	granter := s.network.Validators[0].Address.String()
	grantee := s.network.Validators[1].Address.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		msg1 := types.QueryAuthzGrantMsg{
			Granter: granter,
			Grantee: grantee,
		}
		res1, err := s.xplac.QueryAuthzGrants(msg1).Query()
		s.Require().NoError(err)

		var queryGrantsResponse authz.QueryGrantsResponse
		jsonpb.Unmarshal(strings.NewReader(res1), &queryGrantsResponse)

		s.Require().Equal(1, len(queryGrantsResponse.Grants))

		msg2 := types.QueryAuthzGrantMsg{
			Granter: granter,
		}
		res2, err := s.xplac.QueryAuthzGrants(msg2).Query()
		s.Require().NoError(err)

		var queryGranterGrantsResponse authz.QueryGranterGrantsResponse
		jsonpb.Unmarshal(strings.NewReader(res2), &queryGranterGrantsResponse)

		s.Require().Equal(1, len(queryGranterGrantsResponse.Grants))

		msg3 := types.QueryAuthzGrantMsg{
			Grantee: grantee,
		}
		res3, err := s.xplac.QueryAuthzGrants(msg3).Query()
		s.Require().NoError(err)

		var queryGranteeGrantsResponse authz.QueryGranteeGrantsResponse
		jsonpb.Unmarshal(strings.NewReader(res3), &queryGranteeGrantsResponse)

		s.Require().Equal(1, len(queryGranteeGrantsResponse.Grants))
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func ExecGrant(val *network.Validator, args []string) (sdktestutil.BufferWriter, error) {
	cmd := cli.NewCmdGrantAuthorization()
	clientCtx := val.ClientCtx
	return clitestutil.ExecTestCLICmd(clientCtx, cmd, args)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
