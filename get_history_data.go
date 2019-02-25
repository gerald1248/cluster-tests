package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func getHistoryData() ([]byte, []byte, int, []string, string, bool, error) {
	failed := false
	files, err := filepath.Glob(fmt.Sprintf("%s/*.json", globalOutputdir))
	if err != nil {
		return nil, nil, 0, nil, "", failed, fmt.Errorf("can't glob test files (%s)", err.Error())
	}

	var items []VegaLiteItem
	var durationItems []VegaLiteDurationItem
	var logEntries []string
	var logHead string
	maxTests := 0
	for i, file := range files {
		var record Record
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, nil, 0, nil, "", failed, fmt.Errorf("can't read %s (%s)", file, err.Error())
		}

		err = json.Unmarshal(b, &record)

		if err != nil {
			return nil, nil, 0, nil, "", failed, fmt.Errorf("invalid JSON (%s)", err.Error())
		}

		// populate main result history
		items = append(items, VegaLiteItem{record.Pass, record.Time, "PASS"})
		items = append(items, VegaLiteItem{record.Fail, record.Time, "FAIL"})

		// populate
		durationItems = append(durationItems, VegaLiteDurationItem{record.Duration, record.Time})

		sum := record.Fail + record.Pass
		if sum > maxTests {
			maxTests = sum
		}

		if i == len(files)-1 {
			for _, log := range record.FailLog {
				logEntries = append(logEntries, log)
			}
			for _, log := range record.PassLog {
				logEntries = append(logEntries, log)
			}
			logHead = record.Head
			failed = len(record.FailLog) > 0
		}
	}

	jsonResults, err := json.Marshal(items)
	if err != nil {
		return nil, nil, 0, nil, "", failed, fmt.Errorf("Can't marshal result items (%s)", err.Error())
	}

	jsonDurations, err := json.Marshal(durationItems)
	if err != nil {
		return nil, nil, 0, nil, "", failed, fmt.Errorf("Can't marshal duration items (%s)", err.Error())
	}

	return jsonResults, jsonDurations, maxTests, logEntries, logHead, failed, nil
}
