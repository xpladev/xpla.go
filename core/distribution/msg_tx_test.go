package distribution_test

import (
	"fmt"
	"math/rand"
	"testing"

	mdist "github.com/xpladev/xpla.go/core/distribution"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type TestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *xapp.XplaApp
	queryClient disttypes.QueryClient
}

func (suite *TestSuite) SetupTest() {
	checkTx := false
	app := testutil.Setup(checkTx, 5)
	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	disttypes.RegisterQueryServer(queryHelper, app.DistrKeeper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = disttypes.NewQueryClient(queryHelper)
}

// TestSimulateMsgSetWithdrawAddress tests the normal scenario of a valid message of type TypeMsgSetWithdrawAddress.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func (suite *TestSuite) TestSimulateMsgSetWithdrawAddress() {

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: suite.app.LastBlockHeight() + 1, AppHash: suite.app.LastCommitID().Hash}})

	// execute operation
	op := SimulateMsgSetWithdrawAddress(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.DistrKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, testutil.TestChainId)
	suite.Require().NoError(err)

	var msg disttypes.MsgSetWithdrawAddress
	disttypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal("xpla1ghekyjucln7y67ntx7cf27m9dpuxxemntll57s", msg.DelegatorAddress)
	suite.Require().Equal("xpla1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7zj7eku", msg.WithdrawAddress)
	suite.Require().Equal(disttypes.TypeMsgSetWithdrawAddress, msg.Type())
	suite.Require().Equal(disttypes.ModuleName, msg.Route())
	suite.Require().Len(futureOperations, 0)
}

// TestSimulateMsgWithdrawDelegatorReward tests the normal scenario of a valid message
// of type TypeMsgWithdrawDelegatorReward.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func (suite *TestSuite) TestSimulateMsgWithdrawDelegatorReward() {
	// setup 3 accounts
	s := rand.NewSource(4)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// setup accounts[0] as validator
	validator0 := suite.getTestingValidator(accounts)

	// setup delegation
	delTokens := suite.app.StakingKeeper.TokensFromConsensusPower(suite.ctx, 2)
	validator0, issuedShares := validator0.AddTokensFromDel(delTokens)
	delegator := accounts[1]
	delegation := stakingtypes.NewDelegation(delegator.Address, validator0.GetOperator(), issuedShares)
	suite.app.StakingKeeper.SetDelegation(suite.ctx, delegation)
	suite.app.DistrKeeper.SetDelegatorStartingInfo(suite.ctx, validator0.GetOperator(), delegator.Address, disttypes.NewDelegatorStartingInfo(2, sdk.OneDec(), 200))

	suite.setupValidatorRewards(validator0.GetOperator())

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: suite.app.LastBlockHeight() + 1, AppHash: suite.app.LastCommitID().Hash}})

	// execute operation
	op := SimulateMsgWithdrawDelegatorReward(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.DistrKeeper, suite.app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, testutil.TestChainId)
	suite.Require().NoError(err)

	var msg disttypes.MsgWithdrawDelegatorReward
	disttypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal("xplavaloper1l4s054098kk9hmr5753c6k3m2kw65h68sr7gl7", msg.ValidatorAddress)
	suite.Require().Equal("xpla1d6u7zhjwmsucs678d7qn95uqajd4ucl9vl2hpf", msg.DelegatorAddress)
	suite.Require().Equal(disttypes.TypeMsgWithdrawDelegatorReward, msg.Type())
	suite.Require().Equal(disttypes.ModuleName, msg.Route())
	suite.Require().Len(futureOperations, 0)
}

// TestSimulateMsgFundCommunityPool tests the normal scenario of a valid message of type TypeMsgFundCommunityPool.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func (suite *TestSuite) TestSimulateMsgFundCommunityPool() {
	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: suite.app.LastBlockHeight() + 1, AppHash: suite.app.LastCommitID().Hash}})

	// execute operation
	op := SimulateMsgFundCommunityPool(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.DistrKeeper, suite.app.StakingKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, testutil.TestChainId)
	suite.Require().NoError(err)

	var msg disttypes.MsgFundCommunityPool
	disttypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal("677063145828421287964axpla", msg.Amount.String())
	suite.Require().Equal("xpla1ghekyjucln7y67ntx7cf27m9dpuxxemntll57s", msg.Depositor)
	suite.Require().Equal(disttypes.TypeMsgFundCommunityPool, msg.Type())
	suite.Require().Equal(disttypes.ModuleName, msg.Route())
	suite.Require().Len(futureOperations, 0)
}

