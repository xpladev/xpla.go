package authz_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/xpladev/xpla.go/client"
	mauthz "github.com/xpladev/xpla.go/core/authz"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/cosmos/cosmos-sdk/x/authz/keeper"
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
	queryClient authz.QueryClient
}

// authz message types
var (
	TypeMsgGrant  = sdk.MsgTypeURL(&authz.MsgGrant{})
	TypeMsgRevoke = sdk.MsgTypeURL(&authz.MsgRevoke{})
	TypeMsgExec   = sdk.MsgTypeURL(&authz.MsgExec{})
)

func (suite *TestSuite) SetupTest() {
	checkTx := false
	app := testutil.Setup(checkTx, 5)
	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	authz.RegisterQueryServer(queryHelper, app.AuthzKeeper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = authz.NewQueryClient(queryHelper)

}

func (suite *TestSuite) TestSimulateGrant() {
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 2)
	blockTime := time.Now().UTC()
	ctx := suite.ctx.WithBlockTime(blockTime)

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: tmproto.Header{
			Height:  suite.app.LastBlockHeight() + 1,
			AppHash: suite.app.LastCommitID().Hash,
		},
	})

	granter := accounts[0]
	grantee := accounts[1]

	// execute operation
	op := SimulateMsgGrant(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.AuthzKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, ctx, accounts, testutil.TestChainId)
	suite.Require().NoError(err)

	var msg authz.MsgGrant
	suite.app.AppCodec().UnmarshalJSON(operationMsg.Msg, &msg)
	suite.Require().True(operationMsg.OK)
	suite.Require().Equal(granter.Address.String(), msg.Granter)
	suite.Require().Equal(grantee.Address.String(), msg.Grantee)
	suite.Require().Len(futureOperations, 0)
}

func (suite *TestSuite) TestSimulateRevoke() {
	s := rand.NewSource(2)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{
		Header: tmproto.Header{
			Height:  suite.app.LastBlockHeight() + 1,
			AppHash: suite.app.LastCommitID().Hash,
		}})

	initAmt := suite.app.StakingKeeper.TokensFromConsensusPower(suite.ctx, 200000)
	initCoins := sdk.NewCoins(sdk.NewCoin("axpla", initAmt))

	granter := accounts[0]
	grantee := accounts[1]
	authorization := banktypes.NewSendAuthorization(initCoins)

	err := suite.app.AuthzKeeper.SaveGrant(suite.ctx, grantee.Address, granter.Address, authorization, time.Now().Add(30*time.Hour))
	suite.Require().NoError(err)

	// execute operation
	op := SimulateMsgRevoke(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.AuthzKeeper)
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, testutil.TestChainId)
	suite.Require().NoError(err)

	var msg authz.MsgRevoke
	suite.app.AppCodec().UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal(granter.Address.String(), msg.Granter)
	suite.Require().Equal(grantee.Address.String(), msg.Grantee)
	suite.Require().Equal(banktypes.SendAuthorization{}.MsgTypeURL(), msg.MsgTypeUrl)
	suite.Require().Len(futureOperations, 0)
}

func (suite *TestSuite) TestSimulateExec() {
	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 3)

	// begin a new block
	suite.app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: suite.app.LastBlockHeight() + 1, AppHash: suite.app.LastCommitID().Hash}})

	initAmt := suite.app.StakingKeeper.TokensFromConsensusPower(suite.ctx, 200000)
	initCoins := sdk.NewCoins(sdk.NewCoin("axpla", initAmt))

	granter := accounts[0]
	grantee := accounts[1]
	authorization := banktypes.NewSendAuthorization(initCoins)

	err := suite.app.AuthzKeeper.SaveGrant(suite.ctx, grantee.Address, granter.Address, authorization, time.Now().Add(30*time.Hour))
	suite.Require().NoError(err)

	// execute operation
	op := SimulateMsgExec(suite.app.AccountKeeper, suite.app.BankKeeper, suite.app.AuthzKeeper, suite.app.AppCodec())
	operationMsg, futureOperations, err := op(r, suite.app.BaseApp, suite.ctx, accounts, testutil.TestChainId)
	suite.Require().NoError(err)

	var msg authz.MsgExec

	suite.app.AppCodec().UnmarshalJSON(operationMsg.Msg, &msg)

	suite.Require().True(operationMsg.OK)
	suite.Require().Equal(grantee.Address.String(), msg.Grantee)
	suite.Require().Len(futureOperations, 0)
}

