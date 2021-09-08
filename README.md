# Super Aggregator
An all in one log aggregation server written in Go.

## Building
`go build .`

## Developing
`go run .`

## Running in production
You're going to want to use a production ready Postgres database for running
this in production, as well as to put it behind an HTTPS load balancer. Here's
an example of what that would look like:
```
GIN_MODE=release ./super-aggregator \
  -db.driver=postgres \
  -pg.host=db1.local \
  -pg.user=super-aggregator \
  -pg.password=somesecuresecret

```

Setting GIN_MODE to "release" is important for performance when not developing.

## Roadmap
- Adapt to run in distributed mode (frontend, ingestor, queryor)
- Web GUI: dark mode
- Web GUI: fill in time picker while picked
