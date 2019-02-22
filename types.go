// ingress/egress rule types based on:
// https://github.com/kubernetes/kubernetes/blob/master/pkg/apis/networking/types.go

package main

// Record captures the internal representation of a test run
type Record struct {
	Fail     int      `json:"fail"`
	Pass     int      `json:"pass"`
	Duration int      `json:"duration"`
	Time     string   `json:"time"`
	FailLog  []string `json:"failLog"`
	PassLog  []string `json:"passLog"`
}

// VegaLiteItem wraps a data point in a Vega-Lite config block
type VegaLiteItem struct {
	Tests  int    `json:"tests"`
	Time   string `json:"time"`
	Result string `json:"result"`
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
