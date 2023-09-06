package controller

import (
	"sync"

	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/core/auth"
)

var once sync.Once
var cc *coreController

// Controller is able to contol modules in the core package.
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
	once.Do(
		func() {
			cc = NewCoreController(
				auth.NewCoreModule(),
			)
		})
	return cc
}

func NewCoreController(coreModules ...core.CoreModule) *coreController {
	m := make(map[string]core.CoreModule)
	for _, coreModule := range coreModules {
		m[coreModule.Name()] = coreModule
	}

	return &coreController{
		cores: m,
	}
}

func (c coreController) Get(moduleName string) core.CoreModule {
	return c.cores[moduleName]
}
