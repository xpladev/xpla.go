# Authz module
## Usage
### (Tx) Authz grant
```go
authzGrantMsg := types.AuthzGrantMsg{
    Granter: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    Grantee: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
    AuthorizationType: "send",
    SpendLimit: "1000axpla",		
}

txbytes, err := xplac.AuthzGrant(authzGrantMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```
### (Tx) Authz revoke
```go
authzRevokeMsg := types.AuthzRevokeMsg{
    Granter: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    Grantee: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
    MsgType: "/cosmos.bank.v1beta1.MsgSend",
}

txbytes, err := xplac.AuthzRevoke(authzRevokeMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Authz exec
```go
// Execute by using file
authzExecMsg := types.AuthzExecMsg{
    Grantee: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
    ExecFile: "execFile.json",
}

// Execute by using transaction json string
authzExecMsg := types.AuthzExecMsg{
    Grantee: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
    ExecTxString: `{TRANSACTION_JSON}`,
}

txbytes, err := xplac.AuthzExec(authzExecMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) Authz grants
```go
// Query grants for a granter-grantee pair and optionally a msg-type-url
queryAuthzGrantMsg := types.QueryAuthzGrantMsg {
    Granter: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    Grantee: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
}

// Query authorization grants granted by granter
queryAuthzGrantMsg := types.QueryAuthzGrantMsg {
    Granter: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
}

// Query authorization grants granted to a grantee
queryAuthzGrantMsg := types.QueryAuthzGrantMsg {
    Grantee: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
}

res, err := xplac.QueryAuthzGrants(queryAuthzGrantMsg).Query()
```