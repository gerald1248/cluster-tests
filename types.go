// ingress/egress rule types based on:
// https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/networking/types.go

package main

// Record captures the internal representation of a test run
type Record struct {
	Fail      int            `json:"fail"`
	Pass      int            `json:"pass"`
	Duration  int            `json:"duration"`
	Time      string         `json:"time"`
	FailLog   []string       `json:"failLog"`
	PassLog   []string       `json:"passLog"`
	Head      string         `json:"head"`
	Histogram map[string]int `json:"histogram"`
}

// VegaLiteItem wraps a data point in the main Vega-Lite config
type VegaLiteItem struct {
	Tests  int    `json:"tests"`
	Time   string `json:"time"`
	Result string `json:"result"`
}

// VegaLiteDurationItem wraps a data point in the duration Vega-Lite config
type VegaLiteDurationItem struct {
	Duration float64 `json:"duration"`
	Time     string  `json:"time"`
}

// VegaLiteHistogramItem wraps a data point in the histogram Vega-Lite config
type VegaLiteHistogramItem struct {
	Test  string `json:"test"`
	Count int    `json:"count"`
}

// Metrics describes the return JSON object of a metrics query
type Metrics struct {
	Fail int `json:"fail"`
	Pass int `json:"pass"`
}

// RunTestsParam gathers up the parameters passed to the test runner goroutine
type RunTestsParam struct {
	datadir   string
	outputdir string
	context   string
	retain    int
	errors    bool
	duration  bool
	histogram bool
	cache     bool
}

// ParsedHistory prevents a lengthy list of return values from getHistoryData()
type ParsedHistory struct {
	jsonResults   []byte
	jsonDurations []byte
	jsonHistogram []byte
	maxTests      int
	logEntries    []string
	logHead       string
	lastRecord    Record
}

// MinimalObject is a placeholder struct for Kubernetes manifests
type MinimalObject struct {
	Kind string
}

// Result TODO remove
type Result struct {
	PercentageIsolated          int `json:"percentageIsolated"`
	PercentageNamespaceCoverage int `json:"percentageNamespaceCoverage"`
}
