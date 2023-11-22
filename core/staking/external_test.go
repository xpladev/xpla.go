package staking_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	mstaking "github.com/xpladev/xpla.go/core/staking"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util/testutil"
)

func (s *IntegrationTestSuite) TestStakingTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	tmpVal := sdk.ValAddress(s.accounts[0].Address)

	// create validator
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
	s.xplac.CreateValidator(createValidatorMsg)

	makeCreateValidatorMsg, err := mstaking.MakeCreateValidatorMsg(createValidatorMsg, s.xplac.GetFromAddress(), s.xplac.GetOutputDocument())
	s.Require().NoError(err)

	s.Require().Equal(makeCreateValidatorMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingCreateValidatorMsgType, s.xplac.GetMsgType())

	_, err = s.xplac.CreateValidator(createValidatorMsg).CreateAndSignTx()
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
	s.xplac.EditValidator(editValidatorMsg)

	makeEditValidatorMsg, err := mstaking.MakeEditValidatorMsg(editValidatorMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeEditValidatorMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingEditValidatorMsgType, s.xplac.GetMsgType())

	stakingEditValidatorTxbytes, err := s.xplac.EditValidator(editValidatorMsg).CreateAndSignTx()
	s.Require().NoError(err)

	stakingEditValidatorJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(stakingEditValidatorTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.StakingEditValidatorTxTemplates, string(stakingEditValidatorJsonTxbytes))

	// delegate
	delegateMsg := types.DelegateMsg{
		Amount:  "1000",
		ValAddr: sdk.ValAddress(s.accounts[0].Address).String(),
	}
	s.xplac.Delegate(delegateMsg)

	makeDelegateMsg, err := mstaking.MakeDelegateMsg(delegateMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeDelegateMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingDelegateMsgType, s.xplac.GetMsgType())

	stakingDelegateTxbytes, err := s.xplac.Delegate(delegateMsg).CreateAndSignTx()
	s.Require().NoError(err)

	stakingDelegateJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(stakingDelegateTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.StakingDelegateTxTemplates, string(stakingDelegateJsonTxbytes))

	// unbonding
	unbondMsg := types.UnbondMsg{
		Amount:  "1000",
		ValAddr: sdk.ValAddress(s.accounts[0].Address).String(),
	}
	s.xplac.Unbond(unbondMsg)

	makeUnbondMsg, err := mstaking.MakeUnbondMsg(unbondMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeUnbondMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingUnbondMsgType, s.xplac.GetMsgType())

	stakingUnbondTxbytes, err := s.xplac.Unbond(unbondMsg).CreateAndSignTx()
	s.Require().NoError(err)

	stakingUnbondJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(stakingUnbondTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.StakingUnbondTxTemplates, string(stakingUnbondJsonTxbytes))

	// redelegation
	redelegateMsg := types.RedelegateMsg{
		Amount:     "1000",
		ValSrcAddr: sdk.ValAddress(s.accounts[0].Address).String(),
		ValDstAddr: sdk.ValAddress(s.accounts[1].Address).String(),
	}
	s.xplac.Redelegate(redelegateMsg)

	makeRedelegateMsg, err := mstaking.MakeRedelegateMsg(redelegateMsg, s.xplac.GetFromAddress())
	s.Require().NoError(err)

	s.Require().Equal(makeRedelegateMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingRedelegateMsgType, s.xplac.GetMsgType())

	stakingRedelegateTxbytes, err := s.xplac.Redelegate(redelegateMsg).CreateAndSignTx()
	s.Require().NoError(err)

	stakingRedelegateJsonTxbytes, err := s.xplac.EncodedTxbytesToJsonTx(stakingRedelegateTxbytes)
	s.Require().NoError(err)
	s.Require().Equal(testutil.StakingRedelegateTxTemplates, string(stakingRedelegateJsonTxbytes))
}

