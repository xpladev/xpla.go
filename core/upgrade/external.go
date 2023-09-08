package upgrade

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type UpgradeExternal struct {
	Xplac provider.XplaClient
}

func NewUpgradeExternal(xplac provider.XplaClient) (e UpgradeExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Submit a software upgrade proposal.
func (e UpgradeExternal) SoftwareUpgrade(softwareUpgradeMsg types.SoftwareUpgradeMsg) provider.XplaClient {
	msg, err := MakeProposalSoftwareUpgradeMsg(softwareUpgradeMsg, e.Xplac.GetPrivateKey())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(UpgradeModule).
		WithMsgType(UpgradeProposalSoftwareUpgradeMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Cancel the current software upgrade proposal.
func (e UpgradeExternal) CancelSoftwareUpgrade(cancelSoftwareUpgradeMsg types.CancelSoftwareUpgradeMsg) provider.XplaClient {
	msg, err := MakeCancelSoftwareUpgradeMsg(cancelSoftwareUpgradeMsg, e.Xplac.GetPrivateKey())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(UpgradeModule).
		WithMsgType(UpgradeCancelSoftwareUpgradeMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Block header for height at which a completed upgrade was applied.
func (e UpgradeExternal) UpgradeApplied(appliedMsg types.AppliedMsg) provider.XplaClient {
	msg, err := MakeAppliedMsg(appliedMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(UpgradeModule).
		WithMsgType(UpgradeAppliedMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query the list of module versions.
func (e UpgradeExternal) ModulesVersion(queryModulesVersionMsg ...types.QueryModulesVersionMsg) provider.XplaClient {
	if len(queryModulesVersionMsg) == 0 {
		msg, err := MakeQueryAllModuleVersionMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(UpgradeModule).
			WithMsgType(UpgradeQueryAllModuleVersionsMsgType).
			WithMsg(msg)
	} else if len(queryModulesVersionMsg) == 1 {
		msg, err := MakeQueryModuleVersionMsg(queryModulesVersionMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(UpgradeModule).
			WithMsgType(UpgradeQueryModuleVersionsMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}

// Query upgrade plan(if one exists).
func (e UpgradeExternal) Plan() provider.XplaClient {
	msg, err := MakePlanMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(UpgradeModule).
		WithMsgType(UpgradePlanMsgType).
		WithMsg(msg)
	return e.Xplac
}
