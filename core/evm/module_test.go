package evm_test

import (
	"github.com/xpladev/xpla.go/core/evm"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	c := evm.NewCoreModule()

	// test get name
	s.Require().Equal(evm.EvmModule, c.Name())

	// test tx
	_, err := c.NewTxRouter(nil, "", nil)
	s.Require().Error(err)
}
