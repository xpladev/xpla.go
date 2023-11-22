package wasm

import (
	"encoding/base64"
	"encoding/hex"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/xpladev/xpla.go/core"
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"
)

// (Tx) make msg - store code
func MakeStoreCodeMsg(storeMsg types.StoreMsg, addr sdk.AccAddress) (wasmtypes.MsgStoreCode, error) {
	msg, err := parseStoreCodeArgs(storeMsg, addr)
	if err != nil {
		return wasmtypes.MsgStoreCode{}, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return wasmtypes.MsgStoreCode{}, err
	}

	return msg, nil
}

// (Tx) make msg - instantiate
func MakeInstantiateMsg(instantiateMsg types.InstantiateMsg, addr sdk.AccAddress) (wasmtypes.MsgInstantiateContract, error) {
	if (types.InstantiateMsg{}) == instantiateMsg {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}

	if instantiateMsg.CodeId == "" ||
		instantiateMsg.Amount == "" ||
		instantiateMsg.Label == "" ||
		instantiateMsg.InitMsg == "" {
		return wasmtypes.MsgInstantiateContract{}, util.LogErr(errors.ErrInsufficientParams, "Empty mandatory parameters")
	}

	msg, err := parseInstantiateArgs(instantiateMsg, addr)
	if err != nil {
		return wasmtypes.MsgInstantiateContract{}, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return wasmtypes.MsgInstantiateContract{}, err
	}

	return msg, nil
}

