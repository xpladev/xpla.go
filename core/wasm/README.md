# Wasm module
## Usage
### (Tx) Store code
```go
// can instantiate only store msg sender
storeMsg := types.StoreMsg {
    FilePath: "./wasmcontract.wasm",
    InstantiatePermission: "instantiate-only-sender" //optional. if method is empty string, everybody
}

// can instantiate address
// method and address must be seperated by "."
storeMsg := types.StoreMsg {
    FilePath: "./wasmcontract.wasm",
    InstantiatePermission: "instantiate-only-address.xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7" 
}

txbytes, err := xplac.StoreCode(storeMsg).CreateAndSignTx()
res, _ := xplac.Broadcast(txbytes)
```

### (Tx) Instantiate contract
```go
instantiateMsg := types.InstantiateMsg {
    CodeId: "1",
    Amount: "10",
    Label: "Contract instant",
    InitMsg: `{"owner":"xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7"}`,
}
txbytes, err := xplac.InstantiateContract(instantiateMsg).CreateAndSignTx()
res, _ := xplac.Broadcast(txbytes)
```

### (Tx) Execute contract
```go
executeMsg := types.ExecuteMsg {
    ContractAddress: "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h",
    Amount: "0",
    ExecMsg: `{"execute_method":{"execute_key":"execute_test","execute_value":"xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7"}}`,
}
txbytes, err := xplac.ExecuteContract(executeMsg).CreateAndSignTx()
res, _ := xplac.Broadcast(txbytes)
```

### (Tx) Clear contract admin
```go
clearContractAdminMsg := types.ClearContractAdminMsg{
    ContractAddress: "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h",
}
txbytes, err := xplac.ClearContractAdmin(clearContractAdminMsg).CreateAndSignTx()
res, _ := xplac.Broadcast(txbytes)
```

### (Tx) Set contract admin
```go
setContractAdminMsg := types.SetContractAdminMsg{
    ContractAddress: "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h",
    NewAdmin: "xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7",
}

txbytes, err := xplac.SetContractAdmin(setContractAdminMsg).CreateAndSignTx()
res, _ := xplac.Broadcast(txbytes)
```

### (Tx) Migrate
```go
migrateMsg := types.MigrateMsg{
    ContractAddress: "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h",
    CodeId: "123",
    MigrateMsg: "json encoded migrate msg",
}

txbytes, err := xplac.Migrate(migrateMsg).CreateAndSignTx()
```

### (Query) contract
```go
queryMsg := types.QueryMsg {
    ContractAddress: "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h",
    QueryMsg: `{"query_method":{"query":"query_test"}}`,
}

response, _ := xplac.QueryContract(queryMsg).Query()
```

### (Query) list code
```go
response, _ := xplac.ListCode().Query()
```

### (Query) list contract by code
```go
listContractByCodeMsg := types.ListContractByCodeMsg {
    CodeId: "1",
}

response, err := xplac.ListContractByCode(listContractByCodeMsg).Query()
```

### (Query) Download contract wasm file
```go
downloadMsg := types.DownloadMsg{
    CodeId: "1",
    DownloadFileName: "test",
}

response, err := xplac.Download(downloadMsg).Query()
```

### (Query) code info
```go
codeInfoMsg := types.CodeInfoMsg{
    CodeId: "1",
}

response, err := xplac.CodeInfo(codeInfoMsg).Query()
```

### (Query) Contract info
```go
contractInfoMsg := types.ContractInfoMsg {
    ContractAddress: "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h",
}

response, err := xplac.ContractInfo(contractInfoMsg).Query()
```

### (Query) all contract state
```go
contractStateAllMsg := types.ContractStateAllMsg{
    ContractAddress: "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h",
}

response, err := xplac.ContractStateAll(contractStateAllMsg).Query()
```

### (Query) contract history
```go
contractHistoryMsg := types.ContractHistoryMsg {
    ContractAddress: "xpla14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s525s0h",
}

response, err := xplac.ContractHistory(contractHistoryMsg).Query()
```

### (Query) pinned
```go
response, err := xplac.Pinned().Query()
```

### (Query) Lib wasmvm version
```go
response, err := xplac.LibwasmvmVersion().Query()
```
