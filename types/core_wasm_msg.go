package types

type StoreMsg struct {
	FilePath              string
	InstantiatePermission string
}

type InstantiateMsg struct {
	CodeId  string
	Amount  string
	Label   string
	InitMsg string
	Admin   string
	NoAdmin string
}

type ExecuteMsg struct {
	ContractAddress string
	Amount          string
	ExecMsg         string
}

type ClearContractAdminMsg struct {
	ContractAddress string
}

type SetContractAdminMsg struct {
	NewAdmin        string
	ContractAddress string
}

type MigrateMsg struct {
	ContractAddress string
	CodeId          string
	MigrateMsg      string
}

type QueryMsg struct {
	ContractAddress string
	QueryMsg        string
}

type ListContractByCodeMsg struct {
	CodeId string
}

type DownloadMsg struct {
	CodeId           string
	DownloadFileName string
}

type CodeInfoMsg struct {
	CodeId string
}

type ContractInfoMsg struct {
	ContractAddress string
}

type ContractStateAllMsg struct {
	ContractAddress string
}

type ContractHistoryMsg struct {
	ContractAddress string
}
