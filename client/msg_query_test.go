package client_test

import (
	"fmt"

	"github.com/xpladev/xpla.go/client"
	mauth "github.com/xpladev/xpla.go/core/auth"
	mauthz "github.com/xpladev/xpla.go/core/authz"
	mbank "github.com/xpladev/xpla.go/core/bank"
	mbase "github.com/xpladev/xpla.go/core/base"
	mdist "github.com/xpladev/xpla.go/core/distribution"
	mevidence "github.com/xpladev/xpla.go/core/evidence"
	mevm "github.com/xpladev/xpla.go/core/evm"
	mfeegrant "github.com/xpladev/xpla.go/core/feegrant"
	mgov "github.com/xpladev/xpla.go/core/gov"
	mibc "github.com/xpladev/xpla.go/core/ibc"
	mmint "github.com/xpladev/xpla.go/core/mint"
	mparams "github.com/xpladev/xpla.go/core/params"
	mreward "github.com/xpladev/xpla.go/core/reward"
	mslashing "github.com/xpladev/xpla.go/core/slashing"
	mstaking "github.com/xpladev/xpla.go/core/staking"
	mupgrade "github.com/xpladev/xpla.go/core/upgrade"
	mwasm "github.com/xpladev/xpla.go/core/wasm"
	"github.com/xpladev/xpla.go/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdktestutil "github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govcli "github.com/cosmos/cosmos-sdk/x/gov/client/cli"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var (
	testSolContractAddress = "0x80E123317190cAf36292A04776b0De020136526F"
	testABIPath            = "../util/testutil/test_files/abi.json"
	testBytecodePath       = "../util/testutil/test_files/bytecode.json"
	testIbcClientID        = "07-tendermint-0"
	testIbcConnectionID    = "connection-1"
	testIbcChannelID       = "channel-0"
	testIbcChannelPortId   = "transfer"
	testCWContractAddress  = "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h"
	testWasmFilePath       = "../util/testutil/test_files/cw721_metadata_onchain.wasm"
)

func (s *ClientTestSuite) TestAuth() {
	// auth params
	s.xplac.AuthParams()

	authParamMsg, err := mauth.MakeAuthParamMsg()
	s.Require().NoError(err)

	s.Require().Equal(authParamMsg, s.xplac.GetMsg())
	s.Require().Equal(mauth.AuthModule, s.xplac.GetModule())
	s.Require().Equal(mauth.AuthQueryParamsMsgType, s.xplac.GetMsgType())

	// acc address
	queryAccAddressMsg := types.QueryAccAddressMsg{
		Address: s.accounts[0].Address.String(),
	}
	s.xplac.AccAddress(queryAccAddressMsg)

	accAddressMsg, err := mauth.MakeQueryAccAddressMsg(queryAccAddressMsg)
	s.Require().NoError(err)

	s.Require().Equal(accAddressMsg, s.xplac.GetMsg())
	s.Require().Equal(mauth.AuthModule, s.xplac.GetModule())
	s.Require().Equal(mauth.AuthQueryAccAddressMsgType, s.xplac.GetMsgType())

	// accounts
	s.xplac.Accounts()

	accountsMsg, err := mauth.MakeQueryAccountsMsg()
	s.Require().NoError(err)

	s.Require().Equal(accountsMsg, s.xplac.GetMsg())
	s.Require().Equal(mauth.AuthModule, s.xplac.GetModule())
	s.Require().Equal(mauth.AuthQueryAccountsMsgType, s.xplac.GetMsgType())

	// txs by events
	queryTxsByEventsMsg := types.QueryTxsByEventsMsg{
		Events: "transfer.recipient=" + s.accounts[0].Address.String(),
	}
	s.xplac.TxsByEvents(queryTxsByEventsMsg)
	txsByEventMsg, err := mauth.MakeTxsByEventsMsg(queryTxsByEventsMsg)
	s.Require().NoError(err)

	s.Require().Equal(txsByEventMsg, s.xplac.GetMsg())
	s.Require().Equal(mauth.AuthModule, s.xplac.GetModule())
	s.Require().Equal(mauth.AuthQueryTxsByEventsMsgType, s.xplac.GetMsgType())

	// tx
	queryTxMsg := types.QueryTxMsg{
		Value: s.testTxHash,
	}
	s.xplac.Tx(queryTxMsg)

	txMsg, err := mauth.MakeQueryTxMsg(queryTxMsg)
	s.Require().NoError(err)

	s.Require().Equal(txMsg, s.xplac.GetMsg())
	s.Require().Equal(mauth.AuthModule, s.xplac.GetModule())
	s.Require().Equal(mauth.AuthQueryTxMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestAuthz() {
	// query authz grants
	queryAuthzGrantMsg := types.QueryAuthzGrantMsg{
		Grantee: s.accounts[0].Address.String(),
		Granter: s.accounts[1].Address.String(),
	}
	s.xplac.QueryAuthzGrants(queryAuthzGrantMsg)

	authzGrantsMsg, err := mauthz.MakeQueryAuthzGrantsMsg(queryAuthzGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(authzGrantsMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzQueryGrantMsgType, s.xplac.GetMsgType())

	// grants by grantee
	queryAuthzGrantMsg = types.QueryAuthzGrantMsg{
		Grantee: s.accounts[0].Address.String(),
	}
	s.xplac.QueryAuthzGrants(queryAuthzGrantMsg)

	authzGrantsByGranteeMsg, err := mauthz.MakeQueryAuthzGrantsByGranteeMsg(queryAuthzGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(authzGrantsByGranteeMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzQueryGrantsByGranteeMsgType, s.xplac.GetMsgType())

	// grants by granter
	queryAuthzGrantMsg = types.QueryAuthzGrantMsg{
		Granter: s.accounts[1].Address.String(),
	}
	s.xplac.QueryAuthzGrants(queryAuthzGrantMsg)

	authzGrantsByGranterMsg, err := mauthz.MakeQueryAuthzGrantsByGranterMsg(queryAuthzGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(authzGrantsByGranterMsg, s.xplac.GetMsg())
	s.Require().Equal(mauthz.AuthzModule, s.xplac.GetModule())
	s.Require().Equal(mauthz.AuthzQueryGrantsByGranterMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestBank() {
	// bank all balances
	bankBalancesMsg := types.BankBalancesMsg{
		Address: s.accounts[0].Address.String(),
	}
	s.xplac.BankBalances(bankBalancesMsg)

	makeBankAllBalancesMsg, err := mbank.MakeBankAllBalancesMsg(bankBalancesMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeBankAllBalancesMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankAllBalancesMsgType, s.xplac.GetMsgType())

	// bank balance denom
	bankBalancesMsg = types.BankBalancesMsg{
		Address: s.accounts[0].Address.String(),
		Denom:   types.XplaDenom,
	}
	s.xplac.BankBalances(bankBalancesMsg)

	makeBankBalanceMsg, err := mbank.MakeBankBalanceMsg(bankBalancesMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeBankBalanceMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankBalanceMsgType, s.xplac.GetMsgType())

	// denoms metadata
	s.xplac.DenomMetadata()

	makeDenomsMetaDataMsg, err := mbank.MakeDenomsMetaDataMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeDenomsMetaDataMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankDenomsMetadataMsgType, s.xplac.GetMsgType())

	// denom metadata
	denomMetadataMsg := types.DenomMetadataMsg{
		Denom: types.XplaDenom,
	}
	s.xplac.DenomMetadata(denomMetadataMsg)

	makeDenomMetaDataMsg, err := mbank.MakeDenomMetaDataMsg(denomMetadataMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeDenomMetaDataMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankDenomMetadataMsgType, s.xplac.GetMsgType())

	// total supply
	s.xplac.Total()

	makeTotalSupplyMsg, err := mbank.MakeTotalSupplyMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeTotalSupplyMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankTotalMsgType, s.xplac.GetMsgType())

	// supply of
	totalMsg := types.TotalMsg{
		Denom: types.XplaDenom,
	}
	s.xplac.Total(totalMsg)

	makeSupplyOfMsg, err := mbank.MakeSupplyOfMsg(totalMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeSupplyOfMsg, s.xplac.GetMsg())
	s.Require().Equal(mbank.BankModule, s.xplac.GetModule())
	s.Require().Equal(mbank.BankTotalSupplyOfMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestBase() {
	// node info
	s.xplac.NodeInfo()

	makeBaseNodeInfoMsg, err := mbase.MakeBaseNodeInfoMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeBaseNodeInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseNodeInfoMsgType, s.xplac.GetMsgType())

	// syncing
	s.xplac.Syncing()

	makeBaseSyncingMsg, err := mbase.MakeBaseSyncingMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeBaseSyncingMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseSyncingMsgType, s.xplac.GetMsgType())

	// latest block
	s.xplac.Block()

	makeBaseLatestBlockMsg, err := mbase.MakeBaseLatestBlockMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeBaseLatestBlockMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseLatestBlockMsgtype, s.xplac.GetMsgType())

	// block by height
	blockMsg := types.BlockMsg{
		Height: "1",
	}
	s.xplac.Block(blockMsg)

	makeBaseBlockByheightMsg, err := mbase.MakeBaseBlockByHeightMsg(blockMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeBaseBlockByheightMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseBlockByHeightMsgType, s.xplac.GetMsgType())

	// latest validator set
	s.xplac.ValidatorSet()

	makeLatestValidatorSetMsg, err := mbase.MakeLatestValidatorSetMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeLatestValidatorSetMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseLatestValidatorSetMsgType, s.xplac.GetMsgType())

	// validator set by height
	validatorSetMsg := types.ValidatorSetMsg{
		Height: "1",
	}
	s.xplac.ValidatorSet(validatorSetMsg)

	makeValidatorSetByHeightMsg, err := mbase.MakeValidatorSetByHeightMsg(validatorSetMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeValidatorSetByHeightMsg, s.xplac.GetMsg())
	s.Require().Equal(mbase.Base, s.xplac.GetModule())
	s.Require().Equal(mbase.BaseValidatorSetByHeightMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestDistribution() {
	val := s.network.Validators[0].ValAddress.String()
	// query dist params
	s.xplac.DistributionParams()

	makeQueryDistributionParamsMsg, err := mdist.MakeQueryDistributionParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDistributionParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQueryDistributionParamsMsgType, s.xplac.GetMsgType())

	// validator outstanding rewards
	validatorOutstandingRewardsMsg := types.ValidatorOutstandingRewardsMsg{
		ValidatorAddr: val,
	}
	s.xplac.ValidatorOutstandingRewards(validatorOutstandingRewardsMsg)

	makeValidatorOutstandingRewardsMsg, err := mdist.MakeValidatorOutstandingRewardsMsg(validatorOutstandingRewardsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeValidatorOutstandingRewardsMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionValidatorOutstandingRewardsMsgType, s.xplac.GetMsgType())

	// dist commission
	queryDistCommissionMsg := types.QueryDistCommissionMsg{
		ValidatorAddr: val,
	}
	s.xplac.DistCommission(queryDistCommissionMsg)

	makeQueryDistCommissionMsg, err := mdist.MakeQueryDistCommissionMsg(queryDistCommissionMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDistCommissionMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQueryDistCommissionMsgType, s.xplac.GetMsgType())

	// dist slashes
	queryDistSlashesMsg := types.QueryDistSlashesMsg{
		ValidatorAddr: val,
	}
	s.xplac.DistSlashes(queryDistSlashesMsg)

	makeQueryDistSlashesMsg, err := mdist.MakeQueryDistSlashesMsg(queryDistSlashesMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDistSlashesMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQuerySlashesMsgType, s.xplac.GetMsgType())

	// dist rewards
	queryDistRewardsMsg := types.QueryDistRewardsMsg{
		DelegatorAddr: s.accounts[0].Address.String(),
		ValidatorAddr: val,
	}
	s.xplac.DistRewards(queryDistRewardsMsg)

	makeQueryDistRewwradsMsg, err := mdist.MakeQueryDistRewardsMsg(queryDistRewardsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDistRewwradsMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQueryRewardsMsgType, s.xplac.GetMsgType())

	// total rewards
	queryDistRewardsMsg = types.QueryDistRewardsMsg{
		DelegatorAddr: s.accounts[0].Address.String(),
	}
	s.xplac.DistRewards(queryDistRewardsMsg)

	makeQueryDistTotalRewardsMsg, err := mdist.MakeQueryDistTotalRewardsMsg(queryDistRewardsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryDistTotalRewardsMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQueryTotalRewardsMsgType, s.xplac.GetMsgType())

	// community pool
	s.xplac.CommunityPool()

	makeQueryCommunityPoolMsg, err := mdist.MakeQueryCommunityPoolMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryCommunityPoolMsg, s.xplac.GetMsg())
	s.Require().Equal(mdist.DistributionModule, s.xplac.GetModule())
	s.Require().Equal(mdist.DistributionQueryCommunityPoolMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestEvidence() {
	// all evidence
	s.xplac.QueryEvidence()

	makeQueryAllEvidenceMsg, err := mevidence.MakeQueryAllEvidenceMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryAllEvidenceMsg, s.xplac.GetMsg())
	s.Require().Equal(mevidence.EvidenceModule, s.xplac.GetModule())
	s.Require().Equal(mevidence.EvidenceQueryAllMsgType, s.xplac.GetMsgType())

	// evidence
	queryEvidenceMsg := types.QueryEvidenceMsg{
		Hash: s.testTxHash,
	}
	s.xplac.QueryEvidence(queryEvidenceMsg)

	makeQueryEvidenceMsg, err := mevidence.MakeQueryEvidenceMsg(queryEvidenceMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryEvidenceMsg, s.xplac.GetMsg())
	s.Require().Equal(mevidence.EvidenceModule, s.xplac.GetModule())
	s.Require().Equal(mevidence.EvidenceQueryMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestEvm() {
	// call contract
	callSolContractMsg := types.CallSolContractMsg{
		ContractAddress:      testSolContractAddress,
		ContractFuncCallName: "retrieve",
		ABIJsonFilePath:      testABIPath,
		BytecodeJsonFilePath: testBytecodePath,
		FromByteAddress:      s.accounts[0].PubKey.Address().String(),
	}
	s.xplac.CallSolidityContract(callSolContractMsg)

	makeCallSolContractMsg, err := mevm.MakeCallSolContractMsg(callSolContractMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeCallSolContractMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmCallSolContractMsgType, s.xplac.GetMsgType())

	// tx by hash
	getTransactionByHashMsg := types.GetTransactionByHashMsg{
		TxHash: s.testTxHash,
	}
	s.xplac.GetTransactionByHash(getTransactionByHashMsg)

	makeGetTransactionByHashMsg, err := mevm.MakeGetTransactionByHashMsg(getTransactionByHashMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGetTransactionByHashMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmGetTransactionByHashMsgType, s.xplac.GetMsgType())

	// block by hash or height
	getBlockByHashHeightMsg := types.GetBlockByHashHeightMsg{
		BlockHeight: "1",
	}
	s.xplac.GetBlockByHashOrHeight(getBlockByHashHeightMsg)

	makeGetBlockByHashHeightMsg, err := mevm.MakeGetBlockByHashHeightMsg(getBlockByHashHeightMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGetBlockByHashHeightMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmGetBlockByHashHeightMsgType, s.xplac.GetMsgType())

	// account info
	accountInfoMsg := types.AccountInfoMsg{
		Account: s.accounts[0].PubKey.Address().String(),
	}
	s.xplac.AccountInfo(accountInfoMsg)

	makeQueryAccountInfoMsg, err := mevm.MakeQueryAccountInfoMsg(accountInfoMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryAccountInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmQueryAccountInfoMsgType, s.xplac.GetMsgType())

	// suggest gas price
	s.xplac.SuggestGasPrice()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmSuggestGasPriceMsgType, s.xplac.GetMsgType())

	// eth chain ID
	s.xplac.EthChainID()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmQueryChainIdMsgType, s.xplac.GetMsgType())

	// eth block number
	s.xplac.EthBlockNumber()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmQueryCurrentBlockNumberMsgType, s.xplac.GetMsgType())

	// web3 client version
	s.xplac.Web3ClientVersion()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmWeb3ClientVersionMsgType, s.xplac.GetMsgType())

	// web3 sha3
	web3Sha3Msg := types.Web3Sha3Msg{
		InputParam: "ABC",
	}
	s.xplac.Web3Sha3(web3Sha3Msg)

	makeWeb3Sha3Msg, err := mevm.MakeWeb3Sha3Msg(web3Sha3Msg)
	s.Require().NoError(err)

	s.Require().Equal(makeWeb3Sha3Msg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmWeb3Sha3MsgType, s.xplac.GetMsgType())

	// net version
	s.xplac.NetVersion()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmNetVersionMsgType, s.xplac.GetMsgType())

	// net peer count
	s.xplac.NetPeerCount()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmNetPeerCountMsgType, s.xplac.GetMsgType())

	// net listening
	s.xplac.NetListening()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmNetListeningMsgType, s.xplac.GetMsgType())

	// eth protocol version
	s.xplac.EthProtocolVersion()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthProtocolVersionMsgType, s.xplac.GetMsgType())

	// eth syncing
	s.xplac.EthSyncing()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthSyncingMsgType, s.xplac.GetMsgType())

	// eth accounts
	s.xplac.EthAccounts()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthAccountsMsgType, s.xplac.GetMsgType())

	// eth get block transaction count
	ethGetBlockTransactionCountMsg := types.EthGetBlockTransactionCountMsg{
		BlockHeight: "1",
	}
	s.xplac.EthGetBlockTransactionCount(ethGetBlockTransactionCountMsg)

	makeEthGetBlockTransactionCountMsg, err := mevm.MakeEthGetBlockTransactionCountMsg(ethGetBlockTransactionCountMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthGetBlockTransactionCountMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthGetBlockTransactionCountMsgType, s.xplac.GetMsgType())

	// estimate gas
	invokeSolContractMsg := types.InvokeSolContractMsg{
		ContractAddress:      testSolContractAddress,
		ABIJsonFilePath:      testABIPath,
		BytecodeJsonFilePath: testBytecodePath,
		FromByteAddress:      s.accounts[0].PubKey.Address().String(),
	}
	s.xplac.EstimateGas(invokeSolContractMsg)

	makeEstimateGasSolMsg, err := mevm.MakeEstimateGasSolMsg(invokeSolContractMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEstimateGasSolMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthEstimateGasMsgType, s.xplac.GetMsgType())

	// get tx by block hash and index
	getTransactionByBlockHashAndIndexMsg := types.GetTransactionByBlockHashAndIndexMsg{
		BlockHash: "1",
		Index:     "0",
	}
	s.xplac.EthGetTransactionByBlockHashAndIndex(getTransactionByBlockHashAndIndexMsg)

	makeGetTransactionByBlockHashAndIndexMsg, err := mevm.MakeGetTransactionByBlockHashAndIndexMsg(getTransactionByBlockHashAndIndexMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGetTransactionByBlockHashAndIndexMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmGetTransactionByBlockHashAndIndexMsgType, s.xplac.GetMsgType())

	// tx receipt
	getTransactionReceiptMsg := types.GetTransactionReceiptMsg{
		TransactionHash: s.testTxHash,
	}
	s.xplac.EthGetTransactionReceipt(getTransactionReceiptMsg)

	makeGetTransactionReceiptMsg, err := mevm.MakeGetTransactionReceiptMsg(getTransactionReceiptMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGetTransactionReceiptMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmGetTransactionReceiptMsgType, s.xplac.GetMsgType())

	// eth new fileter
	ethNewFilterMsg := types.EthNewFilterMsg{
		Topics:    []string{"0x20ec56d16231c4d7f761c2533885619489fface85cf6c478868ef1d531b93177"},
		Address:   []string{"0xf7777b36a51fb0b33dd0c5118361AfC94ff7f967"},
		ToBlock:   "latest",
		FromBlock: "earliest",
	}
	s.xplac.EthNewFilter(ethNewFilterMsg)

	makeEthNewFilterMsg, err := mevm.MakeEthNewFilterMsg(ethNewFilterMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthNewFilterMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthNewFilterMsgType, s.xplac.GetMsgType())

	// new block filter
	s.xplac.EthNewBlockFilter()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthNewBlockFilterMsgType, s.xplac.GetMsgType())

	// new pending transaction filter
	s.xplac.EthNewPendingTransactionFilter()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthNewPendingTransactionFilterMsgType, s.xplac.GetMsgType())

	// uninstall filter
	ethUninstallFilterMsg := types.EthUninstallFilterMsg{
		FilterId: "0x168b9d421ecbffa1ac706926c2203454",
	}
	s.xplac.EthUninstallFilter(ethUninstallFilterMsg)

	makeEthUninstallFilterMsg, err := mevm.MakeEthUninstallFilterMsg(ethUninstallFilterMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthUninstallFilterMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthUninstallFilterMsgType, s.xplac.GetMsgType())

	// filter changes
	ethGetFilterChangesMsg := types.EthGetFilterChangesMsg{
		FilterId: "0x168b9d421ecbffa1ac706926c2203454",
	}
	s.xplac.EthGetFilterChanges(ethGetFilterChangesMsg)

	makeEthGetFilterChangesMsg, err := mevm.MakeEthGetFilterChangesMsg(ethGetFilterChangesMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthGetFilterChangesMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthGetFilterChangesMsgType, s.xplac.GetMsgType())

	// eth filter logs
	ethGetFilterLogsMsg := types.EthGetFilterLogsMsg{
		FilterId: "0x168b9d421ecbffa1ac706926c2203454",
	}
	s.xplac.EthGetFilterLogs(ethGetFilterLogsMsg)

	makeEthGetFilterLogsMsg, err := mevm.MakeEthGetFilterLogsMsg(ethGetFilterLogsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthGetFilterLogsMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthGetFilterLogsMsgType, s.xplac.GetMsgType())

	// eth logs
	ethGetLogsMsg := types.EthGetLogsMsg{
		Topics:    []string{"0x20ec56d16231c4d7f761c2533885619489fface85cf6c478868ef1d531b93177"},
		Address:   []string{"0xf7777b36a51fb0b33dd0c5118361AfC94ff7f967"},
		ToBlock:   "latest",
		FromBlock: "latest",
	}
	s.xplac.EthGetLogs(ethGetLogsMsg)

	makeEthGetLogsMsg, err := mevm.MakeEthGetLogsMsg(ethGetLogsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeEthGetLogsMsg, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthGetLogsMsgType, s.xplac.GetMsgType())

	// coinbase
	s.xplac.EthCoinbase()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mevm.EvmModule, s.xplac.GetModule())
	s.Require().Equal(mevm.EvmEthCoinbaseMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestFeegrant() {
	// feegrant
	queryFeeGrantMsg := types.QueryFeeGrantMsg{
		Grantee: s.accounts[0].Address.String(),
		Granter: s.accounts[1].Address.String(),
	}
	s.xplac.QueryFeeGrants(queryFeeGrantMsg)

	makeQueryFeeGrantMsg, err := mfeegrant.MakeQueryFeeGrantMsg(queryFeeGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryFeeGrantMsg, s.xplac.GetMsg())
	s.Require().Equal(mfeegrant.FeegrantModule, s.xplac.GetModule())
	s.Require().Equal(mfeegrant.FeegrantQueryGrantMsgType, s.xplac.GetMsgType())

	// feegrant by grantee
	queryFeeGrantMsg = types.QueryFeeGrantMsg{
		Grantee: s.accounts[0].Address.String(),
	}
	s.xplac.QueryFeeGrants(queryFeeGrantMsg)

	makeQueryFeeGrantsByGranteeMsg, err := mfeegrant.MakeQueryFeeGrantsByGranteeMsg(queryFeeGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryFeeGrantsByGranteeMsg, s.xplac.GetMsg())
	s.Require().Equal(mfeegrant.FeegrantModule, s.xplac.GetModule())
	s.Require().Equal(mfeegrant.FeegrantQueryGrantsByGranteeMsgType, s.xplac.GetMsgType())

	// feegrant by granter
	queryFeeGrantMsg = types.QueryFeeGrantMsg{
		Granter: s.accounts[1].Address.String(),
	}
	s.xplac.QueryFeeGrants(queryFeeGrantMsg)

	makeQueryFeeGrantsByGranterMsg, err := mfeegrant.MakeQueryFeeGrantsByGranterMsg(queryFeeGrantMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryFeeGrantsByGranterMsg, s.xplac.GetMsg())
	s.Require().Equal(mfeegrant.FeegrantModule, s.xplac.GetModule())
	s.Require().Equal(mfeegrant.FeegrantQueryGrantsByGranterMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestGov() {
	val := s.network.Validators[0]

	_, err := MsgSubmitProposal(val.ClientCtx, val.Address.String(),
		"Text Proposal 1", "Where is the title!?", govtypes.ProposalTypeText,
		fmt.Sprintf("--%s=%s", govcli.FlagDeposit, sdk.NewCoin(s.cfg.BondDenom, govtypes.DefaultMinDepositTokens).String()))
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	_, err = MsgVote(val.ClientCtx, val.Address.String(), "1", "yes=0.6,no=0.3,abstain=0.05,no_with_veto=0.05")
	s.Require().NoError(err)
	s.Require().NoError(s.network.WaitForNextBlock())

	// query proposal
	queryProposalMsg := types.QueryProposalMsg{
		ProposalID: "1",
	}
	s.xplac.QueryProposal(queryProposalMsg)

	makeQueryProposalMsg, err := mgov.MakeQueryProposalMsg(queryProposalMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryProposalMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryProposalMsgType, s.xplac.GetMsgType())

	// query proposals
	queryProposalsMsg := types.QueryProposalsMsg{
		Status:    "DepositPeriod",
		Voter:     s.accounts[0].Address.String(),
		Depositor: s.accounts[1].Address.String(),
	}
	s.xplac.QueryProposals(queryProposalsMsg)

	makeQueryProposalsMsg, err := mgov.MakeQueryProposalsMsg(queryProposalsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryProposalsMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryProposalsMsgType, s.xplac.GetMsgType())

	var queryType int
	for i, api := range s.apis {
		if i == 0 {
			s.xplac.WithURL(api)
			queryType = types.QueryLcd
		} else {
			s.xplac.WithGrpc(api)
			queryType = types.QueryGrpc
		}

		// query deposit
		queryDepositMsg := types.QueryDepositMsg{
			ProposalID: "1",
			Depositor:  s.accounts[0].Address.String(),
		}
		s.xplac.QueryDeposit(queryDepositMsg)

		makeQueryDepositMsg, _, err := mgov.MakeQueryDepositMsg(queryDepositMsg, s.xplac.GetGrpcClient(), s.xplac.GetContext(), s.xplac.GetLcdURL(), queryType)
		s.Require().NoError(err)

		s.Require().Equal(makeQueryDepositMsg, s.xplac.GetMsg())
		s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
		s.Require().Equal(mgov.GovQueryDepositRequestMsgType, s.xplac.GetMsgType())

		// query deposits
		queryDepositMsg = types.QueryDepositMsg{
			ProposalID: "1",
		}
		s.xplac.QueryDeposit(queryDepositMsg)

		makeQueryDepositsMsg, _, err := mgov.MakeQueryDepositsMsg(queryDepositMsg, s.xplac.GetGrpcClient(), s.xplac.GetContext(), s.xplac.GetLcdURL(), queryType)
		s.Require().NoError(err)

		s.Require().Equal(makeQueryDepositsMsg, s.xplac.GetMsg())
		s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
		s.Require().Equal(mgov.GovQueryDepositsRequestMsgType, s.xplac.GetMsgType())

		// query vote
		queryVoteMsg := types.QueryVoteMsg{
			ProposalID: "1",
			VoterAddr:  val.Address.String(),
		}
		s.xplac.QueryVote(queryVoteMsg)

		makeQueryVoteMsg, err := mgov.MakeQueryVoteMsg(queryVoteMsg, s.xplac.GetGrpcClient(), s.xplac.GetContext(), s.xplac.GetLcdURL(), queryType)
		s.Require().NoError(err)

		s.Require().Equal(makeQueryVoteMsg, s.xplac.GetMsg())
		s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
		s.Require().Equal(mgov.GovQueryVoteMsgType, s.xplac.GetMsgType())

		// query votes
		queryVoteMsg = types.QueryVoteMsg{
			ProposalID: "1",
		}
		s.xplac.QueryVote(queryVoteMsg)

		makeQueryVotesMsg, _, err := mgov.MakeQueryVotesMsg(queryVoteMsg, s.xplac.GetGrpcClient(), s.xplac.GetContext(), s.xplac.GetLcdURL(), queryType)
		s.Require().NoError(err)

		s.Require().Equal(makeQueryVotesMsg, s.xplac.GetMsg())
		s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
		s.Require().Equal(mgov.GovQueryVotesPassedMsgType, s.xplac.GetMsgType())

		// tally
		tallyMsg := types.TallyMsg{
			ProposalID: "1",
		}
		s.xplac.Tally(tallyMsg)

		makeGovTallyMsg, err := mgov.MakeGovTallyMsg(tallyMsg, s.xplac.GetGrpcClient(), s.xplac.GetContext(), s.xplac.GetLcdURL(), queryType)
		s.Require().NoError(err)

		s.Require().Equal(makeGovTallyMsg, s.xplac.GetMsg())
		s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
		s.Require().Equal(mgov.GovTallyMsgType, s.xplac.GetMsgType())
	}
	s.xplac = client.ResetXplac(s.xplac)

	// gov params
	s.xplac.GovParams()

	s.Require().Equal(nil, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryGovParamsMsgType, s.xplac.GetMsgType())

	// gov params, paramtype voting
	govParamsMsg := types.GovParamsMsg{
		ParamType: "voting",
	}
	s.xplac.GovParams(govParamsMsg)

	makeGovParamsMsg, err := mgov.MakeGovParamsMsg(govParamsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGovParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryGovParamVotingMsgType, s.xplac.GetMsgType())

	// gov params, paramtype tallying
	govParamsMsg = types.GovParamsMsg{
		ParamType: "tallying",
	}
	s.xplac.GovParams(govParamsMsg)

	makeGovParamsMsg, err = mgov.MakeGovParamsMsg(govParamsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGovParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryGovParamTallyingMsgType, s.xplac.GetMsgType())

	// gov params, paramtype deposit
	govParamsMsg = types.GovParamsMsg{
		ParamType: "deposit",
	}
	s.xplac.GovParams(govParamsMsg)

	makeGovParamsMsg, err = mgov.MakeGovParamsMsg(govParamsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeGovParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryGovParamDepositMsgType, s.xplac.GetMsgType())

	// proposer
	proposerMsg := types.ProposerMsg{
		ProposalID: "1",
	}
	s.xplac.Proposer(proposerMsg)

	s.Require().Equal(proposerMsg.ProposalID, s.xplac.GetMsg())
	s.Require().Equal(mgov.GovModule, s.xplac.GetModule())
	s.Require().Equal(mgov.GovQueryProposerMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestMint() {
	// mint params
	s.xplac.MintParams()

	makeQueryMintParamsMsg, err := mmint.MakeQueryMintParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryMintParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mmint.MintModule, s.xplac.GetModule())
	s.Require().Equal(mmint.MintQueryMintParamsMsgType, s.xplac.GetMsgType())

	// inflation
	s.xplac.Inflation()

	makeQueryInflationMsg, err := mmint.MakeQueryInflationMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryInflationMsg, s.xplac.GetMsg())
	s.Require().Equal(mmint.MintModule, s.xplac.GetModule())
	s.Require().Equal(mmint.MintQueryInflationMsgType, s.xplac.GetMsgType())

	// annual provisions
	s.xplac.AnnualProvisions()

	makeQueryAnnualProvisionsMsg, err := mmint.MakeQueryAnnualProvisionsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryAnnualProvisionsMsg, s.xplac.GetMsg())
	s.Require().Equal(mmint.MintModule, s.xplac.GetModule())
	s.Require().Equal(mmint.MintQueryAnnualProvisionsMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestIBC() {
	// client states
	s.xplac.IbcClientStates()

	makeIbcClientStatesMsg, err := mibc.MakeIbcClientStatesMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientStatesMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientStatesMsgType, s.xplac.GetMsgType())

	// client state
	ibcClientStateMsg := types.IbcClientStateMsg{
		ClientId: testIbcClientID,
	}
	s.xplac.IbcClientState(ibcClientStateMsg)

	makeIbcClientStateMsg, err := mibc.MakeIbcClientStateMsg(ibcClientStateMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientStateMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientStateMsgType, s.xplac.GetMsgType())

	// client status
	ibcClientStatusMsg := types.IbcClientStatusMsg{
		ClientId: testIbcClientID,
	}
	s.xplac.IbcClientStatus(ibcClientStatusMsg)

	makeIbcClientStatusMsg, err := mibc.MakeIbcClientStatusMsg(ibcClientStatusMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientStatusMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientStatusMsgType, s.xplac.GetMsgType())

	// client consensus states
	ibcClientConsensusStatesMsg := types.IbcClientConsensusStatesMsg{
		ClientId: testIbcClientID,
	}
	s.xplac.IbcClientConsensusStates(ibcClientConsensusStatesMsg)

	makeIbcClientConsensusStatesMsg, err := mibc.MakeIbcClientConsensusStatesMsg(ibcClientConsensusStatesMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientConsensusStatesMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientConsensusStatesMsgType, s.xplac.GetMsgType())

	// client consensus state heights
	ibcClientConsensusStateHeightsMsg := types.IbcClientConsensusStateHeightsMsg{
		ClientId: testIbcClientID,
	}
	s.xplac.IbcClientConsensusStateHeights(ibcClientConsensusStateHeightsMsg)

	makeIbcClientConsensusStateHeightsMsg, err := mibc.MakeIbcClientConsensusStateHeightsMsg(ibcClientConsensusStateHeightsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientConsensusStateHeightsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientConsensusStateHeightsMsgType, s.xplac.GetMsgType())

	// client consensus state
	ibcClientConsensusStateMsg := types.IbcClientConsensusStateMsg{
		ClientId:     testIbcClientID,
		LatestHeight: false,
		Height:       "1-115",
	}
	s.xplac.IbcClientConsensusState(ibcClientConsensusStateMsg)

	makeIbcClientConsensusStateMsg, err := mibc.MakeIbcClientConsensusStateMsg(ibcClientConsensusStateMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientConsensusStateMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientConsensusStateMsgType, s.xplac.GetMsgType())

	// client params
	s.xplac.IbcClientParams()

	makeIbcClientParamsMsg, err := mibc.MakeIbcClientParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientParamsMsgType, s.xplac.GetMsgType())

	// connections
	s.xplac.IbcConnections()

	makeIbcConnectionConnectionsMsg, err := mibc.MakeIbcConnectionConnectionsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcConnectionConnectionsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcConnectionConnectionsMsgType, s.xplac.GetMsgType())

	// connection
	ibcConnectionMsg := types.IbcConnectionMsg{
		ConnectionId: testIbcConnectionID,
	}
	s.xplac.IbcConnections(ibcConnectionMsg)

	makeIbcConnectionConnectionMsg, err := mibc.MakeIbcConnectionConnectionMsg(ibcConnectionMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcConnectionConnectionMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcConnectionConnectionMsgType, s.xplac.GetMsgType())

	// client connection
	ibcClientConnectionsMsg := types.IbcClientConnectionsMsg{
		ClientId: testIbcClientID,
	}
	s.xplac.IbcClientConnections(ibcClientConnectionsMsg)

	makeIbcConnectionClientConnectionsMsg, err := mibc.MakeIbcConnectionClientConnectionsMsg(ibcClientConnectionsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcConnectionClientConnectionsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcConnectionClientConnectionsMsgType, s.xplac.GetMsgType())

	// channels
	s.xplac.IbcChannels()

	makeIbcChannelChannelsMsg, err := mibc.MakeIbcChannelChannelsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelChannelsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelChannelsMsgType, s.xplac.GetMsgType())

	// channel
	ibcChannelMsg := types.IbcChannelMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannels(ibcChannelMsg)

	makeIbcChannelChannelMsg, err := mibc.MakeIbcChannelChannelMsg(ibcChannelMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelChannelMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelChannelMsgType, s.xplac.GetMsgType())

	// channel connections
	ibcChannelConnectionsMsg := types.IbcChannelConnectionsMsg{
		ConnectionId: testIbcConnectionID,
	}
	s.xplac.IbcChannelConnections(ibcChannelConnectionsMsg)

	makeIbcChannelConnectionsMsg, err := mibc.MakeIbcChannelConnectionsMsg(ibcChannelConnectionsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelConnectionsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelConnectionsMsgType, s.xplac.GetMsgType())

	// channel client state
	ibcChannelClientStateMsg := types.IbcChannelClientStateMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelClientState(ibcChannelClientStateMsg)

	makeIbcChannelClientStateMsg, err := mibc.MakeIbcChannelClientStateMsg(ibcChannelClientStateMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelClientStateMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelClientStateMsgType, s.xplac.GetMsgType())

	// channel packet commitments
	ibcChannelPacketCommitmentsMsg := types.IbcChannelPacketCommitmentsMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelPacketCommitments(ibcChannelPacketCommitmentsMsg)

	makeIbcChannelPacketCommitmentsMsg, err := mibc.MakeIbcChannelPacketCommitmentsMsg(ibcChannelPacketCommitmentsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketCommitmentsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelPacketCommitmentsMsgType, s.xplac.GetMsgType())

	// channel packet commitment
	ibcChannelPacketCommitmentsMsg = types.IbcChannelPacketCommitmentsMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
		Sequence:  "1",
	}
	s.xplac.IbcChannelPacketCommitments(ibcChannelPacketCommitmentsMsg)

	makeIbcChannelPacketCommitmentMsg, err := mibc.MakeIbcChannelPacketCommitmentMsg(ibcChannelPacketCommitmentsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketCommitmentMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelPacketCommitmentMsgType, s.xplac.GetMsgType())

	// packet receipt
	ibcChannelPacketReceiptMsg := types.IbcChannelPacketReceiptMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelPacketReceipt(ibcChannelPacketReceiptMsg)

	makeIbcChannelPacketReceiptMsg, err := mibc.MakeIbcChannelPacketReceiptMsg(ibcChannelPacketReceiptMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketReceiptMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelPacketReceiptMsgType, s.xplac.GetMsgType())

	// packet ack
	ibcChannelPacketAckMsg := types.IbcChannelPacketAckMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelPacketAck(ibcChannelPacketAckMsg)

	makeIbcChannelPacketAckMsg, err := mibc.MakeIbcChannelPacketAckMsg(ibcChannelPacketAckMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketAckMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelPacketAckMsgType, s.xplac.GetMsgType())

	// unreceived packets
	ibcChannelUnreceivedPacketsMsg := types.IbcChannelUnreceivedPacketsMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelUnreceivedPackets(ibcChannelUnreceivedPacketsMsg)

	makeIbcChannelPacketUnreceivedPacketsMsg, err := mibc.MakeIbcChannelPacketUnreceivedPacketsMsg(ibcChannelUnreceivedPacketsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketUnreceivedPacketsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelUnreceivedPacketsMsgType, s.xplac.GetMsgType())

	// unreceived acks
	ibcChannelUnreceivedAcksMsg := types.IbcChannelUnreceivedAcksMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelUnreceivedAcks(ibcChannelUnreceivedAcksMsg)

	makeIbcChannelPacketUnreceivedAcksMsg, err := mibc.MakeIbcChannelPacketUnreceivedAcksMsg(ibcChannelUnreceivedAcksMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketUnreceivedAcksMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelUnreceivedAcksMsgType, s.xplac.GetMsgType())

	// channel next sequence
	ibcChannelNextSequenceMsg := types.IbcChannelNextSequenceMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelNextSequence(ibcChannelNextSequenceMsg)

	makeIbcChannelNextSequenceReceiveMsg, err := mibc.MakeIbcChannelNextSequenceReceiveMsg(ibcChannelNextSequenceMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelNextSequenceReceiveMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelNextSequenceMsgType, s.xplac.GetMsgType())

	// denom traces
	s.xplac.IbcDenomTraces()

	makeIbcTransferDenomTracesMsg, err := mibc.MakeIbcTransferDenomTracesMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferDenomTracesMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferDenomTracesMsgType, s.xplac.GetMsgType())

	// denom trace
	ibcDenomTraceMsg := types.IbcDenomTraceMsg{
		HashDenom: "B249D1E86F588286FEA286AA8364FFCE69EC65604BD7869D824ADE40F00FA25B",
	}
	s.xplac.IbcDenomTraces(ibcDenomTraceMsg)

	makeIbcTransferDenomTraceMsg, err := mibc.MakeIbcTransferDenomTraceMsg(ibcDenomTraceMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferDenomTraceMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferDenomTraceMsgType, s.xplac.GetMsgType())

	// denom trace
	ibcDenomTraceMsg = types.IbcDenomTraceMsg{
		HashDenom: "B249D1E86F588286FEA286AA8364FFCE69EC65604BD7869D824ADE40F00FA25B",
	}
	s.xplac.IbcDenomTrace(ibcDenomTraceMsg)

	makeIbcTransferDenomTraceMsg, err = mibc.MakeIbcTransferDenomTraceMsg(ibcDenomTraceMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferDenomTraceMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferDenomTraceMsgType, s.xplac.GetMsgType())

	// denom hash
	ibcDenomHashMsg := types.IbcDenomHashMsg{
		Trace: testIbcChannelPortId + "/" + testIbcChannelID + "/" + types.XplaDenom,
	}
	s.xplac.IbcDenomHash(ibcDenomHashMsg)

	makeIbcTransferDenomHashMsg, err := mibc.MakeIbcTransferDenomHashMsg(ibcDenomHashMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferDenomHashMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferDenomHashMsgType, s.xplac.GetMsgType())

	// escrow address
	ibcEscrowAddressMsg := types.IbcEscrowAddressMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcEscrowAddress(ibcEscrowAddressMsg)

	makeIbcTransferEscrowAddressMsg, err := mibc.MakeIbcTransferEscrowAddressMsg(ibcEscrowAddressMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferEscrowAddressMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferEscrowAddressMsgType, s.xplac.GetMsgType())

	// escrow address
	s.xplac.IbcTransferParams()

	makeIbcTransferParamsMsg, err := mibc.MakeIbcTransferParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferParamsMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestParams() {
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

func (s *ClientTestSuite) TestReward() {
	// reward params
	s.xplac.RewardParams()

	makeQueryRewardParamsMsg, err := mreward.MakeQueryRewardParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryRewardParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mreward.RewardModule, s.xplac.GetModule())
	s.Require().Equal(mreward.RewardQueryRewardParamsMsgType, s.xplac.GetMsgType())

	// reward pool
	s.xplac.RewardPool()

	makeQueryRewardPoolMsg, err := mreward.MakeQueryRewardPoolMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryRewardPoolMsg, s.xplac.GetMsg())
	s.Require().Equal(mreward.RewardModule, s.xplac.GetModule())
	s.Require().Equal(mreward.RewardQueryRewardPoolMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestSlashing() {
	// slashing params
	s.xplac.SlashingParams()

	makeQuerySlashingParamsMsg, err := mslashing.MakeQuerySlashingParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQuerySlashingParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mslashing.SlashingModule, s.xplac.GetModule())
	s.Require().Equal(mslashing.SlashingQuerySlashingParamsMsgType, s.xplac.GetMsgType())

	// signing infos
	s.xplac.SigningInfos()

	makeQuerySigningInfosMsg, err := mslashing.MakeQuerySigningInfosMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQuerySigningInfosMsg, s.xplac.GetMsg())
	s.Require().Equal(mslashing.SlashingModule, s.xplac.GetModule())
	s.Require().Equal(mslashing.SlashingQuerySigningInfosMsgType, s.xplac.GetMsgType())

	// signing info
	signingInfoMsg := types.SigningInfoMsg{
		ConsPubKey: `{"@type": "/cosmos.crypto.ed25519.PubKey","key": "6RBPm24ckoWhRt8mArcSCnEKvt0FMGvcaMwchfZ3ue8="}`,
	}
	s.xplac.SigningInfos(signingInfoMsg)

	makeQuerySigningInfoMsg, err := mslashing.MakeQuerySigningInfoMsg(signingInfoMsg, s.xplac.GetEncoding())
	s.Require().NoError(err)

	s.Require().Equal(makeQuerySigningInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mslashing.SlashingModule, s.xplac.GetModule())
	s.Require().Equal(mslashing.SlashingQuerySigningInfoMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestStaking() {
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

func (s *ClientTestSuite) TestUpgrade() {
	// upgrade applied
	appliedMsg := types.AppliedMsg{
		UpgradeName: "upgrade name",
	}
	s.xplac.UpgradeApplied(appliedMsg)

	makeAppliedMsg, err := mupgrade.MakeAppliedMsg(appliedMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeAppliedMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradeAppliedMsgType, s.xplac.GetMsgType())

	// modules version
	s.xplac.ModulesVersion()

	makeQueryAllModuleVersionMsg, err := mupgrade.MakeQueryAllModuleVersionMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeQueryAllModuleVersionMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradeQueryAllModuleVersionsMsgType, s.xplac.GetMsgType())

	// module version
	queryModulesVersionMsg := types.QueryModulesVersionMsg{
		ModuleName: "staking",
	}
	s.xplac.ModulesVersion(queryModulesVersionMsg)

	makeQueryModuleVersionMsg, err := mupgrade.MakeQueryModuleVersionMsg(queryModulesVersionMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryModuleVersionMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradeQueryModuleVersionsMsgType, s.xplac.GetMsgType())

	// plan
	s.xplac.Plan()

	makePlanMsg, err := mupgrade.MakePlanMsg()
	s.Require().NoError(err)

	s.Require().Equal(makePlanMsg, s.xplac.GetMsg())
	s.Require().Equal(mupgrade.UpgradeModule, s.xplac.GetModule())
	s.Require().Equal(mupgrade.UpgradePlanMsgType, s.xplac.GetMsgType())
}

func (s *ClientTestSuite) TestWasm() {
	// call contract
	queryMsg := types.QueryMsg{
		ContractAddress: testCWContractAddress,
		QueryMsg:        `{"query_method":{"query":"query_test"}}`,
	}
	s.xplac.QueryContract(queryMsg)

	makeQueryMsg, err := mwasm.MakeQueryMsg(queryMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeQueryMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmQueryContractMsgType, s.xplac.GetMsgType())

	// list code
	s.xplac.ListCode()

	makeListcodeMsg, err := mwasm.MakeListcodeMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeListcodeMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmListCodeMsgType, s.xplac.GetMsgType())

	// list contract by code
	listContractByCodeMsg := types.ListContractByCodeMsg{
		CodeId: "1",
	}
	s.xplac.ListContractByCode(listContractByCodeMsg)

	makeListContractByCodeMsg, err := mwasm.MakeListContractByCodeMsg(listContractByCodeMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeListContractByCodeMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmListContractByCodeMsgType, s.xplac.GetMsgType())

	// download
	downloadMsg := types.DownloadMsg{
		CodeId:           "1",
		DownloadFileName: "./example.json",
	}
	s.xplac.Download(downloadMsg)

	makeDownloadMsg, err := mwasm.MakeDownloadMsg(downloadMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeDownloadMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmDownloadMsgType, s.xplac.GetMsgType())

	// code info
	codeInfoMsg := types.CodeInfoMsg{
		CodeId: "1",
	}
	s.xplac.CodeInfo(codeInfoMsg)

	makeCodeInfoMsg, err := mwasm.MakeCodeInfoMsg(codeInfoMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeCodeInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmCodeInfoMsgType, s.xplac.GetMsgType())

	// contract info
	contractInfoMsg := types.ContractInfoMsg{
		ContractAddress: testCWContractAddress,
	}
	s.xplac.ContractInfo(contractInfoMsg)

	makeContractInfoMsg, err := mwasm.MakeContractInfoMsg(contractInfoMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeContractInfoMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmContractInfoMsgType, s.xplac.GetMsgType())

	// contract state all
	contractStateAllMsg := types.ContractStateAllMsg{
		ContractAddress: testCWContractAddress,
	}
	s.xplac.ContractStateAll(contractStateAllMsg)

	makeContractStateAllMsg, err := mwasm.MakeContractStateAllMsg(contractStateAllMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeContractStateAllMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmContractStateAllMsgType, s.xplac.GetMsgType())

	// contract history
	contractHistoryMsg := types.ContractHistoryMsg{
		ContractAddress: testCWContractAddress,
	}
	s.xplac.ContractHistory(contractHistoryMsg)

	makeContractHistoryMsg, err := mwasm.MakeContractHistoryMsg(contractHistoryMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeContractHistoryMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmContractHistoryMsgType, s.xplac.GetMsgType())

	// pinned
	s.xplac.Pinned()

	makePinnedMsg, err := mwasm.MakePinnedMsg()
	s.Require().NoError(err)

	s.Require().Equal(makePinnedMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmPinnedMsgType, s.xplac.GetMsgType())

	// libwasmvm version
	s.xplac.LibwasmvmVersion()

	makeLibwasmvmVersionMsg, err := mwasm.MakeLibwasmvmVersionMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeLibwasmvmVersionMsg, s.xplac.GetMsg())
	s.Require().Equal(mwasm.WasmModule, s.xplac.GetModule())
	s.Require().Equal(mwasm.WasmLibwasmvmVersionMsgType, s.xplac.GetMsgType())
}

var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(types.XplaDenom, sdk.NewInt(10))).String()),
}

// MsgSubmitProposal creates a tx for submit proposal
func MsgSubmitProposal(clientCtx cmclient.Context, from, title, description, proposalType string, extraArgs ...string) (sdktestutil.BufferWriter, error) {
	args := append([]string{
		fmt.Sprintf("--%s=%s", govcli.FlagTitle, title),
		fmt.Sprintf("--%s=%s", govcli.FlagDescription, description),
		fmt.Sprintf("--%s=%s", govcli.FlagProposalType, proposalType),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdSubmitProposal(), args)
}

// MsgVote votes for a proposal
func MsgVote(clientCtx cmclient.Context, from, id, vote string, extraArgs ...string) (sdktestutil.BufferWriter, error) {
	args := append([]string{
		id,
		vote,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, from),
	}, commonArgs...)

	args = append(args, extraArgs...)

	return clitestutil.ExecTestCLICmd(clientCtx, govcli.NewCmdWeightedVote(), args)
}
