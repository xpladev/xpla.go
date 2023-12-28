package controller

import (
	"sync"

	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/core/auth"
	"github.com/xpladev/xpla.go/core/authz"
	"github.com/xpladev/xpla.go/core/bank"
	"github.com/xpladev/xpla.go/core/base"
	"github.com/xpladev/xpla.go/core/crisis"
	"github.com/xpladev/xpla.go/core/distribution"
	"github.com/xpladev/xpla.go/core/evidence"
	"github.com/xpladev/xpla.go/core/evm"
	"github.com/xpladev/xpla.go/core/feegrant"
	"github.com/xpladev/xpla.go/core/gov"
	"github.com/xpladev/xpla.go/core/ibc"
	"github.com/xpladev/xpla.go/core/mint"
	"github.com/xpladev/xpla.go/core/params"
	"github.com/xpladev/xpla.go/core/reward"
	"github.com/xpladev/xpla.go/core/slashing"
	"github.com/xpladev/xpla.go/core/staking"
	"github.com/xpladev/xpla.go/core/upgrade"
	"github.com/xpladev/xpla.go/core/volunteer"
	"github.com/xpladev/xpla.go/core/wasm"
)

var once sync.Once
var cc *coreController

// Controller is able to control modules in the core package.
// Route Tx & Query logic by message type.
// If need to add new modules of XPLA, insert NewCoreModule in the core controller.
type coreController struct {
	cores map[string]core.CoreModule
}

func init() {
	Controller()
}

// Set core controller only once as singleton, and get core controller.
func Controller() *coreController {
	once.Do(func() {
		cc = NewCoreController(
			auth.NewCoreModule(),
			authz.NewCoreModule(),
			bank.NewCoreModule(),
			base.NewCoreModule(),
			crisis.NewCoreModule(),
			distribution.NewCoreModule(),
			evidence.NewCoreModule(),
			evm.NewCoreModule(),
			feegrant.NewCoreModule(),
			gov.NewCoreModule(),
			ibc.NewCoreModule(),
			mint.NewCoreModule(),
			params.NewCoreModule(),
			reward.NewCoreModule(),
			slashing.NewCoreModule(),
			staking.NewCoreModule(),
			upgrade.NewCoreModule(),
			volunteer.NewCoreModule(),
			wasm.NewCoreModule(),
		)
	})
	return cc
}

// Register routing info of core modules in the hash map.
func NewCoreController(coreModules ...core.CoreModule) *coreController {
	m := make(map[string]core.CoreModule)
	for _, coreModule := range coreModules {
		m[coreModule.Name()] = coreModule
	}

	return &coreController{
		cores: m,
	}
}

// Get info of each module by its name.
func (c coreController) Get(moduleName string) core.CoreModule {
	return c.cores[moduleName]
}
