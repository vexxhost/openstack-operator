FROM golang:1.13 AS builder
RUN go get -v -u github.com/Dev25/mcrouter_exporter && \
    cd $GOPATH/src/github.com/Dev25/mcrouter_exporter && \
    CGO_ENABLED=0 go build -a -installsuffix cgo -o /app

FROM scratch
COPY --from=builder /app /mcrouter_exporter
EXPOSE 9442
ENTRYPOINT ["/mcrouter_exporter"]