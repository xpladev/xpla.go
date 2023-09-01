package gov_test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	mgov "github.com/xpladev/xpla.go/core/gov"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/util"
	"github.com/xpladev/xpla.go/util/testutil"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	xapp "github.com/xpladev/xpla/app"
)

type MockWeightedProposalContent struct {
	n int
}

func (m MockWeightedProposalContent) ContentSimulatorFn() simtypes.ContentSimulatorFn {
	return func(r *rand.Rand, _ sdk.Context, _ []simtypes.Account) simtypes.Content {
		return govtypes.NewTextProposal(
			fmt.Sprintf("title-%d: %s", m.n, simtypes.RandStringOfLength(r, 100)),
			fmt.Sprintf("description-%d: %s", m.n, simtypes.RandStringOfLength(r, 4000)),
		)
	}
}

var initialProposalID = uint64(100000000000000)

// TestSimulateMsgSubmitProposal tests the normal scenario of a valid message of type TypeMsgSubmitProposal.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgSubmitProposal(t *testing.T) {
	app, ctx := createTestApp(false)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash}})

	// execute operation
	op := SimulateMsgSubmitProposal(app.AccountKeeper, app.BankKeeper, app.GovKeeper, MockWeightedProposalContent{3}.ContentSimulatorFn())
	operationMsg, _, err := op(r, app.BaseApp, ctx, accounts, testutil.TestChainId)
	require.NoError(t, err)

	var msg govtypes.MsgSubmitProposal
	govtypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, "xpla1p8wcgrjr4pjju90xg6u9cgq55dxwq8j7zj7eku", msg.Proposer)
	require.Equal(t, "2686011axpla", msg.InitialDeposit.String())
	require.Equal(t, "title-3: ZBSpYuLyYggwexjxusrBqDOTtGTOWeLrQKjLxzIivHSlcxgdXhhuTSkuxKGLwQvuyNhYFmBZHeAerqyNEUzXPFGkqEGqiQWIXnku", msg.GetContent().GetTitle())
	require.Equal(t, "description-3: NJWzHdBNpAXKJPHWQdrGYcAHSctgVlqwqHoLfHsXUdStwfefwzqLuKEhmMyYLdbZrcPgYqjNHxPexsruwEGStAneKbWkQDDIlCWBLSiAASNhZqNFlPtfqPJoxKsgMdzjWqLWdqKQuJqWPMvwPQWZUtVMOTMYKJbfdlZsjdsomuScvDmbDkgRualsxDvRJuCAmPOXitIbcyWsKGSdrEunFAOdmXnsuyFVgJqEjbklvmwrUlsxjRSfKZxGcpayDdgoFcnVSutxjRgOSFzPwidAjubMncNweqpbxhXGchpZUxuFDOtpnhNUycJICRYqsPhPSCjPTWZFLkstHWJxvdPEAyEIxXgLwbNOjrgzmaujiBABBIXvcXpLrbcEWNNQsbjvgJFgJkflpRohHUutvnaUqoopuKjTDaemDeSdqbnOzcfJpcTuAQtZoiLZOoAIlboFDAeGmSNwkvObPRvRWQgWkGkxwtPauYgdkmypLjbqhlHJIQTntgWjXwZdOyYEdQRRLfMSdnxqppqUofqLbLQDUjwKVKfZJUJQPsWIPwIVaSTrmKskoAhvmZyJgeRpkaTfGgrJzAigcxtfshmiDCFkuiluqtMOkidknnTBtumyJYlIsWLnCQclqdVmikUoMOPdPWwYbJxXyqUVicNxFxyqJTenNblyyKSdlCbiXxUiYUiMwXZASYfvMDPFgxniSjWaZTjHkqlJvtBsXqwPpyVxnJVGFWhfSxgOcduoxkiopJvFjMmFabrGYeVtTXLhxVUEiGwYUvndjFGzDVntUvibiyZhfMQdMhgsiuysLMiePBNXifRLMsSmXPkwlPloUbJveCvUlaalhZHuvdkCnkSHbMbmOnrfEGPwQiACiPlnihiaOdbjPqPiTXaHDoJXjSlZmltGqNHHNrcKdlFSCdmVOuvDcBLdSklyGJmcLTbSFtALdGlPkqqecJrpLCXNPWefoTJNgEJlyMEPneVaxxduAAEqQpHWZodWyRkDAxzyMnFMcjSVqeRXLqsNyNtQBbuRvunZflWSbbvXXdkyLikYqutQhLPONXbvhcQZJPSWnOulqQaXmbfFxAkqfYeseSHOQidHwbcsOaMnSrrmGjjRmEMQNuknupMxJiIeVjmgZvbmjPIQTEhQFULQLBMPrxcFPvBinaOPYWGvYGRKxLZdwamfRQQFngcdSlvwjfaPbURasIsGJVHtcEAxnIIrhSriiXLOlbEBLXFElXJFGxHJczRBIxAuPKtBisjKBwfzZFagdNmjdwIRvwzLkFKWRTDPxJCmpzHUcrPiiXXHnOIlqNVoGSXZewdnCRhuxeYGPVTfrNTQNOxZmxInOazUYNTNDgzsxlgiVEHPKMfbesvPHUqpNkUqbzeuzfdrsuLDpKHMUbBMKczKKWOdYoIXoPYtEjfOnlQLoGnbQUCuERdEFaptwnsHzTJDsuZkKtzMpFaZobynZdzNydEeJJHDYaQcwUxcqvwfWwNUsCiLvkZQiSfzAHftYgAmVsXgtmcYgTqJIawstRYJrZdSxlfRiqTufgEQVambeZZmaAyRQbcmdjVUZZCgqDrSeltJGXPMgZnGDZqISrGDOClxXCxMjmKqEPwKHoOfOeyGmqWqihqjINXLqnyTesZePQRqaWDQNqpLgNrAUKulklmckTijUltQKuWQDwpLmDyxLppPVMwsmBIpOwQttYFMjgJQZLYFPmxWFLIeZihkRNnkzoypBICIxgEuYsVWGIGRbbxqVasYnstWomJnHwmtOhAFSpttRYYzBmyEtZXiCthvKvWszTXDbiJbGXMcrYpKAgvUVFtdKUfvdMfhAryctklUCEdjetjuGNfJjajZtvzdYaqInKtFPPLYmRaXPdQzxdSQfmZDEVHlHGEGNSPRFJuIfKLLfUmnHxHnRjmzQPNlqrXgifUdzAGKVabYqvcDeYoTYgPsBUqehrBhmQUgTvDnsdpuhUoxskDdppTsYMcnDIPSwKIqhXDCIxOuXrywahvVavvHkPuaenjLmEbMgrkrQLHEAwrhHkPRNvonNQKqprqOFVZKAtpRSpvQUxMoXCMZLSSbnLEFsjVfANdQNQVwTmGxqVjVqRuxREAhuaDrFgEZpYKhwWPEKBevBfsOIcaZKyykQafzmGPLRAKDtTcJxJVgiiuUkmyMYuDUNEUhBEdoBLJnamtLmMJQgmLiUELIhLpiEvpOXOvXCPUeldLFqkKOwfacqIaRcnnZvERKRMCKUkMABbDHytQqQblrvoxOZkwzosQfDKGtIdfcXRJNqlBNwOCWoQBcEWyqrMlYZIAXYJmLfnjoJepgSFvrgajaBAIksoyeHqgqbGvpAstMIGmIhRYGGNPRIfOQKsGoKgxtsidhTaAePRCBFqZgPDWCIkqOJezGVkjfYUCZTlInbxBXwUAVRsxHTQtJFnnpmMvXDYCVlEmnZBKhmmxQOIQzxFWpJQkQoSAYzTEiDWEOsVLNrbfzeHFRyeYATakQQWmFDLPbVMCJcWjFGJjfqCoVzlbNNEsqxdSmNPjTjHYOkuEMFLkXYGaoJlraLqayMeCsTjWNRDPBywBJLAPVkGQqTwApVVwYAetlwSbzsdHWsTwSIcctkyKDuRWYDQikRqsKTMJchrliONJeaZIzwPQrNbTwxsGdwuduvibtYndRwpdsvyCktRHFalvUuEKMqXbItfGcNGWsGzubdPMYayOUOINjpcFBeESdwpdlTYmrPsLsVDhpTzoMegKrytNVZkfJRPuDCUXxSlSthOohmsuxmIZUedzxKmowKOdXTMcEtdpHaPWgIsIjrViKrQOCONlSuazmLuCUjLltOGXeNgJKedTVrrVCpWYWHyVrdXpKgNaMJVjbXxnVMSChdWKuZdqpisvrkBJPoURDYxWOtpjzZoOpWzyUuYNhCzRoHsMjmmWDcXzQiHIyjwdhPNwiPqFxeUfMVFQGImhykFgMIlQEoZCaRoqSBXTSWAeDumdbsOGtATwEdZlLfoBKiTvodQBGOEcuATWXfiinSjPmJKcWgQrTVYVrwlyMWhxqNbCMpIQNoSMGTiWfPTCezUjYcdWppnsYJihLQCqbNLRGgqrwHuIvsazapTpoPZIyZyeeSueJuTIhpHMEJfJpScshJubJGfkusuVBgfTWQoywSSliQQSfbvaHKiLnyjdSbpMkdBgXepoSsHnCQaYuHQqZsoEOmJCiuQUpJkmfyfbIShzlZpHFmLCsbknEAkKXKfRTRnuwdBeuOGgFbJLbDksHVapaRayWzwoYBEpmrlAxrUxYMUekKbpjPNfjUCjhbdMAnJmYQVZBQZkFVweHDAlaqJjRqoQPoOMLhyvYCzqEuQsAFoxWrzRnTVjStPadhsESlERnKhpEPsfDxNvxqcOyIulaCkmPdambLHvGhTZzysvqFauEgkFRItPfvisehFmoBhQqmkfbHVsgfHXDPJVyhwPllQpuYLRYvGodxKjkarnSNgsXoKEMlaSKxKdcVgvOkuLcfLFfdtXGTclqfPOfeoVLbqcjcXCUEBgAGplrkgsmIEhWRZLlGPGCwKWRaCKMkBHTAcypUrYjWwCLtOPVygMwMANGoQwFnCqFrUGMCRZUGJKTZIGPyldsifauoMnJPLTcDHmilcmahlqOELaAUYDBuzsVywnDQfwRLGIWozYaOAilMBcObErwgTDNGWnwQMUgFFSKtPDMEoEQCTKVREqrXZSGLqwTMcxHfWotDllNkIJPMbXzjDVjPOOjCFuIvTyhXKLyhUScOXvYthRXpPfKwMhptXaxIxgqBoUqzrWbaoLTVpQoottZyPFfNOoMioXHRuFwMRYUiKvcWPkrayyTLOCFJlAyslDameIuqVAuxErqFPEWIScKpBORIuZqoXlZuTvAjEdlEWDODFRregDTqGNoFBIHxvimmIZwLfFyKUfEWAnNBdtdzDmTPXtpHRGdIbuucfTjOygZsTxPjfweXhSUkMhPjMaxKlMIJMOXcnQfyzeOcbWwNbeH", msg.GetContent().GetDescription())
	require.Equal(t, "gov", msg.Route())
	require.Equal(t, govtypes.TypeMsgSubmitProposal, msg.Type())
}

