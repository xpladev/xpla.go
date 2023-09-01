package distribution

import (
	"context"

	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/key"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distcli "github.com/cosmos/cosmos-sdk/x/distribution/client/cli"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/xpladev/xpla/app/params"
)

// Parsing - fund community pool
func parseFundCommunityPoolArgs(fundCommunityPoolMsg types.FundCommunityPoolMsg, privKey key.PrivateKey) (disttypes.MsgFundCommunityPool, error) {
	depositorAddr, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return disttypes.MsgFundCommunityPool{}, util.LogErr(errors.ErrParse, err)
	}

	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(fundCommunityPoolMsg.Amount))
	if err != nil {
		return disttypes.MsgFundCommunityPool{}, util.LogErr(errors.ErrParse, err)
	}

	msg := disttypes.NewMsgFundCommunityPool(amount, depositorAddr)
	return *msg, nil
}

// Parsing - proposal community pool
func parseProposalCommunityPoolSpendArgs(communityPoolSpendMsg types.CommunityPoolSpendMsg, privKey key.PrivateKey, encodingConfig params.EncodingConfig) (govtypes.MsgSubmitProposal, error) {
	var proposal disttypes.CommunityPoolSpendProposalWithDeposit
	var err error

	if communityPoolSpendMsg.JsonFilePath != "" {
		proposal, err = distcli.ParseCommunityPoolSpendProposalWithDeposit(encodingConfig.Marshaler, communityPoolSpendMsg.JsonFilePath)
		if err != nil {
			return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
		}
	} else {
		proposal.Title = communityPoolSpendMsg.Title
		proposal.Description = communityPoolSpendMsg.Description
		proposal.Recipient = communityPoolSpendMsg.Recipient
		proposal.Amount = communityPoolSpendMsg.Amount
		proposal.Deposit = communityPoolSpendMsg.Deposit
	}

	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(proposal.Amount))
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	deposit, err := sdk.ParseCoinsNormalized(util.DenomAdd(proposal.Deposit))
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	from, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}
	recpAddr, err := sdk.AccAddressFromBech32(proposal.Recipient)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	content := disttypes.NewCommunityPoolSpendProposal(proposal.Title, proposal.Description, recpAddr, amount)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, util.LogErr(errors.ErrParse, err)
	}

	return *msg, nil
}

// Parsing - withdraw rewards
func parseWithdrawRewardsArgs(withdrawRewardsMsg types.WithdrawRewardsMsg, privKey key.PrivateKey) ([]sdk.Msg, error) {
	delAddr, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	valAddr, err := sdk.ValAddressFromBech32(withdrawRewardsMsg.ValidatorAddr)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}

	msgs := []sdk.Msg{disttypes.NewMsgWithdrawDelegatorReward(delAddr, valAddr)}
	if withdrawRewardsMsg.Commission {
		msgs = append(msgs, disttypes.NewMsgWithdrawValidatorCommission(valAddr))
	}

	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}
	}

	return msgs, nil
}

// Parsing - withdraw all rewards
func parseWithdrawAllRewardsArgs(privKey key.PrivateKey, grpcConn grpc.ClientConn, ctx context.Context) ([]sdk.Msg, error) {
	delAddr, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}
	queryClient := disttypes.NewQueryClient(grpcConn)
	delValsRes, err := queryClient.DelegatorValidators(
		ctx,
		&disttypes.QueryDelegatorValidatorsRequest{
			DelegatorAddress: delAddr.String(),
		},
	)
	if err != nil {
		return nil, util.LogErr(errors.ErrGrpcRequest, err)
	}

	vals := delValsRes.Validators
	msgs := make([]sdk.Msg, 0, len(vals))
	for _, valAddr := range vals {
		val, err := sdk.ValAddressFromBech32(valAddr)
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}

		msg := disttypes.NewMsgWithdrawDelegatorReward(delAddr, val)
		if err := msg.ValidateBasic(); err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

// Parsing - set withdraw addr
func parseSetWithdrawAddrArgs(setWithdrawAddrMsg types.SetWithdrawAddrMsg, privKey key.PrivateKey) (disttypes.MsgSetWithdrawAddress, error) {
	delAddr, err := util.GetAddrByPrivKey(privKey)
	if err != nil {
		return disttypes.MsgSetWithdrawAddress{}, util.LogErr(errors.ErrParse, err)
	}
	withdrawAddr, err := sdk.AccAddressFromBech32(setWithdrawAddrMsg.WithdrawAddr)
	if err != nil {
		return disttypes.MsgSetWithdrawAddress{}, util.LogErr(errors.ErrParse, err)
	}

	msg := disttypes.NewMsgSetWithdrawAddress(delAddr, withdrawAddr)

	return *msg, nil
}

// Parsing - validator outstanding rewards
func parseValidatorOutstandingRewardsArgs(validatorOutstandingRewardsMsg types.ValidatorOutstandingRewardsMsg) (disttypes.QueryValidatorOutstandingRewardsRequest, error) {
	valAddr, err := sdk.ValAddressFromBech32(validatorOutstandingRewardsMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryValidatorOutstandingRewardsRequest{}, util.LogErr(errors.ErrParse, err)
	}

	return disttypes.QueryValidatorOutstandingRewardsRequest{
		ValidatorAddress: valAddr.String(),
	}, nil
}

// Parsing - commission
func parseQueryDistCommissionArgs(queryDistCommissionMsg types.QueryDistCommissionMsg) (disttypes.QueryValidatorCommissionRequest, error) {
	valAddr, err := sdk.ValAddressFromBech32(queryDistCommissionMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryValidatorCommissionRequest{}, util.LogErr(errors.ErrParse, err)
	}

	return disttypes.QueryValidatorCommissionRequest{
		ValidatorAddress: valAddr.String(),
	}, nil
}

// Parsing - distribution slashes
func parseDistSlashesArgs(queryDistSlashesMsg types.QueryDistSlashesMsg) (disttypes.QueryValidatorSlashesRequest, error) {
	valAddr, err := sdk.ValAddressFromBech32(queryDistSlashesMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryValidatorSlashesRequest{}, util.LogErr(errors.ErrParse, err)
	}
	startHeightNumber, err := util.FromStringToUint64(queryDistSlashesMsg.StartHeight)
	if err != nil {
		return disttypes.QueryValidatorSlashesRequest{}, util.LogErr(errors.ErrParse, err)
	}
	endHeightNumber, err := util.FromStringToUint64(queryDistSlashesMsg.EndHeight)
	if err != nil {
		return disttypes.QueryValidatorSlashesRequest{}, util.LogErr(errors.ErrParse, err)
	}

	pageReq := core.PageRequest

	return disttypes.QueryValidatorSlashesRequest{
		ValidatorAddress: valAddr.String(),
		StartingHeight:   startHeightNumber,
		EndingHeight:     endHeightNumber,
		Pagination:       pageReq,
	}, nil
}

// Parsing - distribution rewards
func parseQueryDistRewardsArgs(queryDistRewardsMsg types.QueryDistRewardsMsg) (disttypes.QueryDelegationRewardsRequest, error) {
	delAddr := queryDistRewardsMsg.DelegatorAddr
	valAddr, err := sdk.ValAddressFromBech32(queryDistRewardsMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryDelegationRewardsRequest{}, util.LogErr(errors.ErrParse, err)
	}

	return disttypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr.String(),
	}, nil
}
