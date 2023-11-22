package authz

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - authz grant
func MakeAuthzGrantMsg(authzGrantMsg types.AuthzGrantMsg, granter sdk.AccAddress) (authz.MsgGrant, error) {
	return parseAuthzGrantArgs(authzGrantMsg, granter)
}

// (Tx) make msg - revoke
func MakeAuthzRevokeMsg(authzRevokeMsg types.AuthzRevokeMsg, granter sdk.AccAddress) (authz.MsgRevoke, error) {
	return parseAuthzRevokeArgs(authzRevokeMsg, granter)
}

// (Tx) make msg - authz execute
func MakeAuthzExecMsg(authzExecMsg types.AuthzExecMsg, encodingConfig params.EncodingConfig) (authz.MsgExec, error) {
	return parseAuthzExecArgs(authzExecMsg, encodingConfig)
}

// (Query) make msg - authz grants
func MakeQueryAuthzGrantsMsg(queryAuthzGrantMsg types.QueryAuthzGrantMsg) (authz.QueryGrantsRequest, error) {
	return parseQueryAuthzGrantsArgs(queryAuthzGrantMsg)
}

// (Query) make msg - authz grants by grantee
func MakeQueryAuthzGrantsByGranteeMsg(queryAuthzGrantMsg types.QueryAuthzGrantMsg) (authz.QueryGranteeGrantsRequest, error) {
	return parseQueryAuthzGrantsByGranteeArgs(queryAuthzGrantMsg)
}

// (Query) make msg - authz grants by granter
func MakeQueryAuthzGrantsByGranterMsg(queryAuthzGrantMsg types.QueryAuthzGrantMsg) (authz.QueryGranterGrantsRequest, error) {
	return parseQueryAuthzGrantsByGranterArgs(queryAuthzGrantMsg)
}