// TestSimulateMsgDeposit tests the normal scenario of a valid message of type TypeMsgDeposit.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgDeposit(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// setup a proposal
	content := govtypes.NewTextProposal("Test", "description")

	submitTime := ctx.BlockHeader().Time
	depositPeriod := app.GovKeeper.GetDepositParams(ctx).MaxDepositPeriod

	proposal, err := govtypes.NewProposal(content, 1, submitTime, submitTime.Add(depositPeriod))
	require.NoError(t, err)

	app.GovKeeper.SetProposal(ctx, proposal)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := SimulateMsgDeposit(app.AccountKeeper, app.BankKeeper, app.GovKeeper)
	operationMsg, _, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg govtypes.MsgDeposit
	govtypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, uint64(1), msg.ProposalId)
	require.Equal(t, "xpla1ghekyjucln7y67ntx7cf27m9dpuxxemntll57s", msg.Depositor)
	require.Equal(t, "560969axpla", msg.Amount.String())
	require.Equal(t, "gov", msg.Route())
	require.Equal(t, govtypes.TypeMsgDeposit, msg.Type())
}

// TestSimulateMsgVote tests the normal scenario of a valid message of type TypeMsgVote.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgVote(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// setup a proposal
	content := govtypes.NewTextProposal("Test", "description")

	submitTime := ctx.BlockHeader().Time
	depositPeriod := app.GovKeeper.GetDepositParams(ctx).MaxDepositPeriod

	proposal, err := govtypes.NewProposal(content, 1, submitTime, submitTime.Add(depositPeriod))
	require.NoError(t, err)

	app.GovKeeper.ActivateVotingPeriod(ctx, proposal)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := SimulateMsgVote(app.AccountKeeper, app.BankKeeper, app.GovKeeper)
	operationMsg, _, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg govtypes.MsgVote
	govtypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, uint64(1), msg.ProposalId)
	require.Equal(t, "xpla1ghekyjucln7y67ntx7cf27m9dpuxxemntll57s", msg.Voter)
	require.Equal(t, govtypes.OptionYes, msg.Option)
	require.Equal(t, "gov", msg.Route())
	require.Equal(t, govtypes.TypeMsgVote, msg.Type())
}

