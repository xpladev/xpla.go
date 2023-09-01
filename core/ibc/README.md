# IBC module
## Usage
### (Query) Client states
```go
// Query IBC client states
res, err = xplac.IbcClientStates().Query()
```

### (Query) Client state by client-ID
```go
// Query IBC client state
ibcClientStatesMsg := types.IbcClientStateMsg{
    ClientId: "07-tendermint-0",
}
res, err = xplac.IbcClientState(ibcClientStateMsg).Query()
```

### (Query) Client status by client-ID
```go
// Query IBC client status
ibcClientStatusMsg := types.IbcClientStatusMsg{
    ClientId: "07-tendermint-0",
}
res, err = xplac.IbcClientStatus(ibcClientStatusMsg).Query()
```

### (Query) Client consensus states
```go
// Query IBC client consensus states
ibcClientConsensusStatesMsg := types.IbcClientConsensusStatesMsg{
    ClientId: "07-tendermint-0",
}
res, err = xplac.IbcClientConsensusStates(ibcClientConsensusStatesMsg).Query()
```

### (Query) Client consensus state heights
```go
// Query IBC client consensus state heights
ibcClientConsensusStateHeightsMsg := types.IbcClientConsensusStateHeightsMsg{
    ClientId: "07-tendermint-0",
}
res, err = xplac.IbcClientConsensusStateHeights(ibcClientConsensusStateHeightsMsg).Query()
```

### (Query) Client consensus state
```go
// Query IBC client consensus state
ibcClientConsensusStateMsg := types.IbcClientConsensusStateMsg{
    ClientId:     "07-tendermint-0",
    LatestHeight: false,
    Height:       "1-115",
}
res, err = xplac.IbcClientConsensusState(ibcClientConsensusStateMsg).Query()
```

### (Query) Client tendermint header
```go
// Query IBC client tendermint header
res, err = xplac.IbcClientHeader().Query()
```

### (Query) Client self consensus state
```go
// Query IBC client self-consensus state
res, err = xplac.IbcClientSelfConsensusState().Query()
```

### (Query) Client params
```go
// Query IBC client params
res, err = xplac.IbcClientParams().Query()
```

### (Query) IBC Connections
```go
// Query IBC all connections
res, err = xplac.IbcConnections().Query()

// Query IBC a connection by retriving connection ID
ibcConnectionMsg := types.IbcConnectionMsg{
    ConnectionId: "connection-1",
}
res, err = xplac.IbcConnections(ibcConnectionMsg).Query()
```

### (Query) A client connections
```go
// Query IBC connections of a client
IbcClientConnectionsMsg := types.IbcClientConnectionsMsg{
    ClientId: "07-tendermint-0",
}
res, err = xplac.IbcClientConnections(IbcClientConnectionsMsg).Query()
```

### (Query) IBC channels
```go
// Query IBC all channels
res, err = xplac.IbcChannels().Query()

// Query IBC a channel
ibcChannelMsg := types.IbcChannelMsg{
    ChannelId: "channel-0",
    PortId:    "transfer",
}
res, err = xplac.IbcChannels(ibcChannelMsg).Query()
```

### (Query) Channel connections
```go
// Query channels by retrieving a connection
IbcChannelConnectionsMsg := types.IbcChannelConnectionsMsg{
    ConnectionId: "connection-0",
}
res, err = xplac.IbcChannelConnections(IbcChannelConnectionsMsg).Query()
```

### (Query) Channel client state
```go
// Query channel client state
ibcChannelClientStateMsg := types.IbcChannelClientStateMsg{
    ChannelId: "channel-0",
    PortId:    "transfer",
}
res, err = xplac.IbcChannelClientState(ibcChannelClientStateMsg).Query()
```

### (Query) Channel packet commitments
```go
// Query IBC channel packet commitments
ibcChannelPacketCommitmentsMsg := types.IbcChannelPacketCommitmentsMsg{
    ChannelId: "channel-0",
    PortId:    "transfer",
}
res, err = xplac.IbcChannelPacketCommitments(ibcChannelPacketCommitmentsMsg).Query()

// Query IBC channel packet commitment by sequence
ibcChannelPacketCommitmentsMsg := types.IbcChannelPacketCommitmentsMsg{
    ChannelId: "channel-0",
    PortId:    "transfer",
    Sequence:  "2",
}
res, err = xplac.IbcChannelPacketCommitments(ibcChannelPacketCommitmentsMsg).Query()
```

### (Query) Channel packet receipt
```go
// Query IBC channel packet receipt
ibcChannelPacketReceiptMsg := types.IbcChannelPacketReceiptMsg{
    ChannelId: "channel-0",
    PortId:    "transfer",
    Sequence:  "2",
}
res, err = xplac.IbcChannelPacketReceipt(ibcChannelPacketReceiptMsg).Query()
```

### (Query) Channel packet ack
```go
// Query IBC channel packet ack
ibcChannelPacketAckMsg := types.IbcChannelPacketAckMsg{
    ChannelId: "channel-0",
    PortId:    "transfer",
    Sequence:  "2",
}
res, err = xplac.IbcChannelPacketAck(ibcChannelPacketAckMsg).Query()
```

### (Query) Channel unreceived packets
```go
// Query unreceived packets
ibcChannelUnreceivedPacketsMsg := types.IbcChannelUnreceivedPacketsMsg{
    ChannelId: "channel-0",
    PortId:    "transfer",
    Sequence:  "2",
}
res, err = xplac.IbcChannelUnreceivedPackets(ibcChannelUnreceivedPacketsMsg).Query()
```

### (Query) Channel unreceived acks
```go
// Query unreceived acks
ibcChannelUnreceivedAcksMsg := types.IbcChannelUnreceivedAcksMsg{
    ChannelId: "channel-0",
    PortId:    "transfer",
    Sequence:  "2",
}
res, err = xplac.IbcChannelUnreceivedAcks(ibcChannelUnreceivedAcksMsg).Query()
```

### (Query) Channel next sequence
```go
// Query channel next sequence receive
ibcChannelNextSequenceMsg := types.IbcChannelNextSequenceMsg{
    ChannelId: "channel-0",
    PortId:    "transfer",
}
res, err = xplac.IbcChannelNextSequence(ibcChannelNextSequenceMsg).Query()
```

### (Query) IBC denom traces
```go
// Query IBC transfer denom traces
res, err = xplac.IbcDenomTraces().Query()

// Query IBC transfer a denom trace
ibcDenomTraceMsg := types.IbcDenomTraceMsg{
    HashDenom: "B249D1E86F588286FEA286AA8364FFCE69EC65604BD7869D824ADE40F00FA25B",
}
res, err = xplac.IbcDenomTraces().Query()
```

### (Query) IBC denom hash
```go
// Make denom hash by trace
ibcDenomHashMsg := types.IbcDenomHashMsg{
    Trace: "[port-id]/[channel-id]/[denom]",
}
res, err = xplac.IbcDenomHash(ibcDenomHashMsg).Query()
```

### (Query) IBC escrow address
```go
// Query escrow address
ibcEscrowAddressMsg := types.IbcEscrowAddressMsg{
    PortId:    "transfer",
    ChannelId: "channel-5",
}
res, err = xplac.IbcEscrowAddress(ibcEscrowAddressMsg).Query()
```

### (Query) IBC transfer params
```go
// Query IBC transfer params
res, err = xplac.IbcTransferParams().Query()
```