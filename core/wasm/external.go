package wasm

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &WasmExternal{}

type WasmExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e WasmExternal) {
	e.Xplac = xplac
	e.Name = WasmModule
	return e
}

func (e WasmExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e WasmExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Upload a wasm binary.
func (e WasmExternal) StoreCode(storeMsg types.StoreMsg) provider.XplaClient {
	msg, err := MakeStoreCodeMsg(storeMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(WasmStoreMsgType, err)
	}

	return e.ToExternal(WasmStoreMsgType, msg)
}

// Instantiate a wasm contract.
func (e WasmExternal) InstantiateContract(instantiageMsg types.InstantiateMsg) provider.XplaClient {
	msg, err := MakeInstantiateMsg(instantiageMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(WasmInstantiateMsgType, err)
	}

	return e.ToExternal(WasmInstantiateMsgType, msg)
}

// Execute a wasm contract.
func (e WasmExternal) ExecuteContract(executeMsg types.ExecuteMsg) provider.XplaClient {
	msg, err := MakeExecuteMsg(executeMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(WasmExecuteMsgType, err)
	}

	return e.ToExternal(WasmExecuteMsgType, msg)
}

// Clears admin for a contract to prevent further migrations.
func (e WasmExternal) ClearContractAdmin(clearContractAdminMsg types.ClearContractAdminMsg) provider.XplaClient {
	msg, err := MakeClearContractAdminMsg(clearContractAdminMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(WasmClearContractAdminMsgType, err)
	}

	return e.ToExternal(WasmClearContractAdminMsgType, msg)
}

// Set new admin for a contract.
func (e WasmExternal) SetContractAdmin(setContractAdminMsg types.SetContractAdminMsg) provider.XplaClient {
	msg, err := MakeSetContractAdmintMsg(setContractAdminMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(WasmSetContractAdminMsgType, err)
	}

	return e.ToExternal(WasmSetContractAdminMsgType, msg)
}

// Migrate a wasm contract to a new code version.
func (e WasmExternal) Migrate(migrateMsg types.MigrateMsg) provider.XplaClient {
	msg, err := MakeMigrateMsg(migrateMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return e.Err(WasmMigrateMsgType, err)
	}

	return e.ToExternal(WasmMigrateMsgType, msg)
}

// Query

// Calls contract with given address with query data and prints the returned result.
func (e WasmExternal) QueryContract(queryMsg types.QueryMsg) provider.XplaClient {
	msg, err := MakeQueryMsg(queryMsg)
	if err != nil {
		return e.Err(WasmQueryContractMsgType, err)
	}

	return e.ToExternal(WasmQueryContractMsgType, msg)
}

// Query list all wasm bytecode on the chain.
func (e WasmExternal) ListCode() provider.XplaClient {
	msg, err := MakeListcodeMsg()
	if err != nil {
		return e.Err(WasmListCodeMsgType, err)
	}

	return e.ToExternal(WasmListCodeMsgType, msg)
}

// Query list wasm all bytecode on the chain for given code ID.
func (e WasmExternal) ListContractByCode(listContractByCodeMsg types.ListContractByCodeMsg) provider.XplaClient {
	msg, err := MakeListContractByCodeMsg(listContractByCodeMsg)
	if err != nil {
		return e.Err(WasmListContractByCodeMsgType, err)
	}

	return e.ToExternal(WasmListContractByCodeMsgType, msg)
}

// Downloads wasm bytecode for given code ID.
func (e WasmExternal) Download(downloadMsg types.DownloadMsg) provider.XplaClient {
	msg, err := MakeDownloadMsg(downloadMsg)
	if err != nil {
		return e.Err(WasmDownloadMsgType, err)
	}

	return e.ToExternal(WasmDownloadMsgType, msg)
}

// Prints out metadata of a code ID.
func (e WasmExternal) CodeInfo(codeInfoMsg types.CodeInfoMsg) provider.XplaClient {
	msg, err := MakeCodeInfoMsg(codeInfoMsg)
	if err != nil {
		return e.Err(WasmCodeInfoMsgType, err)
	}

	return e.ToExternal(WasmCodeInfoMsgType, msg)
}

// Prints out metadata of a contract given its address.
func (e WasmExternal) ContractInfo(contractInfoMsg types.ContractInfoMsg) provider.XplaClient {
	msg, err := MakeContractInfoMsg(contractInfoMsg)
	if err != nil {
		return e.Err(WasmContractInfoMsgType, err)
	}

	return e.ToExternal(WasmContractInfoMsgType, msg)
}

// Prints out all internal state of a contract given its address.
func (e WasmExternal) ContractStateAll(contractStateAllMsg types.ContractStateAllMsg) provider.XplaClient {
	msg, err := MakeContractStateAllMsg(contractStateAllMsg)
	if err != nil {
		return e.Err(WasmContractStateAllMsgType, err)
	}

	return e.ToExternal(WasmContractStateAllMsgType, msg)
}

// Prints out the code history for a contract given its address.
func (e WasmExternal) ContractHistory(contractHistoryMsg types.ContractHistoryMsg) provider.XplaClient {
	msg, err := MakeContractHistoryMsg(contractHistoryMsg)
	if err != nil {
		return e.Err(WasmContractHistoryMsgType, err)
	}

	return e.ToExternal(WasmContractHistoryMsgType, msg)
}

// Query list all pinned code IDs.
func (e WasmExternal) Pinned() provider.XplaClient {
	msg, err := MakePinnedMsg()
	if err != nil {
		return e.Err(WasmPinnedMsgType, err)
	}

	return e.ToExternal(WasmPinnedMsgType, msg)
}

// Get libwasmvm version.
func (e WasmExternal) LibwasmvmVersion() provider.XplaClient {
	msg, err := MakeLibwasmvmVersionMsg()
	if err != nil {
		return e.Err(WasmLibwasmvmVersionMsgType, err)
	}

	return e.ToExternal(WasmLibwasmvmVersionMsgType, msg)
}
