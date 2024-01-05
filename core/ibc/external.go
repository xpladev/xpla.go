package ibc

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &IbcExternal{}

type IbcExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e IbcExternal) {
	e.Xplac = xplac
	e.Name = IbcModule
	return e
}

func (e IbcExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e IbcExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Query

// Query IBC light client states
func (e IbcExternal) IbcClientStates() provider.XplaClient {
	msg, err := MakeIbcClientStatesMsg()
	if err != nil {
		return e.Err(IbcClientStatesMsgType, err)
	}

	return e.ToExternal(IbcClientStatesMsgType, msg)
}

// Query IBC light client state by client ID
func (e IbcExternal) IbcClientState(ibcClientStateMsg types.IbcClientStateMsg) provider.XplaClient {
	msg, err := MakeIbcClientStateMsg(ibcClientStateMsg)
	if err != nil {
		return e.Err(IbcClientStateMsgType, err)
	}

	return e.ToExternal(IbcClientStateMsgType, msg)
}

// Query IBC light client status by client ID
func (e IbcExternal) IbcClientStatus(ibcClientStatusMsg types.IbcClientStatusMsg) provider.XplaClient {
	msg, err := MakeIbcClientStatusMsg(ibcClientStatusMsg)
	if err != nil {
		return e.Err(IbcClientStatusMsgType, err)
	}

	return e.ToExternal(IbcClientStatusMsgType, msg)
}

// Query IBC client consensus states
func (e IbcExternal) IbcClientConsensusStates(ibcClientConsensusStatesMsg types.IbcClientConsensusStatesMsg) provider.XplaClient {
	msg, err := MakeIbcClientConsensusStatesMsg(ibcClientConsensusStatesMsg)
	if err != nil {
		return e.Err(IbcClientConsensusStatesMsgType, err)
	}

	return e.ToExternal(IbcClientConsensusStatesMsgType, msg)
}

// Query IBC client consensus state heights
func (e IbcExternal) IbcClientConsensusStateHeights(ibcClientConsensusStateHeightsMsg types.IbcClientConsensusStateHeightsMsg) provider.XplaClient {
	msg, err := MakeIbcClientConsensusStateHeightsMsg(ibcClientConsensusStateHeightsMsg)
	if err != nil {
		return e.Err(IbcClientConsensusStateHeightsMsgType, err)
	}

	return e.ToExternal(IbcClientConsensusStateHeightsMsgType, msg)
}

// Query IBC client consensus state
func (e IbcExternal) IbcClientConsensusState(ibcClientConsensusStateMsg types.IbcClientConsensusStateMsg) provider.XplaClient {
	msg, err := MakeIbcClientConsensusStateMsg(ibcClientConsensusStateMsg)
	if err != nil {
		return e.Err(IbcClientConsensusStateMsgType, err)
	}

	return e.ToExternal(IbcClientConsensusStateMsgType, msg)
}

// Query IBC client tendermint header
func (e IbcExternal) IbcClientHeader() provider.XplaClient {
	msg, err := MakeIbcClientHeaderMsg(e.Xplac.GetRpc())
	if err != nil {
		return e.Err(IbcClientHeaderMsgType, err)
	}

	return e.ToExternal(IbcClientHeaderMsgType, msg)
}

// Query IBC client self consensus state
func (e IbcExternal) IbcClientSelfConsensusState() provider.XplaClient {
	msg, err := MakeIbcClientSelfConsensusStateMsg(e.Xplac.GetRpc())
	if err != nil {
		return e.Err(IbcClientSelfConsensusStateMsgType, err)
	}

	return e.ToExternal(IbcClientSelfConsensusStateMsgType, msg)
}

// Query IBC client params
func (e IbcExternal) IbcClientParams() provider.XplaClient {
	msg, err := MakeIbcClientParamsMsg()
	if err != nil {
		return e.Err(IbcClientParamsMsgType, err)
	}

	return e.ToExternal(IbcClientParamsMsgType, msg)
}

// Query IBC connections
func (e IbcExternal) IbcConnections(ibcConnectionMsg ...types.IbcConnectionMsg) provider.XplaClient {
	switch {
	case len(ibcConnectionMsg) == 0:
		msg, err := MakeIbcConnectionConnectionsMsg()
		if err != nil {
			return e.Err(IbcConnectionConnectionsMsgType, err)
		}

		return e.ToExternal(IbcConnectionConnectionsMsgType, msg)

	case len(ibcConnectionMsg) == 1:
		msg, err := MakeIbcConnectionConnectionMsg(ibcConnectionMsg[0])
		if err != nil {
			return e.Err(IbcConnectionConnectionMsgType, err)
		}

		return e.ToExternal(IbcConnectionConnectionMsgType, msg)

	default:
		return e.Err(IbcConnectionConnectionMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}

// Query IBC client connections
func (e IbcExternal) IbcClientConnections(ibcClientConnectionsMsg types.IbcClientConnectionsMsg) provider.XplaClient {
	msg, err := MakeIbcConnectionClientConnectionsMsg(ibcClientConnectionsMsg)
	if err != nil {
		return e.Err(IbcConnectionClientConnectionsMsgType, err)
	}

	return e.ToExternal(IbcConnectionClientConnectionsMsgType, msg)
}

// Query IBC channels
func (e IbcExternal) IbcChannels(ibcChannelMsg ...types.IbcChannelMsg) provider.XplaClient {
	switch {
	case len(ibcChannelMsg) == 0:
		msg, err := MakeIbcChannelChannelsMsg()
		if err != nil {
			return e.Err(IbcChannelChannelsMsgType, err)
		}

		return e.ToExternal(IbcChannelChannelsMsgType, msg)

	case len(ibcChannelMsg) == 1:
		msg, err := MakeIbcChannelChannelMsg(ibcChannelMsg[0])
		if err != nil {
			return e.Err(IbcChannelChannelMsgType, err)
		}

		return e.ToExternal(IbcChannelChannelMsgType, msg)

	default:
		return e.Err(IbcChannelChannelMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}

// Query IBC channel connections
func (e IbcExternal) IbcChannelConnections(ibcChannelConnectionsMsg types.IbcChannelConnectionsMsg) provider.XplaClient {
	msg, err := MakeIbcChannelConnectionsMsg(ibcChannelConnectionsMsg)
	if err != nil {
		return e.Err(IbcChannelConnectionsMsgType, err)
	}

	return e.ToExternal(IbcChannelConnectionsMsgType, msg)
}

// Query IBC channel client state
func (e IbcExternal) IbcChannelClientState(ibcChannelClientStateMsg types.IbcChannelClientStateMsg) provider.XplaClient {
	msg, err := MakeIbcChannelClientStateMsg(ibcChannelClientStateMsg)
	if err != nil {
		return e.Err(IbcChannelClientStateMsgType, err)
	}

	return e.ToExternal(IbcChannelClientStateMsgType, msg)
}

// Query IBC channel packet commitments
func (e IbcExternal) IbcChannelPacketCommitments(ibcChannelPacketCommitmentsMsg types.IbcChannelPacketCommitmentsMsg) provider.XplaClient {
	switch {
	case ibcChannelPacketCommitmentsMsg.Sequence == "":
		msg, err := MakeIbcChannelPacketCommitmentsMsg(ibcChannelPacketCommitmentsMsg)
		if err != nil {
			return e.Err(IbcChannelPacketCommitmentsMsgType, err)
		}

		return e.ToExternal(IbcChannelPacketCommitmentsMsgType, msg)

	default:
		msg, err := MakeIbcChannelPacketCommitmentMsg(ibcChannelPacketCommitmentsMsg)
		if err != nil {
			return e.Err(IbcChannelPacketCommitmentMsgType, err)
		}

		return e.ToExternal(IbcChannelPacketCommitmentMsgType, msg)
	}
}

// Query IBC packet receipt
func (e IbcExternal) IbcChannelPacketReceipt(ibcChannelPacketReceiptMsg types.IbcChannelPacketReceiptMsg) provider.XplaClient {
	msg, err := MakeIbcChannelPacketReceiptMsg(ibcChannelPacketReceiptMsg)
	if err != nil {
		return e.Err(IbcChannelPacketReceiptMsgType, err)
	}

	return e.ToExternal(IbcChannelPacketReceiptMsgType, msg)
}

// Query IBC packet ack
func (e IbcExternal) IbcChannelPacketAck(ibcChannelPacketAckMsg types.IbcChannelPacketAckMsg) provider.XplaClient {
	msg, err := MakeIbcChannelPacketAckMsg(ibcChannelPacketAckMsg)
	if err != nil {
		return e.Err(IbcChannelPacketAckMsgType, err)
	}

	return e.ToExternal(IbcChannelPacketAckMsgType, msg)
}

// Query IBC unreceived packets
func (e IbcExternal) IbcChannelUnreceivedPackets(ibcChannelUnreceivedPacketsMsg types.IbcChannelUnreceivedPacketsMsg) provider.XplaClient {
	msg, err := MakeIbcChannelPacketUnreceivedPacketsMsg(ibcChannelUnreceivedPacketsMsg)
	if err != nil {
		return e.Err(IbcChannelUnreceivedPacketsMsgType, err)
	}

	return e.ToExternal(IbcChannelUnreceivedPacketsMsgType, msg)
}

// Query IBC unreceived acks
func (e IbcExternal) IbcChannelUnreceivedAcks(ibcChannelUnreceivedAcksMsg types.IbcChannelUnreceivedAcksMsg) provider.XplaClient {
	msg, err := MakeIbcChannelPacketUnreceivedAcksMsg(ibcChannelUnreceivedAcksMsg)
	if err != nil {
		return e.Err(IbcChannelUnreceivedAcksMsgType, err)
	}

	return e.ToExternal(IbcChannelUnreceivedAcksMsgType, msg)
}

// Query IBC next sequence receive
func (e IbcExternal) IbcChannelNextSequence(ibcChannelNextSequenceMsg types.IbcChannelNextSequenceMsg) provider.XplaClient {
	msg, err := MakeIbcChannelNextSequenceReceiveMsg(ibcChannelNextSequenceMsg)
	if err != nil {
		return e.Err(IbcChannelNextSequenceMsgType, err)
	}

	return e.ToExternal(IbcChannelNextSequenceMsgType, msg)
}

// Query IBC transfer denom traces
func (e IbcExternal) IbcDenomTraces(ibcDenomTraceMsg ...types.IbcDenomTraceMsg) provider.XplaClient {
	switch {
	case len(ibcDenomTraceMsg) == 0:
		msg, err := MakeIbcTransferDenomTracesMsg()
		if err != nil {
			return e.Err(IbcTransferDenomTracesMsgType, err)
		}

		return e.ToExternal(IbcTransferDenomTracesMsgType, msg)

	case len(ibcDenomTraceMsg) == 1:
		msg, err := MakeIbcTransferDenomTraceMsg(ibcDenomTraceMsg[0])
		if err != nil {
			return e.Err(IbcTransferDenomTraceMsgType, err)
		}

		return e.ToExternal(IbcTransferDenomTraceMsgType, msg)

	default:
		return e.Err(IbcTransferDenomTraceMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}

// Query IBC transfer denom trace
func (e IbcExternal) IbcDenomTrace(ibcDenomTraceMsg types.IbcDenomTraceMsg) provider.XplaClient {
	msg, err := MakeIbcTransferDenomTraceMsg(ibcDenomTraceMsg)
	if err != nil {
		return e.Err(IbcTransferDenomTraceMsgType, err)
	}

	return e.ToExternal(IbcTransferDenomTraceMsgType, msg)
}

// Query IBC transfer denom hash
func (e IbcExternal) IbcDenomHash(ibcDenomHashMsg types.IbcDenomHashMsg) provider.XplaClient {
	msg, err := MakeIbcTransferDenomHashMsg(ibcDenomHashMsg)
	if err != nil {
		return e.Err(IbcTransferDenomHashMsgType, err)
	}

	return e.ToExternal(IbcTransferDenomHashMsgType, msg)
}

// Query IBC transfer denom hash
func (e IbcExternal) IbcEscrowAddress(ibcEscrowAddressMsg types.IbcEscrowAddressMsg) provider.XplaClient {
	msg, err := MakeIbcTransferEscrowAddressMsg(ibcEscrowAddressMsg)
	if err != nil {
		return e.Err(IbcTransferEscrowAddressMsgType, err)
	}

	return e.ToExternal(IbcTransferEscrowAddressMsgType, msg)
}

// Query IBC transfer params
func (e IbcExternal) IbcTransferParams() provider.XplaClient {
	msg, err := MakeIbcTransferParamsMsg()
	if err != nil {
		return e.Err(IbcTransferParamsMsgType, err)
	}

	return e.ToExternal(IbcTransferParamsMsgType, msg)
}
