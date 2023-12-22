package auth_test

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"
	"github.com/xpladev/xpla.go/util/testutil/network"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ethermint "github.com/evmos/ethermint/types"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
)

var validatorNumber = 2

type IntegrationTestSuite struct {
	suite.Suite

	xplac      provider.XplaClient
	apis       []string
	accounts   []simtypes.Account
	testTxHash string

	cfg     network.Config
	network network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
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

	kb := s.network.Validators[0].ClientCtx.Keyring
	_, _, err = kb.NewMnemonic("newAccount", keyring.English, sdk.GetConfig().GetFullBIP44Path(), keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)

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

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestQueryParams() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.AuthParams().Query()
		s.Require().NoError(err)

		var authParamsResponse authtypes.QueryParamsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &authParamsResponse)
		s.Require().Equal(uint64(256), authParamsResponse.Params.MaxMemoCharacters)
		s.Require().Equal(uint64(7), authParamsResponse.Params.TxSigLimit)
		s.Require().Equal(uint64(10), authParamsResponse.Params.TxSizeCostPerByte)
		s.Require().Equal(uint64(590), authParamsResponse.Params.SigVerifyCostED25519)
		s.Require().Equal(uint64(1000), authParamsResponse.Params.SigVerifyCostSecp256k1)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestAccAddress() {
	validator := s.network.Validators[0]
	addr := validator.Address.String()

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryAccAddressMsg := types.QueryAccAddressMsg{
			Address: addr,
		}

		res, err := s.xplac.AccAddress(queryAccAddressMsg).Query()
		s.Require().NoError(err)

		var accountResponse authtypes.QueryAccountResponse
		jsonpb.Unmarshal(strings.NewReader(res), &accountResponse)

		var ethAccount ethermint.EthAccount
		proto.Unmarshal(accountResponse.Account.Value, &ethAccount)

		s.Require().Equal("/ethermint.crypto.v1.ethsecp256k1.PubKey", ethAccount.PubKey.TypeUrl)
		s.Require().Equal(addr, ethAccount.Address)
		s.Require().Equal(uint64(2), ethAccount.AccountNumber)

	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestAccounts() {
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		res, err := s.xplac.Accounts().Query()
		s.Require().NoError(err)

		var accountsResponse authtypes.QueryAccountsResponse
		jsonpb.Unmarshal(strings.NewReader(res), &accountsResponse)

		// 2 validators, 7 module accounts, 2 init accounts
		s.Require().Len(accountsResponse.Accounts, 11)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) TestQueryTxByEventAndQueryTx() {
	val := s.network.Validators[0]

	newAccount, err := val.ClientCtx.Keyring.Key("newAccount")
	s.Require().NoError(err)

	out, err := s.createBankMsg(val, newAccount.GetAddress(), sdk.NewCoins(sdk.NewInt64Coin(s.cfg.BondDenom, 1000)))
	s.Require().NoError(err)

	var txRes sdk.TxResponse
	err = val.ClientCtx.Codec.UnmarshalJSON(out.Bytes(), &txRes)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	txHash := txRes.TxHash

	s.xplac.WithRpc(s.network.Validators[0].RPCAddress)
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
		}

		queryTxsByEventsMsg := types.QueryTxsByEventsMsg{
			Events: fmt.Sprintf("tx.fee=%s",
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
		}

		queryTxMsg := types.QueryTxMsg{
			Value: txHash,
		}

		res, err := s.xplac.TxsByEvents(queryTxsByEventsMsg).Query()
		s.Require().NoError(err)

		res1, err := s.xplac.Tx(queryTxMsg).Query()
		s.Require().NoError(err)

		if i == 0 {
			var getTxsEventResponse tx.GetTxsEventResponse
			jsonpb.Unmarshal(strings.NewReader(res), &getTxsEventResponse)

			s.Require().Equal(2, len(getTxsEventResponse.TxResponses))
			s.Require().Equal(2, len(getTxsEventResponse.Txs))

			var txResponse tx.GetTxResponse
			jsonpb.Unmarshal(strings.NewReader(res1), &txResponse)

			s.Require().Equal(txHash, txResponse.TxResponse.TxHash)
		} else {
			var searchTxsResult sdk.SearchTxsResult
			jsonpb.Unmarshal(strings.NewReader(res), &searchTxsResult)

			s.Require().Equal(2, len(searchTxsResult.Txs))

			var txResponse sdk.TxResponse
			jsonpb.Unmarshal(strings.NewReader(res1), &txResponse)

			s.Require().Equal(txHash, txResponse.TxHash)
		}
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) createBankMsg(val *network.Validator, toAddr sdk.AccAddress, amount sdk.Coins, extraFlags ...string) (sdktestutil.BufferWriter, error) {
	flags := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	flags = append(flags, extraFlags...)
	return banktestutil.MsgSendExec(val.ClientCtx, val.Address, toAddr, amount, flags...)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
