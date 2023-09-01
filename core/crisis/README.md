# Crisis module
## Usage
### (Tx) Invariant broken
```go
invariantBrokenMsg := types.InvariantBrokenMsg{
    ModuleName: "module_name",
    InvariantRoute: "invariant_route",
}

txbytes, err := xplac.InvariantBroken(invariantBrokenMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```