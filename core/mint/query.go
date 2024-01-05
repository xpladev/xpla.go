package mint

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	mintv1beta1 "cosmossdk.io/api/cosmos/mint/v1beta1"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

var out []byte
var res proto.Message
var err error

// Query client for mint module.
func QueryMint(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcMint(i)
	} else {
		return queryByLcdMint(i)
	}
}

func queryByGrpcMint(i core.QueryClient) (string, error) {
	queryClient := minttypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Mint parameters
	case i.Ixplac.GetMsgType() == MintQueryMintParamsMsgType:
		convertMsg := i.Ixplac.GetMsg().(minttypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Mint inflation
	case i.Ixplac.GetMsgType() == MintQueryInflationMsgType:
		convertMsg := i.Ixplac.GetMsg().(minttypes.QueryInflationRequest)
		res, err = queryClient.Inflation(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Mint annual provisions
	case i.Ixplac.GetMsgType() == MintQueryAnnualProvisionsMsgType:
		convertMsg := i.Ixplac.GetMsg().(minttypes.QueryAnnualProvisionsRequest)
		res, err = queryClient.AnnualProvisions(
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
	mintParamsLabel           = "params"
	mintInflationLabel        = "inflation"
	mintAnnualProvisionsLabel = "annual_provisions"
)

func queryByLcdMint(i core.QueryClient) (string, error) {
	url := util.MakeQueryLcdUrl(mintv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Mint parameters
	case i.Ixplac.GetMsgType() == MintQueryMintParamsMsgType:
		url = url + mintParamsLabel

	// Mint inflation
	case i.Ixplac.GetMsgType() == MintQueryInflationMsgType:
		url = url + mintInflationLabel

	// Mint annual provisions
	case i.Ixplac.GetMsgType() == MintQueryAnnualProvisionsMsgType:
		url = url + mintAnnualProvisionsLabel

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
