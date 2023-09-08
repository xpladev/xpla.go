package wasm_test

import (
	"encoding/json"
	"math/rand"
	"strings"
	"testing"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/suite"
)

var (
	validatorNumber   = 1
	testBalance       = "1000000000000000000000000000"
	testContractLabel = "test contract"
	testWasmFilePath  = "../../util/testutil/test_files/cw721_metadata_onchain.wasm"
)

type IntegrationTestSuite struct {
	suite.Suite

	xplac        provider.XplaClient
	apis         []string
	accounts     []simtypes.Account
	wasmCodeID   string
	contractAddr string

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

	s.xplac = client.NewXplaClient(testutil.TestChainId)
	s.apis = []string{
		s.network.Validators[0].APIAddress,
		s.network.Validators[0].AppConfig.GRPC.Address,
	}

	xplac := s.xplac.WithPrivateKey(s.accounts[0].PrivKey).WithURL(s.apis[0])

	// store wasm file
	storeMsg := types.StoreMsg{
		FilePath: testWasmFilePath,
	}
	txbytes, err := xplac.StoreCode(storeMsg).CreateAndSignTx()
	s.Require().NoError(err)

	storeTxRes, err := xplac.Broadcast(txbytes)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	queryTxMsg := types.QueryTxMsg{
		Value: storeTxRes.Response.TxHash,
	}
	storeTxQuery, err := xplac.Tx(queryTxMsg).Query()
	s.Require().NoError(err)

	var getTxResponse sdktx.GetTxResponse
	jsonpb.Unmarshal(strings.NewReader(storeTxQuery), &getTxResponse)

	s.wasmCodeID = getTxResponse.TxResponse.Logs[0].Events[1].Attributes[0].Value

	// instantiate contract
	instantiateMsg := types.InstantiateMsg{
		CodeId: s.wasmCodeID,
		Amount: "0",
		Label:  testContractLabel,
		InitMsg: `{
			"name":"cw721-metadata-onchain",
			"symbol":"CW721",
			"minter":"` + s.accounts[0].Address.String() + `"
	   }`,
		Admin: s.accounts[0].Address.String(),
	}
	txbytes, err = xplac.WithSequence("").InstantiateContract(instantiateMsg).CreateAndSignTx()
	s.Require().NoError(err)

	instTxRes, err := xplac.Broadcast(txbytes)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	queryTxMsg = types.QueryTxMsg{
		Value: instTxRes.Response.TxHash,
	}
	instTxQuery, err := xplac.Tx(queryTxMsg).Query()
	s.Require().NoError(err)

	jsonpb.Unmarshal(strings.NewReader(instTxQuery), &getTxResponse)

	s.contractAddr = getTxResponse.TxResponse.Logs[0].Events[0].Attributes[0].Value
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestQueryContract() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryMsg := types.QueryMsg{
			ContractAddress: s.contractAddr,
			QueryMsg:        `{"minter":{}}`,
		}
		res, err := s.xplac.QueryContract(queryMsg).Query()
		s.Require().NoError(err)

		type minterResponse struct {
			Data struct {
				Minter string `json:"minter"`
			} `json:"data"`
		}
		var m minterResponse

		json.Unmarshal([]byte(res), &m)
		s.Require().Equal(s.accounts[0].Address.String(), m.Data.Minter)
	}
	s.xplac = provider.ResetXplac(s.xplac)

}

