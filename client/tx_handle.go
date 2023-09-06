package client

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"os"

	"github.com/xpladev/xpla.go/controller"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	xauthsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
)

// Set message for transaction builder.
// Interface type messages are converted to correct type.
func setTxBuilderMsg(xplac *XplaClient) (cmclient.TxBuilder, error) {
	if xplac.GetErr() != nil {
		return nil, xplac.GetErr()
	}

	builder := xplac.GetEncoding().TxConfig.NewTxBuilder()

	return controller.Controller().Get(xplac.GetModule()).
		NewTxRouter(builder, xplac.GetMsgType(), xplac.GetMsg())
}

// Set information for transaction builder.
func convertAndSetBuilder(xplac *XplaClient, builder cmclient.TxBuilder, gasLimit string, feeAmount string) (cmclient.TxBuilder, error) {
	feeAmountDenomRemove, err := util.FromStringToBigInt(util.DenomRemove(feeAmount))
	if err != nil {
		return nil, err
	}
	feeAmountCoin := sdk.Coin{
		Amount: sdk.NewIntFromBigInt(feeAmountDenomRemove),
		Denom:  types.XplaDenom,
	}
	feeAmountCoins := sdk.NewCoins(feeAmountCoin)

	if xplac.GetTimeoutHeight() != "" {
		h, err := util.FromStringToUint64(xplac.GetTimeoutHeight())
		if err != nil {
			return nil, err
		}
		builder.SetTimeoutHeight(h)
	}
	if types.Memo != "" {
		builder.SetMemo(types.Memo)
		types.Memo = ""
	}
	gasLimitStr, err := util.FromStringToUint64(gasLimit)
	if err != nil {
		return nil, err
	}

	builder.SetGasLimit(gasLimitStr)
	builder.SetFeeAmount(feeAmountCoins)
	builder.SetFeeGranter(xplac.GetFeeGranter())

	return builder, nil
}

// Sign transaction by using given private key.
func txSignRound(xplac *XplaClient,
	sigsV2 []signing.SignatureV2,
	privs []cryptotypes.PrivKey,
	accSeqs []uint64,
	accNums []uint64,
	builder cmclient.TxBuilder) error {

	for i, priv := range privs {
		sigV2 := signing.SignatureV2{
			PubKey: priv.PubKey(),
			Data: &signing.SingleSignatureData{
				SignMode:  xplac.GetSignMode(),
				Signature: nil,
			},
			Sequence: accSeqs[i],
		}
		sigsV2 = append(sigsV2, sigV2)
	}

	err := builder.SetSignatures(sigsV2...)
	if err != nil {
		return util.LogErr(errors.ErrParse, err)
	}

	sigsV2 = []signing.SignatureV2{}
	for i, priv := range privs {
		signerData := xauthsigning.SignerData{
			ChainID:       xplac.GetChainId(),
			AccountNumber: accNums[i],
			Sequence:      accSeqs[i],
		}
		sigV2, err := tx.SignWithPrivKey(
			xplac.GetSignMode(),
			signerData,
			builder,
			priv,
			xplac.GetEncoding().TxConfig,
			accSeqs[i],
		)
		if err != nil {
			return util.LogErr(errors.ErrParse, err)
		}

		sigsV2 = append(sigsV2, sigV2)
	}

	err = builder.SetSignatures(sigsV2...)
	if err != nil {
		return util.LogErr(errors.ErrParse, err)
	}

	return nil
}

// Sign evm transaction by using given private key.
func evmTxSignRound(xplac *XplaClient,
	toAddr common.Address,
	gasPrice *big.Int,
	gasLimit string,
	amount *big.Int,
	invokeByteData []byte,
	chainId *big.Int,
	ethPrivKey *ecdsa.PrivateKey) ([]byte, error) {

	seqU64, err := util.FromStringToUint64(xplac.GetSequence())
	if err != nil {
		return nil, err
	}
	gasLimitStr, err := util.FromStringToUint64(gasLimit)
	if err != nil {
		return nil, err
	}

	tx := evmtypes.NewTransaction(
		seqU64,
		toAddr,
		amount,
		gasLimitStr,
		gasPrice,
		invokeByteData,
	)

	signer := evmtypes.NewEIP155Signer(chainId)

	signedTx, err := evmtypes.SignTx(tx, signer, ethPrivKey)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}
	txbytes, err := signedTx.MarshalJSON()
	if err != nil {
		return nil, util.LogErr(errors.ErrFailedToMarshal, err)
	}

	return txbytes, nil
}

