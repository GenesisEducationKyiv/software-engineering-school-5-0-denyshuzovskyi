## gRPC
protoc --go_out=. --go-grpc_out=. email/email.proto

buf config init

buf build

buf generate

//add deps

buf dep update

buf dep prune

buf lint