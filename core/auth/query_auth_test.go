package auth_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil/network"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankcli "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	ethermint "github.com/evmos/ethermint/types"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
)

var validatorNumber = 2

type IntegrationTestSuite struct {
	suite.Suite

	xplac *client.XplaClient
	apis  []string

	cfg     network.Config
	network *network.Network
}

func NewIntegrationTestSuite(cfg network.Config) *IntegrationTestSuite {
	return &IntegrationTestSuite{cfg: cfg}
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	s.network = network.New(s.T(), s.cfg)
	s.Require().NoError(s.network.WaitForNextBlock())

	kb := s.network.Validators[0].ClientCtx.Keyring
	_, _, err := kb.NewMnemonic("newAccount", keyring.English, sdk.GetConfig().GetFullBIP44Path(), keyring.DefaultBIP39Passphrase, hd.Secp256k1)
	s.Require().NoError(err)

	s.xplac = client.NewTestXplaClient()
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
	s.xplac = client.ResetXplac(s.xplac)
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
		s.Require().Equal(uint64(0), ethAccount.AccountNumber)

	}
	s.xplac = client.ResetXplac(s.xplac)
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

		// 2 validator, 7 module accounts
		s.Require().Len(accountsResponse.Accounts, 9)
	}
	s.xplac = client.ResetXplac(s.xplac)
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

			s.Require().Equal(1, len(getTxsEventResponse.TxResponses))
			s.Require().Equal(1, len(getTxsEventResponse.Txs))

			var txResponse tx.GetTxResponse
			jsonpb.Unmarshal(strings.NewReader(res1), &txResponse)

			s.Require().Equal(txHash, txResponse.TxResponse.TxHash)
		} else {
			var searchTxsResult sdk.SearchTxsResult
			jsonpb.Unmarshal(strings.NewReader(res), &searchTxsResult)

			s.Require().Equal(1, len(searchTxsResult.Txs))

			var txResponse sdk.TxResponse
			jsonpb.Unmarshal(strings.NewReader(res1), &txResponse)

			s.Require().Equal(txHash, txResponse.TxHash)
		}
	}
	s.xplac = client.ResetXplac(s.xplac)
}

func (s *IntegrationTestSuite) createBankMsg(val *network.Validator, toAddr sdk.AccAddress, amount sdk.Coins, extraFlags ...string) (sdktestutil.BufferWriter, error) {
	flags := []string{
		fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees,
			sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	}

	flags = append(flags, extraFlags...)
	return bankcli.MsgSendExec(val.ClientCtx, val.Address, toAddr, amount, flags...)
}

func TestIntegrationTestSuite(t *testing.T) {
	cfg := network.DefaultConfig()
	cfg.NumValidators = validatorNumber
	suite.Run(t, NewIntegrationTestSuite(cfg))
}
