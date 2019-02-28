FROM golang:1.11.1 as builder
WORKDIR /go/src/github.com/gerald1248/cluster-tests/
ADD . ./
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
ENV KUBE_VERSION=v1.12.0
EXPOSE 8080
RUN groupadd app && useradd -g app app && \
  mkdir -p /app/output /app/cluster-tests.d && \
  chmod -R 777 /app && \
  chown app:app /app/cluster-tests.d /app/output
RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get -qq install -y jq apt-transport-https wget curl gnupg
RUN wget -O /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/${KUBE_VERSION}/bin/linux/amd64/kubectl && \
  chmod +x /usr/local/bin/kubectl
RUN mkdir -p /app/output
COPY --from=builder /go/src/github.com/gerald1248/cluster-tests/cluster-tests /usr/local/bin/cluster-tests
COPY --from=builder /go/src/github.com/gerald1248/cluster-tests/cluster-tests.d/* /app/cluster-tests.d/
USER app
CMD ["cluster-tests", "-d=/app/cluster-tests.d", "-o=/app/output"]
