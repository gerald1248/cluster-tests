package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	au "github.com/logrusorgru/aurora"
)

func runTests(datadir string) error {
	tests, err := filepath.Glob(fmt.Sprintf("%s/*test", datadir))
	if err != nil {
		return err
	}

	var successCount, failureCount int
	successCount = 0
	failureCount = 0
	for _, match := range tests {
		t := time.Now()
		basename := filepath.Base(match)

		os.Setenv("USER_NAMESPACES", "default kube-public kube-system")
		os.Setenv("NODES", "minikube")
		os.Setenv("HA_SERVICES", "")

		fmt.Printf("[%s] %s... ", t.Format("2006-01-02 15:04:05"), au.Bold(au.Cyan(basename)))
		var stdout, stderr bytes.Buffer
		cmd := exec.Command("bash", match)
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		// ignore stderr for now

		err = cmd.Run()
		if err != nil {
			message := strings.TrimRight(string(stdout.Bytes()), " \n")
			if len(message) == 0 {
				fmt.Printf("%s\n", au.Bold(au.Red("failed")))
			} else {
				fmt.Printf("%s: %s\n", au.Bold(au.Red("failed")), message)
			}
			failureCount++
		} else {
			fmt.Printf("%s\n", au.Bold(au.Green("ok")))
			successCount++
		}
	}

	total := failureCount + successCount
	plural := "s"
	if total == 1 {
		plural = ""
	}
	fmt.Printf("Ran %d test%s\n", total, plural)

	if failureCount == 0 {
		fmt.Printf("%s\n", au.Bold(au.Green("OK")))
	} else {
		fmt.Printf("%s (failures=%d)\n", au.Bold(au.Red("FAILED")), failureCount)
	}
	return nil
}
