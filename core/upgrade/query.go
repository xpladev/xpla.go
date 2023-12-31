package upgrade

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	upgradev1beta1 "cosmossdk.io/api/cosmos/upgrade/v1beta1"
	cmclient "github.com/cosmos/cosmos-sdk/client"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

var out []byte
var res proto.Message
var err error

// Query client for upgrade module.
func QueryUpgrade(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcUpgrade(i)
	} else {
		return queryByLcdUpgrade(i)
	}
}

func queryByGrpcUpgrade(i core.QueryClient) (string, error) {
	queryClient := upgradetypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Upgrade applied
	case i.Ixplac.GetMsgType() == UpgradeAppliedMsgType:
		convertMsg := i.Ixplac.GetMsg().(upgradetypes.QueryAppliedPlanRequest)
		appliedPlanRes, err := queryClient.AppliedPlan(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

		if appliedPlanRes.Height == 0 {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrNotFound, "applied plan height is 0"))
		}
		headerData, err := appliedReturnBlockheader(appliedPlanRes, i.Ixplac.GetRpc(), i.Ixplac.GetContext())
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(err)
		}
		return string(headerData), nil

	// Upgrade all module versions
	case i.Ixplac.GetMsgType() == UpgradeQueryAllModuleVersionsMsgType ||
		i.Ixplac.GetMsgType() == UpgradeQueryModuleVersionsMsgType:
		convertMsg := i.Ixplac.GetMsg().(upgradetypes.QueryModuleVersionsRequest)
		res, err = queryClient.ModuleVersions(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Upgrade plan
	case i.Ixplac.GetMsgType() == UpgradePlanMsgType:
		convertMsg := i.Ixplac.GetMsg().(upgradetypes.QueryCurrentPlanRequest)
		res, err = queryClient.CurrentPlan(
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
	upgradeAppliedPlanLabel    = "applied_plan"
	upgradeModuleVersionsLabel = "module_versions"
	upgradeCurrentPlanLabel    = "current_plan"
)

func queryByLcdUpgrade(i core.QueryClient) (string, error) {
	url := util.MakeQueryLcdUrl(upgradev1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Upgrade applied
	case i.Ixplac.GetMsgType() == UpgradeAppliedMsgType:
		convertMsg := i.Ixplac.GetMsg().(upgradetypes.QueryAppliedPlanRequest)

		url = url + util.MakeQueryLabels(upgradeAppliedPlanLabel, convertMsg.Name)

	// Upgrade all module versions
	case i.Ixplac.GetMsgType() == UpgradeQueryAllModuleVersionsMsgType ||
		i.Ixplac.GetMsgType() == UpgradeQueryModuleVersionsMsgType:

		url = url + upgradeModuleVersionsLabel

	// Upgrade plan
	case i.Ixplac.GetMsgType() == UpgradePlanMsgType:
		url = url + upgradeCurrentPlanLabel

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

func appliedReturnBlockheader(res *upgradetypes.QueryAppliedPlanResponse, rpcUrl string, ctx context.Context) ([]byte, error) {
	if rpcUrl == "" {
		return nil, types.ErrWrap(types.ErrNotSatisfiedOptions, "need RPC URL")
	}
	clientCtx, err := util.NewClient()
	if err != nil {
		return nil, err
	}

	client, err := cmclient.NewClientFromNode(rpcUrl)
	if err != nil {
		return nil, types.ErrWrap(types.ErrSdkClient, err)
	}
	clientCtx = clientCtx.WithClient(client)

	node, err := clientCtx.GetNode()
	if err != nil {
		return nil, types.ErrWrap(types.ErrSdkClient, err)
	}

	headers, err := node.BlockchainInfo(ctx, res.Height, res.Height)
	if err != nil {
		return nil, types.ErrWrap(types.ErrSdkClient, err)
	}

	if len(headers.BlockMetas) == 0 {
		return nil, types.ErrWrap(types.ErrNotFound, "no headers returns for height", res.Height)
	}

	bytes, err := clientCtx.LegacyAmino.MarshalJSONIndent(headers.BlockMetas[0], "", "  ")
	if err != nil {
		return nil, types.ErrWrap(types.ErrFailedToMarshal, err)
	}

	return bytes, nil
}
