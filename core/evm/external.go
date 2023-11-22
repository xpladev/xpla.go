package evm

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type EvmExternal struct {
	Xplac provider.XplaClient
}

func NewEvmExternal(xplac provider.XplaClient) (e EvmExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Send coind by using evm client.
func (e EvmExternal) EvmSendCoin(sendCoinMsg types.SendCoinMsg) provider.XplaClient {
	msg, err := MakeSendCoinMsg(sendCoinMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmSendCoinMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Deploy soldity contract.
func (e EvmExternal) DeploySolidityContract(deploySolContractMsg types.DeploySolContractMsg) provider.XplaClient {
	msg, err := MakeDeploySolContractMsg(deploySolContractMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmDeploySolContractMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Invoke (as execute) solidity contract.
func (e EvmExternal) InvokeSolidityContract(invokeSolContractMsg types.InvokeSolContractMsg) provider.XplaClient {
	msg, err := MakeInvokeSolContractMsg(invokeSolContractMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmInvokeSolContractMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Call(as query) solidity contract.
func (e EvmExternal) CallSolidityContract(callSolContractMsg types.CallSolContractMsg) provider.XplaClient {
	msg, err := MakeCallSolContractMsg(callSolContractMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmCallSolContractMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query a transaction which is ethereum type information by retrieving hash.
func (e EvmExternal) GetTransactionByHash(getTransactionByHashMsg types.GetTransactionByHashMsg) provider.XplaClient {
	msg, err := MakeGetTransactionByHashMsg(getTransactionByHashMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmGetTransactionByHashMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query a block which is ethereum type information by retrieving hash or block height(as number).
func (e EvmExternal) GetBlockByHashOrHeight(getBlockByHashHeightMsg types.GetBlockByHashHeightMsg) provider.XplaClient {
	msg, err := MakeGetBlockByHashHeightMsg(getBlockByHashHeightMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmGetBlockByHashHeightMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query a account information which includes account address(hex and bech32), balance and etc.
func (e EvmExternal) AccountInfo(accountInfoMsg types.AccountInfoMsg) provider.XplaClient {
	msg, err := MakeQueryAccountInfoMsg(accountInfoMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmQueryAccountInfoMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query suggested gas price.
func (e EvmExternal) SuggestGasPrice() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmSuggestGasPriceMsgType).
		WithMsg(nil)
	return e.Xplac
}

// Query chain ID of ethereum type.
func (e EvmExternal) EthChainID() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmQueryChainIdMsgType).
		WithMsg(nil)
	return e.Xplac
}

// Query latest block height(as number)
func (e EvmExternal) EthBlockNumber() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmQueryCurrentBlockNumberMsgType).
		WithMsg(nil)
	return e.Xplac
}

// Query web3 client version.
func (e EvmExternal) Web3ClientVersion() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmWeb3ClientVersionMsgType).
		WithMsg(nil)
	return e.Xplac
}

// Query web3 sha3.
func (e EvmExternal) Web3Sha3(web3Sha3Msg types.Web3Sha3Msg) provider.XplaClient {
	msg, err := MakeWeb3Sha3Msg(web3Sha3Msg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmWeb3Sha3MsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query current network ID.
func (e EvmExternal) NetVersion() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmNetVersionMsgType).
		WithMsg(nil)
	return e.Xplac
}

// Query the number of peers currently connected to the client.
func (e EvmExternal) NetPeerCount() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmNetPeerCountMsgType).
		WithMsg(nil)
	return e.Xplac
}

// actively listening for network connections.
func (e EvmExternal) NetListening() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmNetListeningMsgType).
		WithMsg(nil)
	return e.Xplac
}

// ethereum protocol version.
func (e EvmExternal) EthProtocolVersion() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthProtocolVersionMsgType).
		WithMsg(nil)
	return e.Xplac
}

// Query the sync status object depending on the details of tendermint's sync protocol.
func (e EvmExternal) EthSyncing() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthSyncingMsgType).
		WithMsg(nil)
	return e.Xplac
}

