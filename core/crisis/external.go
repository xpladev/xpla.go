package crisis

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &CrisisExternal{}

type CrisisExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e CrisisExternal) {
	e.Xplac = xplac
	e.Name = CrisisModule
	return e
}

func (e CrisisExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e CrisisExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Submit proof that an invariant broken to halt the chain.
func (e CrisisExternal) InvariantBroken(invariantBrokenMsg types.InvariantBrokenMsg) provider.XplaClient {
	msg, err := MakeInvariantRouteMsg(invariantBrokenMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(CrisisInvariantBrokenMsgType, err)
	}

	return e.ToExternal(CrisisInvariantBrokenMsgType, msg)
}
