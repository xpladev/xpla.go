package base_test

import (
	"github.com/xpladev/xpla.go/core/base"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	c := base.NewCoreModule()

	// test get name
	s.Require().Equal(base.Base, c.Name())

	// test tx
	_, err := c.NewTxRouter(s.xplac.GetLogger(), nil, "", nil)
	s.Require().Error(err)
}
