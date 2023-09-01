package gov_test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/xpladev/xpla.go/core"
	mgov "github.com/xpladev/xpla.go/core/gov"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type TestSuite struct {
	suite.Suite

	app         *xapp.XplaApp
	ctx         sdk.Context
	queryClient govtypes.QueryClient
}

func (suite *TestSuite) SetupTest() {
	checkTx := false
	app := testutil.Setup(checkTx, 5)
	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	govtypes.RegisterQueryServer(queryHelper, app.GovKeeper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = govtypes.NewQueryClient(queryHelper)
}

func (suite *TestSuite) TestGRPCQueryProposal() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	var (
		req         *govtypes.QueryProposalRequest
		expProposal govtypes.Proposal
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"non existing proposal request",
			func() {
				queryProposalMsg := types.QueryProposalMsg{
					ProposalID: "3",
				}
				msg, _ := mgov.MakeQueryProposalMsg(queryProposalMsg)
				req = &msg
			},
			false,
		},
		{
			"zero proposal id request",
			func() {
				queryProposalMsg := types.QueryProposalMsg{
					ProposalID: "0",
				}
				msg, _ := mgov.MakeQueryProposalMsg(queryProposalMsg)
				req = &msg
			},
			false,
		},
		{
			"valid request",
			func() {
				queryProposalMsg := types.QueryProposalMsg{
					ProposalID: "1",
				}
				msg, _ := mgov.MakeQueryProposalMsg(queryProposalMsg)
				req = &msg

				testProposal := govtypes.NewTextProposal("Proposal", "testing proposal")
				submittedProposal, err := app.GovKeeper.SubmitProposal(ctx, testProposal)
				suite.Require().NoError(err)
				suite.Require().NotEmpty(submittedProposal)

				expProposal = submittedProposal
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			proposalRes, err := queryClient.Proposal(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expProposal.String(), proposalRes.Proposal.String())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(proposalRes)
			}
		})
	}
}

func (suite *TestSuite) TestGRPCQueryProposals() {
	app, ctx, queryClient := suite.app, suite.ctx, suite.queryClient

	m1, _ := key.NewMnemonic()
	k1, _ := key.NewPrivKey(m1)
	addr, _ := util.GetAddrByPrivKey(k1)

	testProposals := []govtypes.Proposal{}

	var (
		req    *govtypes.QueryProposalsRequest
		expRes *govtypes.QueryProposalsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"empty state request",
			func() {
				QueryProposalsMsg := types.QueryProposalsMsg{}
				msg, _ := mgov.MakeQueryProposalsMsg(QueryProposalsMsg)
				req = &msg
			},
			true,
		},
		{
			"request proposals with limit 3",
			func() {
				// create 5 test proposals
				for i := 0; i < 5; i++ {
					num := strconv.Itoa(i + 1)
					testProposal := govtypes.NewTextProposal("Proposal"+num, "testing proposal "+num)
					proposal, err := app.GovKeeper.SubmitProposal(ctx, testProposal)
					suite.Require().NotEmpty(proposal)
					suite.Require().NoError(err)
					testProposals = append(testProposals, proposal)
				}

				pagination := types.Pagination{
					Limit: 3,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				QueryProposalsMsg := types.QueryProposalsMsg{}
				msg, _ := mgov.MakeQueryProposalsMsg(QueryProposalsMsg)
				req = &msg

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals[:3],
				}
			},
			true,
		},
		{
			"request 2nd page with limit 4",
			func() {
				pagination := types.Pagination{
					Limit:  3,
					Offset: 3,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				QueryProposalsMsg := types.QueryProposalsMsg{}
				msg, _ := mgov.MakeQueryProposalsMsg(QueryProposalsMsg)
				req = &msg

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals[3:],
				}
			},
			true,
		},
		{
			"request with limit 2 and count true",
			func() {
				pagination := types.Pagination{
					Limit:      2,
					CountTotal: true,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				QueryProposalsMsg := types.QueryProposalsMsg{}
				msg, _ := mgov.MakeQueryProposalsMsg(QueryProposalsMsg)
				req = &msg

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals[:2],
				}
			},
			true,
		},
		{
			"request with filter of status deposit period",
			func() {
				QueryProposalsMsg := types.QueryProposalsMsg{
					Status: "1",
				}
				msg, _ := mgov.MakeQueryProposalsMsg(QueryProposalsMsg)
				req = &msg

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals,
				}
			},
			true,
		},
		{
			"request with filter of deposit address",
			func() {
				depositCoins := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, app.StakingKeeper.TokensFromConsensusPower(ctx, 20)))
				deposit := govtypes.NewDeposit(testProposals[0].ProposalId, addr, depositCoins)
				app.GovKeeper.SetDeposit(ctx, deposit)

				QueryProposalsMsg := types.QueryProposalsMsg{
					Depositor: addr.String(),
				}
				msg, _ := mgov.MakeQueryProposalsMsg(QueryProposalsMsg)
				req = &msg

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals[:1],
				}
			},
			true,
		},
		{
			"request with filter of deposit address",
			func() {
				testProposals[1].Status = govtypes.StatusVotingPeriod
				app.GovKeeper.SetProposal(ctx, testProposals[1])
				suite.Require().NoError(app.GovKeeper.AddVote(ctx, testProposals[1].ProposalId, addr, govtypes.NewNonSplitVoteOption(govtypes.OptionAbstain)))

				QueryProposalsMsg := types.QueryProposalsMsg{
					Voter: addr.String(),
				}
				msg, _ := mgov.MakeQueryProposalsMsg(QueryProposalsMsg)
				req = &msg

				expRes = &govtypes.QueryProposalsResponse{
					Proposals: testProposals[1:2],
				}
			},
			true,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			proposals, err := queryClient.Proposals(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)

				suite.Require().Len(proposals.GetProposals(), len(expRes.GetProposals()))
				for i := 0; i < len(proposals.GetProposals()); i++ {
					suite.Require().Equal(proposals.GetProposals()[i].String(), expRes.GetProposals()[i].String())
				}

			} else {
				suite.Require().Error(err)
				suite.Require().Nil(proposals)
			}
		})
	}
}

