package client

import (
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/authz"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type TestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *xapp.XplaApp
	queryClient authz.QueryClient
}

const (
	unsignedTxPath  = "../util/testutil/test_files/unsignedTx.json"
	signedTxPath    = "../util/testutil/test_files/signedTx.json"
	signedEVMTxPath = "../util/testutil/test_files/signedEVMTx.json"
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

func (suite *TestSuite) TestSimulateCreateUnsignedTx() {
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 4)
	from := accounts[0]
	to := accounts[1]
	coins := randomSendCoins(r, suite.ctx, from, suite.app.BankKeeper, suite.app.AccountKeeper)

	gasLimitU64, err := util.FromStringToUint64(types.DefaultGasLimit)
	suite.Require().NoError(err)
	gasPriceU64, err := util.FromStringToUint64(types.DefaultGasPrice)
	suite.Require().NoError(err)

	fee := util.FromUint64ToString(util.MulUint64(gasLimitU64, gasPriceU64))

	xplac := NewXplaClient(testutil.TestChainId)
	xplac.WithOptions(
		Options{
			GasLimit:  types.DefaultGasLimit,
			FeeAmount: fee,
		},
	)

	bankSendMsg := types.BankSendMsg{
		FromAddress: from.Address.String(),
		ToAddress:   to.Address.String(),
		Amount:      coins.String(),
	}

	txbytes, err := xplac.BankSend(bankSendMsg).CreateUnsignedTx()
	suite.Require().NoError(err)

	clientCtx, err := util.NewClient()
	suite.Require().NoError(err)

	_, _, newTx, err := readTxAndInitContexts(clientCtx, unsignedTxPath)
	suite.Require().NoError(err)

	newTxbytes, err := xplac.GetEncoding().TxConfig.TxEncoder()(newTx)
	suite.Require().NoError(err)
	suite.Require().Equal(txbytes, newTxbytes)
}

func (suite *TestSuite) TestSimulateSignTx() {
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 4)
	from := accounts[0]

	xplac := NewXplaClient(testutil.TestChainId)
	xplac.WithOptions(
		Options{
			GasLimit:   types.DefaultGasLimit,
			GasPrice:   types.DefaultGasPrice,
			PrivateKey: from.PrivKey,
		},
	)

	signTxMsg := types.SignTxMsg{
		UnsignedFileName: unsignedTxPath,
		SignatureOnly:    false,
		MultisigAddress:  "",
		Overwrite:        false,
		Amino:            false,
	}

	txbytes, err := xplac.SignTx(signTxMsg)
	suite.Require().NoError(err)

	clientCtx, err := util.NewClient()
	suite.Require().NoError(err)

	_, _, newTx, err := readTxAndInitContexts(clientCtx, signedTxPath)
	suite.Require().NoError(err)

	newTxbytes, err := xplac.GetEncoding().TxConfig.TxJSONEncoder()(newTx)
	suite.Require().NoError(err)
	suite.Require().Equal(txbytes, newTxbytes)
}

func (suite *TestSuite) TestSimulateSignatureOnly() {
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 4)
	from := accounts[0]
	to := accounts[1]
	m1 := accounts[2]
	m2 := accounts[3]

	accList := []simtypes.Account{from, to, m1, m2}
	xplac := NewXplaClient(testutil.TestChainId)
	xplac.WithOptions(
		Options{
			GasLimit: types.DefaultGasLimit,
			GasPrice: types.DefaultGasPrice,
		},
	)

	for i, acc := range accList {
		xplac.WithOptions(
			Options{
				PrivateKey: acc.PrivKey,
			},
		)

		signTxMsg := types.SignTxMsg{
			UnsignedFileName: unsignedTxPath,
			SignatureOnly:    true,
			MultisigAddress:  "",
			Overwrite:        false,
			Amino:            false,
		}

		txbytes, err := xplac.SignTx(signTxMsg)
		suite.Require().NoError(err)

		signPath := "../util/testutil/test_files/signature" + util.FromIntToString(i) + ".json"
		jsonByte, err := convertJson(signPath)
		suite.Require().NoError(err)

		suite.Require().Equal(txbytes, jsonByte)
	}
}

