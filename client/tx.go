package client

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"

	mevm "github.com/xpladev/xpla.go/core/evm"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	kmultisig "github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/crypto/types/multisig"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authcli "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
)

// Create and sign a transaction before it is broadcasted to xpla chain.
// Options required for create and sign are stored in the xpla client and reflected when the values of those options exist.
// Create and sign transaction must be needed in order to send transaction to the chain.
func (xplac *XplaClient) CreateAndSignTx() ([]byte, error) {
	var err error
	if xplac.GetErr() != nil {
		return nil, xplac.GetErr()
	}

	xplac, err = GetAccNumAndSeq(xplac)
	if err != nil {
		return nil, err
	}

	if xplac.GetGasAdjustment() == "" {
		xplac.WithGasAdjustment(types.DefaultGasAdjustment)
	}

	if xplac.GetGasPrice() == "" {
		xplac.WithGasPrice(types.DefaultGasPrice)
	}

	if xplac.GetModule() == mevm.EvmModule {
		return xplac.createAndSignEvmTx()

	} else {
		builder, err := setTxBuilderMsg(xplac)
		if err != nil {
			return nil, err
		}

		gasLimit, feeAmount, err := getGasLimitFeeAmount(xplac, builder)
		if err != nil {
			return nil, err
		}

		builder, err = convertAndSetBuilder(xplac, builder, gasLimit, feeAmount)
		if err != nil {
			return nil, err
		}

		// Set default sign mode (DIRECT=1)
		if xplac.GetSignMode() == signing.SignMode_SIGN_MODE_UNSPECIFIED {
			xplac.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
		}

		privs := []cryptotypes.PrivKey{xplac.GetPrivateKey()}

		accNumU64, err := util.FromStringToUint64(xplac.GetAccountNumber())
		if err != nil {
			return nil, err
		}
		accSeqU64, err := util.FromStringToUint64(xplac.GetSequence())
		if err != nil {
			return nil, err
		}
		accNums := []uint64{accNumU64}
		accSeqs := []uint64{accSeqU64}

		var sigsV2 []signing.SignatureV2

		err = txSignRound(xplac, sigsV2, privs, accSeqs, accNums, builder)
		if err != nil {
			return nil, err
		}

		sdkTx := builder.GetTx()
		txBytes, err := xplac.GetEncoding().TxConfig.TxEncoder()(sdkTx)
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}

		if xplac.GetOutputDocument() != "" {
			jsonTx, err := xplac.GetEncoding().TxConfig.TxJSONEncoder()(sdkTx)
			if err != nil {
				return nil, util.LogErr(errors.ErrParse, err)
			}
			err = util.SaveJsonPretty(jsonTx, xplac.GetOutputDocument())
			if err != nil {
				return nil, err
			}

			return jsonTx, nil
		}

		return txBytes, nil
	}

}

// Create transaction with unsigning.
// It returns txbytes of byte type when output document options is nil.
// If not, save the unsigned transaction file which name is "outputDocument"
func (xplac *XplaClient) CreateUnsignedTx() ([]byte, error) {
	if xplac.GetErr() != nil {
		return nil, xplac.GetErr()
	}
	builder, err := setTxBuilderMsg(xplac)
	if err != nil {
		return nil, err
	}

	gasLimit, feeAmount, err := getGasLimitFeeAmount(xplac, builder)
	if err != nil {
		return nil, err
	}

	builder, err = convertAndSetBuilder(xplac, builder, gasLimit, feeAmount)
	if err != nil {
		return nil, err
	}

	sdkTx := builder.GetTx()
	txBytes, err := xplac.GetEncoding().TxConfig.TxEncoder()(sdkTx)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	if xplac.GetOutputDocument() != "" {
		jsonTx, err := xplac.GetEncoding().TxConfig.TxJSONEncoder()(sdkTx)
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}
		err = util.SaveJsonPretty(jsonTx, xplac.GetOutputDocument())
		if err != nil {
			return nil, err
		}

		return jsonTx, nil
	}

	return txBytes, nil
}

