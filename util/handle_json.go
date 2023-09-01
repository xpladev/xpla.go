package util

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/xpladev/xpla.go/types/errors"
)

func JsonMarshalData(jsonData interface{}) ([]byte, error) {
	byteData, err := json.Marshal(jsonData)
	if err != nil {
		return nil, err
	}

	return byteData, nil
}

func JsonMarshalDataIndent(jsonData interface{}) ([]byte, error) {
	byteData, err := json.MarshalIndent(jsonData, "", "    ")
	if err != nil {
		return nil, err
	}

	return byteData, nil
}

func JsonUnmarshalData(jsonStruct interface{}, byteValue []byte) interface{} {
	json.Unmarshal(byteValue, &jsonStruct)

	return jsonStruct
}

func JsonUnmarshal(jsonStruct interface{}, jsonFilePath string) (interface{}, error) {
	jsonData, err := os.Open(jsonFilePath)
	if err != nil {
		return nil, err
	}
	byteValue, err := io.ReadAll(jsonData)
	if err != nil {
		return nil, err
	}
	jsonStruct = JsonUnmarshalData(jsonStruct, byteValue)

	return jsonStruct, nil
}

func SaveJsonPretty(jsonByte []byte, saveTxPath string) error {
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, jsonByte, "", "    ")
	if err != nil {
		return LogErr(errors.ErrFailedToMarshal, err)
	}

	err = os.WriteFile(saveTxPath, prettyJson.Bytes(), 0660)
	if err != nil {
		return LogErr(errors.ErrFailedToMarshal, err)
	}

	return nil
}
