package key

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec/legacy"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/tendermint/crypto/bcrypt"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/armor"
	"github.com/tendermint/tendermint/crypto/xsalsa20symmetric"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

const (
	blockTypePrivKey = "TENDERMINT PRIVATE KEY"
	blockTypeKeyInfo = "TENDERMINT KEY INFO"
	blockTypePubKey  = "TENDERMINT PUBLIC KEY"
	headerVersion    = "version"
	headerType       = "type"
)

var BcryptSecurityParameter = 12
var DefaultEncryptPassphrase = "xplaDefaultPassphrase"

// Encrypt secp-256k1 private key to make armored key.
func EncryptArmorPrivKey(privKey cryptotypes.PrivKey, passphrase string) string {
	return encryptPrivKey(privKey, passphrase)
}

// Encrypt secp-256k1 private key to make armored key without passphrase for keyring. (for test)
func EncryptArmorPrivKeyWithoutPassphrase(privKey cryptotypes.PrivKey) string {
	return encryptPrivKey(privKey, DefaultEncryptPassphrase)
}

// Encrypt private key.
func encryptPrivKey(privKey cryptotypes.PrivKey, passphrase string) string {
	saltBytes := crypto.CRandBytes(16)
	key, err := bcrypt.GenerateFromPassword(saltBytes, []byte(passphrase), BcryptSecurityParameter)
	if err != nil {
		panic(util.LogErr(errors.ErrInvalidRequest, "error generating bcrypt key from passphrase", err))
	}

	key = crypto.Sha256(key) // get 32 bytes
	encodingConfig := util.MakeEncodingConfig()
	privKeyBytes := encodingConfig.Amino.MustMarshal(privKey)

	encBytes := xsalsa20symmetric.EncryptSymmetric(privKeyBytes, key)

	header := map[string]string{
		"kdf":  "bcrypt",
		"salt": fmt.Sprintf("%X", saltBytes),
	}

	header[headerType] = types.DefaultXplaKeyAlgo

	return armor.EncodeArmor(blockTypePrivKey, header, encBytes)
}

// Decrypt armored private key.
func UnarmorDecryptPrivKey(armorStr string, passphrase string) (privKey cryptotypes.PrivKey, algo string, err error) {
	saltBytes, encBytes, header, err := decodingArmoredPrivKey(armorStr)
	if err != nil {
		return privKey, "", err
	}

	privKey, err = decryptPrivKey(saltBytes, encBytes, passphrase)

	if header[headerType] == "" {
		header[headerType] = types.DefaultXplaKeyAlgo
	}

	return privKey, header[headerType], err
}

// Decrypt armored private key without passpharse for keyring. (for test)
func UnarmorDecryptPrivKeyWithoutPassphrase(armorStr string) (privKey cryptotypes.PrivKey, algo string, err error) {
	saltBytes, encBytes, header, err := decodingArmoredPrivKey(armorStr)
	if err != nil {
		return privKey, "", err
	}

	privKey, err = decryptPrivKey(saltBytes, encBytes, DefaultEncryptPassphrase)
	if header[headerType] == "" {
		header[headerType] = types.DefaultXplaKeyAlgo
	}

	return privKey, header[headerType], err
}

// Decode armored private key.
func decodingArmoredPrivKey(armorStr string) ([]byte, []byte, map[string]string, error) {
	blockType, header, encBytes, err := armor.DecodeArmor(armorStr)
	if err != nil {
		return nil, nil, nil, err
	}

	if blockType != blockTypePrivKey {
		return nil, nil, nil, util.LogErr(errors.ErrInvalidRequest, "unrecognized armor type: ", blockType)
	}

	if header["kdf"] != "bcrypt" {
		return nil, nil, nil, util.LogErr(errors.ErrInvalidRequest, "unrecognized KDF type: ", header["kdf"])
	}

	if header["salt"] == "" {
		return nil, nil, nil, util.LogErr(errors.ErrInvalidRequest, "missing salt bytes")
	}

	saltBytes, err := hex.DecodeString(header["salt"])
	if err != nil {
		return nil, nil, nil, util.LogErr(errors.ErrInvalidRequest, "error decoding salt: %v", err.Error())
	}

	return saltBytes, encBytes, header, nil
}

// Decrypt private key.
func decryptPrivKey(saltBytes []byte, encBytes []byte, passphrase string) (privKey cryptotypes.PrivKey, err error) {
	key, err := bcrypt.GenerateFromPassword(saltBytes, []byte(passphrase), BcryptSecurityParameter)
	if err != nil {
		return privKey, util.LogErr(errors.ErrInvalidRequest, err, "error generating bcrypt key from passphrase")
	}

	key = crypto.Sha256(key) // Get 32 bytes

	privKeyBytes, err := xsalsa20symmetric.DecryptSymmetric(encBytes, key)
	if err != nil && err.Error() == "Ciphertext decryption failed" {
		return privKey, util.LogErr(errors.ErrInvalidRequest, "invalid account password")
	} else if err != nil {
		return privKey, util.LogErr(errors.ErrParse, err)
	}

	return legacy.PrivKeyFromBytes(privKeyBytes)
}
