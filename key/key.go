package key

import (
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bip39 "github.com/cosmos/go-bip39"
	evmhd "github.com/evmos/ethermint/crypto/hd"
)

type PrivateKey = cryptotypes.PrivKey

// Make new mnemonic words by using bip39 entropy.
// Mnemonic words are changed every time user run new mnemonic function.
func NewMnemonic() (string, error) {
	// Default number of words (24): This generates a mnemonic directly from the
	// number of words by reading system entropy.
	entropy, err := bip39.NewEntropy(types.DefaultEntropySize)
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}

	return mnemonic, nil
}

// Make new private key.
// The private key generation algorithm uses eth-secp256k1 to use the evm module.
func NewPrivKey(mnemonic string) (cryptotypes.PrivKey, error) {
	algo := evmhd.EthSecp256k1
	derivedPri, err := algo.Derive()(mnemonic, keyring.DefaultBIP39Passphrase, sdk.GetConfig().GetFullBIP44Path())
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	privateKey := algo.Generate()(derivedPri)

	return privateKey, nil
}

// Convert private key to bech32 address.
func Bech32AddrString(p PrivateKey) (string, error) {
	addr, err := sdk.AccAddressFromHex(p.PubKey().Address().String())
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}

	return addr.String(), nil
}

// Convert private key to hex address.
func HexAddrString(p PrivateKey) string {
	return p.PubKey().Address().String()
}
