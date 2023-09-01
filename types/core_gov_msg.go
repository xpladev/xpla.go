package types

type SubmitProposalMsg struct {
	Title       string
	Description string
	Type        string
	Deposit     string
}

type GovDepositMsg struct {
	ProposalID string
	Deposit    string
}

type VoteMsg struct {
	ProposalID string
	Option     string
}

type WeightedVoteMsg struct {
	ProposalID string
	Yes        string
	No         string
	Abstain    string
	NoWithVeto string
}

type QueryProposalMsg struct {
	ProposalID string
}

type QueryProposalsMsg struct {
	Depositor string
	Voter     string
	Status    string
}

type QueryDepositMsg struct {
	ProposalID string
	Depositor  string
}

type TallyMsg struct {
	ProposalID string
}

type GovParamsMsg struct {
	ParamType string
}

type ProposerMsg struct {
	ProposalID string
}

type QueryVoteMsg struct {
	ProposalID string
	VoterAddr  string
}
