FROM golang:1.17 as builder
WORKDIR /build
COPY go.mod go.sum *.go /build/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -trimpath -ldflags "-s -w" -o /build/zabrbl

FROM scratch
COPY --from=builder /build/zabrbl /zabrbl
ENTRYPOINT [ "/zabrbl" ]

