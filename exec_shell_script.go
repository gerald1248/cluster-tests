package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
)

func execShellScript(path string) (string, string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("bash", path)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	return stdout.String(), stderr.String(), err
}

func execTestShellScript(path string) (string, string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("bash")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	testBuffer, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Can't read test file %s\n", path)
		return "", "", nil
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, fmt.Sprintf("%s\n%s\n%s\n", staticTextTestStart, string(testBuffer), staticTextTestEnd))
	}()

	err = cmd.Run()

	return stdout.String(), stderr.String(), err
}
