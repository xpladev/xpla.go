package types

type SoftwareUpgradeMsg struct {
	UpgradeName   string
	Title         string
	Description   string
	UpgradeHeight string
	UpgradeInfo   string
	Deposit       string
}

type CancelSoftwareUpgradeMsg struct {
	Title       string
	Description string
	Deposit     string
}

type AppliedMsg struct {
	UpgradeName string
}

type QueryModulesVersionMsg struct {
	ModuleName string
}
