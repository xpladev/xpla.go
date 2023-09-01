package bank_test

import (
	"math/rand"
	"testing"

	mbank "github.com/xpladev/xpla.go/core/bank"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type TestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *xapp.XplaApp
	queryClient banktypes.QueryClient
}

func (suite *TestSuite) SetupTest() {
	checkTx := false
	app := testutil.Setup(checkTx, 5)
	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	banktypes.RegisterQueryServer(queryHelper, app.BankKeeper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = banktypes.NewQueryClient(queryHelper)
}

// TestSimulateMsgSend tests the normal scenario of a valid message of type TypeMsgSend.
// Abonormal scenarios, where the message is created by an errors, are not tested here.
func (suite *TestSuite) TestSimulateMsgSend() {
	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: suite.app.LastBlockHeight() + 1, AppHash: suite.app.LastCommitID().Hash}})

	// execute operation
	op := SimulateMsgSend(suite.app.AccountKeeper, suite.app.BankKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, testutil.TestChainId)
	suite.Require().NoError(err)

	var msg banktypes.MsgSend
	banktypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal("107242728498552699278axpla", msg.Amount.String())
	suite.Require().Equal("xpla1ghekyjucln7y67ntx7cf27m9dpuxxemntll57s", msg.FromAddress)
	suite.Require().Equal("xpla1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7zj7eku", msg.ToAddress)
	suite.Require().Equal(banktypes.TypeMsgSend, msg.Type())
	suite.Require().Equal(banktypes.ModuleName, msg.Route())
	suite.Require().Len(futureOperations, 0)
}

// SimulateMsgSend tests and runs a single msg send where both
// accounts already exist.
func SimulateMsgSend(ak banktypes.AccountKeeper, bk keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		from, to, coins, skip := randomSendFields(r, ctx, accs, bk, ak)

		// Check send_enabled status of each coin denom
		if err := bk.IsSendEnabledCoins(ctx, coins...); err != nil {
			return simtypes.NoOpMsg(banktypes.ModuleName, banktypes.TypeMsgSend, err.Error()), nil, nil
		}

		if skip {
			return simtypes.NoOpMsg(banktypes.ModuleName, banktypes.TypeMsgSend, "skip all transfers"), nil, nil
		}

		bankSendMsg := types.BankSendMsg{
			FromAddress: from.Address.String(),
			ToAddress:   to.Address.String(),
			Amount:      coins.String(),
		}

		msg, err := mbank.MakeBankSendMsg(bankSendMsg, from.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(banktypes.ModuleName, banktypes.TypeMsgSend, err.Error()), nil, err
		}

		err = sendMsgSend(r, app, bk, ak, &msg, ctx, chainID, []cryptotypes.PrivKey{from.PrivKey})
		if err != nil {
			return simtypes.NoOpMsg(banktypes.ModuleName, msg.Type(), "invalid transfers"), nil, err
		}

		return simtypes.NewOperationMsg(&msg, true, "", nil), nil, nil
	}
}

func (suite *TestSuite) getTestingAccounts(r *rand.Rand, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := suite.app.StakingKeeper.TokensFromConsensusPower(suite.ctx, 200)
	initCoins := sdk.NewCoins(sdk.NewCoin("axpla", initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := suite.app.AccountKeeper.NewAccountWithAddress(suite.ctx, account.Address)
		suite.app.AccountKeeper.SetAccount(suite.ctx, acc)
		suite.Require().NoError(testutil.FundAccount(suite.app.BankKeeper, suite.ctx, account.Address, initCoins))
	}

	return accounts
}

// randomSendFields returns the sender and recipient simulation accounts as well
// as the transferred amount.
func randomSendFields(
	r *rand.Rand, ctx sdk.Context, accs []simtypes.Account, bk keeper.Keeper, ak banktypes.AccountKeeper,
) (simtypes.Account, simtypes.Account, sdk.Coins, bool) {

	from, _ := simtypes.RandomAcc(r, accs)
	to, _ := simtypes.RandomAcc(r, accs)

	// disallow sending money to yourself
	for from.PubKey.Equals(to.PubKey) {
		to, _ = simtypes.RandomAcc(r, accs)
	}

	acc := ak.GetAccount(ctx, from.Address)
	if acc == nil {
		return from, to, nil, true
	}

	spendable := bk.SpendableCoins(ctx, acc.GetAddress())
	sendCoins := simtypes.RandSubsetCoins(r, spendable)
	if sendCoins.Empty() {
		return from, to, nil, true
	}

	return from, to, sendCoins, false
}

// sendMsgSend sends a transaction with a MsgSend from a provided random account.
func sendMsgSend(
	r *rand.Rand, app *baseapp.BaseApp, bk keeper.Keeper, ak banktypes.AccountKeeper,
	msg *banktypes.MsgSend, ctx sdk.Context, chainID string, privkeys []cryptotypes.PrivKey,
) error {

	var (
		fees sdk.Coins
		err  error
	)

	from, err := sdk.AccAddressFromBech32(msg.FromAddress)
	if err != nil {
		return err
	}

	account := ak.GetAccount(ctx, from)
	spendable := bk.SpendableCoins(ctx, account.GetAddress())

	coins, hasNeg := spendable.SafeSub(msg.Amount)
	if !hasNeg {
		fees, err = simtypes.RandomFees(r, ctx, coins)
		if err != nil {
			return err
		}
	}
	txGen := util.MakeEncodingConfig().TxConfig
	tx, err := testutil.GenTx(
		txGen,
		[]sdk.Msg{msg},
		fees,
		testutil.DefaultTestGenTxGas,
		chainID,
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
		privkeys...,
	)
	if err != nil {
		return err
	}

	_, _, err = app.Deliver(txGen.TxEncoder(), tx)
	if err != nil {
		return err
	}

	return nil
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
