package slashing_test

import (
	"math/rand"

	mslashing "github.com/xpladev/xpla.go/core/slashing"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestSlashingTx() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)

	s.xplac.WithPrivateKey(accounts[0].PrivKey)
	// unjail
	s.xplac.Unjail()

	makeUnjailMsg, err := mslashing.MakeUnjailMsg(s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeUnjailMsg, s.xplac.GetMsg())
	s.Require().Equal(mslashing.SlashingModule, s.xplac.GetModule())
	s.Require().Equal(mslashing.SlahsingUnjailMsgType, s.xplac.GetMsgType())

	slashingUnjailTxbytes, err := s.xplac.Unjail().CreateAndSignTx()
	s.Require().NoError(err)

	slashingUnjailJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(slashingUnjailTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.SlashingUnjailTxTemplates, string(slashingUnjailJsonTxbytes))
}

func (s *IntegrationTestSuite) TestSlashing() {
	// slashing params
	s.xplac.SlashingParams()

	makeQuerySlashingParamsMsg, err := mslashing.MakeQuerySlashingParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQuerySlashingParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mslashing.SlashingModule, s.xplac.GetModule())
	s.Require().Equal(mslashing.SlashingQuerySlashingParamsMsgType, s.xplac.GetMsgType())

	// signing infos
	s.xplac.SigningInfos()

	makeQuerySigningInfosMsg, err := mslashing.MakeQuerySigningInfosMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQuerySigningInfosMsg, s.xplac.GetMsg())
	s.Require().Equal(mslashing.SlashingModule, s.xplac.GetModule())
	s.Require().Equal(mslashing.SlashingQuerySigningInfosMsgType, s.xplac.GetMsgType())

	// signing info
	signingInfoMsg := types.SigningInfoMsg{
		ConsPubKey: `{"@type": "/cosmos.crypto.ed25519.PubKey","key": "6RBPm24ckoWhRt8mArcSCnEKvt0FMGvcaMwchfZ3ue8="}`,
	}
	s.xplac.SigningInfos(signingInfoMsg)

	makeQuerySigningInfoMsg, err := mslashing.MakeQuerySigningInfoMsg(signingInfoMsg, s.xplac.GetEncoding())
	s.Require().NoError(err)

	s.Require().Equal(makeQuerySigningInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mslashing.SlashingModule, s.xplac.GetModule())
	s.Require().Equal(mslashing.SlashingQuerySigningInfoMsgType, s.xplac.GetMsgType())
}
