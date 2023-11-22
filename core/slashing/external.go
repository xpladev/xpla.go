package slashing

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type SlashingExternal struct {
	Xplac provider.XplaClient
}

func NewSlashingExternal(xplac provider.XplaClient) (e SlashingExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Unjail validator previously jailed for downtime.
func (e SlashingExternal) Unjail() provider.XplaClient {
	msg, err := MakeUnjailMsg(e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(SlashingModule).
		WithMsgType(SlahsingUnjailMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Query the current slashing parameters.
func (e SlashingExternal) SlashingParams() provider.XplaClient {
	msg, err := MakeQuerySlashingParamsMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(SlashingModule).
		WithMsgType(SlashingQuerySlashingParamsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query a validator's signing information or signing information of all validators.
func (e SlashingExternal) SigningInfos(signingInfoMsg ...types.SigningInfoMsg) provider.XplaClient {
	if len(signingInfoMsg) == 0 {
		msg, err := MakeQuerySigningInfosMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(SlashingModule).
			WithMsgType(SlashingQuerySigningInfosMsgType).
			WithMsg(msg)
	} else if len(signingInfoMsg) == 1 {
		msg, err := MakeQuerySigningInfoMsg(signingInfoMsg[0], e.Xplac.GetEncoding())
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(SlashingModule).
			WithMsgType(SlashingQuerySigningInfoMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}
