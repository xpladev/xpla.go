package auth_test

import (
	mauth "github.com/xpladev/xpla.go/core/auth"
	"github.com/xpladev/xpla.go/types"
)

func (s *IntegrationTestSuite) TestAuth() {
	// auth params
	s.xplac.AuthParams()

	authParamMsg, err := mauth.MakeAuthParamMsg()
	s.Require().NoError(err)

	s.Require().Equal(authParamMsg, s.xplac.GetMsg())
	s.Require().Equal(mauth.AuthModule, s.xplac.GetModule())
	s.Require().Equal(mauth.AuthQueryParamsMsgType, s.xplac.GetMsgType())

	// acc address
	queryAccAddressMsg := types.QueryAccAddressMsg{
		Address: s.network.Validators[0].AdditionalAccount.Address.String(),
	}
	s.xplac.AccAddress(queryAccAddressMsg)

	accAddressMsg, err := mauth.MakeQueryAccAddressMsg(queryAccAddressMsg)
	s.Require().NoError(err)

	s.Require().Equal(accAddressMsg, s.xplac.GetMsg())
	s.Require().Equal(mauth.AuthModule, s.xplac.GetModule())
	s.Require().Equal(mauth.AuthQueryAccAddressMsgType, s.xplac.GetMsgType())

	// accounts
	s.xplac.Accounts()

	accountsMsg, err := mauth.MakeQueryAccountsMsg()
	s.Require().NoError(err)

	s.Require().Equal(accountsMsg, s.xplac.GetMsg())
	s.Require().Equal(mauth.AuthModule, s.xplac.GetModule())
	s.Require().Equal(mauth.AuthQueryAccountsMsgType, s.xplac.GetMsgType())

	// txs by events
	queryTxsByEventsMsg := types.QueryTxsByEventsMsg{
		Events: "transfer.recipient=" + s.network.Validators[0].AdditionalAccount.Address.String(),
	}
	s.xplac.TxsByEvents(queryTxsByEventsMsg)
	txsByEventMsg, err := mauth.MakeTxsByEventsMsg(queryTxsByEventsMsg)
	s.Require().NoError(err)

	s.Require().Equal(txsByEventMsg, s.xplac.GetMsg())
	s.Require().Equal(mauth.AuthModule, s.xplac.GetModule())
	s.Require().Equal(mauth.AuthQueryTxsByEventsMsgType, s.xplac.GetMsgType())

	// tx
	queryTxMsg := types.QueryTxMsg{
		Value: s.testTxHash,
	}
	s.xplac.Tx(queryTxMsg)

	txMsg, err := mauth.MakeQueryTxMsg(queryTxMsg)
	s.Require().NoError(err)

	s.Require().Equal(txMsg, s.xplac.GetMsg())
	s.Require().Equal(mauth.AuthModule, s.xplac.GetModule())
	s.Require().Equal(mauth.AuthQueryTxMsgType, s.xplac.GetMsgType())
}
