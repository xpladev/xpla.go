package mint

import (
	"github.com/xpladev/xpla.go/provider"
)

type MintExternal struct {
	Xplac provider.XplaClient
}

func NewMintExternal(xplac provider.XplaClient) (e MintExternal) {
	e.Xplac = xplac
	return e
}

// Query

// Query the current minting parameters.
func (e MintExternal) MintParams() provider.XplaClient {
	msg, err := MakeQueryMintParamsMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(MintModule).
		WithMsgType(MintQueryMintParamsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query the current minting inflation value.
func (e MintExternal) Inflation() provider.XplaClient {
	msg, err := MakeQueryInflationMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(MintModule).
		WithMsgType(MintQueryInflationMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query the current minting annual provisions value.
func (e MintExternal) AnnualProvisions() provider.XplaClient {
	msg, err := MakeQueryAnnualProvisionsMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(MintModule).
		WithMsgType(MintQueryAnnualProvisionsMsgType).
		WithMsg(msg)
	return e.Xplac
}
