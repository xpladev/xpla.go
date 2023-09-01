package staking_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	mstaking "github.com/xpladev/xpla.go/core/staking"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	tmjson "github.com/tendermint/tendermint/libs/json"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

var normalAccountNum = 3
var denom = "axpla"
var validatorAddress = "xplavaloper10gv4zj9633v6cje6s2sc0a0xl52hjr6f9jp0q7"

const NodeKey = `{"priv_key":{"type":"tendermint/PrivKeyEd25519","value":"F20DGZKfFFCqgXe2AxF6855KrzfqVasdunk2LMG/EBV+U3gf7GVokgm+X8JP0WG1dyzZ7UddnmC9LGpUMRRQmQ=="}}`
const PrivValidatorKey = `{"address":"3C5042645BAD50A98F0A7D567F862E1A861C23C5","pub_key":{"type":"tendermint/PubKeyEd25519","value":"/0bCEBBwUIrjqYr+pKfzHly+SBMjkA/hcCR9oswxnrk="},"priv_key":{"type":"tendermint/PrivKeyEd25519","value":"iks74YM/Di06VI4JPZ3zOxrKfQ0iwwgXhNa6aIzaduf/RsIQEHBQiuOpiv6kp/MeXL5IEyOQD+FwJH2izDGeuQ=="}}`

