FROM golang:1.20 AS builder

ENV GOOS=linux
ENV GOARCH=386

WORKDIR /work

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build -o ./super-aggregator .

FROM alpine:3.14

COPY --from=builder /work/super-aggregator .
COPY ./views ./views

ENTRYPOINT ["./super-aggregator"]
