package types

type FeeGrantMsg struct {
	Grantee     string
	Granter     string
	SpendLimit  string
	Expiration  string
	Period      string
	PeriodLimit string
	AllowedMsg  []string
}

type RevokeFeeGrantMsg struct {
	Grantee string
	Granter string
}

type QueryFeeGrantMsg struct {
	Grantee string
	Granter string
}