// TestSimulateMsgCreateValidator tests the normal scenario of a valid message of type TypeMsgCreateValidator.
// Abonormal scenarios, where the message are created by an errors are not tested here.
func TestSimulateMsgCreateValidator(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, normalAccountNum)
	accounts = getValTargetAccount(t, app, ctx, accounts)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := SimulateMsgCreateValidator(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg stakingtypes.MsgCreateValidator
	stakingtypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	// require.True(t, operationMsg.OK)
	require.Equal(t, "0.110000000000000000", msg.Commission.MaxChangeRate.String())
	require.Equal(t, "0.110000000000000000", msg.Commission.MaxRate.String())
	require.Equal(t, "0.025948405117655148", msg.Commission.Rate.String())
	require.Equal(t, stakingtypes.TypeMsgCreateValidator, msg.Type())
	require.Equal(t, []byte{0xa, 0x20, 0xff, 0x46, 0xc2, 0x10, 0x10, 0x70, 0x50, 0x8a, 0xe3, 0xa9, 0x8a, 0xfe, 0xa4, 0xa7, 0xf3, 0x1e, 0x5c, 0xbe, 0x48, 0x13, 0x23, 0x90, 0xf, 0xe1, 0x70, 0x24, 0x7d, 0xa2, 0xcc, 0x31, 0x9e, 0xb9}, msg.Pubkey.Value)
	require.Equal(t, "xpla10gv4zj9633v6cje6s2sc0a0xl52hjr6f50z40r", msg.DelegatorAddress)
	require.Equal(t, "xplavaloper10gv4zj9633v6cje6s2sc0a0xl52hjr6f9jp0q7", msg.ValidatorAddress)
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgEditValidator tests the normal scenario of a valid message of type TypeMsgEditValidator.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgEditValidator(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// setup accounts[0] as validator
	_ = getTestingValidator0(t, app, ctx, accounts)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := SimulateMsgEditValidator(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg stakingtypes.MsgEditValidator
	stakingtypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "0.057488122528113873", msg.CommissionRate.String())
	require.Equal(t, "wNbeHVIkPZ", msg.Description.Moniker)
	require.Equal(t, "rBqDOTtGTO", msg.Description.Identity)
	require.Equal(t, "XhhuTSkuxK", msg.Description.Website)
	require.Equal(t, "jLxzIivHSl", msg.Description.SecurityContact)
	require.Equal(t, stakingtypes.TypeMsgEditValidator, msg.Type())
	require.Equal(t, "xplavaloper1tnh2q55v8wyygtt9srz5safamzdengsn0rl7kl", msg.ValidatorAddress)
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgDelegate tests the normal scenario of a valid message of type TypeMsgDelegate.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgDelegate(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// setup accounts[0] as validator
	validator0 := getTestingValidator0(t, app, ctx, accounts)
	setupValidatorRewards(app, ctx, validator0.GetOperator())

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := SimulateMsgDelegate(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg stakingtypes.MsgDelegate
	stakingtypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "xpla1ghekyjucln7y67ntx7cf27m9dpuxxemntll57s", msg.DelegatorAddress)
	require.Equal(t, "98100858108421259236", msg.Amount.Amount.String())
	require.Equal(t, "axpla", msg.Amount.Denom)
	require.Equal(t, stakingtypes.TypeMsgDelegate, msg.Type())
	require.Equal(t, "xplavaloper1tnh2q55v8wyygtt9srz5safamzdengsn0rl7kl", msg.ValidatorAddress)
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgUndelegate tests the normal scenario of a valid message of type TypeMsgUndelegate.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgUndelegate(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// setup accounts[0] as validator
	validator0 := getTestingValidator0(t, app, ctx, accounts)

	// setup delegation
	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	validator0, issuedShares := validator0.AddTokensFromDel(delTokens)
	delegator := accounts[1]
	delegation := stakingtypes.NewDelegation(delegator.Address, validator0.GetOperator(), issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)
	app.DistrKeeper.SetDelegatorStartingInfo(ctx, validator0.GetOperator(), delegator.Address, distrtypes.NewDelegatorStartingInfo(2, sdk.OneDec(), 200))

	setupValidatorRewards(app, ctx, validator0.GetOperator())

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := SimulateMsgUndelegate(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg stakingtypes.MsgUndelegate
	stakingtypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "xpla1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7zj7eku", msg.DelegatorAddress)
	require.Equal(t, "280623462081924937", msg.Amount.Amount.String())
	require.Equal(t, "axpla", msg.Amount.Denom)
	require.Equal(t, stakingtypes.TypeMsgUndelegate, msg.Type())
	require.Equal(t, "xplavaloper1tnh2q55v8wyygtt9srz5safamzdengsn0rl7kl", msg.ValidatorAddress)
	require.Len(t, futureOperations, 0)
}

// TestSimulateMsgBeginRedelegate tests the normal scenario of a valid message of type TypeMsgBeginRedelegate.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func TestSimulateMsgBeginRedelegate(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(5)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// setup accounts[0] as validator0 and accounts[1] as validator1
	validator0 := getTestingValidator0(t, app, ctx, accounts)
	validator1 := getTestingValidator1(t, app, ctx, accounts)

	delTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	validator0, issuedShares := validator0.AddTokensFromDel(delTokens)

	// setup accounts[2] as delegator
	delegator := accounts[2]
	delegation := stakingtypes.NewDelegation(delegator.Address, validator1.GetOperator(), issuedShares)
	app.StakingKeeper.SetDelegation(ctx, delegation)
	app.DistrKeeper.SetDelegatorStartingInfo(ctx, validator1.GetOperator(), delegator.Address, distrtypes.NewDelegatorStartingInfo(2, sdk.OneDec(), 200))

	setupValidatorRewards(app, ctx, validator0.GetOperator())
	setupValidatorRewards(app, ctx, validator1.GetOperator())

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := SimulateMsgBeginRedelegate(app.AccountKeeper, app.BankKeeper, app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg stakingtypes.MsgBeginRedelegate
	stakingtypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "xpla12gwd9jchc69wck8dhstxgwz3z8qs8yv6qxgms0", msg.DelegatorAddress)
	require.Equal(t, "489348507626016866", msg.Amount.Amount.String())
	require.Equal(t, "axpla", msg.Amount.Denom)
	require.Equal(t, stakingtypes.TypeMsgBeginRedelegate, msg.Type())
	require.Equal(t, "xplavaloper1h6a7shta7jyc72hyznkys683z98z36e0gre5qc", msg.ValidatorDstAddress)
	require.Equal(t, "xplavaloper17s94pzwhsn4ah25tec27w70n65h5t2sczgd9yh", msg.ValidatorSrcAddress)
	require.Len(t, futureOperations, 0)
}

// SimulateMsgCreateValidator generates a MsgCreateValidator with random values
func SimulateMsgCreateValidator(ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount := accs[normalAccountNum]
		address := sdk.ValAddress(simAccount.Address)

		// ensure the validator doesn't exist already
		_, found := k.GetValidator(ctx, address)
		if found {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, "unable to find validator"), nil, nil
		}

		balance := bk.GetBalance(ctx, simAccount.Address, denom).Amount
		if !balance.IsPositive() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, "balance is negative"), nil, nil
		}

		amount, err := simtypes.RandPositiveInt(r, balance)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, "unable to generate positive amount"), nil, err
		}

		selfDelegation := sdk.NewCoin(denom, amount)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		var fees sdk.Coins

		coins, hasNeg := spendable.SafeSub(sdk.Coins{selfDelegation})
		if !hasNeg {
			fees, err = simtypes.RandomFees(r, ctx, coins)
			if err != nil {
				return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, "unable to generate fees"), nil, err
			}
		}

		maxCommission := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 0, 100)), 2)
		createValidatorMsg := types.CreateValidatorMsg{
			NodeKey:                 NodeKey,
			PrivValidatorKey:        PrivValidatorKey,
			ValidatorAddress:        validatorAddress,
			Moniker:                 simtypes.RandStringOfLength(r, 10),
			Identity:                simtypes.RandStringOfLength(r, 10),
			Website:                 simtypes.RandStringOfLength(r, 10),
			SecurityContact:         simtypes.RandStringOfLength(r, 10),
			Details:                 simtypes.RandStringOfLength(r, 10),
			Amount:                  selfDelegation.String(),
			CommissionRate:          simtypes.RandomDecAmount(r, maxCommission).String(),
			CommissionMaxRate:       maxCommission.String(),
			CommissionMaxChangeRate: simtypes.RandomDecAmount(r, maxCommission).String(),
			MinSelfDelegation:       sdk.OneInt().String(),
		}
		msg, err := mstaking.MakeCreateValidatorMsg(createValidatorMsg, simAccount.PrivKey, "")
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgCreateValidator, "make msg err"), nil, err
		}

		txCtx := simulation.OperationInput{
			App:           app,
			TxGen:         util.MakeEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           msg,
			MsgType:       stakingtypes.TypeMsgCreateValidator,
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			ModuleName:    stakingtypes.ModuleName,
		}

		return testutil.GenAndDeliverTx(txCtx, fees)
	}
}

