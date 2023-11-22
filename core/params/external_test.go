package params_test

import (
	"math/rand"

	mparams "github.com/xpladev/xpla.go/core/params"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestParamsTx() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)

	s.xplac.WithPrivateKey(accounts[0].PrivKey)
	// change params
	paramChangeMsg := types.ParamChangeMsg{
		Title:       "Staking param change",
		Description: "update max validators",
		Changes: []string{
			`{
				"subspace": "staking",
				"key": "MaxValidators",
				"value": 105
			}`,
		},
		Deposit: "1000",
	}
	s.xplac.ParamChange(paramChangeMsg)

	makeProposalParamChangeMsg, err := mparams.MakeProposalParamChangeMsg(paramChangeMsg, s.xplac.GetFromAddress(), s.xplac.GetEncoding())
	s.Require().NoError(err)

	s.Require().Equal(makeProposalParamChangeMsg, s.xplac.GetMsg())
	s.Require().Equal(mparams.ParamsModule, s.xplac.GetModule())
	s.Require().Equal(mparams.ParamsProposalParamChangeMsgType, s.xplac.GetMsgType())

	paramsParamChangeTxbytes, err := s.xplac.ParamChange(paramChangeMsg).CreateAndSignTx()
	s.Require().NoError(err)

	paramsParamChangeJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(paramsParamChangeTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.ParamsParamChangeTxTemplates, string(paramsParamChangeJsonTxbytes))
}

func (s *IntegrationTestSuite) TestParams() {
	// raw params by subspace
	subspaceMsg := types.SubspaceMsg{
		Subspace: "staking",
		Key:      "MaxValidators",
	}
	s.xplac.QuerySubspace(subspaceMsg)

	makeQueryParamsSubspaceMsg, err := mparams.MakeQueryParamsSubspaceMsg(subspaceMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryParamsSubspaceMsg, s.xplac.GetMsg())
	s.Require().Equal(mparams.ParamsModule, s.xplac.GetModule())
	s.Require().Equal(mparams.ParamsQuerySubpsaceMsgType, s.xplac.GetMsgType())
}