// Sign created unsigned transaction.
func (xplac *XplaClient) SignTx(signTxMsg types.SignTxMsg) ([]byte, error) {
	var err error

	if xplac.GetErr() != nil {
		return nil, xplac.GetErr()
	}
	var emptySignTxMsg types.SignTxMsg
	if signTxMsg == emptySignTxMsg {
		return nil, util.LogErr(errors.ErrNotSatisfiedOptions, "need sign tx message of xpla client's option")
	}

	if !signTxMsg.Offline {
		xplac, err = GetAccNumAndSeq(xplac)
		if err != nil {
			return nil, err
		}
	}

	clientCtx, err := util.NewClient()
	if err != nil {
		return nil, err
	}

	fromName := types.XplaToolDefaultName
	err = clientCtx.Keyring.ImportPrivKey(fromName, key.EncryptArmorPrivKey(xplac.GetPrivateKey(), key.DefaultEncryptPassphrase), key.DefaultEncryptPassphrase)
	if err != nil {
		return nil, util.LogErr(errors.ErrKeyNotFound, err)
	}

	clientCtx.WithSignModeStr("direct")

	clientCtx, txFactory, newTx, err := readTxAndInitContexts(clientCtx, signTxMsg.UnsignedFileName)
	if err != nil {
		return nil, err
	}

	txCfg := clientCtx.TxConfig
	txBuilder, err := txCfg.WrapTxBuilder(newTx)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	signatureOnly := signTxMsg.SignatureOnly
	multisig := signTxMsg.MultisigAddress
	if multisig != "" {
		multisigAddr, err := sdk.AccAddressFromBech32(multisig)
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}

		multisigAccNum := uint64(types.DefaultAccNum)
		multisigAccSeq := uint64(types.DefaultAccSeq)
		if !signTxMsg.Offline {
			if xplac.GetLcdURL() == "" && xplac.GetGrpcUrl() == "" {
				return nil, util.LogErr(errors.ErrInvalidRequest, "need LCD or gRPC URL when not offline mode")
			}
			signerAcc, err := xplac.LoadAccount(multisigAddr)
			if err != nil {
				return nil, util.LogErr(errors.ErrParse, err)
			}
			multisigAccNum = signerAcc.GetAccountNumber()
			multisigAccSeq = signerAcc.GetSequence()
		}

		txFactory = txFactory.WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON).
			WithChainID(xplac.GetChainId()).
			WithAccountNumber(multisigAccNum).
			WithSequence(multisigAccSeq)

		if !isTxSigner(multisigAddr, txBuilder.GetTx().GetSigners()) {
			return nil, util.LogErr(errors.ErrParse, err)
		}

		err = tx.Sign(txFactory, fromName, txBuilder, signTxMsg.Overwrite)
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}
		signatureOnly = true
	} else {
		// Set default sign mode (DIRECT=1)
		if xplac.GetSignMode() == signing.SignMode_SIGN_MODE_UNSPECIFIED {
			xplac.WithSignMode(signing.SignMode_SIGN_MODE_DIRECT)
		}

		accNumU64, err := util.FromStringToUint64(xplac.GetAccountNumber())
		if err != nil {
			return nil, err
		}
		accSeqU64, err := util.FromStringToUint64(xplac.GetSequence())
		if err != nil {
			return nil, err
		}

		privs := []cryptotypes.PrivKey{xplac.GetPrivateKey()}
		accNums := []uint64{accNumU64}
		accSeqs := []uint64{accSeqU64}

		var sigsV2 []signing.SignatureV2

		err = txSignRound(xplac, sigsV2, privs, accSeqs, accNums, txBuilder)
		if err != nil {
			return nil, err
		}
	}

	var json []byte
	if signTxMsg.Amino {
		stdTx, err := tx.ConvertTxToStdTx(clientCtx.LegacyAmino, txBuilder.GetTx())
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}
		req := authcli.BroadcastReq{
			Tx:   stdTx,
			Mode: "block|sync|async",
		}
		json, err = clientCtx.LegacyAmino.MarshalJSON(req)
		if err != nil {
			return nil, util.LogErr(errors.ErrFailedToMarshal, err)
		}
	} else {
		json, err = marshalSignatureJSON(txCfg, txBuilder, signatureOnly)
		if err != nil {
			return nil, err
		}
	}

	if xplac.GetOutputDocument() != "" {
		err = util.SaveJsonPretty(json, xplac.GetOutputDocument())
		if err != nil {
			return nil, err
		}

		return json, nil
	}

	return json, nil
}

