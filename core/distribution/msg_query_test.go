package distribution_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strconv"

	"github.com/xpladev/xpla.go/core"
	mdist "github.com/xpladev/xpla.go/core/distribution"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/teststaking"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var (
	pubKeys, _ = CreateTestPubKeys(5)
	valConsPk1 = pubKeys[0]
)

func (suite *TestSuite) TestGRPCParams() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	var (
		params    disttypes.Params
		req       *disttypes.QueryParamsRequest
		expParams disttypes.Params
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {
				params = disttypes.Params{
					CommunityTax:        sdk.NewDecWithPrec(3, 1),
					BaseProposerReward:  sdk.NewDecWithPrec(2, 1),
					BonusProposerReward: sdk.NewDecWithPrec(1, 1),
					WithdrawAddrEnabled: true,
				}
				app.DistrKeeper.SetParams(ctx, params)

				msg, _ := mdist.MakeQueryDistributionParamsMsg()
				req = &msg
				expParams = params
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			paramsRes, err := queryClient.Params(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(paramsRes)
				suite.Require().Equal(paramsRes.Params, expParams)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func (suite *TestSuite) TestGRPCValidatorOutstandingRewards() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	_, valAddrs := suite.generateAddrsAndValAddrs()
	valCommission := sdk.DecCoins{
		sdk.NewDecCoinFromDec("axpla", sdk.NewDec(300)),
	}

	// set outstanding rewards
	app.DistrKeeper.SetValidatorOutstandingRewards(ctx, valAddrs[0], disttypes.ValidatorOutstandingRewards{Rewards: valCommission})
	rewards := app.DistrKeeper.GetValidatorOutstandingRewards(ctx, valAddrs[0])

	var req *disttypes.QueryValidatorOutstandingRewardsRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {
				validatorOutstandingRewardsMsg := types.ValidatorOutstandingRewardsMsg{
					ValidatorAddr: valAddrs[0].String(),
				}
				msg, _ := mdist.MakeValidatorOutstandingRewardsMsg(validatorOutstandingRewardsMsg)
				req = &msg
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			validatorOutstandingRewards, err := queryClient.ValidatorOutstandingRewards(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(rewards, validatorOutstandingRewards.Rewards)
				suite.Require().Equal(valCommission, validatorOutstandingRewards.Rewards.Rewards)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(validatorOutstandingRewards)
			}
		})
	}
}

func (suite *TestSuite) TestGRPCValidatorCommission() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	_, valAddrs := suite.generateAddrsAndValAddrs()

	commission := sdk.DecCoins{
		{Denom: "axpla", Amount: sdk.NewDec(2)},
	}
	app.DistrKeeper.SetValidatorAccumulatedCommission(ctx, valAddrs[0], disttypes.ValidatorAccumulatedCommission{Commission: commission})

	var req *disttypes.QueryValidatorCommissionRequest

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {
				queryDistCommissionMsg := types.QueryDistCommissionMsg{
					ValidatorAddr: valAddrs[0].String(),
				}
				msg, _ := mdist.MakeQueryDistCommissionMsg(queryDistCommissionMsg)
				req = &msg
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			commissionRes, err := queryClient.ValidatorCommission(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(commissionRes)
				suite.Require().Equal(commissionRes.Commission.Commission, commission)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(commissionRes)
			}
		})
	}
}

