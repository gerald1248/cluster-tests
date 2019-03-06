package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	term "github.com/buildkite/terminal"
	au "github.com/logrusorgru/aurora"
)

// callRunTests is extracted from main() as it is called twice
// don't panic on error
func callRunTests(param RunTestsParam) {
	// run tests first
	err := runTests(param)
	if err != nil {
		fmt.Printf("%s: %s\n", au.Bold(au.Red("Error")), err.Error())
		return
	}

	// then write out dashboard as static page
	var pageBuffer string

	history, err := getHistoryData(param.outputdir)
	if err != nil {
		sData := fmt.Sprintf("<p>Can't display dashboard: %s</p>", err.Error())
		pageBuffer = fmt.Sprintf(pageMinimal(param.context, sData))
		return
	}

	timeSummary := fmt.Sprintf("%s (%s)", history.lastRecord.Time, formatDuration(int64(history.lastRecord.Duration)))

	terminal := history.logHead
	terminal += strings.Join(history.logEntries, "\n")
	terminal += "\n\n"
	terminal += fmt.Sprintf("%s", au.Bold(au.Gray(timeSummary)))
	terminal += "\n"

	terminalBytes := []byte(terminal)

	var chart01, chart02, chart03 string
	chart01 = fmt.Sprintf(staticTextVis01, history.jsonResults, history.maxTests)

	if param.duration {
		chart02 = fmt.Sprintf(staticTextVis02, history.jsonDurations)
	}

	if param.histogram {
		chart03 = fmt.Sprintf(staticTextVis03, history.jsonHistogram)
	}

	log := fmt.Sprintf(`<div class="term-container">%s</div>`, string(term.Render(terminalBytes)))

	bgColorClass := "bg-secondary"
	if history.lastRecord.Fail > 0 {
		bgColorClass = "bg-danger"
	}
	pageBuffer = fmt.Sprintf(page(param.context, chart01, chart02, chart03, log, bgColorClass))

	filename := fmt.Sprintf("%s/index.html", param.outputdir)
	err = ioutil.WriteFile(filename, []byte(pageBuffer), 0644)
	if err != nil {
		fmt.Printf("can't write index file %s (%s)", filename, err.Error())
	}
}

func runTests(param RunTestsParam) error {

	// cleanup first
	purgeOutput(param.outputdir, param.retain)

	// update ignore list
	ignoreSet := getIgnoreSet(param.datadir)

	tests, err := filepath.Glob(fmt.Sprintf("%s/test*", param.datadir))
	if err != nil {
		return fmt.Errorf("can't glob test files (%s)", err.Error())
	}

	userNamespaces, _, err := execShellScript(fmt.Sprintf("%s/get_user_namespaces", param.datadir))
	if err != nil {
		return fmt.Errorf("can't determine user namespaces (%s)", err.Error())
	}
	nodes, _, err := execShellScript(fmt.Sprintf("%s/get_nodes", param.datadir))
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
		matchBasename := strings.TrimPrefix(match, fmt.Sprintf("%s/", param.datadir))
		if _, ok := ignoreSet[matchBasename]; ok {
			continue
		}

		t := time.Now()
		basename := filepath.Base(match)

		os.Setenv("USER_NAMESPACES", userNamespaces)
		os.Setenv("NODES", nodes)
		os.Setenv("HA_SERVICES", "")
		os.Setenv("CLUSTER_TESTS_EXIT", "")

		fmt.Printf("[%s] %s... ", t.Format("2006-01-02 15:04:05"), au.Bold(au.Cyan(basename)))

		stdout, stderr, err := execTestShellScript(match)

		if err != nil {
			message := strings.TrimRight(stdout, " \n")

			if param.errors {
				message = fmt.Sprintf("%s %s", stderr, message)
			}

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

	recordFilename := fmt.Sprintf("%s/%d.json", param.outputdir, recordTime.Unix())

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
