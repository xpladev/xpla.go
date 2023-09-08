package evm_test

import (
	mevm "github.com/xpladev/xpla.go/core/evm"
	"github.com/xpladev/xpla.go/types"
)

var (
	testSolContractAddress = "0x80E123317190cAf36292A04776b0De020136526F"
	testABIPath            = "../../util/testutil/test_files/abi.json"
	testBytecodePath       = "../../util/testutil/test_files/bytecode.json"
	testTxHash             = "B6BBBB649F19E8970EF274C0083FE945FD38AD8C524D68BB3FE3A20D72DF03C4"
)

func (s *IntegrationTestSuite) TestEvmTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// send evm coin
	sendCoinMsg := types.SendCoinMsg{
		FromAddress: s.accounts[0].PubKey.Address().String(),
		ToAddress:   s.accounts[1].PubKey.Address().String(),
		Amount:      "1000",
	}
	s.xplac.EvmSendCoin(sendCoinMsg)

	makeSendCoinMsg, err := mevm.MakeSendCoinMsg(sendCoinMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeSendCoinMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmSendCoinMsgType, s.xplac.GetMsgType())

	// deploy solidity contract
	deploySolContractMsg := types.DeploySolContractMsg{
		ABIJsonFilePath:      testABIPath,
		BytecodeJsonFilePath: testBytecodePath,
		Args:                 nil,
	}
	s.xplac.DeploySolidityContract(deploySolContractMsg)

	makeDeploySolContractMsg, err := mevm.MakeDeploySolContractMsg(deploySolContractMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeDeploySolContractMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmDeploySolContractMsgType, s.xplac.GetMsgType())

	// invoke solidity contract
	invokeSolContractMsg := types.InvokeSolContractMsg{
		ContractAddress:      testSolContractAddress,
		ABIJsonFilePath:      testABIPath,
		BytecodeJsonFilePath: testBytecodePath,
		Args:                 nil,
	}
	s.xplac.InvokeSolidityContract(invokeSolContractMsg)

	makeInvokeSolContractMsg, err := mevm.MakeInvokeSolContractMsg(invokeSolContractMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeInvokeSolContractMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmInvokeSolContractMsgType, s.xplac.GetMsgType())
}

