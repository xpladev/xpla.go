package ibc

import (
	"github.com/xpladev/xpla.go/types"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	ibcclient "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
)

// Parsing - IBC client consensus state
func parseIbcClientConsensusStateArgs(ibcClientConsensusStateMsg types.IbcClientConsensusStateMsg) (ibcclient.QueryConsensusStateRequest, error) {
	var height ibcclient.Height
	var err error

	if !ibcClientConsensusStateMsg.LatestHeight {
		if ibcClientConsensusStateMsg.Height == "" {
			return ibcclient.QueryConsensusStateRequest{}, types.ErrWrap(types.ErrInsufficientParams, "must include a second 'Height' argument when 'LatestHeight' is not provided")
		}

		height, err = ibcclient.ParseHeight(ibcClientConsensusStateMsg.Height)
		if err != nil {
			return ibcclient.QueryConsensusStateRequest{}, types.ErrWrap(types.ErrParse, err)
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
		return cmclient.Context{}, types.ErrWrap(types.ErrInsufficientParams, "need a tendermint RPC URL")
	}

	client, err := cmclient.NewClientFromNode(rpcUrl)
	if err != nil {
		return cmclient.Context{}, types.ErrWrap(types.ErrSdkClient, err)
	}

	clientCtx := cmclient.Context{}
	clientCtx = clientCtx.WithClient(client)

	return clientCtx, nil
}
