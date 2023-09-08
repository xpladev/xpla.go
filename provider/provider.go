package provider

import (
	"context"

	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla/app/params"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gogo/protobuf/grpc"
)

// The standard form of XPLA client is interface type.
// XplaClient is endpoint in order to access xpla.go from external packages.
// If new modules are implemeted, external functions that are used to send tx or query state should be
// enrolled in XplaClient interface.
//
// e.g. - enroll bank module
//
//	  type TxMsgProvider interface {
//		...
//		BankSend(types.BankSendMsg) XplaClient
//		...
//	  }
//
//	  type QueryMsgProvider interface {
//		...
//		BankBalances(types.BankBalancesMsg) XplaClient
//		DenomMetadata(...types.DenomMetadataMsg) XplaClient
//		Total(...types.TotalMsg) XplaClient
//		...
//	  }
//
// The return type of these methods must be always the XplaClient because the client uses mehod chaining.
//
// e.g. - create and sign transaction
//
//	txbytes, err := xplac.BankSend(bankSendMsg).CreateAndSignTx()
//
// e.g. - query
//
//	res, err := xplac.BankBalances(bankBalancesMsg).Query()
type XplaClient interface {
	WithProvider
	GetProvider
	TxProvider
	QueryProvider
	BroadcastProvider
	InfoRequestProvider
	TxMsgProvider
	QueryMsgProvider
	HelperProvider
}

// Optional parameters of client.xplaClient.
type Options struct {
	PrivateKey     key.PrivateKey
	AccountNumber  string
	Sequence       string
	BroadcastMode  string
	GasLimit       string
	GasPrice       string
	GasAdjustment  string
	FeeAmount      string
	SignMode       signing.SignMode
	FeeGranter     sdk.AccAddress
	TimeoutHeight  string
	LcdURL         string
	GrpcURL        string
	RpcURL         string
	EvmRpcURL      string
	Pagination     types.Pagination
	OutputDocument string
}

// Methods set params of client.xplaClient.
type WithProvider interface {
	UpdateXplacInCoreModule() XplaClient
	WithOptions(Options) XplaClient
	WithChainId(string) XplaClient
	WithEncoding(params.EncodingConfig) XplaClient
	WithContext(context.Context) XplaClient
	WithPrivateKey(key.PrivateKey) XplaClient
	WithAccountNumber(string) XplaClient
	WithBroadcastMode(string) XplaClient
	WithSequence(string) XplaClient
	WithGasLimit(string) XplaClient
	WithGasPrice(string) XplaClient
	WithGasAdjustment(string) XplaClient
	WithFeeAmount(string) XplaClient
	WithSignMode(signing.SignMode) XplaClient
	WithFeeGranter(sdk.AccAddress) XplaClient
	WithTimeoutHeight(string) XplaClient
	WithURL(string) XplaClient
	WithGrpc(string) XplaClient
	WithRpc(string) XplaClient
	WithEvmRpc(string) XplaClient
	WithPagination(types.Pagination) XplaClient
	WithOutputDocument(string) XplaClient
	WithModule(module string) XplaClient
	WithMsgType(msgType string) XplaClient
	WithMsg(msg interface{}) XplaClient
	WithErr(err error) XplaClient
}

// Methods get params of client.xplaClient.
type GetProvider interface {
	GetChainId() string
	GetPrivateKey() key.PrivateKey
	GetEncoding() params.EncodingConfig
	GetContext() context.Context
	GetLcdURL() string
	GetGrpcUrl() string
	GetGrpcClient() grpc.ClientConn
	GetRpc() string
	GetEvmRpc() string
	GetBroadcastMode() string
	GetAccountNumber() string
	GetSequence() string
	GetGasLimit() string
	GetGasPrice() string
	GetGasAdjustment() string
	GetFeeAmount() string
	GetSignMode() signing.SignMode
	GetFeeGranter() sdk.AccAddress
	GetTimeoutHeight() string
	GetPagination() *query.PageRequest
	GetOutputDocument() string
	GetModule() string
	GetMsg() interface{}
	GetMsgType() string
	GetErr() error
}

// Methods handle transaction.
type TxProvider interface {
	CreateAndSignTx() ([]byte, error)
	CreateUnsignedTx() ([]byte, error)
	SignTx(types.SignTxMsg) ([]byte, error)
	MultiSign(types.TxMultiSignMsg) ([]byte, error)
	EncodeTx(types.EncodeTxMsg) (string, error)
	DecodeTx(types.DecodeTxMsg) (string, error)
	ValidateSignatures(types.ValidateSignaturesMsg) (string, error)
}

// Method handles query functions.
type QueryProvider interface {
	Query() (string, error)
}

// Methods handle functions of broadcasting.
type BroadcastProvider interface {
	Broadcast([]byte) (*types.TxRes, error)
	BroadcastBlock([]byte) (*types.TxRes, error)
	BroadcastAsync([]byte) (*types.TxRes, error)
}

