# Cosmos/Tendermint base
## Usage
### (Query) Node info
```go
// Query node info
res, err = xplac.NodeInfo().Query()
```

### (Query) Syncing
```go
// Query syncing
res, err = xplac.Syncing().Query()
```

### (Query) Blocks
```go
// Query latest block info
res, err = xplac.Block().Query()

// Query block info by height
blockMsg := types.BlockMsg{
    Height: "1",
}
res, err = xplac.Block(blockMsg).Query()
```

### (Query) Validator set
```go
// Query latest validator set
res, err = xplac.ValidatorSet().Query()

// QUery validator set by height
validatorSetMsg := types.ValidatorSetMsg{
    Height: "1",
}
res, err = xplac.ValidatorSet(validatorSetMsg).Query()
```