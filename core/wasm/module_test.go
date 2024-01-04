package wasm_test

import (
	"math/rand"

	"github.com/xpladev/xpla.go/core/wasm"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	account0 := s.network.Validators[0].AdditionalAccount

	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := wasm.NewCoreModule()

	// test get name
	s.Require().Equal(wasm.WasmModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// store code
	storeMsg := types.StoreMsg{
		FilePath:              testWasmFilePath,
		InstantiatePermission: "instantiate-only-sender",
	}

	makeStoreCodeMsg, err := wasm.MakeStoreCodeMsg(storeMsg, account0.Address)
	s.Require().NoError(err)

	testMsg = makeStoreCodeMsg
	txBuilder, err = c.NewTxRouter(s.xplac.GetLogger(), txBuilder, wasm.WasmStoreMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeStoreCodeMsg, txBuilder.GetTx().GetMsgs()[0])

	// instantiate
	instantiateMsg := types.InstantiateMsg{
		CodeId:  "1",
		Amount:  "10",
		Label:   "Contract instant",
		InitMsg: `{"owner":"` + account0.Address.String() + `"}`,
		Admin:   account0.Address.String(),
	}

	makeInstantiateMsg, err := wasm.MakeInstantiateMsg(instantiateMsg, account0.Address)
	s.Require().NoError(err)

	testMsg = makeInstantiateMsg
	txBuilder, err = c.NewTxRouter(s.xplac.GetLogger(), txBuilder, wasm.WasmInstantiateMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeInstantiateMsg, txBuilder.GetTx().GetMsgs()[0])

	// execute
	testCWContractAddress := "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h"
	executeMsg := types.ExecuteMsg{
		ContractAddress: testCWContractAddress,
		Amount:          "0",
		ExecMsg:         `{"execute_method":{"execute_key":"execute_test","execute_value":"execute_val"}}`,
	}

	makeExecuteMsg, err := wasm.MakeExecuteMsg(executeMsg, account0.Address)
	s.Require().NoError(err)

	testMsg = makeExecuteMsg
	txBuilder, err = c.NewTxRouter(s.xplac.GetLogger(), txBuilder, wasm.WasmExecuteMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeExecuteMsg, txBuilder.GetTx().GetMsgs()[0])

	// clear contract admin
	clearContractAdminMsg := types.ClearContractAdminMsg{
		ContractAddress: testCWContractAddress,
	}

	makeClearContractAdminMsg, err := wasm.MakeClearContractAdminMsg(clearContractAdminMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	testMsg = makeClearContractAdminMsg
	txBuilder, err = c.NewTxRouter(s.xplac.GetLogger(), txBuilder, wasm.WasmClearContractAdminMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeClearContractAdminMsg, txBuilder.GetTx().GetMsgs()[0])

	// set contract admin
	setContractAdminMsg := types.SetContractAdminMsg{
		ContractAddress: testCWContractAddress,
	}

	makeSetContractAdmintMsg, err := wasm.MakeSetContractAdmintMsg(setContractAdminMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	testMsg = makeSetContractAdmintMsg
	txBuilder, err = c.NewTxRouter(s.xplac.GetLogger(), txBuilder, wasm.WasmSetContractAdminMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeSetContractAdmintMsg, txBuilder.GetTx().GetMsgs()[0])

	// migrate
	migrateMsg := types.MigrateMsg{
		ContractAddress: testCWContractAddress,
		CodeId:          "2",
		MigrateMsg:      `{}`,
	}

	makeMigrateMsg, err := wasm.MakeMigrateMsg(migrateMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	testMsg = makeMigrateMsg
	txBuilder, err = c.NewTxRouter(s.xplac.GetLogger(), txBuilder, wasm.WasmMigrateMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeMigrateMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(s.xplac.GetLogger(), nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = provider.ResetXplac(s.xplac)
}