// Methods get information from XPLA chain.
type InfoRequestProvider interface {
	LoadAccount(sdk.AccAddress) (authtypes.AccountI, error)
	Simulate(cmclient.TxBuilder) (*sdktx.SimulateResponse, error)
}

// Methods are external functions of each module for sending transaction.
type TxMsgProvider interface {
	// authz
	AuthzGrant(types.AuthzGrantMsg) XplaClient
	AuthzRevoke(types.AuthzRevokeMsg) XplaClient
	AuthzExec(types.AuthzExecMsg) XplaClient

	// bank
	BankSend(types.BankSendMsg) XplaClient

	// crisis
	InvariantBroken(types.InvariantBrokenMsg) XplaClient

	// distribution
	FundCommunityPool(types.FundCommunityPoolMsg) XplaClient
	CommunityPoolSpend(types.CommunityPoolSpendMsg) XplaClient
	WithdrawRewards(types.WithdrawRewardsMsg) XplaClient
	WithdrawAllRewards() XplaClient
	SetWithdrawAddr(types.SetWithdrawAddrMsg) XplaClient

	// evm
	EvmSendCoin(types.SendCoinMsg) XplaClient
	DeploySolidityContract(types.DeploySolContractMsg) XplaClient
	InvokeSolidityContract(types.InvokeSolContractMsg) XplaClient

	// feegrant
	FeeGrant(types.FeeGrantMsg) XplaClient
	RevokeFeeGrant(types.RevokeFeeGrantMsg) XplaClient

	// gov
	SubmitProposal(types.SubmitProposalMsg) XplaClient
	GovDeposit(types.GovDepositMsg) XplaClient
	Vote(types.VoteMsg) XplaClient
	WeightedVote(types.WeightedVoteMsg) XplaClient

	// params
	ParamChange(types.ParamChangeMsg) XplaClient

	// reward
	FundFeeCollector(types.FundFeeCollectorMsg) XplaClient

	// slashing
	Unjail() XplaClient

	// staking
	CreateValidator(types.CreateValidatorMsg) XplaClient
	EditValidator(types.EditValidatorMsg) XplaClient
	Delegate(types.DelegateMsg) XplaClient
	Unbond(types.UnbondMsg) XplaClient
	Redelegate(types.RedelegateMsg) XplaClient

	// upgrade
	SoftwareUpgrade(types.SoftwareUpgradeMsg) XplaClient
	CancelSoftwareUpgrade(types.CancelSoftwareUpgradeMsg) XplaClient

	// wasm
	StoreCode(types.StoreMsg) XplaClient
	InstantiateContract(types.InstantiateMsg) XplaClient
	ExecuteContract(types.ExecuteMsg) XplaClient
	ClearContractAdmin(types.ClearContractAdminMsg) XplaClient
	SetContractAdmin(types.SetContractAdminMsg) XplaClient
	Migrate(types.MigrateMsg) XplaClient
}

