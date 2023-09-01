package feegrant_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	mfeegrant "github.com/xpladev/xpla.go/core/feegrant"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	"github.com/cosmos/cosmos-sdk/x/feegrant/keeper"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type TestSuite struct {
	suite.Suite

	ctx     sdk.Context
	context context.Context
	app     *xapp.XplaApp
	keeper  keeper.Keeper
}

var (
	TypeMsgGrantAllowance  = sdk.MsgTypeURL(&feegrant.MsgGrantAllowance{})
	TypeMsgRevokeAllowance = sdk.MsgTypeURL(&feegrant.MsgRevokeAllowance{})
)

func (suite *TestSuite) SetupTest() {
	checkTx := false
	app := testutil.Setup(checkTx, 5)
	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{
		Time: time.Now(),
	})

	suite.app = app
	suite.ctx = ctx
	suite.context = sdk.WrapSDKContext(ctx)
	suite.keeper = suite.app.FeeGrantKeeper
}

func (suite *TestSuite) TestSimulateMsgGrantAllowance() {
	app, ctx := suite.app, suite.ctx
	require := suite.Require()

	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := SimulateMsgGrantAllowance(app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(err)

	var msg feegrant.MsgGrantAllowance
	suite.app.AppCodec().UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(operationMsg.OK)
	require.Equal(accounts[2].Address.String(), msg.Granter)
	require.Equal(accounts[1].Address.String(), msg.Grantee)
	require.Len(futureOperations, 0)
}

func (suite *TestSuite) TestSimulateMsgRevokeAllowance() {
	app, ctx := suite.app, suite.ctx
	require := suite.Require()

	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: suite.app.LastBlockHeight() + 1, AppHash: suite.app.LastCommitID().Hash}})

	feeAmt := app.StakingKeeper.TokensFromConsensusPower(ctx, 200000)
	feeCoins := sdk.NewCoins(sdk.NewCoin("axpla", feeAmt))

	granter, grantee := accounts[0], accounts[1]

	oneYear := ctx.BlockTime().AddDate(1, 0, 0)
	err := app.FeeGrantKeeper.GrantAllowance(
		ctx,
		granter.Address,
		grantee.Address,
		&feegrant.BasicAllowance{
			SpendLimit: feeCoins,
			Expiration: &oneYear,
		},
	)
	require.NoError(err)

	// execute operation
	op := SimulateMsgRevokeAllowance(app.AccountKeeper, app.BankKeeper, app.FeeGrantKeeper)
	operationMsg, futureOperations, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(err)

	var msg feegrant.MsgRevokeAllowance
	suite.app.AppCodec().UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(operationMsg.OK)
	require.Equal(granter.Address.String(), msg.Granter)
	require.Equal(grantee.Address.String(), msg.Grantee)
	require.Len(futureOperations, 0)
}

// SimulateMsgGrantAllowance generates MsgGrantAllowance with random values.
func SimulateMsgGrantAllowance(ak feegrant.AccountKeeper, bk feegrant.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		granter, _ := simtypes.RandomAcc(r, accs)
		grantee, _ := simtypes.RandomAcc(r, accs)
		if grantee.Address.String() == granter.Address.String() {
			return simtypes.NoOpMsg(feegrant.ModuleName, TypeMsgGrantAllowance, "grantee and granter cannot be same"), nil, nil
		}

		if f, _ := k.GetAllowance(ctx, granter.Address, grantee.Address); f != nil {
			return simtypes.NoOpMsg(feegrant.ModuleName, TypeMsgGrantAllowance, "fee allowance exists"), nil, nil
		}

		account := ak.GetAccount(ctx, granter.Address)

		spendableCoins := bk.SpendableCoins(ctx, account.GetAddress())
		if spendableCoins.Empty() {
			return simtypes.NoOpMsg(feegrant.ModuleName, TypeMsgGrantAllowance, "unable to grant empty coins as SpendLimit"), nil, nil
		}

		oneYear, _ := ctx.BlockTime().AddDate(1, 0, 0).MarshalText()

		grantMsg := types.FeeGrantMsg{
			Grantee:    grantee.Address.String(),
			Granter:    granter.Address.String(),
			SpendLimit: spendableCoins.String(),
			Expiration: string(oneYear),
		}

		msg, err := mfeegrant.MakeFeeGrantMsg(grantMsg, granter.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(feegrant.ModuleName, TypeMsgGrantAllowance, "make msg err"), nil, err
		}

		if err != nil {
			return simtypes.NoOpMsg(feegrant.ModuleName, TypeMsgGrantAllowance, err.Error()), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           util.MakeEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         TypeMsgGrantAllowance,
			Context:         ctx,
			SimAccount:      granter,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      feegrant.ModuleName,
			CoinsSpentInMsg: spendableCoins,
		}

		return testutil.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// SimulateMsgRevokeAllowance generates a MsgRevokeAllowance with random values.
func SimulateMsgRevokeAllowance(ak feegrant.AccountKeeper, bk feegrant.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		hasGrant := false
		var granterAddr sdk.AccAddress
		var granteeAddr sdk.AccAddress
		k.IterateAllFeeAllowances(ctx, func(grant feegrant.Grant) bool {
			granter := sdk.MustAccAddressFromBech32(grant.Granter)
			grantee := sdk.MustAccAddressFromBech32(grant.Grantee)
			granterAddr = granter
			granteeAddr = grantee
			hasGrant = true
			return true
		})

		if !hasGrant {
			return simtypes.NoOpMsg(feegrant.ModuleName, TypeMsgRevokeAllowance, "no grants"), nil, nil
		}
		granter, ok := simtypes.FindAccount(accs, granterAddr)

		if !ok {
			return simtypes.NoOpMsg(feegrant.ModuleName, TypeMsgRevokeAllowance, "Account not found"), nil, nil
		}

		account := ak.GetAccount(ctx, granter.Address)
		spendableCoins := bk.SpendableCoins(ctx, account.GetAddress())

		revokeGrantMsg := types.RevokeFeeGrantMsg{
			Granter: granterAddr.String(),
			Grantee: granteeAddr.String(),
		}

		msg, err := mfeegrant.MakeRevokeFeeGrantMsg(revokeGrantMsg, granter.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(feegrant.ModuleName, TypeMsgRevokeAllowance, "make msg err"), nil, err
		}

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           util.MakeEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         TypeMsgRevokeAllowance,
			Context:         ctx,
			SimAccount:      granter,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      feegrant.ModuleName,
			CoinsSpentInMsg: spendableCoins,
		}

		return testutil.GenAndDeliverTxWithRandFees(txCtx)
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

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
