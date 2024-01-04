package slashing_test

import (
	"strings"
	"testing"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/xpladev/xpla.go/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
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

func (s *IntegrationTestSuite) TestParams() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.SlashingParams().Query()
		s.Require().NoError(err)

		var queryParamsResponse slashingtypes.QueryParamsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryParamsResponse)

		s.Require().Equal(int64(100), queryParamsResponse.Params.SignedBlocksWindow)
		s.Require().Equal("0.500000000000000000", queryParamsResponse.Params.MinSignedPerWindow.String())
		s.Require().Equal("10m0s", queryParamsResponse.Params.DowntimeJailDuration.String())
		s.Require().Equal("0.050000000000000000", queryParamsResponse.Params.SlashFractionDoubleSign.String())
		s.Require().Equal("0.010000000000000000", queryParamsResponse.Params.SlashFractionDowntime.String())
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestSigningInfos() {
	val := s.network.Validators[0]

	valConsAddr1, err := sdk.ConsAddressFromHex(val.PubKey.Address().String())
	s.Require().NoError(err)

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		// a validator signing info
		signingInfoMsg := types.SigningInfoMsg{
			ConsAddr: valConsAddr1.String(),
		}

		res1, err := s.xplac.SigningInfos(signingInfoMsg).Query()
		s.Require().NoError(err)

		var querySigningInfoResponse slashingtypes.QuerySigningInfoResponse
		jsonpb.Unmarshal(strings.NewReader(res1), &querySigningInfoResponse)

		s.Require().Equal(valConsAddr1.String(), querySigningInfoResponse.ValSigningInfo.Address)
		s.Require().Equal(int64(0), querySigningInfoResponse.ValSigningInfo.MissedBlocksCounter)

		// validators signing info
		res2, err := s.xplac.SigningInfos().Query()
		s.Require().NoError(err)

		var querySigningInfosResponse slashingtypes.QuerySigningInfosResponse
		jsonpb.Unmarshal(strings.NewReader(res2), &querySigningInfosResponse)

		s.Require().Equal(validatorNumber, len(querySigningInfosResponse.Info))
		s.Require().Equal(int64(0), querySigningInfosResponse.Info[0].MissedBlocksCounter)
		s.Require().Equal(int64(0), querySigningInfosResponse.Info[1].MissedBlocksCounter)

	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