// Methods are external functions of each module for querying.
type QueryMsgProvider interface {
	// auth
	AuthParams() XplaClient
	AccAddress(types.QueryAccAddressMsg) XplaClient
	Accounts() XplaClient
	TxsByEvents(types.QueryTxsByEventsMsg) XplaClient
	Tx(types.QueryTxMsg) XplaClient

	// authz
	QueryAuthzGrants(types.QueryAuthzGrantMsg) XplaClient

	// bank
	BankBalances(types.BankBalancesMsg) XplaClient
	DenomMetadata(...types.DenomMetadataMsg) XplaClient
	Total(...types.TotalMsg) XplaClient

	// base
	NodeInfo() XplaClient
	Syncing() XplaClient
	Block(...types.BlockMsg) XplaClient
	ValidatorSet(...types.ValidatorSetMsg) XplaClient

	// distribution
	DistributionParams() XplaClient
	ValidatorOutstandingRewards(types.ValidatorOutstandingRewardsMsg) XplaClient
	DistCommission(types.QueryDistCommissionMsg) XplaClient
	DistSlashes(types.QueryDistSlashesMsg) XplaClient
	DistRewards(types.QueryDistRewardsMsg) XplaClient
	CommunityPool() XplaClient

	// evidence
	QueryEvidence(...types.QueryEvidenceMsg) XplaClient

	// evm
	CallSolidityContract(types.CallSolContractMsg) XplaClient
	GetTransactionByHash(types.GetTransactionByHashMsg) XplaClient
	GetBlockByHashOrHeight(types.GetBlockByHashHeightMsg) XplaClient
	AccountInfo(types.AccountInfoMsg) XplaClient
	SuggestGasPrice() XplaClient
	EthChainID() XplaClient
	EthBlockNumber() XplaClient
	Web3ClientVersion() XplaClient
	Web3Sha3(types.Web3Sha3Msg) XplaClient
	NetVersion() XplaClient
	NetPeerCount() XplaClient
	NetListening() XplaClient
	EthProtocolVersion() XplaClient
	EthSyncing() XplaClient
	EthAccounts() XplaClient
	EthGetBlockTransactionCount(types.EthGetBlockTransactionCountMsg) XplaClient
	EstimateGas(types.InvokeSolContractMsg) XplaClient
	EthGetTransactionByBlockHashAndIndex(types.GetTransactionByBlockHashAndIndexMsg) XplaClient
	EthGetTransactionReceipt(types.GetTransactionReceiptMsg) XplaClient
	EthNewFilter(types.EthNewFilterMsg) XplaClient
	EthNewBlockFilter() XplaClient
	EthNewPendingTransactionFilter() XplaClient
	EthUninstallFilter(types.EthUninstallFilterMsg) XplaClient
	EthGetFilterChanges(types.EthGetFilterChangesMsg) XplaClient
	EthGetFilterLogs(types.EthGetFilterLogsMsg) XplaClient
	EthGetLogs(types.EthGetLogsMsg) XplaClient
	EthCoinbase() XplaClient

	// feegrant
	QueryFeeGrants(types.QueryFeeGrantMsg) XplaClient

	// gov
	QueryProposal(types.QueryProposalMsg) XplaClient
	QueryProposals(types.QueryProposalsMsg) XplaClient
	QueryDeposit(types.QueryDepositMsg) XplaClient
	QueryVote(types.QueryVoteMsg) XplaClient
	Tally(types.TallyMsg) XplaClient
	GovParams(...types.GovParamsMsg) XplaClient
	Proposer(types.ProposerMsg) XplaClient

	// mint
	MintParams() XplaClient
	Inflation() XplaClient
	AnnualProvisions() XplaClient

	// ibc
	IbcClientStates() XplaClient
	IbcClientState(types.IbcClientStateMsg) XplaClient
	IbcClientStatus(types.IbcClientStatusMsg) XplaClient
	IbcClientConsensusStates(types.IbcClientConsensusStatesMsg) XplaClient
	IbcClientConsensusStateHeights(types.IbcClientConsensusStateHeightsMsg) XplaClient
	IbcClientConsensusState(types.IbcClientConsensusStateMsg) XplaClient
	IbcClientHeader() XplaClient
	IbcClientSelfConsensusState() XplaClient
	IbcClientParams() XplaClient
	IbcConnections(...types.IbcConnectionMsg) XplaClient
	IbcClientConnections(types.IbcClientConnectionsMsg) XplaClient
	IbcChannels(...types.IbcChannelMsg) XplaClient
	IbcChannelConnections(types.IbcChannelConnectionsMsg) XplaClient
	IbcChannelClientState(types.IbcChannelClientStateMsg) XplaClient
	IbcChannelPacketCommitments(types.IbcChannelPacketCommitmentsMsg) XplaClient
	IbcChannelPacketReceipt(types.IbcChannelPacketReceiptMsg) XplaClient
	IbcChannelPacketAck(types.IbcChannelPacketAckMsg) XplaClient
	IbcChannelUnreceivedPackets(types.IbcChannelUnreceivedPacketsMsg) XplaClient
	IbcChannelUnreceivedAcks(types.IbcChannelUnreceivedAcksMsg) XplaClient
	IbcChannelNextSequence(types.IbcChannelNextSequenceMsg) XplaClient
	IbcDenomTraces(...types.IbcDenomTraceMsg) XplaClient
	IbcDenomTrace(types.IbcDenomTraceMsg) XplaClient
	IbcDenomHash(types.IbcDenomHashMsg) XplaClient
	IbcEscrowAddress(types.IbcEscrowAddressMsg) XplaClient
	IbcTransferParams() XplaClient

	// params
	QuerySubspace(types.SubspaceMsg) XplaClient

	// reward
	RewardParams() XplaClient
	RewardPool() XplaClient

	// slashing
	SlashingParams() XplaClient
	SigningInfos(...types.SigningInfoMsg) XplaClient

	// staking
	QueryValidators(...types.QueryValidatorMsg) XplaClient
	QueryDelegation(types.QueryDelegationMsg) XplaClient
	QueryUnbondingDelegation(types.QueryUnbondingDelegationMsg) XplaClient
	QueryRedelegation(types.QueryRedelegationMsg) XplaClient
	HistoricalInfo(types.HistoricalInfoMsg) XplaClient
	StakingPool() XplaClient
	StakingParams() XplaClient

	// upgrade
	UpgradeApplied(types.AppliedMsg) XplaClient
	ModulesVersion(...types.QueryModulesVersionMsg) XplaClient
	Plan() XplaClient

	// wasm
	QueryContract(types.QueryMsg) XplaClient
	ListCode() XplaClient
	ListContractByCode(types.ListContractByCodeMsg) XplaClient
	Download(types.DownloadMsg) XplaClient
	CodeInfo(types.CodeInfoMsg) XplaClient
	ContractInfo(types.ContractInfoMsg) XplaClient
	ContractStateAll(types.ContractStateAllMsg) XplaClient
	ContractHistory(types.ContractHistoryMsg) XplaClient
	Pinned() XplaClient
	LibwasmvmVersion() XplaClient
}

// Method of helper.
type HelperProvider interface {
	EncodedTxbytesToJsonTx([]byte) ([]byte, error)
}
