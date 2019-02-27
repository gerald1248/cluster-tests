package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

func getMetrics() (int, int, bool, error) {
	files, err := filepath.Glob(fmt.Sprintf("%s/*.json", globalOutputdir))
	if err != nil {
		return 0, 0, false, fmt.Errorf("can't glob output files (%s)", err.Error())
	}

	var passed, failed int
	for i, file := range files {
		if i != len(files)-1 {
			continue
		}

		var record Record
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return 0, 0, false, fmt.Errorf("can't read %s (%s)", file, err.Error())
		}

		err = json.Unmarshal(b, &record)

		if err != nil {
			return 0, 0, false, fmt.Errorf("invalid JSON (%s)", err.Error())
		}

		passed = record.Pass
		failed = record.Fail
	}
	return passed, failed, (failed == 0), nil
}
