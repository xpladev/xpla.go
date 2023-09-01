package client

import (
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
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
)

// Authz module

// Grant authorization to an address.
func (xplac *XplaClient) AuthzGrant(authzGrantMsg types.AuthzGrantMsg) *XplaClient {
	msg, err := mauthz.MakeAuthzGrantMsg(authzGrantMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mauthz.AuthzModule).
		WithMsgType(mauthz.AuthzGrantMsgType).
		WithMsg(msg)
	return xplac
}

// Revoke authorization.
func (xplac *XplaClient) AuthzRevoke(authzRevokeMsg types.AuthzRevokeMsg) *XplaClient {
	msg, err := mauthz.MakeAuthzRevokeMsg(authzRevokeMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mauthz.AuthzModule).
		WithMsgType(mauthz.AuthzRevokeMsgType).
		WithMsg(msg)
	return xplac
}

// Execute transaction on behalf of granter account.
func (xplac *XplaClient) AuthzExec(authzExecMsg types.AuthzExecMsg) *XplaClient {
	msg, err := mauthz.MakeAuthzExecMsg(authzExecMsg, xplac.GetEncoding())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mauthz.AuthzModule).
		WithMsgType(mauthz.AuthzExecMsgType).
		WithMsg(msg)
	return xplac
}

// Bank module

// Send funds from one account to another.
func (xplac *XplaClient) BankSend(bankSendMsg types.BankSendMsg) *XplaClient {
	msg, err := mbank.MakeBankSendMsg(bankSendMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mbank.BankModule).
		WithMsgType(mbank.BankSendMsgType).
		WithMsg(msg)
	return xplac
}

// Crisis module

// Submit proof that an invariant broken to halt the chain.
func (xplac *XplaClient) InvariantBroken(invariantBrokenMsg types.InvariantBrokenMsg) *XplaClient {
	msg, err := mcrisis.MakeInvariantRouteMsg(invariantBrokenMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mcrisis.CrisisModule).
		WithMsgType(mcrisis.CrisisInvariantBrokenMsgType).
		WithMsg(msg)
	return xplac
}

// Distribution module

// Funds the community pool with the specified amount.
func (xplac *XplaClient) FundCommunityPool(fundCommunityPoolMsg types.FundCommunityPoolMsg) *XplaClient {
	msg, err := mdist.MakeFundCommunityPoolMsg(fundCommunityPoolMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mdist.DistributionModule).
		WithMsgType(mdist.DistributionFundCommunityPoolMsgType).
		WithMsg(msg)
	return xplac
}

// Submit a community pool spend proposal.
func (xplac *XplaClient) CommunityPoolSpend(communityPoolSpendMsg types.CommunityPoolSpendMsg) *XplaClient {
	msg, err := mdist.MakeProposalCommunityPoolSpendMsg(communityPoolSpendMsg, xplac.GetPrivateKey(), xplac.GetEncoding())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mdist.DistributionModule).
		WithMsgType(mdist.DistributionProposalCommunityPoolSpendMsgType).
		WithMsg(msg)
	return xplac
}

// Withdraw rewards from a given delegation address, and optionally withdraw validator commission if the delegation address given is a validator operator.
func (xplac *XplaClient) WithdrawRewards(withdrawRewardsMsg types.WithdrawRewardsMsg) *XplaClient {
	msg, err := mdist.MakeWithdrawRewardsMsg(withdrawRewardsMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mdist.DistributionModule).
		WithMsgType(mdist.DistributionWithdrawRewardsMsgType).
		WithMsg(msg)
	return xplac
}

// Withdraw all delegations rewards for a delegator.
func (xplac *XplaClient) WithdrawAllRewards() *XplaClient {
	msg, err := mdist.MakeWithdrawAllRewardsMsg(xplac.GetPrivateKey(), xplac.GetGrpcClient(), xplac.GetContext())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mdist.DistributionModule).
		WithMsgType(mdist.DistributionWithdrawAllRewardsMsgType).
		WithMsg(msg)
	return xplac
}

// Change the default withdraw address for rewards associated with an address.
func (xplac *XplaClient) SetWithdrawAddr(setWithdrawAddrMsg types.SetWithdrawAddrMsg) *XplaClient {
	msg, err := mdist.MakeSetWithdrawAddrMsg(setWithdrawAddrMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mdist.DistributionModule).
		WithMsgType(mdist.DistributionSetWithdrawAddrMsgType).
		WithMsg(msg)
	return xplac
}

// EVM module

// Send coind by using evm client.
func (xplac *XplaClient) EvmSendCoin(sendCoinMsg types.SendCoinMsg) *XplaClient {
	msg, err := mevm.MakeSendCoinMsg(sendCoinMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mevm.EvmModule).
		WithMsgType(mevm.EvmSendCoinMsgType).
		WithMsg(msg)
	return xplac
}

