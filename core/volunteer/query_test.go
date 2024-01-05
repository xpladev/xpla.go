package volunteer_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/spf13/cobra"
	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/xpladev/xpla.go/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingcli "github.com/cosmos/cosmos-sdk/x/staking/client/cli"
	"github.com/stretchr/testify/suite"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"
	volunteercli "github.com/xpladev/xpla/x/volunteer/client/cli"
)

var (
	validatorNumber                    = 2
	registerVolunteerValidatorTestFile = "../../util/testutil/test_files/registerVolunteer.json"
)

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

	_, err := MsgRegisterVolunteerValidatorExec(
		val.ClientCtx,
		registerVolunteerValidatorTestFile,
		val.AdditionalAccount.Address,
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

func (s *IntegrationTestSuite) TestProposalRegisterVolunteerValidator() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryProposalMsg := types.QueryProposalMsg{
			ProposalID: "1",
		}
		res, err := s.xplac.QueryProposal(queryProposalMsg).Query()
		s.Require().NoError(err)

		var queryProposalResponse govtypes.QueryProposalResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryProposalResponse)

		var content govtypes.Content
		s.xplac.GetEncoding().InterfaceRegistry.UnpackAny(queryProposalResponse.Proposal.Content, &content)

		s.Require().Equal("/xpla.volunteer.v1beta1.RegisterVolunteerValidatorProposal", queryProposalResponse.Proposal.Content.TypeUrl)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(types.XplaDenom, sdk.NewInt(10))).String()),
}

func MsgRegisterVolunteerValidatorExec(clientCtx cmclient.Context, proposalFilePath string, from sdk.AccAddress,
	extraArgs ...string,
) (sdktestutil.BufferWriter, error) {
	args := []string{
		proposalFilePath,
		fmt.Sprintf("--%s=%s", stakingcli.FlagMoniker, "moniker"),
		fmt.Sprintf("--%s=%s", stakingcli.FlagIdentity, "identity"),
		fmt.Sprintf("--%s=%s", stakingcli.FlagWebsite, "website"),
		fmt.Sprintf("--%s=%s", stakingcli.FlagSecurityContact, "security"),
		fmt.Sprintf("--%s=%s", stakingcli.FlagDetails, "details"),
		fmt.Sprintf("--%s=%s", stakingcli.FlagPubKey, `{"@type": "/cosmos.crypto.ed25519.PubKey", "key": "2z2yttKfEsLQyQnHYdgKEuky9zB3gscxapn9IyexxWk="}`),
		fmt.Sprintf("--%s=%s", stakingcli.FlagAmount, "10000000axpla"),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from.String()),
		fmt.Sprintf("--%s=%d", flags.FlagGas, 3000000),
	}
	args = append(args, extraArgs...)

	args = append(args, commonArgs...)
	return execTestCLICmd(clientCtx, volunteercli.GetSubmitProposalRegisterVolunteerValidator(), args)
}

func execTestCLICmd(clientCtx cmclient.Context, cmd *cobra.Command, extraArgs []string) (sdktestutil.BufferWriter, error) {
	flags.AddTxFlagsToCmd(cmd)
	cmd.SetArgs(extraArgs)

	_, out := sdktestutil.ApplyMockIO(cmd)
	clientCtx = clientCtx.WithOutput(out)

	ctx := context.Background()
	ctx = context.WithValue(ctx, cmclient.ClientContextKey, &clientCtx)

	if err := cmd.ExecuteContext(ctx); err != nil {
		return out, err
	}

	return out, nil
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	cfg.AdditionalAccount = true
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
