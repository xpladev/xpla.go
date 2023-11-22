package bank

import (
	"github.com/xpladev/xpla.go/provider"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

type BankExternal struct {
	Xplac provider.XplaClient
}

func NewBankExternal(xplac provider.XplaClient) (e BankExternal) {
	e.Xplac = xplac
	return e
}

// Tx

// Send funds from one account to another.
func (e BankExternal) BankSend(bankSendMsg types.BankSendMsg) provider.XplaClient {
	msg, err := MakeBankSendMsg(bankSendMsg)
	if err != nil {
		return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
	}
	e.Xplac.WithModule(BankModule).
		WithMsgType(BankSendMsgType).
		WithMsg(msg)
	return e.Xplac
}

// Query

// Query for account balances by address
func (e BankExternal) BankBalances(bankBalancesMsg types.BankBalancesMsg) provider.XplaClient {
	if bankBalancesMsg.Denom == "" {
		msg, err := MakeBankAllBalancesMsg(bankBalancesMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(BankModule).
			WithMsgType(BankAllBalancesMsgType).
			WithMsg(msg)
	} else {
		msg, err := MakeBankBalanceMsg(bankBalancesMsg)
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(BankModule).
			WithMsgType(BankBalanceMsgType).
			WithMsg(msg)
	}
	return e.Xplac

}

// Query the client metadata for coin denominations.
func (e BankExternal) DenomMetadata(denomMetadataMsg ...types.DenomMetadataMsg) provider.XplaClient {
	if len(denomMetadataMsg) == 0 {
		msg, err := MakeDenomsMetaDataMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(BankModule).
			WithMsgType(BankDenomsMetadataMsgType).
			WithMsg(msg)
	} else if len(denomMetadataMsg) == 1 {
		msg, err := MakeDenomMetaDataMsg(denomMetadataMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(BankModule).
			WithMsgType(BankDenomMetadataMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}

// Query the total supply of coins of the chain.
func (e BankExternal) Total(totalMsg ...types.TotalMsg) provider.XplaClient {
	if len(totalMsg) == 0 {
		msg, err := MakeTotalSupplyMsg()
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(BankModule).
			WithMsgType(BankTotalMsgType).
			WithMsg(msg)
	} else if len(totalMsg) == 1 {
		msg, err := MakeSupplyOfMsg(totalMsg[0])
		if err != nil {
			return provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(err)
		}
		e.Xplac.WithModule(BankModule).
			WithMsgType(BankTotalSupplyOfMsgType).
			WithMsg(msg)
	} else {
		provider.ResetModuleAndMsgXplac(e.Xplac).WithErr(util.LogErr(errors.ErrInvalidRequest, "need only one parameter"))
	}
	return e.Xplac
}
