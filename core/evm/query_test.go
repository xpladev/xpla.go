package evm_test

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"math/rand"
	"strings"
	"testing"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"

	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/suite"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"golang.org/x/crypto/sha3"
)

var (
	validatorNumber          = 1
	testBalance              = "1000000000000000000000000000"
	testABIJsonFilePath      = "../../util/testutil/test_files/abi.json"
	testBytecodeJsonFilePath = "../../util/testutil/test_files/bytecode.json"
)

type IntegrationTestSuite struct {
	suite.Suite

	xplac              provider.XplaClient
	accounts           []simtypes.Account
	evmtestTxHash      string
	evmTestBlockHeight int64
	contractAddr       string

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	src := rand.NewSource(1)
	r := rand.New(src)
	s.accounts = testutil.RandomAccounts(r, 2)

	balanceBigInt, err := util.FromStringToBigInt(testBalance)
	s.Require().NoError(err)

	genesisState := s.cfg.GenesisState

	// add genesis account
	var authGenesis authtypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[authtypes.ModuleName], &authGenesis))

	var genAccounts []authtypes.GenesisAccount

	genAccounts = append(genAccounts, &ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(s.accounts[0].Address, nil, 0, 0),
		CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
	})
	genAccounts = append(genAccounts, &ethermint.EthAccount{
		BaseAccount: authtypes.NewBaseAccount(s.accounts[1].Address, nil, 0, 0),
		CodeHash:    common.BytesToHash(evmtypes.EmptyCodeHash).Hex(),
	})

	accounts, err := authtypes.PackAccounts(genAccounts)
	s.Require().NoError(err)

	authGenesis.Accounts = accounts

	authGenesisBz, err := s.cfg.Codec.MarshalJSON(&authGenesis)
	s.Require().NoError(err)
	genesisState[authtypes.ModuleName] = authGenesisBz

	// add balances
	var bankGenesis banktypes.GenesisState
	s.Require().NoError(s.cfg.Codec.UnmarshalJSON(genesisState[banktypes.ModuleName], &bankGenesis))

	bankGenesis.Balances = []banktypes.Balance{
		{
			Address: s.accounts[0].Address.String(),
			Coins: sdk.Coins{
				sdk.NewCoin(types.XplaDenom, sdk.NewIntFromBigInt(balanceBigInt)),
			},
		},
		{
			Address: s.accounts[1].Address.String(),
			Coins: sdk.Coins{
				sdk.NewCoin(types.XplaDenom, sdk.NewIntFromBigInt(balanceBigInt)),
			},
		},
	}

	bankGenesisBz, err := s.cfg.Codec.MarshalJSON(&bankGenesis)
	s.Require().NoError(err)
	genesisState[banktypes.ModuleName] = bankGenesisBz

	s.cfg.GenesisState = genesisState
	s.network = network.New(s.T(), s.cfg)
	s.Require().NoError(s.network.WaitForNextBlock())

	s.xplac = client.NewXplaClient(testutil.TestChainId).
		WithEvmRpc("http://" + s.network.Validators[0].AppConfig.JSONRPC.Address).
		WithURL(s.network.Validators[0].APIAddress)

	xplac := s.xplac.WithPrivateKey(s.accounts[0].PrivKey)

	// deploy contract
	deploySolContractMsg := types.DeploySolContractMsg{
		ABIJsonFilePath:      testABIJsonFilePath,
		BytecodeJsonFilePath: testBytecodeJsonFilePath,
		Args:                 nil,
	}
	txbytes, err := xplac.DeploySolidityContract(deploySolContractMsg).CreateAndSignTx()
	s.Require().NoError(err)

	_, err = xplac.Broadcast(txbytes)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// request tx event
	txEventsRes, err := xplac.TxsByEvents(types.QueryTxsByEventsMsg{
		Events: "transfer.sender=" + s.accounts[0].Address.String(),
	}).Query()
	s.Require().NoError(err)

	var getTxsEventResponse sdktx.GetTxsEventResponse
	jsonpb.Unmarshal(strings.NewReader(txEventsRes), &getTxsEventResponse)

	// extract evm transaction
	s.evmTestBlockHeight = getTxsEventResponse.TxResponses[0].Height
	s.evmtestTxHash = getTxsEventResponse.TxResponses[0].Logs[0].Events[0].Attributes[1].Value

	// query transaction receipte by evm tx hash
	getTransactionReceiptMsg := types.GetTransactionReceiptMsg{
		TransactionHash: s.evmtestTxHash,
	}
	receiptRes, err := xplac.EthGetTransactionReceipt(getTransactionReceiptMsg).Query()
	s.Require().NoError(err)

	var receipt ethtypes.Receipt
	json.Unmarshal([]byte(receiptRes), &receipt)

	s.contractAddr = receipt.ContractAddress.String()
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestCallSolidityContract() {
	xplac := s.xplac.WithPrivateKey(s.accounts[0].PrivKey)

	// call contract
	callSolContractMsg := types.CallSolContractMsg{
		ContractAddress:      s.contractAddr,
		ContractFuncCallName: "retrieve",
		ABIJsonFilePath:      testABIJsonFilePath,
		BytecodeJsonFilePath: testBytecodeJsonFilePath,
		FromByteAddress:      util.FromStringToByte20Address(s.accounts[0].PubKey.Address().String()).String(),
	}

	res, err := xplac.CallSolidityContract(callSolContractMsg).Query()
	s.Require().NoError(err)

	var callSolContractResponse types.CallSolContractResponse
	json.Unmarshal([]byte(res), &callSolContractResponse)

	s.Require().Equal("0", callSolContractResponse.ContractResponse[0])

	// update contract
	newStoreValue := "123"
	newStoreValueBigInt, err := util.FromStringToBigInt(newStoreValue)
	s.Require().NoError(err)

	var args []interface{}
	args = append(args, newStoreValueBigInt)

	invokeSolContractMsg := types.InvokeSolContractMsg{
		ContractAddress:      s.contractAddr,
		ContractFuncCallName: "store",
		Args:                 args,
		ABIJsonFilePath:      testABIJsonFilePath,
		BytecodeJsonFilePath: testBytecodeJsonFilePath,
		FromByteAddress:      util.FromStringToByte20Address(s.accounts[0].PubKey.Address().String()).String(),
	}

	exeTxbytes, err := xplac.WithSequence("").InvokeSolidityContract(invokeSolContractMsg).CreateAndSignTx()
	s.Require().NoError(err)

	_, err = xplac.Broadcast(exeTxbytes)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// call contract repeat
	callSolContractMsg = types.CallSolContractMsg{
		ContractAddress:      s.contractAddr,
		ContractFuncCallName: "retrieve",
		ABIJsonFilePath:      testABIJsonFilePath,
		BytecodeJsonFilePath: testBytecodeJsonFilePath,
		FromByteAddress:      util.FromStringToByte20Address(s.accounts[0].PubKey.Address().String()).String(),
	}

	res, err = xplac.CallSolidityContract(callSolContractMsg).Query()
	s.Require().NoError(err)

	json.Unmarshal([]byte(res), &callSolContractResponse)

	s.Require().Equal(newStoreValue, callSolContractResponse.ContractResponse[0])
}