// (Tx) make msg - execute
func MakeExecuteMsg(executeMsg types.ExecuteMsg, addr sdk.AccAddress) (wasmtypes.MsgExecuteContract, error) {
	if (types.ExecuteMsg{}) == executeMsg {
		return wasmtypes.MsgExecuteContract{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}

	msg, err := parseExecuteArgs(executeMsg, addr)
	if err != nil {
		return wasmtypes.MsgExecuteContract{}, err
	}

	if err = msg.ValidateBasic(); err != nil {
		return wasmtypes.MsgExecuteContract{}, err
	}

	return msg, nil
}

// (Tx) make msg - clear contract admin
func MakeClearContractAdminMsg(clearContractAdminMsg types.ClearContractAdminMsg, sender sdk.AccAddress) (wasmtypes.MsgClearAdmin, error) {
	return parseClearContractAdminArgs(clearContractAdminMsg, sender)
}

// (Tx) make msg - set contract admin
func MakeSetContractAdmintMsg(setContractAdminMsg types.SetContractAdminMsg, sender sdk.AccAddress) (wasmtypes.MsgUpdateAdmin, error) {
	return parseSetContractAdmintArgs(setContractAdminMsg, sender)
}

// (Tx) make msg - migrate
func MakeMigrateMsg(migrateMsg types.MigrateMsg, sender sdk.AccAddress) (wasmtypes.MsgMigrateContract, error) {
	return parseMigrateArgs(migrateMsg, sender)
}

// (Query) make msg - query contract
func MakeQueryMsg(queryMsg types.QueryMsg) (wasmtypes.QuerySmartContractStateRequest, error) {
	if (types.QueryMsg{}) == queryMsg {
		return wasmtypes.QuerySmartContractStateRequest{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}

	return parseQueryArgs(queryMsg)
}

// (Query) make msg - list code
func MakeListcodeMsg() (wasmtypes.QueryCodesRequest, error) {
	return wasmtypes.QueryCodesRequest{
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - list contract by code
func MakeListContractByCodeMsg(listContractByCodeMsg types.ListContractByCodeMsg) (wasmtypes.QueryContractsByCodeRequest, error) {
	if (types.ListContractByCodeMsg{}) == listContractByCodeMsg {
		return wasmtypes.QueryContractsByCodeRequest{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}
	codeIdU64, err := util.FromStringToUint64(listContractByCodeMsg.CodeId)
	if err != nil {
		return wasmtypes.QueryContractsByCodeRequest{}, util.LogErr(errors.ErrParse, err)
	}
	return wasmtypes.QueryContractsByCodeRequest{
		CodeId:     codeIdU64,
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - download
func MakeDownloadMsg(downloadMsg types.DownloadMsg) ([]interface{}, error) {
	var msgInterfaceSlice []interface{}
	if (types.DownloadMsg{}) == downloadMsg {
		return nil, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}

	codeIdU64, err := util.FromStringToUint64(downloadMsg.CodeId)
	if err != nil {
		return nil, util.LogErr(errors.ErrParse, err)
	}
	msg := wasmtypes.QueryCodeRequest{
		CodeId: codeIdU64,
	}
	msgInterfaceSlice = append(msgInterfaceSlice, msg)
	msgInterfaceSlice = append(msgInterfaceSlice, downloadMsg.DownloadFileName)
	return msgInterfaceSlice, nil
}

// (Query) make msg - code info
func MakeCodeInfoMsg(codeInfoMsg types.CodeInfoMsg) (wasmtypes.QueryCodeRequest, error) {
	if (types.CodeInfoMsg{}) == codeInfoMsg {
		return wasmtypes.QueryCodeRequest{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}
	codeIdU64, err := util.FromStringToUint64(codeInfoMsg.CodeId)
	if err != nil {
		return wasmtypes.QueryCodeRequest{}, util.LogErr(errors.ErrParse, err)
	}

	return wasmtypes.QueryCodeRequest{
		CodeId: codeIdU64,
	}, nil
}

// (Query) make msg - contract info
func MakeContractInfoMsg(contractInfoMsg types.ContractInfoMsg) (wasmtypes.QueryContractInfoRequest, error) {
	if (types.ContractInfoMsg{}) == contractInfoMsg {
		return wasmtypes.QueryContractInfoRequest{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}
	return wasmtypes.QueryContractInfoRequest{
		Address: contractInfoMsg.ContractAddress,
	}, nil
}

// (Query) make msg - contract state all
func MakeContractStateAllMsg(contractStateAllMsg types.ContractStateAllMsg) (wasmtypes.QueryAllContractStateRequest, error) {
	if (types.ContractStateAllMsg{}) == contractStateAllMsg {
		return wasmtypes.QueryAllContractStateRequest{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}
	return wasmtypes.QueryAllContractStateRequest{
		Address:    contractStateAllMsg.ContractAddress,
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - history
func MakeContractHistoryMsg(contractHistoryMsg types.ContractHistoryMsg) (wasmtypes.QueryContractHistoryRequest, error) {
	if (types.ContractHistoryMsg{}) == contractHistoryMsg {
		return wasmtypes.QueryContractHistoryRequest{}, util.LogErr(errors.ErrInsufficientParams, "Empty request or type of parameter is not correct")
	}
	return wasmtypes.QueryContractHistoryRequest{
		Address:    contractHistoryMsg.ContractAddress,
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - pinned
func MakePinnedMsg() (wasmtypes.QueryPinnedCodesRequest, error) {
	return wasmtypes.QueryPinnedCodesRequest{
		Pagination: core.PageRequest,
	}, nil
}

// (Query) make msg - libwasmvm version
func MakeLibwasmvmVersionMsg() (string, error) {
	return parseLibwasmvmVersionArgs()
}

type ArgumentDecoder struct {
	// dec is the default decoder
	dec                func(string) ([]byte, error)
	asciiF, hexF, b64F bool
}

// Make new query decoder.
func NewArgDecoder(def func(string) ([]byte, error)) *ArgumentDecoder {
	return &ArgumentDecoder{dec: def}
}

func (a *ArgumentDecoder) DecodeString(s string) ([]byte, error) {
	found := -1
	for i, v := range []*bool{&a.asciiF, &a.hexF, &a.b64F} {
		if !*v {
			continue
		}
		if found != -1 {

			return nil, util.LogErr(errors.ErrInvalidRequest, "multiple decoding flags used")
		}
		found = i
	}
	switch found {
	case 0:
		return AsciiDecodeString(s)
	case 1:
		return hex.DecodeString(s)
	case 2:
		return base64.StdEncoding.DecodeString(s)
	default:
		return a.dec(s)
	}
}

func AsciiDecodeString(s string) ([]byte, error) {
	return []byte(s), nil
}
