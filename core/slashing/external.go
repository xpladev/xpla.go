package slashing

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &SlashingExternal{}

type SlashingExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e SlashingExternal) {
	e.Xplac = xplac
	e.Name = SlashingModule
	return e
}

func (e SlashingExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e SlashingExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Unjail validator previously jailed for downtime.
func (e SlashingExternal) Unjail() provider.XplaClient {
	msg, err := MakeUnjailMsg(e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(SlashingUnjailMsgType, err)
	}

	return e.ToExternal(SlashingUnjailMsgType, msg)
}

// Query

// Query the current slashing parameters.
func (e SlashingExternal) SlashingParams() provider.XplaClient {
	msg, err := MakeQuerySlashingParamsMsg()
	if err != nil {
		return e.Err(SlashingQuerySlashingParamsMsgType, err)
	}

	return e.ToExternal(SlashingQuerySlashingParamsMsgType, msg)
}

// Query a validator's signing information or signing information of all validators.
func (e SlashingExternal) SigningInfos(signingInfoMsg ...types.SigningInfoMsg) provider.XplaClient {
	switch {
	case len(signingInfoMsg) == 0:
		msg, err := MakeQuerySigningInfosMsg()
		if err != nil {
			return e.Err(SlashingQuerySigningInfosMsgType, err)
		}

		return e.ToExternal(SlashingQuerySigningInfosMsgType, msg)

	case len(signingInfoMsg) == 1:
		msg, err := MakeQuerySigningInfoMsg(signingInfoMsg[0], e.Xplac.GetEncoding())
		if err != nil {
			return e.Err(SlashingQuerySigningInfoMsgType, err)
		}

		return e.ToExternal(SlashingQuerySigningInfoMsgType, msg)

	default:
		return e.Err(SlashingQuerySigningInfoMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}