func (s *IntegrationTestSuite) TestGetTransactionByHash() {
	// get tx which is deployed contract
	getTransactionByHashMsg := types.GetTransactionByHashMsg{
		TxHash: s.evmtestTxHash,
	}
	res, err := s.xplac.GetTransactionByHash(getTransactionByHashMsg).Query()
	s.Require().NoError(err)

	var transaction ethtypes.Transaction
	json.Unmarshal([]byte(res), &transaction)

	// unmarshal contract bytecode file
	bytecode, err := util.BytecodeParsing(testBytecodeJsonFilePath)
	s.Require().NoError(err)
	s.Require().Equal(bytecode, hex.EncodeToString(transaction.Data()))
}

func (s *IntegrationTestSuite) TestGetBlockByHashOrHeight() {
	// get block for test
	blockMsg := types.BlockMsg{
		Height: "1",
	}
	blockRes, err := s.xplac.Block(blockMsg).Query()
	s.Require().NoError(err)

	var resultBlock ctypes.ResultBlock
	json.Unmarshal([]byte(blockRes), &resultBlock)

	blockHash := resultBlock.BlockID.Hash.String()

	// get block by hash
	getBlockByHashHeightMsg := types.GetBlockByHashHeightMsg{
		BlockHash: blockHash,
	}
	hashRes, err := s.xplac.GetBlockByHashOrHeight(getBlockByHashHeightMsg).Query()
	s.Require().NoError(err)

	// get block by number
	getBlockByHashHeightMsg = types.GetBlockByHashHeightMsg{
		BlockHeight: "1",
	}
	numberRes, err := s.xplac.GetBlockByHashOrHeight(getBlockByHashHeightMsg).Query()
	s.Require().NoError(err)

	s.Require().Equal(hashRes, numberRes)
}

