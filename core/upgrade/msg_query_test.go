package upgrade_test

import (
	"context"
	"fmt"
	"testing"

	mupgrade "github.com/xpladev/xpla.go/core/upgrade"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type UpgradeTestSuite struct {
	suite.Suite

	app         *xapp.XplaApp
	ctx         sdk.Context
	queryClient upgradetypes.QueryClient
}

func (suite *UpgradeTestSuite) SetupTest() {
	suite.app = testutil.Setup(false, 5)
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(suite.ctx, suite.app.InterfaceRegistry())
	upgradetypes.RegisterQueryServer(queryHelper, suite.app.UpgradeKeeper)
	suite.queryClient = upgradetypes.NewQueryClient(queryHelper)
}

func (suite *UpgradeTestSuite) TestQueryCurrentPlan() {
	var (
		req         *upgradetypes.QueryCurrentPlanRequest
		expResponse upgradetypes.QueryCurrentPlanResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"without current upgrade plan",
			func() {
				msg, _ := mupgrade.MakePlanMsg()
				req = &msg
				expResponse = upgradetypes.QueryCurrentPlanResponse{}
			},
			true,
		},
		{
			"with current upgrade plan",
			func() {
				plan := upgradetypes.Plan{Name: "test-plan", Height: 5}
				suite.app.UpgradeKeeper.ScheduleUpgrade(suite.ctx, plan)

				msg, _ := mupgrade.MakePlanMsg()
				req = &msg
				expResponse = upgradetypes.QueryCurrentPlanResponse{Plan: &plan}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()

			res, err := suite.queryClient.CurrentPlan(context.Background(), req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(&expResponse, res)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *UpgradeTestSuite) TestAppliedCurrentPlan() {
	var (
		req       *upgradetypes.QueryAppliedPlanRequest
		expHeight int64
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"with non-existent upgrade plan",
			func() {
				appliedMsg := types.AppliedMsg{
					UpgradeName: "foo",
				}
				msg, _ := mupgrade.MakeAppliedMsg(appliedMsg)
				req = &msg
			},
			true,
		},
		{
			"with applied upgrade plan",
			func() {
				expHeight = 5

				planName := "test-plan"
				plan := upgradetypes.Plan{Name: planName, Height: expHeight}
				suite.app.UpgradeKeeper.ScheduleUpgrade(suite.ctx, plan)

				suite.ctx = suite.ctx.WithBlockHeight(expHeight)
				suite.app.UpgradeKeeper.SetUpgradeHandler(planName, func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
					return vm, nil
				})
				suite.app.UpgradeKeeper.ApplyUpgrade(suite.ctx, plan)

				appliedMsg := types.AppliedMsg{
					UpgradeName: planName,
				}
				msg, _ := mupgrade.MakeAppliedMsg(appliedMsg)
				req = &msg
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset

			tc.malleate()

			res, err := suite.queryClient.AppliedPlan(context.Background(), req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expHeight, res.Height)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *UpgradeTestSuite) TestModuleVersions() {
	var req *upgradetypes.QueryModuleVersionsRequest
	testCases := []struct {
		msg      string
		malleate func()
		single   bool
		expPass  bool
	}{
		{
			msg: "test full query",
			malleate: func() {
				msg, _ := mupgrade.MakeQueryAllModuleVersionMsg()
				req = &msg
			},
			single:  false,
			expPass: true,
		},
		{
			msg: "test single module",
			malleate: func() {
				queryModulesVersionMsg := types.QueryModulesVersionMsg{
					ModuleName: "bank",
				}
				msg, _ := mupgrade.MakeQueryModuleVersionMsg(queryModulesVersionMsg)
				req = &msg
			},
			single:  true,
			expPass: true,
		},
		{
			msg: "test non-existent module",
			malleate: func() {
				queryModulesVersionMsg := types.QueryModulesVersionMsg{
					ModuleName: "abcdefg",
				}
				msg, _ := mupgrade.MakeQueryModuleVersionMsg(queryModulesVersionMsg)
				req = &msg
			},
			single:  true,
			expPass: false,
		},
	}

	vm := suite.app.UpgradeKeeper.GetModuleVersionMap(suite.ctx)
	mv := suite.app.UpgradeKeeper.GetModuleVersions(suite.ctx)

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest() // reset
			tc.malleate()

			res, err := suite.queryClient.ModuleVersions(context.Background(), req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)

				if tc.single {
					// test that the single module response is valid
					suite.Require().Len(res.ModuleVersions, 1)
					// make sure we got the right values
					suite.Require().Equal(vm[req.ModuleName], res.ModuleVersions[0].Version)
					suite.Require().Equal(req.ModuleName, res.ModuleVersions[0].Name)
				} else {
					// check that the full response is valid
					suite.Require().NotEmpty(res.ModuleVersions)
					suite.Require().Equal(len(mv), len(res.ModuleVersions))
					for i, v := range res.ModuleVersions {
						suite.Require().Equal(mv[i].Version, v.Version)
						suite.Require().Equal(mv[i].Name, v.Name)
					}
				}
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func TestUpgradeTestSuite(t *testing.T) {
	suite.Run(t, new(UpgradeTestSuite))
}
