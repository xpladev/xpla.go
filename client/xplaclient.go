package client

import (
	"context"

	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/core/auth"
	"github.com/xpladev/xpla.go/core/authz"
	"github.com/xpladev/xpla.go/core/bank"
	"github.com/xpladev/xpla.go/core/base"
	"github.com/xpladev/xpla.go/core/crisis"
	"github.com/xpladev/xpla.go/core/distribution"
	"github.com/xpladev/xpla.go/core/evidence"
	"github.com/xpladev/xpla.go/core/evm"
	"github.com/xpladev/xpla.go/core/feegrant"
	"github.com/xpladev/xpla.go/core/gov"
	"github.com/xpladev/xpla.go/core/ibc"
	"github.com/xpladev/xpla.go/core/mint"
	"github.com/xpladev/xpla.go/core/params"
	"github.com/xpladev/xpla.go/core/reward"
	"github.com/xpladev/xpla.go/core/slashing"
	"github.com/xpladev/xpla.go/core/staking"
	"github.com/xpladev/xpla.go/core/upgrade"
	"github.com/xpladev/xpla.go/core/wasm"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	grpc1 "github.com/gogo/protobuf/grpc"
	paramsapp "github.com/xpladev/xpla/app/params"
	"google.golang.org/grpc"
)

var _ provider.XplaClient = &xplaClient{}

// The xpla client is a client for performing all functions within the xpla.go library.
// The user mandatorily inputs chain ID.
type xplaClient struct {
	chainId        string
	encodingConfig paramsapp.EncodingConfig
	grpc           grpc1.ClientConn
	context        context.Context

	opts provider.Options

	module  string
	msgType string
	msg     interface{}
	err     error

	externalCoreModule
}

// Make new xpla client.
func NewXplaClient(
	chainId string,
) provider.XplaClient {
	var xplac xplaClient
	return xplac.
		WithChainId(chainId).
		WithEncoding(util.MakeEncodingConfig()).
		WithContext(context.Background()).
		UpdateXplacInCoreModule()
}

// Set options of xpla client.
func (xplac *xplaClient) WithOptions(
	options provider.Options,
) provider.XplaClient {
	return xplac.
		WithPrivateKey(options.PrivateKey).
		WithAccountNumber(options.AccountNumber).
		WithBroadcastMode(options.BroadcastMode).
		WithSequence(options.Sequence).
		WithGasLimit(options.GasLimit).
		WithGasPrice(util.DenomRemove(options.GasPrice)).
		WithGasAdjustment(options.GasAdjustment).
		WithFeeAmount(options.FeeAmount).
		WithSignMode(options.SignMode).
		WithFeeGranter(options.FeeGranter).
		WithTimeoutHeight(options.TimeoutHeight).
		WithURL(options.LcdURL).
		WithGrpc(options.GrpcURL).
		WithRpc(options.RpcURL).
		WithEvmRpc(options.EvmRpcURL).
		WithPagination(options.Pagination).
		WithOutputDocument(options.OutputDocument).
		UpdateXplacInCoreModule()
}

// List of core modules.
type externalCoreModule struct {
	auth.AuthExternal
	authz.AuthzExternal
	bank.BankExternal
	base.BaseExternal
	crisis.CrisisExternal
	distribution.DistributionExternal
	evidence.EvidenceExternal
	evm.EvmExternal
	feegrant.FeegrantExternal
	gov.GovExternal
	ibc.IbcExternal
	mint.MintExternal
	params.ParamsExternal
	reward.RewardExternal
	slashing.SlashingExternal
	staking.StakingExternal
	upgrade.UpgradeExternal
	wasm.WasmExternal
}

// Update xpla client if data in the client are changed.
func (xplac *xplaClient) UpdateXplacInCoreModule() provider.XplaClient {
	xplac.externalCoreModule = externalCoreModule{
		auth.NewAuthExternal(xplac),
		authz.NewAuthzExternal(xplac),
		bank.NewBankExternal(xplac),
		base.NewBaseExternal(xplac),
		crisis.NewCrisisExternal(xplac),
		distribution.NewDistributionExternal(xplac),
		evidence.NewEvidenceExternal(xplac),
		evm.NewEvmExternal(xplac),
		feegrant.NewFeegrantExternal(xplac),
		gov.NewGovExternal(xplac),
		ibc.NewIbcExternal(xplac),
		mint.NewMintExternal(xplac),
		params.NewParamsExternal(xplac),
		reward.NewRewardExternal(xplac),
		slashing.NewSlashingExternal(xplac),
		staking.NewStakingExternal(xplac),
		upgrade.NewUpgradeExternal(xplac),
		wasm.NewWasmExternal(xplac),
	}
	return xplac
}

