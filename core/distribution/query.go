package distribution

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	distv1beta1 "cosmossdk.io/api/cosmos/distribution/v1beta1"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
)

var out []byte
var res proto.Message
var err error

// Query client for distribution module.
func QueryDistribution(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcDist(i)
	} else {
		return queryByLcdDist(i)
	}
}

func queryByGrpcDist(i core.QueryClient) (string, error) {
	queryClient := disttypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Distribution params
	case i.Ixplac.GetMsgType() == DistributionQueryDistributionParamsMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Distribution validator outstanding rewards
	case i.Ixplac.GetMsgType() == DistributionValidatorOutstandingRewardsMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryValidatorOutstandingRewardsRequest)
		res, err = queryClient.ValidatorOutstandingRewards(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Distribution commission
	case i.Ixplac.GetMsgType() == DistributionQueryDistCommissionMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryValidatorCommissionRequest)
		res, err = queryClient.ValidatorCommission(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Distribution slashes
	case i.Ixplac.GetMsgType() == DistributionQuerySlashesMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryValidatorSlashesRequest)
		res, err = queryClient.ValidatorSlashes(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Distribution rewards
	case i.Ixplac.GetMsgType() == DistributionQueryRewardsMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryDelegationRewardsRequest)
		res, err = queryClient.DelegationRewards(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Distribution total rewards
	case i.Ixplac.GetMsgType() == DistributionQueryTotalRewardsMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryDelegationTotalRewardsRequest)
		res, err = queryClient.DelegationTotalRewards(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Distribution community pool
	case i.Ixplac.GetMsgType() == DistributionQueryCommunityPoolMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryCommunityPoolRequest)
		res, err = queryClient.CommunityPool(
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
	distParamsLabel             = "params"
	distValidatorLabel          = "validators"
	distDelegatorLabel          = "delegators"
	distOutstandingRewardsLabel = "outstanding_rewards"
	distCommissionLabel         = "commission"
	distSlashesLabel            = "slashes"
	distRewardsLabel            = "rewards"
	distCommunityPoolLabel      = "community_pool"
)

func queryByLcdDist(i core.QueryClient) (string, error) {
	url := util.MakeQueryLcdUrl(distv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Distribution params
	case i.Ixplac.GetMsgType() == DistributionQueryDistributionParamsMsgType:
		url = url + distParamsLabel

	// Distribution validator outstanding rewards
	case i.Ixplac.GetMsgType() == DistributionValidatorOutstandingRewardsMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryValidatorOutstandingRewardsRequest)

		url = url + util.MakeQueryLabels(distValidatorLabel, convertMsg.ValidatorAddress, distOutstandingRewardsLabel)

	// Distribution commission
	case i.Ixplac.GetMsgType() == DistributionQueryDistCommissionMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryValidatorCommissionRequest)

		url = url + util.MakeQueryLabels(distValidatorLabel, convertMsg.ValidatorAddress, distCommissionLabel)

	// Distribution slashes
	case i.Ixplac.GetMsgType() == DistributionQuerySlashesMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryValidatorSlashesRequest)

		url = url + util.MakeQueryLabels(distValidatorLabel, convertMsg.ValidatorAddress, distSlashesLabel)

	// Distribution rewards
	case i.Ixplac.GetMsgType() == DistributionQueryRewardsMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryDelegationRewardsRequest)

		url = url + util.MakeQueryLabels(distDelegatorLabel, convertMsg.DelegatorAddress, distRewardsLabel, convertMsg.ValidatorAddress)

	// Distribution total rewards
	case i.Ixplac.GetMsgType() == DistributionQueryTotalRewardsMsgType:
		convertMsg := i.Ixplac.GetMsg().(disttypes.QueryDelegationTotalRewardsRequest)

		url = url + util.MakeQueryLabels(distDelegatorLabel, convertMsg.DelegatorAddress, distRewardsLabel)

	// Distribution community pool
	case i.Ixplac.GetMsgType() == DistributionQueryCommunityPoolMsgType:
		url = url + distCommunityPoolLabel

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
