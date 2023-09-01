package types

type AuthzGrantMsg struct {
	Grantee           string
	Granter           string
	SpendLimit        string
	Expiration        string
	MsgType           string
	AllowValidators   []string
	DenyValidators    []string
	AuthorizationType string
}

type AuthzRevokeMsg struct {
	Grantee string
	Granter string
	MsgType string
}

type AuthzExecMsg struct {
	Grantee      string
	ExecFile     string
	ExecTxString string
}

type QueryAuthzGrantMsg struct {
	Grantee string
	Granter string
	MsgType string
}
