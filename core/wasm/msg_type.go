package wasm

const (
	WasmModule                    = "wasm"
	WasmStoreMsgType              = "store-code"
	WasmInstantiateMsgType        = "instantiate-contract"
	WasmExecuteMsgType            = "execute-contract"
	WasmClearContractAdminMsgType = "clear-contract-admin"
	WasmSetContractAdminMsgType   = "set-contract-admin"
	WasmMigrateMsgType            = "migrate"
	WasmQueryContractMsgType      = "query-contract"
	WasmListCodeMsgType           = "list-code"
	WasmListContractByCodeMsgType = "list-contract-by-code"
	WasmDownloadMsgType           = "download"
	WasmCodeInfoMsgType           = "code-info"
	WasmContractInfoMsgType       = "contract-info"
	WasmContractStateAllMsgType   = "contract-state-all"
	WasmContractHistoryMsgType    = "contract-history"
	WasmPinnedMsgType             = "pinned"
	WasmLibwasmvmVersionMsgType   = "libwasmvm-version"
)
