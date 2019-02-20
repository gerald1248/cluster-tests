package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// run tests
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	server := flag.String("s", "localhost", "server")
	port := flag.Int("p", 8080, "listen on port")
	datadir := flag.String("d", "./cluster-tests.d", "data directory")
	interval := flag.Int("i", 5, "interval (s)")

	flag.Parse()

	ticker := time.NewTicker(time.Millisecond * 1000 * time.Duration(*interval))
	runTests(*datadir)
	go func() {
		for range ticker.C {
			runTests(*datadir)
		}
	}()

	serve(*server, *port)
	return
}