// Sign created unsigned transaction with multi signatures.
func (xplac *XplaClient) MultiSign(txMultiSignMsg types.TxMultiSignMsg) ([]byte, error) {
	if xplac.GetErr() != nil {
		return nil, xplac.GetErr()
	}

	clientCtx, err := util.NewClient()
	if err != nil {
		return nil, err
	}

	if txMultiSignMsg.KeyringBackend != keyring.BackendFile &&
		txMultiSignMsg.KeyringBackend != keyring.BackendMemory &&
		txMultiSignMsg.KeyringBackend != keyring.BackendTest {
		return nil, util.LogErr(errors.ErrParse, "invalid keyring backend, must be "+util.BackendFile+", "+util.BackendTest+" or "+util.BackendMemory)
	}

	keyringPath := txMultiSignMsg.KeyringPath
	if (txMultiSignMsg.KeyringBackend == keyring.BackendFile ||
		txMultiSignMsg.KeyringBackend == keyring.BackendTest) && txMultiSignMsg.KeyringPath == "" {
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}

		keyringPath = filepath.Join(userHomeDir, ".xpla")
	}

	newKeyring, err := util.NewKeyring(txMultiSignMsg.KeyringBackend, keyringPath)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	clientCtx = clientCtx.WithKeyring(newKeyring)

	parseTx, err := authclient.ReadTxFromFile(clientCtx, txMultiSignMsg.FileName)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	txFactory := util.NewFactory(clientCtx)
	if txFactory.SignMode() == signing.SignMode_SIGN_MODE_UNSPECIFIED {
		txFactory = txFactory.WithSignMode(signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON)
	}
	txFactory = txFactory.
		WithChainID(xplac.GetChainId()).
		WithAccountNumber(uint64(types.DefaultAccNum)).
		WithSequence(uint64(types.DefaultAccSeq))

	txCfg := clientCtx.TxConfig
	txBuilder, err := txCfg.WrapTxBuilder(parseTx)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	multisigInfo, err := getMultisigInfo(clientCtx, txMultiSignMsg.FromName)
	if err != nil {
		return nil, err
	}

	multisigPub := multisigInfo.GetPubKey().(*kmultisig.LegacyAminoPubKey)
	multisigSig := multisig.NewMultisig(len(multisigPub.PubKeys))
	if !txMultiSignMsg.Offline {
		if xplac.GetLcdURL() == "" && xplac.GetGrpcUrl() == "" {
			return nil, util.LogErr(errors.ErrInvalidRequest, "need LCD or gRPC URL when not offline mode")
		}
		multisigAccount, err := xplac.LoadAccount(multisigInfo.GetAddress())
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}
		txFactory = txFactory.
			WithAccountNumber(multisigAccount.GetAccountNumber()).
			WithSequence(multisigAccount.GetSequence())
	}

	for _, sigFile := range txMultiSignMsg.SignatureFiles {
		sigs, err := unmarshalSignatureJSON(clientCtx, sigFile)
		if err != nil {
			return nil, util.LogErr(errors.ErrFailedToUnmarshal, err)
		}

		for _, sig := range sigs {
			data, ok := sig.Data.(*signing.SingleSignatureData)
			if !ok {
				return nil, util.LogErr(errors.ErrParse, "signature data is not single signature")
			}

			if data.SignMode != signing.SignMode_SIGN_MODE_LEGACY_AMINO_JSON {
				continue
			}

			addr, err := sdk.AccAddressFromHex(sig.PubKey.Address().String())
			if err != nil {
				return nil, util.LogErr(errors.ErrParse, err)
			}

			signingData := authsigning.SignerData{
				ChainID:       txFactory.ChainID(),
				AccountNumber: txFactory.AccountNumber(),
				Sequence:      txFactory.Sequence(),
			}

			err = authsigning.VerifySignature(sig.PubKey, signingData, sig.Data, txCfg.SignModeHandler(), txBuilder.GetTx())
			if err != nil {

				return nil, util.LogErr(errors.ErrInvalidRequest, "couldn't verify signature for address", addr.String())
			}

			if err := multisig.AddSignatureV2(multisigSig, sig, multisigPub.GetPubKeys()); err != nil {
				return nil, util.LogErr(errors.ErrParse, err)
			}
		}
	}

	sigV2 := signing.SignatureV2{
		PubKey:   multisigPub,
		Data:     multisigSig,
		Sequence: txFactory.Sequence(),
	}

	err = txBuilder.SetSignatures(sigV2)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	sigOnly := txMultiSignMsg.SignatureOnly
	aminoJson := txMultiSignMsg.Amino

	var json []byte
	if aminoJson {
		stdTx, err := tx.ConvertTxToStdTx(clientCtx.LegacyAmino, txBuilder.GetTx())
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}

		req := authcli.BroadcastReq{
			Tx:   stdTx,
			Mode: "block|sync|async",
		}

		json, err = clientCtx.LegacyAmino.MarshalJSON(req)
		if err != nil {
			return nil, err
		}
	} else {
		json, err = marshalSignatureJSON(txCfg, txBuilder, sigOnly)
		if err != nil {
			return nil, err
		}
	}

	if xplac.GetOutputDocument() != "" {
		err = util.SaveJsonPretty(json, xplac.GetOutputDocument())
		if err != nil {
			return nil, err
		}
	}

	return json, nil
}

