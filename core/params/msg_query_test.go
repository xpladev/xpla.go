package params_test

import (
	"fmt"
	"testing"

	mparams "github.com/xpladev/xpla.go/core/params"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type TestSuite struct {
	suite.Suite

	app         *xapp.XplaApp
	ctx         sdk.Context
	queryClient proposal.QueryClient
}

func (suite *TestSuite) SetupTest() {
	checkTx := false
	app := testutil.Setup(checkTx, 5)
	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	proposal.RegisterQueryServer(queryHelper, app.ParamsKeeper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = proposal.NewQueryClient(queryHelper)
}

type paramJSON struct {
	Param1 int64  `json:"param1,omitempty" yaml:"param1,omitempty"`
	Param2 string `json:"param2,omitempty" yaml:"param2,omitempty"`
}

func validateNoOp(_ interface{}) error { return nil }

func (suite *TestSuite) TestGRPCQueryParams() {
	var (
		req      *proposal.QueryParamsRequest
		expValue string
		space    paramstypes.Subspace
	)
	key := []byte("key")

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"invalid request with subspace not found",
			func() {
				subspaceMsg := types.SubspaceMsg{
					Subspace: "test",
				}
				msg, _ := mparams.MakeQueryParamsSubspaceMsg(subspaceMsg)
				req = &msg
			},
			false,
		},
		{
			"invalid request with subspace and key not found",
			func() {
				subspaceMsg := types.SubspaceMsg{
					Subspace: "test",
					Key:      "key",
				}
				msg, _ := mparams.MakeQueryParamsSubspaceMsg(subspaceMsg)
				req = &msg
			},
			false,
		},
		{
			"success",
			func() {
				space = suite.app.ParamsKeeper.Subspace("test").
					WithKeyTable(paramstypes.NewKeyTable(paramstypes.NewParamSetPair(key, paramJSON{}, validateNoOp)))

				subspaceMsg := types.SubspaceMsg{
					Subspace: "test",
					Key:      "key",
				}
				msg, _ := mparams.MakeQueryParamsSubspaceMsg(subspaceMsg)
				req = &msg
				expValue = ""
			},
			true,
		},
		{
			"update value success",
			func() {
				err := space.Update(suite.ctx, key, []byte(`{"param1":"10241024"}`))
				suite.Require().NoError(err)
				subspaceMsg := types.SubspaceMsg{
					Subspace: "test",
					Key:      "key",
				}
				msg, _ := mparams.MakeQueryParamsSubspaceMsg(subspaceMsg)
				req = &msg
				expValue = `{"param1":"10241024"}`
			},
			true,
		},
	}

	suite.SetupTest()
	ctx := sdk.WrapSDKContext(suite.ctx)

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()

			res, err := suite.queryClient.Params(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expValue, res.Param.Value)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}
		})
	}
}

func TestMintTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
