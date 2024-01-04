package auth

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &AuthExternal{}

type AuthExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e AuthExternal) {
	e.Xplac = xplac
	e.Name = AuthModule
	return e
}

func (e AuthExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e AuthExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Query

// Query the current auth parameters.
func (e AuthExternal) AuthParams() provider.XplaClient {
	msg, err := MakeAuthParamMsg()
	if err != nil {
		return e.Err(AuthQueryParamsMsgType, err)
	}

	return e.ToExternal(AuthQueryParamsMsgType, msg)
}

// Query for account by address.
func (e AuthExternal) AccAddress(queryAccAddresMsg types.QueryAccAddressMsg) provider.XplaClient {
	msg, err := MakeQueryAccAddressMsg(queryAccAddresMsg)
	if err != nil {
		return e.Err(AuthQueryAccAddressMsgType, err)
	}

	return e.ToExternal(AuthQueryAccAddressMsgType, msg)
}

// Query all accounts.
func (e AuthExternal) Accounts() provider.XplaClient {
	msg, err := MakeQueryAccountsMsg()
	if err != nil {
		return e.Err(AuthQueryAccountsMsgType, err)
	}

	return e.ToExternal(AuthQueryAccountsMsgType, msg)
}

// Query for paginated transactions that match a set of events.
func (e AuthExternal) TxsByEvents(txsByEventsMsg types.QueryTxsByEventsMsg) provider.XplaClient {
	msg, err := MakeTxsByEventsMsg(txsByEventsMsg)
	if err != nil {
		return e.Err(AuthQueryTxsByEventsMsgType, err)
	}

	return e.ToExternal(AuthQueryTxsByEventsMsgType, msg)
}

// Query for a transaction by hash <addr>/<seq> combination or comma-separated signatures in a committed block.
func (e AuthExternal) Tx(queryTxMsg types.QueryTxMsg) provider.XplaClient {
	msg, err := MakeQueryTxMsg(queryTxMsg)
	if err != nil {
		return e.Err(AuthQueryTxMsgType, err)
	}

	return e.ToExternal(AuthQueryTxMsgType, msg)
}
