package reward

import (
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	rewardtypes "github.com/xpladev/xpla/x/reward/types"
)

// parsing - fund fee collector
func parseFundFeeCollectorArgs(fundFeeCollectorMsg types.FundFeeCollectorMsg, from sdk.AccAddress) (rewardtypes.MsgFundFeeCollector, error) {
	if fundFeeCollectorMsg.DepositorAddr != from.String() {
		return rewardtypes.MsgFundFeeCollector{}, util.LogErr(errors.ErrAccountNotMatch, "wrong depositor address, not match private key")
	}

	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(fundFeeCollectorMsg.Amount))
	if err != nil {
		return rewardtypes.MsgFundFeeCollector{}, util.LogErr(errors.ErrParse, err)
	}

	addr, err := sdk.AccAddressFromBech32(fundFeeCollectorMsg.DepositorAddr)
	if err != nil {
		return rewardtypes.MsgFundFeeCollector{}, util.LogErr(errors.ErrParse, err)
	}

	msg := rewardtypes.NewMsgFundFeeCollector(amount, addr)

	return *msg, nil
}
