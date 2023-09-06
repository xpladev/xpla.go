package client

import (
	"github.com/xpladev/xpla.go/controller"
	"github.com/xpladev/xpla.go/core"

	mevm "github.com/xpladev/xpla.go/core/evm"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

// Query transactions and xpla blockchain information.
// Execute a query of functions for all modules.
// After module query messages are generated, it receives query messages/information to the xpla client receiver and transmits a query message.
func (xplac *XplaClient) Query() (string, error) {
	if xplac.GetErr() != nil {
		return "", xplac.GetErr()
	}

	if xplac.GetGrpcUrl() == "" && xplac.GetLcdURL() == "" {
		if xplac.GetModule() == mevm.EvmModule {
			if xplac.GetEvmRpc() == "" {
				return "", util.LogErr(errors.ErrNotSatisfiedOptions, "evm JSON-RPC URL must exist")
			}

		} else {
			return "", util.LogErr(errors.ErrNotSatisfiedOptions, "at least one of the gRPC URL or LCD URL must exist for query")
		}
	}
	queryClient := core.NewIXplaClient(xplac, setQueryType(xplac))

	return controller.Controller().Get(xplac.GetModule()).NewQueryRouter(*queryClient)

	// switch {
	// case xplac.GetModule() == mauth.AuthModule:
	// 	return mauth.QueryAuth(*queryClient)

	// case xplac.GetModule() == mauthz.AuthzModule:
	// 	return mauthz.QueryAuthz(*queryClient)

	// case xplac.GetModule() == mbank.BankModule:
	// 	return mbank.QueryBank(*queryClient)

	// case xplac.GetModule() == mbase.Base:
	// 	return mbase.QueryBase(*queryClient)

	// case xplac.GetModule() == mdist.DistributionModule:
	// 	return mdist.QueryDistribution(*queryClient)

	// case xplac.GetModule() == mevidence.EvidenceModule:
	// 	return mevidence.QueryEvidence(*queryClient)

	// case xplac.GetModule() == mevm.EvmModule:
	// 	return mevm.QueryEvm(*queryClient)

	// case xplac.GetModule() == mfeegrant.FeegrantModule:
	// 	return mfeegrant.QueryFeegrant(*queryClient)

	// case xplac.GetModule() == mgov.GovModule:
	// 	return mgov.QueryGov(*queryClient)

	// case xplac.GetModule() == mibc.IbcModule:
	// 	return mibc.QueryIbc(*queryClient)

	// case xplac.GetModule() == mmint.MintModule:
	// 	return mmint.QueryMint(*queryClient)

	// case xplac.GetModule() == mparams.ParamsModule:
	// 	return mparams.QueryParams(*queryClient)

	// case xplac.GetModule() == mreward.RewardModule:
	// 	return mreward.QueryReward(*queryClient)

	// case xplac.GetModule() == mslashing.SlashingModule:
	// 	return mslashing.QuerySlashing(*queryClient)

	// case xplac.GetModule() == mstaking.StakingModule:
	// 	return mstaking.QueryStaking(*queryClient)

	// case xplac.GetModule() == mupgrade.UpgradeModule:
	// 	return mupgrade.QueryUpgrade(*queryClient)

	// case xplac.GetModule() == mwasm.WasmModule:
	// 	return mwasm.QueryWasm(*queryClient)

	// default:
	// 	return "", util.LogErr(errors.ErrInvalidRequest, "invalid module")
	// }
}

func setQueryType(xplac *XplaClient) uint8 {
	// Default query type is gRPC, not LCD.
	if xplac.GetGrpcUrl() != "" {
		return types.QueryGrpc
	} else {
		return types.QueryLcd
	}
}
