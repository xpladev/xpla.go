package slashing

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	slashingv1beta1 "cosmossdk.io/api/cosmos/slashing/v1beta1"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
)

var out []byte
var res proto.Message
var err error

// Query client for slashing module.
func QuerySlashing(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcSlashing(i)
	} else {
		return queryByLcdSlashing(i)
	}
}

func queryByGrpcSlashing(i core.QueryClient) (string, error) {
	queryClient := slashingtypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Slashing parameters
	case i.Ixplac.GetMsgType() == SlashingQuerySlashingParamsMsgType:
		convertMsg := i.Ixplac.GetMsg().(slashingtypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Slashing signing information
	case i.Ixplac.GetMsgType() == SlashingQuerySigningInfosMsgType:
		convertMsg := i.Ixplac.GetMsg().(slashingtypes.QuerySigningInfosRequest)
		res, err = queryClient.SigningInfos(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Slashing signing information
	case i.Ixplac.GetMsgType() == SlashingQuerySigningInfoMsgType:
		convertMsg := i.Ixplac.GetMsg().(slashingtypes.QuerySigningInfoRequest)
		res, err = queryClient.SigningInfo(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = core.PrintProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	slashingParamsLabel       = "params"
	slashingSigningInfosLabel = "signing_infos"
)

func queryByLcdSlashing(i core.QueryClient) (string, error) {
	url := util.MakeQueryLcdUrl(slashingv1beta1.Query_ServiceDesc.Metadata.(string))
	switch {
	// Slashing parameters
	case i.Ixplac.GetMsgType() == SlashingQuerySlashingParamsMsgType:
		url = url + slashingParamsLabel

	// Slashing signing information
	case i.Ixplac.GetMsgType() == SlashingQuerySigningInfosMsgType:
		url = url + slashingSigningInfosLabel

	// Slashing signing information
	case i.Ixplac.GetMsgType() == SlashingQuerySigningInfoMsgType:
		convertMsg := i.Ixplac.GetMsg().(slashingtypes.QuerySigningInfoRequest)

		url = url + util.MakeQueryLabels(slashingSigningInfosLabel, convertMsg.ConsAddress)

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	i.Ixplac.GetHttpMutex().Lock()
	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		i.Ixplac.GetHttpMutex().Unlock()
		return "", err
	}
	i.Ixplac.GetHttpMutex().Unlock()

	return string(out), nil
}
