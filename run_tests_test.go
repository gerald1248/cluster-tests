package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestRunTests(t *testing.T) {
	datadir := "testdata/cluster-tests.d"
	outputdir := "testdata/output"
	context := "test"
	retain := 0
	param := RunTestsParam{
		datadir,
		outputdir,
		context,
		retain,
		false,
		false,
		false,
	}

	prepareTest(param)
	err := runTests(param)

	if err != nil {
		t.Errorf("Sample test should run without errors: %s", err.Error())
	}
}

func TestRunTestsInvalidOutputdir(t *testing.T) {
	datadir := "testdata/cluster-tests.d"
	outputdir := "nonexistent/subfolder"
	context := "test"
	retain := 0
	param := RunTestsParam{
		datadir,
		outputdir,
		context,
		retain,
		false,
		false,
		false,
	}

	prepareTest(param)
	err := runTests(param)

	if err == nil {
		t.Error("Sample test with invalid output directory should report error")
	}
}

func prepareTest(param RunTestsParam) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.json", param.outputdir))
	if err != nil {
		return
	}

	for _, file := range files {
		os.Remove(file)
	}
}
