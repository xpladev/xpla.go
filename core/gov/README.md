# Gov module
## Usage
### (Tx) Submit proposal
```go
submitProposalMsg := types.SubmitProposalMsg{
    Title: "Test proposal",
    Description: "Proposal description",
    Type: "text",
    Deposit: "1000",
}
txbytes, err := xplac.SubmitProposal(submitProposalMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Deposit
```go
govDepositMsg := types.GovDepositMsg {
    ProposalID: "1",
    Deposit: "1000axpla",
}

txbytes, err := xplac.GovDeposit(govDepositMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Vote
```go
voteMsg := types.VoteMsg{
    ProposalID: "1",
    Option: "yes",
}

txbytes, err := xplac.Vote(voteMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Weighted vote
```go
weightedVoteMsg := types.WeightedVoteMsg{
    ProposalID: "1",
    Yes: "0.6",
    No: "0.3",
    Abstain: "0.05",
    NoWithVeto: "0.05",
}

txbytes, err := xplac.WeightedVote(weightedVoteMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```
### (Query) proposal
```go
queryProposalMsg := types.QueryProposalMsg{
    ProposalID: "1",
}

res, err := xplac.QueryProposal(queryProposalMsg).Query()
```

### (Query) proposals
```go
// Query all proposals
queryProposalsMsg := types.QueryProposalsMsg{}

// Set options
queryProposalsMsg := types.QueryProposalsMsg{
    Status: "DepositPeriod",
    Voter: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    Depositor: "xpla19yq7kjcgse7x672faptju0lxmy4cvdlcpmxnyn",
}    

res, err := xplac.QueryProposals(queryProposalsMsg).Query()
```

### (Query) deposit
```go
// Query details of a deposit
queryDepositMsg := types.QueryDepositMsg{
    Depositor: "xpla1e4f6k98es55vxxv2pcfzpsjrf3mvazeyqpw8g9",
    ProposalID: "1",
}

// Query depostis on a proposal
queryDepositMsg := types.QueryDepositMsg{
    ProposalID: "1",
}

res, err := xplac.QueryDeposit(queryDepositMsg).Query()
```

### (Query) Tally
```go
tallyMsg := types.TallyMsg{
    ProposalID: "1",
}

res, err := xplac.Tally(tallyMsg).Query()
```

### (Query) gov params
```go
// Query all parameters of the gov process
res, err := xplac.GovParams().Query()

// Query the parameters (voting|tallying|deposit) of the governance process
govParamsMsg := types.GovParamsMsg{
    ParamType: "tallying",
}
res, err := xplac.GovParams(govParamsMsg).Query()
```

### (Query) proposer
```go
proposerMsg := types.ProposerMsg {
    ProposalID: "1",
}

res, err := xplac.Proposer(proposerMsg).Query()
```

### (Query) Votes
```go
// query details of a single vote
queryVoteMsg := types.QueryVoteMsg{
    ProposalID: "1",
    VoterAddr: "xpla1mss8agu59x3qte67lfn80x6w65thq4283dvn7l",
}

// Query votes on a proposal
queryVoteMsg := types.QueryVoteMsg{
    ProposalID: "1",
}

res, err := xplac.QueryVote(queryVoteMsg).Query()
```