package authz

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &AuthzExternal{}

type AuthzExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e AuthzExternal) {
	e.Xplac = xplac
	e.Name = AuthzModule
	return e
}

func (e AuthzExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e AuthzExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Grant authorization to an address.
func (e AuthzExternal) AuthzGrant(authzGrantMsg types.AuthzGrantMsg) provider.XplaClient {
	msg, err := MakeAuthzGrantMsg(authzGrantMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(AuthzGrantMsgType, err)
	}

	return e.ToExternal(AuthzGrantMsgType, msg)
}

// Revoke authorization.
func (e AuthzExternal) AuthzRevoke(authzRevokeMsg types.AuthzRevokeMsg) provider.XplaClient {
	msg, err := MakeAuthzRevokeMsg(authzRevokeMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(AuthzRevokeMsgType, err)
	}

	return e.ToExternal(AuthzRevokeMsgType, msg)
}

// Execute transaction on behalf of granter account.
func (e AuthzExternal) AuthzExec(authzExecMsg types.AuthzExecMsg) provider.XplaClient {
	msg, err := MakeAuthzExecMsg(authzExecMsg, e.Xplac.GetEncoding())
	if err != nil {
		return e.Err(AuthzExecMsgType, err)
	}

	return e.ToExternal(AuthzExecMsgType, msg)
}

// Query

// Query grants for granter-grantee pair and optionally a msg-type-url.
// Also, it is able to support querying grants granted by granter and granted to a grantee.
func (e AuthzExternal) QueryAuthzGrants(queryAuthzGrantMsg types.QueryAuthzGrantMsg) provider.XplaClient {
	switch {
	case queryAuthzGrantMsg.Grantee != "" && queryAuthzGrantMsg.Granter != "":
		msg, err := MakeQueryAuthzGrantsMsg(queryAuthzGrantMsg)
		if err != nil {
			return e.Err(AuthzQueryGrantMsgType, err)
		}

		return e.ToExternal(AuthzQueryGrantMsgType, msg)

	case queryAuthzGrantMsg.Grantee != "" && queryAuthzGrantMsg.Granter == "":
		msg, err := MakeQueryAuthzGrantsByGranteeMsg(queryAuthzGrantMsg)
		if err != nil {
			return e.Err(AuthzQueryGrantsByGranteeMsgType, err)
		}

		return e.ToExternal(AuthzQueryGrantsByGranteeMsgType, msg)

	case queryAuthzGrantMsg.Grantee == "" && queryAuthzGrantMsg.Granter != "":
		msg, err := MakeQueryAuthzGrantsByGranterMsg(queryAuthzGrantMsg)
		if err != nil {
			return e.Err(AuthzQueryGrantsByGranterMsgType, err)
		}

		return e.ToExternal(AuthzQueryGrantsByGranterMsgType, msg)

	default:
		return e.Err(AuthzQueryGrantMsgType, types.ErrWrap(types.ErrInsufficientParams, "No query grants parameters"))
	}
}