// Set chain ID
func (xplac *xplaClient) WithChainId(chainId string) provider.XplaClient {
	xplac.chainId = chainId
	return xplac.UpdateXplacInCoreModule()
}

// Set encoding configuration
func (xplac *xplaClient) WithEncoding(encodingConfig paramsapp.EncodingConfig) provider.XplaClient {
	xplac.encodingConfig = encodingConfig
	return xplac.UpdateXplacInCoreModule()
}

// Set xpla client context
func (xplac *xplaClient) WithContext(ctx context.Context) provider.XplaClient {
	xplac.context = ctx
	return xplac.UpdateXplacInCoreModule()
}

// Set private key
func (xplac *xplaClient) WithPrivateKey(privateKey key.PrivateKey) provider.XplaClient {
	xplac.opts.PrivateKey = privateKey
	return xplac.UpdateXplacInCoreModule()
}

// Set LCD URL
func (xplac *xplaClient) WithURL(lcdURL string) provider.XplaClient {
	xplac.opts.LcdURL = lcdURL
	return xplac.UpdateXplacInCoreModule()
}

// Set GRPC URL to query or broadcast tx
func (xplac *xplaClient) WithGrpc(grpcUrl string) provider.XplaClient {
	connUrl := util.GrpcUrlParsing(grpcUrl)
	c, err := grpc.Dial(
		connUrl, grpc.WithInsecure(),
	)
	if err != nil {
		xplac.err = err
		return xplac.UpdateXplacInCoreModule()
	}
	xplac.grpc = c
	xplac.opts.GrpcURL = grpcUrl
	return xplac.UpdateXplacInCoreModule()
}

// Set RPC URL of tendermint core
func (xplac *xplaClient) WithRpc(rpcUrl string) provider.XplaClient {
	xplac.opts.RpcURL = rpcUrl
	return xplac.UpdateXplacInCoreModule()
}

// Set RPC URL for evm module
func (xplac *xplaClient) WithEvmRpc(evmRpcUrl string) provider.XplaClient {
	xplac.opts.EvmRpcURL = evmRpcUrl
	return xplac.UpdateXplacInCoreModule()
}

// Set broadcast mode
func (xplac *xplaClient) WithBroadcastMode(broadcastMode string) provider.XplaClient {
	xplac.opts.BroadcastMode = broadcastMode
	return xplac.UpdateXplacInCoreModule()
}

// Set account number
func (xplac *xplaClient) WithAccountNumber(accountNumber string) provider.XplaClient {
	xplac.opts.AccountNumber = accountNumber
	return xplac.UpdateXplacInCoreModule()
}

// Set account sequence
func (xplac *xplaClient) WithSequence(sequence string) provider.XplaClient {
	xplac.opts.Sequence = sequence
	return xplac.UpdateXplacInCoreModule()
}

// Set gas limit
func (xplac *xplaClient) WithGasLimit(gasLimit string) provider.XplaClient {
	xplac.opts.GasLimit = gasLimit
	return xplac.UpdateXplacInCoreModule()
}

// Set Gas price
func (xplac *xplaClient) WithGasPrice(gasPrice string) provider.XplaClient {
	xplac.opts.GasPrice = gasPrice
	return xplac.UpdateXplacInCoreModule()
}

// Set Gas adjustment
func (xplac *xplaClient) WithGasAdjustment(gasAdjustment string) provider.XplaClient {
	xplac.opts.GasAdjustment = gasAdjustment
	return xplac.UpdateXplacInCoreModule()
}

// Set fee amount
func (xplac *xplaClient) WithFeeAmount(feeAmount string) provider.XplaClient {
	xplac.opts.FeeAmount = feeAmount
	return xplac.UpdateXplacInCoreModule()
}

// Set transaction sign mode
func (xplac *xplaClient) WithSignMode(signMode signing.SignMode) provider.XplaClient {
	xplac.opts.SignMode = signMode
	return xplac.UpdateXplacInCoreModule()
}

// Set fee granter
func (xplac *xplaClient) WithFeeGranter(feeGranter sdk.AccAddress) provider.XplaClient {
	xplac.opts.FeeGranter = feeGranter
	return xplac.UpdateXplacInCoreModule()
}

// Set timeout block height
func (xplac *xplaClient) WithTimeoutHeight(timeoutHeight string) provider.XplaClient {
	xplac.opts.TimeoutHeight = timeoutHeight
	return xplac.UpdateXplacInCoreModule()
}