// TestSimulateMsgVoteWeighted tests the normal scenario of a valid message of type TypeMsgVoteWeighted.
// Abonormal scenarios, where the message is created by an errors are not tested here.
func TestSimulateMsgVoteWeighted(t *testing.T) {
	app, ctx := createTestApp(false)
	blockTime := time.Now().UTC()
	ctx = ctx.WithBlockTime(blockTime)

	// setup 3 accounts
	s := rand.NewSource(1)
	r := rand.New(s)
	accounts := getTestingAccounts(t, r, app, ctx, 3)

	// setup a proposal
	content := govtypes.NewTextProposal("Test", "description")

	submitTime := ctx.BlockHeader().Time
	depositPeriod := app.GovKeeper.GetDepositParams(ctx).MaxDepositPeriod

	proposal, err := govtypes.NewProposal(content, 1, submitTime, submitTime.Add(depositPeriod))
	require.NoError(t, err)

	app.GovKeeper.ActivateVotingPeriod(ctx, proposal)

	// begin a new block
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1, AppHash: app.LastCommitID().Hash, Time: blockTime}})

	// execute operation
	op := SimulateMsgVoteWeighted(app.AccountKeeper, app.BankKeeper, app.GovKeeper)
	operationMsg, _, err := op(r, app.BaseApp, ctx, accounts, "")
	require.NoError(t, err)

	var msg govtypes.MsgVoteWeighted
	govtypes.ModuleCdc.UnmarshalJSON(operationMsg.Msg, &msg)

	require.True(t, operationMsg.OK)
	require.Equal(t, uint64(1), msg.ProposalId)
	require.Equal(t, "xpla1ghekyjucln7y67ntx7cf27m9dpuxxemntll57s", msg.Voter)
	require.True(t, len(msg.Options) >= 1)
	require.Equal(t, "gov", msg.Route())
	require.Equal(t, govtypes.TypeMsgVoteWeighted, msg.Type())
}

