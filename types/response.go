package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	evmtypes "github.com/ethereum/go-ethereum/core/types"
)

type TxRes struct {
	Response   *sdk.TxResponse
	EvmReceipt *evmtypes.Receipt
}