// SimulateMsgEditValidator generates a MsgEditValidator with random values
func SimulateMsgEditValidator(ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if len(k.GetAllValidators(ctx)) == 0 {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgEditValidator, "number of validators equal zero"), nil, nil
		}

		val, ok := keeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgEditValidator, "unable to pick a validator"), nil, nil
		}

		newCommissionRate := simtypes.RandomDecAmount(r, val.Commission.MaxRate)

		if err := val.Commission.ValidateNewRate(newCommissionRate, ctx.BlockHeader().Time); err != nil {
			// skip as the commission is invalid
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgEditValidator, "invalid commission rate"), nil, nil
		}

		simAccount, found := simtypes.FindAccount(accs, sdk.AccAddress(val.GetOperator()))
		if !found {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgEditValidator, "unable to find account"), nil, fmt.Errorf("validator %s not found", val.GetOperator())
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		maxCommission := sdk.NewDecWithPrec(int64(simtypes.RandIntBetween(r, 0, 100)), 2)
		editValidatorMsg := types.EditValidatorMsg{
			Website:           simtypes.RandStringOfLength(r, 10),
			SecurityContact:   simtypes.RandStringOfLength(r, 10),
			Identity:          simtypes.RandStringOfLength(r, 10),
			Details:           simtypes.RandStringOfLength(r, 10),
			Moniker:           simtypes.RandStringOfLength(r, 10),
			CommissionRate:    simtypes.RandomDecAmount(r, maxCommission).String(),
			MinSelfDelegation: "1000",
		}

		msg, err := mstaking.MakeEditValidatorMsg(editValidatorMsg, simAccount.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgEditValidator, "make msg err"), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           util.MakeEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      stakingtypes.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return testutil.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgDelegate generates a MsgDelegate with random values
