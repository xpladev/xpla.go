package client_test

import (
	"path/filepath"
	"strings"

	"github.com/evmos/ethermint/crypto/hd"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/gogo/protobuf/jsonpb"
)

func (s *ClientTestSuite) TestBroadcast() {
	from := s.accounts[0]
	to := s.accounts[1]
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)

	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
		} else {
			s.xplac.WithGrpc(api)
			nowSeq, err := util.FromStringToInt(s.xplac.GetSequence())
			s.Require().NoError(err)
			s.xplac.WithSequence(util.FromIntToString(nowSeq + 1))
		}

		// check before send
		bankBalancesMsg := types.BankBalancesMsg{
			Address: to.Address.String(),
		}
		beforeToRes, err := s.xplac.BankBalances(bankBalancesMsg).Query()
		s.Require().NoError(err)

		var beforeQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
		jsonpb.Unmarshal(strings.NewReader(beforeToRes), &beforeQueryAllBalancesResponse)

		// broadcast transaction - bank send
		bankSendMsg := types.BankSendMsg{
			FromAddress: from.Address.String(),
			ToAddress:   to.Address.String(),
			Amount:      testSendAmount,
		}
		txbytes, err := s.xplac.BankSend(bankSendMsg).CreateAndSignTx()
		s.Require().NoError(err)

		_, err = s.xplac.Broadcast(txbytes)
		s.Require().NoError(err)
		s.Require().NoError(s.network.WaitForNextBlock())

		// check after send
		bankBalancesMsg = types.BankBalancesMsg{
			Address: to.Address.String(),
		}
		afterToRes, err := s.xplac.BankBalances(bankBalancesMsg).Query()
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

	from := s.accounts[0]
	to := s.accounts[1]

	s.xplac.
		WithURL(s.apis[0]).
		WithPrivateKey(s.accounts[0].PrivKey)

	modes := []string{"block", "async", "sync", "", invalidBroadcastMode}

	for _, mode := range modes {
		s.xplac.WithBroadcastMode(mode)

		bankSendMsg := types.BankSendMsg{
			FromAddress: from.Address.String(),
			ToAddress:   to.Address.String(),
			Amount:      testSendAmount,
		}
		txbytes, err := s.xplac.BankSend(bankSendMsg).CreateAndSignTx()
		s.Require().NoError(err)

		// if empty mode or invalid mode is changed to "sync"
		_, err = s.xplac.Broadcast(txbytes)
		s.Require().NoError(err)
		s.Require().NoError(s.network.WaitForNextBlock())

		nowSeq, err := util.FromStringToInt(s.xplac.GetSequence())
		s.Require().NoError(err)
		s.xplac.WithSequence(util.FromIntToString(nowSeq + 1))
	}
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *ClientTestSuite) TestBroadcastEVM() {
	from := s.accounts[0]
	to := s.accounts[1]
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey).
		WithURL(s.apis[0]).
		WithEvmRpc("http://" + s.network.Validators[0].AppConfig.JSONRPC.Address)

	// check before send
	bankBalancesMsg := types.BankBalancesMsg{
		Address: to.Address.String(),
	}
	beforeToRes, err := s.xplac.BankBalances(bankBalancesMsg).Query()
	s.Require().NoError(err)

	var beforeQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
	jsonpb.Unmarshal(strings.NewReader(beforeToRes), &beforeQueryAllBalancesResponse)

	// broadcast transaction - evm send coin
	sendCoinMsg := types.SendCoinMsg{
		FromAddress: from.PubKey.Address().String(),
		ToAddress:   to.PubKey.Address().String(),
		Amount:      testSendAmount,
	}
	txbytes, err := s.xplac.EvmSendCoin(sendCoinMsg).CreateAndSignTx()
	s.Require().NoError(err)

	_, err = s.xplac.Broadcast(txbytes)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// check after send
	bankBalancesMsg = types.BankBalancesMsg{
		Address: to.Address.String(),
	}
	afterToRes, err := s.xplac.BankBalances(bankBalancesMsg).Query()
	s.Require().NoError(err)

	var afterQueryAllBalancesResponse banktypes.QueryAllBalancesResponse
	jsonpb.Unmarshal(strings.NewReader(afterToRes), &afterQueryAllBalancesResponse)

	s.Require().Equal(
		testSendAmount,
		afterQueryAllBalancesResponse.Balances[0].Amount.Sub(beforeQueryAllBalancesResponse.Balances[0].Amount).String(),
	)
	s.xplac = provider.ResetXplac(s.xplac)
}

