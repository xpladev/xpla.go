package crisis_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/core/crisis"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := crisis.NewCoreModule()

	// test get name
	s.Require().Equal(crisis.CrisisModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// invariant broken
	invariantBrokenMsg := types.InvariantBrokenMsg{
		ModuleName:     "bank",
		InvariantRoute: "total-supply",
	}
	s.xplac.InvariantBroken(invariantBrokenMsg)

	makeInvariantRouteMsg, err := crisis.MakeInvariantRouteMsg(invariantBrokenMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeInvariantRouteMsg
	txBuilder, err = c.NewTxRouter(txBuilder, crisis.CrisisInvariantBrokenMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeInvariantRouteMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = client.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = 2
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
