package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func getIgnoreSet(datadir string) map[string]int {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/ignore", datadir))
	if err != nil {
		fmt.Printf("can't read ignore list in %s", datadir)
		return map[string]int{}
	}

	files := strings.Split(string(bytes), "\n")

	ignoreSet := map[string]int{}
	for _, file := range files {
		ignoreSet[file] = 1
	}

	return ignoreSet
}
