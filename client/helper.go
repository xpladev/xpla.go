package client

import (
	"github.com/xpladev/xpla.go/util/testutil"
)

func NewTestXplaClient() *XplaClient {
	return NewXplaClient(testutil.TestChainId)
}

func ResetXplac(xplac *XplaClient) *XplaClient {
	return ResetModuleAndMsgXplac(xplac).
		WithOptions(Options{}).
		WithErr(nil)
}

func ResetModuleAndMsgXplac(xplac *XplaClient) *XplaClient {
	return xplac.
		WithModule("").
		WithMsgType("").
		WithMsg(nil)
}
