package slashing_test

import (
	"math/rand"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/core/slashing"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := slashing.NewCoreModule()

	// test get name
	s.Require().Equal(slashing.SlashingModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// unjail
	s.xplac.Unjail()

	makeUnjailMsg, err := slashing.MakeUnjailMsg(s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeUnjailMsg
	txBuilder, err = c.NewTxRouter(txBuilder, slashing.SlahsingUnjailMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeUnjailMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = client.ResetXplac(s.xplac)
}
