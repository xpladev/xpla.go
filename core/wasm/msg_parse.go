package wasm

import (
	"os"
	"strconv"
	"strings"

	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	"github.com/CosmWasm/wasmd/x/wasm/ioutils"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvm "github.com/CosmWasm/wasmvm"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	instantiateByEverybody = "instantiate-everybody"
	instantiateNobody      = "instantiate-nobody"
	instantiateBySender    = "instantiate-only-sender"
	instantiateByAddress   = "instantiate-only-address"
)

// Parsing - store code
func parseStoreCodeArgs(storeMsg types.StoreMsg, sender sdk.AccAddress) (wasmtypes.MsgStoreCode, error) {
	if storeMsg.FilePath == "" {
		return wasmtypes.MsgStoreCode{}, util.LogErr(errors.ErrInsufficientParams, "filepath is empty")
	}

	wasm, err := os.ReadFile(storeMsg.FilePath)
	if err != nil {
		return wasmtypes.MsgStoreCode{}, util.LogErr(errors.ErrParse, err)
	}

	// gzip the wasm file
	if ioutils.IsWasm(wasm) {
		wasm, err = ioutils.GzipIt(wasm)
		if err != nil {
			return wasmtypes.MsgStoreCode{}, util.LogErr(errors.ErrParse, err)
		}

	} else if !ioutils.IsGzip(wasm) {
		return wasmtypes.MsgStoreCode{}, util.LogErr(errors.ErrInvalidRequest, "invalid input file. Use wasm binary or gzip")
	}

	permission, err := instantiatePermission(storeMsg.InstantiatePermission, sender)
	if err != nil {
		return wasmtypes.MsgStoreCode{}, err
	}

	msg := wasmtypes.MsgStoreCode{
		Sender:                sender.String(),
		WASMByteCode:          wasm,
		InstantiatePermission: permission,
	}
	return msg, nil
}

func instantiatePermission(permission string, sender sdk.AccAddress) (*wasmtypes.AccessConfig, error) {
	var permMethod string
	var onlyAddr string

	if strings.Contains(permission, ".") {
		perm := strings.Split(permission, ".")
		permMethod = perm[0]
		onlyAddr = perm[1]
	} else {
		permMethod = permission
		onlyAddr = ""
	}

	switch {
	case permMethod == "" || permMethod == instantiateByEverybody:
		return &wasmtypes.AllowEverybody, nil

	case permMethod == instantiateBySender:
		x := wasmtypes.AccessTypeOnlyAddress.With(sender)
		return &x, nil

	case permMethod == instantiateByAddress:
		if onlyAddr == "" {
			return nil, util.LogErr(errors.ErrInsufficientParams, "invalid permission, empty address")
		}
		addr, err := sdk.AccAddressFromBech32(onlyAddr)
		if err != nil {
			return nil, util.LogErr(errors.ErrParse, err)
		}
		x := wasmtypes.AccessTypeOnlyAddress.With(addr)
		return &x, nil

	case permMethod == instantiateNobody:
		return &wasmtypes.AllowNobody, nil

	default:
		return nil, util.LogErr(errors.ErrInvalidMsgType, "invalid permission type")
	}
}

