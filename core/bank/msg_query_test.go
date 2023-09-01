package bank_test

import (
	"context"
	"fmt"

	mbank "github.com/xpladev/xpla.go/core/bank"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

const (
	fooDenom = "foo"
	barDenom = "bar"
)

func (suite *TestSuite) TestQueryBalance() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	mnemonic, _ := key.NewMnemonic()
	privKey, _ := key.NewPrivKey(mnemonic)
	addr, _ := util.GetAddrByPrivKey(privKey)

	xplaCoins := sdk.NewCoins(sdk.Coin{
		Denom:  "axpla",
		Amount: sdk.NewInt(10),
	})
	account := app.AccountKeeper.NewAccountWithAddress(ctx, addr)

	app.AccountKeeper.SetAccount(ctx, account)
	suite.Require().NoError(testutil.FundAccount(app.BankKeeper, ctx, account.GetAddress(), xplaCoins))

	bankBalancesMsg := types.BankBalancesMsg{}
	msg, _ := mbank.MakeBankAllBalancesMsg(bankBalancesMsg)
	_, err := queryClient.AllBalances(context.Background(), &msg)
	suite.Require().Error(err)

	bankBalancesMsg = types.BankBalancesMsg{
		Address: addr.String(),
	}
	msg, err = mbank.MakeBankAllBalancesMsg(bankBalancesMsg)
	suite.Require().NoError(err)
	res, err := queryClient.AllBalances(context.Background(), &msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	bankBalanceMsg := types.BankBalancesMsg{
		Address: addr.String(),
		Denom:   "axpla",
	}
	msg2, err := mbank.MakeBankBalanceMsg(bankBalanceMsg)
	suite.Require().NoError(err)
	res2, err := queryClient.Balance(context.Background(), &msg2)
	suite.Require().NoError(err)
	suite.Require().NotNil(res2)

	origCoins := sdk.NewCoins(newFooCoin(50), newBarCoin(30))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)

	app.AccountKeeper.SetAccount(ctx, acc)
	suite.Require().NoError(testutil.FundAccount(app.BankKeeper, ctx, acc.GetAddress(), origCoins))

	res, err = queryClient.AllBalances(context.Background(), &msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)
}

func (suite *TestSuite) TestQueryTotalSupply() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient
	expectedTotalSupply := sdk.NewCoins(sdk.NewInt64Coin("test", 400000000))
	suite.
		Require().
		NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, expectedTotalSupply))

	msg, _ := mbank.MakeTotalSupplyMsg()
	res, err := queryClient.TotalSupply(context.Background(), &msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	suite.Require().Equal(expectedTotalSupply, res.Supply)
}

func (suite *TestSuite) TestQueryTotalSupplyOf() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	test1Supply := sdk.NewInt64Coin("test1", 4000000)
	test2Supply := sdk.NewInt64Coin("test2", 700000000)
	expectedTotalSupply := sdk.NewCoins(test1Supply, test2Supply)
	suite.
		Require().
		NoError(app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, expectedTotalSupply))

	totalMsg := types.TotalMsg{}
	msg, _ := mbank.MakeSupplyOfMsg(totalMsg)
	_, err := queryClient.SupplyOf(context.Background(), &msg)
	suite.Require().Error(err)

	totalMsg = types.TotalMsg{
		Denom: "test1",
	}
	msg, _ = mbank.MakeSupplyOfMsg(totalMsg)
	res, err := queryClient.SupplyOf(context.Background(), &msg)
	suite.Require().NoError(err)
	suite.Require().NotNil(res)

	suite.Require().Equal(test1Supply, res.Amount)
}

func (suite *TestSuite) TestQueryDenomMetadataRequest() {
	var (
		req         *banktypes.QueryDenomMetadataRequest
		expMetadata = banktypes.Metadata{}
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"success",
			func() {
				expMetadata = banktypes.Metadata{
					Description: "The native staking token of the XPLA.",
					DenomUnits: []*banktypes.DenomUnit{
						{
							Denom:    "axpla",
							Exponent: 0,
							Aliases:  []string{"attoxpla"},
						},
						{
							Denom:    "xpla",
							Exponent: 18,
							Aliases:  []string{"XPLA"},
						},
					},
					Base:    "axpla",
					Display: "xpla",
				}

				suite.app.BankKeeper.SetDenomMetaData(suite.ctx, expMetadata)
				denomMetadataMsg := types.DenomMetadataMsg{
					Denom: expMetadata.Base,
				}
				msg, _ := mbank.MakeDenomMetaDataMsg(denomMetadataMsg)
				req = &msg
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)

			res, err := suite.queryClient.DenomMetadata(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
				suite.Require().Equal(expMetadata, res.Metadata)
			} else {
				suite.Require().Error(err)
			}
		})
	}
}

func newFooCoin(amt int64) sdk.Coin {
	return sdk.NewInt64Coin(fooDenom, amt)
}

func newBarCoin(amt int64) sdk.Coin {
	return sdk.NewInt64Coin(barDenom, amt)
}
