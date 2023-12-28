# Volunteer module
## Usage
### (Tx) Register volunteer validator
```go
// Register volunteer validator using json file
registerVolunteerValidatorMsg := types.RegisterVolunteerValidatorMsg{	
    ValPubKey: `{"@type": "/cosmos.crypto.ed25519.PubKey","key": "2z2yttKfEsLQyQnHYdgKEuky9zB3gscxapn9IyexxWk="}`,
    Amount: "100000axpla",
    Moniker: "moniker",
    Identity: "identity",
    Website: "website",
    Security: "contact point",
    Details: "details",
    JsonFilePath: "/JSON/FILE/PATH/FOR/PROPOSAL/TO/REGISTER"
}

// Register volunteer validator using string values
registerVolunteerValidatorMsg := types.RegisterVolunteerValidatorMsg{
    Title: "Register volunteer validator",
    Description: "description",
    Deposit: "1000000axpla",
    ValPubKey: `{"@type": "/cosmos.crypto.ed25519.PubKey","key": "2z2yttKfEsLQyQnHYdgKEuky9zB3gscxapn9IyexxWk="}`,
    Amount: "100000axpla",
    Moniker: "moniker",
    Identity: "identity",
    Website: "website",
    Security: "contact point",
    Details: "details",
}

txbytes, err := xplac.RegisterVolunteerValidator(registerVolunteerValidatorMsg).CreateAndSignTx()
```

### (Tx) Unregister volunteer validator
```go
// Unregister volunteer validator using json file
unregisterVolunteerValidatorMsg := types.UnregisterVolunteerValidatorMsg{	
    JsonFilePath: "/JSON/FILE/PATH/FOR/PROPOSAL/TO/UNREGISTER"
}

// Unegister volunteer validator using string values
unregisterVolunteerValidatorMsg := types.UnregisterVolunteerValidatorMsg{
    Title: "Unregister volunteer validator",
    Description: "description",
    Deposit: "1000000",
    ValAddress: "xplavaloper1hggt7sgsegcg3daz0rpa9m8mkmx3qyarse9utx",
}

txbytes, err := xplac.UnregisterVolunteerValidator(unregisterVolunteerValidatorMsg).CreateAndSignTx()
```

### (Query) Volunter validators
```go
// Query volunteer validators
res, err := xplac.QueryVolunteerValidators().Query()
```