// SimulateMsgGrant generates a MsgGrant with random values.
func SimulateMsgGrant(ak authz.AccountKeeper, bk authz.BankKeeper, _ keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {

		granter := accs[0]
		grantee := accs[1]

		granterAddr := granter.Address.String()
		granteeAddr := grantee.Address.String()
		granterAcc := ak.GetAccount(ctx, granter.Address)

		spendLimit := util.FromIntToString(r.Intn(10000000))
		// Expiration time is unix timestamp
		exp := "1695861001"

		authzTypes := []string{"send", "generic", "delegate", "unbond", "redelegate"}
		randomAuthzType := generateRandomAuthorization(r, authzTypes)

		testMsgType := []string{
			"cosmos.authz.v1beta1.MsgGrant",
			"cosmos.gov.v1beta1.MsgVote",
			"cosmos.feegrant.v1beta1.MsgGrantAllowance",
			"cosmos.bank.v1beta1.MsgSend",
		}

		var msgType string
		var allowedValidators []string
		var denyValidators []string
		if randomAuthzType == "generic" {
			msgType = generateRandomMsgType(r, testMsgType)

		} else if randomAuthzType == "delegate" ||
			randomAuthzType == "unbond" ||
			randomAuthzType == "redelegate" {

			mnemonicVal1, _ := key.NewMnemonic()
			kVal1, _ := key.NewPrivKey(mnemonicVal1)
			mnemonicVal2, _ := key.NewMnemonic()
			kVal2, _ := key.NewPrivKey(mnemonicVal2)

			val1 := sdk.ValAddress(kVal1.PubKey().Address())
			val2 := sdk.ValAddress(kVal2.PubKey().Address())

			allowedValidators = append(allowedValidators, val1.String())
			denyValidators = append(denyValidators, val2.String())
		}

		authzGrantMsg := types.AuthzGrantMsg{
			Granter:           granterAddr,
			Grantee:           granteeAddr,
			AuthorizationType: randomAuthzType,
			SpendLimit:        spendLimit,
			Expiration:        exp,
			MsgType:           msgType,
			AllowValidators:   allowedValidators,
			DenyValidators:    denyValidators,
		}

		spendableCoins := bk.SpendableCoins(ctx, granter.Address)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgRevoke, "fee error"), nil, err
		}

		msg, err := mauthz.MakeAuthzGrantMsg(authzGrantMsg, granter.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgGrant, err.Error()), nil, err
		}
		txCfg := util.MakeEncodingConfig().TxConfig
		tx, err := testutil.GenTx(
			txCfg,
			[]sdk.Msg{&msg},
			fees,
			testutil.DefaultTestGenTxGas,
			chainID,
			[]uint64{granterAcc.GetAccountNumber()},
			[]uint64{granterAcc.GetSequence()},
			granter.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgGrant, "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txCfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, sdk.MsgTypeURL(&msg), "unable to deliver tx"), nil, err
		}
		return simtypes.NewOperationMsg(&msg, true, "", nil), nil, err
	}
}

// SimulateMsgRevoke generates a MsgRevoke with random values.
func SimulateMsgRevoke(ak authz.AccountKeeper, bk authz.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		var granterAddr, granteeAddr sdk.AccAddress
		var grant authz.Grant
		hasGrant := false

		k.IterateGrants(ctx, func(granter, grantee sdk.AccAddress, g authz.Grant) bool {
			grant = g
			granterAddr = granter
			granteeAddr = grantee
			hasGrant = true
			return true
		})

		if !hasGrant {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgRevoke, "no grants"), nil, nil
		}

		granterAcc, ok := simtypes.FindAccount(accs, granterAddr)
		if !ok {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgRevoke, "account not found"), nil, util.LogErr(errors.ErrInvalidRequest, "account not found")
		}

		spendableCoins := bk.SpendableCoins(ctx, granterAddr)
		fees, err := simtypes.RandomFees(r, ctx, spendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgRevoke, "fee error"), nil, err
		}

		a := grant.GetAuthorization()

		authzRevokeMsg := types.AuthzRevokeMsg{
			Granter: granterAddr.String(),
			Grantee: granteeAddr.String(),
			MsgType: a.MsgTypeURL(),
		}

		msg, err := mauthz.MakeAuthzRevokeMsg(authzRevokeMsg, granterAcc.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgRevoke, "make msg error"), nil, err
		}

		// msg := authz.NewMsgRevoke(granterAddr, granteeAddr, a.MsgTypeURL())
		txCfg := util.MakeEncodingConfig().TxConfig
		account := ak.GetAccount(ctx, granterAddr)
		tx, err := testutil.GenTx(
			txCfg,
			[]sdk.Msg{&msg},
			fees,
			testutil.DefaultTestGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			granterAcc.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgRevoke, err.Error()), nil, err
		}

		_, _, err = app.Deliver(txCfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgRevoke, "unable to deliver tx"), nil, err
		}

		return simtypes.NewOperationMsg(&msg, true, "", nil), nil, nil
	}
}

