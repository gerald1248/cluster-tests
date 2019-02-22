package main

import (
	"fmt"
	"io/ioutil"
)

func processBytes(byteArray []byte, output *string) (string, error) {

	//preflight with optional conversion from YAMLs
	err := preflightAsset(&byteArray)
	if err != nil {
		return "", fmt.Errorf("input failed preflight check: %v", err)
	}

	// make sure config objects are presented as a list
	err = makeList(&byteArray)
	if err != nil {
		return "", err
	}

	return "", nil
}

func processFile(path string, output *string) (string, error) {
	byteArray, err := ioutil.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("can't read %s: %v", path, err)
	}

	result, err := processBytes(byteArray, output)

	if err != nil {
		return "", fmt.Errorf("can't process %s: %s", path, err)
	}

	return result, nil
}