// Deploy soldity contract.
func (xplac *XplaClient) DeploySolidityContract(deploySolContractMsg types.DeploySolContractMsg) *XplaClient {
	msg, err := mevm.MakeDeploySolContractMsg(deploySolContractMsg)
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mevm.EvmModule).
		WithMsgType(mevm.EvmDeploySolContractMsgType).
		WithMsg(msg)
	return xplac
}

// Invoke (as execute) solidity contract.
func (xplac *XplaClient) InvokeSolidityContract(invokeSolContractMsg types.InvokeSolContractMsg) *XplaClient {
	msg, err := mevm.MakeInvokeSolContractMsg(invokeSolContractMsg)
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mevm.EvmModule).
		WithMsgType(mevm.EvmInvokeSolContractMsgType).
		WithMsg(msg)
	return xplac
}

// Feegrant module

// Grant fee allowance to an address.
func (xplac *XplaClient) FeeGrant(grantMsg types.FeeGrantMsg) *XplaClient {
	msg, err := mfeegrant.MakeFeeGrantMsg(grantMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mfeegrant.FeegrantModule).
		WithMsgType(mfeegrant.FeegrantGrantMsgType).
		WithMsg(msg)
	return xplac
}

// Revoke fee-grant.
func (xplac *XplaClient) RevokeFeeGrant(revokeGrantMsg types.RevokeFeeGrantMsg) *XplaClient {
	msg, err := mfeegrant.MakeRevokeFeeGrantMsg(revokeGrantMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mfeegrant.FeegrantModule).
		WithMsgType(mfeegrant.FeegrantRevokeGrantMsgType).
		WithMsg(msg)
	return xplac
}

// gov module

// Submit a proposal along with an initial deposit.
func (xplac *XplaClient) SubmitProposal(submitProposalMsg types.SubmitProposalMsg) *XplaClient {
	msg, err := mgov.MakeSubmitProposalMsg(submitProposalMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mgov.GovModule).
		WithMsgType(mgov.GovSubmitProposalMsgType).
		WithMsg(msg)
	return xplac
}

// Deposit tokens for an active proposal.
func (xplac *XplaClient) GovDeposit(govDepositMsg types.GovDepositMsg) *XplaClient {
	msg, err := mgov.MakeGovDepositMsg(govDepositMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mgov.GovModule).
		WithMsgType(mgov.GovDepositMsgType).
		WithMsg(msg)
	return xplac
}

// Vote for an active proposal, options: yes/no/no_with_veto/abstain.
func (xplac *XplaClient) Vote(voteMsg types.VoteMsg) *XplaClient {
	msg, err := mgov.MakeVoteMsg(voteMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mgov.GovModule).
		WithMsgType(mgov.GovVoteMsgType).
		WithMsg(msg)
	return xplac
}

// Vote for an active proposal, options: yes/no/no_with_veto/abstain.
func (xplac *XplaClient) WeightedVote(weightedVoteMsg types.WeightedVoteMsg) *XplaClient {
	msg, err := mgov.MakeWeightedVoteMsg(weightedVoteMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mgov.GovModule).
		WithMsgType(mgov.GovWeightedVoteMsgType).
		WithMsg(msg)
	return xplac
}

// Params module

// Submit a parameter change proposal.
func (xplac *XplaClient) ParamChange(paramChangeMsg types.ParamChangeMsg) *XplaClient {
	msg, err := mparams.MakeProposalParamChangeMsg(paramChangeMsg, xplac.GetPrivateKey(), xplac.GetEncoding())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mparams.ParamsModule).
		WithMsgType(mparams.ParamsProposalParamChangeMsgType).
		WithMsg(msg)
	return xplac
}

// Reward module

// Funds the fee collector with the specified amount
func (xplac *XplaClient) FundFeeCollector(fundFeeCollectorMsg types.FundFeeCollectorMsg) *XplaClient {
	msg, err := mreward.MakeFundFeeCollectorMsg(fundFeeCollectorMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mreward.RewardModule).
		WithMsgType(mreward.RewardFundFeeCollectorMsgType).
		WithMsg(msg)
	return xplac
}

// Slashing module

// Unjail validator previously jailed for downtime.
func (xplac *XplaClient) Unjail() *XplaClient {
	msg, err := mslashing.MakeUnjailMsg(xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mslashing.SlashingModule).
		WithMsgType(mslashing.SlahsingUnjailMsgType).
		WithMsg(msg)
	return xplac
}

// Staking module

// Create new validator initialized with a self-delegation to it.
func (xplac *XplaClient) CreateValidator(createValidatorMsg types.CreateValidatorMsg) *XplaClient {
	msg, err := mstaking.MakeCreateValidatorMsg(createValidatorMsg, xplac.GetPrivateKey(), xplac.GetOutputDocument())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mstaking.StakingModule).
		WithMsgType(mstaking.StakingCreateValidatorMsgType).
		WithMsg(msg)
	return xplac
}

