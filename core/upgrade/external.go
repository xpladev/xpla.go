package upgrade

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &UpgradeExternal{}

type UpgradeExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e UpgradeExternal) {
	e.Xplac = xplac
	e.Name = UpgradeModule
	return e
}

func (e UpgradeExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e UpgradeExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Submit a software upgrade proposal.
func (e UpgradeExternal) SoftwareUpgrade(softwareUpgradeMsg types.SoftwareUpgradeMsg) provider.XplaClient {
	msg, err := MakeProposalSoftwareUpgradeMsg(softwareUpgradeMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(UpgradeProposalSoftwareUpgradeMsgType, err)
	}

	return e.ToExternal(UpgradeProposalSoftwareUpgradeMsgType, msg)
}

// Cancel the current software upgrade proposal.
func (e UpgradeExternal) CancelSoftwareUpgrade(cancelSoftwareUpgradeMsg types.CancelSoftwareUpgradeMsg) provider.XplaClient {
	msg, err := MakeCancelSoftwareUpgradeMsg(cancelSoftwareUpgradeMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(UpgradeCancelSoftwareUpgradeMsgType, err)
	}

	return e.ToExternal(UpgradeCancelSoftwareUpgradeMsgType, msg)
}

// Query

// Block header for height at which a completed upgrade was applied.
func (e UpgradeExternal) UpgradeApplied(appliedMsg types.AppliedMsg) provider.XplaClient {
	msg, err := MakeAppliedMsg(appliedMsg)
	if err != nil {
		return e.Err(UpgradeAppliedMsgType, err)
	}

	return e.ToExternal(UpgradeAppliedMsgType, msg)
}

// Query the list of module versions.
func (e UpgradeExternal) ModulesVersion(queryModulesVersionMsg ...types.QueryModulesVersionMsg) provider.XplaClient {
	switch {
	case len(queryModulesVersionMsg) == 0:
		msg, err := MakeQueryAllModuleVersionMsg()
		if err != nil {
			return e.Err(UpgradeQueryAllModuleVersionsMsgType, err)
		}

		return e.ToExternal(UpgradeQueryAllModuleVersionsMsgType, msg)

	case len(queryModulesVersionMsg) == 1:
		msg, err := MakeQueryModuleVersionMsg(queryModulesVersionMsg[0])
		if err != nil {
			return e.Err(UpgradeQueryModuleVersionsMsgType, err)
		}

		return e.ToExternal(UpgradeQueryModuleVersionsMsgType, msg)

	default:
		return e.Err(UpgradeQueryModuleVersionsMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}

// Query upgrade plan(if one exists).
func (e UpgradeExternal) Plan() provider.XplaClient {
	msg, err := MakePlanMsg()
	if err != nil {
		return e.Err(UpgradePlanMsgType, err)
	}

	return e.ToExternal(UpgradePlanMsgType, msg)
}
