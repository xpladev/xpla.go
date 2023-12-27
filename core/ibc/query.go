package ibc

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	cmclient "github.com/cosmos/cosmos-sdk/client"
	ibctransfer "github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"
	ibcclientutils "github.com/cosmos/ibc-go/v4/modules/core/02-client/client/utils"
	ibcclient "github.com/cosmos/ibc-go/v4/modules/core/02-client/types"
	ibcconnection "github.com/cosmos/ibc-go/v4/modules/core/03-connection/types"
	ibcchannel "github.com/cosmos/ibc-go/v4/modules/core/04-channel/types"
)

var out []byte
var res proto.Message
var err error

// Query client for gov module.
func QueryIbc(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcIbc(i)
	} else {
		return queryByLcdIbc(i)
	}
}

func queryByGrpcIbc(i core.QueryClient) (string, error) {
	ibcclientQueryClient := ibcclient.NewQueryClient(i.Ixplac.GetGrpcClient())
	ibcconnectionQueryClient := ibcconnection.NewQueryClient(i.Ixplac.GetGrpcClient())
	ibccchannelQueryClient := ibcchannel.NewQueryClient(i.Ixplac.GetGrpcClient())
	ibctransferQueryClient := ibctransfer.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// IBC client states
	case i.Ixplac.GetMsgType() == IbcClientStatesMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryClientStatesRequest)
		res, err = ibcclientQueryClient.ClientStates(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC client state
	case i.Ixplac.GetMsgType() == IbcClientStateMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryClientStateRequest)
		res, err = ibcclientQueryClient.ClientState(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC client status
	case i.Ixplac.GetMsgType() == IbcClientStatusMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryClientStatusRequest)
		res, err = ibcclientQueryClient.ClientStatus(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC client consensus states
	case i.Ixplac.GetMsgType() == IbcClientConsensusStatesMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryConsensusStatesRequest)
		res, err = ibcclientQueryClient.ConsensusStates(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC client consensus state heights
	case i.Ixplac.GetMsgType() == IbcClientConsensusStateHeightsMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryConsensusStateHeightsRequest)
		res, err = ibcclientQueryClient.ConsensusStateHeights(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC client consensus state
	case i.Ixplac.GetMsgType() == IbcClientConsensusStateMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryConsensusStateRequest)
		res, err = ibcclientQueryClient.ConsensusState(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC client tendermint header
	case i.Ixplac.GetMsgType() == IbcClientHeaderMsgType:

		convertMsg := i.Ixplac.GetMsg().(cmclient.Context)
		header, _, err := ibcclientutils.QueryTendermintHeader(convertMsg)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		res = &header

	// IBC client self consensus state
	case i.Ixplac.GetMsgType() == IbcClientSelfConsensusStateMsgType:

		convertMsg := i.Ixplac.GetMsg().(cmclient.Context)
		state, _, err := ibcclientutils.QuerySelfConsensusState(convertMsg)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

		res = state

	// IBC client params
	case i.Ixplac.GetMsgType() == IbcClientParamsMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryClientParamsRequest)
		res, err = ibcclientQueryClient.ClientParams(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC connection connections
	case i.Ixplac.GetMsgType() == IbcConnectionConnectionsMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcconnection.QueryConnectionsRequest)
		res, err = ibcconnectionQueryClient.Connections(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC connection connection
	case i.Ixplac.GetMsgType() == IbcConnectionConnectionMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcconnection.QueryConnectionRequest)
		res, err = ibcconnectionQueryClient.Connection(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC connection a client connections
	case i.Ixplac.GetMsgType() == IbcConnectionClientConnectionsMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcconnection.QueryClientConnectionsRequest)
		res, err = ibcconnectionQueryClient.ClientConnections(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC channels
	case i.Ixplac.GetMsgType() == IbcChannelChannelsMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryChannelsRequest)
		res, err = ibccchannelQueryClient.Channels(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC a channel
	case i.Ixplac.GetMsgType() == IbcChannelChannelMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryChannelRequest)
		res, err = ibccchannelQueryClient.Channel(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC channel connections
	case i.Ixplac.GetMsgType() == IbcChannelConnectionsMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryConnectionChannelsRequest)
		res, err = ibccchannelQueryClient.ConnectionChannels(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC channel client state
	case i.Ixplac.GetMsgType() == IbcChannelClientStateMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryChannelClientStateRequest)
		res, err = ibccchannelQueryClient.ChannelClientState(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC channel packet commitments
	case i.Ixplac.GetMsgType() == IbcChannelPacketCommitmentsMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryPacketCommitmentsRequest)
		res, err = ibccchannelQueryClient.PacketCommitments(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC channel packet commitment by sequece
	case i.Ixplac.GetMsgType() == IbcChannelPacketCommitmentMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryPacketCommitmentRequest)
		res, err = ibccchannelQueryClient.PacketCommitment(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC channel packet receipt
	case i.Ixplac.GetMsgType() == IbcChannelPacketReceiptMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryPacketReceiptRequest)
		res, err = ibccchannelQueryClient.PacketReceipt(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC channel packet ack
	case i.Ixplac.GetMsgType() == IbcChannelPacketAckMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryPacketAcknowledgementRequest)
		res, err = ibccchannelQueryClient.PacketAcknowledgement(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC channel unreceived packets
	case i.Ixplac.GetMsgType() == IbcChannelUnreceivedPacketsMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryUnreceivedPacketsRequest)
		res, err = ibccchannelQueryClient.UnreceivedPackets(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC channel unreceived acks
	case i.Ixplac.GetMsgType() == IbcChannelUnreceivedAcksMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryUnreceivedAcksRequest)
		res, err = ibccchannelQueryClient.UnreceivedAcks(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC channel next sequence receive
	case i.Ixplac.GetMsgType() == IbcChannelNextSequenceMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryNextSequenceReceiveRequest)
		res, err = ibccchannelQueryClient.NextSequenceReceive(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC transfer denom traces
	case i.Ixplac.GetMsgType() == IbcTransferDenomTracesMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibctransfer.QueryDenomTracesRequest)
		res, err = ibctransferQueryClient.DenomTraces(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC transfer denom trace
	case i.Ixplac.GetMsgType() == IbcTransferDenomTraceMsgType:

		convertMsg := i.Ixplac.GetMsg().(ibctransfer.QueryDenomTraceRequest)
		res, err = ibctransferQueryClient.DenomTrace(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC transfer denom hash
	case i.Ixplac.GetMsgType() == IbcTransferDenomHashMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibctransfer.QueryDenomHashRequest)
		res, err = ibctransferQueryClient.DenomHash(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// IBC transfer escrow address
	case i.Ixplac.GetMsgType() == IbcTransferEscrowAddressMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.IbcEscrowAddressMsg)

		addr := ibctransfer.GetEscrowAddress(convertMsg.PortId, convertMsg.ChannelId)
		return addr.String(), nil

	// IBC transfer params
	case i.Ixplac.GetMsgType() == IbcTransferParamsMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibctransfer.QueryParamsRequest)
		res, err = ibctransferQueryClient.Params(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = core.PrintProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	ibcclientClientStatesLabel          = "client_states"
	ibcclientClientStatusLabel          = "client_status"
	ibcclientClientConsensusStatesLabel = "consensus_states"
	ibcclientHeightsLabel               = "heights"
	ibcclientRevisionNumberLabel        = "revision"
	ibcclientRevisionHeightLabel        = "height"

	ibcconnectionConnectionsLabel       = "connections"
	ibcconnectionClientConnectionsLabel = "client_connections"

	ibcchannelChannelsLabel          = "channels"
	ibcchannelPortsLabel             = "ports"
	ibcchannelClientStateLabel       = "client_state"
	ibcchannelPacketCommitmentsLabel = "packet_commitments"
	ibcchannelUnreceivedPacketsLabel = "unreceived_packets"
	ibcchannelUnreceivedAcksLabel    = "unreceived_acks"
	ibcchannelPacketReceiptLabel     = "packet_receipts"
	ibcchannelPacketAckLabel         = "packet_acks"
	ibcchannelNextSequenceLabel      = "next_sequence"

	ibctransferDenomTracesLabel   = "denom_traces"
	ibctransferDenomHashesLabel   = "denom_hashes"
	ibctransferEscrowAddressLabel = "escrow_address"
)