func (s *IntegrationTestSuite) TestListCode() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.ListCode().Query()
		s.Require().NoError(err)

		var queryCodeResponse wasmtypes.QueryCodesResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryCodeResponse)

		s.Require().Equal(uint64(1), queryCodeResponse.CodeInfos[0].CodeID)
		s.Require().Equal(s.accounts[0].Address.String(), queryCodeResponse.CodeInfos[0].Creator)
		s.Require().Equal("2DD26686622A5BF5A94DF201867C82E638E3A139E3FDE30B5B8D33F37AF1CD89", queryCodeResponse.CodeInfos[0].DataHash.String())
		s.Require().Equal("", queryCodeResponse.CodeInfos[0].InstantiatePermission.Address)
		s.Require().Equal(wasmtypes.AccessTypeEverybody, queryCodeResponse.CodeInfos[0].InstantiatePermission.Permission)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestListContractByCode() {
	for i, api := range s.apis {
		// LCD cannot support list contract by code, replace code info
		if i == 0 {
			s.xplac.WithURL(api)

			listContractByCodeMsg := types.ListContractByCodeMsg{
				CodeId: s.wasmCodeID,
			}
			res, err := s.xplac.ListContractByCode(listContractByCodeMsg).Query()
			s.Require().NoError(err)

			var queryCodeResponse wasmtypes.QueryCodeResponse
			jsonpb.Unmarshal(strings.NewReader(res), &queryCodeResponse)

			s.Require().Equal(uint64(1), queryCodeResponse.CodeID)
			s.Require().Equal(s.accounts[0].Address.String(), queryCodeResponse.Creator)

		} else {
			s.xplac.WithGrpc(api)

			listContractByCodeMsg := types.ListContractByCodeMsg{
				CodeId: s.wasmCodeID,
			}
			res, err := s.xplac.ListContractByCode(listContractByCodeMsg).Query()
			s.Require().NoError(err)

			var queryContractsByCodeResponse wasmtypes.QueryContractsByCodeResponse
			jsonpb.Unmarshal(strings.NewReader(res), &queryContractsByCodeResponse)

			s.Require().Len(queryContractsByCodeResponse.Contracts, 1)
			s.Require().Equal(s.contractAddr, queryContractsByCodeResponse.Contracts[0])
		}

	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestCodeInfo() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		codeInfoMsg := types.CodeInfoMsg{
			CodeId: s.wasmCodeID,
		}
		res, err := s.xplac.CodeInfo(codeInfoMsg).Query()
		s.Require().NoError(err)

		var queryCodeResponse wasmtypes.QueryCodeResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryCodeResponse)

		s.Require().Equal(uint64(1), queryCodeResponse.CodeID)
		s.Require().Equal(s.accounts[0].Address.String(), queryCodeResponse.Creator)
		s.Require().Equal("2DD26686622A5BF5A94DF201867C82E638E3A139E3FDE30B5B8D33F37AF1CD89", queryCodeResponse.DataHash.String())
		s.Require().Equal("", queryCodeResponse.InstantiatePermission.Address)
		s.Require().Equal(wasmtypes.AccessTypeEverybody, queryCodeResponse.InstantiatePermission.Permission)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestContractInfo() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		contractInfoMsg := types.ContractInfoMsg{
			ContractAddress: s.contractAddr,
		}
		res, err := s.xplac.ContractInfo(contractInfoMsg).Query()
		s.Require().NoError(err)

		var queryContractInfoResponse wasmtypes.QueryContractInfoResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryContractInfoResponse)

		codeIdU64, err := util.FromStringToUint64(s.wasmCodeID)
		s.Require().NoError(err)

		s.Require().Equal(codeIdU64, queryContractInfoResponse.CodeID)
		s.Require().Equal(s.contractAddr, queryContractInfoResponse.Address)
		s.Require().Equal(s.accounts[0].Address.String(), queryContractInfoResponse.Admin)
		s.Require().Equal(testContractLabel, queryContractInfoResponse.Label)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestContractStateAll() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		contractStateAllMsg := types.ContractStateAllMsg{
			ContractAddress: s.contractAddr,
		}
		res, err := s.xplac.ContractStateAll(contractStateAllMsg).Query()
		s.Require().NoError(err)

		var queryAllContractStateResponse wasmtypes.QueryAllContractStateResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryAllContractStateResponse)

		s.Require().Equal("636F6E74726163745F696E666F", queryAllContractStateResponse.Models[0].Key.String())
		s.Require().Equal([]byte(`{"contract":"crates.io:cw721-metadata-onchain","version":"0.15.0"}`), queryAllContractStateResponse.Models[0].Value)
		s.Require().Equal("6D696E746572", queryAllContractStateResponse.Models[1].Key.String())
		s.Require().Equal([]byte(`"xpla1l03kma4vv9qcvhgcxf2ga0rnv7dqcumaxexssh"`), queryAllContractStateResponse.Models[1].Value)
		s.Require().Equal("6E66745F696E666F", queryAllContractStateResponse.Models[2].Key.String())
		s.Require().Equal([]byte(`{"name":"cw721-metadata-onchain","symbol":"CW721"}`), queryAllContractStateResponse.Models[2].Value)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestContractHistory() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		contractHistoryMsg := types.ContractHistoryMsg{
			ContractAddress: s.contractAddr,
		}
		res, err := s.xplac.ContractHistory(contractHistoryMsg).Query()
		s.Require().NoError(err)

		var queryContractHistoryResponse wasmtypes.QueryContractHistoryResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryContractHistoryResponse)

		codeIdU64, err := util.FromStringToUint64(s.wasmCodeID)
		s.Require().NoError(err)

		s.Require().Equal(codeIdU64, queryContractHistoryResponse.Entries[0].CodeID)
		s.Require().Equal(wasmtypes.ContractCodeHistoryOperationTypeInit, queryContractHistoryResponse.Entries[0].Operation)
		s.Require().Equal([]byte(`{"name":"cw721-metadata-onchain","symbol":"CW721","minter":"`+s.accounts[0].Address.String()+`"}`), queryContractHistoryResponse.Entries[0].Msg.Bytes())
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestPinned() {
	// have not proposal of the pinned
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.Pinned().Query()
		s.Require().NoError(err)

		var queryPinnedCodesResponse wasmtypes.QueryPinnedCodesResponse
		jsonpb.Unmarshal(strings.NewReader(res), &queryPinnedCodesResponse)

		s.Require().Len(queryPinnedCodesResponse.CodeIDs, 0)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestLibWasmvmVersion() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.LibwasmvmVersion().Query()
		s.Require().NoError(err)
		s.Require().Equal("1.0.1", res)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
