.PHONY: azure-service-bus redis ocpp message sqlc proto dev test test-coverage start stop

azure-service-bus:
	docker compose -f ./azure-service-bus/docker-compose.yaml up -d

redis:
	docker compose -f ./redis/docker-compose.yaml up -d

ocpp:
	go run -v ./cmd/ocpp/main.go

message:
	go run -v ./cmd/message/main.go

sqlc:
	cd ./service/ocpp/db && sqlc generate

proto:
	protoc --go_out=./pkg --go_opt=paths=source_relative \
		--go-grpc_out=./pkg --go-grpc_opt=paths=source_relative \
		api/proto/ocpp/v1/ocpp.proto

test:
	gotestsum
	
test-coverage:
	gotestsum -- -coverprofile=cover.out ./...

dev:
	@echo "1. Run 'make start' to start the Azure Service Bus and Redis server"
	@echo "2. Run 'make ocpp' to start the OCPP machine"

start:
	$(MAKE) azure-service-bus
	$(MAKE) redis
	$(MAKE) message

stop:
	@echo "Stopping all services..."
	docker compose -f ./azure-service-bus/docker-compose.yaml down
	docker compose -f ./redis/docker-compose.yaml down
	@echo "All services stopped."