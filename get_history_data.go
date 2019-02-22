package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func getHistoryData() ([]byte, int, []string, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.json", globalOutputdir))
	if err != nil {
		return nil, 0, nil, fmt.Errorf("can't glob test files (%s)", err.Error())
	}

	var items []VegaLiteItem
	var logEntries []string
	maxTests := 0
	for i, file := range files {
		var record Record
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, 0, nil, fmt.Errorf("can't read %s (%s)", file, err.Error())
		}

		err = json.Unmarshal(b, &record)

		if err != nil {
			return nil, 0, nil, fmt.Errorf("invalid JSON (%s)", err.Error())
		}

		// populate
		items = append(items, VegaLiteItem{record.Fail, record.Time, "FAIL"})
		items = append(items, VegaLiteItem{record.Pass, record.Time, "PASS"})

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
		}
	}

	json, err := json.Marshal(items)
	if err != nil {
		return nil, 0, nil, fmt.Errorf("Can't marshal Vega Lite items (%s)", err.Error())
	}

	return json, maxTests, logEntries, nil
}