func (s *IntegrationTestSuite) TestAccountInfo() {
	accountInfoMsg := types.AccountInfoMsg{
		Account: s.accounts[0].PubKey.Address().String(),
	}
	res, err := s.xplac.AccountInfo(accountInfoMsg).Query()
	s.Require().NoError(err)

	var accountInfoResponse types.AccountInfoResponse
	json.Unmarshal([]byte(res), &accountInfoResponse)

	s.Require().Equal(s.accounts[0].PubKey.Address().String(), strings.ToUpper(accountInfoResponse.Account[2:]))
	s.Require().Equal(s.accounts[0].Address.String(), accountInfoResponse.Bech32Account)
}

func (s *IntegrationTestSuite) TestSuggestGasPrice() {
	res, err := s.xplac.SuggestGasPrice().Query()
	s.Require().NoError(err)

	var queryParamsResponse types.SuggestGasPriceResponse
	json.Unmarshal([]byte(res), &queryParamsResponse)

	s.Require().NotEqual(big.NewInt(0), queryParamsResponse.GasPrice)
	s.Require().NotEqual(big.NewInt(0), queryParamsResponse.GasTipCap)
}

func (s *IntegrationTestSuite) TestEthChainID() {
	res, err := s.xplac.EthChainID().Query()
	s.Require().NoError(err)

	var ethChainIdResponse types.EthChainIdResponse
	json.Unmarshal([]byte(res), &ethChainIdResponse)

	s.Require().Equal(big.NewInt(47), ethChainIdResponse.ChainID)
}

func (s *IntegrationTestSuite) TestEthBlockNumber() {
	// get latest block height by cosmos LCD
	blockRes, err := s.xplac.Block().Query()
	s.Require().NoError(err)

	type blockHeightStruct struct {
		Block struct {
			Header struct {
				Height string `json:"height"`
			} `json:"header"`
		} `json:"block"`
	}
	var b blockHeightStruct
	json.Unmarshal([]byte(blockRes), &b)

	// get latest block number by eth client
	res, err := s.xplac.EthBlockNumber().Query()
	s.Require().NoError(err)

	var ethBlockNumberResponse types.EthBlockNumberResponse
	json.Unmarshal([]byte(res), &ethBlockNumberResponse)

	heightU64, err := util.FromStringToUint64(b.Block.Header.Height)
	s.Require().NoError(err)
	s.Require().Equal(heightU64, ethBlockNumberResponse.BlockNumber)
}

func (s *IntegrationTestSuite) TestWeb3ClientVersion() {
	res, err := s.xplac.Web3ClientVersion().Query()
	s.Require().NoError(err)

	var web3ClientVersionResponse types.Web3ClientVersionResponse
	json.Unmarshal([]byte(res), &web3ClientVersionResponse)

	s.Require().NotEqual("", web3ClientVersionResponse.Web3ClientVersion)
}

