package crisis_test

import (
	"math/rand"

	mcrisis "github.com/xpladev/xpla.go/core/crisis"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCrisisTx() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)

	s.xplac.WithPrivateKey(accounts[0].PrivKey)
	// invariant broken
	invariantBrokenMsg := types.InvariantBrokenMsg{
		ModuleName:     "bank",
		InvariantRoute: "total-supply",
	}
	s.xplac.InvariantBroken(invariantBrokenMsg)

	makeInvariantRouteMsg, err := mcrisis.MakeInvariantRouteMsg(invariantBrokenMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeInvariantRouteMsg, s.xplac.GetMsg())
	s.Require().Equal(mcrisis.CrisisModule, s.xplac.GetModule())
	s.Require().Equal(mcrisis.CrisisInvariantBrokenMsgType, s.xplac.GetMsgType())

	crisisInvariantBrokenTxbytes, err := s.xplac.InvariantBroken(invariantBrokenMsg).CreateAndSignTx()
	s.Require().NoError(err)

	crisisInvariantBrokenJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(crisisInvariantBrokenTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.CrisisInvariantBrokenTxTemplates, string(crisisInvariantBrokenJsonTxbytes))
}
