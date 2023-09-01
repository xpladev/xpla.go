package bank

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Parsing - bank send
func parseBankSendArgs(bankSendMsg types.BankSendMsg, privKey key.PrivateKey) (banktypes.MsgSend, error) {
	denom := types.XplaDenom

	if bankSendMsg.FromAddress == "" || bankSendMsg.ToAddress == "" || bankSendMsg.Amount == "" {
		return banktypes.MsgSend{}, util.LogErr(errors.ErrInsufficientParams, "no parameters")
	}

	amountBigInt, ok := sdk.NewIntFromString(util.DenomRemove(bankSendMsg.Amount))
	if !ok {
		return banktypes.MsgSend{}, util.LogErr(errors.ErrInvalidRequest, "Wrong amount parameter")
	}

	msg := banktypes.MsgSend{
		FromAddress: bankSendMsg.FromAddress,
		ToAddress:   bankSendMsg.ToAddress,
		Amount:      sdk.NewCoins(sdk.NewCoin(denom, amountBigInt)),
	}

	return msg, nil

}

// Parsing - all balances
func parseBankAllBalancesArgs(bankBalancesMsg types.BankBalancesMsg) (banktypes.QueryAllBalancesRequest, error) {
	addr, err := sdk.AccAddressFromBech32(bankBalancesMsg.Address)
	if err != nil {
		return banktypes.QueryAllBalancesRequest{}, util.LogErr(errors.ErrInvalidRequest, err)
	}

	params := *banktypes.NewQueryAllBalancesRequest(addr, core.PageRequest)
	return params, nil
}

// Parsing - balance
func parseBankBalanceArgs(bankBalancesMsg types.BankBalancesMsg) (banktypes.QueryBalanceRequest, error) {
	addr, err := sdk.AccAddressFromBech32(bankBalancesMsg.Address)
	if err != nil {
		return banktypes.QueryBalanceRequest{}, util.LogErr(errors.ErrInvalidRequest, err)
	}

	params := *banktypes.NewQueryBalanceRequest(addr, bankBalancesMsg.Denom)
	return params, nil
}
