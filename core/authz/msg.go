package authz

import (
	"github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - authz grant
func MakeAuthzGrantMsg(authzGrantMsg types.AuthzGrantMsg, privKey key.PrivateKey) (authz.MsgGrant, error) {
	return parseAuthzGrantArgs(authzGrantMsg, privKey)
}

// (Tx) make msg - revoke
func MakeAuthzRevokeMsg(authzRevokeMsg types.AuthzRevokeMsg, privKey key.PrivateKey) (authz.MsgRevoke, error) {
	return parseAuthzRevokeArgs(authzRevokeMsg, privKey)
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
