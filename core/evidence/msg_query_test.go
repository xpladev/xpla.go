package evidence_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/xpladev/xpla.go/core"
	mevidence "github.com/xpladev/xpla.go/core/evidence"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/evidence/exported"
	evidencetypes "github.com/cosmos/cosmos-sdk/x/evidence/types"
	"github.com/stretchr/testify/suite"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	evidencekeeper "github.com/cosmos/cosmos-sdk/x/evidence/keeper"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type TestSuite struct {
	suite.Suite

	ctx         sdk.Context
	app         *xapp.XplaApp
	queryClient evidencetypes.QueryClient
}

func (suite *TestSuite) SetupTest() {
	checkTx := false
	app := testutil.Setup(checkTx, 5)
	ctx := app.BaseApp.NewContext(checkTx, tmproto.Header{})

	evidenceKeeper := evidencekeeper.NewKeeper(
		app.AppCodec(), app.GetKey(evidencetypes.StoreKey), app.StakingKeeper, app.SlashingKeeper,
	)
	router := evidencetypes.NewRouter()
	router = router.AddRoute(evidencetypes.RouteEquivocation, testEquivocationHandler(*evidenceKeeper))
	evidenceKeeper.SetRouter(router)

	app.EvidenceKeeper = *evidenceKeeper

	queryHelper := baseapp.NewQueryServerTestHelper(ctx, app.InterfaceRegistry())
	evidencetypes.RegisterQueryServer(queryHelper, app.EvidenceKeeper)

	suite.app = app
	suite.ctx = ctx
	suite.queryClient = evidencetypes.NewQueryClient(queryHelper)
}

func (suite *TestSuite) TestQueryEvidence() {
	var (
		req      *evidencetypes.QueryEvidenceRequest
		evidence []exported.Evidence
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		posttests func(res *evidencetypes.QueryEvidenceResponse)
	}{
		{
			"invalid request with empty evidence hash",
			func() {

				queryEvidenceMsg := types.QueryEvidenceMsg{
					Hash: tmbytes.HexBytes{}.String(),
				}
				msg, _ := mevidence.MakeQueryEvidenceMsg(queryEvidenceMsg)
				req = &msg
			},
			false,
			func(res *evidencetypes.QueryEvidenceResponse) {},
		},
		{
			"success",
			func() {
				numEvidence := 100
				evidence = suite.populateEvidence(suite.ctx, numEvidence)
				queryEvidenceMsg := types.QueryEvidenceMsg{
					Hash: evidence[0].Hash().String(),
				}
				msg, _ := mevidence.MakeQueryEvidenceMsg(queryEvidenceMsg)
				req = &msg
			},
			true,
			func(res *evidencetypes.QueryEvidenceResponse) {
				var evi exported.Evidence
				err := suite.app.InterfaceRegistry().UnpackAny(res.Evidence, &evi)
				suite.Require().NoError(err)
				suite.Require().NotNil(evi)
				suite.Require().Equal(evi, evidence[0])
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)

			res, err := suite.queryClient.Evidence(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}

			tc.posttests(res)
		})
	}
}

func (suite *TestSuite) TestQueryAllEvidence() {
	var (
		req *evidencetypes.QueryAllEvidenceRequest
	)

	testCases := []struct {
		msg       string
		malleate  func()
		expPass   bool
		posttests func(res *evidencetypes.QueryAllEvidenceResponse)
	}{
		{
			"success without evidence",
			func() {
				msg, _ := mevidence.MakeQueryAllEvidenceMsg()
				req = &msg
			},
			true,
			func(res *evidencetypes.QueryAllEvidenceResponse) {
				suite.Require().Empty(res.Evidence)
			},
		},
		{
			"success",
			func() {
				numEvidence := 100
				_ = suite.populateEvidence(suite.ctx, numEvidence)

				pagination := types.Pagination{
					PageKey:    "",
					CountTotal: false,
					Limit:      50,
				}
				pr, _ := core.ReadPageRequest(pagination)
				core.PageRequest = pr

				msg, _ := mevidence.MakeQueryAllEvidenceMsg()
				req = &msg
			},
			true,
			func(res *evidencetypes.QueryAllEvidenceResponse) {
				suite.Equal(len(res.Evidence), 50)
				suite.NotNil(res.Pagination.NextKey)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Case %s", tc.msg), func() {
			suite.SetupTest()

			tc.malleate()
			ctx := sdk.WrapSDKContext(suite.ctx)

			res, err := suite.queryClient.AllEvidence(ctx, req)

			if tc.expPass {
				suite.Require().NoError(err)
				suite.Require().NotNil(res)
			} else {
				suite.Require().Error(err)
				suite.Require().Nil(res)
			}

			tc.posttests(res)
		})
	}
}

func (suite *TestSuite) populateEvidence(ctx sdk.Context, numEvidence int) []exported.Evidence {
	evidence := make([]exported.Evidence, numEvidence)

	for i := 0; i < numEvidence; i++ {
		pk := ed25519.GenPrivKey()

		evidence[i] = &evidencetypes.Equivocation{
			Height:           11,
			Power:            100,
			Time:             time.Now().UTC(),
			ConsensusAddress: sdk.ConsAddress(pk.PubKey().Address().Bytes()).String(),
		}

		suite.Nil(suite.app.EvidenceKeeper.SubmitEvidence(ctx, evidence[i]))
	}

	return evidence
}

func testEquivocationHandler(_ interface{}) evidencetypes.Handler {
	return func(ctx sdk.Context, e exported.Evidence) error {
		if err := e.ValidateBasic(); err != nil {
			return err
		}

		ee, ok := e.(*evidencetypes.Equivocation)
		if !ok {
			return fmt.Errorf("unexpected evidence type: %T", e)
		}
		if ee.Height%2 == 0 {
			return fmt.Errorf("unexpected even evidence height: %d", ee.Height)
		}

		return nil
	}
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
