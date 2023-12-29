package wasm_test

import (
	mwasm "github.com/xpladev/xpla.go/core/wasm"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

var (
	testCWContractAddress = "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h"
)

func (s *IntegrationTestSuite) TestWasmTx() {
	account0 := s.network.Validators[0].AdditionalAccount

	s.xplac.WithPrivateKey(account0.PrivKey)
	// store code
	storeMsg := types.StoreMsg{
		FilePath:              testWasmFilePath,
		InstantiatePermission: "instantiate-only-sender",
	}
	s.xplac.StoreCode(storeMsg)

	makeStoreCodeMsg, err := mwasm.MakeStoreCodeMsg(storeMsg, account0.Address)
	s.Require().NoError(err)

	s.Require().Equal(makeStoreCodeMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmStoreMsgType, s.xplac.GetMsgType())

	_, err = s.xplac.StoreCode(storeMsg).CreateAndSignTx()
	s.Require().NoError(err)

	// instantiate
	instantiateMsg := types.InstantiateMsg{
		CodeId:  "1",
		Amount:  "10",
		Label:   "Contract instant",
		InitMsg: `{"owner":"` + account0.Address.String() + `"}`,
		Admin:   account0.Address.String(),
	}
	s.xplac.InstantiateContract(instantiateMsg)

	makeInstantiateMsg, err := mwasm.MakeInstantiateMsg(instantiateMsg, account0.Address)
	s.Require().NoError(err)

	s.Require().Equal(makeInstantiateMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmInstantiateMsgType, s.xplac.GetMsgType())

	wasmInstantiateContractTxbytes, err := s.xplac.InstantiateContract(instantiateMsg).CreateAndSignTx()
	s.Require().NoError(err)

	wasmInstantiateContractJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(wasmInstantiateContractTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.WasmInstantiateContractTxTemplates, string(wasmInstantiateContractJsonTxbytes))

	// execute
	executeMsg := types.ExecuteMsg{
		ContractAddress: testCWContractAddress,
		Amount:          "0",
		ExecMsg:         `{"execute_method":{"execute_key":"execute_test","execute_value":"execute_val"}}`,
	}
	s.xplac.ExecuteContract(executeMsg)

	makeExecuteMsg, err := mwasm.MakeExecuteMsg(executeMsg, account0.Address)
	s.Require().NoError(err)

	s.Require().Equal(makeExecuteMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmExecuteMsgType, s.xplac.GetMsgType())

	wasmExecuteContractTxbytes, err := s.xplac.ExecuteContract(executeMsg).CreateAndSignTx()
	s.Require().NoError(err)

	wasmExecuteContractJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(wasmExecuteContractTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.WasmExecuteContractTxTemplates, string(wasmExecuteContractJsonTxbytes))

	// clear contract admin
	clearContractAdminMsg := types.ClearContractAdminMsg{
		ContractAddress: testCWContractAddress,
	}
	s.xplac.ClearContractAdmin(clearContractAdminMsg)

	makeClearContractAdminMsg, err := mwasm.MakeClearContractAdminMsg(clearContractAdminMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeClearContractAdminMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmClearContractAdminMsgType, s.xplac.GetMsgType())

	wasmClearContractAdminTxbytes, err := s.xplac.ClearContractAdmin(clearContractAdminMsg).CreateAndSignTx()
	s.Require().NoError(err)

	wasmClearContractAdminJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(wasmClearContractAdminTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.WasmClearContractAdminTxTemplates, string(wasmClearContractAdminJsonTxbytes))

	// set contract admin
	setContractAdminMsg := types.SetContractAdminMsg{
		ContractAddress: testCWContractAddress,
	}
	s.xplac.SetContractAdmin(setContractAdminMsg)

	makeSetContractAdmintMsg, err := mwasm.MakeSetContractAdmintMsg(setContractAdminMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeSetContractAdmintMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmSetContractAdminMsgType, s.xplac.GetMsgType())

	wasmSetContractAdminTxbytes, err := s.xplac.SetContractAdmin(setContractAdminMsg).CreateAndSignTx()
	s.Require().NoError(err)

	wasmSetContractAdminJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(wasmSetContractAdminTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.WasmSetContractAdminTxTemplates, string(wasmSetContractAdminJsonTxbytes))

	// migrate
	migrateMsg := types.MigrateMsg{
		ContractAddress: testCWContractAddress,
		CodeId:          "2",
		MigrateMsg:      `{}`,
	}
	s.xplac.Migrate(migrateMsg)

	makeMigrateMsg, err := mwasm.MakeMigrateMsg(migrateMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeMigrateMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmMigrateMsgType, s.xplac.GetMsgType())

	wasmMigrateTxbytes, err := s.xplac.Migrate(migrateMsg).CreateAndSignTx()
	s.Require().NoError(err)

	wasmMigrateJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(wasmMigrateTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.WasmMigrateTxTemplates, string(wasmMigrateJsonTxbytes))
}

func (s *IntegrationTestSuite) TestWasm() {
	// call contract
	queryMsg := types.QueryMsg{
		ContractAddress: testCWContractAddress,
		QueryMsg:        `{"query_method":{"query":"query_test"}}`,
	}
	s.xplac.QueryContract(queryMsg)

	makeQueryMsg, err := mwasm.MakeQueryMsg(queryMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmQueryContractMsgType, s.xplac.GetMsgType())

	// list code
	s.xplac.ListCode()

	makeListcodeMsg, err := mwasm.MakeListcodeMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeListcodeMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmListCodeMsgType, s.xplac.GetMsgType())

	// list contract by code
	listContractByCodeMsg := types.ListContractByCodeMsg{
		CodeId: "1",
	}
	s.xplac.ListContractByCode(listContractByCodeMsg)

	makeListContractByCodeMsg, err := mwasm.MakeListContractByCodeMsg(listContractByCodeMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeListContractByCodeMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmListContractByCodeMsgType, s.xplac.GetMsgType())

	// download
	downloadMsg := types.DownloadMsg{
		CodeId:           "1",
		DownloadFileName: "./example.json",
	}
	s.xplac.Download(downloadMsg)

	makeDownloadMsg, err := mwasm.MakeDownloadMsg(downloadMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeDownloadMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmDownloadMsgType, s.xplac.GetMsgType())

	// code info
	codeInfoMsg := types.CodeInfoMsg{
		CodeId: "1",
	}
	s.xplac.CodeInfo(codeInfoMsg)

	makeCodeInfoMsg, err := mwasm.MakeCodeInfoMsg(codeInfoMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeCodeInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmCodeInfoMsgType, s.xplac.GetMsgType())

	// contract info
	contractInfoMsg := types.ContractInfoMsg{
		ContractAddress: testCWContractAddress,
	}
	s.xplac.ContractInfo(contractInfoMsg)

	makeContractInfoMsg, err := mwasm.MakeContractInfoMsg(contractInfoMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeContractInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmContractInfoMsgType, s.xplac.GetMsgType())

	// contract state all
	contractStateAllMsg := types.ContractStateAllMsg{
		ContractAddress: testCWContractAddress,
	}
	s.xplac.ContractStateAll(contractStateAllMsg)

	makeContractStateAllMsg, err := mwasm.MakeContractStateAllMsg(contractStateAllMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeContractStateAllMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmContractStateAllMsgType, s.xplac.GetMsgType())

	// contract history
	contractHistoryMsg := types.ContractHistoryMsg{
		ContractAddress: testCWContractAddress,
	}
	s.xplac.ContractHistory(contractHistoryMsg)

	makeContractHistoryMsg, err := mwasm.MakeContractHistoryMsg(contractHistoryMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeContractHistoryMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmContractHistoryMsgType, s.xplac.GetMsgType())

	// pinned
	s.xplac.Pinned()

	makePinnedMsg, err := mwasm.MakePinnedMsg()
	s.Require().NoError(err)

	s.Require().Equal(makePinnedMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmPinnedMsgType, s.xplac.GetMsgType())

	// libwasmvm version
	s.xplac.LibwasmvmVersion()

	makeLibwasmvmVersionMsg, err := mwasm.MakeLibwasmvmVersionMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeLibwasmvmVersionMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmLibwasmvmVersionMsgType, s.xplac.GetMsgType())
}
