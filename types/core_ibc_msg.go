package types

type IbcClientStateMsg struct {
	ClientId string
}

type IbcClientStatusMsg struct {
	ClientId string
}

type IbcClientConsensusStatesMsg struct {
	ClientId string
}

type IbcClientConsensusStateHeightsMsg struct {
	ClientId string
}

type IbcClientConsensusStateMsg struct {
	ClientId     string
	Height       string
	LatestHeight bool
}

type IbcConnectionMsg struct {
	ConnectionId string
}

type IbcClientConnectionsMsg struct {
	ClientId string
}

type IbcChannelMsg struct {
	PortId    string
	ChannelId string
}

type IbcChannelConnectionsMsg struct {
	ConnectionId string
}

type IbcChannelClientStateMsg struct {
	ChannelId string
	PortId    string
}

type IbcChannelPacketCommitmentsMsg struct {
	ChannelId string
	PortId    string
	Sequence  string
}

type IbcChannelPacketReceiptMsg struct {
	ChannelId string
	PortId    string
	Sequence  string
}

type IbcChannelPacketAckMsg struct {
	ChannelId string
	PortId    string
	Sequence  string
}

type IbcChannelUnreceivedPacketsMsg struct {
	ChannelId string
	PortId    string
	Sequence  string
}

type IbcChannelUnreceivedAcksMsg struct {
	ChannelId string
	PortId    string
	Sequence  string
}

type IbcChannelNextSequenceMsg struct {
	ChannelId string
	PortId    string
}

type IbcDenomTraceMsg struct {
	HashDenom string
}

type IbcDenomHashMsg struct {
	Trace string
}

type IbcEscrowAddressMsg struct {
	ChannelId string
	PortId    string
}
