# Super Aggregator
An all in one log aggregation server written in Go.

The aggregator comes with a web GUI and a storage backend which enables users
both to ship logs to it through JSON payloads over HTTP, and to view and sort
logs according to application, time range and other metadata.

## Roadmap
- Adapt to run in distributed mode (frontend, ingestor, queryor)
- Implement cache to reduce database pressure
- Web GUI: fill in time picker while picked
- Web GUI: improve look and feel
- Web GUI: dark mode

## Posting logs
You could post some example logs to the aggregator like so:
```
curl localhost:9000/api/entry -d '[{
    "timestamp": "2020-09-09T09:09:09+02:00",
    "application": "myapp",
    "host": "my-server",
    "environment": "dev",
    "log_line": "something just happened"
}]'
```

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

### Running in containerized environments
Container images for super-aggregator are hosted on [docker hub](https://hub.docker.com/repository/docker/j18e/super-aggregator).
```
docker run \
    -e GIN_MODE=release \
    -p 9000:9000 \
    j18e/super-aggregator:latest \
    -pg.host=db1.local \
    -pg.user=super-aggregator \
    -pg.password=somesecuresecret
```

## Continuous integration
Commits to the main branch trigger builds using Github actions, which if
successful, push built images to Docker Hub.
