package upgrade_test

import (
	"math/rand"

	mupgrade "github.com/xpladev/xpla.go/core/upgrade"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestUpgradeTx() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)

	s.xplac.WithPrivateKey(accounts[0].PrivKey)
	// software upgrade
	softwareUpgradeMsg := types.SoftwareUpgradeMsg{
		UpgradeName:   "Upgrade Name",
		Title:         "Upgrade Title",
		Description:   "Upgrade Description",
		UpgradeHeight: "6000",
		UpgradeInfo:   `{"upgrade_info":"INFO"}`,
		Deposit:       "1000",
	}
	s.xplac.SoftwareUpgrade(softwareUpgradeMsg)

	makeProposalSoftwareUpgradeMsg, err := mupgrade.MakeProposalSoftwareUpgradeMsg(softwareUpgradeMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeProposalSoftwareUpgradeMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradeProposalSoftwareUpgradeMsgType, s.xplac.GetMsgType())

	upgradeSoftwareUpgradeTxbytes, err := s.xplac.SoftwareUpgrade(softwareUpgradeMsg).CreateAndSignTx()
	s.Require().NoError(err)

	upgradeSoftwareUpgradeJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(upgradeSoftwareUpgradeTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.UpgradeSoftwareUpgradeTxTemplates, string(upgradeSoftwareUpgradeJsonTxbytes))

	// cancel software upgrade
	cancelSoftwareUpgradeMsg := types.CancelSoftwareUpgradeMsg{
		Title:       "Cancel software upgrade",
		Description: "Cancel software upgrade description",
		Deposit:     "1000",
	}
	s.xplac.CancelSoftwareUpgrade(cancelSoftwareUpgradeMsg)

	makeCancelSoftwareUpgradeMsg, err := mupgrade.MakeCancelSoftwareUpgradeMsg(cancelSoftwareUpgradeMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeCancelSoftwareUpgradeMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradeCancelSoftwareUpgradeMsgType, s.xplac.GetMsgType())

	upgradeCancelSoftwareUpgradeTxbytes, err := s.xplac.CancelSoftwareUpgrade(cancelSoftwareUpgradeMsg).CreateAndSignTx()
	s.Require().NoError(err)

	upgradeCancelSoftwareUpgradeJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(upgradeCancelSoftwareUpgradeTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.UpgradeCancelSoftwareUpgradeTxTemplates, string(upgradeCancelSoftwareUpgradeJsonTxbytes))

}

func (s *IntegrationTestSuite) TestUpgrade() {
	// upgrade applied
	appliedMsg := types.AppliedMsg{
		UpgradeName: "upgrade name",
	}
	s.xplac.UpgradeApplied(appliedMsg)

	makeAppliedMsg, err := mupgrade.MakeAppliedMsg(appliedMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeAppliedMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradeAppliedMsgType, s.xplac.GetMsgType())

	// modules version
	s.xplac.ModulesVersion()

	makeQueryAllModuleVersionMsg, err := mupgrade.MakeQueryAllModuleVersionMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryAllModuleVersionMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradeQueryAllModuleVersionsMsgType, s.xplac.GetMsgType())

	// module version
	queryModulesVersionMsg := types.QueryModulesVersionMsg{
		ModuleName: "staking",
	}
	s.xplac.ModulesVersion(queryModulesVersionMsg)

	makeQueryModuleVersionMsg, err := mupgrade.MakeQueryModuleVersionMsg(queryModulesVersionMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryModuleVersionMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradeQueryModuleVersionsMsgType, s.xplac.GetMsgType())

	// plan
	s.xplac.Plan()

	makePlanMsg, err := mupgrade.MakePlanMsg()
	s.Require().NoError(err)

	s.Require().Equal(makePlanMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradePlanMsgType, s.xplac.GetMsgType())
}
