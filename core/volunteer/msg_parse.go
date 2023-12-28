package volunteer

import (
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	paramsapp "github.com/xpladev/xpla/app/params"
	volunteercli "github.com/xpladev/xpla/x/volunteer/client/cli"
	volunteertypes "github.com/xpladev/xpla/x/volunteer/types"
)

// Parsing - create validator
func parseRegisterVolunteerValidatorArgs(
	registerVolunteerValidatorMsg types.RegisterVolunteerValidatorMsg,
	encodingConfig paramsapp.EncodingConfig,
	from sdk.AccAddress,
) (govtypes.MsgSubmitProposal, error) {
	var proposal volunteertypes.RegisterVolunteerValidatorProposalWithDeposit
	var err error

	if registerVolunteerValidatorMsg.JsonFilePath != "" {
		proposal, err = volunteercli.ParseRegisterVolunteerValidatorProposalWithDeposit(encodingConfig.Codec, registerVolunteerValidatorMsg.JsonFilePath)
		if err != nil {
			return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
		}
	} else {
		proposal.Title = registerVolunteerValidatorMsg.Title
		proposal.Description = registerVolunteerValidatorMsg.Description
		proposal.Deposit = util.DenomAdd(registerVolunteerValidatorMsg.Deposit)
	}

	deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	amount, err := sdk.ParseCoinNormalized(util.DenomAdd(registerVolunteerValidatorMsg.Amount))
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	var pubKey cryptotypes.PubKey
	if err := encodingConfig.Codec.UnmarshalInterfaceJSON([]byte(registerVolunteerValidatorMsg.ValPubKey), &pubKey); err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	stakingDescription := stakingtypes.NewDescription(
		registerVolunteerValidatorMsg.Moniker,
		registerVolunteerValidatorMsg.Identity,
		registerVolunteerValidatorMsg.Website,
		registerVolunteerValidatorMsg.Security,
		registerVolunteerValidatorMsg.Details,
	)

	content, err := volunteertypes.NewRegisterVolunteerValidatorProposal(
		proposal.Title,
		proposal.Description,
		from,
		sdk.ValAddress(from),
		pubKey,
		amount,
		stakingDescription,
	)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	return *msg, nil
}

// Parsing - unregister volunteer validator
func parseUnregisterVolunteerValidatorArgs(unregisterVolunteerValidatorMsg types.UnregisterVolunteerValidatorMsg, encodingConfig paramsapp.EncodingConfig, from sdk.AccAddress) (govtypes.MsgSubmitProposal, error) {
	var proposal volunteertypes.UnregisterVolunteerValidatorProposalWithDeposit
	var err error

	if unregisterVolunteerValidatorMsg.JsonFilePath != "" {
		proposal, err = volunteercli.ParseUnregisterVolunteerValidatorProposalWithDeposit(encodingConfig.Codec, unregisterVolunteerValidatorMsg.JsonFilePath)
		if err != nil {
			return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
		}
	} else {
		proposal.Title = unregisterVolunteerValidatorMsg.Title
		proposal.Description = unregisterVolunteerValidatorMsg.Description
		proposal.Deposit = util.DenomAdd(unregisterVolunteerValidatorMsg.Deposit)
		proposal.ValidatorAddress = unregisterVolunteerValidatorMsg.ValAddress
	}

	deposit, err := sdk.ParseCoinsNormalized(proposal.Deposit)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	valAddr, err := sdk.ValAddressFromBech32(proposal.ValidatorAddress)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	content := volunteertypes.NewUnregisterVolunteerValidatorProposal(proposal.Title, proposal.Description, valAddr)
	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	return *msg, nil
}
