package auth

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

type AuthExternal struct {
	Xplac provider.XplaClient
}

func NewAuthExternal(xplac provider.XplaClient) (e AuthExternal) {
	e.Xplac = xplac
	return e
}

// Query

// Query the current auth parameters.
func (e AuthExternal) AuthParams() provider.XplaClient {
	msg, err := MakeAuthParamMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(AuthModule).
		WithMsgType(AuthQueryParamsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query for account by address.
func (e AuthExternal) AccAddress(queryAccAddresMsg types.QueryAccAddressMsg) provider.XplaClient {
	msg, err := MakeQueryAccAddressMsg(queryAccAddresMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(AuthModule).
		WithMsgType(AuthQueryAccAddressMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query all accounts.
func (e AuthExternal) Accounts() provider.XplaClient {
	msg, err := MakeQueryAccountsMsg()
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(AuthModule).
		WithMsgType(AuthQueryAccountsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query for paginated transactions that match a set of events.
func (e AuthExternal) TxsByEvents(txsByEventsMsg types.QueryTxsByEventsMsg) provider.XplaClient {
	msg, err := MakeTxsByEventsMsg(txsByEventsMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(AuthModule).
		WithMsgType(AuthQueryTxsByEventsMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query for a transaction by hash <addr>/<seq> combination or comma-separated signatures in a committed block.
func (e AuthExternal) Tx(queryTxMsg types.QueryTxMsg) provider.XplaClient {
	msg, err := MakeQueryTxMsg(queryTxMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(AuthModule).
		WithMsgType(AuthQueryTxMsgType).
		WithMsg(msg)
	return e.Xplac
}