// Query all eth accounts.
func (e EvmExternal) EthAccounts() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthAccountsMsgType).
		WithMsg(nil)
	return e.Xplac
}

// Query the number of transaction a given block.
func (e EvmExternal) EthGetBlockTransactionCount(ethGetBlockTransactionCountMsg types.EthGetBlockTransactionCountMsg) provider.XplaClient {
	if ethGetBlockTransactionCountMsg.BlockHash == "" && ethGetBlockTransactionCountMsg.BlockHeight == "" {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInsufficientParams, "cannot query, without block hash or height parameter"))
	}

	if ethGetBlockTransactionCountMsg.BlockHash != "" && ethGetBlockTransactionCountMsg.BlockHeight != "" {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "select only one parameter, block hash OR height"))
	}

	msg, err := MakeEthGetBlockTransactionCountMsg(ethGetBlockTransactionCountMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthGetBlockTransactionCountMsgType).
		WithMsg(msg)

	return e.Xplac
}

// Query estimate gas.
func (e EvmExternal) EstimateGas(invokeSolContractMsg types.InvokeSolContractMsg) provider.XplaClient {
	msg, err := MakeEstimateGasSolMsg(invokeSolContractMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthEstimateGasMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query transaction by block hash and index.
func (e EvmExternal) EthGetTransactionByBlockHashAndIndex(getTransactionByBlockHashAndIndexMsg types.GetTransactionByBlockHashAndIndexMsg) provider.XplaClient {
	msg, err := MakeGetTransactionByBlockHashAndIndexMsg(getTransactionByBlockHashAndIndexMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmGetTransactionByBlockHashAndIndexMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query transaction receipt.
func (e EvmExternal) EthGetTransactionReceipt(getTransactionReceiptMsg types.GetTransactionReceiptMsg) provider.XplaClient {
	msg, err := MakeGetTransactionReceiptMsg(getTransactionReceiptMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmGetTransactionReceiptMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query filter ID by eth new filter.
func (e EvmExternal) EthNewFilter(ethNewFilterMsg types.EthNewFilterMsg) provider.XplaClient {
	msg, err := MakeEthNewFilterMsg(ethNewFilterMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthNewFilterMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query filter ID by eth new block filter.
func (e EvmExternal) EthNewBlockFilter() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthNewBlockFilterMsgType).
		WithMsg(nil)
	return e.Xplac
}

// Query filter ID by eth new pending transaction filter.
func (e EvmExternal) EthNewPendingTransactionFilter() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthNewPendingTransactionFilterMsgType).
		WithMsg(nil)
	return e.Xplac
}

// Uninstall filter.
func (e EvmExternal) EthUninstallFilter(ethUninstallFilterMsg types.EthUninstallFilterMsg) provider.XplaClient {
	msg, err := MakeEthUninstallFilterMsg(ethUninstallFilterMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthUninstallFilterMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query filter changes.
func (e EvmExternal) EthGetFilterChanges(ethGetFilterChangesMsg types.EthGetFilterChangesMsg) provider.XplaClient {
	msg, err := MakeEthGetFilterChangesMsg(ethGetFilterChangesMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthGetFilterChangesMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query filter logs.
func (e EvmExternal) EthGetFilterLogs(ethGetFilterLogsMsg types.EthGetFilterLogsMsg) provider.XplaClient {
	msg, err := MakeEthGetFilterLogsMsg(ethGetFilterLogsMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthGetFilterLogsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Get logs.
func (e EvmExternal) EthGetLogs(ethGetLogsMsg types.EthGetLogsMsg) provider.XplaClient {
	msg, err := MakeEthGetLogsMsg(ethGetLogsMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthGetLogsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query coinbase.
func (e EvmExternal) EthCoinbase() provider.XplaClient {
	e.Xplac.WithModule(EvmModule).
		WithMsgType(EvmEthCoinbaseMsgType).
		WithMsg(nil)
	return e.Xplac
}
