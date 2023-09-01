package staking_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/xpladev/xpla.go/core"
	mstaking "github.com/xpladev/xpla.go/core/staking"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type KeeperTestSuite struct {
	suite.Suite

	app         *xapp.XplaApp
	ctx         sdk.Context
	addrs       []sdk.AccAddress
	vals        []stakingtypes.Validator
	queryClient stakingtypes.QueryClient
}

func (suite *KeeperTestSuite) SetupTest() {
	app := testutil.Setup(false, 5)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	querier := keeper.Querier{Keeper: app.StakingKeeper}

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	stakingtypes.RegisterQueryServer(queryHelper, querier)
	queryClient := stakingtypes.NewQueryClient(queryHelper)

	addrs, _, validators := createValidators(suite.T(), ctx, app, []int64{9, 8, 7})
	header := tmproto.Header{
		ChainID: "HelloChain",
		Height:  5,
	}

	// sort a copy of the validators, so that original validators does not
	// have its order changed
	sortedVals := make([]stakingtypes.Validator, len(validators))
	copy(sortedVals, validators)
	hi := stakingtypes.NewHistoricalInfo(header, sortedVals, app.StakingKeeper.PowerReduction(ctx))
	app.StakingKeeper.SetHistoricalInfo(ctx, 5, &hi)

	suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals = app, ctx, queryClient, addrs, validators
}