// SimulateMsgSubmitProposal simulates creating a msg Submit Proposal
// voting on the proposal, and subsequently slashing the proposal. It is implemented using
// future operations.
func SimulateMsgSubmitProposal(
	ak govtypes.AccountKeeper, bk govtypes.BankKeeper, k keeper.Keeper, contentSim simtypes.ContentSimulatorFn,
) simtypes.Operation {
	// The states are:
	// column 1: All validators vote
	// column 2: 90% vote
	// column 3: 75% vote
	// column 4: 40% vote
	// column 5: 15% vote
	// column 6: noone votes
	// All columns sum to 100 for simplicity, values chosen by @valardragon semi-arbitrarily,
	// feel free to change.
	numVotesTransitionMatrix, _ := simulation.CreateTransitionMatrix([][]int{
		{20, 10, 0, 0, 0, 0},
		{55, 50, 20, 10, 0, 0},
		{25, 25, 30, 25, 30, 15},
		{0, 15, 30, 25, 30, 30},
		{0, 0, 20, 30, 30, 30},
		{0, 0, 0, 10, 10, 25},
	})

	statePercentageArray := []float64{1, .9, .75, .4, .15, 0}
	curNumVotesState := 1

	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		// 1) submit proposal now
		content := contentSim(r, ctx, accs)
		if content == nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgSubmitProposal, "content is nil"), nil, nil
		}

		simAccount, _ := simtypes.RandomAcc(r, accs)
		deposit, skip, err := randomDeposit(r, ctx, ak, bk, k, simAccount.Address)
		switch {
		case skip:
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgSubmitProposal, "skip deposit"), nil, nil
		case err != nil:
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgSubmitProposal, "unable to generate deposit"), nil, err
		}

		submitProposalMsg := types.SubmitProposalMsg{
			Title:       content.GetTitle(),
			Description: content.GetDescription(),
			Type:        "text",
			Deposit:     deposit.String(),
		}
		msg, err := mgov.MakeSubmitProposalMsg(submitProposalMsg, simAccount.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, msg.Type(), "unable to generate a submit proposal msg"), nil, err
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		var fees sdk.Coins
		coins, hasNeg := spendable.SafeSub(deposit)
		if !hasNeg {
			fees, err = simtypes.RandomFees(r, ctx, coins)
			if err != nil {
				return simtypes.NoOpMsg(govtypes.ModuleName, msg.Type(), "unable to generate fees"), nil, err
			}
		}

		txGen := util.MakeEncodingConfig().TxConfig
		tx, err := testutil.GenTx(
			txGen,
			[]sdk.Msg{&msg},
			fees,
			testutil.DefaultTestGenTxGas,
			chainID,
			[]uint64{account.GetAccountNumber()},
			[]uint64{account.GetSequence()},
			simAccount.PrivKey,
		)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, msg.Type(), "unable to generate mock tx"), nil, err
		}

		_, _, err = app.Deliver(txGen.TxEncoder(), tx)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, msg.Type(), "unable to deliver tx"), nil, err
		}

		opMsg := simtypes.NewOperationMsg(&msg, true, "", nil)

		// get the submitted proposal ID
		proposalID, err := k.GetProposalID(ctx)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, msg.Type(), "unable to generate proposalID"), nil, err
		}

		// 2) Schedule operations for votes
		// 2.1) first pick a number of people to vote.
		curNumVotesState = numVotesTransitionMatrix.NextState(r, curNumVotesState)
		numVotes := int(math.Ceil(float64(len(accs)) * statePercentageArray[curNumVotesState]))

		// 2.2) select who votes and when
		whoVotes := r.Perm(len(accs))

		// didntVote := whoVotes[numVotes:]
		whoVotes = whoVotes[:numVotes]
		votingPeriod := k.GetVotingParams(ctx).VotingPeriod

		fops := make([]simtypes.FutureOperation, numVotes+1)
		for i := 0; i < numVotes; i++ {
			whenVote := ctx.BlockHeader().Time.Add(time.Duration(r.Int63n(int64(votingPeriod.Seconds()))) * time.Second)
			fops[i] = simtypes.FutureOperation{
				BlockTime: whenVote,
				Op:        operationSimulateMsgVote(ak, bk, k, accs[whoVotes[i]], int64(proposalID)),
			}
		}

		return opMsg, fops, nil
	}
}

