package feegrant

import (
	"github.com/xpladev/xpla.go/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
)

// (Tx) make msg - fee grant
func MakeFeeGrantMsg(feeGrantMsg types.FeeGrantMsg, granter sdk.AccAddress) (feegrant.MsgGrantAllowance, error) {
	return parseFeeGrantArgs(feeGrantMsg, granter)
}

// (Tx) make msg - fee grant revoke
func MakeRevokeFeeGrantMsg(revokeFeeGrantMsg types.RevokeFeeGrantMsg, granter sdk.AccAddress) (feegrant.MsgRevokeAllowance, error) {
	return parseRevokeFeeGrantArgs(revokeFeeGrantMsg, granter)
}

// (Query) make msg - query fee grants
func MakeQueryFeeGrantMsg(queryFeeGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowanceRequest, error) {
	return parseQueryFeeGrantArgs(queryFeeGrantMsg)
}

// (Query) make msg - fee grants by grantee
func MakeQueryFeeGrantsByGranteeMsg(queryFeeGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowancesRequest, error) {
	return parseQueryFeeGrantsByGranteeArgs(queryFeeGrantMsg)
}

// (Query) make msg - fee grants by granter
func MakeQueryFeeGrantsByGranterMsg(queryFeeGrantMsg types.QueryFeeGrantMsg) (feegrant.QueryAllowancesByGranterRequest, error) {
	return parseQueryFeeGrantsByGranterArgs(queryFeeGrantMsg)
}
