cluster-tests
=============

`cluster-tests` is a test runner intended for in-cluster use.

<img src="images/cluster-tests-screenshot.png" alt="dashboard screenshot"/>

The available endpoints are:

| Endpoint        | Description      |
| --------------- | ---------------- |
| /               | Show dashboard   |
| /health         | Health endpoint  |
| /api/v1/metrics | Metrics endpoint |

Build
-----
The build steps are the following:
```
$ go mod download
$ go get
$ go vet
$ go test -v
$ go build -o cluster-tests .
```

`make build` will run these steps in a two-stage docker build process.

