package bank

import (
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// (Tx) make msg - bank send
func MakeBankSendMsg(bankSendMsg types.BankSendMsg) (banktypes.MsgSend, error) {
	return parseBankSendArgs(bankSendMsg)
}

// (Query) make msg - all balances
func MakeBankAllBalancesMsg(bankBalancesMsg types.BankBalancesMsg) (banktypes.QueryAllBalancesRequest, error) {
	if (types.BankBalancesMsg{}) == bankBalancesMsg {
		return banktypes.QueryAllBalancesRequest{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}

	return parseBankAllBalancesArgs(bankBalancesMsg)
}

// (Query) make msg - balance
func MakeBankBalanceMsg(bankBalancesMsg types.BankBalancesMsg) (banktypes.QueryBalanceRequest, error) {
	if (types.BankBalancesMsg{}) == bankBalancesMsg {
		return banktypes.QueryBalanceRequest{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}

	return parseBankBalanceArgs(bankBalancesMsg)
}

// (Query) make msg - denominations metadata
func MakeDenomsMetaDataMsg() (banktypes.QueryDenomsMetadataRequest, error) {
	return banktypes.QueryDenomsMetadataRequest{}, nil
}

// (Query) make msg - denomination metadata
func MakeDenomMetaDataMsg(denomMetadataMsg types.DenomMetadataMsg) (banktypes.QueryDenomMetadataRequest, error) {
	return banktypes.QueryDenomMetadataRequest{
		Denom: denomMetadataMsg.Denom,
	}, nil
}

// (Query) make msg - total supply
func MakeTotalSupplyMsg() (banktypes.QueryTotalSupplyRequest, error) {
	return banktypes.QueryTotalSupplyRequest{Pagination: core.PageRequest}, nil
}

// (Query) make msg - supply of
func MakeSupplyOfMsg(totalMsg types.TotalMsg) (banktypes.QuerySupplyOfRequest, error) {
	return banktypes.QuerySupplyOfRequest{Denom: totalMsg.Denom}, nil
}
