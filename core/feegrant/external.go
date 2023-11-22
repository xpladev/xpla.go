package feegrant

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type FeegrantExternal struct {
	Xplac provider.XplaClient
}

func NewFeegrantExternal(xplac provider.XplaClient) (e FeegrantExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Grant fee allowance to an address.
func (e FeegrantExternal) FeeGrant(grantMsg types.FeeGrantMsg) provider.XplaClient {
	msg, err := MakeFeeGrantMsg(grantMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(FeegrantModule).
		WithMsgType(FeegrantGrantMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Revoke fee-grant.
func (e FeegrantExternal) RevokeFeeGrant(revokeGrantMsg types.RevokeFeeGrantMsg) provider.XplaClient {
	msg, err := MakeRevokeFeeGrantMsg(revokeGrantMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(FeegrantModule).
		WithMsgType(FeegrantRevokeGrantMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Query details of fee grants.
func (e FeegrantExternal) QueryFeeGrants(queryFeeGrantMsg types.QueryFeeGrantMsg) provider.XplaClient {
	if queryFeeGrantMsg.Grantee != "" && queryFeeGrantMsg.Granter != "" {
		msg, err := MakeQueryFeeGrantMsg(queryFeeGrantMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(FeegrantModule).
			WithMsgType(FeegrantQueryGrantMsgType).
			WithMsg(msg)
	} else if queryFeeGrantMsg.Grantee != "" && queryFeeGrantMsg.Granter == "" {
		msg, err := MakeQueryFeeGrantsByGranteeMsg(queryFeeGrantMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(FeegrantModule).
			WithMsgType(FeegrantQueryGrantsByGranteeMsgType).
			WithMsg(msg)
	} else if queryFeeGrantMsg.Grantee == "" && queryFeeGrantMsg.Granter != "" {
		msg, err := MakeQueryFeeGrantsByGranterMsg(queryFeeGrantMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(FeegrantModule).
			WithMsgType(FeegrantQueryGrantsByGranterMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInsufficientParams, "no query grants parameters"))
	}

	return e.Xplac
}
