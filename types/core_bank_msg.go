package types

type BankSendMsg struct {
	FromAddress string
	ToAddress   string
	Amount      string
}

type BankBalancesMsg struct {
	Address string
	Denom   string
}

type DenomMetadataMsg struct {
	Denom string
}

type TotalMsg struct {
	Denom string
}
