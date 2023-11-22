package feegrant_test

import (
	"math/rand"

	"github.com/xpladev/xpla.go/core/feegrant"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := feegrant.NewCoreModule()

	// test get name
	s.Require().Equal(feegrant.FeegrantModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// feegrant
	feeGrantMsg := types.FeeGrantMsg{
		Granter:    accounts[0].Address.String(),
		Grantee:    accounts[1].Address.String(),
		SpendLimit: "1000",
		// Period:      "3600",
		// PeriodLimit: "10",
		Expiration: "2100-01-01T23:59:59+00:00",
	}

	makeFeeGrantMsg, err := feegrant.MakeFeeGrantMsg(feeGrantMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	testMsg = makeFeeGrantMsg
	txBuilder, err = c.NewTxRouter(txBuilder, feegrant.FeegrantGrantMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeFeeGrantMsg, txBuilder.GetTx().GetMsgs()[0])

	// revoke feegrant
	revokeFeeGrantMsg := types.RevokeFeeGrantMsg{
		Granter: accounts[0].Address.String(),
		Grantee: accounts[1].Address.String(),
	}

	makeRevokeFeeGrantMsg, err := feegrant.MakeRevokeFeeGrantMsg(revokeFeeGrantMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	testMsg = makeRevokeFeeGrantMsg
	txBuilder, err = c.NewTxRouter(txBuilder, feegrant.FeegrantRevokeGrantMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeRevokeFeeGrantMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = provider.ResetXplac(s.xplac)
}
