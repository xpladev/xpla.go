package bank

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
)

var _ core.External = &BankExternal{}

type BankExternal struct {
	Xplac provider.XplaClient
	Name  string
}

func NewExternal(xplac provider.XplaClient) (e BankExternal) {
	e.Xplac = xplac
	e.Name = BankModule
	return e
}

func (e BankExternal) ToExternal(msgType string, msg interface{}) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithModule(e.Name).
		WithMsgType(msgType).
		WithMsg(msg)
}

func (e BankExternal) Err(msgType string, err error) provider.XplaClient {
	return provider.ResetModuleAndMsgXplac(e.Xplac).
		WithErr(
			e.Xplac.GetLogger().Err(err,
				types.LogMsg("module", e.Name),
				types.LogMsg("msg", msgType)),
		)
}

// Tx

// Send funds from one account to another.
func (e BankExternal) BankSend(bankSendMsg types.BankSendMsg) provider.XplaClient {
	msg, err := MakeBankSendMsg(bankSendMsg)
	if err != nil {
		return e.Err(BankSendMsgType, err)
	}

	return e.ToExternal(BankSendMsgType, msg)
}

// Query

// Query for account balances by address
func (e BankExternal) BankBalances(bankBalancesMsg types.BankBalancesMsg) provider.XplaClient {
	switch {
	case bankBalancesMsg.Denom == "":
		msg, err := MakeBankAllBalancesMsg(bankBalancesMsg)
		if err != nil {
			return e.Err(BankAllBalancesMsgType, err)
		}

		return e.ToExternal(BankAllBalancesMsgType, msg)

	default:
		msg, err := MakeBankBalanceMsg(bankBalancesMsg)
		if err != nil {
			return e.Err(BankBalanceMsgType, err)
		}

		return e.ToExternal(BankBalanceMsgType, msg)
	}
}

// Query the client metadata for coin denominations.
func (e BankExternal) DenomMetadata(denomMetadataMsg ...types.DenomMetadataMsg) provider.XplaClient {
	switch {
	case len(denomMetadataMsg) == 0:
		msg, err := MakeDenomsMetaDataMsg()
		if err != nil {
			return e.Err(BankDenomsMetadataMsgType, err)
		}

		return e.ToExternal(BankDenomsMetadataMsgType, msg)

	case len(denomMetadataMsg) == 1:
		msg, err := MakeDenomMetaDataMsg(denomMetadataMsg[0])
		if err != nil {
			return e.Err(BankDenomMetadataMsgType, err)
		}

		return e.ToExternal(BankDenomMetadataMsgType, msg)

	default:
		return e.Err(BankDenomMetadataMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}

// Query the total supply of coins of the chain.
func (e BankExternal) Total(totalMsg ...types.TotalMsg) provider.XplaClient {
	if len(totalMsg) == 0 {
		msg, err := MakeTotalSupplyMsg()
		if err != nil {
			return e.Err(BankTotalMsgType, err)
		}

		return e.ToExternal(BankTotalMsgType, msg)

	} else if len(totalMsg) == 1 {
		msg, err := MakeSupplyOfMsg(totalMsg[0])
		if err != nil {
			return e.Err(BankTotalSupplyOfMsgType, err)
		}

		return e.ToExternal(BankTotalSupplyOfMsgType, msg)

	} else {
		return e.Err(BankTotalSupplyOfMsgType, types.ErrWrap(types.ErrInvalidRequest, "need only one parameter"))
	}
}
