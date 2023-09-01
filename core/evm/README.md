# EVM module
## Usage
### (Tx) Send coin
```go
sendCoinMsg := types.SendCoinMsg{
    Amount: "10000",
    FromAddress: "0x6577385b5d959644ae31263208a88E921273C774",
    ToAddress: "0xF9AC4736D8034F2CB3BFF22A977CD8759934F090",
}

txbytes, err := xplac.EvmSendCoin(sendCoinMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Deploy solidity contract
```go
// Select input type of ABI and bytecode.
// It is also possible to input the entire abi and bytecode as a string type, but you can also enter a file path.
// ABI type must be json and bytecode file must be compiled file on Remix IDE.

// Constructor input arguments
var args []interface{}
owner := []common.Address{
    common.HexToAddress("0xC9F0A2b814d389088a508E31fBa483E8C4372CC2"),
    common.HexToAddress("0x41776240700C033A75A2872EF0AD32b4911e13B1"),
    common.HexToAddress("0xaaC4758A943B2692F2daE0DE8d402aD7045A8DfB"),
}
required := big.NewInt(2)

args = append(args, owner)
args = append(args, required)

deploySolContractMsg := types.DeploySolContractMsg{
    //ABI: `{ ABI json string type }`
    ABIJsonFilePath: "./abi.json",
    // Bytecode: "60806040523480156100......",
    BytecodeJsonFilePath: "./bytecode.json",
    Args: args,
}

