package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	term "github.com/buildkite/terminal"
)

// PostStruct wraps minimal POST requests
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
  <div class="col-sm-10">show dashboard</div>
</div>
<div class="row">
  <div class="col-sm-2"><a href="/health/">/health/</a></div>
  <div class="col-sm-10">health endpoint</div>
</div>
<div class="row">
  <div class="col-sm-2"><a href="/api/v1/metrics/">/api/v1/metrics/</a></div>
  <div class="col-sm-10">metrics endpoint</div>
</div>`)
	fmt.Fprintf(w, pageMinimal(globalContext, buffer))
	return
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "{\"status\":\"ok\"}")
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	passed, failed, _, err := getMetrics()
	if err != nil {
		fmt.Fprintf(w, "{}")
		return
	}

	result := "OK"
	if failed > 0 {
		result = "FAILED"
	}
	fmt.Fprintf(w, "{\"pass\":%d,\"fail\":%d,\"result\":\"%s\"}", passed, failed, result)
}

func handleGet(w *http.ResponseWriter, r *http.Request) {

	vegaLiteDataBytes, vegaLiteDurationBytes, vegaLiteHistogramBytes, maxTests, logEntries, logHead, failed, err := getHistoryData()
	if err != nil {
		sData := fmt.Sprintf("<p>Can't display dashboard: %s</p>", err.Error())
		fmt.Fprintf(*w, pageMinimal(globalContext, sData))
		return
	}

	terminal := logHead + strings.Join(logEntries, "\n")

	terminalBytes := []byte(terminal)

	var chart01 = fmt.Sprintf(staticTextVis01, vegaLiteDataBytes, maxTests)
	var chart02 = fmt.Sprintf(staticTextVis02, vegaLiteDurationBytes)
	var chart03 = fmt.Sprintf(staticTextVis03, vegaLiteHistogramBytes)

	log := fmt.Sprintf(`<div class="term-container">%s</div>`, string(term.Render(terminalBytes)))

	bgColorClass := "bg-secondary"
	if failed {
		bgColorClass = "bg-danger"
	}
	fmt.Fprintf(*w, page(globalContext, chart01, chart02, chart03, log, bgColorClass))
}

func handlePost(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(*w, "Can't read POST request body '%s': %s", body, err)
		return
	}
}
