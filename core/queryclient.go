package core

import (
	"context"

	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/gogo/protobuf/grpc"
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla/app/params"
)

// Query internal XPLA client
type QueryClient struct {
	Ixplac    ModuleClient
	QueryType uint8
}

type ModuleClient interface {
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
}

func NewIXplaClient(moduleClient ModuleClient, qt uint8) *QueryClient {
	return &QueryClient{Ixplac: moduleClient, QueryType: qt}
}

// Print protobuf message by using cosmos sdk codec.
func PrintProto(i QueryClient, toPrint proto.Message) ([]byte, error) {
	out, err := i.Ixplac.GetEncoding().Marshaler.MarshalJSON(toPrint)
	if err != nil {
		return nil, util.LogErr(errors.ErrFailedToMarshal, err)
	}
	return out, nil
}

// Print object by using cosmos sdk legacy amino.
func PrintObjectLegacy(i QueryClient, toPrint interface{}) ([]byte, error) {
	out, err := i.Ixplac.GetEncoding().Amino.MarshalJSON(toPrint)
	if err != nil {
		return nil, util.LogErr(errors.ErrFailedToMarshal, err)
	}
	return out, nil
}

// For auth module and gov module, make cosmos sdk client for querying.
func ClientForQuery(i QueryClient) (cmclient.Context, error) {
	client, err := cmclient.NewClientFromNode(i.Ixplac.GetRpc())
	if err != nil {
		return cmclient.Context{}, util.LogErr(errors.ErrSdkClient, err)
	}

	clientCtx, err := util.NewClient()
	if err != nil {
		return cmclient.Context{}, err
	}

	clientCtx = clientCtx.
		WithNodeURI(i.Ixplac.GetRpc()).
		WithClient(client)

	return clientCtx, nil
}
