package authz_test

import (
	"context"
	"fmt"
	"time"

	mauthz "github.com/xpladev/xpla.go/core/authz"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
)

func (suite *TestSuite) TestGRPCQueryAuthorization() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	var (
		req              *authz.QueryGrantsRequest
		expAuthorization authz.Authorization
	)

	m1, _ := key.NewMnemonic()
	m2, _ := key.NewMnemonic()

	k1, _ := key.NewPrivKey(m1)
	k2, _ := key.NewPrivKey(m2)

	granter, _ := util.GetAddrByPrivKey(k1)
	grantee, _ := util.GetAddrByPrivKey(k2)

	testCases := []struct {
		msg      string
		malleate func(require *require.Assertions)
		expError string
		postTest func(require *require.Assertions, res *authz.QueryGrantsResponse)
	}{
		{
			"Success",
			func(require *require.Assertions) {
				now := ctx.BlockHeader().Time
				newCoins := sdk.NewCoins(sdk.NewInt64Coin("axpla", 100))
				expAuthorization = &banktypes.SendAuthorization{SpendLimit: newCoins}
				err := app.AuthzKeeper.SaveGrant(ctx, grantee, granter, expAuthorization, now.Add(time.Hour))
				require.NoError(err)

				queryAuthzGrantMsg := types.QueryAuthzGrantMsg{
					Granter: granter.String(),
					Grantee: grantee.String(),
					MsgType: expAuthorization.MsgTypeURL(),
				}

				msg, err := mauthz.MakeQueryAuthzGrantsMsg(queryAuthzGrantMsg)
				require.NoError(err)
				req = &msg
			},
			"",
			func(require *require.Assertions, res *authz.QueryGrantsResponse) {
				var auth authz.Authorization
				require.Equal(1, len(res.Grants))
				err := suite.app.InterfaceRegistry().UnpackAny(res.Grants[0].Authorization, &auth)
				require.NoError(err)
				require.NotNil(auth)
				require.Equal(auth.String(), expAuthorization.String())
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			require := suite.Require()
			tc.malleate(require)
			result, err := queryClient.Grants(context.Background(), req)
			if tc.expError == "" {
				require.NoError(err)
			} else {
				require.Error(err)
				require.Contains(err.Error(), tc.expError)
			}
			tc.postTest(require, result)
		})
	}
}

func (suite *TestSuite) TestGRPCQueryAuthorizations() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	var (
		req              *authz.QueryGrantsRequest
		expAuthorization authz.Authorization
	)

	m1, _ := key.NewMnemonic()
	m2, _ := key.NewMnemonic()

	k1, _ := key.NewPrivKey(m1)
	k2, _ := key.NewPrivKey(m2)

	granter, _ := util.GetAddrByPrivKey(k1)
	grantee, _ := util.GetAddrByPrivKey(k2)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		postTest func(res *authz.QueryGrantsResponse)
	}{
		{
			"Success",
			func() {
				now := ctx.BlockHeader().Time
				newCoins := sdk.NewCoins(sdk.NewInt64Coin("axpla", 100))
				expAuthorization = &banktypes.SendAuthorization{SpendLimit: newCoins}
				err := app.AuthzKeeper.SaveGrant(ctx, grantee, granter, expAuthorization, now.Add(time.Hour))
				suite.Require().NoError(err)

				queryAuthzGrantMsg := types.QueryAuthzGrantMsg{
					Granter: granter.String(),
					Grantee: grantee.String(),
				}

				msg, err := mauthz.MakeQueryAuthzGrantsMsg(queryAuthzGrantMsg)
				suite.Require().NoError(err)
				req = &msg
			},
			true,
			func(res *authz.QueryGrantsResponse) {
				var auth authz.Authorization
				suite.Require().Equal(1, len(res.Grants))
				err := suite.app.InterfaceRegistry().UnpackAny(res.Grants[0].Authorization, &auth)
				suite.Require().NoError(err)
				suite.Require().NotNil(auth)
				suite.Require().Equal(auth.String(), expAuthorization.String())
			},
		},
	}
	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()
			result, err := queryClient.Grants(context.Background(), req)
			if testCase.expPass {
				suite.Require().NoError(err)
			} else {
				suite.Require().Error(err)
			}
			testCase.postTest(result)
		})
	}
}

