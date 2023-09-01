package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	DefaultTokens                  = sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)
	DefaultAmount                  = DefaultTokens.String() + XplaDenom
	DefaultCommissionRate          = "0.1"
	DefaultCommissionMaxRate       = "0.2"
	DefaultCommissionMaxChangeRate = "0.01"
	DefaultMinSelfDelegation       = "1"
	DefaultHomeDir                 = ".xpla"
	Memo                           string
)

type CreateValidatorMsg struct {
	NodeKey                 string
	PrivValidatorKey        string
	ValidatorAddress        string
	HomeDir                 string
	Website                 string
	SecurityContact         string
	Identity                string
	Moniker                 string
	Amount                  string
	Details                 string
	CommissionRate          string
	CommissionMaxRate       string
	CommissionMaxChangeRate string
	MinSelfDelegation       string
	ServerIp                string
}

type EditValidatorMsg struct {
	Website           string
	SecurityContact   string
	Identity          string
	Details           string
	Moniker           string
	CommissionRate    string
	MinSelfDelegation string
}

type DelegateMsg struct {
	Amount  string
	ValAddr string
}

type RedelegateMsg struct {
	Amount     string
	ValSrcAddr string
	ValDstAddr string
}

type UnbondMsg struct {
	Amount  string
	ValAddr string
}

type QueryValidatorMsg struct {
	ValidatorAddr string
}

type QueryDelegationMsg struct {
	DelegatorAddr string
	ValidatorAddr string
}

type QueryUnbondingDelegationMsg struct {
	DelegatorAddr string
	ValidatorAddr string
}

type QueryRedelegationMsg struct {
	DelegatorAddr    string
	SrcValidatorAddr string
	DstValidatorAddr string
}

type HistoricalInfoMsg struct {
	Height string
}
