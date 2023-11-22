package evm

import (
	"github.com/xpladev/xpla.go/types"
	"github.com/xpladev/xpla.go/types/errors"
	"github.com/xpladev/xpla.go/util"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rpc"
)

// Parsing - send coin
func parseSendCoinArgs(sendCoinMsg types.SendCoinMsg) (types.SendCoinMsg, error) {
	sendCoinMsg.Amount = util.DenomRemove(sendCoinMsg.Amount)
	return sendCoinMsg, nil
}

// Parsing - deploy solidity contract
func parseDeploySolContractArgs(deploySolContractMsg types.DeploySolContractMsg) (ContractInfo, error) {
	var err error
	bytecode := deploySolContractMsg.Bytecode
	if deploySolContractMsg.BytecodeJsonFilePath != "" {
		bytecode, err = util.BytecodeParsing(deploySolContractMsg.BytecodeJsonFilePath)
		if err != nil {
			return ContractInfo{}, util.LogErr(errors.ErrParse, err)
		}
	}

	abi := deploySolContractMsg.ABI
	if deploySolContractMsg.ABIJsonFilePath != "" {
		abi, err = util.AbiParsing(deploySolContractMsg.ABIJsonFilePath)
		if err != nil {
			return ContractInfo{}, util.LogErr(errors.ErrParse, err)
		}
	}

	if deploySolContractMsg.Args == nil || len(deploySolContractMsg.Args) == 0 {
		Args = nil
	} else {
		Args = deploySolContractMsg.Args
	}

	return ContractInfo{
		Abi:      abi,
		Bytecode: bytecode,
	}, nil
}

// Parsing - invoke solidity contract
func parseInvokeSolContractArgs(invokeSolContractMsg types.InvokeSolContractMsg) (types.InvokeSolContractMsg, error) {
	var err error
	if invokeSolContractMsg.BytecodeJsonFilePath != "" {
		invokeSolContractMsg.Bytecode, err = util.BytecodeParsing(invokeSolContractMsg.BytecodeJsonFilePath)
		if err != nil {
			return types.InvokeSolContractMsg{}, util.LogErr(errors.ErrParse, err)
		}
	}

	if invokeSolContractMsg.ABIJsonFilePath != "" {
		invokeSolContractMsg.ABI, err = util.AbiParsing(invokeSolContractMsg.ABIJsonFilePath)
		if err != nil {
			return types.InvokeSolContractMsg{}, util.LogErr(errors.ErrParse, err)
		}
	}
	if invokeSolContractMsg.Args == nil || len(invokeSolContractMsg.Args) == 0 {
		Args = nil
	} else {
		Args = invokeSolContractMsg.Args
	}

	invokeSolContractMsg.ContractAddress = util.FromStringToTypeHexString(invokeSolContractMsg.ContractAddress)

	return invokeSolContractMsg, nil
}

// Parsing - call solidity contract
func parseCallSolContractArgs(callSolContractMsg types.CallSolContractMsg) (CallSolContractParseMsg, error) {
	var err error
	bytecode := callSolContractMsg.Bytecode
	if callSolContractMsg.BytecodeJsonFilePath != "" {
		bytecode, err = util.BytecodeParsing(callSolContractMsg.BytecodeJsonFilePath)
		if err != nil {
			return CallSolContractParseMsg{}, util.LogErr(errors.ErrParse, err)
		}
	}

	abi := callSolContractMsg.ABI
	if callSolContractMsg.ABIJsonFilePath != "" {
		abi, err = util.AbiParsing(callSolContractMsg.ABIJsonFilePath)
		if err != nil {
			return CallSolContractParseMsg{}, util.LogErr(errors.ErrParse, err)
		}
	}

	if callSolContractMsg.Args == nil || len(callSolContractMsg.Args) == 0 {
		Args = nil
	} else {
		Args = callSolContractMsg.Args
	}
	callByteData, err := util.GetAbiPack(callSolContractMsg.ContractFuncCallName, abi, bytecode, Args...)
	if err != nil {
		return CallSolContractParseMsg{}, util.LogErr(errors.ErrParse, err)
	}

	fromAddr := util.FromStringToByte20Address(callSolContractMsg.FromByteAddress)
	toAddr := util.FromStringToByte20Address(callSolContractMsg.ContractAddress)
	value, err := util.FromStringToBigInt("0")
	if err != nil {
		return CallSolContractParseMsg{}, util.LogErr(errors.ErrParse, err)
	}

	msg := ethereum.CallMsg{
		From:  fromAddr,
		To:    &toAddr,
		Value: value,
		Data:  callByteData,
	}

	callSolContractParseMsg := CallSolContractParseMsg{
		CallMsg:  msg,
		CallName: callSolContractMsg.ContractFuncCallName,
		ABI:      abi,
		Bytecode: bytecode,
	}

	return callSolContractParseMsg, nil
}