// Edit an existing validator account.
func (xplac *XplaClient) EditValidator(editValidatorMsg types.EditValidatorMsg) *XplaClient {
	msg, err := mstaking.MakeEditValidatorMsg(editValidatorMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mstaking.StakingModule).
		WithMsgType(mstaking.StakingEditValidatorMsgType).
		WithMsg(msg)
	return xplac
}

// Delegate liquid tokens to a validator.
func (xplac *XplaClient) Delegate(delegateMsg types.DelegateMsg) *XplaClient {
	msg, err := mstaking.MakeDelegateMsg(delegateMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mstaking.StakingModule).
		WithMsgType(mstaking.StakingDelegateMsgType).
		WithMsg(msg)
	return xplac
}

// Unbond shares from a validator.
func (xplac *XplaClient) Unbond(unbondMsg types.UnbondMsg) *XplaClient {
	msg, err := mstaking.MakeUnbondMsg(unbondMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mstaking.StakingModule).
		WithMsgType(mstaking.StakingUnbondMsgType).
		WithMsg(msg)
	return xplac
}

// Redelegate illiquid tokens from one validator to another.
func (xplac *XplaClient) Redelegate(redelegateMsg types.RedelegateMsg) *XplaClient {
	msg, err := mstaking.MakeRedelegateMsg(redelegateMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mstaking.StakingModule).
		WithMsgType(mstaking.StakingRedelegateMsgType).
		WithMsg(msg)
	return xplac
}

// Upgrade module

// Submit a software upgrade proposal.
func (xplac *XplaClient) SoftwareUpgrade(softwareUpgradeMsg types.SoftwareUpgradeMsg) *XplaClient {
	msg, err := mupgrade.MakeProposalSoftwareUpgradeMsg(softwareUpgradeMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mupgrade.UpgradeModule).
		WithMsgType(mupgrade.UpgradeProposalSoftwareUpgradeMsgType).
		WithMsg(msg)
	return xplac
}

// Cancel the current software upgrade proposal.
func (xplac *XplaClient) CancelSoftwareUpgrade(cancelSoftwareUpgradeMsg types.CancelSoftwareUpgradeMsg) *XplaClient {
	msg, err := mupgrade.MakeCancelSoftwareUpgradeMsg(cancelSoftwareUpgradeMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mupgrade.UpgradeModule).
		WithMsgType(mupgrade.UpgradeCancelSoftwareUpgradeMsgType).
		WithMsg(msg)
	return xplac
}

// Wasm module

// Upload a wasm binary.
func (xplac *XplaClient) StoreCode(storeMsg types.StoreMsg) *XplaClient {
	addr, err := util.GetAddrByPrivKey(xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	msg, err := mwasm.MakeStoreCodeMsg(storeMsg, addr)
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mwasm.WasmModule).
		WithMsgType(mwasm.WasmStoreMsgType).
		WithMsg(msg)
	return xplac
}

// Instantiate a wasm contract.
func (xplac *XplaClient) InstantiateContract(instantiageMsg types.InstantiateMsg) *XplaClient {
	addr, err := util.GetAddrByPrivKey(xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	msg, err := mwasm.MakeInstantiateMsg(instantiageMsg, addr)
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mwasm.WasmModule).
		WithMsgType(mwasm.WasmInstantiateMsgType).
		WithMsg(msg)
	return xplac
}

// Execute a wasm contract.
func (xplac *XplaClient) ExecuteContract(executeMsg types.ExecuteMsg) *XplaClient {
	addr, err := util.GetAddrByPrivKey(xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	msg, err := mwasm.MakeExecuteMsg(executeMsg, addr)
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mwasm.WasmModule).
		WithMsgType(mwasm.WasmExecuteMsgType).
		WithMsg(msg)
	return xplac
}

// Clears admin for a contract to prevent further migrations.
func (xplac *XplaClient) ClearContractAdmin(clearContractAdminMsg types.ClearContractAdminMsg) *XplaClient {
	msg, err := mwasm.MakeClearContractAdminMsg(clearContractAdminMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mwasm.WasmModule).
		WithMsgType(mwasm.WasmClearContractAdminMsgType).
		WithMsg(msg)
	return xplac
}

// Set new admin for a contract.
func (xplac *XplaClient) SetContractAdmin(setContractAdminMsg types.SetContractAdminMsg) *XplaClient {
	msg, err := mwasm.MakeSetContractAdmintMsg(setContractAdminMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mwasm.WasmModule).
		WithMsgType(mwasm.WasmSetContractAdminMsgType).
		WithMsg(msg)
	return xplac
}

// Migrate a wasm contract to a new code version.
func (xplac *XplaClient) Migrate(migrateMsg types.MigrateMsg) *XplaClient {
	msg, err := mwasm.MakeMigrateMsg(migrateMsg, xplac.GetPrivateKey())
	if err != nil {
		return ResetModuleAndMsgXplac(xplac).WithErr(err)
	}
	xplac.WithModule(mwasm.WasmModule).
		WithMsgType(mwasm.WasmMigrateMsgType).
		WithMsg(msg)
	return xplac
}
