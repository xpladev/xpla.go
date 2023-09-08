package evidence

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type EvidenceExternal struct {
	Xplac provider.XplaClient
}

func NewEvidenceExternal(xplac provider.XplaClient) (e EvidenceExternal) {
	e.Xplac = xplac
	return e
}

// Query

// Query for evidence by hash or for all (paginated) submitted evidence.
func (e EvidenceExternal) QueryEvidence(queryEvidenceMsg ...types.QueryEvidenceMsg) provider.XplaClient {
	if len(queryEvidenceMsg) == 0 {
		msg, err := MakeQueryAllEvidenceMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(EvidenceModule).
			WithMsgType(EvidenceQueryAllMsgType).
			WithMsg(msg)
	} else if len(queryEvidenceMsg) == 1 {
		msg, err := MakeQueryEvidenceMsg(queryEvidenceMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(EvidenceModule).
			WithMsgType(EvidenceQueryMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}