func (s *IntegrationTestSuite) TestStaking() {
	val := s.network.Validators[0].ValAddress.String()
	val2 := s.network.Validators[1].ValAddress.String()
	// query validators
	s.xplac.QueryValidators()

	makeQueryValidatorsMsg, err := mstaking.MakeQueryValidatorsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryValidatorsMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryValidatorsMsgType, s.xplac.GetMsgType())

	// query validator
	queryValidatorMsg := types.QueryValidatorMsg{
		ValidatorAddr: val,
	}

	s.xplac.QueryValidators(queryValidatorMsg)

	makeQueryValidatorMsg, err := mstaking.MakeQueryValidatorMsg(queryValidatorMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryValidatorMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryValidatorMsgType, s.xplac.GetMsgType())

	// delegation
	queryDelegationMsg := types.QueryDelegationMsg{
		DelegatorAddr: s.accounts[0].Address.String(),
		ValidatorAddr: val,
	}

	s.xplac.QueryDelegation(queryDelegationMsg)

	makeQueryDelegationMsg, err := mstaking.MakeQueryDelegationMsg(queryDelegationMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDelegationMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryDelegationMsgType, s.xplac.GetMsgType())

	// delegations
	queryDelegationMsg = types.QueryDelegationMsg{
		DelegatorAddr: s.accounts[0].Address.String(),
	}

	s.xplac.QueryDelegation(queryDelegationMsg)

	makeQueryDelegationsMsg, err := mstaking.MakeQueryDelegationsMsg(queryDelegationMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDelegationsMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryDelegationsMsgType, s.xplac.GetMsgType())

	// delegations to
	queryDelegationMsg = types.QueryDelegationMsg{
		ValidatorAddr: val,
	}

	s.xplac.QueryDelegation(queryDelegationMsg)

	makeQueryDelegationsToMsg, err := mstaking.MakeQueryDelegationsToMsg(queryDelegationMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDelegationsToMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryDelegationsToMsgType, s.xplac.GetMsgType())

	// unbonding delegation
	queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg{
		DelegatorAddr: s.accounts[0].Address.String(),
		ValidatorAddr: val,
	}

	s.xplac.QueryUnbondingDelegation(queryUnbondingDelegationMsg)

	makeQueryUnbondingDelegationMsg, err := mstaking.MakeQueryUnbondingDelegationMsg(queryUnbondingDelegationMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryUnbondingDelegationMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryUnbondingDelegationMsgType, s.xplac.GetMsgType())

	// unbonding delegations
	queryUnbondingDelegationMsg = types.QueryUnbondingDelegationMsg{
		DelegatorAddr: s.accounts[0].Address.String(),
	}

	s.xplac.QueryUnbondingDelegation(queryUnbondingDelegationMsg)

	makeQueryUnbondingDelegationsMsg, err := mstaking.MakeQueryUnbondingDelegationsMsg(queryUnbondingDelegationMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryUnbondingDelegationsMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryUnbondingDelegationsMsgType, s.xplac.GetMsgType())

	// unbonding delegations from
	queryUnbondingDelegationMsg = types.QueryUnbondingDelegationMsg{
		ValidatorAddr: val,
	}

	s.xplac.QueryUnbondingDelegation(queryUnbondingDelegationMsg)

	makeQueryUnbondingDelegationsFromMsg, err := mstaking.MakeQueryUnbondingDelegationsFromMsg(queryUnbondingDelegationMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryUnbondingDelegationsFromMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryUnbondingDelegationsFromMsgType, s.xplac.GetMsgType())

	// redelegation
	queryRedelegationMsg := types.QueryRedelegationMsg{
		DelegatorAddr:    s.accounts[0].Address.String(),
		SrcValidatorAddr: val,
		DstValidatorAddr: val2,
	}

	s.xplac.QueryRedelegation(queryRedelegationMsg)

	makeQueryRedelegationMsg, err := mstaking.MakeQueryRedelegationMsg(queryRedelegationMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryRedelegationMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryRedelegationMsgType, s.xplac.GetMsgType())

	// redelegations
	queryRedelegationMsg = types.QueryRedelegationMsg{
		DelegatorAddr: s.accounts[0].Address.String(),
	}

	s.xplac.QueryRedelegation(queryRedelegationMsg)

	makeQueryRedelegationsMsg, err := mstaking.MakeQueryRedelegationsMsg(queryRedelegationMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryRedelegationsMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryRedelegationsMsgType, s.xplac.GetMsgType())

	// redelegations from
	queryRedelegationMsg = types.QueryRedelegationMsg{
		SrcValidatorAddr: val,
	}

	s.xplac.QueryRedelegation(queryRedelegationMsg)

	makeQueryRedelegationsFromMsg, err := mstaking.MakeQueryRedelegationsFromMsg(queryRedelegationMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryRedelegationsFromMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryRedelegationsFromMsgType, s.xplac.GetMsgType())

	// hiestorcal info
	historicalInfoMsg := types.HistoricalInfoMsg{
		Height: "1",
	}

	s.xplac.HistoricalInfo(historicalInfoMsg)

	makeHistoricalInfoMsg, err := mstaking.MakeHistoricalInfoMsg(historicalInfoMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeHistoricalInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingHistoricalInfoMsgType, s.xplac.GetMsgType())

	// staking pool
	s.xplac.StakingPool()

	makeQueryStakingPoolMsg, err := mstaking.MakeQueryStakingPoolMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryStakingPoolMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryStakingPoolMsgType, s.xplac.GetMsgType())

	// staking params
	s.xplac.StakingParams()

	makeQueryStakingParamsMsg, err := mstaking.MakeQueryStakingParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryStakingParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingQueryStakingParamsMsgType, s.xplac.GetMsgType())
}