func (suite *TestSuite) TestGRPCQueryGranterGrants() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	var (
		req              *authz.QueryGranterGrantsRequest
		expAuthorization authz.Authorization
	)

	m1, _ := key.NewMnemonic()
	m2, _ := key.NewMnemonic()

	k1, _ := key.NewPrivKey(m1)
	k2, _ := key.NewPrivKey(m2)

	granter, _ := util.GetAddrByPrivKey(k1)
	grantee, _ := util.GetAddrByPrivKey(k2)

	testCases := []struct {
		msg      string
		preRun   func()
		expError bool
		parse    func()
		postTest func(res *authz.QueryGranterGrantsResponse)
	}{
		{
			"valid case, single authorization",
			func() {
				now := ctx.BlockHeader().Time
				newCoins := sdk.NewCoins(sdk.NewInt64Coin("axpla", 100))
				expAuthorization = &banktypes.SendAuthorization{SpendLimit: newCoins}
				err := app.AuthzKeeper.SaveGrant(ctx, grantee, granter, expAuthorization, now.Add(time.Hour))
				suite.Require().NoError(err)
			},
			false,
			func() {
				queryAuthzGrantMsg := types.QueryAuthzGrantMsg{
					Granter: granter.String(),
				}

				msg, err := mauthz.MakeQueryAuthzGrantsByGranterMsg(queryAuthzGrantMsg)
				suite.Require().NoError(err)
				req = &msg
			},
			func(res *authz.QueryGranterGrantsResponse) {
				var auth authz.Authorization
				suite.Require().Equal(1, len(res.Grants))
				err := suite.app.InterfaceRegistry().UnpackAny(res.Grants[0].Authorization, &auth)
				suite.Require().NoError(err)
				suite.Require().NotNil(auth)
				suite.Require().Equal(auth.String(), expAuthorization.String())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.preRun()
			tc.parse()
			result, err := queryClient.GranterGrants(context.Background(), req)
			if tc.expError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
			tc.postTest(result)
		})
	}
}

func (suite *TestSuite) TestGRPCQueryGranteeGrants() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	var (
		req              *authz.QueryGranteeGrantsRequest
		expAuthorization authz.Authorization
	)

	m1, _ := key.NewMnemonic()
	m2, _ := key.NewMnemonic()

	k1, _ := key.NewPrivKey(m1)
	k2, _ := key.NewPrivKey(m2)

	granter, _ := util.GetAddrByPrivKey(k1)
	grantee, _ := util.GetAddrByPrivKey(k2)

	testCases := []struct {
		msg      string
		preRun   func()
		expError bool
		parse    func()
		postTest func(res *authz.QueryGranteeGrantsResponse)
	}{
		{
			"valid case, single authorization",
			func() {
				now := ctx.BlockHeader().Time
				newCoins := sdk.NewCoins(sdk.NewInt64Coin("axpla", 100))
				expAuthorization = &banktypes.SendAuthorization{SpendLimit: newCoins}
				err := app.AuthzKeeper.SaveGrant(ctx, grantee, granter, expAuthorization, now.Add(time.Hour))
				suite.Require().NoError(err)
			},
			false,
			func() {
				queryAuthzGrantMsg := types.QueryAuthzGrantMsg{
					Grantee: grantee.String(),
				}

				msg, err := mauthz.MakeQueryAuthzGrantsByGranteeMsg(queryAuthzGrantMsg)
				suite.Require().NoError(err)
				req = &msg
			},
			func(res *authz.QueryGranteeGrantsResponse) {
				var auth authz.Authorization
				suite.Require().Equal(1, len(res.Grants))
				err := suite.app.InterfaceRegistry().UnpackAny(res.Grants[0].Authorization, &auth)
				suite.Require().NoError(err)
				suite.Require().NotNil(auth)
				suite.Require().Equal(auth.String(), expAuthorization.String())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.preRun()
			tc.parse()
			result, err := queryClient.GranteeGrants(context.Background(), req)
			if tc.expError {
				suite.Require().Error(err)
			} else {
				suite.Require().NoError(err)
			}
			tc.postTest(result)
		})
	}
}
