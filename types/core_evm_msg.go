package types

import (
	"math/big"

	ethtypes "github.com/ethereum/go-ethereum/core/types"
)

type SendCoinMsg struct {
	Amount      string
	FromAddress string
	ToAddress   string
}

type DeploySolContractMsg struct {
	ABI                  string
	Bytecode             string
	ABIJsonFilePath      string
	BytecodeJsonFilePath string
	Args                 []interface{}
}

type InvokeSolContractMsg struct {
	ContractAddress      string
	ContractFuncCallName string
	Args                 []interface{}
	ABI                  string
	Bytecode             string
	ABIJsonFilePath      string
	BytecodeJsonFilePath string
	FromByteAddress      string
}

type CallSolContractMsg struct {
	ContractAddress      string
	ContractFuncCallName string
	Args                 []interface{}
	ABI                  string
	Bytecode             string
	ABIJsonFilePath      string
	BytecodeJsonFilePath string
	FromByteAddress      string
}

type GetTransactionByHashMsg struct {
	TxHash string
}

type GetBlockByHashHeightMsg struct {
	BlockHash   string
	BlockHeight string
}

type AccountInfoMsg struct {
	Account string
}

type Web3Sha3Msg struct {
	InputParam string
}

type EthGetBlockTransactionCountMsg struct {
	BlockHash   string
	BlockHeight string
}

type GetTransactionByBlockHashAndIndexMsg struct {
	BlockHash string
	Index     string
}

type GetTransactionReceiptMsg struct {
	TransactionHash string
}

type EthNewFilterMsg struct {
	FromBlock string
	ToBlock   string
	Address   []string
	Topics    []string
}

type EthGetLogsMsg struct {
	BlockHash string
	FromBlock string
	ToBlock   string
	Address   []string
	Topics    []string
}

type EthUninstallFilterMsg struct {
	FilterId string
}

type EthGetFilterChangesMsg struct {
	FilterId string
}

type EthGetFilterLogsMsg struct {
	FilterId string
}

// Responses
type CallSolContractResponse struct {
	ContractResponse []string `json:"contract_response"`
}

type BlockResponse struct {
	BlockHeader  *ethtypes.Header      `json:"blockHeader"`
	Transactions ethtypes.Transactions `json:"transactions"`
	Uncles       []*ethtypes.Header    `json:"uncles"`
}

type AccountInfoResponse struct {
	Account                 string   `json:"account"`
	Bech32Account           string   `json:"bech32_account"`
	Balance                 *big.Int `json:"balance"`
	Nonce                   uint64   `json:"nonce"`
	Storage                 string   `json:"storage"`
	Code                    string   `json:"code"`
	PendingBalance          *big.Int `json:"pending_balance"`
	PendingNonce            uint64   `json:"pending_nonce"`
	PendingStorage          string   `json:"pending_storage"`
	PendingCode             string   `json:"pending_code"`
	PendingTransactionCount uint     `json:"pending_transaction_count"`
}

type SuggestGasPriceResponse struct {
	GasPrice  *big.Int `json:"gas_price"`
	GasTipCap *big.Int `json:"gas_tip_cap"`
}

type EthChainIdResponse struct {
	ChainID *big.Int `json:"chain_id"`
}

type EthBlockNumberResponse struct {
	BlockNumber uint64 `json:"block_number"`
}

type Web3ClientVersionResponse struct {
	Web3ClientVersion string `json:"web3_clientVersion"`
}

type Web3Sha3Response struct {
	Web3Sha3 string `json:"web3_sha3"`
}

type NetVersionResponse struct {
	NetVersion string `json:"net_version"`
}

type NetPeerCountResponse struct {
	NetPeerCount int `json:"net_peerCount"`
}

type NetListeningResponse struct {
	NetListening bool `json:"net_listening"`
}

type EthProtocolVersionResponse struct {
	EthProtocolVersionHex string   `json:"eth_protocolVersion_hex"`
	EthProtocolVersion    *big.Int `json:"eth_protocolVersion"`
}

type EthSyncingResponse struct {
	EthSyncing bool `json:"eth_syncing"`
}

type EthAccountsResponse struct {
	EthAccounts []string `json:"eth_accounts"`
}

type EthGetBlockTransactionCountResponse struct {
	EthGetBlockTransactionCountHex string   `json:"eth_getBlockTransactionCount_hex"`
	EthGetBlockTransactionCount    *big.Int `json:"eth_getBlockTransactionCount"`
}

type EstimateGasResponse struct {
	EstimateGas uint64 `json:"eth_estimateGas"`
}

type EthNewFilterResponse struct {
	NewFilter interface{} `json:"eth_newFilter"`
}

type EthNewBlockFilterResponse struct {
	NewBlockFilter interface{} `json:"eth_newBlockFilter"`
}

type EthNewPendingTransactionFilterResponse struct {
	NewPendingTransactionFilter interface{} `json:"eth_newPendingTransactionFilter"`
}

type EthUninstallFilterResponse struct {
	UninstallFilter bool `json:"eth_uninstallFilter"`
}

type EthGetFilterChangesResponse struct {
	GetFilterChanges []string `json:"eth_getFilterChanges"`
}

type EthGetFilterLogsResponse struct {
	GetFilterLogs []string `json:"eth_getFilterLogs"`
}

type EthGetLogsResponse struct {
	GetLogs interface{} `json:"eth_getLogs"`
}

type EthCoinbaseResponse struct {
	Coinbase string `json:"eth_coinbase"`
}
