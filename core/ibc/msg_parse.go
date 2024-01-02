package ibc

import (
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	ibcclient "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
)

// Parsing - IBC client consensus state
func parseIbcClientConsensusStateArgs(ibcClientConsensusStateMsg types.IbcClientConsensusStateMsg) (ibcclient.QueryConsensusStateRequest, error) {
	var height ibcclient.Height
	var err error

	if !ibcClientConsensusStateMsg.LatestHeight {
		if ibcClientConsensusStateMsg.Height == "" {
			return ibcclient.QueryConsensusStateRequest{}, util.LogErr(errors.ErrParse, "must include a second 'Height' argument when 'LatestHeight' is not provided")
		}

		height, err = ibcclient.ParseHeight(ibcClientConsensusStateMsg.Height)
		if err != nil {
			return ibcclient.QueryConsensusStateRequest{}, util.LogErr(errors.ErrParse, err)
		}
	}

	return ibcclient.QueryConsensusStateRequest{
		ClientId:       ibcClientConsensusStateMsg.ClientId,
		RevisionNumber: height.GetRevisionNumber(),
		RevisionHeight: height.GetRevisionHeight(),
		LatestHeight:   ibcClientConsensusStateMsg.LatestHeight,
	}, nil
}

// Parsing - cosmos client for IBC client
func parseCmclientForIbcClientArgs(rpcUrl string) (cmclient.Context, error) {
	if rpcUrl == "" {
		return cmclient.Context{}, util.LogErr(errors.ErrInvalidRequest, "need a tendermint RPC URL")
	}

	client, err := cmclient.NewClientFromNode(rpcUrl)
	if err != nil {
		return cmclient.Context{}, util.LogErr(errors.ErrSdkClient, err)
	}

	clientCtx := cmclient.Context{}
	clientCtx = clientCtx.WithClient(client)

	return clientCtx, nil
}
