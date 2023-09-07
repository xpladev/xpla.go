package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
)

type coreModule struct{}

func NewCoreModule() core.CoreModule {
	return &coreModule{}
}

func (c *coreModule) Name() string {
	return WasmModule
}

func (c *coreModule) NewTxRouter(builder cmclient.TxBuilder, msgType string, msg interface{}) (cmclient.TxBuilder, error) {
	switch {
	case msgType == WasmStoreMsgType:
		convertMsg := msg.(wasm.MsgStoreCode)
		builder.SetMsgs(&convertMsg)

	case msgType == WasmInstantiateMsgType:
		convertMsg := msg.(wasm.MsgInstantiateContract)
		builder.SetMsgs(&convertMsg)

	case msgType == WasmExecuteMsgType:
		convertMsg := msg.(wasm.MsgExecuteContract)
		builder.SetMsgs(&convertMsg)

	case msgType == WasmClearContractAdminMsgType:
		convertMsg := msg.(wasm.MsgClearAdmin)
		builder.SetMsgs(&convertMsg)

	case msgType == WasmSetContractAdminMsgType:
		convertMsg := msg.(wasm.MsgUpdateAdmin)
		builder.SetMsgs(&convertMsg)

	case msgType == WasmMigrateMsgType:
		convertMsg := msg.(wasm.MsgMigrateContract)
		builder.SetMsgs(&convertMsg)

	default:
		return nil, util.LogErr(errors.ErrInvalidMsgType, msgType)
	}

	return builder, nil
}

func (c *coreModule) NewQueryRouter(q core.QueryClient) (string, error) {
	return QueryWasm(q)
}
