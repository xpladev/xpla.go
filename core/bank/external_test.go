package bank_test

import (
	mbank "github.com/xpladev/xpla.go/core/bank"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestBankTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// bank send
	bankSendMsg := types.BankSendMsg{
		FromAddress: s.accounts[0].Address.String(),
		ToAddress:   s.accounts[1].Address.String(),
		Amount:      "1000",
	}
	s.xplac.BankSend(bankSendMsg)

	makeBankSendMsg, err := mbank.MakeBankSendMsg(bankSendMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeBankSendMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankSendMsgType, s.xplac.GetMsgType())

	bankSendTxbytes, err := s.xplac.BankSend(bankSendMsg).CreateAndSignTx()
	s.Require().NoError(err)

	bankSendJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(bankSendTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.BankSendTxTemplates, string(bankSendJsonTxbytes))
}

func (s *IntegrationTestSuite) TestBank() {
	// bank all balances
	bankBalancesMsg := types.BankBalancesMsg{
		Address: s.accounts[0].Address.String(),
	}
	s.xplac.BankBalances(bankBalancesMsg)

	makeBankAllBalancesMsg, err := mbank.MakeBankAllBalancesMsg(bankBalancesMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeBankAllBalancesMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankAllBalancesMsgType, s.xplac.GetMsgType())

	// bank balance denom
	bankBalancesMsg = types.BankBalancesMsg{
		Address: s.accounts[0].Address.String(),
		Denom:   types.XplaDenom,
	}
	s.xplac.BankBalances(bankBalancesMsg)

	makeBankBalanceMsg, err := mbank.MakeBankBalanceMsg(bankBalancesMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeBankBalanceMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankBalanceMsgType, s.xplac.GetMsgType())

	// denoms metadata
	s.xplac.DenomMetadata()

	makeDenomsMetaDataMsg, err := mbank.MakeDenomsMetaDataMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeDenomsMetaDataMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankDenomsMetadataMsgType, s.xplac.GetMsgType())

	// denom metadata
	denomMetadataMsg := types.DenomMetadataMsg{
		Denom: types.XplaDenom,
	}
	s.xplac.DenomMetadata(denomMetadataMsg)

	makeDenomMetaDataMsg, err := mbank.MakeDenomMetaDataMsg(denomMetadataMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeDenomMetaDataMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankDenomMetadataMsgType, s.xplac.GetMsgType())

	// total supply
	s.xplac.Total()

	makeTotalSupplyMsg, err := mbank.MakeTotalSupplyMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeTotalSupplyMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankTotalMsgType, s.xplac.GetMsgType())

	// supply of
	totalMsg := types.TotalMsg{
		Denom: types.XplaDenom,
	}
	s.xplac.Total(totalMsg)

	makeSupplyOfMsg, err := mbank.MakeSupplyOfMsg(totalMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeSupplyOfMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankTotalSupplyOfMsgType, s.xplac.GetMsgType())
}
