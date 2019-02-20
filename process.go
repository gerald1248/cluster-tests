package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func processBytes(byteArray []byte, output *string) (string, error) {

	//preflight with optional conversion from YAMLs
	err := preflightAsset(&byteArray)
	if err != nil {
		return "", errors.New(fmt.Sprintf("input failed preflight check: %v", err))
	}

	// make sure config objects are presented as a list
	err = makeList(&byteArray)
	if err != nil {
		return "", err
	}

	var apiObjectSet ApiObjectSet

	if err = json.Unmarshal(byteArray, &apiObjectSet); err != nil {
		return "", errors.New(fmt.Sprintf("can't unmarshal data: %v", err))
	}

	return string(byteArray), nil
}

func processFile(path string, output *string) (string, error) {
	byteArray, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.New(fmt.Sprintf("can't read %s: %v", path, err))
	}

	result, err := processBytes(byteArray, output)

	if err != nil {
		return "", errors.New(fmt.Sprintf("can't process %s: %s", path, err))
	}

	return result, nil
}
