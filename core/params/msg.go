package params

import (
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramsproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/xpladev/xpla/app/params"
)

// (Tx) make msg - param change
func MakeProposalParamChangeMsg(paramChangeMsg types.ParamChangeMsg, privKey key.PrivateKey, encodingConfig params.EncodingConfig) (govtypes.MsgSubmitProposal, error) {
	return parseProposalParamChangeArgs(paramChangeMsg, privKey, encodingConfig)
}

// (Query) make msg - subspace
func MakeQueryParamsSubspaceMsg(subspaceMsg types.SubspaceMsg) (paramsproposal.QueryParamsRequest, error) {
	return paramsproposal.QueryParamsRequest{
		Subspace: subspaceMsg.Subspace,
		Key:      subspaceMsg.Key,
	}, nil
}
