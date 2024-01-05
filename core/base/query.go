package base

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	tmv1beta1 "cosmossdk.io/api/cosmos/base/tendermint/v1beta1"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
)

var out []byte
var res proto.Message
var err error

// Query client for bank module.
func QueryBase(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcBase(i)
	} else {
		return queryByLcdBase(i)
	}

}

func queryByGrpcBase(i core.QueryClient) (string, error) {
	serviceClient := tmservice.NewServiceClient(i.Ixplac.GetGrpcClient())

	switch {
	// Node info
	case i.Ixplac.GetMsgType() == BaseNodeInfoMsgType:
		convertMsg := i.Ixplac.GetMsg().(tmservice.GetNodeInfoRequest)
		res, err = serviceClient.GetNodeInfo(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Syncing
	case i.Ixplac.GetMsgType() == BaseSyncingMsgType:
		convertMsg := i.Ixplac.GetMsg().(tmservice.GetSyncingRequest)
		res, err = serviceClient.GetSyncing(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Latest block
	case i.Ixplac.GetMsgType() == BaseLatestBlockMsgtype:
		if i.Ixplac.GetRpc() != "" {
			var height *int64
			return queryBlockByRpc(i, height)

		} else {
			convertMsg := i.Ixplac.GetMsg().(tmservice.GetLatestBlockRequest)
			res, err = serviceClient.GetLatestBlock(
				i.Ixplac.GetContext(),
				&convertMsg,
			)
			if err != nil {
				return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
			}
		}

	// Block by height
	case i.Ixplac.GetMsgType() == BaseBlockByHeightMsgType:
		convertMsg := i.Ixplac.GetMsg().(tmservice.GetBlockByHeightRequest)
		if i.Ixplac.GetRpc() != "" {
			height := &convertMsg.Height
			return queryBlockByRpc(i, height)

		} else {
			res, err = serviceClient.GetBlockByHeight(
				i.Ixplac.GetContext(),
				&convertMsg,
			)
			if err != nil {
				return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
			}
		}

	// Latest validator set
	case i.Ixplac.GetMsgType() == BaseLatestValidatorSetMsgType:
		convertMsg := i.Ixplac.GetMsg().(tmservice.GetLatestValidatorSetRequest)
		res, err = serviceClient.GetLatestValidatorSet(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Validator set by height
	case i.Ixplac.GetMsgType() == BaseValidatorSetByHeightMsgType:
		convertMsg := i.Ixplac.GetMsg().(tmservice.GetValidatorSetByHeightRequest)
		res, err = serviceClient.GetValidatorSetByHeight(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	default:
		return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrInvalidMsgType, i.Ixplac.GetMsgType()))
	}

	out, err = core.PrintProto(i, res)
	if err != nil {
		return "", i.Ixplac.GetLogger().Err(err)
	}

	return string(out), nil
}

const (
	baseNodeInfoLabel      = "node_info"
	baseSyncingLabel       = "syncing"
	baseBlocksLabel        = "blocks"
	baseLatestLabel        = "latest"
	baseValidatorsetsLabel = "validatorsets"
)

func queryByLcdBase(i core.QueryClient) (string, error) {
	url := util.MakeQueryLcdUrl(tmv1beta1.Service_ServiceDesc.Metadata.(string))

	switch {
	// Node info
	case i.Ixplac.GetMsgType() == BaseNodeInfoMsgType:
		url = url + baseNodeInfoLabel

	// Syncing
	case i.Ixplac.GetMsgType() == BaseSyncingMsgType:
		url = url + baseSyncingLabel

	// Latest block
	case i.Ixplac.GetMsgType() == BaseLatestBlockMsgtype:
		url = util.MakeQueryLabels("/", baseBlocksLabel, baseLatestLabel)

	// Block by height
	case i.Ixplac.GetMsgType() == BaseBlockByHeightMsgType:
		convertMsg := i.Ixplac.GetMsg().(tmservice.GetBlockByHeightRequest)
		url = util.MakeQueryLabels("/", baseBlocksLabel, util.FromInt64ToString(convertMsg.Height))

	// Latest validator set
	case i.Ixplac.GetMsgType() == BaseLatestValidatorSetMsgType:
		url = url + util.MakeQueryLabels(baseValidatorsetsLabel, baseLatestLabel)

	// Validator set by height
	case i.Ixplac.GetMsgType() == BaseValidatorSetByHeightMsgType:
		convertMsg := i.Ixplac.GetMsg().(tmservice.GetValidatorSetByHeightRequest)
		url = url + util.MakeQueryLabels(baseValidatorsetsLabel, util.FromInt64ToString(convertMsg.Height))

	default:
		return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrInvalidMsgType, i.Ixplac.GetMsgType()))
	}

	i.Ixplac.GetHttpMutex().Lock()
	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		i.Ixplac.GetHttpMutex().Unlock()
		return "", i.Ixplac.GetLogger().Err(err)
	}
	i.Ixplac.GetHttpMutex().Unlock()

	return string(out), nil
}

func queryBlockByRpc(i core.QueryClient, height *int64) (string, error) {
	client, err := cmclient.NewClientFromNode(i.Ixplac.GetRpc())
	if err != nil {
		return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
	}
	res, err := client.Block(i.Ixplac.GetContext(), height)
	if err != nil {
		return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
	}
	out, err := core.PrintObjectLegacy(i, res)
	if err != nil {
		return "", i.Ixplac.GetLogger().Err(err)
	}
	return string(out), nil
}
