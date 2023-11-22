package crisis

import (
	"github.com/xpladev/xpla.go/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// (Tx) make msg - invariant broken
func MakeInvariantRouteMsg(invariantBrokenMsg types.InvariantBrokenMsg, senderAddr sdk.AccAddress) (crisistypes.MsgVerifyInvariant, error) {
	return parseInvariantBrokenArgs(invariantBrokenMsg, senderAddr)
}
