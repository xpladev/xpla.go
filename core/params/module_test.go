package params_test

import (
	"math/rand"

	"github.com/xpladev/xpla.go/client"
	"github.com/xpladev/xpla.go/core/params"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := params.NewCoreModule()

	// test get name
	s.Require().Equal(params.ParamsModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

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

	makeProposalParamChangeMsg, err := params.MakeProposalParamChangeMsg(paramChangeMsg, s.xplac.GetPrivateKey(), s.xplac.GetEncoding())
	s.Require().NoError(err)

	testMsg = makeProposalParamChangeMsg
	txBuilder, err = c.NewTxRouter(txBuilder, params.ParamsProposalParamChangeMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeProposalParamChangeMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = client.ResetXplac(s.xplac)
}