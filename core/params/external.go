package params

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

type ParamsExternal struct {
	Xplac provider.XplaClient
}

func NewParamsExternal(xplac provider.XplaClient) (e ParamsExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Submit a parameter change proposal.
func (e ParamsExternal) ParamChange(paramChangeMsg types.ParamChangeMsg) provider.XplaClient {
	msg, err := MakeProposalParamChangeMsg(paramChangeMsg, e.Xplac.GetFromAddress(), e.Xplac.GetEncoding())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(ParamsModule).
		WithMsgType(ParamsProposalParamChangeMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Query for raw parameters by subspace and key.
func (e ParamsExternal) QuerySubspace(subspaceMsg types.SubspaceMsg) provider.XplaClient {
	msg, err := MakeQueryParamsSubspaceMsg(subspaceMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(ParamsModule).
		WithMsgType(ParamsQuerySubpsaceMsgType).
		WithMsg(msg)

	return e.Xplac
}
