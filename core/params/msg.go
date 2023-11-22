package params

import (
	"github.com/xpladev/xpla.go/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - param change
func MakeProposalParamChangeMsg(paramChangeMsg types.ParamChangeMsg, from sdk.AccAddress, encodingConfig params.EncodingConfig) (govtypes.MsgSubmitProposal, error) {
	return parseProposalParamChangeArgs(paramChangeMsg, from, encodingConfig)
}

// (Query) make msg - subspace
func MakeQueryParamsSubspaceMsg(subspaceMsg types.SubspaceMsg) (paramsproposal.QueryParamsRequest, error) {
	return paramsproposal.QueryParamsRequest{
		Subspace: subspaceMsg.Subspace,
		Key:      subspaceMsg.Key,
	}, nil
}
