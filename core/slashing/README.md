# Slashing module
## Usage
### (Tx) Unjail validator
```go
txbytes, err := xplac.Unjail().CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) slashing params
```go
res, err := xplac.SlashingParams().Query()
```

### (Query) slashing signing infos
```go
// Query a validator's signing information by using public key
signingInfoMsg := types.SigningInfoMsg{
    ConsPubKey: `{"@type": "/cosmos.crypto.ed25519.PubKey","key": "6RBPm24ckoWhRt8mArcSCnEKvt0FMGvcaMwchfZ3ue8="}`,
}

// Query a validator's signing information by using bech32 address
signingInfoMsg := types.SigningInfoMsg{
    ConsAddr: "xplavalcons1v9jz99h7dsf50fgwr3wr2v8d73dfc3m8qvuaah",
}

res, err := xplac.SigningInfos(signingInfoMsg).Query()

// Query signing information of all validators
res, err := xplac.SigningInfos().Query()
```
