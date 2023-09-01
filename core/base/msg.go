package base

import (
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
)

// (Query) make msg - node info
func MakeBaseNodeInfoMsg() (tmservice.GetNodeInfoRequest, error) {
	return tmservice.GetNodeInfoRequest{}, nil
}

// (Query) make msg - syncing
func MakeBaseSyncingMsg() (tmservice.GetSyncingRequest, error) {
	return tmservice.GetSyncingRequest{}, nil
}

// (Query) make msg - latest block
func MakeBaseLatestBlockMsg() (tmservice.GetLatestBlockRequest, error) {
	return tmservice.GetLatestBlockRequest{}, nil
}

// (Query) make msg - get block by height
func MakeBaseBlockByHeightMsg(blockMsg types.BlockMsg) (tmservice.GetBlockByHeightRequest, error) {
	heighti64, err := util.FromStringToInt64(blockMsg.Height)
	if err != nil {
		return tmservice.GetBlockByHeightRequest{}, err
	}
	return tmservice.GetBlockByHeightRequest{
		Height: heighti64,
	}, nil
}

// (Query) make msg - latest validator set
func MakeLatestValidatorSetMsg() (tmservice.GetLatestValidatorSetRequest, error) {
	return tmservice.GetLatestValidatorSetRequest{}, nil
}

// (Query) make msg - latest validator set
func MakeValidatorSetByHeightMsg(validatorSetMsg types.ValidatorSetMsg) (tmservice.GetValidatorSetByHeightRequest, error) {
	heighti64, err := util.FromStringToInt64(validatorSetMsg.Height)
	if err != nil {
		return tmservice.GetValidatorSetByHeightRequest{}, err
	}
	return tmservice.GetValidatorSetByHeightRequest{
		Height: heighti64,
	}, nil
}
