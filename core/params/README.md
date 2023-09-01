# Params module
## Usage
### (Tx) Proposal parameter change
```go
// Input arguments
paramChangeMsg := types.ParamChangeMsg{
    Title: "Staking param change",
    Description: "update max validators",
    Changes: []string{
        `{
            "subspace": "staking",
            "key": "MaxValidators",
            "value": 105
        }`,
    },
    Deposit: "1000",
}

// Input json file
paramChangeMsg := types.ParamChangeMsg{
    JsonFilePath: "./proposal.json"
}

txbytes, err := xplac.ParamChange(paramChangeMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```
### (Query) Query params by using subspace
```go
subspaceMsg := types.SubspaceMsg{
    Subspace: "staking",
    Key: "MaxValidators",
}

res, err := xplac.QuerySubspace(subspaceMsg).Query()
```
