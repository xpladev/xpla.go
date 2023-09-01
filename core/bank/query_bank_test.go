package bank_test

import (
	"strings"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/suite"
	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/xpladev/xpla.go/util/testutil/network"
)

var (
	validatorNumber = 2
)

type IntegrationTestSuite struct {
	suite.Suite

	xplac *client.XplaClient
	apis  []string

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	genesisState := s.cfg.GenesisState
	var bankGenesis banktypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[banktypes.ModuleName], &bankGenesis))

	bankGenesisBz, err := s.cfg.Codec.MarshalJSON(&bankGenesis)
	s.Require().NoError(err)
	genesisState[banktypes.ModuleName] = bankGenesisBz
	s.cfg.GenesisState = genesisState

	s.network = network.New(s.T(), s.cfg)
	s.Require().NoError(s.network.WaitForNextBlock())

	s.xplac = client.NewTestXplaClient()
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
	s.xplac = client.ResetXplac(s.xplac)
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
	s.xplac = client.ResetXplac(s.xplac)
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
	s.xplac = client.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
