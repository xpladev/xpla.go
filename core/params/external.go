package params

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &ParamsExternal{}

type ParamsExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e ParamsExternal) {
	e.Xplac = xplac
	e.Name = ParamsModule
	return e
}

func (e ParamsExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e ParamsExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Submit a parameter change proposal.
func (e ParamsExternal) ParamChange(paramChangeMsg types.ParamChangeMsg) provider.XplaClient {
	msg, err := MakeProposalParamChangeMsg(paramChangeMsg, e.Xplac.GetFromAddress(), e.Xplac.GetEncoding())
	if err != nil {
		return e.Err(ParamsProposalParamChangeMsgType, err)
	}

	return e.ToExternal(ParamsProposalParamChangeMsgType, msg)
}

// Query

// Query for raw parameters by subspace and key.
func (e ParamsExternal) QuerySubspace(subspaceMsg types.SubspaceMsg) provider.XplaClient {
	msg, err := MakeQueryParamsSubspaceMsg(subspaceMsg)
	if err != nil {
		return e.Err(ParamsQuerySubpsaceMsgType, err)
	}

	return e.ToExternal(ParamsQuerySubpsaceMsgType, msg)
}
