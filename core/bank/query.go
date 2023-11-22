package bank

import (
	"github.com/gogo/protobuf/proto"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	bankv1beta1 "cosmossdk.io/api/cosmos/bank/v1beta1"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

var out []byte
var res proto.Message
var err error

// Query client for bank module.
func QueryBank(i core.QueryClient) (string, error) {
	if i.QueryType == types.QueryGrpc {
		return queryByGrpcBank(i)
	} else {
		return queryByLcdBank(i)
	}
}

func queryByGrpcBank(i core.QueryClient) (string, error) {
	queryClient := banktypes.NewQueryClient(i.Ixplac.GetGrpcClient())

	switch {
	// Bank balances
	case i.Ixplac.GetMsgType() == BankAllBalancesMsgType:
		convertMsg := i.Ixplac.GetMsg().(banktypes.QueryAllBalancesRequest)
		res, err = queryClient.AllBalances(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Bank balance
	case i.Ixplac.GetMsgType() == BankBalanceMsgType:
		convertMsg := i.Ixplac.GetMsg().(banktypes.QueryBalanceRequest)
		res, err = queryClient.Balance(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Bank denominations metadata
	case i.Ixplac.GetMsgType() == BankDenomsMetadataMsgType:
		convertMsg := i.Ixplac.GetMsg().(banktypes.QueryDenomsMetadataRequest)
		res, err = queryClient.DenomsMetadata(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Bank denomination metadata
	case i.Ixplac.GetMsgType() == BankDenomMetadataMsgType:
		convertMsg := i.Ixplac.GetMsg().(banktypes.QueryDenomMetadataRequest)
		res, err = queryClient.DenomMetadata(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Bank total
	case i.Ixplac.GetMsgType() == BankTotalMsgType:
		convertMsg := i.Ixplac.GetMsg().(banktypes.QueryTotalSupplyRequest)
		res, err = queryClient.TotalSupply(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	// Bank total supply
	case i.Ixplac.GetMsgType() == BankTotalSupplyOfMsgType:
		convertMsg := i.Ixplac.GetMsg().(banktypes.QuerySupplyOfRequest)
		res, err = queryClient.SupplyOf(
			i.Ixplac.GetContext(),
			&convertMsg,
		)
		if err != nil {
			return "", util.LogErr(errors.ErrGrpcRequest, err)
		}

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	out, err = core.PrintProto(i, res)
	if err != nil {
		return "", err
	}

	return string(out), nil
}

const (
	bankBalancesLabel      = "balances"
	bankDenomMetadataLabel = "denoms_metadata"
	bankSupplyLabel        = "supply"
)

func queryByLcdBank(i core.QueryClient) (string, error) {
	url := util.MakeQueryLcdUrl(bankv1beta1.Query_ServiceDesc.Metadata.(string))

	switch {
	// Bank balances
	case i.Ixplac.GetMsgType() == BankAllBalancesMsgType:
		convertMsg := i.Ixplac.GetMsg().(banktypes.QueryAllBalancesRequest)
		url = url + util.MakeQueryLabels(bankBalancesLabel, convertMsg.Address)

	// Bank balance
	case i.Ixplac.GetMsgType() == BankBalanceMsgType:
		// not supported now.
		convertMsg := i.Ixplac.GetMsg().(banktypes.QueryBalanceRequest)
		url = url + util.MakeQueryLabels(bankBalancesLabel, convertMsg.Address, convertMsg.Denom)

	// Bank denominations metadata
	case i.Ixplac.GetMsgType() == BankDenomsMetadataMsgType:
		url = url + bankDenomMetadataLabel

	// Bank denomination metadata
	case i.Ixplac.GetMsgType() == BankDenomMetadataMsgType:
		convertMsg := i.Ixplac.GetMsg().(banktypes.QueryDenomMetadataRequest)
		url = url + util.MakeQueryLabels(bankDenomMetadataLabel, convertMsg.Denom)

	// Bank total
	case i.Ixplac.GetMsgType() == BankTotalMsgType:
		url = url + bankSupplyLabel

	// Bank total supply
	case i.Ixplac.GetMsgType() == BankTotalSupplyOfMsgType:
		convertMsg := i.Ixplac.GetMsg().(banktypes.QuerySupplyOfRequest)
		url = url + util.MakeQueryLabels(bankSupplyLabel, convertMsg.Denom)

	default:
		return "", util.LogErr(errors.ErrInvalidMsgType, i.Ixplac.GetMsgType())
	}

	i.Ixplac.GetHttpMutex().Lock()
	out, err := util.CtxHttpClient("GET", i.Ixplac.GetLcdURL()+url, nil, i.Ixplac.GetContext())
	if err != nil {
		i.Ixplac.GetHttpMutex().Unlock()
		return "", err
	}
	i.Ixplac.GetHttpMutex().Unlock()

	return string(out), nil
}
