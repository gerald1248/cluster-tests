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
	fmt.Fprintf(w, page("cluster ID", "", buffer))
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
		fmt.Fprintf(w, page("cluster ID", "", sData))
		return
	}

	percentageIsolated := 0
	percentageNamespaceCoverage := 0
	fmt.Fprintf(w, "{\"percentageIsolated\":%d,\"percentageNamespaceCoverage\":%d}", percentageIsolated, percentageNamespaceCoverage)
}

func handleGet(w *http.ResponseWriter, r *http.Request) {

	vegaLiteDataBytes, maxTests, logEntries, err := getHistoryData()
	if err != nil {
		// empty byte array returned
		fmt.Sprintf("can't fetch JSON data (%s)", err.Error())
	}

	terminal := strings.Join(logEntries, "\n")

	terminalBytes := []byte(terminal)

	adjust := int(maxTests / 2)

	var chart = fmt.Sprintf(`
    <div id="vis"></div>

    <script type="text/javascript">
	var yourVlSpec = {
		"$schema": "https://vega.github.io/schema/vega-lite/v2.0.json",
		"data": {
		  "values": %s
		},
		"width": 400,
		"height": 200,
		"mark": "area",
		"encoding": {
		  "x": {
			"field": "time",
			"type": "temporal"
		  },
		  "y": {
			"field": "tests",
			"type": "quantitative",
			"scale": {
			  "domain": [
				0,
				%d
			  ]
			}
		  },
		  "color": {
			"field": "result",
			"type": "nominal",
			"legend": {
				"labelColor": "#fff",
				"titleColor": "#fff"
			},
			"scale": {
			  "domain": [
				"PASS",
				"FAIL"
			  ],
			  "range": [
				"#2ECC40",
				"#FF4136"
			  ]
			}
		  }
		},
		"config": {
		  "axis": {
			"labelFont": "sans-serif",
			"titleFont": "sans-serif",
			"labelColor": "white",
			"titleColor": "white"
		  },
		  "axisX": {
			"labelAngle": 0
		  }
		}
	  }
	  vegaEmbed('#vis', yourVlSpec);
	</script>`, vegaLiteDataBytes, maxTests+adjust) // alternatively, use staticTextVegaLiteData

	log := fmt.Sprintf(`<div class="term-container">%s</div>`, string(term.Render(terminalBytes)))

	// TODO: fetch actual context
	fmt.Fprintf(*w, page(globalContext, chart, log))
}

func handlePost(w *http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(*w, "Can't read POST request body '%s': %s", body, err)
		return
	}
}
