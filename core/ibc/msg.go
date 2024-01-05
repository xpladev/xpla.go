package ibc

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	ibctransfer "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibcclient "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	ibcconnection "github.com/cosmos/ibc-go/v4/modules/core/03-connection/types"
	ibcchannel "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
)

// (Query) make msg - IBC client states
func MakeIbcClientStatesMsg() (ibcclient.QueryClientStatesRequest, error) {
	return ibcclient.QueryClientStatesRequest{
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - IBC client state by client ID
func MakeIbcClientStateMsg(ibcClientStatesMsg types.IbcClientStateMsg) (ibcclient.QueryClientStateRequest, error) {
	return ibcclient.QueryClientStateRequest{
		ClientId: ibcClientStatesMsg.ClientId,
	}, nil
}

// (Query) make msg - IBC client status by client ID
func MakeIbcClientStatusMsg(ibcClientStatusMsg types.IbcClientStatusMsg) (ibcclient.QueryClientStatusRequest, error) {
	return ibcclient.QueryClientStatusRequest{
		ClientId: ibcClientStatusMsg.ClientId,
	}, nil
}

// (Query) make msg - IBC client consensus states
func MakeIbcClientConsensusStatesMsg(ibcClientConsensusStatesMsg types.IbcClientConsensusStatesMsg) (ibcclient.QueryConsensusStatesRequest, error) {
	return ibcclient.QueryConsensusStatesRequest{
		ClientId:   ibcClientConsensusStatesMsg.ClientId,
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - IBC client consensus state heights
func MakeIbcClientConsensusStateHeightsMsg(ibcClientConsensusStateHeightsMsg types.IbcClientConsensusStateHeightsMsg) (ibcclient.QueryConsensusStateHeightsRequest, error) {
	return ibcclient.QueryConsensusStateHeightsRequest{
		ClientId:   ibcClientConsensusStateHeightsMsg.ClientId,
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - IBC client consensus state
func MakeIbcClientConsensusStateMsg(ibcClientConsensusStateMsg types.IbcClientConsensusStateMsg) (ibcclient.QueryConsensusStateRequest, error) {
	return parseIbcClientConsensusStateArgs(ibcClientConsensusStateMsg)
}

// (Query) make msg - IBC client tendermint header
func MakeIbcClientHeaderMsg(rpcUrl string) (cmclient.Context, error) {
	return parseCmclientForIbcClientArgs(rpcUrl)
}

// (Query) make msg - IBC client self consensus state
func MakeIbcClientSelfConsensusStateMsg(rpcUrl string) (cmclient.Context, error) {
	return parseCmclientForIbcClientArgs(rpcUrl)
}

// (Query) make msg - IBC client params
func MakeIbcClientParamsMsg() (ibcclient.QueryClientParamsRequest, error) {
	return ibcclient.QueryClientParamsRequest{}, nil
}

// (Query) make msg - IBC connection connetions
func MakeIbcConnectionConnectionsMsg() (ibcconnection.QueryConnectionsRequest, error) {
	return ibcconnection.QueryConnectionsRequest{
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - IBC connection connection
func MakeIbcConnectionConnectionMsg(ibcConnectionMsg types.IbcConnectionMsg) (ibcconnection.QueryConnectionRequest, error) {
	return ibcconnection.QueryConnectionRequest{
		ConnectionId: ibcConnectionMsg.ConnectionId,
	}, nil
}

// (Query) make msg - IBC client connections
func MakeIbcConnectionClientConnectionsMsg(ibcClientConnectionsMsg types.IbcClientConnectionsMsg) (ibcconnection.QueryClientConnectionsRequest, error) {
	return ibcconnection.QueryClientConnectionsRequest{
		ClientId: ibcClientConnectionsMsg.ClientId,
	}, nil
}

// (Query) make msg - IBC channels
func MakeIbcChannelChannelsMsg() (ibcchannel.QueryChannelsRequest, error) {
	return ibcchannel.QueryChannelsRequest{
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - IBC a channel
func MakeIbcChannelChannelMsg(ibcChannelMsg types.IbcChannelMsg) (ibcchannel.QueryChannelRequest, error) {
	return ibcchannel.QueryChannelRequest{
		PortId:    ibcChannelMsg.PortId,
		ChannelId: ibcChannelMsg.ChannelId,
	}, nil
}

// (Query) make msg - IBC channel connections
func MakeIbcChannelConnectionsMsg(ibcChannelConnectionsMsg types.IbcChannelConnectionsMsg) (ibcchannel.QueryConnectionChannelsRequest, error) {
	return ibcchannel.QueryConnectionChannelsRequest{
		Connection: ibcChannelConnectionsMsg.ConnectionId,
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - IBC channel connections
func MakeIbcChannelClientStateMsg(ibcChannelClientStateMsg types.IbcChannelClientStateMsg) (ibcchannel.QueryChannelClientStateRequest, error) {
	return ibcchannel.QueryChannelClientStateRequest{
		ChannelId: ibcChannelClientStateMsg.ChannelId,
		PortId:    ibcChannelClientStateMsg.PortId,
	}, nil
}

// (Query) make msg - IBC channel packet commitments
func MakeIbcChannelPacketCommitmentsMsg(ibcChannelPacketCommitmentsMsg types.IbcChannelPacketCommitmentsMsg) (ibcchannel.QueryPacketCommitmentsRequest, error) {
	return ibcchannel.QueryPacketCommitmentsRequest{
		Pagination: core.PageRequest,
		ChannelId:  ibcChannelPacketCommitmentsMsg.ChannelId,
		PortId:     ibcChannelPacketCommitmentsMsg.PortId,
	}, nil
}

// (Query) make msg - IBC channel packet commitment
func MakeIbcChannelPacketCommitmentMsg(ibcChannelPacketCommitmentsMsg types.IbcChannelPacketCommitmentsMsg) (ibcchannel.QueryPacketCommitmentRequest, error) {
	seq, err := util.FromStringToUint64(ibcChannelPacketCommitmentsMsg.Sequence)
	if err != nil {
		return ibcchannel.QueryPacketCommitmentRequest{}, types.ErrWrap(types.ErrConvert, err)
	}
	return ibcchannel.QueryPacketCommitmentRequest{
		ChannelId: ibcChannelPacketCommitmentsMsg.ChannelId,
		PortId:    ibcChannelPacketCommitmentsMsg.PortId,
		Sequence:  seq,
	}, nil
}

// (Query) make msg - IBC channel packet receipt
func MakeIbcChannelPacketReceiptMsg(ibcChannelPacketReceiptMsg types.IbcChannelPacketReceiptMsg) (ibcchannel.QueryPacketReceiptRequest, error) {
	seq, err := util.FromStringToUint64(ibcChannelPacketReceiptMsg.Sequence)
	if err != nil {
		return ibcchannel.QueryPacketReceiptRequest{}, types.ErrWrap(types.ErrConvert, err)
	}
	return ibcchannel.QueryPacketReceiptRequest{
		ChannelId: ibcChannelPacketReceiptMsg.ChannelId,
		PortId:    ibcChannelPacketReceiptMsg.PortId,
		Sequence:  seq,
	}, nil
}

// (Query) make msg - IBC channel packet ack
func MakeIbcChannelPacketAckMsg(ibcChannelPacketAckMsg types.IbcChannelPacketAckMsg) (ibcchannel.QueryPacketAcknowledgementRequest, error) {
	seq, err := util.FromStringToUint64(ibcChannelPacketAckMsg.Sequence)
	if err != nil {
		return ibcchannel.QueryPacketAcknowledgementRequest{}, types.ErrWrap(types.ErrConvert, err)
	}
	return ibcchannel.QueryPacketAcknowledgementRequest{
		ChannelId: ibcChannelPacketAckMsg.ChannelId,
		PortId:    ibcChannelPacketAckMsg.PortId,
		Sequence:  seq,
	}, nil
}

// (Query) make msg - IBC channel unreceived packets
func MakeIbcChannelPacketUnreceivedPacketsMsg(ibcChannelUnreceivedPacketsMsg types.IbcChannelUnreceivedPacketsMsg) (ibcchannel.QueryUnreceivedPacketsRequest, error) {
	seq, err := util.FromStringToUint64(ibcChannelUnreceivedPacketsMsg.Sequence)
	if err != nil {
		return ibcchannel.QueryUnreceivedPacketsRequest{}, types.ErrWrap(types.ErrConvert, err)
	}
	return ibcchannel.QueryUnreceivedPacketsRequest{
		ChannelId:                 ibcChannelUnreceivedPacketsMsg.ChannelId,
		PortId:                    ibcChannelUnreceivedPacketsMsg.PortId,
		PacketCommitmentSequences: []uint64{seq},
	}, nil
}

// (Query) make msg - IBC channel unreceived acks
func MakeIbcChannelPacketUnreceivedAcksMsg(ibcChannelUnreceivedAcksMsg types.IbcChannelUnreceivedAcksMsg) (ibcchannel.QueryUnreceivedAcksRequest, error) {
	seq, err := util.FromStringToUint64(ibcChannelUnreceivedAcksMsg.Sequence)
	if err != nil {
		return ibcchannel.QueryUnreceivedAcksRequest{}, types.ErrWrap(types.ErrConvert, err)
	}
	return ibcchannel.QueryUnreceivedAcksRequest{
		ChannelId:          ibcChannelUnreceivedAcksMsg.ChannelId,
		PortId:             ibcChannelUnreceivedAcksMsg.PortId,
		PacketAckSequences: []uint64{seq},
	}, nil
}

// (Query) make msg - IBC channel next sequence receive
func MakeIbcChannelNextSequenceReceiveMsg(ibcChannelNextSequenceMsg types.IbcChannelNextSequenceMsg) (ibcchannel.QueryNextSequenceReceiveRequest, error) {
	return ibcchannel.QueryNextSequenceReceiveRequest{
		ChannelId: ibcChannelNextSequenceMsg.ChannelId,
		PortId:    ibcChannelNextSequenceMsg.PortId,
	}, nil
}

// (Query) make msg - IBC transfer denom traces
func MakeIbcTransferDenomTracesMsg() (ibctransfer.QueryDenomTracesRequest, error) {
	return ibctransfer.QueryDenomTracesRequest{
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - IBC transfer denom trace
func MakeIbcTransferDenomTraceMsg(ibcDenomTraceMsg types.IbcDenomTraceMsg) (ibctransfer.QueryDenomTraceRequest, error) {
	return ibctransfer.QueryDenomTraceRequest{
		Hash: ibcDenomTraceMsg.HashDenom,
	}, nil
}

// (Query) make msg - IBC transfer denom hash
func MakeIbcTransferDenomHashMsg(ibcDenomHashMsg types.IbcDenomHashMsg) (ibctransfer.QueryDenomHashRequest, error) {
	return ibctransfer.QueryDenomHashRequest{
		Trace: ibcDenomHashMsg.Trace,
	}, nil
}

// (Query) make msg - IBC transfer escrow address
func MakeIbcTransferEscrowAddressMsg(ibcEscrowAddressMsg types.IbcEscrowAddressMsg) (types.IbcEscrowAddressMsg, error) {
	return ibcEscrowAddressMsg, nil
}

// (Query) make msg - IBC transfer params
func MakeIbcTransferParamsMsg() (ibctransfer.QueryParamsRequest, error) {
	return ibctransfer.QueryParamsRequest{}, nil
}
