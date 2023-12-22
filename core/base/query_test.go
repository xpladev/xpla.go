package base_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/stretchr/testify/suite"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/xpladev/xpla.go/util/testutil/network"
)

var validatorNumber = 1

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

func (s *IntegrationTestSuite) TestNodeInfo() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.NodeInfo().Query()
		s.Require().NoError(err)

		var getNodeInfoResponse tmservice.GetNodeInfoResponse
		jsonpb.Unmarshal(strings.NewReader(res), &getNodeInfoResponse)

		s.Require().Equal(testutil.TestChainId, getNodeInfoResponse.DefaultNodeInfo.Network)
		s.Require().Equal("node0", getNodeInfoResponse.DefaultNodeInfo.Moniker)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestSyncing() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.Syncing().Query()
		s.Require().NoError(err)

		var getSyncingResponse tmservice.GetSyncingResponse
		jsonpb.Unmarshal(strings.NewReader(res), &getSyncingResponse)

		s.Require().False(getSyncingResponse.Syncing)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestBlock() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res1, err := s.xplac.Block().Query()
		s.Require().NoError(err)

		queryBlockHeight := 1
		blockMsg := types.BlockMsg{
			Height: util.FromIntToString(queryBlockHeight),
		}
		res2, err := s.xplac.Block(blockMsg).Query()
		s.Require().NoError(err)

		if i == 0 {
			var resultBlock ctypes.ResultBlock
			json.Unmarshal([]byte(res1), &resultBlock)

			s.Require().Equal(testutil.TestChainId, resultBlock.Block.ChainID)

			json.Unmarshal([]byte(res2), &resultBlock)

			s.Require().Equal(testutil.TestChainId, resultBlock.Block.ChainID)

		} else {
			var getLatestBlockResponse tmservice.GetLatestBlockResponse
			jsonpb.Unmarshal(strings.NewReader(res1), &getLatestBlockResponse)

			s.Require().Equal(testutil.TestChainId, getLatestBlockResponse.Block.Header.ChainID)

			var getBlockByHeightResponse tmservice.GetBlockByHeightResponse
			jsonpb.Unmarshal(strings.NewReader(res2), &getBlockByHeightResponse)

			s.Require().Equal(testutil.TestChainId, getBlockByHeightResponse.Block.Header.ChainID)
		}

	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestValidatorSet() {
	val := s.network.Validators[0]

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.ValidatorSet().Query()
		s.Require().NoError(err)

		var getLatestValidatorSet tmservice.GetLatestValidatorSetResponse
		jsonpb.Unmarshal(strings.NewReader(res), &getLatestValidatorSet)

		var pubkey cryptotypes.PubKey
		s.xplac.GetEncoding().InterfaceRegistry.UnpackAny(getLatestValidatorSet.Validators[0].PubKey, &pubkey)

		s.Require().Equal(val.PubKey, pubkey)

		validatorSetMsg := types.ValidatorSetMsg{
			Height: "1",
		}

		res, err = s.xplac.ValidatorSet(validatorSetMsg).Query()
		s.Require().NoError(err)

		var getValidatorSetByHeightResponse tmservice.GetValidatorSetByHeightResponse
		jsonpb.Unmarshal(strings.NewReader(res), &getValidatorSetByHeightResponse)

		s.xplac.GetEncoding().InterfaceRegistry.UnpackAny(getValidatorSetByHeightResponse.Validators[0].PubKey, &pubkey)

		s.Require().Equal(val.PubKey, pubkey)

	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.ChainID = testutil.TestChainId
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
