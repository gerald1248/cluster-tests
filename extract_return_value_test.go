package main

import (
	"fmt"
	"testing"
)

func TestExtractReturnValue(t *testing.T) {
	testVal := 99
	i, err := extractReturnValue(fmt.Errorf("exit status %d", testVal))

	if err != nil {
		t.Errorf("return value extraction failed")
		return
	}

	if i != testVal {
		t.Errorf("expected extracted value %d, got %d", testVal, i)
	}
}
