package crisis

import (
	"github.com/xpladev/xpla.go/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
)

// Parsing - invariant broken
func parseInvariantBrokenArgs(invariantBrokenMsg types.InvariantBrokenMsg, senderAddr sdk.AccAddress) (crisistypes.MsgVerifyInvariant, error) {
	if invariantBrokenMsg.ModuleName == "" || invariantBrokenMsg.InvariantRoute == "" {
		return crisistypes.MsgVerifyInvariant{}, types.ErrWrap(types.ErrInsufficientParams, "invalid module name or invariant route")
	}
	msg := crisistypes.NewMsgVerifyInvariant(senderAddr, invariantBrokenMsg.ModuleName, invariantBrokenMsg.InvariantRoute)

	return *msg, nil
}
