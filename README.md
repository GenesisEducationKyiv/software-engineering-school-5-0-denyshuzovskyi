# nimbus-notify
Service that allows to subscribe to regular emails with weather updates

### Demo
Deployed version: http://db35m6zjaamdj.cloudfront.net/

**Email sending currently works only for pre-approved recipients due to the use of a sandbox domain**

### Task
https://github.com/mykhailo-hrynko/se-school-5


### Useful commands
migrate -source file://./migrations -database "postgres://user:password@localhost:5432/nimbus-notify?sslmode=disable" up

mockery init nimbus-notify/internal/service
mockery

go test ./...

docker build -f api-server.Dockerfile -t nimbus-notify .
docker-compose up

golangci-lint run
golangci-lint fmt

