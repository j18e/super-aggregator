# Super Aggregator
An all in one log aggregation server written in Go.

## Building
`go build .`

## Developing
You can run the development environment either in Docker or outside of Docker.
Outside Docker:
```
docker-compose up postgres -d
go run . -pg.password=supersecret
```

In Docker:
```
docker-compose up --build
```
The docker-compose environment mounts the local `views/` folder so that updates
to template files will be applied by the Gin HTTP server immediately.

## Running in production
You're going to want to use a production ready Postgres database for running
this in production, as well as to put it behind an HTTPS load balancer. Here's
an example of what that would look like:
```
GIN_MODE=release ./super-aggregator \
  -pg.host=db1.local \
  -pg.user=super-aggregator \
  -pg.password=somesecuresecret
```
Note that the Postgres database will need to be available, as will the `views/`
folder which contains HTML templates.

Setting GIN_MODE to "release" is important for performance when not developing.

## Roadmap
- Adapt to run in distributed mode (frontend, ingestor, queryor)
- Implement cache to reduce database pressure
- Web GUI: fill in time picker while picked
- Web GUI: improve look and feel
- Web GUI: dark mode
