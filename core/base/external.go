package base

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &BaseExternal{}

type BaseExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e BaseExternal) {
	e.Xplac = xplac
	e.Name = Base
	return e
}

func (e BaseExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e BaseExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Query

// Query node info
func (e BaseExternal) NodeInfo() provider.XplaClient {
	msg, err := MakeBaseNodeInfoMsg()
	if err != nil {
		return e.Err(BaseNodeInfoMsgType, err)
	}

	return e.ToExternal(BaseNodeInfoMsgType, msg)
}

// Query syncing
func (e BaseExternal) Syncing() provider.XplaClient {
	msg, err := MakeBaseSyncingMsg()
	if err != nil {
		return e.Err(BaseSyncingMsgType, err)
	}

	return e.ToExternal(BaseSyncingMsgType, msg)
}

// Query block
func (e BaseExternal) Block(blockMsg ...types.BlockMsg) provider.XplaClient {
	switch {
	case len(blockMsg) == 0:
		msg, err := MakeBaseLatestBlockMsg()
		if err != nil {
			return e.Err(BaseLatestBlockMsgtype, err)
		}

		return e.ToExternal(BaseLatestBlockMsgtype, msg)

	case len(blockMsg) == 1:
		msg, err := MakeBaseBlockByHeightMsg(blockMsg[0])
		if err != nil {
			return e.Err(BaseBlockByHeightMsgType, err)
		}

		return e.ToExternal(BaseBlockByHeightMsgType, msg)

	default:
		return e.Err(BaseBlockByHeightMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}

// Query validator set
func (e BaseExternal) ValidatorSet(validatorSetMsg ...types.ValidatorSetMsg) provider.XplaClient {
	switch {
	case len(validatorSetMsg) == 0:
		msg, err := MakeLatestValidatorSetMsg()
		if err != nil {
			return e.Err(BaseLatestValidatorSetMsgType, err)
		}

		return e.ToExternal(BaseLatestValidatorSetMsgType, msg)

	case len(validatorSetMsg) == 1:
		msg, err := MakeValidatorSetByHeightMsg(validatorSetMsg[0])
		if err != nil {
			return e.Err(BaseValidatorSetByHeightMsgType, err)
		}

		return e.ToExternal(BaseValidatorSetByHeightMsgType, msg)

	default:
		return e.Err(BaseValidatorSetByHeightMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}
