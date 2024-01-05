package util

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
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

func JsonUnmarshalData(jsonStruct interface{}, byteValue []byte) (interface{}, error) {
	err := json.Unmarshal(byteValue, &jsonStruct)
	if err != nil {
		return nil, err
	}

	return jsonStruct, nil
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

	return JsonUnmarshalData(jsonStruct, byteValue)
}

func SaveJsonPretty(jsonByte []byte, saveTxPath string) error {
	var prettyJson bytes.Buffer
	err := json.Indent(&prettyJson, jsonByte, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(saveTxPath, prettyJson.Bytes(), 0660)
	if err != nil {
		return err
	}

	return nil
}