// SimulateMsgDeposit generates a MsgDeposit with random values.
func SimulateMsgDeposit(ak govtypes.AccountKeeper, bk govtypes.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		proposalID, ok := randomProposalID(r, k, ctx, govtypes.StatusDepositPeriod)
		if !ok {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgDeposit, "unable to generate proposalID"), nil, nil
		}

		deposit, skip, err := randomDeposit(r, ctx, ak, bk, k, simAccount.Address)
		switch {
		case skip:
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgDeposit, "skip deposit"), nil, nil
		case err != nil:
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgDeposit, "unable to generate deposit"), nil, err
		}

		govDepositMsg := types.GovDepositMsg{
			ProposalID: util.FromUint64ToString(proposalID),
			Deposit:    deposit.String(),
		}
		msg, err := mgov.MakeGovDepositMsg(govDepositMsg, simAccount.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgDeposit, "make msg err"), nil, err
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		var fees sdk.Coins
		coins, hasNeg := spendable.SafeSub(deposit)
		if !hasNeg {
			fees, err = simtypes.RandomFees(r, ctx, coins)
			if err != nil {
				return simtypes.NoOpMsg(govtypes.ModuleName, msg.Type(), "unable to generate fees"), nil, err
			}
		}

		txCtx := simulation.OperationInput{
			App:           app,
			TxGen:         util.MakeEncodingConfig().TxConfig,
			Cdc:           nil,
			Msg:           &msg,
			MsgType:       msg.Type(),
			Context:       ctx,
			SimAccount:    simAccount,
			AccountKeeper: ak,
			ModuleName:    govtypes.ModuleName,
		}

		return testutil.GenAndDeliverTx(txCtx, fees)
	}
}

// SimulateMsgVote generates a MsgVote with random values.
func SimulateMsgVote(ak govtypes.AccountKeeper, bk govtypes.BankKeeper, k keeper.Keeper) simtypes.Operation {
	return operationSimulateMsgVote(ak, bk, k, simtypes.Account{}, -1)
}

