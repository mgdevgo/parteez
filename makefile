build:
	go mod download && CGO_ENABLED=0 GOOS=linux go build -o ./bin/parteez ./cmd/parteez/main.go

run:
	go run ./cmd/server/main.go

lint:
	golangci-lint run

protoc:
	protoc \
		--go_out=./internal/api/ \
		--go_opt=paths=source_relative \
		--go-grpc_out=./internal/api/ \
		--go-grpc_opt=paths=source_relative \
		./api/account.proto

buf:
	easyp generate

buf-lint:
	easyp lint

swagger:
	swag init -g ./cmd/parteez/main.go -o ./api -outputTypes json
	swag init -g ./cmd/parteez/main.go -o ./docs -outputTypes go
	swag fmt

.PHONY:
.SILENT:
.DEFAULT_GOAL := run