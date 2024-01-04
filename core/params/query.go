package params

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	paramsv1beta1 "cosmossdk.io/api/cosmos/params/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
)

var out []byte
var res proto.Message
var err error

// Query client for params module.
func QueryParams(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcParams(i)
	} else {
		return queryByLcdParams(i)
	}

}

func queryByGrpcParams(i core.QueryClient) (string, error) {
	queryClient := proposal.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Params subspace
	case i.Ixplac.GetMsgType() == ParamsQuerySubpsaceMsgType:
		convertMsg := i.Ixplac.GetMsg().(proposal.QueryParamsRequest)
		res, err = queryClient.Params(
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
	paramsParamsLabel = "params"
)

func queryByLcdParams(i core.QueryClient) (string, error) {
	url := util.MakeQueryLcdUrl(paramsv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Params subspace
	case i.Ixplac.GetMsgType() == ParamsQuerySubpsaceMsgType:
		convertMsg := i.Ixplac.GetMsg().(proposal.QueryParamsRequest)

		parsedSubspace := convertMsg.Subspace
		parsedKey := convertMsg.Key

		subspace := "?subspace=" + parsedSubspace
		key := "&key=" + parsedKey

		url = url + paramsParamsLabel + subspace + key

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
