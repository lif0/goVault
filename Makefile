SERVER_APP_NAME=goVault-server
CLI_APP_NAME=goVault-cli

build-server:
	go build -o ${SERVER_APP_NAME} cmd/server/main.go

build-cli:
	go build -o ${CLI_APP_NAME} cmd/cli/main.go

run-server: build-server
	./${SERVER_APP_NAME}

run-server-with-config: build-server
	CONFIG_PATH=config.yml ./${SERVER_APP_NAME}

run-server-in-docker:
	docker build . && docker compose up

run-cli: build-cli
	./${CLI_APP_NAME} $(ARGS)

test:
	go test ./...

test_cover:
	go test ./internal/... -coverprofile=coverage.out && go tool cover -html=coverage.out