func (s *IntegrationTestSuite) TestWeb3Sha3() {
	// query sha3
	testValue := "testweb3Sha3"
	web3Sha3Msg := types.Web3Sha3Msg{
		InputParam: testValue,
	}
	res, err := s.xplac.Web3Sha3(web3Sha3Msg).Query()
	s.Require().NoError(err)

	var web3Sha3Response types.Web3Sha3Response
	json.Unmarshal([]byte(res), &web3Sha3Response)

	// make hash by using keccak256
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write([]byte(testValue))

	hashValue := hex.EncodeToString(hasher.Sum(nil))

	s.Require().Equal(hashValue, web3Sha3Response.Web3Sha3[2:])
}

func (s *IntegrationTestSuite) TestNetVersion() {
	res, err := s.xplac.NetVersion().Query()
	s.Require().NoError(err)

	var netVersionResponse types.NetVersionResponse
	json.Unmarshal([]byte(res), &netVersionResponse)

	s.Require().Equal("47", netVersionResponse.NetVersion)
}

func (s *IntegrationTestSuite) TestNetPeerCount() {
	res, err := s.xplac.NetPeerCount().Query()
	s.Require().NoError(err)

	var netPeerCountResponse types.NetPeerCountResponse
	json.Unmarshal([]byte(res), &netPeerCountResponse)

	s.Require().Equal(0, netPeerCountResponse.NetPeerCount)
}

func (s *IntegrationTestSuite) TestNetListening() {
	res, err := s.xplac.NetListening().Query()
	s.Require().NoError(err)

	var netListeningResponse types.NetListeningResponse
	json.Unmarshal([]byte(res), &netListeningResponse)

	s.Require().Equal(true, netListeningResponse.NetListening)
}

func (s *IntegrationTestSuite) TestEthProtocolVersion() {
	res, err := s.xplac.EthProtocolVersion().Query()
	s.Require().NoError(err)

	var ethProtocolVersionResponse types.EthProtocolVersionResponse
	json.Unmarshal([]byte(res), &ethProtocolVersionResponse)

	s.Require().Equal(big.NewInt(65), ethProtocolVersionResponse.EthProtocolVersion)
}

func (s *IntegrationTestSuite) TestEthSyncing() {
	res, err := s.xplac.EthSyncing().Query()
	s.Require().NoError(err)

	var ethSyncingResponse types.EthSyncingResponse
	json.Unmarshal([]byte(res), &ethSyncingResponse)

	s.Require().Equal(false, ethSyncingResponse.EthSyncing)
}

func (s *IntegrationTestSuite) TestEthAccounts() {
	res, err := s.xplac.EthAccounts().Query()
	s.Require().NoError(err)

	var ethAccountsResponse types.EthAccountsResponse
	json.Unmarshal([]byte(res), &ethAccountsResponse)

	s.Require().Len(ethAccountsResponse.EthAccounts, 1)
}

func (s *IntegrationTestSuite) TestEthGetBlockTransactionCount() {
	testBlockHeight := util.FromInt64ToString(s.evmTestBlockHeight)
	blockRes, err := s.xplac.Block(types.BlockMsg{Height: testBlockHeight}).Query()
	s.Require().NoError(err)

	var resultBlock ctypes.ResultBlock
	json.Unmarshal([]byte(blockRes), &resultBlock)

	// block height
	ethGetBlockTransactionCountMsg := types.EthGetBlockTransactionCountMsg{
		BlockHeight: testBlockHeight,
	}
	res, err := s.xplac.EthGetBlockTransactionCount(ethGetBlockTransactionCountMsg).Query()
	s.Require().NoError(err)

	var ethGetBlockTransactionCountResponse types.EthGetBlockTransactionCountResponse
	json.Unmarshal([]byte(res), &ethGetBlockTransactionCountResponse)

	s.Require().Equal(big.NewInt(1), ethGetBlockTransactionCountResponse.EthGetBlockTransactionCount)

	// block hash
	ethGetBlockTransactionCountMsg = types.EthGetBlockTransactionCountMsg{
		BlockHash: resultBlock.BlockID.Hash.String(),
	}
	res, err = s.xplac.EthGetBlockTransactionCount(ethGetBlockTransactionCountMsg).Query()
	s.Require().NoError(err)

	json.Unmarshal([]byte(res), &ethGetBlockTransactionCountResponse)

	s.Require().Equal(big.NewInt(1), ethGetBlockTransactionCountResponse.EthGetBlockTransactionCount)
}