// Create and sign transaction of evm.
func (xplac *XplaClient) createAndSignEvmTx() ([]byte, error) {
	ethPrivKey, err := toECDSA(xplac.GetPrivateKey())
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	chainId, err := util.ConvertEvmChainId(xplac.GetChainId())
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	if xplac.GetOutputDocument() != "" {
		util.LogInfo("no create output document as tx of evm")
	}

	gasPrice, err := util.FromStringToBigInt(xplac.GetGasPrice())
	if err != nil {
		return nil, err
	}

	switch {
	case xplac.GetMsgType() == mevm.EvmSendCoinMsgType:
		gasLimit := xplac.GetGasLimit()
		if gasLimit == "" {
			gasLimitU64, err := util.FromStringToUint64(util.DefaultEvmGasLimit)
			if err != nil {
				return nil, util.LogErr(errors.ErrParse, err)
			}
			gasLimitAdjustment, err := util.GasLimitAdjustment(gasLimitU64, xplac.GetGasAdjustment())
			if err != nil {
				return nil, util.LogErr(errors.ErrParse, err)
			}
			gasLimit = gasLimitAdjustment
		}

		convertMsg, ok := xplac.GetMsg().(types.SendCoinMsg)
		if !ok {
			return nil, util.LogErr(errors.ErrParse, "invalid msg")
		}
		toAddr := util.FromStringToByte20Address(convertMsg.ToAddress)
		amount, err := util.FromStringToBigInt(convertMsg.Amount)
		if err != nil {
			return nil, err
		}

		return evmTxSignRound(xplac, toAddr, gasPrice, gasLimit, amount, nil, chainId, ethPrivKey)

	case xplac.GetMsgType() == mevm.EvmDeploySolContractMsgType:
		gasLimit := xplac.GetGasLimit()
		if gasLimit == "" {
			gasLimit = "0"
		}

		convertMsg, ok := xplac.GetMsg().(mevm.ContractInfo)
		if !ok {
			return nil, util.LogErr(errors.ErrParse, "invalid msg")
		}
		nonce, err := util.FromStringToBigInt(xplac.GetSequence())
		if err != nil {
			return nil, err
		}

		value, err := util.FromStringToBigInt(util.DefaultSolidityValue)
		if err != nil {
			return nil, err
		}

		gasLimitU64, err := util.FromStringToUint64(gasLimit)
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}

		tx := mevm.DeploySolTx{
			ChainId:  chainId,
			Nonce:    nonce,
			Value:    value,
			GasLimit: gasLimitU64,
			GasPrice: gasPrice,
			ABI:      convertMsg.Abi,
			Bytecode: convertMsg.Bytecode,
		}

		txbytes, err := util.JsonMarshalData(tx)
		if err != nil {
			return nil, util.LogErr(errors.ErrFailedToMarshal, err)
		}

		return txbytes, nil

	case xplac.GetMsgType() == mevm.EvmInvokeSolContractMsgType:
		convertMsg, ok := xplac.GetMsg().(types.InvokeSolContractMsg)
		if !ok {
			return nil, util.LogErr(errors.ErrParse, "invalid msg")
		}
		var invokeByteData []byte
		invokeByteData, err = util.GetAbiPack(convertMsg.ContractFuncCallName, convertMsg.ABI, convertMsg.Bytecode, mevm.Args...)
		mevm.Args = nil
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}

		toAddr := util.FromStringToByte20Address(convertMsg.ContractAddress)
		amount, err := util.FromStringToBigInt("0")
		if err != nil {
			return nil, util.LogErr(errors.ErrEvmRpcRequest, err)
		}

		gasLimit := xplac.GetGasLimit()
		if gasLimit == "" {
			estimateGas, err := xplac.EstimateGas(convertMsg).Query()
			if err != nil {
				return nil, err
			}
			var estimateGasResponse types.EstimateGasResponse
			json.Unmarshal([]byte(estimateGas), &estimateGasResponse)
			xplac.WithMsgType(mevm.EvmInvokeSolContractMsgType)

			gasLimitAdjustment, err := util.GasLimitAdjustment(estimateGasResponse.EstimateGas, xplac.GetGasAdjustment())
			if err != nil {
				return nil, util.LogErr(errors.ErrParse, err)
			}
			gasLimit = gasLimitAdjustment
		}

		return evmTxSignRound(xplac, toAddr, gasPrice, gasLimit, amount, invokeByteData, chainId, ethPrivKey)

	default:
		return nil, util.LogErr(errors.ErrInvalidMsgType, "invalid EVM message type")
	}
}

