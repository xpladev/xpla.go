package client_test

import (
	"context"
	"fmt"
	"math/rand"
	"testing"

	"github.com/xpladev/xpla.go/client"
	mbank "github.com/xpladev/xpla.go/core/bank"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"

	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/ethereum/go-ethereum/common"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var (
	testBroadcastMode  = "sync"
	testTimeoutHeight  = "1000"
	testLcdUrl         = "https://cube-lcd.xpla.dev"
	testGrpcUrl        = "https://cube-grpc.xpla.dev"
	testRpcUrl         = "https://cube-rpc.xpla.dev"
	testEvmRpcUrl      = "https://cube-evm-rpc.xpla.dev"
	testOutputDocument = "./document.json"
)

func TestNewXplaClient(t *testing.T) {
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := testutil.RandomAccounts(r, 2)

	from := accounts[0]
	feegranter := accounts[1]
	gasLimitU64, err := util.FromStringToUint64(types.DefaultGasLimit)
	assert.NoError(t, err)
	gasPriceU64, err := util.FromStringToUint64(types.DefaultGasPrice)
	assert.NoError(t, err)

	feeAmount := util.FromUint64ToString(util.MulUint64(
		gasLimitU64,
		gasPriceU64,
	))
	testPagination := types.Pagination{
		PageKey:    "",
		Offset:     0,
		Limit:      0,
		CountTotal: false,
		Reverse:    true,
	}

	newClientOption := provider.Options{
		PrivateKey:     from.PrivKey,
		AccountNumber:  util.FromIntToString(types.DefaultAccNum),
		Sequence:       util.FromIntToString(types.DefaultAccSeq),
		BroadcastMode:  testBroadcastMode,
		GasLimit:       types.DefaultGasLimit,
		GasPrice:       types.DefaultGasPrice,
		GasAdjustment:  types.DefaultGasAdjustment,
		FeeAmount:      feeAmount,
		SignMode:       signing.SignMode_SIGN_MODE_DIRECT,
		FeeGranter:     feegranter.Address,
		TimeoutHeight:  testTimeoutHeight,
		LcdURL:         testLcdUrl,
		GrpcURL:        testGrpcUrl,
		RpcURL:         testRpcUrl,
		EvmRpcURL:      testEvmRpcUrl,
		Pagination:     testPagination,
		OutputDocument: testOutputDocument,
	}

	xplac := client.NewXplaClient(testutil.TestChainId).WithOptions(newClientOption)
	xplac.Total()

	totalMsg, err := mbank.MakeTotalSupplyMsg()
	assert.NoError(t, err)

	assert.Equal(t, testutil.TestChainId, xplac.GetChainId())
	assert.Equal(t, from.PrivKey, xplac.GetPrivateKey())
	assert.Equal(t, context.Background(), xplac.GetContext())
	assert.Equal(t, testLcdUrl, xplac.GetLcdURL())
	assert.Equal(t, testGrpcUrl, xplac.GetGrpcUrl())
	assert.Equal(t, testRpcUrl, xplac.GetRpc())
	assert.Equal(t, testEvmRpcUrl, xplac.GetEvmRpc())
	assert.Equal(t, testBroadcastMode, xplac.GetBroadcastMode())
	assert.Equal(t, util.FromIntToString(types.DefaultAccNum), xplac.GetAccountNumber())
	assert.Equal(t, util.FromIntToString(types.DefaultAccSeq), xplac.GetSequence())
	assert.Equal(t, types.DefaultGasLimit, xplac.GetGasLimit())
	assert.Equal(t, types.DefaultGasPrice, xplac.GetGasPrice())
	assert.Equal(t, types.DefaultGasAdjustment, xplac.GetGasAdjustment())
	assert.Equal(t, feeAmount, xplac.GetFeeAmount())
	assert.Equal(t, signing.SignMode_SIGN_MODE_DIRECT, xplac.GetSignMode())
	assert.Equal(t, feegranter.Address, xplac.GetFeeGranter())
	assert.Equal(t, testTimeoutHeight, xplac.GetTimeoutHeight())
	assert.Equal(t, testPagination.Reverse, xplac.GetPagination().Reverse)
	assert.Equal(t, testOutputDocument, xplac.GetOutputDocument())
	assert.Equal(t, mbank.BankModule, xplac.GetModule())
	assert.Equal(t, mbank.BankTotalMsgType, xplac.GetMsgType())
	assert.Equal(t, mbank.BankTotalMsgType, xplac.GetMsgType())
	assert.Equal(t, totalMsg, xplac.GetMsg())
}

var (
	validatorNumber = 2
	testSendAmount  = "1000"
)

type ClientTestSuite struct {
	suite.Suite

	xplac      provider.XplaClient
	apis       []string
	accounts   []simtypes.Account
	testTxHash string

	cfg     network.Config
	network network.Network
}

func NewClientTestSuite(cfg network.Config) *ClientTestSuite {
	return &ClientTestSuite{cfg: cfg}
}

func (s *ClientTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	src := rand.NewSource(1)
	r := rand.New(src)
	s.accounts = testutil.RandomAccounts(r, 2)
	s.testTxHash = "B6BBBB649F19E8970EF274C0083FE945FD38AD8C524D68BB3FE3A20D72DF03C4"

	balanceBigInt, err := util.FromStringToBigInt("1000000000000000000000000000")
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

	val := s.network.Validators[0]
	newAddr := s.accounts[0].Address

	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	s.xplac = client.NewXplaClient(testutil.TestChainId)

	s.apis = []string{
		s.network.Validators[0].APIAddress,
		s.network.Validators[0].AppConfig.GRPC.Address,
	}
}

func (s *ClientTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func TestClientTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewClientTestSuite(cfg))
}
