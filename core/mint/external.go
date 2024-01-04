package mint

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &MintExternal{}

type MintExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e MintExternal) {
	e.Xplac = xplac
	e.Name = MintModule
	return e
}

func (e MintExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e MintExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Query

// Query the current minting parameters.
func (e MintExternal) MintParams() provider.XplaClient {
	msg, err := MakeQueryMintParamsMsg()
	if err != nil {
		return e.Err(MintQueryMintParamsMsgType, err)
	}

	return e.ToExternal(MintQueryMintParamsMsgType, msg)
}

// Query the current minting inflation value.
func (e MintExternal) Inflation() provider.XplaClient {
	msg, err := MakeQueryInflationMsg()
	if err != nil {
		return e.Err(MintQueryInflationMsgType, err)
	}

	return e.ToExternal(MintQueryInflationMsgType, msg)
}

// Query the current minting annual provisions value.
func (e MintExternal) AnnualProvisions() provider.XplaClient {
	msg, err := MakeQueryAnnualProvisionsMsg()
	if err != nil {
		return e.Err(MintQueryAnnualProvisionsMsgType, err)
	}

	return e.ToExternal(MintQueryAnnualProvisionsMsgType, msg)
}
