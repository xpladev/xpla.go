package base

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type BaseExternal struct {
	Xplac provider.XplaClient
}

func NewBaseExternal(xplac provider.XplaClient) (e BaseExternal) {
	e.Xplac = xplac
	return e
}

// Query

// Query node info
func (e BaseExternal) NodeInfo() provider.XplaClient {
	msg, err := MakeBaseNodeInfoMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(Base).
		WithMsgType(BaseNodeInfoMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query syncing
func (e BaseExternal) Syncing() provider.XplaClient {
	msg, err := MakeBaseSyncingMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(Base).
		WithMsgType(BaseSyncingMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query block
func (e BaseExternal) Block(blockMsg ...types.BlockMsg) provider.XplaClient {
	if len(blockMsg) == 0 {
		msg, err := MakeBaseLatestBlockMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(Base).
			WithMsgType(BaseLatestBlockMsgtype).
			WithMsg(msg)
	} else if len(blockMsg) == 1 {
		msg, err := MakeBaseBlockByHeightMsg(blockMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(Base).
			WithMsgType(BaseBlockByHeightMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}

// Query validator set
func (e BaseExternal) ValidatorSet(validatorSetMsg ...types.ValidatorSetMsg) provider.XplaClient {
	if len(validatorSetMsg) == 0 {
		msg, err := MakeLatestValidatorSetMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(Base).
			WithMsgType(BaseLatestValidatorSetMsgType).
			WithMsg(msg)
	} else if len(validatorSetMsg) == 1 {
		msg, err := MakeValidatorSetByHeightMsg(validatorSetMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(Base).
			WithMsgType(BaseValidatorSetByHeightMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}
