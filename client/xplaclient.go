package client

import (
	"context"

	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	grpc1 "github.com/gogo/protobuf/grpc"
	"github.com/xpladev/xpla/app/params"
	"google.golang.org/grpc"
)

// The xpla client is a client for performing all functions within the xpla.go library.
// The user mandatorily inputs chain ID.
type XplaClient struct {
	chainId        string
	encodingConfig params.EncodingConfig
	grpc           grpc1.ClientConn
	context        context.Context

	opts Options

	module  string
	msgType string
	msg     interface{}
	err     error
}

// Optional parameters of xpla client.
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

// Make new xpla client.
func NewXplaClient(
	chainId string,
) *XplaClient {
	var xplac XplaClient
	return xplac.
		WithChainId(chainId).
		WithEncoding(util.MakeEncodingConfig()).
		WithContext(context.Background())
}

// Set options of xpla client.
func (xplac *XplaClient) WithOptions(
	options Options,
) *XplaClient {
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
		WithOutputDocument(options.OutputDocument)
}

// Set chain ID
func (xplac *XplaClient) WithChainId(chainId string) *XplaClient {
	xplac.chainId = chainId
	return xplac
}

// Set encoding configuration
func (xplac *XplaClient) WithEncoding(encodingConfig params.EncodingConfig) *XplaClient {
	xplac.encodingConfig = encodingConfig
	return xplac
}

// Set xpla client context
func (xplac *XplaClient) WithContext(ctx context.Context) *XplaClient {
	xplac.context = ctx
	return xplac
}

// Set private key
func (xplac *XplaClient) WithPrivateKey(privateKey key.PrivateKey) *XplaClient {
	xplac.opts.PrivateKey = privateKey
	return xplac
}

// Set LCD URL
func (xplac *XplaClient) WithURL(lcdURL string) *XplaClient {
	xplac.opts.LcdURL = lcdURL
	return xplac
}

// Set GRPC URL to query or broadcast tx
func (xplac *XplaClient) WithGrpc(grpcUrl string) *XplaClient {
	connUrl := util.GrpcUrlParsing(grpcUrl)
	c, err := grpc.Dial(
		connUrl, grpc.WithInsecure(),
	)
	if err != nil {
		xplac.err = err
		return xplac
	}
	xplac.grpc = c
	xplac.opts.GrpcURL = grpcUrl
	return xplac
}

// Set RPC URL of tendermint core
func (xplac *XplaClient) WithRpc(rpcUrl string) *XplaClient {
	xplac.opts.RpcURL = rpcUrl
	return xplac
}

// Set RPC URL for evm module
func (xplac *XplaClient) WithEvmRpc(evmRpcUrl string) *XplaClient {
	xplac.opts.EvmRpcURL = evmRpcUrl
	return xplac
}

// Set broadcast mode
func (xplac *XplaClient) WithBroadcastMode(broadcastMode string) *XplaClient {
	xplac.opts.BroadcastMode = broadcastMode
	return xplac
}

// Set account number
func (xplac *XplaClient) WithAccountNumber(accountNumber string) *XplaClient {
	xplac.opts.AccountNumber = accountNumber
	return xplac
}

// Set account sequence
func (xplac *XplaClient) WithSequence(sequence string) *XplaClient {
	xplac.opts.Sequence = sequence
	return xplac
}

// Set gas limit
func (xplac *XplaClient) WithGasLimit(gasLimit string) *XplaClient {
	xplac.opts.GasLimit = gasLimit
	return xplac
}

// Set Gas price
func (xplac *XplaClient) WithGasPrice(gasPrice string) *XplaClient {
	xplac.opts.GasPrice = gasPrice
	return xplac
}

// Set Gas adjustment
func (xplac *XplaClient) WithGasAdjustment(gasAdjustment string) *XplaClient {
	xplac.opts.GasAdjustment = gasAdjustment
	return xplac
}

// Set fee amount
func (xplac *XplaClient) WithFeeAmount(feeAmount string) *XplaClient {
	xplac.opts.FeeAmount = feeAmount
	return xplac
}

// Set transaction sign mode
func (xplac *XplaClient) WithSignMode(signMode signing.SignMode) *XplaClient {
	xplac.opts.SignMode = signMode
	return xplac
}

