package bank_test

import (
	"math/rand"

	"github.com/xpladev/xpla.go/core/bank"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := bank.NewCoreModule()

	// test get name
	s.Require().Equal(bank.BankModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// bank send
	bankSendMsg := types.BankSendMsg{
		FromAddress: accounts[0].Address.String(),
		ToAddress:   accounts[1].Address.String(),
		Amount:      "1000",
	}

	makeBankSendMsg, err := bank.MakeBankSendMsg(bankSendMsg)
	s.Require().NoError(err)

	testMsg = makeBankSendMsg
	txBuilder, err = c.NewTxRouter(txBuilder, bank.BankSendMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeBankSendMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = provider.ResetXplac(s.xplac)
}
