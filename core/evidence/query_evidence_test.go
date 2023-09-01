package evidence_test

import (
	"strings"
	"testing"

	"github.com/xpladev/xpla.go/client"

	"github.com/gogo/protobuf/jsonpb"

	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/stretchr/testify/suite"
	"github.com/xpladev/xpla.go/util/testutil/network"
)

var validatorNumber = 2

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

func (s IntegrationTestSuite) TestEvidence() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.QueryEvidence().Query()
		s.Require().NoError(err)

		var queryAllEvidenceResponse evidencetypes.QueryAllEvidenceResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryAllEvidenceResponse)

		s.Require().Equal(0, len(queryAllEvidenceResponse.Evidence))
	}
	s.xplac = client.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