// SimulateMsgVoteWeighted generates a MsgVoteWeighted with random values.
func SimulateMsgVoteWeighted(ak govtypes.AccountKeeper, bk govtypes.BankKeeper, k keeper.Keeper) simtypes.Operation {
	simAccount := simtypes.Account{}
	proposalIDInt := -1
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if simAccount.Equals(simtypes.Account{}) {
			simAccount, _ = simtypes.RandomAcc(r, accs)
		}

		var proposalID uint64

		switch {
		case proposalIDInt < 0:
			var ok bool
			proposalID, ok = randomProposalID(r, k, ctx, govtypes.StatusVotingPeriod)
			if !ok {
				return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVoteWeighted, "unable to generate proposalID"), nil, nil
			}
		default:
			proposalID = uint64(proposalIDInt)
		}

		options := randomWeightedVotingOptions(r)

		var yes string
		var no string
		var abstain string
		var noWithVeto string

		for _, option := range options {
			if option.Option == govtypes.OptionYes {
				yes = option.Weight.String()

			} else if option.Option == govtypes.OptionNo {
				no = option.Weight.String()

			} else if option.Option == govtypes.OptionAbstain {
				abstain = option.Weight.String()

			} else if option.Option == govtypes.OptionNoWithVeto {
				noWithVeto = option.Weight.String()

			} else {
				return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVoteWeighted, "invalid option type"), nil, nil
			}
		}

		weightedVoteMsg := types.WeightedVoteMsg{
			ProposalID: util.FromUint64ToString(proposalID),
			Yes:        yes,
			No:         no,
			Abstain:    abstain,
			NoWithVeto: noWithVeto,
		}
		msg, err := mgov.MakeWeightedVoteMsg(weightedVoteMsg, simAccount.PrivKey)
		if err != nil {
			return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVoteWeighted, "make msg err"), nil, err
		}

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           util.MakeEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             &msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      govtypes.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return testutil.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// returns context and an app with updated mint keeper
func createTestApp(isCheckTx bool) (*xapp.XplaApp, sdk.Context) {
	app := testutil.Setup(isCheckTx, 5)

	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})

	p := minttypes.Params{
		MintDenom:           "axpla",
		InflationRateChange: sdk.NewDecWithPrec(13, 2),
		InflationMax:        sdk.NewDecWithPrec(20, 2),
		InflationMin:        sdk.NewDecWithPrec(7, 2),
		GoalBonded:          sdk.NewDecWithPrec(67, 2),
		BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5 second block times
	}

	app.MintKeeper.SetParams(ctx, p)
	app.MintKeeper.SetMinter(ctx, minttypes.DefaultInitialMinter())

	return app, ctx
}

func getTestingAccounts(t *testing.T, r *rand.Rand, app *xapp.XplaApp, ctx sdk.Context, n int) []simtypes.Account {
	accounts := simtypes.RandomAccounts(r, n)

	initAmt := app.StakingKeeper.TokensFromConsensusPower(ctx, 200)
	initCoins := sdk.NewCoins(sdk.NewCoin("axpla", initAmt))

	// add coins to the accounts
	for _, account := range accounts {
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, account.Address)
		app.AccountKeeper.SetAccount(ctx, acc)
		require.NoError(t, testutil.FundAccount(app.BankKeeper, ctx, account.Address, initCoins))
	}

	return accounts
}

// Pick a random deposit with a random denomination with a
// deposit amount between (0, min(balance, minDepositAmount))
// This is to simulate multiple users depositing to get the
// proposal above the minimum deposit amount
func randomDeposit(r *rand.Rand, ctx sdk.Context,
	ak govtypes.AccountKeeper, bk govtypes.BankKeeper, k keeper.Keeper, addr sdk.AccAddress,
) (deposit sdk.Coins, skip bool, err error) {
	account := ak.GetAccount(ctx, addr)
	spendable := bk.SpendableCoins(ctx, account.GetAddress())

	if spendable.Empty() {
		return nil, true, nil // skip
	}

	minDeposit := sdk.NewCoins(sdk.NewCoin("axpla", govtypes.DefaultMinDepositTokens))
	denomIndex := r.Intn(len(minDeposit))
	denom := minDeposit[denomIndex].Denom

	depositCoins := spendable.AmountOf(denom)
	if depositCoins.IsZero() {
		return nil, true, nil
	}

	maxAmt := depositCoins
	if maxAmt.GT(minDeposit[denomIndex].Amount) {
		maxAmt = minDeposit[denomIndex].Amount
	}

	amount, err := simtypes.RandPositiveInt(r, maxAmt)
	if err != nil {
		return nil, false, err
	}

	return sdk.Coins{sdk.NewCoin(denom, amount)}, false, nil
}

