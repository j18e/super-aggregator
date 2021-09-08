FROM golang:1.17 AS builder

ENV GOARCH amd64
ENV GOOS linux

WORKDIR /work

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build -o super-aggregator

FROM alpine:3.14

COPY --from=builder /work/super-aggregator .

ENTRYPOINT ["./super-aggregator"]
