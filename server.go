package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"
)

type PostStruct struct {
	Buffer string
}

func serve(server string, port int) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/api/", apiHandler)
	mux.HandleFunc("/api/v1/", apiHandler)
	mux.HandleFunc("/health/", healthHandler)
	mux.HandleFunc("/api/v1/metrics/", metricsHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", server, port), mux))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handlePost(&w, r)
	case "GET":
		handleGet(&w, r)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	buffer := fmt.Sprintf(`
<div class="row">
  <div class="col-sm-2"><a href="/">/</a></div>
  <div class="col-sm-10">show graph</div>
</div>
<div class="row">
  <div class="col-sm-2"><a href="/health/">/health/</a></div>
  <div class="col-sm-10">health endpoint</div>
</div>
<div class="row">
  <div class="col-sm-2"><a href="/api/v1/metrics/">/api/v1/metrics/</a></div>
  <div class="col-sm-10">metrics endpoint</div>
</div>`)
	fmt.Fprintf(w, page(buffer))
	return
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"status\":\"ok\"}")
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	var buffer, output string

	output = "dot"

	_, err := processBytes([]byte(buffer), &output)
	if err != nil {
		sData := fmt.Sprintf("<p>Can't process input data: %s</p>", err)
		fmt.Fprintf(w, page(sData))
		return
	}

	percentageIsolated := 0
	percentageNamespaceCoverage := 0
	fmt.Fprintf(w, "{\"percentageIsolated\":%d,\"percentageNamespaceCoverage\":%d}", percentageIsolated, percentageNamespaceCoverage)
}

func handleGet(w *http.ResponseWriter, r *http.Request) {
	var buffer, dot, output string

	output = "dot"

	dot, err := processBytes([]byte(buffer), &output)
	if err != nil {
		sData := fmt.Sprintf("<p>Can't process input data: %s</p>", err)
		fmt.Fprintf(*w, page(sData))
		return
	}

	cmd := exec.Command("dot", "-Tsvg")
	cmd.Stdin = strings.NewReader(dot)
	var svg bytes.Buffer
	cmd.Stdout = &svg
	err = cmd.Run()
	if err != nil {
		sDot := fmt.Sprintf("<p>Graphviz conversion failed: %s</p>", err)
		fmt.Fprintf(*w, page(sDot))
		return
	}

	percentageIsolated := 0
	percentageCovered := 0

	colorClassIsolation := "progress-bar-success"
	if percentageIsolated < 50 {
		colorClassIsolation = "progress-bar-danger"
	} else if percentageIsolated < 75 {
		colorClassIsolation = "progress-bar-warning"
	}

	colorClassCoverage := "progress-bar-success"
	if percentageCovered < 50 {
		colorClassCoverage = "progress-bar-danger"
	} else if percentageCovered < 75 {
		colorClassCoverage = "progress-bar-warning"
	}

	buffer = fmt.Sprintf(`
<div>%s</div>
<div class="progress">
  <div class="progress-bar %s" style="width: %d%%%%" role="progressbar" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="100">%d%%%% isolation</div>
</div>
<div class="progress">
  <div class="progress-bar %s" style="width: %d%%%%" role="progressbar" aria-valuenow="%d" aria-valuemin="0" aria-valuemax="100">%d%%%% namespace coverage</div>
</div>`,
		strings.Replace(svg.String(), "Times,serif", "sans-serif", -1),
		colorClassIsolation,
		percentageIsolated,
		percentageIsolated,
		percentageIsolated,
		colorClassCoverage,
		percentageCovered,
		percentageCovered,
		percentageCovered)
	fmt.Fprintf(*w, page(buffer))
}

func handlePost(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(*w, "Can't read POST request body '%s': %s", body, err)
		return
	}
}