func queryByLcdIbc(i core.QueryClient) (string, error) {
	var url string
	ibcclientUrl := "/ibc/core/client/v1/"
	ibcconnectionUrl := "/ibc/core/connection/v1/"
	ibcchannelUrl := "/ibc/core/channel/v1/"
	ibctransferUrl := "/ibc/apps/transfer/v1/"

	switch {
	// IBC client states
	case i.Ixplac.GetMsgType() == IbcClientStatesMsgType:
		url = ibcclientUrl + ibcclientClientStatesLabel

	// IBC client state
	case i.Ixplac.GetMsgType() == IbcClientStateMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryClientStateRequest)

		url = ibcclientUrl + util.MakeQueryLabels(ibcclientClientStatesLabel, convertMsg.ClientId)

	// IBC client status
	case i.Ixplac.GetMsgType() == IbcClientStatusMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryClientStatusRequest)

		url = ibcclientUrl + util.MakeQueryLabels(ibcclientClientStatusLabel, convertMsg.ClientId)

	// IBC client consensus states
	case i.Ixplac.GetMsgType() == IbcClientConsensusStatesMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryConsensusStatesRequest)

		url = ibcclientUrl + util.MakeQueryLabels(ibcclientClientConsensusStatesLabel, convertMsg.ClientId)

	// IBC client consensus state heights
	case i.Ixplac.GetMsgType() == IbcClientConsensusStateHeightsMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryConsensusStateHeightsRequest)

		url = ibcclientUrl + util.MakeQueryLabels(ibcclientClientConsensusStatesLabel, convertMsg.ClientId, ibcclientHeightsLabel)

	// IBC client consensus state height
	case i.Ixplac.GetMsgType() == IbcClientConsensusStateMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcclient.QueryConsensusStateRequest)

		url = ibcclientUrl + util.MakeQueryLabels(
			ibcclientClientConsensusStatesLabel,
			convertMsg.ClientId,
			ibcclientRevisionNumberLabel,
			util.FromUint64ToString(convertMsg.RevisionNumber),
			ibcclientRevisionHeightLabel,
			util.FromUint64ToString(convertMsg.RevisionHeight),
		)

	// IBC client tendermint header
	case i.Ixplac.GetMsgType() == IbcClientHeaderMsgType:
		return "", util.LogErr(errors.ErrNotSupport, "unsupported querying IBC client tendermint header by using LCD")

	// IBC client self consensus state
	case i.Ixplac.GetMsgType() == IbcClientSelfConsensusStateMsgType:
		return "", util.LogErr(errors.ErrNotSupport, "unsupported querying IBC client self consensus state by using LCD")

	// IBC client params
	case i.Ixplac.GetMsgType() == IbcClientParamsMsgType:
		url = "/ibc/client/v1/params"

	// IBC connection connections
	case i.Ixplac.GetMsgType() == IbcConnectionConnectionsMsgType:
		url = ibcconnectionUrl + ibcconnectionConnectionsLabel

	// IBC connection connection
	case i.Ixplac.GetMsgType() == IbcConnectionConnectionMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcconnection.QueryConnectionRequest)

		url = ibcconnectionUrl + util.MakeQueryLabels(ibcconnectionConnectionsLabel, convertMsg.ConnectionId)

	// IBC connection a client connections
	case i.Ixplac.GetMsgType() == IbcConnectionClientConnectionsMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcconnection.QueryClientConnectionsRequest)

		url = ibcconnectionUrl + util.MakeQueryLabels(ibcconnectionClientConnectionsLabel, convertMsg.ClientId)

	// IBC channels
	case i.Ixplac.GetMsgType() == IbcChannelChannelsMsgType:
		url = ibcchannelUrl + ibcchannelChannelsLabel

	// IBC a channel
	case i.Ixplac.GetMsgType() == IbcChannelChannelMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryChannelRequest)

		url = ibcchannelUrl + util.MakeQueryLabels(ibcchannelChannelsLabel, convertMsg.ChannelId, ibcchannelPortsLabel, convertMsg.PortId)

	// IBC channel connections
	case i.Ixplac.GetMsgType() == IbcChannelConnectionsMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryConnectionChannelsRequest)

		url = ibcchannelUrl + util.MakeQueryLabels(ibcconnectionConnectionsLabel, convertMsg.Connection, ibcchannelChannelsLabel)

	// IBC channel client state
	case i.Ixplac.GetMsgType() == IbcChannelClientStateMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryChannelClientStateRequest)

		url = ibcchannelUrl + util.MakeQueryLabels(ibcchannelChannelsLabel, convertMsg.ChannelId, ibcchannelPortsLabel, convertMsg.PortId, ibcchannelClientStateLabel)

	// IBC channel packet commitments
	case i.Ixplac.GetMsgType() == IbcChannelPacketCommitmentsMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryPacketCommitmentsRequest)

		url = ibcchannelUrl + util.MakeQueryLabels(ibcchannelChannelsLabel, convertMsg.ChannelId, ibcchannelPortsLabel, convertMsg.PortId, ibcchannelPacketCommitmentsLabel)

	// IBC channel packet commitment by sequece
	case i.Ixplac.GetMsgType() == IbcChannelPacketCommitmentMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryPacketCommitmentRequest)

		url = ibcchannelUrl + util.MakeQueryLabels(ibcchannelChannelsLabel, convertMsg.ChannelId, ibcchannelPortsLabel, convertMsg.PortId, ibcchannelPacketCommitmentsLabel, util.FromUint64ToString(convertMsg.Sequence))

	// IBC channel packet receipt
	case i.Ixplac.GetMsgType() == IbcChannelPacketReceiptMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryPacketReceiptRequest)

		url = ibcchannelUrl + util.MakeQueryLabels(ibcchannelChannelsLabel, convertMsg.ChannelId, ibcchannelPortsLabel, convertMsg.PortId, ibcchannelPacketReceiptLabel, util.FromUint64ToString(convertMsg.Sequence))

	// IBC channel packet ack
	case i.Ixplac.GetMsgType() == IbcChannelPacketAckMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryPacketAcknowledgementRequest)

		url = ibcchannelUrl + util.MakeQueryLabels(ibcchannelChannelsLabel, convertMsg.ChannelId, ibcchannelPortsLabel, convertMsg.PortId, ibcchannelPacketAckLabel, util.FromUint64ToString(convertMsg.Sequence))

	// IBC channel unreceived packets
	case i.Ixplac.GetMsgType() == IbcChannelUnreceivedPacketsMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryUnreceivedPacketsRequest)

		url = ibcchannelUrl + util.MakeQueryLabels(
			ibcchannelChannelsLabel,
			convertMsg.ChannelId,
			ibcchannelPortsLabel,
			convertMsg.PortId,
			ibcchannelPacketCommitmentsLabel,
			util.FromUint64ToString(convertMsg.PacketCommitmentSequences[0]),
			ibcchannelUnreceivedPacketsLabel,
		)

	// IBC channel unreceived acks
	case i.Ixplac.GetMsgType() == IbcChannelUnreceivedAcksMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryUnreceivedAcksRequest)

		url = ibcchannelUrl + util.MakeQueryLabels(
			ibcchannelChannelsLabel,
			convertMsg.ChannelId,
			ibcchannelPortsLabel,
			convertMsg.PortId,
			ibcchannelPacketCommitmentsLabel,
			util.FromUint64ToString(convertMsg.PacketAckSequences[0]),
			ibcchannelUnreceivedAcksLabel,
		)

	// IBC channel next sequence receive
	case i.Ixplac.GetMsgType() == IbcChannelNextSequenceMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibcchannel.QueryNextSequenceReceiveRequest)

		url = ibcchannelUrl + util.MakeQueryLabels(ibcchannelChannelsLabel, convertMsg.ChannelId, ibcchannelPortsLabel, convertMsg.PortId, ibcchannelNextSequenceLabel)

	// IBC transfer denom traces
	case i.Ixplac.GetMsgType() == IbcTransferDenomTracesMsgType:
		url = ibctransferUrl + ibctransferDenomTracesLabel

	// IBC transfer denom trace
	case i.Ixplac.GetMsgType() == IbcTransferDenomTraceMsgType:
		convertMsg := i.Ixplac.GetMsg().(ibctransfer.QueryDenomTraceRequest)

		url = ibctransferUrl + util.MakeQueryLabels(ibctransferDenomTracesLabel, convertMsg.Hash)

	// IBC transfer denom hash
	case i.Ixplac.GetMsgType() == IbcTransferDenomHashMsgType:
		return "", util.LogErr(errors.ErrNotSupport, "unsupported querying denom hash by using LCD")

	// IBC transfer escrow address
	case i.Ixplac.GetMsgType() == IbcTransferEscrowAddressMsgType:
		convertMsg := i.Ixplac.GetMsg().(types.IbcEscrowAddressMsg)

		url = ibctransferUrl + util.MakeQueryLabels(ibcchannelChannelsLabel, convertMsg.ChannelId, ibcchannelPortsLabel, convertMsg.PortId, ibctransferEscrowAddressLabel)

	// IBC transfer params
	case i.Ixplac.GetMsgType() == IbcTransferParamsMsgType:
		url = ibctransferUrl + "/params"

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	i.Ixplac.GetHttpMutex().Lock()
	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		i.Ixplac.GetHttpMutex().Unlock()
		return "", err
	}
	i.Ixplac.GetHttpMutex().Unlock()

	return string(out), nil

}
