package evm

import (
	"math/big"

	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

// Query client for evm module.
func QueryEvm(i core.QueryClient) (string, error) {
	evmClient, err := util.NewEvmClient(i.Ixplac.GetEvmRpc(), i.Ixplac.GetContext())
	if err != nil {
		return "", err
	}

	gasAdj := i.Ixplac.GetGasAdjustment()
	if i.Ixplac.GetGasAdjustment() == "" {
		gasAdj = types.DefaultGasAdjustment
	}

	gasLimit := i.Ixplac.GetGasLimit()
	if i.Ixplac.GetGasLimit() == "" {
		gasLimitU64, err := util.FromStringToUint64(util.DefaultEvmQueryGasLimit)
		if err != nil {
			return "", util.LogErr(errors.ErrParse, err)
		}
		gasLimitAdjustment, err := util.GasLimitAdjustment(gasLimitU64, gasAdj)
		if err != nil {
			return "", util.LogErr(errors.ErrParse, err)
		}
		gasLimit = gasLimitAdjustment
	}

	gasPrice := i.Ixplac.GetGasPrice()
	if i.Ixplac.GetGasPrice() == "" {
		gasPrice = types.DefaultGasPrice
	}

	gasPriceBigInt, err := util.FromStringToBigInt(gasPrice)
	if err != nil {
		return "", err
	}

	switch {
	// Evm call contract
	case i.Ixplac.GetMsgType() == EvmCallSolContractMsgType:
		convertMsg := i.Ixplac.GetMsg().(CallSolContractParseMsg)

		gasLimitU64, err := util.FromStringToUint64(gasLimit)
		if err != nil {
			return "", util.LogErr(errors.ErrParse, err)
		}
		convertMsg.CallMsg.Gas = gasLimitU64
		convertMsg.CallMsg.GasPrice = gasPriceBigInt

		res, err := evmClient.Client.CallContract(evmClient.Ctx, convertMsg.CallMsg, nil)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		result, err := util.GetAbiUnpack(convertMsg.CallName, convertMsg.ABI, convertMsg.Bytecode, res)
		if err != nil {
			return "", util.LogErr(errors.ErrParse, err)
		}

		var callSolContractResponse types.CallSolContractResponse
		for _, res := range result {
			callSolContractResponse.ContractResponse = append(callSolContractResponse.ContractResponse, util.ToString(res, ""))
		}

		return jsonReturn(callSolContractResponse)

	// Evm transaction by hash
	case i.Ixplac.GetMsgType() == EvmGetTransactionByHashMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.GetTransactionByHashMsg)
		commonTxHash := util.FromStringHexToHash(convertMsg.TxHash)
		tx, isPending, err := evmClient.Client.TransactionByHash(evmClient.Ctx, commonTxHash)
		if isPending {
			return "", util.LogErr(errors.ErrNotFound, "tx is pending..")
		}
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		return jsonReturn(tx)

	// Evm block by hash or height
	case i.Ixplac.GetMsgType() == EvmGetBlockByHashHeightMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.GetBlockByHashHeightMsg)
		var block *ethtypes.Block
		var blockResponse types.BlockResponse

		if convertMsg.BlockHash != "" {
			commonBlockHash := util.FromStringHexToHash(convertMsg.BlockHash)
			block, err = evmClient.Client.BlockByHash(evmClient.Ctx, commonBlockHash)
			if err != nil {
				return "", util.LogErr(errors.ErrEvmRpcRequest, err)
			}
		} else {
			blockNumber, err := util.FromStringToBigInt(convertMsg.BlockHeight)
			if err != nil {
				return "", err
			}

			block, err = evmClient.Client.BlockByNumber(evmClient.Ctx, blockNumber)
			if err != nil {
				return "", util.LogErr(errors.ErrEvmRpcRequest, err)
			}
		}

		txs := block.Body().Transactions
		uncles := block.Body().Uncles

		blockResponse.BlockHeader = block.Header()
		blockResponse.Transactions = txs
		blockResponse.Uncles = uncles

		return jsonReturn(blockResponse)

	// Evm account information
	case i.Ixplac.GetMsgType() == EvmQueryAccountInfoMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.AccountInfoMsg)
		account := util.FromStringToByte20Address(convertMsg.Account)

		balance, err := evmClient.Client.BalanceAt(evmClient.Ctx, account, nil)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}
		currentNonce, err := evmClient.Client.NonceAt(evmClient.Ctx, account, nil)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}
		storage, err := evmClient.Client.StorageAt(evmClient.Ctx, account, util.FromStringHexToHash("0"), nil)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}
		code, err := evmClient.Client.CodeAt(evmClient.Ctx, account, nil)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}
		pendingBalance, err := evmClient.Client.PendingBalanceAt(evmClient.Ctx, account)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}
		pendingNonce, err := evmClient.Client.PendingNonceAt(evmClient.Ctx, account)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}
		pendingStorage, err := evmClient.Client.PendingStorageAt(evmClient.Ctx, account, util.FromStringHexToHash("0"))
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}
		pendingCode, err := evmClient.Client.PendingCodeAt(evmClient.Ctx, account)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}
		pendingTransactionCount, err := evmClient.Client.PendingTransactionCount(evmClient.Ctx)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		bech32Addr, err := util.FromByte20AddressToCosmosAddr(account)
		if err != nil {
			return "", util.LogErr(errors.ErrParse, err)
		}

		var accountInfoResponse types.AccountInfoResponse

		accountInfoResponse.Account = account.Hex()
		accountInfoResponse.Bech32Account = bech32Addr.String()
		accountInfoResponse.Balance = balance
		accountInfoResponse.Nonce = currentNonce
		accountInfoResponse.Storage = string(storage)
		accountInfoResponse.Code = string(code)
		accountInfoResponse.PendingBalance = pendingBalance
		accountInfoResponse.PendingNonce = pendingNonce
		accountInfoResponse.PendingStorage = string(pendingStorage)
		accountInfoResponse.PendingCode = string(pendingCode)
		accountInfoResponse.PendingTransactionCount = pendingTransactionCount

		return jsonReturn(accountInfoResponse)

	// Evm suggest gas price
	case i.Ixplac.GetMsgType() == EvmSuggestGasPriceMsgType:
		gasPrice, err := evmClient.Client.SuggestGasPrice(evmClient.Ctx)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		gasTipCap, err := evmClient.Client.SuggestGasTipCap(evmClient.Ctx)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var suggestGasPriceResponse types.SuggestGasPriceResponse
		suggestGasPriceResponse.GasPrice = gasPrice
		suggestGasPriceResponse.GasTipCap = gasTipCap

		return jsonReturn(suggestGasPriceResponse)

	// Evm chain ID
	case i.Ixplac.GetMsgType() == EvmQueryChainIdMsgType:
		chainId, err := evmClient.Client.ChainID(evmClient.Ctx)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var ethChainIdResponse types.EthChainIdResponse
		ethChainIdResponse.ChainID = chainId

		return jsonReturn(ethChainIdResponse)

	// Evm latest block height
	case i.Ixplac.GetMsgType() == EvmQueryCurrentBlockNumberMsgType:
		blockNumber, err := evmClient.Client.BlockNumber(evmClient.Ctx)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var ethBlockNumberResponse types.EthBlockNumberResponse
		ethBlockNumberResponse.BlockNumber = blockNumber

		return jsonReturn(ethBlockNumberResponse)

	// Web3 client version
	case i.Ixplac.GetMsgType() == EvmWeb3ClientVersionMsgType:
		var result string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "web3_clientVersion")
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var web3ClientVersionResponse types.Web3ClientVersionResponse
		web3ClientVersionResponse.Web3ClientVersion = result

		return jsonReturn(web3ClientVersionResponse)

	// Web3 sha
	case i.Ixplac.GetMsgType() == EvmWeb3Sha3MsgType:
		convertMsg := i.Ixplac.GetMsg().(types.Web3Sha3Msg)

		var result string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "web3_sha3", convertMsg.InputParam)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var web3Sha3Response types.Web3Sha3Response
		web3Sha3Response.Web3Sha3 = result

		return jsonReturn(web3Sha3Response)

	// network ID
	case i.Ixplac.GetMsgType() == EvmNetVersionMsgType:
		var result string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "net_version")
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var netVersionResponse types.NetVersionResponse
		netVersionResponse.NetVersion = result

		return jsonReturn(netVersionResponse)

	// the number of peers
	case i.Ixplac.GetMsgType() == EvmNetPeerCountMsgType:
		var result int
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "net_peerCount")
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var netPeerCountResponse types.NetPeerCountResponse
		netPeerCountResponse.NetPeerCount = result

		return jsonReturn(netPeerCountResponse)

	// actively listening for network connections
	case i.Ixplac.GetMsgType() == EvmNetListeningMsgType:
		var result bool
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "net_listening")
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var netListeningResponse types.NetListeningResponse
		netListeningResponse.NetListening = result

		return jsonReturn(netListeningResponse)

	// eth protocol version
	case i.Ixplac.GetMsgType() == EvmEthProtocolVersionMsgType:
		resultBigInt := big.NewInt(0)

		var result string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_protocolVersion")
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		if result != "" {
			resultBigInt = util.From0xHexStringToIBignt(result)
		}

		var ethProtocolVersionResponse types.EthProtocolVersionResponse
		ethProtocolVersionResponse.EthProtocolVersionHex = result
		ethProtocolVersionResponse.EthProtocolVersion = resultBigInt

		return jsonReturn(ethProtocolVersionResponse)

	// eth syncing status
	case i.Ixplac.GetMsgType() == EvmEthSyncingMsgType:
		var result bool
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_syncing")
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var ethSyncingResponse types.EthSyncingResponse
		ethSyncingResponse.EthSyncing = result

		return jsonReturn(ethSyncingResponse)

	// eth all accounts
	case i.Ixplac.GetMsgType() == EvmEthAccountsMsgType:
		var result []string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_accounts")
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var ethAccountsResponse types.EthAccountsResponse
		ethAccountsResponse.EthAccounts = result

		return jsonReturn(ethAccountsResponse)

	// the number of transaction a given block
	case i.Ixplac.GetMsgType() == EvmEthGetBlockTransactionCountMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.EthGetBlockTransactionCountMsg)
		resultBigInt := big.NewInt(0)

		var result string
		if convertMsg.BlockHash != "" {
			err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_getBlockTransactionCountByHash", convertMsg.BlockHash)
			if err != nil {
				return "", util.LogErr(errors.ErrEvmRpcRequest, err)
			}
		}

		if convertMsg.BlockHeight != "" {
			err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_getBlockTransactionCountByNumber", convertMsg.BlockHeight)
			if err != nil {
				return "", util.LogErr(errors.ErrEvmRpcRequest, err)
			}
		}

		if result != "" {
			resultBigInt = util.From0xHexStringToIBignt(result)
		} else {
			result = "not found"
		}

		var ethGetBlockTransactionCountResponse types.EthGetBlockTransactionCountResponse
		ethGetBlockTransactionCountResponse.EthGetBlockTransactionCountHex = result
		ethGetBlockTransactionCountResponse.EthGetBlockTransactionCount = resultBigInt

		return jsonReturn(ethGetBlockTransactionCountResponse)

	// Evm call contract
	case i.Ixplac.GetMsgType() == EvmEthEstimateGasMsgType:
		convertMsg := i.Ixplac.GetMsg().(CallSolContractParseMsg)
		convertMsg.CallMsg.Gas = 0
		convertMsg.CallMsg.GasPrice = gasPriceBigInt

		res, err := evmClient.Client.EstimateGas(evmClient.Ctx, convertMsg.CallMsg)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var estimateGasResponse types.EstimateGasResponse
		estimateGasResponse.EstimateGas = res

		return jsonReturn(estimateGasResponse)

	// get transaction by block hash and index
	case i.Ixplac.GetMsgType() == EvmGetTransactionByBlockHashAndIndexMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.GetTransactionByBlockHashAndIndexMsg)

		blockHash := util.FromStringHexToHash(convertMsg.BlockHash)

		indexU64, err := util.FromStringToUint64(convertMsg.Index)
		if err != nil {
			return "", util.LogErr(errors.ErrParse, err)
		}
		index := indexU64

		res, err := evmClient.Client.TransactionInBlock(evmClient.Ctx, blockHash, uint(index))
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		return jsonReturn(res)

	// get transaction receipt
	case i.Ixplac.GetMsgType() == EvmGetTransactionReceiptMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.GetTransactionReceiptMsg)

		transactionHash := util.FromStringHexToHash(convertMsg.TransactionHash)

		res, err := evmClient.Client.TransactionReceipt(evmClient.Ctx, transactionHash)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		return jsonReturn(res)

	// get filter ID by eth new filter
	case i.Ixplac.GetMsgType() == EvmEthNewFilterMsgType:
		convertMsg := i.Ixplac.GetMsg().(EthNewFilterParseMsg)

		var result interface{}
		err = evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_newFilter", convertMsg)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		var ethNewFilterResponse types.EthNewFilterResponse
		ethNewFilterResponse.NewFilter = result

		return jsonReturn(ethNewFilterResponse)

	// get transaction receipt
	case i.Ixplac.GetMsgType() == EvmEthNewBlockFilterMsgType:

		var result interface{}
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_newBlockFilter")
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		ethNewBlockFilterResponse := types.EthNewBlockFilterResponse{
			NewBlockFilter: result,
		}

		return jsonReturn(ethNewBlockFilterResponse)

	// get transaction receipt
	case i.Ixplac.GetMsgType() == EvmEthNewPendingTransactionFilterMsgType:

		var result interface{}
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_newPendingTransactionFilter")
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		ethNewPendingTransactionFilterResponse := types.EthNewPendingTransactionFilterResponse{
			NewPendingTransactionFilter: result,
		}

		return jsonReturn(ethNewPendingTransactionFilterResponse)

	// uninstall filter
	case i.Ixplac.GetMsgType() == EvmEthUninstallFilterMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.EthUninstallFilterMsg)

		var result bool
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_uninstallFilter", convertMsg.FilterId)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		ethUninstallFilterResponse := types.EthUninstallFilterResponse{
			UninstallFilter: result,
		}

		return jsonReturn(ethUninstallFilterResponse)

	// get filter changes
	case i.Ixplac.GetMsgType() == EvmEthGetFilterChangesMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.EthGetFilterChangesMsg)

		var result []string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_getFilterChanges", convertMsg.FilterId)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		ethGetFilterChangesResponse := types.EthGetFilterChangesResponse{
			GetFilterChanges: result,
		}

		return jsonReturn(ethGetFilterChangesResponse)

	// get filter logs
	case i.Ixplac.GetMsgType() == EvmEthGetFilterLogsMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.EthGetFilterLogsMsg)

		var result []string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_getFilterLogs", convertMsg.FilterId)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		ethGetFilterLogsResponse := types.EthGetFilterLogsResponse{
			GetFilterLogs: result,
		}

		return jsonReturn(ethGetFilterLogsResponse)

	// get logs
	case i.Ixplac.GetMsgType() == EvmEthGetLogsMsgType:
		convertMsg := i.Ixplac.GetMsg().(EthNewFilterParseMsg)

		var result interface{}
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_getLogs", convertMsg)
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		ethGetLogsResponse := types.EthGetLogsResponse{
			GetLogs: result,
		}

		return jsonReturn(ethGetLogsResponse)

	// get coinbase
	case i.Ixplac.GetMsgType() == EvmEthCoinbaseMsgType:

		var result string
		err := evmClient.RpcClient.CallContext(evmClient.Ctx, &result, "eth_coinbase")
		if err != nil {
			return "", util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		ethCoinbaseResponse := types.EthCoinbaseResponse{
			Coinbase: result,
		}

		return jsonReturn(ethCoinbaseResponse)

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}
}

func jsonReturn(value interface{}) (string, error) {
	json, err := util.JsonMarshalDataIndent(value)
	if err != nil {
		return "", util.LogErr(errors.ErrFailedToMarshal, err)
	}

	return string(json), nil
}