func (suite *TestSuite) TestGRPCValidatorSlashes() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	_, valAddrs := suite.generateAddrsAndValAddrs()

	slashes := []disttypes.ValidatorSlashEvent{
		disttypes.NewValidatorSlashEvent(3, sdk.NewDecWithPrec(5, 1)),
		disttypes.NewValidatorSlashEvent(5, sdk.NewDecWithPrec(5, 1)),
		disttypes.NewValidatorSlashEvent(7, sdk.NewDecWithPrec(5, 1)),
		disttypes.NewValidatorSlashEvent(9, sdk.NewDecWithPrec(5, 1)),
	}

	for i, slash := range slashes {
		app.DistrKeeper.SetValidatorSlashEvent(ctx, valAddrs[0], uint64(i+2), 0, slash)
	}

	var (
		req    *disttypes.QueryValidatorSlashesRequest
		expRes *disttypes.QueryValidatorSlashesResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"Ending height lesser than start height request",
			func() {

				queryDistSlashesMsg := types.QueryDistSlashesMsg{
					ValidatorAddr: valAddrs[1].String(),
					StartHeight:   "10",
					EndHeight:     "1",
				}
				msg, _ := mdist.MakeQueryDistSlashesMsg(queryDistSlashesMsg)
				req = &msg

				expRes = &disttypes.QueryValidatorSlashesResponse{}
			},
			false,
		},
		{
			"no slash event validator request",
			func() {
				queryDistSlashesMsg := types.QueryDistSlashesMsg{
					ValidatorAddr: valAddrs[1].String(),
					StartHeight:   "1",
					EndHeight:     "10",
				}
				msg, _ := mdist.MakeQueryDistSlashesMsg(queryDistSlashesMsg)
				req = &msg
				expRes = &disttypes.QueryValidatorSlashesResponse{}
			},
			true,
		},
		{
			"request slashes with offset 2 and limit 2",
			func() {
				pagination := types.Pagination{
					Offset: 2,
					Limit:  2,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				queryDistSlashesMsg := types.QueryDistSlashesMsg{
					ValidatorAddr: valAddrs[0].String(),
					StartHeight:   "1",
					EndHeight:     "10",
				}
				msg, _ := mdist.MakeQueryDistSlashesMsg(queryDistSlashesMsg)
				req = &msg

				expRes = &disttypes.QueryValidatorSlashesResponse{
					Slashes: slashes[2:],
				}
			},
			true,
		},
		{
			"request slashes with page limit 3 and count total",
			func() {
				pagination := types.Pagination{
					Limit:      3,
					CountTotal: true,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				queryDistSlashesMsg := types.QueryDistSlashesMsg{
					ValidatorAddr: valAddrs[0].String(),
					StartHeight:   "1",
					EndHeight:     "10",
				}
				msg, _ := mdist.MakeQueryDistSlashesMsg(queryDistSlashesMsg)
				req = &msg

				expRes = &disttypes.QueryValidatorSlashesResponse{
					Slashes: slashes[:3],
				}
			},
			true,
		},
		{
			"request slashes with page limit 4 and count total",
			func() {
				pagination := types.Pagination{
					Limit:      4,
					CountTotal: true,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				queryDistSlashesMsg := types.QueryDistSlashesMsg{
					ValidatorAddr: valAddrs[0].String(),
					StartHeight:   "1",
					EndHeight:     "10",
				}
				msg, _ := mdist.MakeQueryDistSlashesMsg(queryDistSlashesMsg)
				req = &msg

				expRes = &disttypes.QueryValidatorSlashesResponse{
					Slashes: slashes[:4],
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			slashesRes, err := queryClient.ValidatorSlashes(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.GetSlashes(), slashesRes.GetSlashes())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(slashesRes)
			}
		})
	}
}

func (suite *TestSuite) TestGRPCDelegationRewards() {
	app, ctx := suite.app, suite.ctx
	addrs, valAddrs := suite.generateAddrsAndValAddrs()

	tstaking := teststaking.NewHelper(suite.T(), ctx, app.StakingKeeper)
	tstaking.Commission = stakingtypes.NewCommissionRates(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), sdk.NewDec(0))
	tstaking.CreateValidator(valAddrs[0], valConsPk1, sdk.NewInt(100), true)

	staking.EndBlocker(ctx, app.StakingKeeper)
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	disttypes.RegisterQueryServer(queryHelper, app.DistrKeeper)
	queryClient := disttypes.NewQueryClient(queryHelper)

	val := app.StakingKeeper.Validator(ctx, valAddrs[0])

	initial := int64(10)
	tokens := sdk.DecCoins{{Denom: "axpla", Amount: sdk.NewDec(initial)}}
	app.DistrKeeper.AllocateTokensToValidator(ctx, val, tokens)

	// test command delegation rewards grpc
	var (
		req    *disttypes.QueryDelegationRewardsRequest
		expRes *disttypes.QueryDelegationRewardsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty delegator request",
			func() {
				queryDistRewardsMsg := types.QueryDistRewardsMsg{
					ValidatorAddr: valAddrs[0].String(),
					DelegatorAddr: "",
				}
				msg, _ := mdist.MakeQueryDistRewardsMsg(queryDistRewardsMsg)
				req = &msg
			},
			false,
		},
		{
			"empty validator request",
			func() {
				queryDistRewardsMsg := types.QueryDistRewardsMsg{
					ValidatorAddr: "",
					DelegatorAddr: addrs[1].String(),
				}
				msg, _ := mdist.MakeQueryDistRewardsMsg(queryDistRewardsMsg)
				req = &msg
			},
			false,
		},
		{
			"request with wrong delegator and validator",
			func() {
				queryDistRewardsMsg := types.QueryDistRewardsMsg{
					ValidatorAddr: valAddrs[1].String(),
					DelegatorAddr: addrs[1].String(),
				}
				msg, _ := mdist.MakeQueryDistRewardsMsg(queryDistRewardsMsg)
				req = &msg
			},
			false,
		},
		{
			"success",
			func() {
				queryDistRewardsMsg := types.QueryDistRewardsMsg{
					ValidatorAddr: valAddrs[0].String(),
					DelegatorAddr: addrs[0].String(),
				}
				msg, _ := mdist.MakeQueryDistRewardsMsg(queryDistRewardsMsg)
				req = &msg

				expRes = &disttypes.QueryDelegationRewardsResponse{
					Rewards: sdk.DecCoins{{Denom: "axpla", Amount: sdk.NewDec(initial / 2)}},
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			rewards, err := queryClient.DelegationRewards(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes, rewards)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(rewards)
			}
		})
	}
}

