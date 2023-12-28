package volunteer

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

type VolunteerExternal struct {
	Xplac provider.XplaClient
}

func NewVolunteerExternal(xplac provider.XplaClient) (e VolunteerExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Register new volunteer validator.
func (e VolunteerExternal) RegisterVolunteerValidator(registerVolunteerValidatorMsg types.RegisterVolunteerValidatorMsg) provider.XplaClient {
	msg, err := MakeRegisterVolunteerValidatorMsg(registerVolunteerValidatorMsg, e.Xplac.GetEncoding(), e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(VolunteerModule).
		WithMsgType(VolunteerRegisterVolunteerValidatorMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Unregister a volunteer validator.
func (e VolunteerExternal) UnregisterVolunteerValidator(unregisterVolunteerValidatorMsg types.UnregisterVolunteerValidatorMsg) provider.XplaClient {
	msg, err := MakeUnregisterVolunteerValidatorMsg(unregisterVolunteerValidatorMsg, e.Xplac.GetEncoding(), e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(VolunteerModule).
		WithMsgType(VolunteerUnregisterVolunteerValidatorMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Query volunteer validators.
func (e VolunteerExternal) QueryVolunteerValidators() provider.XplaClient {
	msg, err := MakeQueryVolunteerValidatorsMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(VolunteerModule).
		WithMsgType(VolunteerQueryValidatorsMsgType).
		WithMsg(msg)
	return e.Xplac
}
