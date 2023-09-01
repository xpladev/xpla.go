package client_test

import (
	"github.com/xpladev/xpla.go/types"

	mauthz "github.com/xpladev/xpla.go/core/authz"
	mbank "github.com/xpladev/xpla.go/core/bank"
	mcrisis "github.com/xpladev/xpla.go/core/crisis"
	mdist "github.com/xpladev/xpla.go/core/distribution"
	mevm "github.com/xpladev/xpla.go/core/evm"
	mfeegrant "github.com/xpladev/xpla.go/core/feegrant"
	mgov "github.com/xpladev/xpla.go/core/gov"
	mparams "github.com/xpladev/xpla.go/core/params"
	mreward "github.com/xpladev/xpla.go/core/reward"
	mslashing "github.com/xpladev/xpla.go/core/slashing"
	mstaking "github.com/xpladev/xpla.go/core/staking"
	mupgrade "github.com/xpladev/xpla.go/core/upgrade"
	mwasm "github.com/xpladev/xpla.go/core/wasm"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (s *ClientTestSuite) TestAuthzTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// authz grant
	authzGrantMsg := types.AuthzGrantMsg{
		Granter:           s.accounts[0].Address.String(),
		Grantee:           s.accounts[1].Address.String(),
		AuthorizationType: "send",
		SpendLimit:        "1000",
	}
	s.xplac.AuthzGrant(authzGrantMsg)

	makeAuthzGrantMsg, err := mauthz.MakeAuthzGrantMsg(authzGrantMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeAuthzGrantMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzGrantMsgType, s.xplac.GetMsgType())

	// authz revoke
	authzRevokeMsg := types.AuthzRevokeMsg{
		Granter: s.accounts[0].Address.String(),
		Grantee: s.accounts[1].Address.String(),
		MsgType: "/cosmos.bank.v1beta1.MsgSend",
	}
	s.xplac.AuthzRevoke(authzRevokeMsg)

	makeAuthzRevokeMsg, err := mauthz.MakeAuthzRevokeMsg(authzRevokeMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeAuthzRevokeMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzRevokeMsgType, s.xplac.GetMsgType())

	// authz exec
	// e.g. bank send
	bankSendMsg := types.BankSendMsg{
		FromAddress: s.accounts[0].Address.String(),
		ToAddress:   s.accounts[1].Address.String(),
		Amount:      "1000",
	}

	txbytesBankSend, err := s.xplac.BankSend(bankSendMsg).CreateAndSignTx()
	s.Require().NoError(err)

	sdkTx, err := s.xplac.GetEncoding().TxConfig.TxDecoder()(txbytesBankSend)
	s.Require().NoError(err)

	jsonTx, err := s.xplac.GetEncoding().TxConfig.TxJSONEncoder()(sdkTx)
	s.Require().NoError(err)

	authzExecMsg := types.AuthzExecMsg{
		Grantee:      s.accounts[1].Address.String(),
		ExecTxString: string(jsonTx),
	}
	s.xplac.AuthzExec(authzExecMsg)

	makeAuthzExecMsg, err := mauthz.MakeAuthzExecMsg(authzExecMsg, s.xplac.GetEncoding())
	s.Require().NoError(err)

	s.Require().Equal(makeAuthzExecMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzExecMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestBankTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// bank send
	bankSendMsg := types.BankSendMsg{
		FromAddress: s.accounts[0].Address.String(),
		ToAddress:   s.accounts[1].Address.String(),
		Amount:      "1000",
	}
	s.xplac.BankSend(bankSendMsg)

	makeBankSendMsg, err := mbank.MakeBankSendMsg(bankSendMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeBankSendMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankSendMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestCrisisTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// invariant broken
	invariantBrokenMsg := types.InvariantBrokenMsg{
		ModuleName:     "bank",
		InvariantRoute: "total-supply",
	}
	s.xplac.InvariantBroken(invariantBrokenMsg)

	makeInvariantRouteMsg, err := mcrisis.MakeInvariantRouteMsg(invariantBrokenMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeInvariantRouteMsg, s.xplac.GetMsg())
	s.Require().Equal(mcrisis.CrisisModule, s.xplac.GetModule())
	s.Require().Equal(mcrisis.CrisisInvariantBrokenMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestDistributionTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// fund community pool
	fundCommunityPoolMsg := types.FundCommunityPoolMsg{
		Amount: "1000",
	}
	s.xplac.FundCommunityPool(fundCommunityPoolMsg)

	makeFundCommunityPoolMsg, err := mdist.MakeFundCommunityPoolMsg(fundCommunityPoolMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeFundCommunityPoolMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionFundCommunityPoolMsgType, s.xplac.GetMsgType())

	// community pool spend
	communityPoolSpendMsg := types.CommunityPoolSpendMsg{
		Title:       "community pool spend",
		Description: "pay me",
		Recipient:   s.accounts[0].Address.String(),
		Amount:      "1000",
		Deposit:     "1000",
	}
	s.xplac.CommunityPoolSpend(communityPoolSpendMsg)

	makeProposalCommunityPoolSpendMsg, err := mdist.MakeProposalCommunityPoolSpendMsg(communityPoolSpendMsg, s.xplac.GetPrivateKey(), s.xplac.GetEncoding())
	s.Require().NoError(err)

	s.Require().Equal(makeProposalCommunityPoolSpendMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionProposalCommunityPoolSpendMsgType, s.xplac.GetMsgType())

	// withdraw rewards
	withdrawRewardsMsg := types.WithdrawRewardsMsg{
		DelegatorAddr: s.accounts[0].Address.String(),
		ValidatorAddr: s.network.Validators[0].ValAddress.String(),
		Commission:    true,
	}
	s.xplac.WithdrawRewards(withdrawRewardsMsg)

	makeWithdrawRewardsMsg, err := mdist.MakeWithdrawRewardsMsg(withdrawRewardsMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeWithdrawRewardsMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionWithdrawRewardsMsgType, s.xplac.GetMsgType())

	// set withdraw address
	setWithdrawAddrMsg := types.SetWithdrawAddrMsg{
		WithdrawAddr: s.accounts[0].Address.String(),
	}
	s.xplac.SetWithdrawAddr(setWithdrawAddrMsg)

	makeSetWithdrawAddrMsg, err := mdist.MakeSetWithdrawAddrMsg(setWithdrawAddrMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeSetWithdrawAddrMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionSetWithdrawAddrMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestEvmTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// send evm coin
	sendCoinMsg := types.SendCoinMsg{
		FromAddress: s.accounts[0].PubKey.Address().String(),
		ToAddress:   s.accounts[1].PubKey.Address().String(),
		Amount:      "1000",
	}
	s.xplac.EvmSendCoin(sendCoinMsg)

	makeSendCoinMsg, err := mevm.MakeSendCoinMsg(sendCoinMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeSendCoinMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmSendCoinMsgType, s.xplac.GetMsgType())

	// deploy solidity contract
	deploySolContractMsg := types.DeploySolContractMsg{
		ABIJsonFilePath:      testABIPath,
		BytecodeJsonFilePath: testBytecodePath,
		Args:                 nil,
	}
	s.xplac.DeploySolidityContract(deploySolContractMsg)

	makeDeploySolContractMsg, err := mevm.MakeDeploySolContractMsg(deploySolContractMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeDeploySolContractMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmDeploySolContractMsgType, s.xplac.GetMsgType())

	// invoke solidity contract
	invokeSolContractMsg := types.InvokeSolContractMsg{
		ContractAddress:      testSolContractAddress,
		ABIJsonFilePath:      testABIPath,
		BytecodeJsonFilePath: testBytecodePath,
		Args:                 nil,
	}
	s.xplac.InvokeSolidityContract(invokeSolContractMsg)

	makeInvokeSolContractMsg, err := mevm.MakeInvokeSolContractMsg(invokeSolContractMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeInvokeSolContractMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmInvokeSolContractMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestFeegrantTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// feegrant
	feeGrantMsg := types.FeeGrantMsg{
		Granter:    s.accounts[0].Address.String(),
		Grantee:    s.accounts[1].Address.String(),
		SpendLimit: "1000",
		// Period:      "3600",
		// PeriodLimit: "10",
		Expiration: "2100-01-01T23:59:59+00:00",
	}
	s.xplac.FeeGrant(feeGrantMsg)

	makeFeeGrantMsg, err := mfeegrant.MakeFeeGrantMsg(feeGrantMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeFeeGrantMsg, s.xplac.GetMsg())
	s.Require().Equal(mfeegrant.FeegrantModule, s.xplac.GetModule())
	s.Require().Equal(mfeegrant.FeegrantGrantMsgType, s.xplac.GetMsgType())

	// revoke feegrant
	revokeFeeGrantMsg := types.RevokeFeeGrantMsg{
		Granter: s.accounts[0].Address.String(),
		Grantee: s.accounts[1].Address.String(),
	}
	s.xplac.RevokeFeeGrant(revokeFeeGrantMsg)

	makeRevokeFeeGrantMsg, err := mfeegrant.MakeRevokeFeeGrantMsg(revokeFeeGrantMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeRevokeFeeGrantMsg, s.xplac.GetMsg())
	s.Require().Equal(mfeegrant.FeegrantModule, s.xplac.GetModule())
	s.Require().Equal(mfeegrant.FeegrantRevokeGrantMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestGovTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// submit proposal
	submitProposalMsg := types.SubmitProposalMsg{
		Title:       "Test proposal",
		Description: "Proposal description",
		Type:        "text",
		Deposit:     "1000",
	}
	s.xplac.SubmitProposal(submitProposalMsg)

	makeSubmitProposalMsg, err := mgov.MakeSubmitProposalMsg(submitProposalMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeSubmitProposalMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovSubmitProposalMsgType, s.xplac.GetMsgType())

	// deposit
	govDepositMsg := types.GovDepositMsg{
		ProposalID: "1",
		Deposit:    "1000",
	}
	s.xplac.GovDeposit(govDepositMsg)

	makeGovDepositMsg, err := mgov.MakeGovDepositMsg(govDepositMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeGovDepositMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovDepositMsgType, s.xplac.GetMsgType())

	// vote
	voteMsg := types.VoteMsg{
		ProposalID: "1",
		Option:     "yes",
	}
	s.xplac.Vote(voteMsg)

	makeVoteMsg, err := mgov.MakeVoteMsg(voteMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeVoteMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovVoteMsgType, s.xplac.GetMsgType())

	// weighted vote
	weightedVoteMsg := types.WeightedVoteMsg{
		ProposalID: "1",
		Yes:        "0.6",
		No:         "0.3",
		Abstain:    "0.05",
		NoWithVeto: "0.05",
	}
	s.xplac.WeightedVote(weightedVoteMsg)

	makeWeightedVoteMsg, err := mgov.MakeWeightedVoteMsg(weightedVoteMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeWeightedVoteMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovWeightedVoteMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestParamsTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
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

	makeProposalParamChangeMsg, err := mparams.MakeProposalParamChangeMsg(paramChangeMsg, s.xplac.GetPrivateKey(), s.xplac.GetEncoding())
	s.Require().NoError(err)

	s.Require().Equal(makeProposalParamChangeMsg, s.xplac.GetMsg())
	s.Require().Equal(mparams.ParamsModule, s.xplac.GetModule())
	s.Require().Equal(mparams.ParamsProposalParamChangeMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestRewardTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// fund fee collector
	fundFeeCollectorMsg := types.FundFeeCollectorMsg{
		DepositorAddr: s.accounts[0].Address.String(),
		Amount:        "1000",
	}
	s.xplac.FundFeeCollector(fundFeeCollectorMsg)

	makeFundFeeCollectorMsg, err := mreward.MakeFundFeeCollectorMsg(fundFeeCollectorMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeFundFeeCollectorMsg, s.xplac.GetMsg())
	s.Require().Equal(mreward.RewardModule, s.xplac.GetModule())
	s.Require().Equal(mreward.RewardFundFeeCollectorMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestSlashingTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// unjail
	s.xplac.Unjail()

	makeUnjailMsg, err := mslashing.MakeUnjailMsg(s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeUnjailMsg, s.xplac.GetMsg())
	s.Require().Equal(mslashing.SlashingModule, s.xplac.GetModule())
	s.Require().Equal(mslashing.SlahsingUnjailMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestStakingTx() {
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

	makeCreateValidatorMsg, err := mstaking.MakeCreateValidatorMsg(createValidatorMsg, s.xplac.GetPrivateKey(), s.xplac.GetOutputDocument())
	s.Require().NoError(err)

	s.Require().Equal(makeCreateValidatorMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingCreateValidatorMsgType, s.xplac.GetMsgType())

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

	makeEditValidatorMsg, err := mstaking.MakeEditValidatorMsg(editValidatorMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeEditValidatorMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingEditValidatorMsgType, s.xplac.GetMsgType())

	// delegate
	delegateMsg := types.DelegateMsg{
		Amount:  "1000",
		ValAddr: s.network.Validators[0].ValAddress.String(),
	}
	s.xplac.Delegate(delegateMsg)

	makeDelegateMsg, err := mstaking.MakeDelegateMsg(delegateMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeDelegateMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingDelegateMsgType, s.xplac.GetMsgType())

	// unbonding
	unbondMsg := types.UnbondMsg{
		Amount:  "1000",
		ValAddr: s.network.Validators[0].ValAddress.String(),
	}
	s.xplac.Unbond(unbondMsg)

	makeUnbondMsg, err := mstaking.MakeUnbondMsg(unbondMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeUnbondMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingUnbondMsgType, s.xplac.GetMsgType())

	// redelegation
	redelegateMsg := types.RedelegateMsg{
		Amount:     "1000",
		ValSrcAddr: s.network.Validators[0].ValAddress.String(),
		ValDstAddr: s.network.Validators[1].ValAddress.String(),
	}
	s.xplac.Redelegate(redelegateMsg)

	makeRedelegateMsg, err := mstaking.MakeRedelegateMsg(redelegateMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeRedelegateMsg, s.xplac.GetMsg())
	s.Require().Equal(mstaking.StakingModule, s.xplac.GetModule())
	s.Require().Equal(mstaking.StakingRedelegateMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestUpgradeTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// software upgrade
	softwareUpgradeMsg := types.SoftwareUpgradeMsg{
		UpgradeName:   "Upgrade Name",
		Title:         "Upgrade Title",
		Description:   "Upgrade Description",
		UpgradeHeight: "6000",
		UpgradeInfo:   `{"upgrade_info":"INFO"}`,
		Deposit:       "1000",
	}
	s.xplac.SoftwareUpgrade(softwareUpgradeMsg)

	makeProposalSoftwareUpgradeMsg, err := mupgrade.MakeProposalSoftwareUpgradeMsg(softwareUpgradeMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeProposalSoftwareUpgradeMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradeProposalSoftwareUpgradeMsgType, s.xplac.GetMsgType())

	// cancel software upgrade
	cancelSoftwareUpgradeMsg := types.CancelSoftwareUpgradeMsg{
		Title:       "Cancel software upgrade",
		Description: "Cancel software upgrade description",
		Deposit:     "1000",
	}
	s.xplac.CancelSoftwareUpgrade(cancelSoftwareUpgradeMsg)

	makeCancelSoftwareUpgradeMsg, err := mupgrade.MakeCancelSoftwareUpgradeMsg(cancelSoftwareUpgradeMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeCancelSoftwareUpgradeMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradeCancelSoftwareUpgradeMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestWasmTx() {
	s.xplac.WithPrivateKey(s.accounts[0].PrivKey)
	// store code
	storeMsg := types.StoreMsg{
		FilePath:              testWasmFilePath,
		InstantiatePermission: "instantiate-only-sender",
	}
	s.xplac.StoreCode(storeMsg)

	makeStoreCodeMsg, err := mwasm.MakeStoreCodeMsg(storeMsg, s.accounts[0].Address)
	s.Require().NoError(err)

	s.Require().Equal(makeStoreCodeMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmStoreMsgType, s.xplac.GetMsgType())

	// instantiate
	instantiateMsg := types.InstantiateMsg{
		CodeId:  "1",
		Amount:  "10",
		Label:   "Contract instant",
		InitMsg: `{"owner":"` + s.accounts[0].Address.String() + `"}`,
		Admin:   s.accounts[0].Address.String(),
	}
	s.xplac.InstantiateContract(instantiateMsg)

	makeInstantiateMsg, err := mwasm.MakeInstantiateMsg(instantiateMsg, s.accounts[0].Address)
	s.Require().NoError(err)

	s.Require().Equal(makeInstantiateMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmInstantiateMsgType, s.xplac.GetMsgType())

	// execute
	executeMsg := types.ExecuteMsg{
		ContractAddress: testCWContractAddress,
		Amount:          "0",
		ExecMsg:         `{"execute_method":{"execute_key":"execute_test","execute_value":"execute_val"}}`,
	}
	s.xplac.ExecuteContract(executeMsg)

	makeExecuteMsg, err := mwasm.MakeExecuteMsg(executeMsg, s.accounts[0].Address)
	s.Require().NoError(err)

	s.Require().Equal(makeExecuteMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmExecuteMsgType, s.xplac.GetMsgType())

	// clear contract admin
	clearContractAdminMsg := types.ClearContractAdminMsg{
		ContractAddress: testCWContractAddress,
	}
	s.xplac.ClearContractAdmin(clearContractAdminMsg)

	makeClearContractAdminMsg, err := mwasm.MakeClearContractAdminMsg(clearContractAdminMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeClearContractAdminMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmClearContractAdminMsgType, s.xplac.GetMsgType())

	// set contract admin
	setContractAdminMsg := types.SetContractAdminMsg{
		ContractAddress: testCWContractAddress,
	}
	s.xplac.SetContractAdmin(setContractAdminMsg)

	makeSetContractAdmintMsg, err := mwasm.MakeSetContractAdmintMsg(setContractAdminMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeSetContractAdmintMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmSetContractAdminMsgType, s.xplac.GetMsgType())

	// migrate
	migrateMsg := types.MigrateMsg{
		ContractAddress: testCWContractAddress,
		CodeId:          "2",
		MigrateMsg:      `{}`,
	}
	s.xplac.Migrate(migrateMsg)

	makeMigrateMsg, err := mwasm.MakeMigrateMsg(migrateMsg, s.xplac.GetPrivateKey())
	s.Require().NoError(err)

	s.Require().Equal(makeMigrateMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmMigrateMsgType, s.xplac.GetMsgType())
}
