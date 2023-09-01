package types

type SubspaceMsg struct {
	Subspace string
	Key      string
}

type ParamChangeMsg struct {
	Title        string
	Description  string
	Changes      []string
	Deposit      string
	JsonFilePath string
}
