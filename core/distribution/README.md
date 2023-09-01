# Distribution module
## Usage
### (Tx) Fund community pool
```go
fundCommunityPoolMsg := types.FundCommunityPoolMsg {
    Amount: "1000",
}

txbytes, err := xplac.FundCommunityPool(fundCommunityPoolMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```
### (Tx) Proposal community pool spend
```go
// Input data
communityPoolSpendMsg := types.CommunityPoolSpendMsg{
    Title: "community pool spend",
    Description: "pay me",
    Recipient: "xpla1ka84cuec6339t8s4nh3sp5zf2fre6dh2v2g9mp",
    Amount: "10000",
    Deposit: "1000",
}

// Use json file
communityPoolSpendMsg := types.CommunityPoolSpendMsg{
    JsonFilePath: "./proposal.json"
}

txbytes, err := xplac.CommunityPoolSpend(communityPoolSpendMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Withdraw rewards
```go
withdrawRewardsMsg := types.WithdrawRewardsMsg{
    DelegatorAddr: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    ValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}
txbytes, err := xplac.WithdrawRewards(withdrawRewardsMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Withdraw all rewards
```go
txbytes, err := xplac.WithdrawAllRewards().CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Set withdraw address
```go
setWithdrawAddrMsg := types.SetwithdrawAddrMsg {
    WithdrawAddr: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
}
txbytes, err := xplac.SetWithdrawAddr(setWithdrawAddrMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) distribution parameters
```go
res, err := xplac.DistributionParams().Query()
```

### (Query) validator outstanding rewards
```go
validatorOutstandingRewardsMsg := types.ValidatorOutstandingRewardsMsg {
    ValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}
res, err := xplac.ValidatorOutstandingRewards(validatorOutstandingRewardsMsg).Query()
```

### (Query) commission
```go
queryDistCommissionMsg := types.QueryDistCommissionMsg {
    ValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}
res, err := xplac.DistCommission(queryDistCommissionMsg).Query()
```

### (Query) validator slashes
```go
queryDistSlashesMsg := types.QueryDistSlashesMsg{
    ValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}

res, err := xplac.DistSlashes(queryDistSlashesMsg).Query()
```
### (Query) rewards
```go
// reward for a validator
queryDistRewardsMsg := types.QueryDistRewardsMsg{
    ValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
    DelegatorAddr: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
}

// total rewards
queryDistRewardsMsg := types.QueryDistRewardsMsg{
    DelegatorAddr: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
}
res, err := xplac.DistRewards(queryDistRewardsMsg).Query()
```
### (Query) community pool
```go
res, err := xplac.CommunityPool().Query()
```