package slashing

import (
	"github.com/xpladev/xpla.go/types"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/xpladev/xpla/app/params"
)

// Parsing - unjail
func parseUnjailArgs(addr sdk.AccAddress) (slashingtypes.MsgUnjail, error) {
	msg := slashingtypes.NewMsgUnjail(sdk.ValAddress(addr))

	return *msg, nil
}

// Parsing - signing info
func parseQuerySigingInfoArgs(signingInfoMsg types.SigningInfoMsg, xplacEncodingConfig params.EncodingConfig) (slashingtypes.QuerySigningInfoRequest, error) {
	if signingInfoMsg.ConsPubKey != "" {
		var pk cryptotypes.PubKey
		err := xplacEncodingConfig.Codec.UnmarshalInterfaceJSON([]byte(signingInfoMsg.ConsPubKey), &pk)
		if err != nil {
			return slashingtypes.QuerySigningInfoRequest{}, types.ErrWrap(types.ErrFailedToUnmarshal, err)
		}

		return slashingtypes.QuerySigningInfoRequest{
			ConsAddress: sdk.ConsAddress(pk.Address()).String(),
		}, nil
	} else if signingInfoMsg.ConsAddr != "" {
		return slashingtypes.QuerySigningInfoRequest{
			ConsAddress: signingInfoMsg.ConsAddr,
		}, nil
	} else {
		return slashingtypes.QuerySigningInfoRequest{}, types.ErrWrap(types.ErrInsufficientParams, "need at least one input")
	}
}