func SimulateMsgDelegate(ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if len(k.GetAllValidators(ctx)) == 0 {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgDelegate, "number of validators equal zero"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		val, ok := keeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgDelegate, "unable to pick a validator"), nil, nil
		}

		if val.InvalidExRate() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgDelegate, "validator's invalid echange rate"), nil, nil
		}

		amount := bk.GetBalance(ctx, simAccount.Address, denom).Amount
		if !amount.IsPositive() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgDelegate, "balance is negative"), nil, nil
		}

		amount, err := simtypes.RandPositiveInt(r, amount)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgDelegate, "unable to generate positive amount"), nil, err
		}

		bondAmt := sdk.NewCoin(denom, amount)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		var fees sdk.Coins

		coins, hasNeg := spendable.SafeSub(sdk.Coins{bondAmt})
		if !hasNeg {
			fees, err = simtypes.RandomFees(r, ctx, coins)
			if err != nil {
				return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgDelegate, "unable to generate fees"), nil, err
			}
		}
		delegateMsg := types.DelegateMsg{
			Amount:  bondAmt.String(),
			ValAddr: val.GetOperator().String(),
		}
		msg, err := mstaking.MakeDelegateMsg(delegateMsg, simAccount.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgDelegate, "make msg err"), nil, err
		}

		txCtx := simulation.OperationInput{
			App:           app,
			TxGen:         util.MakeEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           &msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			ModuleName:    stakingtypes.ModuleName,
		}

		return testutil.GenAndDeliverTx(txCtx, fees)
	}
}