func (s *IntegrationTestSuite) TestEvm() {
	// call contract
	callSolContractMsg := types.CallSolContractMsg{
		ContractAddress:      testSolContractAddress,
		ContractFuncCallName: "retrieve",
		ABIJsonFilePath:      testABIPath,
		BytecodeJsonFilePath: testBytecodePath,
		FromByteAddress:      s.accounts[0].PubKey.Address().String(),
	}
	s.xplac.CallSolidityContract(callSolContractMsg)

	makeCallSolContractMsg, err := mevm.MakeCallSolContractMsg(callSolContractMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeCallSolContractMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmCallSolContractMsgType, s.xplac.GetMsgType())

	// tx by hash
	getTransactionByHashMsg := types.GetTransactionByHashMsg{
		TxHash: testTxHash,
	}
	s.xplac.GetTransactionByHash(getTransactionByHashMsg)

	makeGetTransactionByHashMsg, err := mevm.MakeGetTransactionByHashMsg(getTransactionByHashMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGetTransactionByHashMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmGetTransactionByHashMsgType, s.xplac.GetMsgType())

	// block by hash or height
	getBlockByHashHeightMsg := types.GetBlockByHashHeightMsg{
		BlockHeight: "1",
	}
	s.xplac.GetBlockByHashOrHeight(getBlockByHashHeightMsg)

	makeGetBlockByHashHeightMsg, err := mevm.MakeGetBlockByHashHeightMsg(getBlockByHashHeightMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGetBlockByHashHeightMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmGetBlockByHashHeightMsgType, s.xplac.GetMsgType())

	// account info
	accountInfoMsg := types.AccountInfoMsg{
		Account: s.accounts[0].PubKey.Address().String(),
	}
	s.xplac.AccountInfo(accountInfoMsg)

	makeQueryAccountInfoMsg, err := mevm.MakeQueryAccountInfoMsg(accountInfoMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryAccountInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmQueryAccountInfoMsgType, s.xplac.GetMsgType())

	// suggest gas price
	s.xplac.SuggestGasPrice()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmSuggestGasPriceMsgType, s.xplac.GetMsgType())

	// eth chain ID
	s.xplac.EthChainID()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmQueryChainIdMsgType, s.xplac.GetMsgType())

	// eth block number
	s.xplac.EthBlockNumber()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmQueryCurrentBlockNumberMsgType, s.xplac.GetMsgType())

	// web3 client version
	s.xplac.Web3ClientVersion()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmWeb3ClientVersionMsgType, s.xplac.GetMsgType())

	// web3 sha3
	web3Sha3Msg := types.Web3Sha3Msg{
		InputParam: "ABC",
	}
	s.xplac.Web3Sha3(web3Sha3Msg)

	makeWeb3Sha3Msg, err := mevm.MakeWeb3Sha3Msg(web3Sha3Msg)
	s.Require().NoError(err)

	s.Require().Equal(makeWeb3Sha3Msg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmWeb3Sha3MsgType, s.xplac.GetMsgType())

	// net version
	s.xplac.NetVersion()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmNetVersionMsgType, s.xplac.GetMsgType())

	// net peer count
	s.xplac.NetPeerCount()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmNetPeerCountMsgType, s.xplac.GetMsgType())

	// net listening
	s.xplac.NetListening()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmNetListeningMsgType, s.xplac.GetMsgType())

	// eth protocol version
	s.xplac.EthProtocolVersion()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthProtocolVersionMsgType, s.xplac.GetMsgType())

	// eth syncing
	s.xplac.EthSyncing()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthSyncingMsgType, s.xplac.GetMsgType())

	// eth accounts
	s.xplac.EthAccounts()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthAccountsMsgType, s.xplac.GetMsgType())

	// eth get block transaction count
	ethGetBlockTransactionCountMsg := types.EthGetBlockTransactionCountMsg{
		BlockHeight: "1",
	}
	s.xplac.EthGetBlockTransactionCount(ethGetBlockTransactionCountMsg)

	makeEthGetBlockTransactionCountMsg, err := mevm.MakeEthGetBlockTransactionCountMsg(ethGetBlockTransactionCountMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthGetBlockTransactionCountMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthGetBlockTransactionCountMsgType, s.xplac.GetMsgType())

	// estimate gas
	invokeSolContractMsg := types.InvokeSolContractMsg{
		ContractAddress:      testSolContractAddress,
		ABIJsonFilePath:      testABIPath,
		BytecodeJsonFilePath: testBytecodePath,
		FromByteAddress:      s.accounts[0].PubKey.Address().String(),
	}
	s.xplac.EstimateGas(invokeSolContractMsg)

	makeEstimateGasSolMsg, err := mevm.MakeEstimateGasSolMsg(invokeSolContractMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEstimateGasSolMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthEstimateGasMsgType, s.xplac.GetMsgType())

	// get tx by block hash and index
	getTransactionByBlockHashAndIndexMsg := types.GetTransactionByBlockHashAndIndexMsg{
		BlockHash: "1",
		Index:     "0",
	}
	s.xplac.EthGetTransactionByBlockHashAndIndex(getTransactionByBlockHashAndIndexMsg)

	makeGetTransactionByBlockHashAndIndexMsg, err := mevm.MakeGetTransactionByBlockHashAndIndexMsg(getTransactionByBlockHashAndIndexMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGetTransactionByBlockHashAndIndexMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmGetTransactionByBlockHashAndIndexMsgType, s.xplac.GetMsgType())

	// tx receipt
	getTransactionReceiptMsg := types.GetTransactionReceiptMsg{
		TransactionHash: testTxHash,
	}
	s.xplac.EthGetTransactionReceipt(getTransactionReceiptMsg)

	makeGetTransactionReceiptMsg, err := mevm.MakeGetTransactionReceiptMsg(getTransactionReceiptMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGetTransactionReceiptMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmGetTransactionReceiptMsgType, s.xplac.GetMsgType())

	// eth new fileter
	ethNewFilterMsg := types.EthNewFilterMsg{
		Topics:    []string{"0x20ec56d16231c4d7f761c2533885619489fface85cf6c478868ef1d531b93177"},
		Address:   []string{"0xf7777b36a51fb0b33dd0c5118361AfC94ff7f967"},
		ToBlock:   "latest",
		FromBlock: "earliest",
	}
	s.xplac.EthNewFilter(ethNewFilterMsg)

	makeEthNewFilterMsg, err := mevm.MakeEthNewFilterMsg(ethNewFilterMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthNewFilterMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthNewFilterMsgType, s.xplac.GetMsgType())

	// new block filter
	s.xplac.EthNewBlockFilter()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthNewBlockFilterMsgType, s.xplac.GetMsgType())

	// new pending transaction filter
	s.xplac.EthNewPendingTransactionFilter()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthNewPendingTransactionFilterMsgType, s.xplac.GetMsgType())

	// uninstall filter
	ethUninstallFilterMsg := types.EthUninstallFilterMsg{
		FilterId: "0x168b9d421ecbffa1ac706926c2203454",
	}
	s.xplac.EthUninstallFilter(ethUninstallFilterMsg)

	makeEthUninstallFilterMsg, err := mevm.MakeEthUninstallFilterMsg(ethUninstallFilterMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthUninstallFilterMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthUninstallFilterMsgType, s.xplac.GetMsgType())

	// filter changes
	ethGetFilterChangesMsg := types.EthGetFilterChangesMsg{
		FilterId: "0x168b9d421ecbffa1ac706926c2203454",
	}
	s.xplac.EthGetFilterChanges(ethGetFilterChangesMsg)

	makeEthGetFilterChangesMsg, err := mevm.MakeEthGetFilterChangesMsg(ethGetFilterChangesMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthGetFilterChangesMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthGetFilterChangesMsgType, s.xplac.GetMsgType())

	// eth filter logs
	ethGetFilterLogsMsg := types.EthGetFilterLogsMsg{
		FilterId: "0x168b9d421ecbffa1ac706926c2203454",
	}
	s.xplac.EthGetFilterLogs(ethGetFilterLogsMsg)

	makeEthGetFilterLogsMsg, err := mevm.MakeEthGetFilterLogsMsg(ethGetFilterLogsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthGetFilterLogsMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthGetFilterLogsMsgType, s.xplac.GetMsgType())

	// eth logs
	ethGetLogsMsg := types.EthGetLogsMsg{
		Topics:    []string{"0x20ec56d16231c4d7f761c2533885619489fface85cf6c478868ef1d531b93177"},
		Address:   []string{"0xf7777b36a51fb0b33dd0c5118361AfC94ff7f967"},
		ToBlock:   "latest",
		FromBlock: "latest",
	}
	s.xplac.EthGetLogs(ethGetLogsMsg)

	makeEthGetLogsMsg, err := mevm.MakeEthGetLogsMsg(ethGetLogsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthGetLogsMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthGetLogsMsgType, s.xplac.GetMsgType())

	// coinbase
	s.xplac.EthCoinbase()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthCoinbaseMsgType, s.xplac.GetMsgType())
}
