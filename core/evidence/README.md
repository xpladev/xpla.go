# Evidence module
## Usage
### (Query) Evidence 
```go
queryEvidenceMsg := types.QueryEvidenceMsg{
    Hash: "DF0C23E8634E480F84B9D5674A7CDC9816466DEC28A3358F73260F68D28D7660"
}

res, err := xplac.QueryEvidence(queryEvidenceMsg).Query()
```

### (Query) All evidence
```go
res, err := xplac.QueryEvidence().Query()
```
