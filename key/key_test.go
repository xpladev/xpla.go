package key

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMnemonic(t *testing.T) {
	_, err := NewMnemonic()
	assert.NoError(t, err)
}

func TestNewPrivKey(t *testing.T) {
	mnemonic, err := NewMnemonic()
	assert.NoError(t, err)

	// Only Secp256k1 is supported
	_, err = NewPrivKey(mnemonic)
	assert.NoError(t, err)
}

func TestConvertPrivKeyToBech32Addr(t *testing.T) {
	mnemonic, err := NewMnemonic()
	assert.NoError(t, err)

	PrivateKey, err := NewPrivKey(mnemonic)
	assert.NoError(t, err)

	_, err = Bech32AddrString(PrivateKey)
	assert.NoError(t, err)
}

func TestConvertPrivKeyToHexAddr(t *testing.T) {
	mnemonic, err := NewMnemonic()
	assert.NoError(t, err)

	PrivateKey, err := NewPrivKey(mnemonic)
	assert.NoError(t, err)

	addrMyKey := PrivateKey.PubKey().Address().String()

	addr := HexAddrString(PrivateKey)
	require.Equal(t, addrMyKey, addr)
}

func TestEncryptDecryptPrivKeyArmor(t *testing.T) {
	mnemonic, err := NewMnemonic()
	assert.NoError(t, err)

	PrivateKey, err := NewPrivKey(mnemonic)
	assert.NoError(t, err)

	armor1 := EncryptArmorPrivKey(PrivateKey, DefaultEncryptPassphrase)
	armor2 := EncryptArmorPrivKeyWithoutPassphrase(PrivateKey)

	pk1, algo1, err := UnarmorDecryptPrivKey(armor1, DefaultEncryptPassphrase)
	assert.NoError(t, err)
	pk2, algo2, err := UnarmorDecryptPrivKeyWithoutPassphrase(armor2)
	assert.NoError(t, err)

	require.Equal(t, pk1, pk2)
	require.Equal(t, algo1, algo2)
}
