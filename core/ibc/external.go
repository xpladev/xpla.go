package ibc

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type IbcExternal struct {
	Xplac provider.XplaClient
}

func NewIbcExternal(xplac provider.XplaClient) (e IbcExternal) {
	e.Xplac = xplac
	return e
}

// Query

// Query IBC light client states
func (e IbcExternal) IbcClientStates() provider.XplaClient {
	msg, err := MakeIbcClientStatesMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcClientStatesMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC light client state by client ID
func (e IbcExternal) IbcClientState(ibcClientStateMsg types.IbcClientStateMsg) provider.XplaClient {
	msg, err := MakeIbcClientStateMsg(ibcClientStateMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcClientStateMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC light client status by client ID
func (e IbcExternal) IbcClientStatus(ibcClientStatusMsg types.IbcClientStatusMsg) provider.XplaClient {
	msg, err := MakeIbcClientStatusMsg(ibcClientStatusMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcClientStatusMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC client consensus states
func (e IbcExternal) IbcClientConsensusStates(ibcClientConsensusStatesMsg types.IbcClientConsensusStatesMsg) provider.XplaClient {
	msg, err := MakeIbcClientConsensusStatesMsg(ibcClientConsensusStatesMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcClientConsensusStatesMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC client consensus state heights
func (e IbcExternal) IbcClientConsensusStateHeights(ibcClientConsensusStateHeightsMsg types.IbcClientConsensusStateHeightsMsg) provider.XplaClient {
	msg, err := MakeIbcClientConsensusStateHeightsMsg(ibcClientConsensusStateHeightsMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcClientConsensusStateHeightsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC client consensus state
func (e IbcExternal) IbcClientConsensusState(ibcClientConsensusStateMsg types.IbcClientConsensusStateMsg) provider.XplaClient {
	msg, err := MakeIbcClientConsensusStateMsg(ibcClientConsensusStateMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcClientConsensusStateMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC client tendermint header
func (e IbcExternal) IbcClientHeader() provider.XplaClient {
	msg, err := MakeIbcClientHeaderMsg(e.Xplac.GetRpc())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcClientHeaderMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC client self consensus state
func (e IbcExternal) IbcClientSelfConsensusState() provider.XplaClient {
	msg, err := MakeIbcClientSelfConsensusStateMsg(e.Xplac.GetRpc())
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcClientSelfConsensusStateMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC client params
func (e IbcExternal) IbcClientParams() provider.XplaClient {
	msg, err := MakeIbcClientParamsMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcClientParamsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC connections
func (e IbcExternal) IbcConnections(ibcConnectionMsg ...types.IbcConnectionMsg) provider.XplaClient {
	if len(ibcConnectionMsg) == 0 {
		msg, err := MakeIbcConnectionConnectionsMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(IbcModule).
			WithMsgType(IbcConnectionConnectionsMsgType).
			WithMsg(msg)
	} else if len(ibcConnectionMsg) == 1 {
		msg, err := MakeIbcConnectionConnectionMsg(ibcConnectionMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(IbcModule).
			WithMsgType(IbcConnectionConnectionMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}

// Query IBC client connections
func (e IbcExternal) IbcClientConnections(ibcClientConnectionsMsg types.IbcClientConnectionsMsg) provider.XplaClient {
	msg, err := MakeIbcConnectionClientConnectionsMsg(ibcClientConnectionsMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcConnectionClientConnectionsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC channels
func (e IbcExternal) IbcChannels(ibcChannelMsg ...types.IbcChannelMsg) provider.XplaClient {
	if len(ibcChannelMsg) == 0 {
		msg, err := MakeIbcChannelChannelsMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(IbcModule).
			WithMsgType(IbcChannelChannelsMsgType).
			WithMsg(msg)
	} else if len(ibcChannelMsg) == 1 {
		msg, err := MakeIbcChannelChannelMsg(ibcChannelMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(IbcModule).
			WithMsgType(IbcChannelChannelMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}

// Query IBC channel connections
func (e IbcExternal) IbcChannelConnections(ibcChannelConnectionsMsg types.IbcChannelConnectionsMsg) provider.XplaClient {
	msg, err := MakeIbcChannelConnectionsMsg(ibcChannelConnectionsMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcChannelConnectionsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC channel client state
func (e IbcExternal) IbcChannelClientState(ibcChannelClientStateMsg types.IbcChannelClientStateMsg) provider.XplaClient {
	msg, err := MakeIbcChannelClientStateMsg(ibcChannelClientStateMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcChannelClientStateMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC channel packet commitments
func (e IbcExternal) IbcChannelPacketCommitments(ibcChannelPacketCommitmentsMsg types.IbcChannelPacketCommitmentsMsg) provider.XplaClient {
	if ibcChannelPacketCommitmentsMsg.Sequence == "" {
		msg, err := MakeIbcChannelPacketCommitmentsMsg(ibcChannelPacketCommitmentsMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(IbcModule).
			WithMsgType(IbcChannelPacketCommitmentsMsgType).
			WithMsg(msg)
	} else {
		msg, err := MakeIbcChannelPacketCommitmentMsg(ibcChannelPacketCommitmentsMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(IbcModule).
			WithMsgType(IbcChannelPacketCommitmentMsgType).
			WithMsg(msg)
	}
	return e.Xplac
}

// Query IBC packet receipt
func (e IbcExternal) IbcChannelPacketReceipt(ibcChannelPacketReceiptMsg types.IbcChannelPacketReceiptMsg) provider.XplaClient {
	msg, err := MakeIbcChannelPacketReceiptMsg(ibcChannelPacketReceiptMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcChannelPacketReceiptMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC packet ack
func (e IbcExternal) IbcChannelPacketAck(ibcChannelPacketAckMsg types.IbcChannelPacketAckMsg) provider.XplaClient {
	msg, err := MakeIbcChannelPacketAckMsg(ibcChannelPacketAckMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcChannelPacketAckMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC unreceived packets
func (e IbcExternal) IbcChannelUnreceivedPackets(ibcChannelUnreceivedPacketsMsg types.IbcChannelUnreceivedPacketsMsg) provider.XplaClient {
	msg, err := MakeIbcChannelPacketUnreceivedPacketsMsg(ibcChannelUnreceivedPacketsMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcChannelUnreceivedPacketsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC unreceived acks
func (e IbcExternal) IbcChannelUnreceivedAcks(ibcChannelUnreceivedAcksMsg types.IbcChannelUnreceivedAcksMsg) provider.XplaClient {
	msg, err := MakeIbcChannelPacketUnreceivedAcksMsg(ibcChannelUnreceivedAcksMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcChannelUnreceivedAcksMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC next sequence receive
func (e IbcExternal) IbcChannelNextSequence(ibcChannelNextSequenceMsg types.IbcChannelNextSequenceMsg) provider.XplaClient {
	msg, err := MakeIbcChannelNextSequenceReceiveMsg(ibcChannelNextSequenceMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcChannelNextSequenceMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC transfer denom traces
func (e IbcExternal) IbcDenomTraces(ibcDenomTraceMsg ...types.IbcDenomTraceMsg) provider.XplaClient {
	if len(ibcDenomTraceMsg) == 0 {
		msg, err := MakeIbcTransferDenomTracesMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(IbcModule).
			WithMsgType(IbcTransferDenomTracesMsgType).
			WithMsg(msg)
	} else if len(ibcDenomTraceMsg) == 1 {
		msg, err := MakeIbcTransferDenomTraceMsg(ibcDenomTraceMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(IbcModule).
			WithMsgType(IbcTransferDenomTraceMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}

// Query IBC transfer denom trace
func (e IbcExternal) IbcDenomTrace(ibcDenomTraceMsg types.IbcDenomTraceMsg) provider.XplaClient {
	msg, err := MakeIbcTransferDenomTraceMsg(ibcDenomTraceMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcTransferDenomTraceMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC transfer denom hash
func (e IbcExternal) IbcDenomHash(ibcDenomHashMsg types.IbcDenomHashMsg) provider.XplaClient {
	msg, err := MakeIbcTransferDenomHashMsg(ibcDenomHashMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcTransferDenomHashMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC transfer denom hash
func (e IbcExternal) IbcEscrowAddress(ibcEscrowAddressMsg types.IbcEscrowAddressMsg) provider.XplaClient {
	msg, err := MakeIbcTransferEscrowAddressMsg(ibcEscrowAddressMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcTransferEscrowAddressMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query IBC transfer params
func (e IbcExternal) IbcTransferParams() provider.XplaClient {
	msg, err := MakeIbcTransferParamsMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(IbcModule).
		WithMsgType(IbcTransferParamsMsgType).
		WithMsg(msg)
	return e.Xplac
}
