package evidence_test

import (
	"github.com/xpladev/xpla.go/core/evidence"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	c := evidence.NewCoreModule()

	// test get name
	s.Require().Equal(evidence.EvidenceModule, c.Name())

	// test tx
	_, err := c.NewTxRouter(s.xplac.GetLogger(), nil, "", nil)
	s.Require().Error(err)
}