func (suite *TestSuite) TestGRPCCommunityPool() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	addrs, _ := suite.generateAddrsAndValAddrs()

	var (
		req     *disttypes.QueryCommunityPoolRequest
		expPool *disttypes.QueryCommunityPoolResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"valid request empty community pool",
			func() {

				msg, _ := mdist.MakeQueryCommunityPoolMsg()
				req = &msg

				expPool = &disttypes.QueryCommunityPoolResponse{}
			},
			true,
		},
		{
			"valid request",
			func() {
				amount := sdk.NewCoins(sdk.NewInt64Coin("axpla", 100))
				suite.Require().NoError(testutil.FundAccount(app.BankKeeper, ctx, addrs[0], amount))

				err := app.DistrKeeper.FundCommunityPool(ctx, amount, addrs[0])
				suite.Require().Nil(err)

				msg, _ := mdist.MakeQueryCommunityPoolMsg()
				req = &msg

				expPool = &disttypes.QueryCommunityPoolResponse{Pool: sdk.NewDecCoinsFromCoins(amount...)}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			pool, err := queryClient.CommunityPool(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expPool, pool)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(pool)
			}
		})
	}
}

func (suite *TestSuite) generateAddrsAndValAddrs() ([]sdk.AccAddress, []sdk.ValAddress) {
	accNum := 2
	accAmt := sdk.NewInt(1000000000)
	testAddrs := createRandomAccounts(accNum)

	initCoins := sdk.NewCoins(sdk.NewCoin(suite.app.StakingKeeper.BondDenom(suite.ctx), accAmt))

	for _, addr := range testAddrs {
		suite.initAccountWithCoins(addr, initCoins)
	}

	valAddrs := ConvertAddrsToValAddrs(testAddrs)

	return testAddrs, valAddrs
}

func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey()
		testAddrs[i] = sdk.AccAddress(pk.PubKey().Address())
	}

	return testAddrs
}

func (suite *TestSuite) initAccountWithCoins(addr sdk.AccAddress, coins sdk.Coins) {
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	if err != nil {
		panic(err)
	}

	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, coins)
	if err != nil {
		panic(err)
	}
}

func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addrs))

	for i, addr := range addrs {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs
}

func CreateTestPubKeys(numPubKeys int) ([]cryptotypes.PubKey, error) {
	var publicKeys []cryptotypes.PubKey
	var buffer bytes.Buffer

	// start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") // base pubkey string
		buffer.WriteString(numString)                                                       // adding on final two digits to make pubkeys unique

		p, err := NewPubKeyFromHex(buffer.String())
		if err != nil {
			return nil, err
		}
		publicKeys = append(publicKeys, p)
		buffer.Reset()
	}

	return publicKeys, nil
}

func NewPubKeyFromHex(pk string) (cryptotypes.PubKey, error) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	if len(pkBytes) != ed25519.PubKeySize {
		return nil, util.LogErr(errors.ErrInvalidRequest, "invalid pubkey size")
	}
	return &ed25519.PubKey{Key: pkBytes}, nil
}
