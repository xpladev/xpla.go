package types

type FundCommunityPoolMsg struct {
	Amount string
}

type CommunityPoolSpendMsg struct {
	Title        string
	Description  string
	Recipient    string
	Amount       string
	Deposit      string
	JsonFilePath string
}

type WithdrawRewardsMsg struct {
	DelegatorAddr string
	ValidatorAddr string
	Commission    bool
}

type SetWithdrawAddrMsg struct {
	WithdrawAddr string
}

type ValidatorOutstandingRewardsMsg struct {
	ValidatorAddr string
}

type QueryDistCommissionMsg struct {
	ValidatorAddr string
}

type QueryDistSlashesMsg struct {
	ValidatorAddr string
	StartHeight   string
	EndHeight     string
}

type QueryDistRewardsMsg struct {
	ValidatorAddr string
	DelegatorAddr string
}
