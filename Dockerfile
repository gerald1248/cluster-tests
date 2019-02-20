FROM golang:1.11.1 as builder
WORKDIR /go/src/github.com/gerald1248/cluster-tests/
COPY * ./
ENV CGO_ENABLED 0
ENV GOOS linux
ENV GO111MODULE on
RUN \
  go mod download && \
  go get && \
  go vet && \
  go test -v -cover && \
  go build -o cluster-tests .

FROM ubuntu:18.10
WORKDIR /app/
EXPOSE 8080
ENV CLUSTER_TESTS_BLACKLIST default,kube,flux
RUN apt-get update && \
  DEBIAN_FRONTEND=noninteractive apt-get -qq install curl
COPY --from=builder /go/src/github.com/gerald1248/cluster-tests/cluster-tests /usr/bin/
USER 1000
CMD ["cluster-tests", "-s=true"]
