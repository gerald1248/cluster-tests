package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	au "github.com/logrusorgru/aurora"
)

var globalDatadir, globalOutputdir, globalContext string

// run tests
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	server := flag.String("s", "localhost", "server")
	port := flag.Int("p", 8080, "listen on port")
	datadir := flag.String("d", "cluster-tests.d", "data directory")
	outputdir := flag.String("o", "output", "output directory")
	interval := flag.Int("i", 30, "interval (s)")
	retain := flag.Int("r", 2, "retain (d)")

	globalDatadir = *datadir
	globalOutputdir = *outputdir

	context, _, err := execShellScript(fmt.Sprintf("%s/get_context", *datadir))
	if err != nil {
		fmt.Printf("%s: can't fetch context (%s)\n", au.Bold(au.Red("Error")), err.Error())
	}

	globalContext = context

	flag.Parse()

	ticker := time.NewTicker(time.Millisecond * 1000 * time.Duration(*interval))

	// trigger initial run
	go func() {
		err = runTests(*datadir, *outputdir, *retain)
		if err != nil {
			fmt.Printf("%s: %s\n", au.Bold(au.Red("Error")), err.Error())
		}
	}()

	// schedule subsequent runs
	go func() {
		for range ticker.C {
			err = runTests(*datadir, *outputdir, *retain)
			if err != nil {
				fmt.Printf("%s: %s\n", au.Bold(au.Red("Error")), err.Error())
			}
		}
	}()

	serve(*server, *port)
	return
}