func (suite *TestSuite) TestSimulateCreateAndSignTx() {
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 4)
	from := accounts[0]
	to := accounts[1]
	coins := randomSendCoins(r, suite.ctx, from, suite.app.BankKeeper, suite.app.AccountKeeper)

	xplac := NewXplaClient(testutil.TestChainId)
	xplac.WithPrivateKey(from.PrivKey)

	bankSendMsg := types.BankSendMsg{
		FromAddress: from.Address.String(),
		ToAddress:   to.Address.String(),
		Amount:      coins.String(),
	}

	txbytes, err := xplac.BankSend(bankSendMsg).CreateAndSignTx()
	suite.Require().NoError(err)

	clientCtx, err := util.NewClient()
	suite.Require().NoError(err)

	_, _, newTx, err := readTxAndInitContexts(clientCtx, signedTxPath)
	suite.Require().NoError(err)

	newTxbytes, err := xplac.GetEncoding().TxConfig.TxEncoder()(newTx)
	suite.Require().NoError(err)
	suite.Require().Equal(txbytes, newTxbytes)
}

func (suite *TestSuite) TestSimulateEVMCreateAndSignTx() {
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := suite.getTestingAccounts(r, 4)
	from := accounts[0]

	testABIJsonFilePath := "../util/testutil/test_files/abi.json"
	testBytecodeJsonFilePath := "../util/testutil/test_files/bytecode.json"

	xplac := NewXplaClient(testutil.TestChainId)
	xplac.WithPrivateKey(from.PrivKey)

	// deploy
	deploySolContractMsg := types.DeploySolContractMsg{
		ABIJsonFilePath:      testABIJsonFilePath,
		BytecodeJsonFilePath: testBytecodeJsonFilePath,
		Args:                 nil,
	}
	txbytes, err := xplac.DeploySolidityContract(deploySolContractMsg).CreateAndSignTx()
	suite.Require().NoError(err)

	file, err := os.ReadFile(signedEVMTxPath)
	suite.Require().NoError(err)

	suite.Require().Equal(txbytes, file)
}

func (suite *TestSuite) TestSimulateEncodeAndDecodeTx() {
	xplac := NewXplaClient(testutil.TestChainId)

	encodeTxMsg := types.EncodeTxMsg{
		FileName: unsignedTxPath,
	}

	encodeRes, err := xplac.EncodeTx(encodeTxMsg)
	suite.Require().NoError(err)
	encoded := "Cp8BCpwBChwvY29zbW9zLmJhbmsudjFiZXRhMS5Nc2dTZW5kEnwKK3hwbGExbDAza21hNHZ2OXFjdmhnY3hmMmdhMHJudjdkcWN1bWF4ZXhzc2gSK3hwbGExaDR4MmpsbnFrenEyazh3cmZ6dnR0bDlwNWdjZmZ6NHhlNWNqMmMaIAoFYXhwbGESFzgzNjIyODM5MDIyNDM4MjI3MDU2NTgyEiMSIQobCgVheHBsYRISMjEyNTAwMDAwMDAwMDAwMDAwEJChDw=="
	suite.Require().Equal(encodeRes, encoded)

	decodeTxMsg := types.DecodeTxMsg{
		EncodedByteString: encoded,
	}

	decodeRes, err := xplac.DecodeTx(decodeTxMsg)
	suite.Require().NoError(err)

	clientCtx, err := util.NewClient()
	suite.Require().NoError(err)

	_, _, newTx, err := readTxAndInitContexts(clientCtx, unsignedTxPath)
	suite.Require().NoError(err)

	newTxbytes, err := xplac.GetEncoding().TxConfig.TxJSONEncoder()(newTx)
	suite.Require().NoError(err)
	suite.Require().Equal([]byte(decodeRes), newTxbytes)
}

func (suite *TestSuite) TestSimulateValidateSignature() {
	xplac := NewXplaClient(testutil.TestChainId)

	validateSignaturesMsg := types.ValidateSignaturesMsg{
		FileName: signedTxPath,
		Offline:  true,
	}

	res, err := xplac.ValidateSignatures(validateSignaturesMsg)
	suite.Require().NoError(err)
	suite.Require().Equal(res, "success validate")
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

func randomSendCoins(
	r *rand.Rand, ctx sdk.Context, account simtypes.Account, bk bankkeeper.Keeper, ak banktypes.AccountKeeper,
) sdk.Coins {
	acc := ak.GetAccount(ctx, account.Address)
	if acc == nil {
		return nil
	}

	spendable := bk.SpendableCoins(ctx, acc.GetAddress())
	sendCoins := simtypes.RandSubsetCoins(r, spendable)
	if sendCoins.Empty() {
		return nil
	}

	return sendCoins
}

func convertJson(filePath string) ([]byte, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	temp := strings.Replace(string(bytes), " ", "", -1)
	temp = strings.Replace(temp, "\n", "", -1)
	return []byte(temp), nil
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