// Encode transaction by using base64.
func (xplac *XplaClient) EncodeTx(encodeTxMsg types.EncodeTxMsg) (string, error) {
	if xplac.GetErr() != nil {
		return "", xplac.GetErr()
	}
	clientCtx, err := util.NewClient()
	if err != nil {
		return "", err
	}

	tx, err := authclient.ReadTxFromFile(clientCtx, encodeTxMsg.FileName)
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}

	txbytes, err := xplac.GetEncoding().TxConfig.TxEncoder()(tx)
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}

	txbytesBase64 := base64.StdEncoding.EncodeToString(txbytes)

	return txbytesBase64, nil
}

// Decode transaction which encoded by base64
func (xplac *XplaClient) DecodeTx(decodeTxMsg types.DecodeTxMsg) (string, error) {
	if xplac.GetErr() != nil {
		return "", xplac.GetErr()
	}
	txbytes, err := base64.StdEncoding.DecodeString(decodeTxMsg.EncodedByteString)
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}

	tx, err := xplac.GetEncoding().TxConfig.TxDecoder()(txbytes)
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}

	json, err := xplac.GetEncoding().TxConfig.TxJSONEncoder()(tx)
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}

	return string(json), nil
}

// Validate signature
func (xplac *XplaClient) ValidateSignatures(validateSignaturesMsg types.ValidateSignaturesMsg) (string, error) {
	if xplac.GetErr() != nil {
		return "", xplac.GetErr()
	}
	resBool := true
	clientCtx, err := util.NewClient()
	if err != nil {
		return "", err
	}
	stdTx, err := authclient.ReadTxFromFile(clientCtx, validateSignaturesMsg.FileName)
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}

	sigTx := stdTx.(authsigning.SigVerifiableTx)
	signModeHandler := clientCtx.TxConfig.SignModeHandler()

	signers := sigTx.GetSigners()

	sigs, err := sigTx.GetSignaturesV2()
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}

	if len(sigs) != len(signers) {
		resBool = false
	}

	for i, sig := range sigs {
		var (
			PubKey         = sig.PubKey
			multisigHeader string
			multiSigMsg    string
			sigAddr        = sdk.AccAddress(PubKey.Address())
			sigSanity      = "OK"
		)

		if i >= len(signers) || !sigAddr.Equals(signers[i]) {
			sigSanity = "ERROR: signature does not match its respective signer"
			resBool = false
		}

		if !validateSignaturesMsg.Offline && resBool {
			accNum, accSeq, err := clientCtx.AccountRetriever.GetAccountNumberSequence(clientCtx, sigAddr)
			if err != nil {
				return "", util.LogErr(errors.ErrSdkClient, err)
			}

			signingData := authsigning.SignerData{
				ChainID:       validateSignaturesMsg.ChainID,
				AccountNumber: accNum,
				Sequence:      accSeq,
			}
			err = authsigning.VerifySignature(PubKey, signingData, sig.Data, signModeHandler, sigTx)
			if err != nil {
				return "", util.LogErr(errors.ErrParse, err)
			}
		}

		util.LogInfo(i, ":", sigAddr.String(), "[", sigSanity, "]", multisigHeader, multiSigMsg)
	}

	if resBool {
		return "success validate", nil
	} else {
		return "signature validation failed", nil
	}
}
