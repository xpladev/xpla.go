package client

import (
	"encoding/json"
	"time"

	mevm "github.com/xpladev/xpla.go/core/evm"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	txtypes "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
)

// Broadcast generated transactions.
// Broadcast responses, excluding evm, are delivered as "TxResponse" of the entire response structure of the xpla client.
// Support broadcast by using LCD and gRPC at the same time. Default method is gRPC.
func broadcastTx(xplac *xplaClient, txBytes []byte, mode txtypes.BroadcastMode) (*types.TxRes, error) {
	broadcastReq := txtypes.BroadcastTxRequest{
		TxBytes: txBytes,
		Mode:    mode,
	}

	if xplac.GetGrpcUrl() == "" {
		reqBytes, err := json.Marshal(broadcastReq)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrFailedToMarshal, err))
		}

		xplac.GetHttpMutex().Lock()
		out, err := util.CtxHttpClient("POST", xplac.GetLcdURL()+broadcastUrl, reqBytes, xplac.GetContext())
		if err != nil {
			xplac.GetHttpMutex().Unlock()
			return nil, xplac.GetLogger().Err(err)
		}
		xplac.GetHttpMutex().Unlock()

		var broadcastTxResponse txtypes.BroadcastTxResponse
		err = xplac.GetEncoding().Codec.UnmarshalJSON(out, &broadcastTxResponse)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrFailedToUnmarshal, err))
		}

		txResponse := broadcastTxResponse.TxResponse
		if txResponse.Code != 0 {
			return &xplaTxRes, xplac.GetLogger().Err(types.ErrWrap(types.ErrTxFailed, "with code", txResponse.Code, ":", txResponse.RawLog))
		}

		xplaTxRes.Response = txResponse
	} else {
		txClient := txtypes.NewServiceClient(xplac.GetGrpcClient())
		txResponse, err := txClient.BroadcastTx(xplac.GetContext(), &broadcastReq)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrGrpcRequest, err))
		}
		xplaTxRes.Response = txResponse.TxResponse
	}

	return &xplaTxRes, nil
}

// Broadcast generated transactions of ethereum type.
// Broadcast responses, including evm, are delivered as "TxResponse".
func broadcastTxEvm(xplac *xplaClient, txBytes []byte, broadcastMode string, evmClient *util.EvmClient) (*types.TxRes, error) {
	switch {
	case xplac.GetMsgType() == mevm.EvmSendCoinMsgType ||
		xplac.GetMsgType() == mevm.EvmInvokeSolContractMsgType:
		var signedTx evmtypes.Transaction
		err := signedTx.UnmarshalJSON(txBytes)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrFailedToUnmarshal, err))
		}

		err = evmClient.Client.SendTransaction(evmClient.Ctx, &signedTx)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrEvmRpcRequest, err))
		}

		return checkEvmBroadcastMode(xplac, broadcastMode, evmClient, &signedTx)

	case xplac.GetMsgType() == mevm.EvmDeploySolContractMsgType:
		var deployTx mevm.DeploySolTx

		err := json.Unmarshal(txBytes, &deployTx)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrFailedToUnmarshal, err))
		}

		ethPrivKey, err := toECDSA(xplac.GetPrivateKey())
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrParse, err))
		}

		contractAuth, err := bind.NewKeyedTransactorWithChainID(ethPrivKey, deployTx.ChainId)
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrInsufficientParams, err))
		}
		contractAuth.Nonce = deployTx.Nonce
		contractAuth.Value = deployTx.Value
		contractAuth.GasLimit = deployTx.GasLimit
		contractAuth.GasPrice = deployTx.GasPrice

		metadata := util.GetBindMetaData(deployTx.ABI, deployTx.Bytecode)
		parsedAbi, err := metadata.GetAbi()
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrParse, err))
		}
		if parsedAbi == nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrParse, "GetABI returned nil"))
		}
		parsedBytecode := common.FromHex(metadata.Bin)

		var transaction *evmtypes.Transaction
		if mevm.Args == nil {
			_, transaction, _, err = bind.DeployContract(contractAuth, *parsedAbi, parsedBytecode, evmClient.Client)
		} else {
			_, transaction, _, err = bind.DeployContract(contractAuth, *parsedAbi, parsedBytecode, evmClient.Client, mevm.Args...)
			mevm.Args = nil
		}
		if err != nil {
			return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrEvmRpcRequest, err))
		}

		return checkEvmBroadcastMode(xplac, broadcastMode, evmClient, transaction)

	default:
		return nil, xplac.GetLogger().Err(types.ErrWrap(types.ErrInvalidMsgType, "invalid EVM msg type:", xplac.GetMsgType()))
	}
}

// Handle evm broadcast mode.
// Similarly, determine broadcast mode included in the options of xpla client.
func checkEvmBroadcastMode(xplac *xplaClient, broadcastMode string, evmClient *util.EvmClient, tx *evmtypes.Transaction) (*types.TxRes, error) {
	// Wait tx receipt (Broadcast Block)
	if broadcastMode == "block" {
		receipt, err := waitTxReceipt(evmClient, tx)
		if err != nil {
			return nil, xplac.GetLogger().Err(err)
		}
		xplaTxRes.EvmReceipt = receipt
		return &xplaTxRes, nil
	} else {
		return nil, nil
	}
}

// If broadcast mode is "block", client waits transaction receipt of evm.
// The wait time is in seconds and is set inside the library as timeout count.
// When the timeout occurs, no longer wait for the transaction receipt.
func waitTxReceipt(evmClient *util.EvmClient, signedTx *evmtypes.Transaction) (*evmtypes.Receipt, error) {
	count := util.DefaultEvmTxReceiptTimeout
	for {
		receipt, err := evmClient.Client.TransactionReceipt(evmClient.Ctx, signedTx.Hash())
		if err != nil {
			count = count - 1
			if count < 0 {
				return nil, types.ErrWrap(types.ErrEvmRpcRequest, "cannot receive the transaction receipt in count time is", count)
			}
			time.Sleep(time.Second * 1)
		} else {
			return receipt, nil
		}
	}
}
