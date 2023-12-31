package authz

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	authzv1beta1 "cosmossdk.io/api/cosmos/authz/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

var out []byte
var res proto.Message
var err error

// Query client for authz module.
func QueryAuthz(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcAuthz(i)
	} else {
		return queryByLcdAuthz(i)
	}
}

func queryByGrpcAuthz(i core.QueryClient) (string, error) {
	queryClient := authz.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Authz grant
	case i.Ixplac.GetMsgType() == AuthzQueryGrantMsgType:
		convertMsg := i.Ixplac.GetMsg().(authz.QueryGrantsRequest)
		res, err = queryClient.Grants(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Authz grant by grantee
	case i.Ixplac.GetMsgType() == AuthzQueryGrantsByGranteeMsgType:
		convertMsg := i.Ixplac.GetMsg().(authz.QueryGranteeGrantsRequest)
		res, err = queryClient.GranteeGrants(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", i.Ixplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}

	// Authz grant by granter
	case i.Ixplac.GetMsgType() == AuthzQueryGrantsByGranterMsgType:
		convertMsg := i.Ixplac.GetMsg().(authz.QueryGranterGrantsRequest)
		res, err = queryClient.GranterGrants(
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
	authzGrantsLabel = "grants"
)

func queryByLcdAuthz(i core.QueryClient) (string, error) {

	url := util.MakeQueryLcdUrl(authzv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Authz grant
	case i.Ixplac.GetMsgType() == AuthzQueryGrantMsgType:
		convertMsg := i.Ixplac.GetMsg().(authz.QueryGrantsRequest)
		parsedGranter := convertMsg.Granter
		parsedGrantee := convertMsg.Grantee

		granter := "?granter=" + parsedGranter
		grantee := "&grantee=" + parsedGrantee

		url = url + authzGrantsLabel + granter + grantee

	// Authz grant by grantee
	case i.Ixplac.GetMsgType() == AuthzQueryGrantsByGranteeMsgType:
		convertMsg := i.Ixplac.GetMsg().(authz.QueryGranteeGrantsRequest)
		grantee := convertMsg.Grantee

		url = url + util.MakeQueryLabels(authzGrantsLabel, "grantee", grantee)

	// Authz grant by granter
	case i.Ixplac.GetMsgType() == AuthzQueryGrantsByGranterMsgType:
		convertMsg := i.Ixplac.GetMsg().(authz.QueryGranterGrantsRequest)
		granter := convertMsg.Granter

		url = url + util.MakeQueryLabels(authzGrantsLabel, "granter", granter)

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
