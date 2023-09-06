package types

type EncodeTxMsg struct {
	FileName string
}

type DecodeTxMsg struct {
	EncodedByteString string
}

type ValidateSignaturesMsg struct {
	FileName string
	ChainID  string
	Offline  bool
}

type SignTxMsg struct {
	UnsignedFileName string
	SignatureOnly    bool
	MultisigAddress  string
	Overwrite        bool
	Amino            bool
	Offline          bool
}

type TxMultiSignMsg struct {
	FileName       string
	GenerateOnly   bool
	FromName       string
	Offline        bool
	SignatureFiles []string
	SignatureOnly  bool
	Amino          bool
	KeyringPath    string
	KeyringBackend string
}

type QueryAccAddressMsg struct {
	Address string
}

type QueryTxsByEventsMsg struct {
	Events string
	Page   string
	Limit  string
}

type QueryTxMsg struct {
	Value string
	Type  string
}
