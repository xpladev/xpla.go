# Bank module
## Usage
### (Tx) Bank coin send
```go
// from address, to address, coin amount(axpla)
bankSendMsg := types.BankSendMsg {
    FromAddress: "xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7", 
    ToAddress: "xpla13trl452wgle9qxpxhse9605k9x0399cmkfzn7g", 
    Amount: "10",
}
txbytes, err := xplac.BankSend(bankSendMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) Bank all balances & denom balance
```go
// All balances
bankBalancesMsg := types.BankBalancesMsg {
    Address: "xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7",		
}
response, err := xplac.BankBalances(bankBalancesMsg).Query()	

// denom balances
bankBalancesMsg := types.BankBalancesMsg {
    Address: "xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7",		
    Denom: "axpla"
}
response, err := xplac.BankBalances(bankBalancesMsg).Query()	
```

### (Query) Bank denom metadata
```go
// All metadata
denomMetadataMsg := types.DenomMetadataMsg{}
response, err := xplac.DenomMetadata(denomMetadataMsg).Query()

// client metadata for denom
denomMetadataMsg := types.DenomMetadataMsg{
    Denom: "axpla",
}
response, err := xplac.DenomMetadata(denomMetadataMsg).Query()
```

### (Query) Bank total supply
```go
// Total supply
response, err := xplac.Total().Query()

// Total supply of denom
totalMsg := types.TotalMsg {
    Denom: "axpla",
}
response, err := xplac.Total(totalMsg).Query()
```