func parseEthGetBlockTransactionCountArgs(ethGetBlockTransactionCountMsg types.EthGetBlockTransactionCountMsg) (types.EthGetBlockTransactionCountMsg, error) {
	if ethGetBlockTransactionCountMsg.BlockHash != "" {
		ethGetBlockTransactionCountMsg.BlockHash = util.FromStringToTypeHexString(ethGetBlockTransactionCountMsg.BlockHash)
	}

	return ethGetBlockTransactionCountMsg, nil
}

// Parsing - block by hash or height
func parseGetBlockByHashHeightArgs(getBlockByHashHeightMsg types.GetBlockByHashHeightMsg) (types.GetBlockByHashHeightMsg, error) {
	if getBlockByHashHeightMsg.BlockHash != "" && getBlockByHashHeightMsg.BlockHeight != "" {
		return types.GetBlockByHashHeightMsg{}, util.LogErr(errors.ErrInvalidRequest, "need only one parameter, hash or height")
	}

	return getBlockByHashHeightMsg, nil
}

// Parsing - sol contract estimate gas
func parseEstimateGasSolArgs(invokeSolContractMsg types.InvokeSolContractMsg) (CallSolContractParseMsg, error) {
	var err error
	bytecode := invokeSolContractMsg.Bytecode
	if invokeSolContractMsg.BytecodeJsonFilePath != "" {
		bytecode, err = util.BytecodeParsing(invokeSolContractMsg.BytecodeJsonFilePath)
		if err != nil {
			return CallSolContractParseMsg{}, util.LogErr(errors.ErrParse, err)
		}
	}

	abi := invokeSolContractMsg.ABI
	if invokeSolContractMsg.ABIJsonFilePath != "" {
		abi, err = util.AbiParsing(invokeSolContractMsg.ABIJsonFilePath)
		if err != nil {
			return CallSolContractParseMsg{}, util.LogErr(errors.ErrParse, err)
		}
	}
	invokeSolContractMsg.ContractAddress = util.FromStringToTypeHexString(invokeSolContractMsg.ContractAddress)

	var callByteData []byte
	if invokeSolContractMsg.Args == nil || len(invokeSolContractMsg.Args) == 0 {
		callByteData, err = util.GetAbiPack(invokeSolContractMsg.ContractFuncCallName, abi, bytecode)
	} else {
		callByteData, err = util.GetAbiPack(invokeSolContractMsg.ContractFuncCallName, abi, bytecode, invokeSolContractMsg.Args...)
	}

	if err != nil {
		return CallSolContractParseMsg{}, util.LogErr(errors.ErrParse, err)
	}

	fromAddr := util.FromStringToByte20Address(invokeSolContractMsg.FromByteAddress)
	toAddr := util.FromStringToByte20Address(invokeSolContractMsg.ContractAddress)
	value, err := util.FromStringToBigInt("0")
	if err != nil {
		return CallSolContractParseMsg{}, util.LogErr(errors.ErrParse, err)
	}

	msg := ethereum.CallMsg{
		From:  fromAddr,
		To:    &toAddr,
		Value: value,
		Data:  callByteData,
	}

	callSolContractParseMsg := CallSolContractParseMsg{
		CallMsg: msg,
	}

	return callSolContractParseMsg, nil
}

