package client_test

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/evmos/ethermint/crypto/hd"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/jsonpb"
)

func (s *ClientTestSuite) TestBroadcast() {
	from := s.network.Validators[0].AdditionalAccount
	to := s.network.Validators[1].AdditionalAccount
	xplac := s.xplac.WithPrivateKey(s.network.Validators[0].AdditionalAccount.PrivKey)

	for i, api := range s.apis {
		if i == 0 {
			xplac.WithURL(api)
		} else {
			xplac.WithGrpc(api)
			nowSeq, err := util.FromStringToInt(xplac.GetSequence())
			s.Require().NoError(err)
			xplac.WithSequence(util.FromIntToString(nowSeq + 1))
		}

		// check before send
		bankBalancesMsg := types.BankBalancesMsg{
			Address: to.Address.String(),
		}
		beforeToRes, err := xplac.BankBalances(bankBalancesMsg).Query()
		s.Require().NoError(err)

		var beforeQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
		jsonpb.Unmarshal(strings.NewReader(beforeToRes), &beforeQueryAllBalancesResponse)

		// broadcast transaction - bank send
		bankSendMsg := types.BankSendMsg{
			FromAddress: from.Address.String(),
			ToAddress:   to.Address.String(),
			Amount:      testSendAmount,
		}
		txbytes, err := xplac.BankSend(bankSendMsg).CreateAndSignTx()
		s.Require().NoError(err)

		_, err = xplac.Broadcast(txbytes)
		s.Require().NoError(err)
		s.Require().NoError(s.network.WaitForNextBlock())

		// check after send
		bankBalancesMsg = types.BankBalancesMsg{
			Address: to.Address.String(),
		}
		afterToRes, err := xplac.BankBalances(bankBalancesMsg).Query()
		s.Require().NoError(err)

		var afterQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
		jsonpb.Unmarshal(strings.NewReader(afterToRes), &afterQueryAllBalancesResponse)

		s.Require().Equal(
			testSendAmount,
			afterQueryAllBalancesResponse.Balances[0].Amount.Sub(beforeQueryAllBalancesResponse.Balances[0].Amount).String(),
		)
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *ClientTestSuite) TestBroadcastMode() {
	invalidBroadcastMode := "invalid-mode"

	from := s.network.Validators[1].AdditionalAccount
	to := s.network.Validators[0].AdditionalAccount

	xplac := s.xplac.WithURL(s.apis[0]).
		WithPrivateKey(s.network.Validators[1].AdditionalAccount.PrivKey)

	modes := []string{"block", "async", "sync", "", invalidBroadcastMode}

	for _, mode := range modes {
		xplac.WithBroadcastMode(mode)

		bankSendMsg := types.BankSendMsg{
			FromAddress: from.Address.String(),
			ToAddress:   to.Address.String(),
			Amount:      testSendAmount,
		}
		txbytes, err := xplac.BankSend(bankSendMsg).CreateAndSignTx()
		s.Require().NoError(err)

		// if empty mode or invalid mode is changed to "sync"
		_, err = xplac.Broadcast(txbytes)
		s.Require().NoError(err)
		s.Require().NoError(s.network.WaitForNextBlock())

		nowSeq, err := util.FromStringToInt(xplac.GetSequence())
		s.Require().NoError(err)
		xplac.WithSequence(util.FromIntToString(nowSeq + 1))
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *ClientTestSuite) TestBroadcastEVM() {
	from := s.network.Validators[2].AdditionalAccount
	to := s.network.Validators[0].AdditionalAccount
	xplac := s.xplac.WithPrivateKey(s.network.Validators[2].AdditionalAccount.PrivKey).
		WithURL(s.apis[0]).
		WithEvmRpc("http://" + s.network.Validators[0].AppConfig.JSONRPC.Address)

	// check before send
	bankBalancesMsg := types.BankBalancesMsg{
		Address: to.Address.String(),
	}
	beforeToRes, err := xplac.BankBalances(bankBalancesMsg).Query()
	s.Require().NoError(err)

	var beforeQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
	jsonpb.Unmarshal(strings.NewReader(beforeToRes), &beforeQueryAllBalancesResponse)

	// broadcast transaction - evm send coin
	sendCoinMsg := types.SendCoinMsg{
		FromAddress: from.PubKey.Address().String(),
		ToAddress:   to.PubKey.Address().String(),
		Amount:      testSendAmount,
	}
	txbytes, err := xplac.EvmSendCoin(sendCoinMsg).CreateAndSignTx()
	s.Require().NoError(err)

	_, err = xplac.Broadcast(txbytes)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// check after send
	bankBalancesMsg = types.BankBalancesMsg{
		Address: to.Address.String(),
	}
	afterToRes, err := xplac.BankBalances(bankBalancesMsg).Query()
	s.Require().NoError(err)

	var afterQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
	jsonpb.Unmarshal(strings.NewReader(afterToRes), &afterQueryAllBalancesResponse)

	s.Require().Equal(
		testSendAmount,
		afterQueryAllBalancesResponse.Balances[0].Amount.Sub(beforeQueryAllBalancesResponse.Balances[0].Amount).String(),
	)
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *ClientTestSuite) TestBroadcastSolidityContract() {
	xplac := s.xplac.WithPrivateKey(s.network.Validators[3].AdditionalAccount.PrivKey).
		WithURL(s.apis[0]).
		WithEvmRpc("http://" + s.network.Validators[0].AppConfig.JSONRPC.Address)

	testABIJsonFilePath := "../util/testutil/test_files/abi.json"
	testBytecodeJsonFilePath := "../util/testutil/test_files/bytecode.json"

	deploySolContractMsg := types.DeploySolContractMsg{
		ABIJsonFilePath:      testABIJsonFilePath,
		BytecodeJsonFilePath: testBytecodeJsonFilePath,
		Args:                 nil,
	}
	txbytes, err := xplac.DeploySolidityContract(deploySolContractMsg).CreateAndSignTx()
	s.Require().NoError(err)

	_, err = xplac.Broadcast(txbytes)
	s.Require().NoError(err)

	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *ClientTestSuite) TestMultiSignature() {
	xplac := s.xplac.WithURL(s.apis[0])
	rootDir := s.network.Validators[0].Dir
	key1 := s.network.Validators[0].AdditionalAccount
	key2 := s.network.Validators[1].AdditionalAccount

	key1Name := "key1"
	key2Name := "key2"
	multiKeyName := "multiKey"

	// gen multisig account
	kb, err := keyring.New(types.XplaToolDefaultName, keyring.BackendTest, rootDir, nil, hd.EthSecp256k1Option())
	s.Require().NoError(err)

	armor1, err := key.EncryptArmorPrivKey(key1.PrivKey, key.DefaultEncryptPassphrase)
	s.Require().NoError(err)

	err = kb.ImportPrivKey(
		key1Name,
		armor1,
		key.DefaultEncryptPassphrase,
	)
	s.Require().NoError(err)

	armor2, err := key.EncryptArmorPrivKey(key2.PrivKey, key.DefaultEncryptPassphrase)
	s.Require().NoError(err)

	err = kb.ImportPrivKey(
		key2Name,
		armor2,
		key.DefaultEncryptPassphrase,
	)
	s.Require().NoError(err)

	var pks []cryptotypes.PubKey
	multisigThreshold := 2

	k1, err := kb.Key(key1Name)
	s.Require().NoError(err)
	pks = append(pks, k1.GetPubKey())

	k2, err := kb.Key(key2Name)
	s.Require().NoError(err)
	pks = append(pks, k2.GetPubKey())

	multiKeyInfo, err := kb.SaveMultisig(multiKeyName, multisig.NewLegacyAminoPubKey(multisigThreshold, pks))
	s.Require().NoError(err)

	// send coin to multisig account
	val := s.network.Validators[0]
	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		multiKeyInfo.GetAddress(),
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10000000000))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// create unsigned tx
	xplac.WithPrivateKey(key2.PrivKey)
	unsignedTxPath := filepath.Join(rootDir, "unsignedTx.json")
	xplac.WithOutputDocument(unsignedTxPath)

	bankSendMsg2 := types.BankSendMsg{
		FromAddress: multiKeyInfo.GetAddress().String(),
		ToAddress:   key1.Address.String(),
		Amount:      "10",
	}

	_, err = xplac.BankSend(bankSendMsg2).CreateUnsignedTx()
	s.Require().NoError(err)

	// create signature of key1
	xplac.WithPrivateKey(key1.PrivKey)
	signature1Path := filepath.Join(rootDir, "signature1.json")
	xplac.WithOutputDocument(signature1Path)

	signTxMsg1 := types.SignTxMsg{
		UnsignedFileName: unsignedTxPath,
		SignatureOnly:    true,
		MultisigAddress:  multiKeyInfo.GetAddress().String(),
		Overwrite:        false,
		Amino:            false,
		Offline:          false,
	}
	_, err = xplac.SignTx(signTxMsg1)
	s.Require().NoError(err)

	// create signature of key2
	xplac.WithPrivateKey(key2.PrivKey)
	signature2Path := filepath.Join(rootDir, "signature2.json")
	xplac.WithOutputDocument(signature2Path)

	signTxMsg2 := types.SignTxMsg{
		UnsignedFileName: unsignedTxPath,
		SignatureOnly:    true,
		MultisigAddress:  multiKeyInfo.GetAddress().String(),
		Overwrite:        false,
		Amino:            false,
		Offline:          false,
	}
	_, err = xplac.SignTx(signTxMsg2)
	s.Require().NoError(err)

	// create multisigned transaction
	xplac.WithOutputDocument("")
	txMultiSignMsg := types.TxMultiSignMsg{
		FileName:     unsignedTxPath,
		GenerateOnly: true,
		FromName:     multiKeyInfo.GetName(),
		Offline:      false,
		SignatureFiles: []string{
			signature1Path,
			signature2Path,
		},
		KeyringBackend: "test",
		KeyringPath:    rootDir,
	}
	multiSignTx, err := xplac.MultiSign(txMultiSignMsg)
	s.Require().NoError(err)

	// broadcast test
	sdkTx, err := xplac.GetEncoding().TxConfig.TxJSONDecoder()(multiSignTx)
	s.Require().NoError(err)
	txBytes, err := xplac.GetEncoding().TxConfig.TxEncoder()(sdkTx)
	s.Require().NoError(err)

	_, err = xplac.Broadcast(txBytes)

	// generate error insufficient funds
	// multisig tx is normal
	s.Require().Error(err)

	s.xplac = provider.ResetXplac(s.xplac)
}
