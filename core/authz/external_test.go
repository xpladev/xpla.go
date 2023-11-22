package authz_test

import (
	mauthz "github.com/xpladev/xpla.go/core/authz"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestAuthzTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// authz grant
	authzGrantMsg := types.AuthzGrantMsg{
		Granter:           s.accounts[0].Address.String(),
		Grantee:           s.accounts[1].Address.String(),
		AuthorizationType: "send",
		SpendLimit:        "1000",
	}
	s.xplac.AuthzGrant(authzGrantMsg)

	makeAuthzGrantMsg, err := mauthz.MakeAuthzGrantMsg(authzGrantMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeAuthzGrantMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzGrantMsgType, s.xplac.GetMsgType())

	authzGrantTxbytes, err := s.xplac.AuthzGrant(authzGrantMsg).CreateAndSignTx()
	s.Require().NoError(err)

	authzGrantjsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(authzGrantTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.AuthzGrantTxTemplates, string(authzGrantjsonTxbytes))

	// authz revoke
	authzRevokeMsg := types.AuthzRevokeMsg{
		Granter: s.accounts[0].Address.String(),
		Grantee: s.accounts[1].Address.String(),
		MsgType: "/cosmos.bank.v1beta1.MsgSend",
	}
	s.xplac.AuthzRevoke(authzRevokeMsg)

	makeAuthzRevokeMsg, err := mauthz.MakeAuthzRevokeMsg(authzRevokeMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeAuthzRevokeMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzRevokeMsgType, s.xplac.GetMsgType())

	authzRevokeTxbytes, err := s.xplac.AuthzRevoke(authzRevokeMsg).CreateAndSignTx()
	s.Require().NoError(err)

	authzRevokeJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(authzRevokeTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.AuthzRevokeTxTemplates, string(authzRevokeJsonTxbytes))

	// authz exec
	// e.g. bank send
	bankSendMsg := types.BankSendMsg{
		FromAddress: s.accounts[0].Address.String(),
		ToAddress:   s.accounts[1].Address.String(),
		Amount:      "1000",
	}

	txbytesBankSend, err := s.xplac.BankSend(bankSendMsg).CreateAndSignTx()
	s.Require().NoError(err)

	bankSendJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(txbytesBankSend)
	s.Require().NoError(err)

	authzExecMsg := types.AuthzExecMsg{
		Grantee:      s.accounts[1].Address.String(),
		ExecTxString: string(bankSendJsonTxbytes),
	}
	s.xplac.AuthzExec(authzExecMsg)

	makeAuthzExecMsg, err := mauthz.MakeAuthzExecMsg(authzExecMsg, s.xplac.GetEncoding())
	s.Require().NoError(err)

	s.Require().Equal(makeAuthzExecMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzExecMsgType, s.xplac.GetMsgType())

	authzExecTxbytes, err := s.xplac.AuthzExec(authzExecMsg).CreateAndSignTx()
	s.Require().NoError(err)

	authzExecJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(authzExecTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.AuthzExecTxTemplates, string(authzExecJsonTxbytes))
}

func (s *IntegrationTestSuite) TestAuthz() {
	// query authz grants
	queryAuthzGrantMsg := types.QueryAuthzGrantMsg{
		Grantee: s.accounts[0].Address.String(),
		Granter: s.accounts[1].Address.String(),
	}
	s.xplac.QueryAuthzGrants(queryAuthzGrantMsg)

	authzGrantsMsg, err := mauthz.MakeQueryAuthzGrantsMsg(queryAuthzGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(authzGrantsMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzQueryGrantMsgType, s.xplac.GetMsgType())

	// grants by grantee
	queryAuthzGrantMsg = types.QueryAuthzGrantMsg{
		Grantee: s.accounts[0].Address.String(),
	}
	s.xplac.QueryAuthzGrants(queryAuthzGrantMsg)

	authzGrantsByGranteeMsg, err := mauthz.MakeQueryAuthzGrantsByGranteeMsg(queryAuthzGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(authzGrantsByGranteeMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzQueryGrantsByGranteeMsgType, s.xplac.GetMsgType())

	// grants by granter
	queryAuthzGrantMsg = types.QueryAuthzGrantMsg{
		Granter: s.accounts[1].Address.String(),
	}
	s.xplac.QueryAuthzGrants(queryAuthzGrantMsg)

	authzGrantsByGranterMsg, err := mauthz.MakeQueryAuthzGrantsByGranterMsg(queryAuthzGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(authzGrantsByGranterMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzQueryGrantsByGranterMsgType, s.xplac.GetMsgType())
}