// Parsing - get transaction receipt
func parseEthNewFilterArgs(ethNewFilterMsg types.EthNewFilterMsg) (EthNewFilterParseMsg, error) {
	var fromBlock rpc.BlockNumber
	var toBlock rpc.BlockNumber
	var addresses []common.Address
	var topicsHash []common.Hash
	var topics []interface{}

	if ethNewFilterMsg.FromBlock == "latest" ||
		ethNewFilterMsg.FromBlock == "" {
		fromBlock = rpc.LatestBlockNumber

	} else if ethNewFilterMsg.FromBlock == "earliest" {
		fromBlock = rpc.EarliestBlockNumber

	} else if ethNewFilterMsg.FromBlock == "pending" {
		fromBlock = rpc.PendingBlockNumber

	} else {
		return EthNewFilterParseMsg{}, util.LogErr(errors.ErrInvalidMsgType, "invalid from/to block type, (latest/earliest/pending)")
	}

	if ethNewFilterMsg.ToBlock == "latest" ||
		ethNewFilterMsg.ToBlock == "" {
		toBlock = rpc.LatestBlockNumber

	} else if ethNewFilterMsg.ToBlock == "earliest" {
		toBlock = rpc.EarliestBlockNumber

	} else if ethNewFilterMsg.ToBlock == "pending" {
		toBlock = rpc.PendingBlockNumber

	} else {
		return EthNewFilterParseMsg{}, util.LogErr(errors.ErrInvalidMsgType, "invalid from/to block type, (latest/earliest/pending)")
	}

	if len(ethNewFilterMsg.Address) != 0 {
		for _, address := range ethNewFilterMsg.Address {
			byte20Addr := util.FromStringToByte20Address(address)
			addresses = append(addresses, byte20Addr)
		}
	} else {
		addresses = []common.Address{util.FromStringToByte20Address("0x0")}
	}

	if len(ethNewFilterMsg.Topics) != 0 {
		for _, topic := range ethNewFilterMsg.Topics {
			commonHashTopic := util.FromStringHexToHash(topic)
			topicsHash = append(topicsHash, commonHashTopic)
		}
		topics = append(topics, topicsHash)
	} else {
		topics = append(topics, []common.Hash{})
	}

	varInput := EthNewFilterParseMsg{
		FromBlock: &fromBlock,
		ToBlock:   &toBlock,
		Addresses: addresses,
		Topics:    topics,
	}

	return varInput, nil
}

// Parsing - get logs
func parseEthGetLogsArgs(ethGetLogsMsg types.EthGetLogsMsg) (EthNewFilterParseMsg, error) {
	var blockHash common.Hash
	var fromBlock rpc.BlockNumber
	var toBlock rpc.BlockNumber
	var addresses []common.Address
	var topicsHash []common.Hash
	var topics []interface{}

	if (ethGetLogsMsg.FromBlock != "" || ethGetLogsMsg.ToBlock != "") &&
		ethGetLogsMsg.BlockHash != "" {
		return EthNewFilterParseMsg{}, util.LogErr(errors.ErrInvalidRequest, "cannot specify both BlockHash and FromBlock/ToBlock, choose one or the other")
	}

	if ethGetLogsMsg.FromBlock == "latest" ||
		ethGetLogsMsg.FromBlock == "" {
		fromBlock = rpc.LatestBlockNumber

	} else if ethGetLogsMsg.FromBlock == "earliest" {
		fromBlock = rpc.EarliestBlockNumber

	} else if ethGetLogsMsg.FromBlock == "pending" {
		fromBlock = rpc.PendingBlockNumber

	} else {
		return EthNewFilterParseMsg{}, util.LogErr(errors.ErrInvalidMsgType, "invalid from/to block type, (latest/earliest/pending)")
	}

	if ethGetLogsMsg.ToBlock == "latest" ||
		ethGetLogsMsg.ToBlock == "" {
		toBlock = rpc.LatestBlockNumber

	} else if ethGetLogsMsg.ToBlock == "earliest" {
		toBlock = rpc.EarliestBlockNumber

	} else if ethGetLogsMsg.ToBlock == "pending" {
		toBlock = rpc.PendingBlockNumber

	} else {
		return EthNewFilterParseMsg{}, util.LogErr(errors.ErrInvalidMsgType, "invalid from/to block type, (latest/earliest/pending)")
	}

	if len(ethGetLogsMsg.Address) != 0 {
		for _, address := range ethGetLogsMsg.Address {
			byte20Addr := util.FromStringToByte20Address(address)
			addresses = append(addresses, byte20Addr)
		}
	} else {
		addresses = []common.Address{util.FromStringToByte20Address("0x0")}
	}

	if len(ethGetLogsMsg.Topics) != 0 {
		for _, topic := range ethGetLogsMsg.Topics {
			commonHashTopic := util.FromStringHexToHash(topic)
			topicsHash = append(topicsHash, commonHashTopic)
		}
		topics = append(topics, topicsHash)
	} else {
		topics = append(topics, []common.Hash{})
	}

	if ethGetLogsMsg.BlockHash != "" {
		blockHash = util.FromStringHexToHash(ethGetLogsMsg.BlockHash)

		varInput := EthNewFilterParseMsg{
			BlockHash: &blockHash,
			FromBlock: nil,
			ToBlock:   nil,
			Addresses: addresses,
			Topics:    topics,
		}

		return varInput, nil
	}

	varInput := EthNewFilterParseMsg{
		BlockHash: nil,
		FromBlock: &fromBlock,
		ToBlock:   &toBlock,
		Addresses: addresses,
		Topics:    topics,
	}

	return varInput, nil
}
