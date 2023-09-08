package authz_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/authz/client/cli"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"
)

var validatorNumber = 2

type IntegrationTestSuite struct {
	suite.Suite

	xplac      provider.XplaClient
	apis       []string
	accounts   []simtypes.Account
	testTxHash string

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	src := rand.NewSource(1)
	r := rand.New(src)
	s.accounts = testutil.RandomAccounts(r, 2)
	s.testTxHash = "B6BBBB649F19E8970EF274C0083FE945FD38AD8C524D68BB3FE3A20D72DF03C4"

	balanceBigInt, err := util.FromStringToBigInt("1000000000000000000000000000")
	s.Require().NoError(err)

	genesisState := s.cfg.GenesisState

	// add genesis account
	var authGenesis authtypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[authtypes.ModuleName], &authGenesis))

	var genAccounts []authtypes.GenesisAccount

	genAccounts = append(genAccounts, &ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(s.accounts[0].Address, nil, 0, 0),
		CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
	})
	genAccounts = append(genAccounts, &ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(s.accounts[1].Address, nil, 0, 0),
		CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
	})

	accounts, err := authtypes.PackAccounts(genAccounts)
	s.Require().NoError(err)

	authGenesis.Accounts = accounts

	authGenesisBz, err := s.cfg.Codec.MarshalJSON(&authGenesis)
	s.Require().NoError(err)
	genesisState[authtypes.ModuleName] = authGenesisBz

	// add balances
	var bankGenesis banktypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[banktypes.ModuleName], &bankGenesis))

	bankGenesis.Balances = []banktypes.Balance{
		{
			Address: s.accounts[0].Address.String(),
			Coins: sdk.Coins{
				sdk.NewCoin(types.XplaDenom, sdk.NewIntFromBigInt(balanceBigInt)),
			},
		},
		{
			Address: s.accounts[1].Address.String(),
			Coins: sdk.Coins{
				sdk.NewCoin(types.XplaDenom, sdk.NewIntFromBigInt(balanceBigInt)),
			},
		},
	}

	bankGenesisBz, err := s.cfg.Codec.MarshalJSON(&bankGenesis)
	s.Require().NoError(err)
	genesisState[banktypes.ModuleName] = bankGenesisBz

	s.cfg.GenesisState = genesisState
	s.network = network.New(s.T(), s.cfg)
	s.Require().NoError(s.network.WaitForNextBlock())

	val1 := s.network.Validators[0]
	val2 := s.network.Validators[1]

	_, err = ExecGrant(
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

	val := s.network.Validators[0]
	newAddr := s.accounts[0].Address

	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
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
