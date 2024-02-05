# demo-temporal
Demo application built with Temporal in Golang

## Running

1. Clone [github.com/temporalio/docker-compose](https://github.com/temporalio/docker-compose) in a separate directory.
2. Change to the directory with the above repository and run `docker compose up -d`.
3. Change back to this repository's directory and run `docker compose up -d`.
4. Check that the API is up and running by curling its health check.
   It should return successfully with a status of "ok".
   ```
   $ curl http://localhost:8081/health
   {"status":"ok"}
   ```
