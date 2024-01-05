package feegrant

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &FeegrantExternal{}

type FeegrantExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e FeegrantExternal) {
	e.Xplac = xplac
	e.Name = FeegrantModule
	return e
}

func (e FeegrantExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e FeegrantExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Grant fee allowance to an address.
func (e FeegrantExternal) FeeGrant(grantMsg types.FeeGrantMsg) provider.XplaClient {
	msg, err := MakeFeeGrantMsg(grantMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(FeegrantGrantMsgType, err)
	}

	return e.ToExternal(FeegrantGrantMsgType, msg)
}

// Revoke fee-grant.
func (e FeegrantExternal) RevokeFeeGrant(revokeGrantMsg types.RevokeFeeGrantMsg) provider.XplaClient {
	msg, err := MakeRevokeFeeGrantMsg(revokeGrantMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(FeegrantRevokeGrantMsgType, err)
	}

	return e.ToExternal(FeegrantRevokeGrantMsgType, msg)
}

// Query

// Query details of fee grants.
func (e FeegrantExternal) QueryFeeGrants(queryFeeGrantMsg types.QueryFeeGrantMsg) provider.XplaClient {
	switch {
	case queryFeeGrantMsg.Grantee != "" && queryFeeGrantMsg.Granter != "":
		msg, err := MakeQueryFeeGrantMsg(queryFeeGrantMsg)
		if err != nil {
			return e.Err(FeegrantQueryGrantMsgType, err)
		}

		return e.ToExternal(FeegrantQueryGrantMsgType, msg)

	case queryFeeGrantMsg.Grantee != "" && queryFeeGrantMsg.Granter == "":
		msg, err := MakeQueryFeeGrantsByGranteeMsg(queryFeeGrantMsg)
		if err != nil {
			return e.Err(FeegrantQueryGrantsByGranteeMsgType, err)
		}

		return e.ToExternal(FeegrantQueryGrantsByGranteeMsgType, msg)

	case queryFeeGrantMsg.Grantee == "" && queryFeeGrantMsg.Granter != "":
		msg, err := MakeQueryFeeGrantsByGranterMsg(queryFeeGrantMsg)
		if err != nil {
			return e.Err(FeegrantQueryGrantsByGranterMsgType, err)
		}

		return e.ToExternal(FeegrantQueryGrantsByGranterMsgType, msg)

	default:
		return e.Err(FeegrantQueryGrantsByGranterMsgType, types.ErrWrap(types.ErrInsufficientParams, "no query grants parameters"))
	}
}
