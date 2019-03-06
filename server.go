package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// PostStruct wraps minimal POST requests
type PostStruct struct {
	Buffer string
}

func serve(server string, port int, outputdir string, context string) {
	http.Handle("/", handler(outputdir, context))
	http.Handle("/api/", apiHandler(outputdir, context))
	http.Handle("/api/v1/", apiHandler(outputdir, context))
	http.Handle("/health/", healthHandler(outputdir, context))
	http.Handle("/api/v1/metrics/", metricsHandler(outputdir, context))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", server, port), nil))
}

func handler(outputdir string, context string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			indexfile := fmt.Sprintf("%s/index.html", outputdir)
			buffer, err := ioutil.ReadFile(indexfile)

			if err != nil {
				fmt.Printf("Can't read index file %s\n", indexfile)
				sData := fmt.Sprintf("<p>Can't display dashboard: %s</p>", err.Error())
				fmt.Fprintf(w, pageMinimal(context, sData))
				return
			}

			fmt.Fprintf(w, string(buffer))

		}
	})
}

func metricsHandler(outputdir string, context string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		passed, failed, _, err := getMetrics(outputdir)
		if err != nil {
			fmt.Fprintf(w, "{}")
			return
		}

		result := "OK"
		if failed > 0 {
			result = "FAILED"
		}
		fmt.Fprintf(w, "{\"pass\":%d,\"fail\":%d,\"result\":\"%s\"}", passed, failed, result)
	})
}

func apiHandler(outputdir string, context string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(w, pageMinimal(context, buffer))
		return
	})
}

func healthHandler(outputdir string, context string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "{\"status\":\"ok\"}")
	})
}

func formatDuration(duration int64) string {
	seconds := duration / 1000000
	minutes := seconds / 60
	secondsRemaining := seconds % 60
	return fmt.Sprintf("%02dm:%02ds", minutes, secondsRemaining)
}
