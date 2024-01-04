package evidence

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &EvidenceExternal{}

type EvidenceExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e EvidenceExternal) {
	e.Xplac = xplac
	e.Name = EvidenceModule
	return e
}

func (e EvidenceExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e EvidenceExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Query

// Query for evidence by hash or for all (paginated) submitted evidence.
func (e EvidenceExternal) QueryEvidence(queryEvidenceMsg ...types.QueryEvidenceMsg) provider.XplaClient {
	switch {

	case len(queryEvidenceMsg) == 0:
		msg, err := MakeQueryAllEvidenceMsg()
		if err != nil {
			return e.Err(EvidenceQueryAllMsgType, err)
		}

		return e.ToExternal(EvidenceQueryAllMsgType, msg)

	case len(queryEvidenceMsg) == 1:
		msg, err := MakeQueryEvidenceMsg(queryEvidenceMsg[0])
		if err != nil {
			return e.Err(EvidenceQueryMsgType, err)
		}

		return e.ToExternal(EvidenceQueryMsgType, msg)

	default:
		return e.Err(EvidenceQueryMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}
