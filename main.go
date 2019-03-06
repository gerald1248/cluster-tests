package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// run tests defined in datadir
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
	interval := flag.Int64("i", 3600, "interval (s)")
	retain := flag.Int("r", 2, "retain (d)")
	errors := flag.Bool("e", false, "output stderr")
	name := flag.String("n", "", "context name")
	duration := flag.Bool("-duration", true, "display duration chart")
	histogram := flag.Bool("-histogram", false, "display histogram")

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

	ticker := time.NewTicker(time.Millisecond * time.Duration(1000) * time.Duration(*interval))

	runTestsParam := RunTestsParam{*datadir, *outputdir, context, *retain, *errors, *duration, *histogram}

	// trigger initial run
	go func() {
		callRunTests(runTestsParam)
	}()

	// schedule subsequent runs
	go func() {
		for range ticker.C {
			callRunTests(runTestsParam)
		}
	}()

	serve(*server, *port, *outputdir, context)
	return
}