// SimulateMsgUndelegate generates a MsgUndelegate with random values
func SimulateMsgUndelegate(ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// get random validator
		validator, ok := keeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "validator is not ok"), nil, nil
		}

		valAddr := validator.GetOperator()
		delegations := k.GetValidatorDelegations(ctx, validator.GetOperator())
		if delegations == nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "keeper does have any delegation entries"), nil, nil
		}

		// get random delegator from validator
		delegation := delegations[r.Intn(len(delegations))]
		delAddr := delegation.GetDelegatorAddr()

		if k.HasMaxUnbondingDelegationEntries(ctx, delAddr, valAddr) {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "keeper does have a max unbonding delegation entries"), nil, nil
		}

		totalBond := validator.TokensFromShares(delegation.GetShares()).TruncateInt()
		if !totalBond.IsPositive() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "total bond is negative"), nil, nil
		}

		unbondAmt, err := simtypes.RandPositiveInt(r, totalBond)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "invalid unbond amount"), nil, err
		}

		if unbondAmt.IsZero() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "unbond amount is zero"), nil, nil
		}

		// need to retrieve the simulation account associated with delegation to retrieve PrivKey
		var simAccount simtypes.Account

		for _, simAcc := range accs {
			if simAcc.Address.Equals(delAddr) {
				simAccount = simAcc
				break
			}
		}
		// if simaccount.PrivKey == nil, delegation address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "account private key is nil"), nil, fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)
		}

		account := ak.GetAccount(ctx, delAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		unbondAmtCoin := sdk.NewCoin(denom, unbondAmt)
		unbondMsg := types.UnbondMsg{
			Amount:  unbondAmtCoin.String(),
			ValAddr: valAddr.String(),
		}
		msg, err := mstaking.MakeUnbondMsg(unbondMsg, simAccount.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgUndelegate, "make msg err"), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           util.MakeEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      stakingtypes.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return testutil.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgBeginRedelegate generates a MsgBeginRedelegate with random values
func SimulateMsgBeginRedelegate(ak stakingtypes.AccountKeeper, bk stakingtypes.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// get random source validator
		srcVal, ok := keeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "unable to pick validator"), nil, nil
		}

		srcAddr := srcVal.GetOperator()
		delegations := k.GetValidatorDelegations(ctx, srcAddr)
		if delegations == nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "keeper does have any delegation entries"), nil, nil
		}

		// get random delegator from src validator
		delegation := delegations[r.Intn(len(delegations))]
		delAddr := delegation.GetDelegatorAddr()

		if k.HasReceivingRedelegation(ctx, delAddr, srcAddr) {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "receveing redelegation is not allowed"), nil, nil // skip
		}

		// get random destination validator
		destVal, ok := keeper.RandomValidator(r, k, ctx)
		if !ok {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "unable to pick validator"), nil, nil
		}

		destAddr := destVal.GetOperator()
		if srcAddr.Equals(destAddr) || destVal.InvalidExRate() || k.HasMaxRedelegationEntries(ctx, delAddr, srcAddr, destAddr) {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "checks failed"), nil, nil
		}

		totalBond := srcVal.TokensFromShares(delegation.GetShares()).TruncateInt()
		if !totalBond.IsPositive() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "total bond is negative"), nil, nil
		}

		redAmt, err := simtypes.RandPositiveInt(r, totalBond)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "unable to generate positive amount"), nil, err
		}

		if redAmt.IsZero() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "amount is zero"), nil, nil
		}

		redAmtCoin := sdk.NewCoin(denom, redAmt)

		// check if the shares truncate to zero
		shares, err := srcVal.SharesFromTokens(redAmt)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "invalid shares"), nil, err
		}

		if srcVal.TokensFromShares(shares).TruncateInt().IsZero() {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "shares truncate to zero"), nil, nil // skip
		}

		// need to retrieve the simulation account associated with delegation to retrieve PrivKey
		var simAccount simtypes.Account

		for _, simAcc := range accs {
			if simAcc.Address.Equals(delAddr) {
				simAccount = simAcc
				break
			}
		}

		// if simaccount.PrivKey == nil, delegation address does not exist in accs. Return error
		if simAccount.PrivKey == nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "account private key is nil"), nil, fmt.Errorf("delegation addr: %s does not exist in simulation accounts", delAddr)
		}

		account := ak.GetAccount(ctx, delAddr)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		redelegateMsg := types.RedelegateMsg{
			Amount:     redAmtCoin.String(),
			ValSrcAddr: srcAddr.String(),
			ValDstAddr: destAddr.String(),
		}
		msg, err := mstaking.MakeRedelegateMsg(redelegateMsg, simAccount.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(stakingtypes.ModuleName, stakingtypes.TypeMsgBeginRedelegate, "make msg err"), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           util.MakeEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      stakingtypes.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return testutil.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*xapp.XplaApp, sdk.Context) {
	app := testutil.Setup(isCheckTx, 5)
	mp := minttypes.Params{
		MintDenom:           denom,
		InflationRateChange: sdk.NewDecWithPrec(13, 2),
		InflationMax:        sdk.NewDecWithPrec(20, 2),
		InflationMin:        sdk.NewDecWithPrec(7, 2),
		GoalBonded:          sdk.NewDecWithPrec(67, 2),
		BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
	}

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	app.MintKeeper.SetParams(ctx, mp)
	app.MintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())
	sp := stakingtypes.DefaultParams()
	sp.BondDenom = denom
	app.StakingKeeper.SetParams(ctx, sp)

	return app, ctx
}

func getTestingAccounts(t *testing.T, r *rand.Rand, app *xapp.XplaApp, ctx sdk.Context, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := app.StakingKeeper.TokensFromConsensusPower(ctx, 200)
	initCoins := sdk.NewCoins(sdk.NewCoin(denom, initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, account.Address)
		app.AccountKeeper.SetAccount(ctx, acc)
		require.NoError(t, testutil.FundAccount(app.BankKeeper, ctx, account.Address, initCoins))
	}

	return accounts
}

