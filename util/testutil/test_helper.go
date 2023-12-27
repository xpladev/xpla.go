package testutil

import (
	"encoding/json"

	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/go-bip39"
	evmhd "github.com/evmos/ethermint/crypto/hd"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	dbm "github.com/tendermint/tm-db"
	xapp "github.com/xpladev/xpla/app"
	"github.com/xpladev/xpla/app/helpers"
)

const (
	DefaultTestGenTxGas = 1000000
	TestChainId         = "cube_47-5"
)

func Setup(isCheckTx bool, invCheckPeriod uint) *xapp.XplaApp {
	app, genesisState := setup(!isCheckTx, invCheckPeriod)
	if !isCheckTx {
		// InitChain must be called to stop deliverState from being nil
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ChainId:         TestChainId,
				ConsensusParams: helpers.DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

func setup(withGenesis bool, invCheckPeriod uint) (*xapp.XplaApp, xapp.GenesisState) {
	db := dbm.NewMemDB()
	encCdc := xapp.MakeTestEncodingConfig()
	app := xapp.NewXplaApp(
		log.NewNopLogger(),
		db,
		nil,
		true,
		map[int64]bool{},
		xapp.DefaultNodeHome,
		invCheckPeriod,
		encCdc,
		xapp.GetEnabledProposals(),
		helpers.EmptyAppOptions{},
		[]wasm.Option{},
	)
	if withGenesis {
		return app, xapp.NewDefaultGenesisState()
	}

	return app, xapp.GenesisState{}
}

// FundAccount is a utility function that funds an account by minting and
// sending the coins to the address. This should be used for testing purposes
// only!
//
// TODO: Instead of using the mint module account, which has the
// permission of minting, create a "faucet" account. (@fdymylja)
func FundAccount(bankKeeper bankkeeper.Keeper, ctx sdk.Context, addr sdk.AccAddress, amounts sdk.Coins) error {
	if err := bankKeeper.MintCoins(ctx, minttypes.ModuleName, amounts); err != nil {
		return err
	}

	return bankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, addr, amounts)
}

func NewTestMnemonic(entropy []byte) (string, error) {
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}

	return mnemonic, nil
}

func NewTestEthSecpPrivKey(mnemonic string) (cryptotypes.PrivKey, error) {
	algo := evmhd.EthSecp256k1
	derivedPri, err := algo.Derive()(mnemonic, keyring.DefaultBIP39Passphrase, sdk.GetConfig().GetFullBIP44Path())
	if err != nil {
		return nil, err
	}

	privateKey := algo.Generate()(derivedPri)

	return privateKey, nil
}

func NewTestSecpPrivKey(mnemonic string) (cryptotypes.PrivKey, error) {
	algo := hd.Secp256k1
	derivedPri, err := algo.Derive()(mnemonic, keyring.DefaultBIP39Passphrase, sdk.GetConfig().GetFullBIP44Path())
	if err != nil {
		return nil, err
	}

	privateKey := algo.Generate()(derivedPri)

	return privateKey, nil
}
