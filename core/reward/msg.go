package reward

import (
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"

	rewardtypes "github.com/xpladev/xpla/x/reward/types"
)

// (Tx) make msg - Fund fee collector
func MakeFundFeeCollectorMsg(fundFeeCollectorMsg types.FundFeeCollectorMsg, privKey key.PrivateKey) (rewardtypes.MsgFundFeeCollector, error) {
	return parseFundFeeCollectorArgs(fundFeeCollectorMsg, privKey)
}

// (Query) make msg - query reward params
func MakeQueryRewardParamsMsg() (rewardtypes.QueryParamsRequest, error) {
	return rewardtypes.QueryParamsRequest{}, nil
}

// (Query) make msg - query reward pool
func MakeQueryRewardPoolMsg() (rewardtypes.QueryPoolRequest, error) {
	return rewardtypes.QueryPoolRequest{}, nil
}
