package upgrade

import (
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

// Parsing - software upgrade
func parseProposalSoftwareUpgradeArgs(softwareUpgradeMsg types.SoftwareUpgradeMsg, from sdk.AccAddress) (govtypes.MsgSubmitProposal, error) {
	heightI64, err := util.FromStringToInt64(softwareUpgradeMsg.UpgradeHeight)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}
	plan := upgradetypes.Plan{
		Name:   softwareUpgradeMsg.UpgradeName,
		Height: heightI64,
		Info:   softwareUpgradeMsg.UpgradeInfo,
	}
	content := upgradetypes.NewSoftwareUpgradeProposal(
		softwareUpgradeMsg.Title,
		softwareUpgradeMsg.Description,
		plan,
	)

	deposit, err := sdk.ParseCoinsNormalized(util.DenomAdd(softwareUpgradeMsg.Deposit))
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	return *msg, nil
}

// Parsing - cancel software upgrade
func parseCancelSoftwareUpgradeArgs(cancelSoftwareUpgradeMsg types.CancelSoftwareUpgradeMsg, from sdk.AccAddress) (govtypes.MsgSubmitProposal, error) {
	deposit, err := sdk.ParseCoinsNormalized(util.DenomAdd(cancelSoftwareUpgradeMsg.Deposit))
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}
	content := upgradetypes.NewCancelSoftwareUpgradeProposal(
		cancelSoftwareUpgradeMsg.Deposit,
		cancelSoftwareUpgradeMsg.Description,
	)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	return *msg, nil
}
