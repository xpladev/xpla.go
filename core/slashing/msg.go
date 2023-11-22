package slashing

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - unjail
func MakeUnjailMsg(addr sdk.AccAddress) (slashingtypes.MsgUnjail, error) {
	return parseUnjailArgs(addr)
}

// (Query) make msg - slahsing params
func MakeQuerySlashingParamsMsg() (slashingtypes.QueryParamsRequest, error) {
	return slashingtypes.QueryParamsRequest{}, nil
}

// (Query) make msg - signing infos
func MakeQuerySigningInfosMsg() (slashingtypes.QuerySigningInfosRequest, error) {
	return slashingtypes.QuerySigningInfosRequest{
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - signing info
func MakeQuerySigningInfoMsg(signingInfoMsg types.SigningInfoMsg, xplacEncodingConfig params.EncodingConfig) (slashingtypes.QuerySigningInfoRequest, error) {
	return parseQuerySigingInfoArgs(signingInfoMsg, xplacEncodingConfig)
}
