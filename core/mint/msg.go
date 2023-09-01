package mint

import (
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// (Query) make msg - mint params
func MakeQueryMintParamsMsg() (minttypes.QueryParamsRequest, error) {
	return minttypes.QueryParamsRequest{}, nil
}

// (Query) make msg - inflation
func MakeQueryInflationMsg() (minttypes.QueryInflationRequest, error) {
	return minttypes.QueryInflationRequest{}, nil
}

// (Query) make msg - annual provisions
func MakeQueryAnnualProvisionsMsg() (minttypes.QueryAnnualProvisionsRequest, error) {
	return minttypes.QueryAnnualProvisionsRequest{}, nil
}
