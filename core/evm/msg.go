package evm

import (
	"github.com/xpladev/xpla.go/types"
)

// (Tx) make msg - send coin
func MakeSendCoinMsg(sendCoinMsg types.SendCoinMsg) (types.SendCoinMsg, error) {
	return parseSendCoinArgs(sendCoinMsg)
}

// (Tx) make msg - deploy solidity contract
func MakeDeploySolContractMsg(deploySolContractMsg types.DeploySolContractMsg) (ContractInfo, error) {
	return parseDeploySolContractArgs(deploySolContractMsg)
}

// (Tx) make msg - invoke solidity contract
func MakeInvokeSolContractMsg(InvokeSolContractMsg types.InvokeSolContractMsg) (types.InvokeSolContractMsg, error) {
	return parseInvokeSolContractArgs(InvokeSolContractMsg)
}

// (Query) make msg - call solidity contract
func MakeCallSolContractMsg(callSolContractMsg types.CallSolContractMsg) (CallSolContractParseMsg, error) {
	return parseCallSolContractArgs(callSolContractMsg)
}

// (Query) make msg - transaction by hash
func MakeGetTransactionByHashMsg(getTransactionByHashMsg types.GetTransactionByHashMsg) (types.GetTransactionByHashMsg, error) {
	return getTransactionByHashMsg, nil
}

// (Query) make msg - block by hash or height
func MakeGetBlockByHashHeightMsg(getBlockByHashHeightMsg types.GetBlockByHashHeightMsg) (types.GetBlockByHashHeightMsg, error) {
	return parseGetBlockByHashHeightArgs(getBlockByHashHeightMsg)
}

// (Query) make msg - account info
func MakeQueryAccountInfoMsg(accountInfoMsg types.AccountInfoMsg) (types.AccountInfoMsg, error) {
	return accountInfoMsg, nil
}

// (Query) make msg - web3 sha3
func MakeWeb3Sha3Msg(web3Sha3Msg types.Web3Sha3Msg) (types.Web3Sha3Msg, error) {
	return web3Sha3Msg, nil
}

// (Query) make msg - get transaction count of the block number
func MakeEthGetBlockTransactionCountMsg(ethGetBlockTransactionCountMsg types.EthGetBlockTransactionCountMsg) (types.EthGetBlockTransactionCountMsg, error) {
	return parseEthGetBlockTransactionCountArgs(ethGetBlockTransactionCountMsg)
}

// (Query) make msg - sol contract estimate gas
func MakeEstimateGasSolMsg(invokeSolContractMsg types.InvokeSolContractMsg) (CallSolContractParseMsg, error) {
	return parseEstimateGasSolArgs(invokeSolContractMsg)
}

// (Query) make msg - get transaction by block hash and index
func MakeGetTransactionByBlockHashAndIndexMsg(getTransactionByBlockHashAndIndexMsg types.GetTransactionByBlockHashAndIndexMsg) (types.GetTransactionByBlockHashAndIndexMsg, error) {
	return getTransactionByBlockHashAndIndexMsg, nil
}

// (Query) make msg - get transaction receipt
func MakeGetTransactionReceiptMsg(getTransactionReceiptMsg types.GetTransactionReceiptMsg) (types.GetTransactionReceiptMsg, error) {
	return getTransactionReceiptMsg, nil
}

// (Query) make msg - eth new filter
func MakeEthNewFilterMsg(ethNewFilterMsg types.EthNewFilterMsg) (EthNewFilterParseMsg, error) {
	return parseEthNewFilterArgs(ethNewFilterMsg)
}

// (Query) make msg - eth uninstall filter
func MakeEthUninstallFilterMsg(ethUninstallFilter types.EthUninstallFilterMsg) (types.EthUninstallFilterMsg, error) {
	return ethUninstallFilter, nil
}

// (Query) make msg - eth get filter changes
func MakeEthGetFilterChangesMsg(ethGetFilterChangesMsg types.EthGetFilterChangesMsg) (types.EthGetFilterChangesMsg, error) {
	return ethGetFilterChangesMsg, nil
}

// (Query) make msg - eth get filter logs
func MakeEthGetFilterLogsMsg(ethGetFilterLogsMsg types.EthGetFilterLogsMsg) (types.EthGetFilterLogsMsg, error) {
	return ethGetFilterLogsMsg, nil
}

// (Query) make msg - eth get logs
func MakeEthGetLogsMsg(ethGetLogsMsg types.EthGetLogsMsg) (EthNewFilterParseMsg, error) {
	return parseEthGetLogsArgs(ethGetLogsMsg)
}
