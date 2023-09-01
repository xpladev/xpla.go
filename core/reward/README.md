# Reward module
## Usage
### (Tx) Fund fee collector for reward module
```go
// fund fee collector test
fundFeeCollectorMsg := types.FundFeeCollectorMsg{
    DepositorAddr: "xpla1j55tymfdys9n7k0dq6xmyd4hgfelp9jghzympt",
    Amount:        "1000",
}

txbytes, err := xplac.FundFeeCollector(fundFeeCollectorMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) Query reward module params
```go
res, err := xplac.RewardParams().Query()
```

### (Query) Query pool amount of reward module
```go
res, err := xplac.RewardPool().Query()
```