func (suite *TestSuite) TestGRPCQueryParams() {
	queryClient := suite.queryClient

	var (
		req    *govtypes.QueryParamsRequest
		expRes *govtypes.QueryParamsResponse
	)

	testCases := []struct {
		msg      string
		malleate func()
		expPass  bool
	}{
		{
			"deposit params request",
			func() {
				govParamsMsg := types.GovParamsMsg{
					ParamType: "deposit",
				}
				msg, _ := mgov.MakeGovParamsMsg(govParamsMsg)
				req = &msg

				expRes = &govtypes.QueryParamsResponse{
					DepositParams: govtypes.DefaultDepositParams(),
					TallyParams:   govtypes.NewTallyParams(sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0)),
				}
			},
			true,
		},
		{
			"voting params request",
			func() {
				govParamsMsg := types.GovParamsMsg{
					ParamType: "voting",
				}
				msg, _ := mgov.MakeGovParamsMsg(govParamsMsg)
				req = &msg
				expRes = &govtypes.QueryParamsResponse{
					VotingParams: govtypes.DefaultVotingParams(),
					TallyParams:  govtypes.NewTallyParams(sdk.NewDec(0), sdk.NewDec(0), sdk.NewDec(0)),
				}
			},
			true,
		},
		{
			"tally params request",
			func() {
				govParamsMsg := types.GovParamsMsg{
					ParamType: "tallying",
				}
				msg, _ := mgov.MakeGovParamsMsg(govParamsMsg)
				req = &msg
				expRes = &govtypes.QueryParamsResponse{
					TallyParams: govtypes.DefaultTallyParams(),
				}
			},
			true,
		},
		{
			"invalid request",
			func() {
				govParamsMsg := types.GovParamsMsg{
					ParamType: "invalid",
				}
				msg, _ := mgov.MakeGovParamsMsg(govParamsMsg)
				req = &msg
				expRes = &govtypes.QueryParamsResponse{}
			},
			false,
		},
	}

	for _, testCase := range testCases {
		suite.Run(fmt.Sprintf("Case %s", testCase.msg), func() {
			testCase.malleate()

			params, err := queryClient.Params(context.Background(), req)

			if testCase.expPass {
				suite.Require().NoError(err)
				suite.Require().Equal(expRes.GetDepositParams(), params.GetDepositParams())
				suite.Require().Equal(expRes.GetVotingParams(), params.GetVotingParams())
				suite.Require().Equal(expRes.GetTallyParams(), params.GetTallyParams())
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(params)
			}
		})
	}
}
