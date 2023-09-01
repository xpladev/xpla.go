package slashing_test

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	mslashing "github.com/xpladev/xpla.go/core/slashing"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

// TestSimulateMsgUnjail tests the normal scenario of a valid message of type types.MsgUnjail.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func TestSimulateMsgUnjail(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// setup accounts[0] as validator0
	validator0 := getTestingValidator0(t, app, ctx, accounts)

	// setup validator0 by consensus address
	app.StakingKeeper.SetValidatorByConsAddr(ctx, validator0)
	val0ConsAddress, err := validator0.GetConsAddr()
	require.NoError(t, err)
	info := slashingtypes.NewValidatorSigningInfo(val0ConsAddress, int64(4), int64(3),
		time.Unix(2, 0), false, int64(10))
	app.SlashingKeeper.SetValidatorSigningInfo(ctx, val0ConsAddress, info)

	// put validator0 in jail
	app.StakingKeeper.Jail(ctx, val0ConsAddress)

	// setup self delegation
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	validator0, issuedShares := validator0.AddTokensFromDel(delTokens)
	val0AccAddress, err := sdk.ValAddressFromBech32(validator0.OperatorAddress)
	require.NoError(t, err)
	selfDelegation := stakingtypes.NewDelegation(val0AccAddress.Bytes(), validator0.GetOperator(), issuedShares)
	app.StakingKeeper.SetDelegation(ctx, selfDelegation)
	app.DistrKeeper.SetDelegatorStartingInfo(ctx, validator0.GetOperator(), val0AccAddress.Bytes(), distrtypes.NewDelegatorStartingInfo(2, sdk.OneDec(), 200))

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := SimulateMsgUnjail(app.AccountKeeper, app.BankKeeper, app.SlashingKeeper, app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, testutil.TestChainId)
	require.NoError(t, err)

	var msg slashingtypes.MsgUnjail
	slashingtypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, slashingtypes.TypeMsgUnjail, msg.Type())
	require.Equal(t, "xplavaloper1tnh2q55v8wyygtt9srz5safamzdengsn0rl7kl", msg.ValidatorAddr)
	require.Len(t, futureOperations, 0)
}

