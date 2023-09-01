# Auth module
## Usage
### (Query) auth params
```go
response, err := xplac.AuthParams().Query()
```

### (Query) account address
```go
queryAccAddressMsg := types.QueryAccAddressMsg{
    Address: "xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7",
}
response, err := xplac.AccAddress(queryAccAddressMsg).Query()
```

### (Query) accounts
```go
response, err := xplac.Accounts().Query()
```

### (Query) Txs by events
```go
queryTxsByEventsMsg := types.QueryTxsByEventsMsg{
    Events: "transfer.recipient=xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7",
}
response, err := xplac.TxsByEvents(queryTxsByEventsMsg).Query()
```

### (Query) tx
```go
// Retrieve by using hash
queryTxMsg := types.QueryTxMsg{
    // Default type is "hash" including empty type.
    Value: "B6BBBB649F19E8970EF274C0083FE945FD38AD8C524D68BB3FE3A20D72DF03C4",
}
response, err := xplac.Tx(queryTxMsg).Query()

// Retrieve by using signature
queryTxMsg := types.QueryTxMsg{
    Type:  "signature",
    Value: "4fmwN0Qp084qpfNpm1XV22YOwnjrGYWIuyRgGgj7f3Mv2ECsQ0ZY/9MqOaZ9TGB3slQQQNNiiBf9eR2ACad/pgE=",
}

// Retrieve by using account sequence
queryTxMsg := types.QueryTxMsg{
		Type:  "acc_seq",
		Value: "xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7/5", //<addr>/<sequence>
	}
response, err := xplac.Tx(queryTxMsg).Query()
```