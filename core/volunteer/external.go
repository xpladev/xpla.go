package volunteer

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &VolunteerExternal{}

type VolunteerExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e VolunteerExternal) {
	e.Xplac = xplac
	e.Name = VolunteerModule
	return e
}

func (e VolunteerExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e VolunteerExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Register new volunteer validator.
func (e VolunteerExternal) RegisterVolunteerValidator(registerVolunteerValidatorMsg types.RegisterVolunteerValidatorMsg) provider.XplaClient {
	msg, err := MakeRegisterVolunteerValidatorMsg(registerVolunteerValidatorMsg, e.Xplac.GetEncoding(), e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(VolunteerRegisterVolunteerValidatorMsgType, err)
	}

	return e.ToExternal(VolunteerRegisterVolunteerValidatorMsgType, msg)
}

// Unregister a volunteer validator.
func (e VolunteerExternal) UnregisterVolunteerValidator(unregisterVolunteerValidatorMsg types.UnregisterVolunteerValidatorMsg) provider.XplaClient {
	msg, err := MakeUnregisterVolunteerValidatorMsg(unregisterVolunteerValidatorMsg, e.Xplac.GetEncoding(), e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(VolunteerUnregisterVolunteerValidatorMsgType, err)
	}

	return e.ToExternal(VolunteerUnregisterVolunteerValidatorMsgType, msg)
}

// Query

// Query volunteer validators.
func (e VolunteerExternal) QueryVolunteerValidators() provider.XplaClient {
	msg, err := MakeQueryVolunteerValidatorsMsg()
	if err != nil {
		return e.Err(VolunteerQueryValidatorsMsgType, err)
	}

	return e.ToExternal(VolunteerQueryValidatorsMsgType, msg)
}
