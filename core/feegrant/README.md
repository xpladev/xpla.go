# Feegrant module
## Usage
### (Tx) Grant
```go
grantMsg := types.GrantMsg {
    Granter: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    Grantee: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
    SpendLimit: "1000",
    
    // select options as below
    Period: "3600",
    PeriodLimit: "10",
    Expiration: "2100-01-01T23:59:59+00:00",
}
txbytes, err := xplac.Grant(grantMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Revoke grant
```go
revokeGrantMsg := types.RevokeGrantMsg{
    Granter: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    Grantee: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
}
txbytes, err := xplac.RevokeGrant(revokeGrantMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) grants
```go
// Query details of single grant
queryGrantMsg := types.QueryGrantMsg{
    Granter: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    Grantee: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
}

// Query all grants of a grantee
queryGrantMsg := types.QueryGrantMsg{
    Grantee: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
}

// Query all grants of a granter
queryGrantMsg := types.QueryGrantMsg{
    Granter: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
}

res, err := xplac.QueryGrant(queryGrantMsg).Query()
```