// SimulateMsgExec generates a MsgExec with random values.
func SimulateMsgExec(ak authz.AccountKeeper, bk authz.BankKeeper, k keeper.Keeper, cdc cdctypes.AnyUnpacker) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		hasGrant := false
		var targetGrant authz.Grant
		var granterAddr sdk.AccAddress
		var granteeAddr sdk.AccAddress
		k.IterateGrants(ctx, func(granter, grantee sdk.AccAddress, grant authz.Grant) bool {
			targetGrant = grant
			granterAddr = granter
			granteeAddr = grantee
			hasGrant = true
			return true
		})

		if !hasGrant {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, "no grant found"), nil, nil
		}

		if targetGrant.Expiration.Before(ctx.BlockHeader().Time) {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, "grant expired"), nil, nil
		}

		grantee, ok := simtypes.FindAccount(accs, granteeAddr)
		if !ok {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgRevoke, "Account not found"), nil, util.LogErr(errors.ErrInvalidRequest, "grantee account not found")
		}

		if _, ok := simtypes.FindAccount(accs, granterAddr); !ok {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgRevoke, "Account not found"), nil, util.LogErr(errors.ErrInvalidRequest, "granter account not found")
		}

		granterspendableCoins := bk.SpendableCoins(ctx, granterAddr)
		coins := simtypes.RandSubsetCoins(r, granterspendableCoins)
		// Check send_enabled status of each sent coin denom
		if err := bk.IsSendEnabledCoins(ctx, coins...); err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, err.Error()), nil, nil
		}

		xplac := client.NewXplaClient(chainID)
		xplac.WithPrivateKey(grantee.PrivKey)
		bankSendMsg := types.BankSendMsg{
			FromAddress: granteeAddr.String(),
			ToAddress:   granterAddr.String(),
			Amount:      coins.String(),
		}
		txbytes, err := xplac.BankSend(bankSendMsg).CreateAndSignTx()
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, "create and sign tx error"), nil, err
		}

		encodingConfig := xplac.GetEncoding()
		txDecode, err := encodingConfig.TxConfig.TxDecoder()(txbytes)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, "cannot encode tx"), nil, err
		}

		json, err := encodingConfig.TxConfig.TxJSONEncoder()(txDecode)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, "cannot encode tx"), nil, err
		}

		authzExecMsg := types.AuthzExecMsg{
			Grantee:      granteeAddr.String(),
			ExecTxString: string(json),
		}
		msgExec, err := mauthz.MakeAuthzExecMsg(authzExecMsg, encodingConfig)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, "msg error"), nil, err
		}

		granteeSpendableCoins := bk.SpendableCoins(ctx, granteeAddr)
		fees, err := simtypes.RandomFees(r, ctx, granteeSpendableCoins)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, "fee error"), nil, err
		}

		txCfg := encodingConfig.TxConfig
		granteeAcc := ak.GetAccount(ctx, granteeAddr)
		tx, err := testutil.GenTx(
			txCfg,
			[]sdk.Msg{&msgExec},
			fees,
			testutil.DefaultTestGenTxGas,
			chainID,
			[]uint64{granteeAcc.GetAccountNumber()},
			[]uint64{granteeAcc.GetSequence()},
			grantee.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, err.Error()), nil, err
		}

		_, _, err = app.Deliver(txCfg.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, err.Error()), nil, err
		}

		err = msgExec.UnpackInterfaces(cdc)
		if err != nil {
			return simtypes.NoOpMsg(authz.ModuleName, TypeMsgExec, "unmarshal error"), nil, err
		}
		return simtypes.NewOperationMsg(&msgExec, true, "success", nil), nil, nil
	}
}

func generateRandomAuthorization(r *rand.Rand, authzTypes []string) string {
	return authzTypes[r.Intn(len(authzTypes))]
}

func generateRandomMsgType(r *rand.Rand, msgType []string) string {
	return msgType[r.Intn(len(msgType))]
}

func (suite *TestSuite) getTestingAccounts(r *rand.Rand, n int) []simtypes.Account {
	accounts := testutil.RandomAccounts(r, n)

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
