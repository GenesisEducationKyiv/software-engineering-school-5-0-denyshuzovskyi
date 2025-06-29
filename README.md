# nimbus-notify
Service that allows to subscribe to regular emails with weather updates

### Demo
Deployed version: http://db35m6zjaamdj.cloudfront.net/

**Email sending currently works only for pre-approved recipients due to the use of a sandbox domain**

### Task
https://github.com/mykhailo-hrynko/se-school-5


### Useful commands

#### Migrate
```shell
  migrate -source file://./migrations -database "postgres://user:password@localhost:5432/nimbus-notify?sslmode=disable" up
```

#### Mockery
```shell
  mockery init github.com/GenesisEducationKyiv/software-engineering-school-5-0-denyshuzovskyi/internal/service/weather
  mockery
```

#### Tests
```shell
  go test ./...
```

#### Docker
```shell
    docker build -f api-server.Dockerfile -t nimbus-notify .
    docker-compose up
```

#### Golangci-lint
```shell
    golangci-lint run --config ./.golangci.yaml --verbose  
    golangci-lint run --config ./.golangci.yaml --fix --verbose
    golangci-lint run --config ./.golangci.yaml --enable-only staticcheck --fix --verbose
    
    golangci-lint fmt --config ./.golangci.yaml --verbose
```

