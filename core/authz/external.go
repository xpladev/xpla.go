package authz

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type AuthzExternal struct {
	Xplac provider.XplaClient
}

func NewAuthzExternal(xplac provider.XplaClient) (e AuthzExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Grant authorization to an address.
func (e AuthzExternal) AuthzGrant(authzGrantMsg types.AuthzGrantMsg) provider.XplaClient {
	msg, err := MakeAuthzGrantMsg(authzGrantMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(AuthzModule).
		WithMsgType(AuthzGrantMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Revoke authorization.
func (e AuthzExternal) AuthzRevoke(authzRevokeMsg types.AuthzRevokeMsg) provider.XplaClient {
	msg, err := MakeAuthzRevokeMsg(authzRevokeMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(AuthzModule).
		WithMsgType(AuthzRevokeMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Execute transaction on behalf of granter account.
func (e AuthzExternal) AuthzExec(authzExecMsg types.AuthzExecMsg) provider.XplaClient {
	msg, err := MakeAuthzExecMsg(authzExecMsg, e.Xplac.GetEncoding())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(AuthzModule).
		WithMsgType(AuthzExecMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Query grants for granter-grantee pair and optionally a msg-type-url.
// Also, it is able to support querying grants granted by granter and granted to a grantee.
func (e AuthzExternal) QueryAuthzGrants(queryAuthzGrantMsg types.QueryAuthzGrantMsg) provider.XplaClient {
	if queryAuthzGrantMsg.Grantee != "" && queryAuthzGrantMsg.Granter != "" {
		msg, err := MakeQueryAuthzGrantsMsg(queryAuthzGrantMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(AuthzModule).
			WithMsgType(AuthzQueryGrantMsgType).
			WithMsg(msg)
	} else if queryAuthzGrantMsg.Grantee != "" && queryAuthzGrantMsg.Granter == "" {
		msg, err := MakeQueryAuthzGrantsByGranteeMsg(queryAuthzGrantMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(AuthzModule).
			WithMsgType(AuthzQueryGrantsByGranteeMsgType).
			WithMsg(msg)
	} else if queryAuthzGrantMsg.Grantee == "" && queryAuthzGrantMsg.Granter != "" {
		msg, err := MakeQueryAuthzGrantsByGranterMsg(queryAuthzGrantMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(AuthzModule).
			WithMsgType(AuthzQueryGrantsByGranterMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInsufficientParams, "No query grants parameters"))
	}
	return e.Xplac
}
