package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func getHistoryData() (ParsedHistory, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.json", globalOutputdir))
	if err != nil {
		return ParsedHistory{}, fmt.Errorf("can't glob output files (%s)", err.Error())
	}

	lastRecord := Record{}

	cumulative := map[string]int{}
	var items []VegaLiteItem
	var durationItems []VegaLiteDurationItem
	var histogramItems []VegaLiteHistogramItem
	var logEntries []string
	var logHead string
	maxTests := 0
	for i, file := range files {
		var record Record
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return ParsedHistory{}, fmt.Errorf("can't read %s (%s)", file, err.Error())
		}

		err = json.Unmarshal(b, &record)

		if err != nil {
			return ParsedHistory{}, fmt.Errorf("invalid JSON (%s)", err.Error())
		}

		// populate main result history
		items = append(items, VegaLiteItem{record.Pass, record.Time, "PASS"})
		items = append(items, VegaLiteItem{record.Fail, record.Time, "FAIL"})

		// populate duration history
		durationItems = append(durationItems, VegaLiteDurationItem{float64(record.Duration) / 1000000.0, record.Time})

		// update histogram map
		for k, v := range record.Histogram {
			// update histogram
			if recordValue, ok := cumulative[k]; ok {
				cumulative[k] = recordValue + v
			} else {
				cumulative[k] = v
			}
		}

		sum := record.Fail + record.Pass
		if sum > maxTests {
			maxTests = sum
		}

		if i == len(files)-1 {
			lastRecord = record
			for _, log := range record.FailLog {
				logEntries = append(logEntries, log)
			}
			for _, log := range record.PassLog {
				logEntries = append(logEntries, log)
			}
			logHead = record.Head
		}
	}

	// transfer histogram map to vega-lite array
	for k, v := range cumulative {
		histogramItems = append(histogramItems, VegaLiteHistogramItem{k, v})
	}

	jsonResults, err := json.Marshal(items)
	if err != nil {
		return ParsedHistory{}, fmt.Errorf("Can't marshal result items (%s)", err.Error())
	}

	jsonDurations, err := json.Marshal(durationItems)
	if err != nil {
		return ParsedHistory{}, fmt.Errorf("Can't marshal duration items (%s)", err.Error())
	}

	jsonHistogram, err := json.Marshal(histogramItems)
	if err != nil {
		return ParsedHistory{}, fmt.Errorf("Can't marshal histogram items (%s)", err.Error())
	}

	history := ParsedHistory{
		jsonResults,
		jsonDurations,
		jsonHistogram,
		maxTests,
		logEntries,
		logHead,
		lastRecord,
	}
	return history, nil
}