// Set pagination
func (xplac *xplaClient) WithPagination(pagination types.Pagination) provider.XplaClient {
	emptyPagination := types.Pagination{}
	if pagination != emptyPagination {
		pageReq, err := core.ReadPageRequest(pagination)
		if err != nil {
			xplac.err = err
		}
		core.PageRequest = pageReq
	} else {
		core.PageRequest = core.DefaultPagination()
	}

	return xplac.UpdateXplacInCoreModule()
}

// Set output document name
func (xplac *xplaClient) WithOutputDocument(outputDocument string) provider.XplaClient {
	xplac.opts.OutputDocument = outputDocument
	return xplac.UpdateXplacInCoreModule()
}

// Set module name
func (xplac *xplaClient) WithModule(module string) provider.XplaClient {
	xplac.module = module
	return xplac.UpdateXplacInCoreModule()
}

// Set message type of modules
func (xplac *xplaClient) WithMsgType(msgType string) provider.XplaClient {
	xplac.msgType = msgType
	return xplac.UpdateXplacInCoreModule()
}

// Set message
func (xplac *xplaClient) WithMsg(msg interface{}) provider.XplaClient {
	xplac.msg = msg
	return xplac.UpdateXplacInCoreModule()
}

// Set error
func (xplac *xplaClient) WithErr(err error) provider.XplaClient {
	xplac.err = err
	return xplac.UpdateXplacInCoreModule()
}

// Get chain ID
func (xplac *xplaClient) GetChainId() string {
	return xplac.chainId
}

// Get private key
func (xplac *xplaClient) GetPrivateKey() key.PrivateKey {
	return xplac.opts.PrivateKey
}

// Get encoding configuration
func (xplac *xplaClient) GetEncoding() paramsapp.EncodingConfig {
	return xplac.encodingConfig
}

// Get xpla client context
func (xplac *xplaClient) GetContext() context.Context {
	return xplac.context
}

// Get LCD URL
func (xplac *xplaClient) GetLcdURL() string {
	return xplac.opts.LcdURL
}

// Get GRPC URL to query or broadcast tx
func (xplac *xplaClient) GetGrpcUrl() string {
	return xplac.opts.GrpcURL
}

// Get GRPC client connector
func (xplac *xplaClient) GetGrpcClient() grpc1.ClientConn {
	return xplac.grpc
}

// Get RPC URL of tendermint core
func (xplac *xplaClient) GetRpc() string {
	return xplac.opts.RpcURL
}

// Get RPC URL for evm module
func (xplac *xplaClient) GetEvmRpc() string {
	return xplac.opts.EvmRpcURL
}

// Get broadcast mode
func (xplac *xplaClient) GetBroadcastMode() string {
	return xplac.opts.BroadcastMode
}

// Get account number
func (xplac *xplaClient) GetAccountNumber() string {
	return xplac.opts.AccountNumber
}

// Get account sequence
func (xplac *xplaClient) GetSequence() string {
	return xplac.opts.Sequence
}

// Get gas limit
func (xplac *xplaClient) GetGasLimit() string {
	return xplac.opts.GasLimit
}

// Get Gas price
func (xplac *xplaClient) GetGasPrice() string {
	return xplac.opts.GasPrice
}

// Get Gas adjustment
func (xplac *xplaClient) GetGasAdjustment() string {
	return xplac.opts.GasAdjustment
}

// Get fee amount
func (xplac *xplaClient) GetFeeAmount() string {
	return xplac.opts.FeeAmount
}

// Get transaction sign mode
func (xplac *xplaClient) GetSignMode() signing.SignMode {
	return xplac.opts.SignMode
}

// Get fee granter
func (xplac *xplaClient) GetFeeGranter() sdk.AccAddress {
	return xplac.opts.FeeGranter
}

// Get timeout block height
func (xplac *xplaClient) GetTimeoutHeight() string {
	return xplac.opts.TimeoutHeight
}

// Get pagination
func (xplac *xplaClient) GetPagination() *query.PageRequest {
	return core.PageRequest
}

// Get output document name
func (xplac *xplaClient) GetOutputDocument() string {
	return xplac.opts.OutputDocument
}

// Get module name
func (xplac *xplaClient) GetModule() string {
	return xplac.module
}

// Get message type of modules
func (xplac *xplaClient) GetMsgType() string {
	return xplac.msgType
}

// Get message
func (xplac *xplaClient) GetMsg() interface{} {
	return xplac.msg
}

// Get error
func (xplac *xplaClient) GetErr() error {
	return xplac.err
}
