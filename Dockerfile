FROM golang:1.17-alpine3.14 as builder

ADD . /go/src/coinconv

WORKDIR /go/src/coinconv

RUN apk update \
    apk add --no-cache gcc musl-dev
    
RUN go build -o /tmp/coinconv ./cmd/coinconv

FROM alpine:3.14

RUN apk update \
    apk add --no-cache ca-certificates

COPY --from=builder /tmp/coinconv /usr/bin/coinconv

RUN chmod +x /usr/bin/coinconv

ENTRYPOINT [ "coinconv" ]