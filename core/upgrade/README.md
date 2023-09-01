# Upgrade module
## Usage
### (Tx) Proposal software upgrade
```go
softwareUpgradeMsg := types.SoftwareUpgradeMsg{
    UpgradeName: "Upgrade Name",
    Title: "Upgrade Title",
    Description: "Upgrade Description",
    UpgradeHeight:"6000",
    UpgradeInfo: `{"upgrade_info":"INFO"}`,
    Deposit: "1000",
}

txbytes, err := xplac.SoftwareUpgrade(softwareUpgradeMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Proposal cancel software upgrade
```go
cancelSoftwareUpgradeMsg := types.CancelSoftwareUpgradeMsg {
    Title: "Cancel software upgrade",
    Description: "Cancel software upgrade description",
    Deposit: "1000",
}

txbytes, err := xplac.CancelSoftwareUpgrade(cancelSoftwareUpgradeMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) Upgrade applied
```go
appliedMsg := types.AppliedMsg{
    UpgradeName: "upgrade-name",
}

res, err := xplac.UpgradeApplied(appliedMsg).Query()
```

### (Query) Modules version
```go
// Query specific module name
queryModulesVersionMsg := types.QueryModulesVersionMsg{
    ModuleName: "auth",
}

res, err := xplac.ModulesVersion(queryModulesVersionMsg).Query()

// Query all modules version
res, err := xplac.ModulesVersion().Query()
```

### (Query) Upgrade plan
```go
res, err := xplac.Plan().Query()
```