func operationSimulateMsgVote(ak govtypes.AccountKeeper, bk govtypes.BankKeeper, k keeper.Keeper,
	simAccount simtypes.Account, proposalIDInt int64) simtypes.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		if simAccount.Equals(simtypes.Account{}) {
			simAccount, _ = simtypes.RandomAcc(r, accs)
		}

		var proposalID uint64

		switch {
		case proposalIDInt < 0:
			var ok bool
			proposalID, ok = randomProposalID(r, k, ctx, govtypes.StatusVotingPeriod)
			if !ok {
				return simtypes.NoOpMsg(govtypes.ModuleName, govtypes.TypeMsgVote, "unable to generate proposalID"), nil, nil
			}
		default:
			proposalID = uint64(proposalIDInt)
		}

		option := randomVotingOption(r)
		msg := govtypes.NewMsgVote(simAccount.Address, proposalID, option)

		account := ak.GetAccount(ctx, simAccount.Address)
		spendable := bk.SpendableCoins(ctx, account.GetAddress())

		txCtx := simulation.OperationInput{
			R:               r,
			App:             app,
			TxGen:           util.MakeEncodingConfig().TxConfig,
			Cdc:             nil,
			Msg:             msg,
			MsgType:         msg.Type(),
			Context:         ctx,
			SimAccount:      simAccount,
			AccountKeeper:   ak,
			Bankkeeper:      bk,
			ModuleName:      govtypes.ModuleName,
			CoinsSpentInMsg: spendable,
		}

		return testutil.GenAndDeliverTxWithRandFees(txCtx)
	}
}

// Pick a random proposal ID between the initial proposal ID
// (defined in gov GenesisState) and the latest proposal ID
// that matches a given Status.
// It does not provide a default ID.
func randomProposalID(r *rand.Rand, k keeper.Keeper,
	ctx sdk.Context, status govtypes.ProposalStatus) (proposalID uint64, found bool) {
	proposalID, _ = k.GetProposalID(ctx)

	switch {
	case proposalID > initialProposalID:
		// select a random ID between [initialProposalID, proposalID]
		proposalID = uint64(simtypes.RandIntBetween(r, int(initialProposalID), int(proposalID)))

	default:
		// This is called on the first call to this funcion
		// in order to update the global variable
		initialProposalID = proposalID
	}

	proposal, ok := k.GetProposal(ctx, proposalID)
	if !ok || proposal.Status != status {
		return proposalID, false
	}

	return proposalID, true
}

// Pick a random voting option
func randomVotingOption(r *rand.Rand) govtypes.VoteOption {
	switch r.Intn(4) {
	case 0:
		return govtypes.OptionYes
	case 1:
		return govtypes.OptionAbstain
	case 2:
		return govtypes.OptionNo
	case 3:
		return govtypes.OptionNoWithVeto
	default:
		panic("invalid vote option")
	}
}

// Pick a random weighted voting options
func randomWeightedVotingOptions(r *rand.Rand) govtypes.WeightedVoteOptions {
	w1 := r.Intn(100 + 1)
	w2 := r.Intn(100 - w1 + 1)
	w3 := r.Intn(100 - w1 - w2 + 1)
	w4 := 100 - w1 - w2 - w3
	weightedVoteOptions := govtypes.WeightedVoteOptions{}
	if w1 > 0 {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionYes,
			Weight: sdk.NewDecWithPrec(int64(w1), 2),
		})
	}
	if w2 > 0 {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionAbstain,
			Weight: sdk.NewDecWithPrec(int64(w2), 2),
		})
	}
	if w3 > 0 {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionNo,
			Weight: sdk.NewDecWithPrec(int64(w3), 2),
		})
	}
	if w4 > 0 {
		weightedVoteOptions = append(weightedVoteOptions, govtypes.WeightedVoteOption{
			Option: govtypes.OptionNoWithVeto,
			Weight: sdk.NewDecWithPrec(int64(w4), 2),
		})
	}
	return weightedVoteOptions
}