func (s *IntegrationTestSuite) TestEstimateGas() {
	var args []interface{}
	args = append(args, big.NewInt(1))

	invokeSolContractMsg := types.InvokeSolContractMsg{
		ContractAddress:      s.contractAddr,
		ContractFuncCallName: "store",
		Args:                 args,
		ABIJsonFilePath:      testABIJsonFilePath,
		BytecodeJsonFilePath: testBytecodeJsonFilePath,
		FromByteAddress:      util.FromStringToByte20Address(s.accounts[0].PubKey.Address().String()).String(),
	}
	res, err := s.xplac.EstimateGas(invokeSolContractMsg).Query()
	s.Require().NoError(err)

	var estimateGasResponse types.EstimateGasResponse
	json.Unmarshal([]byte(res), &estimateGasResponse)

	s.Require().NotEqual(uint64(0), estimateGasResponse.EstimateGas)
}

func (s *IntegrationTestSuite) TestEthGetTransactionByBlockHashAndIndex() {
	testBlockHeight := util.FromInt64ToString(s.evmTestBlockHeight)
	blockRes, err := s.xplac.Block(types.BlockMsg{Height: testBlockHeight}).Query()
	s.Require().NoError(err)

	var resultBlock ctypes.ResultBlock
	json.Unmarshal([]byte(blockRes), &resultBlock)

	getTransactionByBlockHashAndIndexMsg := types.GetTransactionByBlockHashAndIndexMsg{
		BlockHash: resultBlock.BlockID.Hash.String(),
		Index:     "0",
	}
	res, err := s.xplac.EthGetTransactionByBlockHashAndIndex(getTransactionByBlockHashAndIndexMsg).Query()
	s.Require().NoError(err)

	var transaction ethtypes.Transaction
	json.Unmarshal([]byte(res), &transaction)

	// unmarshal contract bytecode file
	bytecode, err := util.BytecodeParsing(testBytecodeJsonFilePath)
	s.Require().NoError(err)
	s.Require().Equal(bytecode, hex.EncodeToString(transaction.Data()))
}

func (s *IntegrationTestSuite) TestGetTransactionReceipt() {
	getTransactionReceiptMsg := types.GetTransactionReceiptMsg{
		TransactionHash: s.evmtestTxHash,
	}
	res, err := s.xplac.EthGetTransactionReceipt(getTransactionReceiptMsg).Query()
	s.Require().NoError(err)

	var receipt ethtypes.Receipt
	json.Unmarshal([]byte(res), &receipt)

	s.Require().Equal(s.evmtestTxHash, receipt.TxHash.String())
}

func (s *IntegrationTestSuite) TestEthNewFileter() {
	ethNewFilterMsg := types.EthNewFilterMsg{
		Topics:    []string{"0x20ec56d16231c4d7f761c2533885619489fface85cf6c478868ef1d531b93177"},
		Address:   []string{"0xf7777b36a51fb0b33dd0c5118361AfC94ff7f967"},
		ToBlock:   "latest",
		FromBlock: "earliest",
	}

	res, err := s.xplac.EthNewFilter(ethNewFilterMsg).Query()
	s.Require().NoError(err)

	var ethNewFilterResponse types.EthNewFilterResponse
	json.Unmarshal([]byte(res), &ethNewFilterResponse)

	s.Require().NotEqual("", ethNewFilterResponse.NewFilter)
}

func (s *IntegrationTestSuite) TestEthNewBlockFilter() {
	res, err := s.xplac.EthNewBlockFilter().Query()
	s.Require().NoError(err)

	var ethNewBlockFilterResponse types.EthNewBlockFilterResponse
	json.Unmarshal([]byte(res), &ethNewBlockFilterResponse)

	s.Require().NotEqual("", ethNewBlockFilterResponse.NewBlockFilter)
}

