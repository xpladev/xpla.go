package auth

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	authv1beta1 "cosmossdk.io/api/cosmos/auth/v1beta1"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

var out []byte
var res proto.Message
var err error

// Query client for auth module.
func QueryAuth(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcAuth(i)
	} else {
		return queryByLcdAuth(i)
	}
}

func queryByGrpcAuth(i core.QueryClient) (string, error) {
	queryClient := authtypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Auth params
	case i.Ixplac.GetMsgType() == AuthQueryParamsMsgType:
		convertMsg := i.Ixplac.GetMsg().(authtypes.QueryParamsRequest)
		res, err = queryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Auth account
	case i.Ixplac.GetMsgType() == AuthQueryAccAddressMsgType:
		convertMsg := i.Ixplac.GetMsg().(authtypes.QueryAccountRequest)
		res, err = queryClient.Account(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Auth accounts
	case i.Ixplac.GetMsgType() == AuthQueryAccountsMsgType:
		convertMsg := i.Ixplac.GetMsg().(authtypes.QueryAccountsRequest)
		res, err = queryClient.Accounts(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Auth tx by event
	case i.Ixplac.GetMsgType() == AuthQueryTxsByEventsMsgType:
		if i.Ixplac.GetRpc() == "" {
			return "", util.LogErr(errors.ErrNotSatisfiedOptions, "query txs by events, need RPC URL when txs methods")
		}
		convertMsg := i.Ixplac.GetMsg().(QueryTxsByEventParseMsg)
		clientCtx, err := core.ClientForQuery(i)
		if err != nil {
			return "", err
		}

		res, err = authtx.QueryTxsByEvents(clientCtx, convertMsg.TmEvents, convertMsg.Page, convertMsg.Limit, "")
		if err != nil {
			return "", util.LogErr(errors.ErrRpcRequest, err)
		}

	// Auth tx
	case i.Ixplac.GetMsgType() == AuthQueryTxMsgType:
		if i.Ixplac.GetRpc() == "" {
			return "", util.LogErr(errors.ErrNotSatisfiedOptions, "auth query tx msg, need RPC URL when txs methods")
		}
		convertMsg := i.Ixplac.GetMsg().(QueryTxParseMsg)

		clientCtx, err := core.ClientForQuery(i)
		if err != nil {
			return "", err
		}

		if convertMsg.TxType == "hash" {
			res, err = authtx.QueryTx(clientCtx, convertMsg.TmEvents[0])
			if err != nil {
				return "", util.LogErr(errors.ErrRpcRequest, err)
			}
		} else {
			res, err = authtx.QueryTxsByEvents(clientCtx, convertMsg.TmEvents, rest.DefaultPage, rest.DefaultLimit, "")
			if err != nil {
				return "", util.LogErr(errors.ErrRpcRequest, err)
			}
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
	authParamsLabel   = "params"
	authAccountsLabel = "accounts"
	authTxsLabel      = "txs"
)

func queryByLcdAuth(i core.QueryClient) (string, error) {

	url := util.MakeQueryLcdUrl(authv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Auth params
	case i.Ixplac.GetMsgType() == AuthQueryParamsMsgType:
		url = url + authParamsLabel

	// Auth account
	case i.Ixplac.GetMsgType() == AuthQueryAccAddressMsgType:
		convertMsg := i.Ixplac.GetMsg().(authtypes.QueryAccountRequest)
		url = url + util.MakeQueryLabels(authAccountsLabel, convertMsg.Address)

	// Auth accounts
	case i.Ixplac.GetMsgType() == AuthQueryAccountsMsgType:
		url = url + authAccountsLabel

	// Auth tx by event
	case i.Ixplac.GetMsgType() == AuthQueryTxsByEventsMsgType:
		convertMsg := i.Ixplac.GetMsg().(QueryTxsByEventParseMsg)

		if len(convertMsg.TmEvents) > 1 {
			return "", util.LogErr(errors.ErrNotSupport, "support only one event on the LCD")
		}

		parsedEvent := convertMsg.TmEvents[0]
		parsedPage := convertMsg.Page
		parsedLimit := convertMsg.Limit

		events := "?events=" + parsedEvent
		page := "&pagination.page=" + util.FromIntToString(parsedPage)
		limit := "&pagination.limit=" + util.FromIntToString(parsedLimit)

		url = "/cosmos/tx/v1beta1/"
		url = url + authTxsLabel + events + page + limit

	// Auth tx
	case i.Ixplac.GetMsgType() == AuthQueryTxMsgType:
		convertMsg := i.Ixplac.GetMsg().(QueryTxParseMsg)

		if len(convertMsg.TmEvents) > 1 {
			return "", util.LogErr(errors.ErrNotSupport, "support only one event on the LCD")
		}

		parsedValue := convertMsg.TmEvents
		parsedTxType := convertMsg.TxType

		url = "/cosmos/tx/v1beta1/"
		if parsedTxType == "hash" {
			url = url + util.MakeQueryLabels(authTxsLabel, parsedValue[0])

		} else if parsedTxType == "signature" {
			// inactivate
			return "", util.LogErr(errors.ErrNotSupport, "inactivate GetTxEvent('signature') when using LCD because of sometimes generating parsing error that based64 encoded signature has '='")
			// events := "?events=" + parsedValue
			// page := "&pagination.page=" + util.FromIntToString(rest.DefaultPage)
			// limit := "&pagination.limit=" + util.FromIntToString(rest.DefaultLimit)

			// url = url + authTxsLabel + events + page + limit
		} else {
			events := "?events=" + parsedValue[0]
			page := "&pagination.page=" + util.FromIntToString(rest.DefaultPage)
			limit := "&pagination.limit=" + util.FromIntToString(rest.DefaultLimit)

			url = url + authTxsLabel + events + page + limit
		}

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
