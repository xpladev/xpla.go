package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
)

type coreModule struct{}

func NewCoreModule() core.CoreModule {
	return &coreModule{}
}

func (c *coreModule) Name() string {
	return WasmModule
}

func (c *coreModule) NewTxRouter(logger types.Logger, builder cmclient.TxBuilder, msgType string, msg interface{}) (cmclient.TxBuilder, error) {
	switch {
	case msgType == WasmStoreMsgType:
		convertMsg := msg.(wasm.MsgStoreCode)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == WasmInstantiateMsgType:
		convertMsg := msg.(wasm.MsgInstantiateContract)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == WasmExecuteMsgType:
		convertMsg := msg.(wasm.MsgExecuteContract)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == WasmClearContractAdminMsgType:
		convertMsg := msg.(wasm.MsgClearAdmin)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == WasmSetContractAdminMsgType:
		convertMsg := msg.(wasm.MsgUpdateAdmin)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	case msgType == WasmMigrateMsgType:
		convertMsg := msg.(wasm.MsgMigrateContract)
		err := builder.SetMsgs(&convertMsg)
		if err != nil {
			return nil, logger.Err(err)
		}

	default:
		return nil, logger.Err(types.ErrWrap(types.ErrInvalidMsgType, msgType))
	}

	return builder, nil
}

func (c *coreModule) NewQueryRouter(q core.QueryClient) (string, error) {
	return QueryWasm(q)
}
