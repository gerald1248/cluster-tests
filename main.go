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

	server := flag.String("s", "", "server")
	port := flag.Int("p", 8080, "listen on port")
	datadir := flag.String("d", "cluster-tests.d", "data directory")
	outputdir := flag.String("o", "output", "output directory")
	interval := flag.Int("i", 3600, "interval (s)")
	retain := flag.Int("r", 2, "retain (d)")
	errors := flag.Bool("e", false, "output stderr")
	name := flag.String("n", "", "context name")

	globalDatadir = *datadir
	globalOutputdir = *outputdir

	flag.Parse()

	// identify context/cluster name
	var context string
	// option #1: name param
	if len(*name) > 0 {
		context = *name
	} else { // option #2: context is available
		buffer, _, err := execShellScript(fmt.Sprintf("%s/get_context", *datadir))
		if err == nil {
			context = buffer
		} else {
			context = os.Getenv("CLUSTER_TESTS_CONTEXT") // option #3: custom context variable
			if len(context) == 0 {
				context = os.Getenv("KUBERNETES_PORT_443_TCP_ADDR") // option #4: IP address
			}
		}
	}

	globalContext = context
	ticker := time.NewTicker(time.Millisecond * 1000 * time.Duration(*interval))

	// trigger initial run
	go func() {
		err := runTests(*datadir, *outputdir, *retain, *errors)
		if err != nil {
			fmt.Printf("%s: %s\n", au.Bold(au.Red("Error")), err.Error())
		}
	}()

	// schedule subsequent runs
	go func() {
		for range ticker.C {
			err := runTests(*datadir, *outputdir, *retain, *errors)
			if err != nil {
				fmt.Printf("%s: %s\n", au.Bold(au.Red("Error")), err.Error())
			}
		}
	}()

	serve(*server, *port)
	return
}