func (s *ClientTestSuite) TestMultiSignature() {
	s.xplac.WithURL(s.apis[0])
	rootDir := s.network.Validators[0].Dir
	key1 := s.accounts[0]
	key2 := s.accounts[1]

	key1Name := "key1"
	key2Name := "key2"
	multiKeyName := "multiKey"

	// gen multisig account
	kb, err := keyring.New(types.XplaToolDefaultName, keyring.BackendTest, rootDir, nil, hd.EthSecp256k1Option())
	s.Require().NoError(err)

	err = kb.ImportPrivKey(
		key1Name,
		key.EncryptArmorPrivKey(key1.PrivKey, key.DefaultEncryptPassphrase),
		key.DefaultEncryptPassphrase,
	)
	s.Require().NoError(err)

	err = kb.ImportPrivKey(
		key2Name,
		key.EncryptArmorPrivKey(key2.PrivKey, key.DefaultEncryptPassphrase),
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
	s.xplac.WithPrivateKey(key2.PrivKey)

	bankSendMsg := types.BankSendMsg{
		FromAddress: key2.Address.String(),
		ToAddress:   multiKeyInfo.GetAddress().String(),
		Amount:      "10000000000axpla",
	}
	bankSendTxbytes, err := s.xplac.BankSend(bankSendMsg).CreateAndSignTx()
	s.Require().NoError(err)

	_, err = s.xplac.Broadcast(bankSendTxbytes)
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// create unsigned tx
	unsignedTxPath := filepath.Join(rootDir, "unsignedTx.json")
	s.xplac.WithOutputDocument(unsignedTxPath)

	bankSendMsg2 := types.BankSendMsg{
		FromAddress: multiKeyInfo.GetAddress().String(),
		ToAddress:   key1.Address.String(),
		Amount:      "10",
	}

	_, err = s.xplac.BankSend(bankSendMsg2).CreateUnsignedTx()
	s.Require().NoError(err)

	// create signature of key1
	s.xplac.WithPrivateKey(key1.PrivKey)
	signature1Path := filepath.Join(rootDir, "signature1.json")
	s.xplac.WithOutputDocument(signature1Path)

	signTxMsg1 := types.SignTxMsg{
		UnsignedFileName: unsignedTxPath,
		SignatureOnly:    true,
		MultisigAddress:  multiKeyInfo.GetAddress().String(),
		Overwrite:        false,
		Amino:            false,
		Offline:          false,
	}
	_, err = s.xplac.SignTx(signTxMsg1)
	s.Require().NoError(err)

	// create signature of key2
	s.xplac.WithPrivateKey(key2.PrivKey)
	signature2Path := filepath.Join(rootDir, "signature2.json")
	s.xplac.WithOutputDocument(signature2Path)

	signTxMsg2 := types.SignTxMsg{
		UnsignedFileName: unsignedTxPath,
		SignatureOnly:    true,
		MultisigAddress:  multiKeyInfo.GetAddress().String(),
		Overwrite:        false,
		Amino:            false,
		Offline:          false,
	}
	_, err = s.xplac.SignTx(signTxMsg2)
	s.Require().NoError(err)

	// create multisigned transaction
	s.xplac.WithOutputDocument("")
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
	multiSignTx, err := s.xplac.MultiSign(txMultiSignMsg)
	s.Require().NoError(err)

	// broadcast test
	sdkTx, err := s.xplac.GetEncoding().TxConfig.TxJSONDecoder()(multiSignTx)
	s.Require().NoError(err)
	txBytes, err := s.xplac.GetEncoding().TxConfig.TxEncoder()(sdkTx)
	s.Require().NoError(err)

	_, err = s.xplac.Broadcast(txBytes)

	// generate error insufficient funds
	// multisig tx is normal
	s.Require().Error(err)
	s.Require().Equal(
		`code 8 : tx failed - [with code 5 : 10000000000axpla is smaller than 133715200000000000axpla: insufficient funds: insufficient funds]`,
		err.Error(),
	)

	s.xplac = provider.ResetXplac(s.xplac)
}
