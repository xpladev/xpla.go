# Xpla client
The xpla client is a client for performing all functions within the xpla.go library.
The user mandatorily inputs chain ID.

## Ready to run Xpla client
```go
mnemonic, err := key.NewMnemonic()
if err != nil {
    fmt.Println(err)
}

priKey, err := key.NewPrivKey(mnemonic)
if err != nil {
    fmt.Println(err)
}

// Can check addr (string type)
addr, err := key.Bech32AddrString(priKey)

// Create new XPLA client
xplac := client.NewXplaClient("chain-id")

// Set private key
xplac = xplac.WithOptions(client.Options{PrivateKey: priKey})
```
### Set URLs for xpla client
```go
// Need LCD URL when broadcast transactions
xplac := client.NewXplaClient(
    "chain-id",    
).WithOptions(
    client.Options{
        LcdURL: "http://localhost:1317",
    }
)

// Need GRPC URL to query or broadcast tx
xplac := client.NewXplaClient(
    "chain-id",    
).WithOptions(
    client.Options{
        GrpcURL: "http://localhost:9090",
    }
)

// Need tendermint RPC URL when only "Query tx" methods
// i.e. xplad query tx, xplad query txs
xplac := client.NewXplaClient(
    "chain-id",    
).WithOptions(
    client.Options{
        RpcURL: "http://localhost:26657",
    }
)

// Need EVM RPC URL when use evm module
xplac := client.NewXplaClient(
    "chain-id",    
).WithOptions(
    client.Options{
        EvmRpcURL: "http://localhost:8545",
    }
)
```

### Optional parameters of xpla client
```go
type Options struct {    
    // Set private key
    PrivateKey     key.PrivateKey
    // Set account number of address
    AccountNumber  string
    // Set account sequence of address
    Sequence       string
    // Broadcast mode (sync,async,block)
    BroadcastMode  string	
    // Transaction gas limit
    GasLimit       string
    // Transaction gas price
    GasPrice       string
    // Transaction gas limit adjustment
    GasAdjustment  string
    // Transaction fee amount
    FeeAmount      string
    // Transaction sign mode
    SignMode       signing.SignMode
    // Set fee granter of transaction builder
    FeeGranter     sdk.AccAddress
    // Set timeout height of transaction builder
    TimeoutHeight  string
    // LCD URL
    LcdURL         string
    // GRPC URL
    GrpcURL        string
    // Tendermint RPC URL
    RpcURL         string
    // Ethereum VM RPC URL
    EvmRpcURL      string
    // Set user want pagination option
    Pagination     types.Pagination
    // Set output document name when created transaction with json file
    // "Generate only" is same that OutputDocument is not empty string 
    OutputDocument string	
}
```

## Handle transactions
### Create and sign tx
```go
// Create signed transaction by using msg.
// e.g. Send coin of bank module
bankSendMsg := types.BankSendMsg {
    FromAddress: "xpla1g8ku0mt75j4p8luxzku6dkcxxvnc0tt352z0k9",
    ToAddress: "xpla1j3dtjvchp7ec3nnn6357jv8v8f29akx6p2u78g",
    Amount: "1000",
}

txbytes, err := xplac.BankSend(bankSendMsg).CreateAndSignTx()
```

### Create unsigned tx
```go
// Create unsigned transaction by using msg.
// e.g. Send coin of bank module
bankSendMsg := types.BankSendMsg {
    FromAddress: "xpla1g8ku0mt75j4p8luxzku6dkcxxvnc0tt352z0k9",
    ToAddress: "xpla1j3dtjvchp7ec3nnn6357jv8v8f29akx6p2u78g",
    Amount: "1000axpla",
}

txbytes, err := xplac.BankSend(bankSendMsg).CreateUnsignedTx()
```

### Sign tx
```go
addr, err := key.Bech32AddrString(priKey)
// Sign transaction with local transaction file.
signTxMsg := types.SignTxMsg{
    FileName:    "./unsignedTx.json",    
    FromAddress: addr,
}	

txbytes, err := xplac.SignTx(signTxMsg)
```

### Multisign tx
```go
// Multi sign transaction with local transaction file.
// It is able to sign when local keyring file and signature file exist.
txMultiSignMsg := types.TxMultiSignMsg{
    FileName: "./unsignedTx.json",
    GenerateOnly: true,
    FromName: "mykey",
    Offline: true,
    SignatureFiles: []string{"signatureFiles.json", ...},	
}

res, err := xplac.MultiSign(txMultiSignMsg)
```

### Encode tx
```go
// Encoding transaction by using base64
encodeTxMsg := types.EncodeTxMsg {
    FileName: "./unsignedTx.json",
}

res, err := xplac.EncodeTx(encodeTxMsg)
```

### Decode tx
```go
// Decoding transaction
decodeTxMsg := types.DecodeTxMsg{
    EncodedByteString: "CvwBCvkBCiUvY29zbW9zLmdvdi52MWJldGE......",
}

res, err := xplac.DecodeTx(decodeTxMsg)
```

### Validate signatures of tx
```go
validateSignaturesMsg := types.ValidateSignaturesMsg{
    FileName: "./signedTx.json",
    Offline: true,
}

res, err := xplac.ValidateSignatures(validateSignaturesMsg)
```