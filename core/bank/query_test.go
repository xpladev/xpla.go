package bank_test

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"
)

var (
	validatorNumber = 2
)

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

func (s *IntegrationTestSuite) TestAllBalancesAndBalance() {
	validator := s.network.Validators[0]
	addr := validator.Address.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		bankBalancesMsg := types.BankBalancesMsg{
			Address: addr,
		}

		res, err := s.xplac.BankBalances(bankBalancesMsg).Query()
		s.Require().NoError(err)

		var allBalancesResponse banktypes.QueryAllBalancesResponse
		jsonpb.Unmarshal(strings.NewReader(res), &allBalancesResponse)

		bal1, err := util.FromStringToBigInt("400000000000000000000")
		s.Require().NoError(err)
		bal2, err := util.FromStringToBigInt("1000000000000000000000")
		s.Require().NoError(err)

		denom1 := "axpla"
		denom2 := "node0token"

		s.Require().Equal(bal1, allBalancesResponse.Balances[0].Amount.BigInt())
		s.Require().Equal(bal2, allBalancesResponse.Balances[1].Amount.BigInt())
		s.Require().Equal(denom1, allBalancesResponse.Balances[0].Denom)
		s.Require().Equal(denom2, allBalancesResponse.Balances[1].Denom)

		if i == 1 {
			// LCD not supported
			bankBalancesMsg = types.BankBalancesMsg{
				Address: addr,
				Denom:   denom1,
			}
			res, err = s.xplac.BankBalances(bankBalancesMsg).Query()
			s.Require().NoError(err)

			var balanceResponse banktypes.QueryBalanceResponse
			jsonpb.Unmarshal(strings.NewReader(res), &balanceResponse)

			s.Require().Equal(denom1, balanceResponse.Balance.Denom)
			s.Require().Equal(bal1, balanceResponse.Balance.Amount.BigInt())
		}
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestDenomMetadata() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.DenomMetadata().Query()
		s.Require().NoError(err)

		var denomsMetadataResponse banktypes.QueryDenomsMetadataResponse
		jsonpb.Unmarshal(strings.NewReader(res), &denomsMetadataResponse)

		s.Require().Equal(2, len(denomsMetadataResponse.Metadatas))
		s.Require().Equal(types.XplaDenom, denomsMetadataResponse.Metadatas[0].Base)
		s.Require().Equal("node0token", denomsMetadataResponse.Metadatas[1].Base)

		denomMetadataMsg := types.DenomMetadataMsg{
			Denom: types.XplaDenom,
		}
		res, err = s.xplac.DenomMetadata(denomMetadataMsg).Query()
		s.Require().NoError(err)

		var denomMetadataResponse banktypes.QueryDenomMetadataResponse
		jsonpb.Unmarshal(strings.NewReader(res), &denomMetadataResponse)

		s.Require().Equal(types.XplaDenom, denomMetadataResponse.Metadata.Base)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestBankTotal() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		bal2, err := util.FromStringToBigInt("1000000000000000000000")
		s.Require().NoError(err)
		bal3, err := util.FromStringToBigInt("1000000000000000000000")
		s.Require().NoError(err)

		denom1 := "axpla"
		denom2 := "node0token"
		denom3 := "node1token"

		res, err := s.xplac.Total().Query()
		s.Require().NoError(err)

		var totalSupplyResponse banktypes.QueryTotalSupplyResponse
		jsonpb.Unmarshal(strings.NewReader(res), &totalSupplyResponse)

		s.Require().Equal(denom1, totalSupplyResponse.Supply[0].Denom)
		s.Require().Equal(denom2, totalSupplyResponse.Supply[1].Denom)
		s.Require().Equal(bal2, totalSupplyResponse.Supply[1].Amount.BigInt())
		s.Require().Equal(denom3, totalSupplyResponse.Supply[2].Denom)
		s.Require().Equal(bal3, totalSupplyResponse.Supply[2].Amount.BigInt())

		totalMsg := types.TotalMsg{
			Denom: denom1,
		}
		res, err = s.xplac.Total(totalMsg).Query()
		s.Require().NoError(err)

		var supplyOfResponse banktypes.QuerySupplyOfResponse
		jsonpb.Unmarshal(strings.NewReader(res), &supplyOfResponse)

		s.Require().Equal(denom1, supplyOfResponse.Amount.Denom)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
