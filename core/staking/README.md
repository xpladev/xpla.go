# Staking module
## Usage
### (Tx) Create validator
```go
// Create validator using json file in the home directory
createValidatorMsg := types.CreateValidatorMsg{	
    ValidatorAddress: "xplavaloper10gv4zj9633v6cje6s2sc0a0xl52hjr6f9jp0q7",
    Website: "website",
    HomeDir: "/ABSPATH/.xpla",
    SecurityContact: "contact point",
    Identity: "identity",
    NodeID: "nodeid",
    ChainID: "chainid",
    Moniker: "moniker",
    Details: "details",
    Amount: "100000",
}

// Create validator using string values
createValidatorMsg := types.CreateValidatorMsg{
    NodeKey:                 NodeKey,
    PrivValidatorKey:        PrivValidatorKey,
    ValidatorAddress:        "xplavaloper10gv4zj9633v6cje6s2sc0a0xl52hjr6f9jp0q7",
    Moniker:                 "moniker",
    Identity:                "identity",
    Website:                 "website",
    SecurityContact:         "securityContact",
    Details:                 "details",
    Amount:                  "amount",
    CommissionRate:          "commissionRate",
    CommissionMaxRate:       "commissionMaxRate",
    CommissionMaxChangeRate: "commissionMaxChangeRate",
    MinSelfDelegation:       "minSelfDelegation",
}

txbytes, err := xplac.CreateValidator(createValidatorMsg).CreateAndSignTx()
```
### (Tx) Edit validator
```go
editValidatorMsg := types.EditValidatorMsg{		
    Website: "website",
    SecurityContact: "securitycontact",
    Identity: "identity",    
    Details: "details",
    Moniker: "moniker",
    CommissionRate: "commissionRate",	
    MinSelfDelegation: "minSelfDelegation",
}
txbytes, err := xplac.EditValidator(editValidatorMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```
### (Tx) Delegate
```go
delegateMsg := types.DelegateMsg{
    Amount: "1000axpla",
    ValAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}
txbytes, err := xplac.Delegate(delegateMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Unbond
```go
unbondMsg := types.UnbondMsg{
    Amount: "100axpla",
    ValAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}
txbytes, err := xplac.Unbond(unbondMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Redelegate
```go
redelegateMsg := types.RedelegateMsg{
    Amount: "100axpla",
    ValSrcAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
    ValDstAddr: "xplavaloper1r7tdqs8zgtkty2u06yp5nw6dc6c9hvz9ak98r5",
}
txbytes, err := xplac.Redelegate(redelegateMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) validators
```go
// Query all validators
res, err := xplac.QueryValidators().Query()

// Query validator by retrieving validator address
queryValidatorMsg := types.QueryValidatorMsg {
    ValidatorAddr: "xplavaloper13trl452wgle9qxpxhse9605k9x0399cm85pf34",
}
res, err := xplac.QueryValidators(queryValidatorMsg).Query()
```
### (Query) delegation
```go
// Query a delegation based on address and validator address 
queryDelegationMsg := types.QueryDelegationMsg {
    DelegatorAddr: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    ValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}

// Query all delegations made by one delegator 
queryDelegationMsg := types.QueryDelegationMsg {
    DelegatorAddr: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
}

// Query all delegations made to one validator
queryDelegationMsg := types.QueryDelegationMsg {
    ValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}

res, err := xplac.QueryDelegation(queryDelegationMsg).Query()
```

### (Query) unbonding delegation
```go
// Query an unbonding-delegation record based on delegator and validator address
queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg {
    DelegatorAddr: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    ValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}

// Query all unbonding-delegations records for one delegator
queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg {
    DelegatorAddr: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
}

// Query all unbonding delegatations from a validator
queryUnbondingDelegationMsg := types.QueryUnbondingDelegationMsg {    
    ValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}

res, err := xplac.QueryUnbondingDelegation(queryUnbondingDelegationMsg).Query()
```
### (Query) Redelegations
```go
// Query a redelegation record based on delegator and a source and destination validator
queryRedelegationMsg := types.QueryRedelegationMsg{
    DelegatorAddr: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",    	
    SrcValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
    DstValidatorAddr: "xplavaloper1r7tdqs8zgtkty2u06yp5nw6dc6c9hvz9ak98r5",
}

// Query all redelegations records for one delegator
queryRedelegationMsg := types.QueryRedelegationMsg{
    DelegatorAddr: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",    		
}

// Query all outgoing redelegations from a validator
queryRedelegationMsg := types.QueryRedelegationMsg{
    SrcValidatorAddr: "xplavaloper19yq7kjcgse7x672faptju0lxmy4cvdlcsx9ftw",
}

res, err := xplac.QueryRedelegation(queryRedelegationMsg).Query()
```

### (Query) Historical info
```go
historicalInfoMsg := types.HistoricalInfoMsg{
    Height: "2000",
}

res, err := xplac.HistoricalInfo(historicalInfoMsg).Query()
```

### (Query) staking pool
```go
res, err := xplac.StakingPool().Query()
```

### (Query) staking params
```go
res, err := xplac.StakingParams().Query()
```
