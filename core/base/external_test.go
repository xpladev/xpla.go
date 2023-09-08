package base_test

import (
	mbase "github.com/xpladev/xpla.go/core/base"
	"github.com/xpladev/xpla.go/types"
)

func (s *IntegrationTestSuite) TestBase() {
	// node info
	s.xplac.NodeInfo()

	makeBaseNodeInfoMsg, err := mbase.MakeBaseNodeInfoMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeBaseNodeInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseNodeInfoMsgType, s.xplac.GetMsgType())

	// syncing
	s.xplac.Syncing()

	makeBaseSyncingMsg, err := mbase.MakeBaseSyncingMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeBaseSyncingMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseSyncingMsgType, s.xplac.GetMsgType())

	// latest block
	s.xplac.Block()

	makeBaseLatestBlockMsg, err := mbase.MakeBaseLatestBlockMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeBaseLatestBlockMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseLatestBlockMsgtype, s.xplac.GetMsgType())

	// block by height
	blockMsg := types.BlockMsg{
		Height: "1",
	}
	s.xplac.Block(blockMsg)

	makeBaseBlockByheightMsg, err := mbase.MakeBaseBlockByHeightMsg(blockMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeBaseBlockByheightMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseBlockByHeightMsgType, s.xplac.GetMsgType())

	// latest validator set
	s.xplac.ValidatorSet()

	makeLatestValidatorSetMsg, err := mbase.MakeLatestValidatorSetMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeLatestValidatorSetMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseLatestValidatorSetMsgType, s.xplac.GetMsgType())

	// validator set by height
	validatorSetMsg := types.ValidatorSetMsg{
		Height: "1",
	}
	s.xplac.ValidatorSet(validatorSetMsg)

	makeValidatorSetByHeightMsg, err := mbase.MakeValidatorSetByHeightMsg(validatorSetMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeValidatorSetByHeightMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseValidatorSetByHeightMsgType, s.xplac.GetMsgType())
}