// Set fee granter
func (xplac *XplaClient) WithFeeGranter(feeGranter sdk.AccAddress) *XplaClient {
	xplac.opts.FeeGranter = feeGranter
	return xplac
}

// Set timeout block height
func (xplac *XplaClient) WithTimeoutHeight(timeoutHeight string) *XplaClient {
	xplac.opts.TimeoutHeight = timeoutHeight
	return xplac
}

// Set pagination
func (xplac *XplaClient) WithPagination(pagination types.Pagination) *XplaClient {
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

	return xplac
}

// Set output document name
func (xplac *XplaClient) WithOutputDocument(outputDocument string) *XplaClient {
	xplac.opts.OutputDocument = outputDocument
	return xplac
}

// Set module name
func (xplac *XplaClient) WithModule(module string) *XplaClient {
	xplac.module = module
	return xplac
}

// Set message type of modules
func (xplac *XplaClient) WithMsgType(msgType string) *XplaClient {
	xplac.msgType = msgType
	return xplac
}

// Set message
func (xplac *XplaClient) WithMsg(msg interface{}) *XplaClient {
	xplac.msg = msg
	return xplac
}

// Set error
func (xplac *XplaClient) WithErr(err error) *XplaClient {
	xplac.err = err
	return xplac
}

// Get chain ID
func (xplac *XplaClient) GetChainId() string {
	return xplac.chainId
}

// Get private key
func (xplac *XplaClient) GetPrivateKey() key.PrivateKey {
	return xplac.opts.PrivateKey
}

// Get encoding configuration
func (xplac *XplaClient) GetEncoding() params.EncodingConfig {
	return xplac.encodingConfig
}

// Get xpla client context
func (xplac *XplaClient) GetContext() context.Context {
	return xplac.context
}

// Get LCD URL
func (xplac *XplaClient) GetLcdURL() string {
	return xplac.opts.LcdURL
}

// Get GRPC URL to query or broadcast tx
func (xplac *XplaClient) GetGrpcUrl() string {
	return xplac.opts.GrpcURL
}

// Get GRPC client connector
func (xplac *XplaClient) GetGrpcClient() grpc1.ClientConn {
	return xplac.grpc
}

// Get RPC URL of tendermint core
func (xplac *XplaClient) GetRpc() string {
	return xplac.opts.RpcURL
}

// Get RPC URL for evm module
func (xplac *XplaClient) GetEvmRpc() string {
	return xplac.opts.EvmRpcURL
}

// Get broadcast mode
func (xplac *XplaClient) GetBroadcastMode() string {
	return xplac.opts.BroadcastMode
}

// Get account number
func (xplac *XplaClient) GetAccountNumber() string {
	return xplac.opts.AccountNumber
}

// Get account sequence
func (xplac *XplaClient) GetSequence() string {
	return xplac.opts.Sequence
}

// Get gas limit
func (xplac *XplaClient) GetGasLimit() string {
	return xplac.opts.GasLimit
}

// Get Gas price
func (xplac *XplaClient) GetGasPrice() string {
	return xplac.opts.GasPrice
}

// Get Gas adjustment
func (xplac *XplaClient) GetGasAdjustment() string {
	return xplac.opts.GasAdjustment
}

// Get fee amount
func (xplac *XplaClient) GetFeeAmount() string {
	return xplac.opts.FeeAmount
}

// Get transaction sign mode
func (xplac *XplaClient) GetSignMode() signing.SignMode {
	return xplac.opts.SignMode
}

// Get fee granter
func (xplac *XplaClient) GetFeeGranter() sdk.AccAddress {
	return xplac.opts.FeeGranter
}

// Get timeout block height
func (xplac *XplaClient) GetTimeoutHeight() string {
	return xplac.opts.TimeoutHeight
}

// Get pagination
func (xplac *XplaClient) GetPagination() *query.PageRequest {
	return core.PageRequest
}

// Get output document name
func (xplac *XplaClient) GetOutputDocument() string {
	return xplac.opts.OutputDocument
}

// Get module name
func (xplac *XplaClient) GetModule() string {
	return xplac.module
}

// Get message type of modules
func (xplac *XplaClient) GetMsgType() string {
	return xplac.msgType
}

// Get message
func (xplac *XplaClient) GetMsg() interface{} {
	return xplac.msg
}

// Get error
func (xplac *XplaClient) GetErr() error {
	return xplac.err
}