// SimulateMsgSetWithdrawAddress generates a MsgSetWithdrawAddress with random values.
func SimulateMsgSetWithdrawAddress(ak disttypes.AccountKeeper, bk disttypes.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if !k.GetWithdrawAddrEnabled(ctx) {
			return simtypes.NoOpMsg(disttypes.ModuleName, disttypes.TypeMsgSetWithdrawAddress, "withdrawal is not enabled"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		simToAccount, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		setwithdrawAddrMsg := types.SetWithdrawAddrMsg{
			WithdrawAddr: simToAccount.Address.String(),
		}
		msg, err := mdist.MakeSetWithdrawAddrMsg(setwithdrawAddrMsg, simAccount.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(disttypes.ModuleName, disttypes.TypeMsgSetWithdrawAddress, "make msg err"), nil, err
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
			ModuleName:      disttypes.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return testutil.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgWithdrawDelegatorReward generates a MsgWithdrawDelegatorReward with random values.
func SimulateMsgWithdrawDelegatorReward(ak disttypes.AccountKeeper, bk disttypes.BankKeeper, k keeper.Keeper, sk stakingkeeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		delegations := sk.GetAllDelegatorDelegations(ctx, simAccount.Address)
		if len(delegations) == 0 {
			return simtypes.NoOpMsg(disttypes.ModuleName, disttypes.TypeMsgWithdrawDelegatorReward, "number of delegators equal 0"), nil, nil
		}

		delegation := delegations[r.Intn(len(delegations))]

		validator := sk.Validator(ctx, delegation.GetValidatorAddr())
		if validator == nil {
			return simtypes.NoOpMsg(disttypes.ModuleName, disttypes.TypeMsgWithdrawDelegatorReward, "validator is nil"), nil, fmt.Errorf("validator %s not found", delegation.GetValidatorAddr())
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		// msg := types.NewMsgWithdrawDelegatorReward(simAccount.Address, validator.GetOperator())

		withdrawRewardsMsg := types.WithdrawRewardsMsg{
			DelegatorAddr: simAccount.Address.String(),
			ValidatorAddr: validator.GetOperator().String(),
			Commission:    true,
		}

		msg, err := mdist.MakeWithdrawRewardsMsg(withdrawRewardsMsg, simAccount.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(disttypes.ModuleName, disttypes.TypeMsgWithdrawDelegatorReward, "make msg err"), nil, err
		}

		targetMsg := msg[0]
		targetMsgType := disttypes.TypeMsgWithdrawDelegatorReward
		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           util.MakeEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             targetMsg,
			MsgType:         targetMsgType,
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      disttypes.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return testutil.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgFundCommunityPool simulates MsgFundCommunityPool execution where
// a random account sends a random amount of its funds to the community pool.
func SimulateMsgFundCommunityPool(ak disttypes.AccountKeeper, bk disttypes.BankKeeper, k keeper.Keeper, sk stakingkeeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		funder, _ := simtypes.RandomAcc(r, accs)

		account := ak.GetAccount(ctx, funder.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		fundAmount := simtypes.RandSubsetCoins(r, spendable)
		if fundAmount.Empty() {
			return simtypes.NoOpMsg(disttypes.ModuleName, disttypes.TypeMsgFundCommunityPool, "fund amount is empty"), nil, nil
		}

		var (
			fees sdk.Coins
			err  error
		)

		coins, hasNeg := spendable.SafeSub(fundAmount)
		if !hasNeg {
			fees, err = simtypes.RandomFees(r, ctx, coins)
			if err != nil {
				return simtypes.NoOpMsg(disttypes.ModuleName, disttypes.TypeMsgFundCommunityPool, "unable to generate fees"), nil, err
			}
		}

		fundCommunityPoolMsg := types.FundCommunityPoolMsg{
			Amount: coins.String(),
		}

		msg, err := mdist.MakeFundCommunityPoolMsg(fundCommunityPoolMsg, funder.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(disttypes.ModuleName, disttypes.TypeMsgFundCommunityPool, "make msg err"), nil, err
		}

		txCtx := simulation.OperationInput{
			App:           app,
			TxGen:         util.MakeEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           &msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    funder,
			AccountKeeper: ak,
			ModuleName:    disttypes.ModuleName,
		}

		return testutil.GenAndDeliverTx(txCtx, fees)
	}
}

func (suite *TestSuite) getTestingAccounts(r *rand.Rand, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := suite.app.StakingKeeper.TokensFromConsensusPower(suite.ctx, 200000)
	initCoins := sdk.NewCoins(sdk.NewCoin("axpla", initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, account.Address)
		suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
		suite.Require().NoError(testutil.FundAccount(suite.app.BankKeeper, suite.ctx, account.Address, initCoins))
	}

	return accounts
}

func (suite *TestSuite) getTestingValidator(accounts []simtypes.Account) stakingtypes.Validator {
	n := 0
	commission := stakingtypes.NewCommission(sdk.ZeroDec(), sdk.OneDec(), sdk.OneDec())

	require := suite.Require()
	account := accounts[n]
	valPubKey := account.PubKey
	valAddr := sdk.ValAddress(account.PubKey.Address().Bytes())
	validator, err := stakingtypes.NewValidator(valAddr, valPubKey, stakingtypes.
		Description{})
	require.NoError(err)
	validator, err = validator.SetInitialCommission(commission)
	require.NoError(err)
	validator.DelegatorShares = sdk.NewDec(100)
	validator.Tokens = sdk.NewInt(1000000)

	suite.app.StakingKeeper.SetValidator(suite.ctx, validator)

	return validator
}

func (suite *TestSuite) setupValidatorRewards(valAddress sdk.ValAddress) {
	decCoins := sdk.DecCoins{sdk.NewDecCoinFromDec(sdk.DefaultBondDenom, sdk.OneDec())}
	historicalRewards := disttypes.NewValidatorHistoricalRewards(decCoins, 2)
	suite.app.DistrKeeper.SetValidatorHistoricalRewards(suite.ctx, valAddress, 2, historicalRewards)
	// setup current revards
	currentRewards := disttypes.NewValidatorCurrentRewards(decCoins, 3)
	suite.app.DistrKeeper.SetValidatorCurrentRewards(suite.ctx, valAddress, currentRewards)

}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
