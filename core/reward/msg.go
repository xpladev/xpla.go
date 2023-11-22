package reward

import (
	"github.com/xpladev/xpla.go/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	rewardtypes "github.com/xpladev/xpla/x/reward/types"
)

// (Tx) make msg - Fund fee collector
func MakeFundFeeCollectorMsg(fundFeeCollectorMsg types.FundFeeCollectorMsg, from sdk.AccAddress) (rewardtypes.MsgFundFeeCollector, error) {
	return parseFundFeeCollectorArgs(fundFeeCollectorMsg, from)
}

// (Query) make msg - query reward params
func MakeQueryRewardParamsMsg() (rewardtypes.QueryParamsRequest, error) {
	return rewardtypes.QueryParamsRequest{}, nil
}

// (Query) make msg - query reward pool
func MakeQueryRewardPoolMsg() (rewardtypes.QueryPoolRequest, error) {
	return rewardtypes.QueryPoolRequest{}, nil
}