func (suite *KeeperTestSuite) TestGRPCQueryValidator() {
	app, ctx, queryClient, vals := suite.app, suite.ctx, suite.queryClient, suite.vals
	validator, found := app.StakingKeeper.GetValidator(ctx, vals[0].GetOperator())
	suite.True(found)
	var req *stakingtypes.QueryValidatorRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request",
			func() {
				queryValidatorMsg := types.QueryValidatorMsg{
					ValidatorAddr: vals[0].OperatorAddress,
				}
				msg, _ := mstaking.MakeQueryValidatorMsg(queryValidatorMsg)
				req = &msg
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Validator(context.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.True(validator.Equal(&res.Validator))
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegation() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc, addrAcc1 := addrs[0], addrs[1]
	addrVal := vals[0].OperatorAddress
	valAddr, err := sdk.ValAddressFromBech32(addrVal)
	suite.NoError(err)
	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrAcc, valAddr)
	suite.True(found)
	var req *stakingtypes.QueryDelegationRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"invalid validator, delegator pair",
			func() {
				queryDelegationMsg := types.QueryDelegationMsg{
					DelegatorAddr: addrAcc1.String(),
					ValidatorAddr: addrVal,
				}
				msg, _ := mstaking.MakeQueryDelegationMsg(queryDelegationMsg)
				req = &msg
			},
			false,
		},
		{
			"valid request",
			func() {
				queryDelegationMsg := types.QueryDelegationMsg{
					DelegatorAddr: addrAcc.String(),
					ValidatorAddr: addrVal,
				}
				msg, _ := mstaking.MakeQueryDelegationMsg(queryDelegationMsg)
				req = &msg
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Delegation(context.Background(), req)
			if tc.expPass {
				suite.Equal(delegation.ValidatorAddress, res.DelegationResponse.Delegation.ValidatorAddress)
				suite.Equal(delegation.DelegatorAddress, res.DelegationResponse.Delegation.DelegatorAddress)
				suite.Equal(sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), res.DelegationResponse.Balance)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc := addrs[0]
	addrVal1 := vals[0].OperatorAddress
	valAddr, err := sdk.ValAddressFromBech32(addrVal1)
	suite.NoError(err)
	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrAcc, valAddr)
	suite.True(found)
	var req *stakingtypes.QueryDelegatorDelegationsRequest

	testCases := []struct {
		msg       string
		malleate  func()
		onSuccess func(suite *KeeperTestSuite, response *stakingtypes.QueryDelegatorDelegationsResponse)
		expErr    bool
	}{
		{
			"valid request with no delegations",
			func() {
				queryDelegationMsg := types.QueryDelegationMsg{
					DelegatorAddr: addrs[4].String(),
				}
				msg, _ := mstaking.MakeQueryDelegationsMsg(queryDelegationMsg)
				req = &msg
			},
			func(suite *KeeperTestSuite, response *stakingtypes.QueryDelegatorDelegationsResponse) {
				suite.Equal(uint64(0), response.Pagination.Total)
				suite.Len(response.DelegationResponses, 0)
			},
			false,
		},
		{
			"valid request",
			func() {
				pagination := types.Pagination{
					Limit:      1,
					CountTotal: true,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				queryDelegationMsg := types.QueryDelegationMsg{
					DelegatorAddr: addrAcc.String(),
				}
				msg, _ := mstaking.MakeQueryDelegationsMsg(queryDelegationMsg)
				req = &msg
			},
			func(suite *KeeperTestSuite, response *stakingtypes.QueryDelegatorDelegationsResponse) {
				suite.Equal(uint64(2), response.Pagination.Total)
				suite.Len(response.DelegationResponses, 1)
				suite.Equal(sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), response.DelegationResponses[0].Balance)
			},
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorDelegations(context.Background(), req)
			if tc.expErr {
				suite.Error(err)
			} else {
				suite.NoError(err)
				tc.onSuccess(suite, res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryValidatorDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc := addrs[0]
	addrVal1 := vals[1].OperatorAddress
	valAddrs := testutil.ConvertAddrsToValAddrs(addrs)
	addrVal2 := valAddrs[4]
	valAddr, err := sdk.ValAddressFromBech32(addrVal1)
	suite.NoError(err)
	delegation, found := app.StakingKeeper.GetDelegation(ctx, addrAcc, valAddr)
	suite.True(found)

	var req *stakingtypes.QueryValidatorDelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"invalid validator delegator pair",
			func() {
				queryDelegationMsg := types.QueryDelegationMsg{
					ValidatorAddr: addrVal2.String(),
				}
				msg, _ := mstaking.MakeQueryDelegationsToMsg(queryDelegationMsg)
				req = &msg
			},
			false,
			false,
		},
		{
			"valid request",
			func() {
				pagination := types.Pagination{
					Limit:      1,
					CountTotal: true,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				queryDelegationMsg := types.QueryDelegationMsg{
					ValidatorAddr: addrVal1,
				}
				msg, _ := mstaking.MakeQueryDelegationsToMsg(queryDelegationMsg)
				req = &msg
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.ValidatorDelegations(context.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.DelegationResponses, 1)
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(2), res.Pagination.Total)
				suite.Equal(addrVal1, res.DelegationResponses[0].Delegation.ValidatorAddress)
				suite.Equal(sdk.NewCoin(sdk.DefaultBondDenom, delegation.Shares.TruncateInt()), res.DelegationResponses[0].Balance)
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.DelegationResponses)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryUnbondingDelegation() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc2 := addrs[1]
	addrVal2 := vals[1].OperatorAddress

	unbondingTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	valAddr, err1 := sdk.ValAddressFromBech32(addrVal2)
	suite.NoError(err1)
	_, err := app.StakingKeeper.Undelegate(ctx, addrAcc2, valAddr, unbondingTokens.ToDec())
	suite.NoError(err)

	unbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrAcc2, valAddr)
	suite.True(found)
	var req *stakingtypes.QueryUnbondingDelegationRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"invalid request",
			func() {
				queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{}
				msg, _ := mstaking.MakeQueryUnbondingDelegationMsg(queryUnbondingDelegationMsg)
				req = &msg
			},
			false,
		},
		{
			"valid request",
			func() {
				queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{
					DelegatorAddr: addrAcc2.String(),
					ValidatorAddr: addrVal2,
				}
				msg, _ := mstaking.MakeQueryUnbondingDelegationMsg(queryUnbondingDelegationMsg)
				req = &msg
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.UnbondingDelegation(context.Background(), req)
			if tc.expPass {
				suite.NotNil(res)
				suite.Equal(unbond, res.Unbond)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryDelegatorUnbondingDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc, addrAcc1 := addrs[0], addrs[1]
	addrVal, addrVal2 := vals[0].OperatorAddress, vals[1].OperatorAddress

	unbondingTokens := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	valAddr1, err1 := sdk.ValAddressFromBech32(addrVal)
	suite.NoError(err1)
	_, err := app.StakingKeeper.Undelegate(ctx, addrAcc, valAddr1, unbondingTokens.ToDec())
	suite.NoError(err)
	valAddr2, err1 := sdk.ValAddressFromBech32(addrVal2)
	suite.NoError(err1)
	_, err = app.StakingKeeper.Undelegate(ctx, addrAcc, valAddr2, unbondingTokens.ToDec())
	suite.NoError(err)

	unbond, found := app.StakingKeeper.GetUnbondingDelegation(ctx, addrAcc, valAddr1)
	suite.True(found)
	var req *stakingtypes.QueryDelegatorUnbondingDelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"empty request",
			func() {
				queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{}
				msg, _ := mstaking.MakeQueryUnbondingDelegationsMsg(queryUnbondingDelegationMsg)
				req = &msg
			},
			false,
			true,
		},
		{
			"invalid request",
			func() {
				queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{
					DelegatorAddr: addrAcc1.String(),
				}
				msg, _ := mstaking.MakeQueryUnbondingDelegationsMsg(queryUnbondingDelegationMsg)
				req = &msg
			},
			false,
			false,
		},
		{
			"valid request",
			func() {
				pagination := types.Pagination{
					Limit:      1,
					CountTotal: true,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{
					DelegatorAddr: addrAcc.String(),
				}
				msg, _ := mstaking.MakeQueryUnbondingDelegationsMsg(queryUnbondingDelegationMsg)
				req = &msg
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.DelegatorUnbondingDelegations(context.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.NotNil(res.Pagination.NextKey)
				suite.Equal(uint64(2), res.Pagination.Total)
				suite.Len(res.UnbondingResponses, 1)
				suite.Equal(unbond, res.UnbondingResponses[0])
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.UnbondingResponses)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryPoolParameters() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	bondDenom := sdk.DefaultBondDenom

	// Query pool
	res, err := queryClient.Pool(context.Background(), &stakingtypes.QueryPoolRequest{})
	suite.NoError(err)
	bondedPool := app.StakingKeeper.GetBondedPool(ctx)
	notBondedPool := app.StakingKeeper.GetNotBondedPool(ctx)
	suite.Equal(app.BankKeeper.GetBalance(ctx, notBondedPool.GetAddress(), bondDenom).Amount, res.Pool.NotBondedTokens)
	suite.Equal(app.BankKeeper.GetBalance(ctx, bondedPool.GetAddress(), bondDenom).Amount, res.Pool.BondedTokens)

	msg, _ := mstaking.MakeQueryStakingParamsMsg()
	// Query Params
	resp, err := queryClient.Params(context.Background(), &msg)
	suite.NoError(err)
	suite.Equal(app.StakingKeeper.GetParams(ctx), resp.Params)
}

func (suite *KeeperTestSuite) TestGRPCQueryHistoricalInfo() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	hi, found := app.StakingKeeper.GetHistoricalInfo(ctx, 5)
	suite.True(found)

	var req *stakingtypes.QueryHistoricalInfoRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				historicalInfoMsg := types.HistoricalInfoMsg{}
				msg, _ := mstaking.MakeHistoricalInfoMsg(historicalInfoMsg)
				req = &msg
			},
			false,
		},
		{
			"invalid request with negative height",
			func() {
				historicalInfoMsg := types.HistoricalInfoMsg{
					Height: "-1",
				}
				msg, _ := mstaking.MakeHistoricalInfoMsg(historicalInfoMsg)
				req = &msg
			},
			false,
		},
		{
			"valid request with old height",
			func() {
				historicalInfoMsg := types.HistoricalInfoMsg{
					Height: "4",
				}
				msg, _ := mstaking.MakeHistoricalInfoMsg(historicalInfoMsg)
				req = &msg
			},
			false,
		},
		{
			"valid request with current height",
			func() {
				historicalInfoMsg := types.HistoricalInfoMsg{
					Height: "5",
				}
				msg, _ := mstaking.MakeHistoricalInfoMsg(historicalInfoMsg)
				req = &msg
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.HistoricalInfo(context.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.NotNil(res)
				suite.True(hi.Equal(res.Hist))
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryRedelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals

	addrAcc, addrAcc1 := addrs[0], addrs[1]
	valAddrs := testutil.ConvertAddrsToValAddrs(addrs)
	val1, val2, val3, val4 := vals[0], vals[1], valAddrs[3], valAddrs[4]
	delAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 1)
	_, err := app.StakingKeeper.Delegate(ctx, addrAcc1, delAmount, stakingtypes.Unbonded, val1, true)
	suite.NoError(err)
	applyValidatorSetUpdates(suite.T(), ctx, app.StakingKeeper, -1)

	rdAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 1)
	_, err = app.StakingKeeper.BeginRedelegation(ctx, addrAcc1, val1.GetOperator(), val2.GetOperator(), rdAmount.ToDec())
	suite.NoError(err)
	applyValidatorSetUpdates(suite.T(), ctx, app.StakingKeeper, -1)

	redel, found := app.StakingKeeper.GetRedelegation(ctx, addrAcc1, val1.GetOperator(), val2.GetOperator())
	suite.True(found)

	var req *stakingtypes.QueryRedelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
		expErr   bool
	}{
		{
			"request redelegations for non existent addr",
			func() {
				queryRedelegationMsg := types.QueryRedelegationMsg{
					DelegatorAddr: addrAcc.String(),
				}
				msg, _ := mstaking.MakeQueryRedelegationMsg(queryRedelegationMsg)
				req = &msg
			},
			false,
			false,
		},
		{
			"request redelegations with non existent pairs",
			func() {
				queryRedelegationMsg := types.QueryRedelegationMsg{
					DelegatorAddr:    addrAcc.String(),
					SrcValidatorAddr: val3.String(),
					DstValidatorAddr: val4.String(),
				}
				msg, _ := mstaking.MakeQueryRedelegationMsg(queryRedelegationMsg)
				req = &msg
			},
			false,
			true,
		},
		{
			"request redelegations with delegatoraddr, sourceValAddr, destValAddr",
			func() {
				pagination := types.Pagination{
					Limit:      1,
					CountTotal: true,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				queryRedelegationMsg := types.QueryRedelegationMsg{
					DelegatorAddr:    addrAcc1.String(),
					SrcValidatorAddr: val1.OperatorAddress,
					DstValidatorAddr: val2.OperatorAddress,
				}
				msg, _ := mstaking.MakeQueryRedelegationMsg(queryRedelegationMsg)
				req = &msg
			},
			true,
			false,
		},
		{
			"request redelegations with delegatoraddr and sourceValAddr",
			func() {
				pagination := types.Pagination{
					Limit:      1,
					CountTotal: true,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				queryRedelegationMsg := types.QueryRedelegationMsg{
					DelegatorAddr:    addrAcc1.String(),
					SrcValidatorAddr: val1.OperatorAddress,
				}
				msg, _ := mstaking.MakeQueryRedelegationMsg(queryRedelegationMsg)
				req = &msg
			},
			true,
			false,
		},
		{
			"query redelegations with sourceValAddr only",
			func() {
				pagination := types.Pagination{
					Limit:      1,
					CountTotal: true,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				queryRedelegationMsg := types.QueryRedelegationMsg{
					SrcValidatorAddr: val1.GetOperator().String(),
				}
				msg, _ := mstaking.MakeQueryRedelegationMsg(queryRedelegationMsg)
				req = &msg
			},
			true,
			false,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.Redelegations(context.Background(), req)
			if tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Len(res.RedelegationResponses, len(redel.Entries))
				suite.Equal(redel.DelegatorAddress, res.RedelegationResponses[0].Redelegation.DelegatorAddress)
				suite.Equal(redel.ValidatorSrcAddress, res.RedelegationResponses[0].Redelegation.ValidatorSrcAddress)
				suite.Equal(redel.ValidatorDstAddress, res.RedelegationResponses[0].Redelegation.ValidatorDstAddress)
				suite.Len(redel.Entries, len(res.RedelegationResponses[0].Entries))
			} else if !tc.expPass && !tc.expErr {
				suite.NoError(err)
				suite.Nil(res.RedelegationResponses)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestGRPCQueryValidatorUnbondingDelegations() {
	app, ctx, queryClient, addrs, vals := suite.app, suite.ctx, suite.queryClient, suite.addrs, suite.vals
	addrAcc1, _ := addrs[0], addrs[1]
	val1 := vals[0]

	// undelegate
	undelAmount := app.StakingKeeper.TokensFromConsensusPower(ctx, 2)
	_, err := app.StakingKeeper.Undelegate(ctx, addrAcc1, val1.GetOperator(), undelAmount.ToDec())
	suite.NoError(err)
	applyValidatorSetUpdates(suite.T(), ctx, app.StakingKeeper, -1)

	var req *stakingtypes.QueryValidatorUnbondingDelegationsRequest
	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty request",
			func() {
				queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{}
				msg, _ := mstaking.MakeQueryUnbondingDelegationsFromMsg(queryUnbondingDelegationMsg)
				req = &msg
			},
			false,
		},
		{
			"valid request",
			func() {
				pagination := types.Pagination{
					Limit:      1,
					CountTotal: true,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{
					ValidatorAddr: val1.GetOperator().String(),
				}
				msg, _ := mstaking.MakeQueryUnbondingDelegationsFromMsg(queryUnbondingDelegationMsg)
				req = &msg
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			res, err := queryClient.ValidatorUnbondingDelegations(context.Background(), req)
			if tc.expPass {
				suite.NoError(err)
				suite.Equal(uint64(1), res.Pagination.Total)
				suite.Equal(1, len(res.UnbondingResponses))
				suite.Equal(res.UnbondingResponses[0].ValidatorAddress, val1.OperatorAddress)
			} else {
				suite.Error(err)
				suite.Nil(res)
			}
		})
	}
}

func createValidators(t *testing.T, ctx sdk.Context, app *xapp.XplaApp, powers []int64) ([]sdk.AccAddress, []sdk.ValAddress, []stakingtypes.Validator) {
	addrs := testutil.AddTestAddrsIncremental(app, ctx, 5, app.StakingKeeper.TokensFromConsensusPower(ctx, 300))
	valAddrs := testutil.ConvertAddrsToValAddrs(addrs)
	pks := testutil.CreateTestPubKeys(5)
	cdc := util.MakeEncodingConfig().Marshaler
	app.StakingKeeper = keeper.NewKeeper(
		cdc,
		app.GetKey(stakingtypes.StoreKey),
		app.AccountKeeper,
		app.BankKeeper,
		app.GetSubspace(stakingtypes.ModuleName),
	)

	val1 := teststaking.NewValidator(t, valAddrs[0], pks[0])
	val2 := teststaking.NewValidator(t, valAddrs[1], pks[1])
	vals := []stakingtypes.Validator{val1, val2}

	app.StakingKeeper.SetValidator(ctx, val1)
	app.StakingKeeper.SetValidator(ctx, val2)
	app.StakingKeeper.SetValidatorByConsAddr(ctx, val1)
	app.StakingKeeper.SetValidatorByConsAddr(ctx, val2)
	app.StakingKeeper.SetNewValidatorByPowerIndex(ctx, val1)
	app.StakingKeeper.SetNewValidatorByPowerIndex(ctx, val2)

	_, err := app.StakingKeeper.Delegate(ctx, addrs[0], app.StakingKeeper.TokensFromConsensusPower(ctx, powers[0]), stakingtypes.Unbonded, val1, true)
	require.NoError(t, err)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[1], app.StakingKeeper.TokensFromConsensusPower(ctx, powers[1]), stakingtypes.Unbonded, val2, true)
	require.NoError(t, err)
	_, err = app.StakingKeeper.Delegate(ctx, addrs[0], app.StakingKeeper.TokensFromConsensusPower(ctx, powers[2]), stakingtypes.Unbonded, val2, true)
	require.NoError(t, err)
	applyValidatorSetUpdates(t, ctx, app.StakingKeeper, -1)

	return addrs, valAddrs, vals
}

func applyValidatorSetUpdates(t *testing.T, ctx sdk.Context, k keeper.Keeper, expectedUpdatesLen int) []abci.ValidatorUpdate {
	updates, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
	require.NoError(t, err)
	if expectedUpdatesLen >= 0 {
		require.Equal(t, expectedUpdatesLen, len(updates), "%v", updates)
	}
	return updates
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