txbytes, err := xplac.DeploySolidityContract(deploySolContractMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Invoke(execute) solidity contract
```go
// When invoked, the arguments to be entered into the solidity contract are listed as []interface{}.
var args []interface{}
args = append(a, big.NewInt(2))

// Need contract address and invoke function name.
// Also, same as deployment, need ABI and bytecode.
invokeSolContractMsg := types.InvokeSolContractMsg{
    ContractAddress: "0xBe0AE9A424771C0D68D942A04994a97f928b0821",
    ContractFuncCallName: "store",
    Args: args,
    ABIJsonFilePath: "./abi.json",
    BytecodeJsonFilePath: "./bytecode.json",
}

txbytes, err := xplac.InvokeSolidityContract(invokeSolContractMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) Call solidity contract
```go
callSolContractMsg := types.CallSolContractMsg{
    ContractAddress: "0x80E123317190cAf36292A04776b0De020136526F",
    ContractFuncCallName: "retrieve",
    // Args: nil, // input params if needed to call
    ABIJsonFilePath: "./abi.json",
    BytecodeJsonFilePath: "./bytecode.json",
    FromByteAddress: "0xC9F0A2b814d389088a508E31fBa483E8C4372CC2"
}

res, err := xplac.CallSolidityContract(callSolContractMsg).Query()
```

### (Query) Get transaction by hash
```go
getTransactionByHashMsg := types.GetTransactionByHashMsg {
    TxHash: "556c60576f9af3e4ae7d7fb28f8376e96803c4d9ff02eda6aacb86925f170d09",
}

res, err := xplac.GetTransactionByHash(getTransactionByHashMsg).Query()
```

### (Query) Get Block by hash or height
```go
// Query block by hash
getBlockByHashHeightMsg := types.GetBlockByHashHeightMsg {
    BlockHash: "0xe083b9b3a8b5df69394f55d34cfdfa46e70743a812d7433aba0adf3b7fcecd21",
}

// Query block by height
getBlockByHashHeightMsg := types.GetBlockByHashHeightMsg {
    BlockHeight: "8",
}

res, err := xplac.GetBlockByHashOrHeight(getBlockByHashHeightMsg).Query()
```

### (Query) Account info
```go
// Query account info of user account or contract
// Response of query includes account address(Hex and Bech32), balances and etc. 
// Including Info list
//   - "account" : account address
//   - "bech32_account" : account address of Bech32
//   - "balance" : balances of the account (eth_getBalance)
//   - "nonce" : account nonce as sequence of tendermint based blockchain (eth_getTransactionCount)
//   - "storage" : the storage address for a given account (eth_getStorageAt)
//   - "code" : the contract code of the given account (eth_getCode)
//   - "pending_balance" : the axpla balance of the given account in the pending state (eth_getBalance of the pending state)
//   - "pending_nonce" : the account nonce of the given account in the pending state (eth_getTransactionCount of the pending state)
//   - "pending_storage" : the value of key in the contract storage of the given account in the pending state (eth_getStorageAt of the pending state)
//   - "pending_code" : the contract code of the given account in the pending state (eth_getCode of the pending state)
//   - "pending_transaction_count" : the total number of transactions in the pending state (eth_getBlockTransactionCountByNumber of the pending state)

// so, the xpla.go would not support some RPC APIs as "eth_getBalance", "eth_getTransactionCount", "eth_getStorageAt" and "eth_getCode" because the function is AccountInfo includes these.

accountInfoMsg := types.AccountInfoMsg{
    Account: "0xCa8582862B82867C4Bb9E926682dD75820dE6013",
}

res, err := xplac.AccountInfo(accountInfoMsg).Query()
```

### (Query) Suggest gas price
```go
res, err := xplac.SuggestGasPrice().Query()
```

### (Query) ETH chain ID
```go
res, err := xplac.EthChainID().Query()
```

### (Query) Latest block number
```go
res, err := xplac.EthBlockNumber().Query()
```

### (Query) Web3 client version
```go
res, err = xplac.Web3ClientVersion().Query()
```

### (Query) Web3 SHA3 (return Keccak-256)
```go
web3Sha3Msg := types.Web3Sha3Msg{
    InputParam: "web3-sha3-test",
}

res, err = xplac.Web3Sha3(web3Sha3Msg).Query()
```

### (Query) Network ID
```go
res, err = xplac.NetVersion().Query()
```

### (Query) Network peer count
```go
res, err = xplac.NetPeerCount().Query()
```

### (Query) Network listening
```go
res, err = xplac.NetListening().Query()
```

### (Query) Ethereum protocol version
```go
res, err = xplac.EthProtocolVersion().Query()
```

### (Query) Ethereum syncing
```go
res, err = xplac.EthSyncing().Query()
```

### (Query) Eth accounts
```go
res, err = xplac.EthAccounts().Query()
```

### (Query) The number of transactions in a given block 
```go
// using block height(=number)
e := types.EthGetBlockTransactionCountMsg{
    BlockHeight: "5440",
}

// using block hash
e := types.EthGetBlockTransactionCountMsg{
    BlockHeight: "0x46b3031b22f065f933331dc032ccd34404282ccf7e4fcd54e02d1f808abc112c"
}

res, err = xplac.EthGetBlockTransactionCount(e).Query()
```

### (Query) Estimate gas to contract
```go
var args []interface{}
args = append(args, big.NewInt(6151212))

// invoke message to estimate
invokeSolContractMsg := types.InvokeSolContractMsg{
    ContractAddress:      c.ContractAddress,
    ContractFuncCallName: "store",
    Args:                 args,
    ABIJsonFilePath:      "./testfiles/abi.json",
    BytecodeJsonFilePath: "./testfiles/bytecode.json",
}

res, err = xplac.EstimateGas(invokeSolContractMsg).Query()
```

### (Query) Get transaction by block hash and index
```go
getTransactionByBlockHashAndIndexMsg := types.GetTransactionByBlockHashAndIndexMsg{
    BlockHash: "0x7f562573c1b0ca6fc3a83246372a5d57f917a4c654c91b65ebd756dec4989d0f",
    Index:     "0",
}

res, err = xplac.EthGetTransactionByBlockHashAndIndex(getTransactionByBlockHashAndIndexMsg).Query()
```

### (Query) Get transaction receipt
```go
getTransactionReceiptMsg := types.GetTransactionReceiptMsg{
    TransactionHash: "0x20ec56d16231c4d7f761c2533885619489fface85cf6c478868ef1d531b93177",
}

res, err = xplac.EthGetTransactionReceipt(getTransactionReceiptMsg).Query()
```

### (Query) New filter
```go
ethNewFilterMsg := types.EthNewFilterMsg{
    Topics:    []string{"0x20ec56d16231c4d7f761c2533885619489fface85cf6c478868ef1d531b93177"},
    Address:   []string{"0xf7777b36a51fb0b33dd0c5118361AfC94ff7f967"},
    ToBlock:   "latest",
    FromBlock: "earliest",
}

res, err = xplac.EthNewFilter(ethNewFilterMsg).Query()
```
 
### (Query) New block filter
```go
res, err = xplac.EthNewBlockFilter().Query()
```

### (Query) New pending transaction filter
```go
res, err = xplac.EthNewPendingTransactionFilter().Query()
```

### (Query) Uninstall filter
```go
ethUninsatllFilterMsg := types.EthUninsatllFilterMsg{
    FilterId: "0x168b9d421ecbffa1ac706926c2203454",
}

res, err = xplac.EthUninstallFilter(ethUninsatllFilterMsg).Query()
```

### (Query) Get filter changes
```go
ethGetFilterChangesMsg := types.EthGetFilterChangesMsg{
    FilterId: "0x9852d91813fb44da471436722e02965e",
}

res, err = xplac.EthGetFilterChanges(ethGetFilterChangesMsg).Query()
```

### (Query) Get logs
```go
ethGetLogsMsg := types.EthGetLogsMsg{
    Topics:  []string{"0x20ec56d16231c4d7f761c2533885619489fface85cf6c478868ef1d531b93177"},
    Address: []string{"0xf7777b36a51fb0b33dd0c5118361AfC94ff7f967"},
    ToBlock: "latest",
    FromBlock: "latest",
    // BlockHash: "0x46b3031b22f065f933331dc032ccd34404282ccf7e4fcd54e02d1f808abc112c",
}

res, err = xplac.EthGetLogs(ethGetLogsMsg).Query()
```

### (Query) Coinbase
```go
res, err = xplac.EthCoinbase().Query()
```