package core

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	"github.com/gogo/protobuf/proto"
)

// Query internal XPLA client
type QueryClient struct {
	Ixplac    provider.XplaClient
	QueryType uint8
}

func NewIxplaClient(moduleClient provider.XplaClient, qt uint8) *QueryClient {
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
