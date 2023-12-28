package volunteer

import (
	"github.com/xpladev/xpla.go/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramsapp "github.com/xpladev/xpla/app/params"
	volunteertypes "github.com/xpladev/xpla/x/volunteer/types"
)

// (Tx) make msg - register volunteer validator
func MakeRegisterVolunteerValidatorMsg(registerVolunteerValidatorMsg types.RegisterVolunteerValidatorMsg, encodingConfig paramsapp.EncodingConfig, from sdk.AccAddress) (govtypes.MsgSubmitProposal, error) {
	return parseRegisterVolunteerValidatorArgs(registerVolunteerValidatorMsg, encodingConfig, from)
}

// (Tx) make msg - unregister volunteer validator
func MakeUnregisterVolunteerValidatorMsg(unregisterVolunteerValidatorMsg types.UnregisterVolunteerValidatorMsg, encodingConfig paramsapp.EncodingConfig, from sdk.AccAddress) (govtypes.MsgSubmitProposal, error) {
	return parseUnregisterVolunteerValidatorArgs(unregisterVolunteerValidatorMsg, encodingConfig, from)
}

// (Query) make msg - query volunteer validators
func MakeQueryVolunteerValidatorsMsg() (volunteertypes.QueryVolunteerValidatorsRequest, error) {
	return volunteertypes.QueryVolunteerValidatorsRequest{}, nil
}
