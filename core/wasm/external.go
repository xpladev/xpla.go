package wasm

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

type WasmExternal struct {
	Xplac provider.XplaClient
}

func NewWasmExternal(xplac provider.XplaClient) (e WasmExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Upload a wasm binary.
func (e WasmExternal) StoreCode(storeMsg types.StoreMsg) provider.XplaClient {
	msg, err := MakeStoreCodeMsg(storeMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmStoreMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Instantiate a wasm contract.
func (e WasmExternal) InstantiateContract(instantiageMsg types.InstantiateMsg) provider.XplaClient {
	msg, err := MakeInstantiateMsg(instantiageMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmInstantiateMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Execute a wasm contract.
func (e WasmExternal) ExecuteContract(executeMsg types.ExecuteMsg) provider.XplaClient {
	msg, err := MakeExecuteMsg(executeMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmExecuteMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Clears admin for a contract to prevent further migrations.
func (e WasmExternal) ClearContractAdmin(clearContractAdminMsg types.ClearContractAdminMsg) provider.XplaClient {
	msg, err := MakeClearContractAdminMsg(clearContractAdminMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmClearContractAdminMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Set new admin for a contract.
func (e WasmExternal) SetContractAdmin(setContractAdminMsg types.SetContractAdminMsg) provider.XplaClient {
	msg, err := MakeSetContractAdmintMsg(setContractAdminMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmSetContractAdminMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Migrate a wasm contract to a new code version.
func (e WasmExternal) Migrate(migrateMsg types.MigrateMsg) provider.XplaClient {
	msg, err := MakeMigrateMsg(migrateMsg, e.Xplac.GetFromAddress())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmMigrateMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Calls contract with given address with query data and prints the returned result.
func (e WasmExternal) QueryContract(queryMsg types.QueryMsg) provider.XplaClient {
	msg, err := MakeQueryMsg(queryMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmQueryContractMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query list all wasm bytecode on the chain.
func (e WasmExternal) ListCode() provider.XplaClient {
	msg, err := MakeListcodeMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmListCodeMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query list wasm all bytecode on the chain for given code ID.
func (e WasmExternal) ListContractByCode(listContractByCodeMsg types.ListContractByCodeMsg) provider.XplaClient {
	msg, err := MakeListContractByCodeMsg(listContractByCodeMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmListContractByCodeMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Downloads wasm bytecode for given code ID.
func (e WasmExternal) Download(downloadMsg types.DownloadMsg) provider.XplaClient {
	msg, err := MakeDownloadMsg(downloadMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmDownloadMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Prints out metadata of a code ID.
func (e WasmExternal) CodeInfo(codeInfoMsg types.CodeInfoMsg) provider.XplaClient {
	msg, err := MakeCodeInfoMsg(codeInfoMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmCodeInfoMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Prints out metadata of a contract given its address.
func (e WasmExternal) ContractInfo(contractInfoMsg types.ContractInfoMsg) provider.XplaClient {
	msg, err := MakeContractInfoMsg(contractInfoMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmContractInfoMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Prints out all internal state of a contract given its address.
func (e WasmExternal) ContractStateAll(contractStateAllMsg types.ContractStateAllMsg) provider.XplaClient {
	msg, err := MakeContractStateAllMsg(contractStateAllMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmContractStateAllMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Prints out the code history for a contract given its address.
func (e WasmExternal) ContractHistory(contractHistoryMsg types.ContractHistoryMsg) provider.XplaClient {
	msg, err := MakeContractHistoryMsg(contractHistoryMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmContractHistoryMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query list all pinned code IDs.
func (e WasmExternal) Pinned() provider.XplaClient {
	msg, err := MakePinnedMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmPinnedMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Get libwasmvm version.
func (e WasmExternal) LibwasmvmVersion() provider.XplaClient {
	msg, err := MakeLibwasmvmVersionMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(WasmModule).
		WithMsgType(WasmLibwasmvmVersionMsgType).
		WithMsg(msg)
	return e.Xplac
}
