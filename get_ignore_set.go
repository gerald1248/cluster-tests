package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	au "github.com/logrusorgru/aurora"
)

func getIgnoreSet(datadir string) map[string]int {
	bytes, err := ioutil.ReadFile(fmt.Sprintf("%s/ignore", datadir))
	if err != nil {
		fmt.Printf("%s: can't read ignore list in %s\n", au.Bold(au.Red("Error")), datadir)
		return map[string]int{}
	}

	files := strings.Split(string(bytes), "\n")

	ignoreSet := map[string]int{}
	for _, file := range files {
		ignoreSet[file] = 1
	}

	return ignoreSet
}
