package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	au "github.com/logrusorgru/aurora"
)

func runTests(datadir string, outputdir string, retain int) error {
	// cleanup first
	purgeOutput(outputdir, retain)

	tests, err := filepath.Glob(fmt.Sprintf("%s/test*", datadir))
	if err != nil {
		return fmt.Errorf("can't glob test files (%s)", err.Error())
	}

	userNamespaces, _, err := execShellScript(fmt.Sprintf("%s/get_user_namespaces", datadir))
	if err != nil {
		return fmt.Errorf("can't determine user namespaces (%s)", err.Error())
	}
	nodes, _, err := execShellScript(fmt.Sprintf("%s/get_nodes", datadir))
	if err != nil {
		return fmt.Errorf("can't fetch cluster nodes (%s)", err.Error())
	}

	var successCount, failureCount, maxCount int
	successCount = 0
	failureCount = 0
	maxCount = 0

	var record Record
	record.Histogram = map[string]int{} // initialise map

	startTime := time.Now()

	for _, match := range tests {
		t := time.Now()
		basename := filepath.Base(match)

		os.Setenv("USER_NAMESPACES", userNamespaces)
		os.Setenv("NODES", nodes)
		os.Setenv("HA_SERVICES", "")
		os.Setenv("CLUSTER_TESTS_EXIT", "")

		fmt.Printf("[%s] %s... ", t.Format("2006-01-02 15:04:05"), au.Bold(au.Cyan(basename)))

		stdout, _, err := execShellScript(match)

		if err != nil {
			message := strings.TrimRight(stdout, " \n")
			if len(message) == 0 {
				fmt.Printf("%s\n", au.Bold(au.Red("failed")))
			} else {
				fmt.Printf("%s: %s\n", au.Bold(au.Red("failed")), message)
			}
			failureCount++

			// append to failure log
			record.FailLog = append(record.FailLog, fmt.Sprintf("%s %s %s", basename, au.Bold(au.Red("failed")), au.Bold(au.Cyan(message))))

			// update histogram
			if value, ok := record.Histogram[basename]; ok {
				record.Histogram[basename] = value + 1
			} else {
				record.Histogram[basename] = 1
			}
		} else {
			fmt.Printf("%s\n", au.Bold(au.Green("ok")))
			successCount++

			record.PassLog = append(record.PassLog, fmt.Sprintf("%s %s", basename, au.Bold(au.Green("ok"))))
		}
	}

	recordTime := time.Now()

	record.Time = fmt.Sprintf("%s", recordTime.Format("2006-01-02 15:04:05"))
	record.Fail = failureCount
	record.Pass = successCount

	record.Duration = int(time.Since(startTime).Nanoseconds() / 1000)

	recordFilename := fmt.Sprintf("%s/%d.json", globalOutputdir, recordTime.Unix())

	total := failureCount + successCount
	if total > maxCount {
		maxCount = total
	}
	plural := "s"
	if total == 1 {
		plural = ""
	}
	fmt.Printf("Ran %d test%s\n", total, plural)

	if failureCount == 0 {
		record.Head = fmt.Sprintf("%s\n", au.Bold(au.Green("OK")))
	} else {
		record.Head = fmt.Sprintf("%s (failures=%d)\n", au.Bold(au.Red("FAILED")), failureCount)
	}
	fmt.Printf("%s", record.Head)

	recordJSON, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("can't marshal record (%s)", err.Error())
	}
	err = ioutil.WriteFile(recordFilename, recordJSON, 0644)
	if err != nil {
		return fmt.Errorf("can't write record to file %s (%s)", recordFilename, err.Error())
	}

	return nil
}