func getValTargetAccount(t *testing.T, app *xapp.XplaApp, ctx sdk.Context, accounts []simtypes.Account) []simtypes.Account {
	m1 := "congress warfare much manage hollow birth federal chronic wreck key borrow produce entry inform popular jewel dragon cup list impulse frame unlock genuine want"
	k1 := secp256k1.GenPrivKeyFromSecret([]byte(m1))
	addr := sdk.AccAddress(k1.PubKey().Address())

	consprivKey := tmed25519.GenPrivKeyFromSecret([]byte(m1))

	var pvKey mstaking.FilePVKey

	pvKey.PrivKey = consprivKey
	pvKey.PubKey = consprivKey.PubKey()
	pvKey.Address = consprivKey.PubKey().Address()
	jsonByte, _ := tmjson.Marshal(pvKey)
	require.Equal(t, PrivValidatorKey, string(jsonByte))

	consprivKeyIn := &ed25519.PrivKey{Key: pvKey.PrivKey.Bytes()}

	var specAccount simtypes.Account
	specAccount.PrivKey = k1
	specAccount.PubKey = k1.PubKey()
	specAccount.Address = addr
	specAccount.ConsKey = consprivKeyIn

	initAmt := app.StakingKeeper.TokensFromConsensusPower(ctx, 200)
	initCoins := sdk.NewCoins(sdk.NewCoin(denom, initAmt))

	acc := app.AccountKeeper.NewAccountWithAddress(ctx, specAccount.Address)
	app.AccountKeeper.SetAccount(ctx, acc)
	require.NoError(t, testutil.FundAccount(app.BankKeeper, ctx, specAccount.Address, initCoins))

	accounts = append(accounts, specAccount)

	return accounts
}

func getTestingValidator0(t *testing.T, app *xapp.XplaApp, ctx sdk.Context, accounts []simtypes.Account) stakingtypes.Validator {
	commission0 := stakingtypes.NewCommission(sdk.ZeroDec(), sdk.OneDec(), sdk.OneDec())
	return getTestingValidator(t, app, ctx, accounts, commission0, 0)
}

func getTestingValidator1(t *testing.T, app *xapp.XplaApp, ctx sdk.Context, accounts []simtypes.Account) stakingtypes.Validator {
	commission1 := stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec())
	return getTestingValidator(t, app, ctx, accounts, commission1, 1)
}

func getTestingValidator(t *testing.T, app *xapp.XplaApp, ctx sdk.Context, accounts []simtypes.Account, commission stakingtypes.Commission, n int) stakingtypes.Validator {
	account := accounts[n]
	valPubKey := account.PubKey
	valAddr := sdk.ValAddress(account.PubKey.Address().Bytes())
	validator := teststaking.NewValidator(t, valAddr, valPubKey)
	validator, err := validator.SetInitialCommission(commission)
	require.NoError(t, err)

	validator.DelegatorShares = sdk.NewDec(100)
	validator.Tokens = app.StakingKeeper.TokensFromConsensusPower(ctx, 100)

	app.StakingKeeper.SetValidator(ctx, validator)

	return validator
}

func setupValidatorRewards(app *xapp.XplaApp, ctx sdk.Context, valAddress sdk.ValAddress) {
	decCoins := sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.OneDec())}
	historicalRewards := distrtypes.NewValidatorHistoricalRewards(decCoins, 2)
	app.DistrKeeper.SetValidatorHistoricalRewards(ctx, valAddress, 2, historicalRewards)
	// setup current revards
	currentRewards := distrtypes.NewValidatorCurrentRewards(decCoins, 3)
	app.DistrKeeper.SetValidatorCurrentRewards(ctx, valAddress, currentRewards)

}
