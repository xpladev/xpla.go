package util

import (
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func AbiParsing(jsonFilePath string) (string, error) {
	f, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return "", err
	}
	return string(f), nil
}

type bytecodeParsingStruct struct {
}

func BytecodeParsing(jsonFilePath string) (string, error) {
	var bytecodeStruct bytecodeParsingStruct
	jsonData, err := JsonUnmarshal(bytecodeStruct, jsonFilePath)
	if err != nil {
		return "", err
	}
	bytecode := jsonData.(map[string]interface{})["object"].(string)

	return bytecode, nil
}

// For invoke(as execute) contract, parameters are packed by using ABI.
func GetAbiPack(callName string, abi string, bytecode string, args ...interface{}) ([]byte, error) {
	metadata := GetBindMetaData(abi, bytecode)
	contractAbi, err := metadata.GetAbi()
	if err != nil {
		return nil, err
	}

	var abiByteData []byte
	if args == nil {
		abiByteData, err = contractAbi.Pack(callName)
		if err != nil {
			return nil, err
		}
	} else {
		abiByteData, err = contractAbi.Pack(callName, args...)
		if err != nil {
			return nil, err
		}
	}

	return abiByteData, nil
}

// After call(as query) solidity contract, the response of chain is unpacked by ABI.
func GetAbiUnpack(callName string, abi string, bytecode string, data []byte) ([]interface{}, error) {
	metadata := GetBindMetaData(abi, bytecode)
	contractAbi, err := metadata.GetAbi()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	unpacked, err := contractAbi.Unpack(callName, data)
	if err != nil {
		return nil, err
	}

	return unpacked, nil
}

func GetBindMetaData(abi, bytecode string) *bind.MetaData {
	return &bind.MetaData{
		ABI: abi,
		Bin: bytecode,
	}
}
