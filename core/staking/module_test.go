package staking_test

import (
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xpladev/xpla.go/core/staking"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestCoreModule() {
	src := rand.NewSource(1)
	r := rand.New(src)
	accounts := testutil.RandomAccounts(r, 2)
	s.xplac.WithPrivateKey(accounts[0].PrivKey)

	c := staking.NewCoreModule()

	// test get name
	s.Require().Equal(staking.StakingModule, c.Name())

	// test tx
	var testMsg interface{}
	txBuilder := s.xplac.GetEncoding().TxConfig.NewTxBuilder()

	// create validator
	tmpVal := sdk.ValAddress(accounts[0].Address)
	createValidatorMsg := types.CreateValidatorMsg{
		NodeKey:                 `{"priv_key":{"type":"tendermint/PrivKeyEd25519","value":"F20DGZKfFFCqgXe2AxF6855KrzfqVasdunk2LMG/EBV+U3gf7GVokgm+X8JP0WG1dyzZ7UddnmC9LGpUMRRQmQ=="}}`,
		PrivValidatorKey:        `{"address":"3C5042645BAD50A98F0A7D567F862E1A861C23C5","pub_key":{"type":"tendermint/PubKeyEd25519","value":"/0bCEBBwUIrjqYr+pKfzHly+SBMjkA/hcCR9oswxnrk="},"priv_key":{"type":"tendermint/PrivKeyEd25519","value":"iks74YM/Di06VI4JPZ3zOxrKfQ0iwwgXhNa6aIzaduf/RsIQEHBQiuOpiv6kp/MeXL5IEyOQD+FwJH2izDGeuQ=="}}`,
		ValidatorAddress:        tmpVal.String(),
		Moniker:                 "moniker",
		Identity:                "identity",
		Website:                 "website",
		SecurityContact:         "securityContact",
		Details:                 "details",
		Amount:                  "1000000000axpla",
		CommissionRate:          "",
		CommissionMaxRate:       "",
		CommissionMaxChangeRate: "",
		MinSelfDelegation:       "",
	}

	makeCreateValidatorMsg, err := staking.MakeCreateValidatorMsg(createValidatorMsg, s.xplac.GetPrivateKey(), s.xplac.GetOutputDocument())
	s.Require().NoError(err)

	testMsg = makeCreateValidatorMsg
	txBuilder, err = c.NewTxRouter(txBuilder, staking.StakingCreateValidatorMsgType, testMsg)
	s.Require().NoError(err)

	// edit validator
	editValidatorMsg := types.EditValidatorMsg{
		Moniker:           "moniker",
		Identity:          "identity",
		Website:           "website",
		SecurityContact:   "securityContact",
		Details:           "details",
		CommissionRate:    "",
		MinSelfDelegation: "",
	}

	makeEditValidatorMsg, err := staking.MakeEditValidatorMsg(editValidatorMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeEditValidatorMsg
	txBuilder, err = c.NewTxRouter(txBuilder, staking.StakingEditValidatorMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeEditValidatorMsg, txBuilder.GetTx().GetMsgs()[0])

	// delegate
	delegateMsg := types.DelegateMsg{
		Amount:  "1000",
		ValAddr: sdk.ValAddress(accounts[0].Address).String(),
	}

	makeDelegateMsg, err := staking.MakeDelegateMsg(delegateMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeDelegateMsg
	txBuilder, err = c.NewTxRouter(txBuilder, staking.StakingDelegateMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeDelegateMsg, txBuilder.GetTx().GetMsgs()[0])

	// unbonding
	unbondMsg := types.UnbondMsg{
		Amount:  "1000",
		ValAddr: sdk.ValAddress(accounts[0].Address).String(),
	}

	makeUnbondMsg, err := staking.MakeUnbondMsg(unbondMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeUnbondMsg
	txBuilder, err = c.NewTxRouter(txBuilder, staking.StakingUnbondMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeUnbondMsg, txBuilder.GetTx().GetMsgs()[0])

	// redelegation
	redelegateMsg := types.RedelegateMsg{
		Amount:     "1000",
		ValSrcAddr: sdk.ValAddress(accounts[0].Address).String(),
		ValDstAddr: sdk.ValAddress(accounts[1].Address).String(),
	}

	makeRedelegateMsg, err := staking.MakeRedelegateMsg(redelegateMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	testMsg = makeRedelegateMsg
	txBuilder, err = c.NewTxRouter(txBuilder, staking.StakingRedelegateMsgType, testMsg)
	s.Require().NoError(err)
	s.Require().Equal(&makeRedelegateMsg, txBuilder.GetTx().GetMsgs()[0])

	// invalid tx msg type
	_, err = c.NewTxRouter(nil, "invalid message type", nil)
	s.Require().Error(err)

	s.xplac = provider.ResetXplac(s.xplac)
}
