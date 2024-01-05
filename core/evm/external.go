package evm

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &EvmExternal{}

type EvmExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e EvmExternal) {
	e.Xplac = xplac
	e.Name = EvmModule
	return e
}

func (e EvmExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e EvmExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Send coind by using evm client.
func (e EvmExternal) EvmSendCoin(sendCoinMsg types.SendCoinMsg) provider.XplaClient {
	msg, err := MakeSendCoinMsg(sendCoinMsg)
	if err != nil {
		return e.Err(EvmSendCoinMsgType, err)
	}

	return e.ToExternal(EvmSendCoinMsgType, msg)
}

// Deploy soldity contract.
func (e EvmExternal) DeploySolidityContract(deploySolContractMsg types.DeploySolContractMsg) provider.XplaClient {
	msg, err := MakeDeploySolContractMsg(deploySolContractMsg)
	if err != nil {
		return e.Err(EvmDeploySolContractMsgType, err)
	}

	return e.ToExternal(EvmDeploySolContractMsgType, msg)
}

// Invoke (as execute) solidity contract.
func (e EvmExternal) InvokeSolidityContract(invokeSolContractMsg types.InvokeSolContractMsg) provider.XplaClient {
	msg, err := MakeInvokeSolContractMsg(invokeSolContractMsg)
	if err != nil {
		return e.Err(EvmInvokeSolContractMsgType, err)
	}

	return e.ToExternal(EvmInvokeSolContractMsgType, msg)
}

// Query

// Call(as query) solidity contract.
func (e EvmExternal) CallSolidityContract(callSolContractMsg types.CallSolContractMsg) provider.XplaClient {
	msg, err := MakeCallSolContractMsg(callSolContractMsg)
	if err != nil {
		return e.Err(EvmCallSolContractMsgType, err)
	}

	return e.ToExternal(EvmCallSolContractMsgType, msg)
}

// Query a transaction which is ethereum type information by retrieving hash.
func (e EvmExternal) GetTransactionByHash(getTransactionByHashMsg types.GetTransactionByHashMsg) provider.XplaClient {
	msg, err := MakeGetTransactionByHashMsg(getTransactionByHashMsg)
	if err != nil {
		return e.Err(EvmGetTransactionByHashMsgType, err)
	}

	return e.ToExternal(EvmGetTransactionByHashMsgType, msg)
}

// Query a block which is ethereum type information by retrieving hash or block height(as number).
func (e EvmExternal) GetBlockByHashOrHeight(getBlockByHashHeightMsg types.GetBlockByHashHeightMsg) provider.XplaClient {
	msg, err := MakeGetBlockByHashHeightMsg(getBlockByHashHeightMsg)
	if err != nil {
		return e.Err(EvmGetBlockByHashHeightMsgType, err)
	}

	return e.ToExternal(EvmGetBlockByHashHeightMsgType, msg)
}

// Query a account information which includes account address(hex and bech32), balance and etc.
func (e EvmExternal) AccountInfo(accountInfoMsg types.AccountInfoMsg) provider.XplaClient {
	msg, err := MakeQueryAccountInfoMsg(accountInfoMsg)
	if err != nil {
		return e.Err(EvmQueryAccountInfoMsgType, err)
	}

	return e.ToExternal(EvmQueryAccountInfoMsgType, msg)
}

// Query suggested gas price.
func (e EvmExternal) SuggestGasPrice() provider.XplaClient {
	return e.ToExternal(EvmSuggestGasPriceMsgType, nil)
}

// Query chain ID of ethereum type.
func (e EvmExternal) EthChainID() provider.XplaClient {
	return e.ToExternal(EvmQueryChainIdMsgType, nil)
}

// Query latest block height(as number)
func (e EvmExternal) EthBlockNumber() provider.XplaClient {
	return e.ToExternal(EvmQueryCurrentBlockNumberMsgType, nil)
}

// Query web3 client version.
func (e EvmExternal) Web3ClientVersion() provider.XplaClient {
	return e.ToExternal(EvmWeb3ClientVersionMsgType, nil)
}

// Query web3 sha3.
func (e EvmExternal) Web3Sha3(web3Sha3Msg types.Web3Sha3Msg) provider.XplaClient {
	msg, err := MakeWeb3Sha3Msg(web3Sha3Msg)
	if err != nil {
		return e.Err(EvmWeb3Sha3MsgType, err)
	}

	return e.ToExternal(EvmWeb3Sha3MsgType, msg)
}

// Query current network ID.
func (e EvmExternal) NetVersion() provider.XplaClient {
	return e.ToExternal(EvmNetVersionMsgType, nil)
}

// Query the number of peers currently connected to the client.
func (e EvmExternal) NetPeerCount() provider.XplaClient {
	return e.ToExternal(EvmNetPeerCountMsgType, nil)
}

// actively listening for network connections.
func (e EvmExternal) NetListening() provider.XplaClient {
	return e.ToExternal(EvmNetListeningMsgType, nil)
}

// ethereum protocol version.
func (e EvmExternal) EthProtocolVersion() provider.XplaClient {
	return e.ToExternal(EvmEthProtocolVersionMsgType, nil)
}

// Query the sync status object depending on the details of tendermint's sync protocol.
func (e EvmExternal) EthSyncing() provider.XplaClient {
	return e.ToExternal(EvmEthSyncingMsgType, nil)
}

// Query all eth accounts.
func (e EvmExternal) EthAccounts() provider.XplaClient {
	return e.ToExternal(EvmEthAccountsMsgType, nil)
}

// Query the number of transaction a given block.
func (e EvmExternal) EthGetBlockTransactionCount(ethGetBlockTransactionCountMsg types.EthGetBlockTransactionCountMsg) provider.XplaClient {
	if ethGetBlockTransactionCountMsg.BlockHash == "" && ethGetBlockTransactionCountMsg.BlockHeight == "" {
		return e.Err(EvmEthGetBlockTransactionCountMsgType, types.ErrWrap(types.ErrInsufficientParams, "cannot query, without block hash or height parameter"))
	}

	if ethGetBlockTransactionCountMsg.BlockHash != "" && ethGetBlockTransactionCountMsg.BlockHeight != "" {
		return e.Err(EvmEthGetBlockTransactionCountMsgType, types.ErrWrap(types.ErrInvalidRequest, "select only one parameter, block hash OR height"))
	}

	msg, err := MakeEthGetBlockTransactionCountMsg(ethGetBlockTransactionCountMsg)
	if err != nil {
		return e.Err(EvmEthGetBlockTransactionCountMsgType, err)
	}

	return e.ToExternal(EvmEthGetBlockTransactionCountMsgType, msg)
}

// Query estimate gas.
func (e EvmExternal) EstimateGas(invokeSolContractMsg types.InvokeSolContractMsg) provider.XplaClient {
	msg, err := MakeEstimateGasSolMsg(invokeSolContractMsg)
	if err != nil {
		return e.Err(EvmEthEstimateGasMsgType, err)
	}

	return e.ToExternal(EvmEthEstimateGasMsgType, msg)
}

// Query transaction by block hash and index.
func (e EvmExternal) EthGetTransactionByBlockHashAndIndex(getTransactionByBlockHashAndIndexMsg types.GetTransactionByBlockHashAndIndexMsg) provider.XplaClient {
	msg, err := MakeGetTransactionByBlockHashAndIndexMsg(getTransactionByBlockHashAndIndexMsg)
	if err != nil {
		return e.Err(EvmGetTransactionByBlockHashAndIndexMsgType, err)
	}

	return e.ToExternal(EvmGetTransactionByBlockHashAndIndexMsgType, msg)
}

// Query transaction receipt.
func (e EvmExternal) EthGetTransactionReceipt(getTransactionReceiptMsg types.GetTransactionReceiptMsg) provider.XplaClient {
	msg, err := MakeGetTransactionReceiptMsg(getTransactionReceiptMsg)
	if err != nil {
		return e.Err(EvmGetTransactionReceiptMsgType, err)
	}

	return e.ToExternal(EvmGetTransactionReceiptMsgType, msg)
}

// Query filter ID by eth new filter.
func (e EvmExternal) EthNewFilter(ethNewFilterMsg types.EthNewFilterMsg) provider.XplaClient {
	msg, err := MakeEthNewFilterMsg(ethNewFilterMsg)
	if err != nil {
		return e.Err(EvmEthNewFilterMsgType, err)
	}

	return e.ToExternal(EvmEthNewFilterMsgType, msg)
}

// Query filter ID by eth new block filter.
func (e EvmExternal) EthNewBlockFilter() provider.XplaClient {
	return e.ToExternal(EvmEthNewBlockFilterMsgType, nil)
}

// Query filter ID by eth new pending transaction filter.
func (e EvmExternal) EthNewPendingTransactionFilter() provider.XplaClient {
	return e.ToExternal(EvmEthNewPendingTransactionFilterMsgType, nil)
}

// Uninstall filter.
func (e EvmExternal) EthUninstallFilter(ethUninstallFilterMsg types.EthUninstallFilterMsg) provider.XplaClient {
	msg, err := MakeEthUninstallFilterMsg(ethUninstallFilterMsg)
	if err != nil {
		return e.Err(EvmEthUninstallFilterMsgType, err)
	}

	return e.ToExternal(EvmEthUninstallFilterMsgType, msg)
}

// Query filter changes.
func (e EvmExternal) EthGetFilterChanges(ethGetFilterChangesMsg types.EthGetFilterChangesMsg) provider.XplaClient {
	msg, err := MakeEthGetFilterChangesMsg(ethGetFilterChangesMsg)
	if err != nil {
		return e.Err(EvmEthGetFilterChangesMsgType, err)
	}

	return e.ToExternal(EvmEthGetFilterChangesMsgType, msg)
}

// Query filter logs.
func (e EvmExternal) EthGetFilterLogs(ethGetFilterLogsMsg types.EthGetFilterLogsMsg) provider.XplaClient {
	msg, err := MakeEthGetFilterLogsMsg(ethGetFilterLogsMsg)
	if err != nil {
		return e.Err(EvmEthGetFilterLogsMsgType, err)
	}

	return e.ToExternal(EvmEthGetFilterLogsMsgType, msg)
}

// Get logs.
func (e EvmExternal) EthGetLogs(ethGetLogsMsg types.EthGetLogsMsg) provider.XplaClient {
	msg, err := MakeEthGetLogsMsg(ethGetLogsMsg)
	if err != nil {
		return e.Err(EvmEthGetLogsMsgType, err)
	}

	return e.ToExternal(EvmEthGetLogsMsgType, msg)
}

// Query coinbase.
func (e EvmExternal) EthCoinbase() provider.XplaClient {
	return e.ToExternal(EvmEthCoinbaseMsgType, nil)
}
