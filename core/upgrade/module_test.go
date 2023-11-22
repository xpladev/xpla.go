package upgrade_test

import (
	"math/rand"

	"github.com/xpladev/xpla.go/core/upgrade"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := upgrade.NewCoreModule()

	// test get name
	s.Require().Equal(upgrade.UpgradeModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// software upgrade
	softwareUpgradeMsg := types.SoftwareUpgradeMsg{
		UpgradeName:   "Upgrade Name",
		Title:         "Upgrade Title",
		Description:   "Upgrade Description",
		UpgradeHeight: "6000",
		UpgradeInfo:   `{"upgrade_info":"INFO"}`,
		Deposit:       "1000",
	}

	makeProposalSoftwareUpgradeMsg, err := upgrade.MakeProposalSoftwareUpgradeMsg(softwareUpgradeMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	testMsg = makeProposalSoftwareUpgradeMsg
	txBuilder, err = c.NewTxRouter(txBuilder, upgrade.UpgradeProposalSoftwareUpgradeMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeProposalSoftwareUpgradeMsg, txBuilder.GetTx().GetMsgs()[0])

	// cancel software upgrade
	cancelSoftwareUpgradeMsg := types.CancelSoftwareUpgradeMsg{
		Title:       "Cancel software upgrade",
		Description: "Cancel software upgrade description",
		Deposit:     "1000",
	}

	makeCancelSoftwareUpgradeMsg, err := upgrade.MakeCancelSoftwareUpgradeMsg(cancelSoftwareUpgradeMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	testMsg = makeCancelSoftwareUpgradeMsg
	txBuilder, err = c.NewTxRouter(txBuilder, upgrade.UpgradeCancelSoftwareUpgradeMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeCancelSoftwareUpgradeMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = provider.ResetXplac(s.xplac)
}
