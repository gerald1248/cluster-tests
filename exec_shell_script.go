package main

import (
	"bytes"
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
