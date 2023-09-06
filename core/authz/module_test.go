package authz_test

import (
	"math/rand"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/core/authz"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := authz.NewCoreModule()

	// test get name
	s.Require().Equal(authz.AuthzModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// authz grant
	authzGrantMsg := types.AuthzGrantMsg{
		Granter:           accounts[0].Address.String(),
		Grantee:           accounts[1].Address.String(),
		AuthorizationType: "send",
		SpendLimit:        "1000",
	}

	makeAuthzGrantMsg, err := authz.MakeAuthzGrantMsg(authzGrantMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeAuthzGrantMsg
	txBuilder, err = c.NewTxRouter(txBuilder, authz.AuthzGrantMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeAuthzGrantMsg, txBuilder.GetTx().GetMsgs()[0])

	// authz revoke
	authzRevokeMsg := types.AuthzRevokeMsg{
		Granter: accounts[0].Address.String(),
		Grantee: accounts[1].Address.String(),
		MsgType: "/cosmos.bank.v1beta1.MsgSend",
	}

	makeAuthzRevokeMsg, err := authz.MakeAuthzRevokeMsg(authzRevokeMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeAuthzRevokeMsg
	txBuilder, err = c.NewTxRouter(txBuilder, authz.AuthzRevokeMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeAuthzRevokeMsg, txBuilder.GetTx().GetMsgs()[0])

	// authz exec
	// e.g. bank send
	bankSendMsg := types.BankSendMsg{
		FromAddress: accounts[0].Address.String(),
		ToAddress:   accounts[1].Address.String(),
		Amount:      "1000",
	}

	txbytesBankSend, err := s.xplac.BankSend(bankSendMsg).CreateAndSignTx()
	s.Require().NoError(err)

	bankSendJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(txbytesBankSend)
	s.Require().NoError(err)

	authzExecMsg := types.AuthzExecMsg{
		Grantee:      accounts[1].Address.String(),
		ExecTxString: string(bankSendJsonTxbytes),
	}

	makeAuthzExecMsg, err := authz.MakeAuthzExecMsg(authzExecMsg, s.xplac.GetEncoding())
	s.Require().NoError(err)

	testMsg = makeAuthzExecMsg
	txBuilder, err = c.NewTxRouter(txBuilder, authz.AuthzExecMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeAuthzExecMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = client.ResetXplac(s.xplac)
}