// Parsing - instantiate
func parseInstantiateArgs(
	instantiateMsgData types.InstantiateMsg,
	sender sdk.AccAddress) (wasmtypes.MsgInstantiateContract, error) {

	rawCodeID := instantiateMsgData.CodeId
	if rawCodeID == "" {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr(errors.ErrInsufficientParams, "no code ID")
	}

	// get the id of the code to instantiate
	codeID, err := strconv.ParseUint(rawCodeID, 10, 64)
	if err != nil {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr(errors.ErrParse, err)
	}

	amountStr := instantiateMsgData.Amount
	if amountStr == "" {
		amountStr = "0"
	}
	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(amountStr))
	if err != nil {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr(errors.ErrInvalidRequest, "amount:", err)
	}

	label := instantiateMsgData.Label
	if label == "" {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr(errors.ErrInsufficientParams, "label is required on all contracts")
	}

	initMsg := instantiateMsgData.InitMsg
	if initMsg == "" {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr(errors.ErrInsufficientParams, "no Init Message")
	}

	adminStr := instantiateMsgData.Admin

	noAdminBool := true
	noAdminStr := instantiateMsgData.NoAdmin
	if noAdminStr == "true" {
		noAdminBool = true
	} else if noAdminStr == "" || noAdminStr == "false" {
		noAdminBool = false
	} else {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr(errors.ErrInvalidMsgType, "noAdmin parameter must set \"true\" or \"false\"")
	}

	// ensure sensible admin is set (or explicitly immutable)
	if adminStr == "" && !noAdminBool {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr(errors.ErrInvalidRequest, "you must set an admin or explicitly pass --no-admin to make it immutible (wasmd issue #719)")
	}
	if adminStr != "" && noAdminBool {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr(errors.ErrInvalidRequest, "you set an admin and passed --no-admin, those cannot both be true")
	}

	// build and sign the transaction, then broadcast to Tendermint
	msg := wasmtypes.MsgInstantiateContract{
		Sender: sender.String(),
		CodeID: codeID,
		Label:  label,
		Funds:  amount,
		Msg:    []byte(initMsg),
		Admin:  adminStr,
	}
	return msg, nil
}

// Parsing - execute
func parseExecuteArgs(executeMsgData types.ExecuteMsg,
	sender sdk.AccAddress) (wasmtypes.MsgExecuteContract, error) {
	amountStr := executeMsgData.Amount
	if amountStr == "" {
		amountStr = "0"
	}
	amount, err := sdk.ParseCoinsNormalized(util.DenomAdd(amountStr))
	if err != nil {
		return wasmtypes.MsgExecuteContract{}, util.LogErr(errors.ErrInvalidRequest, "amount:", err)
	}

	return wasmtypes.MsgExecuteContract{
		Sender:   sender.String(),
		Contract: executeMsgData.ContractAddress,
		Funds:    amount,
		Msg:      []byte(executeMsgData.ExecMsg),
	}, nil
}

// Parsing - clear contract admin
func parseClearContractAdminArgs(clearContractAdminMsg types.ClearContractAdminMsg, sender sdk.AccAddress) (wasmtypes.MsgClearAdmin, error) {
	return wasmtypes.MsgClearAdmin{
		Sender:   sender.String(),
		Contract: clearContractAdminMsg.ContractAddress,
	}, nil
}

// Parsing - set contract admin
func parseSetContractAdmintArgs(setContractAdminMsg types.SetContractAdminMsg, sender sdk.AccAddress) (wasmtypes.MsgUpdateAdmin, error) {
	return wasmtypes.MsgUpdateAdmin{
		Sender:   sender.String(),
		Contract: setContractAdminMsg.ContractAddress,
		NewAdmin: setContractAdminMsg.NewAdmin,
	}, nil
}

// Parsing - migrate
func parseMigrateArgs(migrateMsg types.MigrateMsg, sender sdk.AccAddress) (wasmtypes.MsgMigrateContract, error) {
	codeIdU64, err := util.FromStringToUint64(migrateMsg.CodeId)
	if err != nil {
		return wasmtypes.MsgMigrateContract{}, util.LogErr(errors.ErrParse, err)
	}
	return wasmtypes.MsgMigrateContract{
		Sender:   sender.String(),
		Contract: migrateMsg.ContractAddress,
		CodeID:   codeIdU64,
		Msg:      []byte(migrateMsg.MigrateMsg),
	}, nil
}

// Parsing - query contract
func parseQueryArgs(queryMsgData types.QueryMsg) (wasmtypes.QuerySmartContractStateRequest, error) {
	decoder := NewArgDecoder(AsciiDecodeString)

	queryData, err := decoder.DecodeString(queryMsgData.QueryMsg)
	if err != nil {
		return wasmtypes.QuerySmartContractStateRequest{}, util.LogErr(errors.ErrParse, err)
	}

	return wasmtypes.QuerySmartContractStateRequest{
		Address:   queryMsgData.ContractAddress,
		QueryData: queryData,
	}, nil
}

// Parsing - libwasmvm version
func parseLibwasmvmVersionArgs() (string, error) {
	version, err := wasmvm.LibwasmvmVersion()
	if err != nil {
		return "", util.LogErr(errors.ErrParse, err)
	}
	return version, nil
}
