package mint_test

import (
	"context"
	"testing"

	mmint "github.com/xpladev/xpla.go/core/mint"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type TestSuite struct {
	suite.Suite

	app         *xapp.XplaApp
	ctx         sdk.Context
	queryClient minttypes.QueryClient
}

func (suite *TestSuite) SetupTest() {
	checkTx := false
	app := testutil.Setup(checkTx, 5)
	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	minttypes.RegisterQueryServer(queryHelper, app.MintKeeper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = minttypes.NewQueryClient(queryHelper)
}

func (suite *TestSuite) TestGRPCParams() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	msgMint, _ := mmint.MakeQueryMintParamsMsg()
	params, err := queryClient.Params(context.Background(), &msgMint)
	suite.Require().NoError(err)
	suite.Require().Equal(params.Params, app.MintKeeper.GetParams(ctx))

	msgInflation, _ := mmint.MakeQueryInflationMsg()
	inflation, err := queryClient.Inflation(context.Background(), &msgInflation)
	suite.Require().NoError(err)
	suite.Require().Equal(inflation.Inflation, app.MintKeeper.GetMinter(ctx).Inflation)

	msgAnnualProvisions, _ := mmint.MakeQueryAnnualProvisionsMsg()
	annualProvisions, err := queryClient.AnnualProvisions(context.Background(), &msgAnnualProvisions)
	suite.Require().NoError(err)
	suite.Require().Equal(annualProvisions.AnnualProvisions, app.MintKeeper.GetMinter(ctx).AnnualProvisions)
}

func TestMintTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
