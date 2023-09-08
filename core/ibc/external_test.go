package ibc_test

import (
	mibc "github.com/xpladev/xpla.go/core/ibc"
	"github.com/xpladev/xpla.go/types"
)

var (
	testIbcClientID      = "07-tendermint-0"
	testIbcConnectionID  = "connection-1"
	testIbcChannelID     = "channel-0"
	testIbcChannelPortId = "transfer"
)

func (s *IntegrationTestSuite) TestIBC() {
	// client states
	s.xplac.IbcClientStates()

	makeIbcClientStatesMsg, err := mibc.MakeIbcClientStatesMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientStatesMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientStatesMsgType, s.xplac.GetMsgType())

	// client state
	ibcClientStateMsg := types.IbcClientStateMsg{
		ClientId: testIbcClientID,
	}
	s.xplac.IbcClientState(ibcClientStateMsg)

	makeIbcClientStateMsg, err := mibc.MakeIbcClientStateMsg(ibcClientStateMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientStateMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientStateMsgType, s.xplac.GetMsgType())

	// client status
	ibcClientStatusMsg := types.IbcClientStatusMsg{
		ClientId: testIbcClientID,
	}
	s.xplac.IbcClientStatus(ibcClientStatusMsg)

	makeIbcClientStatusMsg, err := mibc.MakeIbcClientStatusMsg(ibcClientStatusMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientStatusMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientStatusMsgType, s.xplac.GetMsgType())

	// client consensus states
	ibcClientConsensusStatesMsg := types.IbcClientConsensusStatesMsg{
		ClientId: testIbcClientID,
	}
	s.xplac.IbcClientConsensusStates(ibcClientConsensusStatesMsg)

	makeIbcClientConsensusStatesMsg, err := mibc.MakeIbcClientConsensusStatesMsg(ibcClientConsensusStatesMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientConsensusStatesMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientConsensusStatesMsgType, s.xplac.GetMsgType())

	// client consensus state heights
	ibcClientConsensusStateHeightsMsg := types.IbcClientConsensusStateHeightsMsg{
		ClientId: testIbcClientID,
	}
	s.xplac.IbcClientConsensusStateHeights(ibcClientConsensusStateHeightsMsg)

	makeIbcClientConsensusStateHeightsMsg, err := mibc.MakeIbcClientConsensusStateHeightsMsg(ibcClientConsensusStateHeightsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientConsensusStateHeightsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientConsensusStateHeightsMsgType, s.xplac.GetMsgType())

	// client consensus state
	ibcClientConsensusStateMsg := types.IbcClientConsensusStateMsg{
		ClientId:     testIbcClientID,
		LatestHeight: false,
		Height:       "1-115",
	}
	s.xplac.IbcClientConsensusState(ibcClientConsensusStateMsg)

	makeIbcClientConsensusStateMsg, err := mibc.MakeIbcClientConsensusStateMsg(ibcClientConsensusStateMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientConsensusStateMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientConsensusStateMsgType, s.xplac.GetMsgType())

	// client params
	s.xplac.IbcClientParams()

	makeIbcClientParamsMsg, err := mibc.MakeIbcClientParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcClientParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcClientParamsMsgType, s.xplac.GetMsgType())

	// connections
	s.xplac.IbcConnections()

	makeIbcConnectionConnectionsMsg, err := mibc.MakeIbcConnectionConnectionsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcConnectionConnectionsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcConnectionConnectionsMsgType, s.xplac.GetMsgType())

	// connection
	ibcConnectionMsg := types.IbcConnectionMsg{
		ConnectionId: testIbcConnectionID,
	}
	s.xplac.IbcConnections(ibcConnectionMsg)

	makeIbcConnectionConnectionMsg, err := mibc.MakeIbcConnectionConnectionMsg(ibcConnectionMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcConnectionConnectionMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcConnectionConnectionMsgType, s.xplac.GetMsgType())

	// client connection
	ibcClientConnectionsMsg := types.IbcClientConnectionsMsg{
		ClientId: testIbcClientID,
	}
	s.xplac.IbcClientConnections(ibcClientConnectionsMsg)

	makeIbcConnectionClientConnectionsMsg, err := mibc.MakeIbcConnectionClientConnectionsMsg(ibcClientConnectionsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcConnectionClientConnectionsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcConnectionClientConnectionsMsgType, s.xplac.GetMsgType())

	// channels
	s.xplac.IbcChannels()

	makeIbcChannelChannelsMsg, err := mibc.MakeIbcChannelChannelsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelChannelsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelChannelsMsgType, s.xplac.GetMsgType())

	// channel
	ibcChannelMsg := types.IbcChannelMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannels(ibcChannelMsg)

	makeIbcChannelChannelMsg, err := mibc.MakeIbcChannelChannelMsg(ibcChannelMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelChannelMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelChannelMsgType, s.xplac.GetMsgType())

	// channel connections
	ibcChannelConnectionsMsg := types.IbcChannelConnectionsMsg{
		ConnectionId: testIbcConnectionID,
	}
	s.xplac.IbcChannelConnections(ibcChannelConnectionsMsg)

	makeIbcChannelConnectionsMsg, err := mibc.MakeIbcChannelConnectionsMsg(ibcChannelConnectionsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelConnectionsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelConnectionsMsgType, s.xplac.GetMsgType())

	// channel client state
	ibcChannelClientStateMsg := types.IbcChannelClientStateMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelClientState(ibcChannelClientStateMsg)

	makeIbcChannelClientStateMsg, err := mibc.MakeIbcChannelClientStateMsg(ibcChannelClientStateMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelClientStateMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelClientStateMsgType, s.xplac.GetMsgType())

	// channel packet commitments
	ibcChannelPacketCommitmentsMsg := types.IbcChannelPacketCommitmentsMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelPacketCommitments(ibcChannelPacketCommitmentsMsg)

	makeIbcChannelPacketCommitmentsMsg, err := mibc.MakeIbcChannelPacketCommitmentsMsg(ibcChannelPacketCommitmentsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketCommitmentsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelPacketCommitmentsMsgType, s.xplac.GetMsgType())

	// channel packet commitment
	ibcChannelPacketCommitmentsMsg = types.IbcChannelPacketCommitmentsMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
		Sequence:  "1",
	}
	s.xplac.IbcChannelPacketCommitments(ibcChannelPacketCommitmentsMsg)

	makeIbcChannelPacketCommitmentMsg, err := mibc.MakeIbcChannelPacketCommitmentMsg(ibcChannelPacketCommitmentsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketCommitmentMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelPacketCommitmentMsgType, s.xplac.GetMsgType())

	// packet receipt
	ibcChannelPacketReceiptMsg := types.IbcChannelPacketReceiptMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelPacketReceipt(ibcChannelPacketReceiptMsg)

	makeIbcChannelPacketReceiptMsg, err := mibc.MakeIbcChannelPacketReceiptMsg(ibcChannelPacketReceiptMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketReceiptMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelPacketReceiptMsgType, s.xplac.GetMsgType())

	// packet ack
	ibcChannelPacketAckMsg := types.IbcChannelPacketAckMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelPacketAck(ibcChannelPacketAckMsg)

	makeIbcChannelPacketAckMsg, err := mibc.MakeIbcChannelPacketAckMsg(ibcChannelPacketAckMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketAckMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelPacketAckMsgType, s.xplac.GetMsgType())

	// unreceived packets
	ibcChannelUnreceivedPacketsMsg := types.IbcChannelUnreceivedPacketsMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelUnreceivedPackets(ibcChannelUnreceivedPacketsMsg)

	makeIbcChannelPacketUnreceivedPacketsMsg, err := mibc.MakeIbcChannelPacketUnreceivedPacketsMsg(ibcChannelUnreceivedPacketsMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketUnreceivedPacketsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelUnreceivedPacketsMsgType, s.xplac.GetMsgType())

	// unreceived acks
	ibcChannelUnreceivedAcksMsg := types.IbcChannelUnreceivedAcksMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelUnreceivedAcks(ibcChannelUnreceivedAcksMsg)

	makeIbcChannelPacketUnreceivedAcksMsg, err := mibc.MakeIbcChannelPacketUnreceivedAcksMsg(ibcChannelUnreceivedAcksMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelPacketUnreceivedAcksMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelUnreceivedAcksMsgType, s.xplac.GetMsgType())

	// channel next sequence
	ibcChannelNextSequenceMsg := types.IbcChannelNextSequenceMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcChannelNextSequence(ibcChannelNextSequenceMsg)

	makeIbcChannelNextSequenceReceiveMsg, err := mibc.MakeIbcChannelNextSequenceReceiveMsg(ibcChannelNextSequenceMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcChannelNextSequenceReceiveMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcChannelNextSequenceMsgType, s.xplac.GetMsgType())

	// denom traces
	s.xplac.IbcDenomTraces()

	makeIbcTransferDenomTracesMsg, err := mibc.MakeIbcTransferDenomTracesMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferDenomTracesMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferDenomTracesMsgType, s.xplac.GetMsgType())

	// denom trace
	ibcDenomTraceMsg := types.IbcDenomTraceMsg{
		HashDenom: "B249D1E86F588286FEA286AA8364FFCE69EC65604BD7869D824ADE40F00FA25B",
	}
	s.xplac.IbcDenomTraces(ibcDenomTraceMsg)

	makeIbcTransferDenomTraceMsg, err := mibc.MakeIbcTransferDenomTraceMsg(ibcDenomTraceMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferDenomTraceMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferDenomTraceMsgType, s.xplac.GetMsgType())

	// denom trace
	ibcDenomTraceMsg = types.IbcDenomTraceMsg{
		HashDenom: "B249D1E86F588286FEA286AA8364FFCE69EC65604BD7869D824ADE40F00FA25B",
	}
	s.xplac.IbcDenomTrace(ibcDenomTraceMsg)

	makeIbcTransferDenomTraceMsg, err = mibc.MakeIbcTransferDenomTraceMsg(ibcDenomTraceMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferDenomTraceMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferDenomTraceMsgType, s.xplac.GetMsgType())

	// denom hash
	ibcDenomHashMsg := types.IbcDenomHashMsg{
		Trace: testIbcChannelPortId + "/" + testIbcChannelID + "/" + types.XplaDenom,
	}
	s.xplac.IbcDenomHash(ibcDenomHashMsg)

	makeIbcTransferDenomHashMsg, err := mibc.MakeIbcTransferDenomHashMsg(ibcDenomHashMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferDenomHashMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferDenomHashMsgType, s.xplac.GetMsgType())

	// escrow address
	ibcEscrowAddressMsg := types.IbcEscrowAddressMsg{
		ChannelId: testIbcChannelID,
		PortId:    testIbcChannelPortId,
	}
	s.xplac.IbcEscrowAddress(ibcEscrowAddressMsg)

	makeIbcTransferEscrowAddressMsg, err := mibc.MakeIbcTransferEscrowAddressMsg(ibcEscrowAddressMsg)
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferEscrowAddressMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferEscrowAddressMsgType, s.xplac.GetMsgType())

	// escrow address
	s.xplac.IbcTransferParams()

	makeIbcTransferParamsMsg, err := mibc.MakeIbcTransferParamsMsg()
	s.Require().NoError(err)

	s.Require().Equal(makeIbcTransferParamsMsg, s.xplac.GetMsg())
	s.Require().Equal(mibc.IbcModule, s.xplac.GetModule())
	s.Require().Equal(mibc.IbcTransferParamsMsgType, s.xplac.GetMsgType())
}
