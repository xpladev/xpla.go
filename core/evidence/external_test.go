package evidence_test

import (
	mevidence "github.com/xpladev/xpla.go/core/evidence"
	"github.com/xpladev/xpla.go/types"
)

func (s *IntegrationTestSuite) TestEvidence() {
	testTxHash := "B6BBBB649F19E8970EF274C0083FE945FD38AD8C524D68BB3FE3A20D72DF03C4"

	// all evidence
	s.xplac.QueryEvidence()

	makeQueryAllEvidenceMsg, err := mevidence.MakeQueryAllEvidenceMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryAllEvidenceMsg, s.xplac.GetMsg())
	s.Require().Equal(mevidence.EvidenceModule, s.xplac.GetModule())
	s.Require().Equal(mevidence.EvidenceQueryAllMsgType, s.xplac.GetMsgType())

	// evidence
	queryEvidenceMsg := types.QueryEvidenceMsg{
		Hash: testTxHash,
	}
	s.xplac.QueryEvidence(queryEvidenceMsg)

	makeQueryEvidenceMsg, err := mevidence.MakeQueryEvidenceMsg(queryEvidenceMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryEvidenceMsg, s.xplac.GetMsg())
	s.Require().Equal(mevidence.EvidenceModule, s.xplac.GetModule())
	s.Require().Equal(mevidence.EvidenceQueryMsgType, s.xplac.GetMsgType())
}