// Read transaction file and make standard transaction.
func readTxAndInitContexts(clientCtx cmclient.Context, filename string) (cmclient.Context, tx.Factory, sdk.Tx, error) {
	stdTx, err := authclient.ReadTxFromFile(clientCtx, filename)
	if err != nil {
		return clientCtx, tx.Factory{}, nil, util.LogErr(errors.ErrParse, err)
	}

	txFactory := util.NewFactory(clientCtx)

	return clientCtx, txFactory, stdTx, nil
}

// Marshal signature type JSON.
func marshalSignatureJSON(txConfig cmclient.TxConfig, txBldr cmclient.TxBuilder, signatureOnly bool) ([]byte, error) {
	parsedTx := txBldr.GetTx()
	if signatureOnly {
		sigs, err := parsedTx.GetSignaturesV2()
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}
		return txConfig.MarshalSignatureJSON(sigs)
	}

	return txConfig.TxJSONEncoder()(parsedTx)
}

// Unmarshal signature type JSON.
func unmarshalSignatureJSON(clientCtx cmclient.Context, filename string) (sigs []signing.SignatureV2, err error) {
	var bytes []byte
	if bytes, err = os.ReadFile(filename); err != nil {
		return
	}
	return clientCtx.TxConfig.UnmarshalSignatureJSON(bytes)
}

// The secp-256k1 private key converts ECDSA privatkey for using evm module.
func toECDSA(privKey key.PrivateKey) (*ecdsa.PrivateKey, error) {
	return ethcrypto.ToECDSA(privKey.Bytes())
}

// Get multiple signatures information. It returns keyring of cosmos sdk.
func getMultisigInfo(clientCtx cmclient.Context, name string) (keyring.Info, error) {
	kb := clientCtx.Keyring
	multisigInfo, err := kb.Key(name)
	if err != nil {
		return nil, util.LogErr(errors.ErrKeyNotFound, "error getting keybase multisig account", err)
	}
	if multisigInfo.GetType() != keyring.TypeMulti {
		return nil, util.LogErr(errors.ErrInvalidMsgType, name, "must be of type", keyring.TypeMulti, ":", multisigInfo.GetType())
	}

	return multisigInfo, nil
}

// Calculate gas limit and fee amount
func getGasLimitFeeAmount(xplac *XplaClient, builder cmclient.TxBuilder) (string, string, error) {
	gasLimit := xplac.GetGasLimit()
	if xplac.GetGasLimit() == "" {
		if xplac.GetLcdURL() == "" && xplac.GetGrpcUrl() == "" {
			gasLimit = types.DefaultGasLimit
		} else {
			simulate, err := xplac.Simulate(builder)
			if err != nil {
				return "", "", err
			}
			gasLimitAdjustment, err := util.GasLimitAdjustment(simulate.GasInfo.GasUsed, xplac.GetGasAdjustment())
			if err != nil {
				return "", "", err
			}
			gasLimit = gasLimitAdjustment
		}
	}

	feeAmount := xplac.GetFeeAmount()
	if xplac.GetFeeAmount() == "" {
		gasLimitBigInt, err := util.FromStringToBigInt(gasLimit)
		if err != nil {
			return "", "", err
		}

		gasPriceBigInt, err := util.FromStringToBigInt(xplac.GetGasPrice())
		if err != nil {
			return "", "", err
		}

		feeAmountBigInt := util.MulBigInt(gasLimitBigInt, gasPriceBigInt)
		feeAmount = util.FromBigIntToString(feeAmountBigInt)
	}

	return gasLimit, feeAmount, nil
}

// check user = signer
func isTxSigner(user sdk.AccAddress, signers []sdk.AccAddress) bool {
	for _, s := range signers {
		if bytes.Equal(user.Bytes(), s.Bytes()) {
			return true
		}
	}

	return false
}

// Get account number and sequence
func GetAccNumAndSeq(xplac *XplaClient) (*XplaClient, error) {
	if xplac.GetAccountNumber() == "" || xplac.GetSequence() == "" {
		if xplac.GetLcdURL() == "" && xplac.GetGrpcUrl() == "" {
			xplac.WithAccountNumber(util.FromUint64ToString(types.DefaultAccNum))
			xplac.WithSequence(util.FromUint64ToString(types.DefaultAccSeq))
		} else {
			account, err := xplac.LoadAccount(sdk.AccAddress(xplac.GetPrivateKey().PubKey().Address()))
			if err != nil {
				return nil, err
			}
			xplac.WithAccountNumber(util.FromUint64ToString(account.GetAccountNumber()))
			xplac.WithSequence(util.FromUint64ToString(account.GetSequence()))
		}
	}
	return xplac, nil
}
