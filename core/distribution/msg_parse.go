package distribution

import (
	"context"

	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"

	sdk "github.com/cosmos/cosmos-sdk/types"
	distcli "github.com/cosmos/cosmos-sdk/x/distribution/client/cli"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/gogo/protobuf/grpc"
	"github.com/xpladev/xpla/app/params"
)

// Parsing - fund community pool
func parseFundCommunityPoolArgs(fundCommunityPoolMsg types.FundCommunityPoolMsg, depositorAddr sdk.AccAddress) (disttypes.MsgFundCommunityPool, error) {
	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(fundCommunityPoolMsg.Amount))
	if err != nil {
		return disttypes.MsgFundCommunityPool{}, types.ErrWrap(types.ErrParse, err)
	}

	msg := disttypes.NewMsgFundCommunityPool(amount, depositorAddr)
	return *msg, nil
}

// Parsing - proposal community pool
func parseProposalCommunityPoolSpendArgs(communityPoolSpendMsg types.CommunityPoolSpendMsg, from sdk.AccAddress, encodingConfig params.EncodingConfig) (govtypes.MsgSubmitProposal, error) {
	var proposal disttypes.CommunityPoolSpendProposalWithDeposit
	var err error

	if communityPoolSpendMsg.JsonFilePath != "" {
		proposal, err = distcli.ParseCommunityPoolSpendProposalWithDeposit(encodingConfig.Codec, communityPoolSpendMsg.JsonFilePath)
		if err != nil {
			return govtypes.MsgSubmitProposal{}, types.ErrWrap(types.ErrParse, err)
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
		return govtypes.MsgSubmitProposal{}, types.ErrWrap(types.ErrParse, err)
	}

	deposit, err := sdk.ParseCoinsNormalized(util.DenomAdd(proposal.Deposit))
	if err != nil {
		return govtypes.MsgSubmitProposal{}, types.ErrWrap(types.ErrParse, err)
	}

	recpAddr, err := sdk.AccAddressFromBech32(proposal.Recipient)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, types.ErrWrap(types.ErrParse, err)
	}

	content := disttypes.NewCommunityPoolSpendProposal(proposal.Title, proposal.Description, recpAddr, amount)

	msg, err := govtypes.NewMsgSubmitProposal(content, deposit, from)
	if err != nil {
		return govtypes.MsgSubmitProposal{}, types.ErrWrap(types.ErrInvalidRequest, err)
	}

	return *msg, nil
}

// Parsing - withdraw rewards
func parseWithdrawRewardsArgs(withdrawRewardsMsg types.WithdrawRewardsMsg, delAddr sdk.AccAddress) ([]sdk.Msg, error) {
	valAddr, err := sdk.ValAddressFromBech32(withdrawRewardsMsg.ValidatorAddr)
	if err != nil {
		return nil, types.ErrWrap(types.ErrParse, err)
	}

	msgs := []sdk.Msg{disttypes.NewMsgWithdrawDelegatorReward(delAddr, valAddr)}
	if withdrawRewardsMsg.Commission {
		msgs = append(msgs, disttypes.NewMsgWithdrawValidatorCommission(valAddr))
	}

	for _, msg := range msgs {
		if err := msg.ValidateBasic(); err != nil {
			return nil, types.ErrWrap(types.ErrInvalidRequest, err)
		}
	}

	return msgs, nil
}

// Parsing - withdraw all rewards
func parseWithdrawAllRewardsArgs(delAddr sdk.AccAddress, grpcConn grpc.ClientConn, ctx context.Context) ([]sdk.Msg, error) {
	queryClient := disttypes.NewQueryClient(grpcConn)
	delValsRes, err := queryClient.DelegatorValidators(
		ctx,
		&disttypes.QueryDelegatorValidatorsRequest{
			DelegatorAddress: delAddr.String(),
		},
	)
	if err != nil {
		return nil, types.ErrWrap(types.ErrGrpcRequest, err)
	}

	vals := delValsRes.Validators
	msgs := make([]sdk.Msg, 0, len(vals))
	for _, valAddr := range vals {
		val, err := sdk.ValAddressFromBech32(valAddr)
		if err != nil {
			return nil, types.ErrWrap(types.ErrParse, err)
		}

		msg := disttypes.NewMsgWithdrawDelegatorReward(delAddr, val)
		if err := msg.ValidateBasic(); err != nil {
			return nil, types.ErrWrap(types.ErrInvalidRequest, err)
		}
		msgs = append(msgs, msg)
	}

	return msgs, nil
}

// Parsing - set withdraw addr
func parseSetWithdrawAddrArgs(setWithdrawAddrMsg types.SetWithdrawAddrMsg, delAddr sdk.AccAddress) (disttypes.MsgSetWithdrawAddress, error) {
	withdrawAddr, err := sdk.AccAddressFromBech32(setWithdrawAddrMsg.WithdrawAddr)
	if err != nil {
		return disttypes.MsgSetWithdrawAddress{}, types.ErrWrap(types.ErrParse, err)
	}

	msg := disttypes.NewMsgSetWithdrawAddress(delAddr, withdrawAddr)

	return *msg, nil
}

// Parsing - validator outstanding rewards
func parseValidatorOutstandingRewardsArgs(validatorOutstandingRewardsMsg types.ValidatorOutstandingRewardsMsg) (disttypes.QueryValidatorOutstandingRewardsRequest, error) {
	valAddr, err := sdk.ValAddressFromBech32(validatorOutstandingRewardsMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryValidatorOutstandingRewardsRequest{}, types.ErrWrap(types.ErrParse, err)
	}

	return disttypes.QueryValidatorOutstandingRewardsRequest{
		ValidatorAddress: valAddr.String(),
	}, nil
}

// Parsing - commission
func parseQueryDistCommissionArgs(queryDistCommissionMsg types.QueryDistCommissionMsg) (disttypes.QueryValidatorCommissionRequest, error) {
	valAddr, err := sdk.ValAddressFromBech32(queryDistCommissionMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryValidatorCommissionRequest{}, types.ErrWrap(types.ErrParse, err)
	}

	return disttypes.QueryValidatorCommissionRequest{
		ValidatorAddress: valAddr.String(),
	}, nil
}

// Parsing - distribution slashes
func parseDistSlashesArgs(queryDistSlashesMsg types.QueryDistSlashesMsg) (disttypes.QueryValidatorSlashesRequest, error) {
	valAddr, err := sdk.ValAddressFromBech32(queryDistSlashesMsg.ValidatorAddr)
	if err != nil {
		return disttypes.QueryValidatorSlashesRequest{}, types.ErrWrap(types.ErrParse, err)
	}
	startHeightNumber, err := util.FromStringToUint64(queryDistSlashesMsg.StartHeight)
	if err != nil {
		return disttypes.QueryValidatorSlashesRequest{}, types.ErrWrap(types.ErrConvert, err)
	}
	endHeightNumber, err := util.FromStringToUint64(queryDistSlashesMsg.EndHeight)
	if err != nil {
		return disttypes.QueryValidatorSlashesRequest{}, types.ErrWrap(types.ErrConvert, err)
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
		return disttypes.QueryDelegationRewardsRequest{}, types.ErrWrap(types.ErrParse, err)
	}

	return disttypes.QueryDelegationRewardsRequest{
		DelegatorAddress: delAddr,
		ValidatorAddress: valAddr.String(),
	}, nil
}
