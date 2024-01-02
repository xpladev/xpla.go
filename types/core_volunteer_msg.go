package types

type RegisterVolunteerValidatorMsg struct {
	Title        string
	Description  string
	Deposit      string
	Amount       string
	ValPubKey    string
	Moniker      string
	Identity     string
	Website      string
	Security     string
	Details      string
	JsonFilePath string
}

type UnregisterVolunteerValidatorMsg struct {
	Title        string
	Description  string
	Deposit      string
	ValAddress   string
	JsonFilePath string
}