func (s *IntegrationTestSuite) TestEthNewPendingTransactionFilter() {
	res, err := s.xplac.EthNewPendingTransactionFilter().Query()
	s.Require().NoError(err)

	var ethNewPendingTransactionFilterResponse types.EthNewPendingTransactionFilterResponse
	json.Unmarshal([]byte(res), &ethNewPendingTransactionFilterResponse)

	s.Require().NotEqual("", ethNewPendingTransactionFilterResponse.NewPendingTransactionFilter)
}

func (s *IntegrationTestSuite) TestUninstallFilter() {
	filterRes, err := s.xplac.EthNewBlockFilter().Query()
	s.Require().NoError(err)

	var ethNewBlockFilterResponse types.EthNewBlockFilterResponse
	json.Unmarshal([]byte(filterRes), &ethNewBlockFilterResponse)

	ethUninstallFilterMsg := types.EthUninstallFilterMsg{
		FilterId: ethNewBlockFilterResponse.NewBlockFilter.(string),
	}
	res, err := s.xplac.EthUninstallFilter(ethUninstallFilterMsg).Query()
	s.Require().NoError(err)

	var ethUninstallFilterResponse types.EthUninstallFilterResponse
	json.Unmarshal([]byte(res), &ethUninstallFilterResponse)

	s.Require().Equal(true, ethUninstallFilterResponse.UninstallFilter)
}

func (s *IntegrationTestSuite) TestEthGetFilterChanges() {
	filterRes, err := s.xplac.EthNewBlockFilter().Query()
	s.Require().NoError(err)

	var ethNewBlockFilterResponse types.EthNewBlockFilterResponse
	json.Unmarshal([]byte(filterRes), &ethNewBlockFilterResponse)

	ethGetFilterChangesMsg := types.EthGetFilterChangesMsg{
		FilterId: ethNewBlockFilterResponse.NewBlockFilter.(string),
	}
	res, err := s.xplac.EthGetFilterChanges(ethGetFilterChangesMsg).Query()
	s.Require().NoError(err)

	var ethGetFilterChangesResponse types.EthGetFilterChangesResponse
	json.Unmarshal([]byte(res), &ethGetFilterChangesResponse)

	s.Require().Len(ethGetFilterChangesResponse.GetFilterChanges, 0)
}

func (s *IntegrationTestSuite) TestEthGetFilterLogs() {
	ethNewFilterMsg := types.EthNewFilterMsg{
		Topics:    []string{"0x20ec56d16231c4d7f761c2533885619489fface85cf6c478868ef1d531b93177"},
		Address:   []string{"0xf7777b36a51fb0b33dd0c5118361AfC94ff7f967"},
		ToBlock:   "latest",
		FromBlock: "earliest",
	}

	newFilterRes, err := s.xplac.EthNewFilter(ethNewFilterMsg).Query()
	s.Require().NoError(err)

	var ethNewFilterResponse types.EthNewFilterResponse
	json.Unmarshal([]byte(newFilterRes), &ethNewFilterResponse)

	ethGetFilterLogsMsg := types.EthGetFilterLogsMsg{
		FilterId: ethNewFilterResponse.NewFilter.(string),
	}
	res, err := s.xplac.EthGetFilterLogs(ethGetFilterLogsMsg).Query()
	s.Require().NoError(err)

	var ethGetFilterLogsResponse types.EthGetFilterLogsResponse
	json.Unmarshal([]byte(res), &ethGetFilterLogsResponse)

	s.T().Log(ethGetFilterLogsResponse.GetFilterLogs)
	s.Require().Equal(0, len(ethGetFilterLogsResponse.GetFilterLogs))
}

func (s *IntegrationTestSuite) TestEthCoinbase() {
	res, err := s.xplac.EthCoinbase().Query()
	s.Require().NoError(err)

	var ethCoinbaseResponse types.EthCoinbaseResponse
	json.Unmarshal([]byte(res), &ethCoinbaseResponse)

	s.Require().NotEqual("", ethCoinbaseResponse.Coinbase[2:])
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
