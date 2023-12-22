package params_test

import (
	"strings"
	"testing"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"
)

var (
	validatorNumber = 2
	maxValidators   = 100
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

	genesisState := s.cfg.GenesisState
	var stakingGenesis stakingtypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[stakingtypes.ModuleName], &stakingGenesis))

	stakingGenesis.Params.MaxValidators = uint32(maxValidators)
	stakingGenesisBz, err := s.cfg.Codec.MarshalJSON(&stakingGenesis)
	s.Require().NoError(err)

	genesisState[stakingtypes.ModuleName] = stakingGenesisBz
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

func (s *IntegrationTestSuite) TestQuerySubspace() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}
		subspaceMsg := types.SubspaceMsg{
			Subspace: "staking",
			Key:      "MaxValidators",
		}

		res, err := s.xplac.QuerySubspace(subspaceMsg).Query()
		s.Require().NoError(err)

		var queryParamsResponse proposal.QueryParamsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryParamsResponse)

		s.Require().Equal("staking", queryParamsResponse.Param.Subspace)
		s.Require().Equal("MaxValidators", queryParamsResponse.Param.Key)
		s.Require().Equal(util.FromIntToString(maxValidators), queryParamsResponse.Param.Value)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