// SimulateMsgUnjail generates a MsgUnjail with random values
func SimulateMsgUnjail(ak slashingtypes.AccountKeeper, bk slashingtypes.BankKeeper, k slashingkeeper.Keeper, sk stakingkeeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		validator, ok := stakingkeeper.RandomValidator(r, sk, ctx)
		if !ok {
			return simtypes.NoOpMsg(slashingtypes.ModuleName, slashingtypes.TypeMsgUnjail, "validator is not ok"), nil, nil // skip
		}

		simAccount, found := simtypes.FindAccount(accs, sdk.AccAddress(validator.GetOperator()))
		if !found {
			return simtypes.NoOpMsg(slashingtypes.ModuleName, slashingtypes.TypeMsgUnjail, "unable to find account"), nil, nil // skip
		}

		if !validator.IsJailed() {
			// TODO: due to this condition this message is almost, if not always, skipped !
			return simtypes.NoOpMsg(slashingtypes.ModuleName, slashingtypes.TypeMsgUnjail, "validator is not jailed"), nil, nil
		}

		consAddr, err := validator.GetConsAddr()
		if err != nil {
			return simtypes.NoOpMsg(slashingtypes.ModuleName, slashingtypes.TypeMsgUnjail, "unable to get validator consensus key"), nil, err
		}
		info, found := k.GetValidatorSigningInfo(ctx, consAddr)
		if !found {
			return simtypes.NoOpMsg(slashingtypes.ModuleName, slashingtypes.TypeMsgUnjail, "unable to find validator signing info"), nil, nil // skip
		}

		selfDel := sk.Delegation(ctx, simAccount.Address, validator.GetOperator())
		if selfDel == nil {
			return simtypes.NoOpMsg(slashingtypes.ModuleName, slashingtypes.TypeMsgUnjail, "self delegation is nil"), nil, nil // skip
		}

		account := ak.GetAccount(ctx, sdk.AccAddress(validator.GetOperator()))
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fees, err := simtypes.RandomFees(r, ctx, spendable)
		if err != nil {
			return simtypes.NoOpMsg(slashingtypes.ModuleName, slashingtypes.TypeMsgUnjail, "unable to generate fees"), nil, err
		}

		msg, err := mslashing.MakeUnjailMsg(simAccount.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(slashingtypes.ModuleName, slashingtypes.TypeMsgUnjail, "make msg err"), nil, err
		}

		txGen := util.MakeEncodingConfig().TxConfig
		tx, err := testutil.GenTx(
			txGen,
			[]sdk.Msg{&msg},
			fees,
			testutil.DefaultTestGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(slashingtypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, res, err := app.Deliver(txGen.TxEncoder(), tx)

		// result should fail if:
		// - validator cannot be unjailed due to tombstone
		// - validator is still in jailed period
		// - self delegation too low
		if info.Tombstoned ||
			ctx.BlockHeader().Time.Before(info.JailedUntil) ||
			validator.TokensFromShares(selfDel.GetShares()).TruncateInt().LT(validator.GetMinSelfDelegation()) {
			if res != nil && err == nil {
				if info.Tombstoned {
					return simtypes.NewOperationMsg(&msg, true, "", nil), nil, errors.New("validator should not have been unjailed if validator tombstoned")
				}
				if ctx.BlockHeader().Time.Before(info.JailedUntil) {
					return simtypes.NewOperationMsg(&msg, true, "", nil), nil, errors.New("validator unjailed while validator still in jail period")
				}
				if validator.TokensFromShares(selfDel.GetShares()).TruncateInt().LT(validator.GetMinSelfDelegation()) {
					return simtypes.NewOperationMsg(&msg, true, "", nil), nil, errors.New("validator unjailed even though self-delegation too low")
				}
			}
			// msg failed as expected
			return simtypes.NewOperationMsg(&msg, false, "", nil), nil, nil
		}

		if err != nil {
			return simtypes.NoOpMsg(slashingtypes.ModuleName, msg.Type(), "unable to deliver tx"), nil, errors.New(res.Log)
		}

		return simtypes.NewOperationMsg(&msg, true, "", nil), nil, nil
	}
}

func getTestingAccounts(t *testing.T, r *rand.Rand, app *xapp.XplaApp, ctx sdk.Context, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := app.StakingKeeper.TokensFromConsensusPower(ctx, 200000)
	initCoins := sdk.NewCoins(sdk.NewCoin("axpla", initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, account.Address)
		app.AccountKeeper.SetAccount(ctx, acc)
		require.NoError(t, testutil.FundAccount(app.BankKeeper, ctx, account.Address, initCoins))
	}

	return accounts
}

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*xapp.XplaApp, sdk.Context) {
	app := testutil.Setup(isCheckTx, 5)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, minttypes.DefaultParams())
	app.MintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())

	return app, ctx
}

func getTestingValidator0(t *testing.T, app *xapp.XplaApp, ctx sdk.Context, accounts []simtypes.Account) stakingtypes.Validator {
	commission0 := stakingtypes.NewCommission(sdk.ZeroDec(), sdk.OneDec(), sdk.OneDec())
	return getTestingValidator(t, app, ctx, accounts, commission0, 0)
}

func getTestingValidator(t *testing.T, app *xapp.XplaApp, ctx sdk.Context, accounts []simtypes.Account, commission stakingtypes.Commission, n int) stakingtypes.Validator {
	account := accounts[n]
	valPubKey := account.ConsKey.PubKey()
	valAddr := sdk.ValAddress(account.PubKey.Address().Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, valPubKey, stakingtypes.Description{})
	require.NoError(t, err)
	validator, err = validator.SetInitialCommission(commission)
	require.NoError(t, err)

	validator.DelegatorShares = sdk.NewDec(100)
	validator.Tokens = sdk.NewInt(1000000)

	app.StakingKeeper.SetValidator(ctx, validator)

	return validator
}
