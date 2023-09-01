package testutil

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
)

// RandomAccounts generates n random accounts
func RandomAccounts(r *rand.Rand, n int) []simtypes.Account {
	accs := make([]simtypes.Account, n)

	for i := 0; i < n; i++ {
		// don't need that much entropy for simulation
		privkeySeed := make([]byte, 16)
		r.Read(privkeySeed)

		mnemonic, _ := NewTestMnemonic(privkeySeed)
		ethsecpPrivKey, _ := NewTestEthSecpPrivKey(mnemonic)

		accs[i].PrivKey = ethsecpPrivKey
		accs[i].PubKey = accs[i].PrivKey.PubKey()
		accs[i].Address = sdk.AccAddress(accs[i].PubKey.Address())

		accs[i].ConsKey = ed25519.GenPrivKeyFromSecret(privkeySeed)
	}

	return accs
}

func RandomSecp256k1Accounts(r *rand.Rand, n int) []simtypes.Account {
	accs := make([]simtypes.Account, n)

	for i := 0; i < n; i++ {
		// don't need that much entropy for simulation
		privkeySeed := make([]byte, 16)
		r.Read(privkeySeed)

		mnemonic, _ := NewTestMnemonic(privkeySeed)
		secpPrivKey, _ := NewTestSecpPrivKey(mnemonic)

		accs[i].PrivKey = secpPrivKey
		accs[i].PubKey = accs[i].PrivKey.PubKey()
		accs[i].Address = sdk.AccAddress(accs[i].PubKey.Address())

		accs[i].ConsKey = ed25519.GenPrivKeyFromSecret(privkeySeed)
	}

	return accs
}
