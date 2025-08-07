.PHONY: azure-service-bus redis ocpp send-message sqlc dev test test-coverage start stop

azure-service-bus:
	docker compose -f ./azure-service-bus/docker-compose.yaml up -d

redis:
	docker compose -f ./redis/docker-compose.yaml up -d

ocpp:
	go run -v ./cmd/ocpp/main.go

send-message:
	go run -v ./cmd/azure-service-bus/main.go

sqlc:
	cd ./service/ocpp/db && sqlc generate

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
	$(MAKE) send-message

stop:
	@echo "Stopping all services..."
	docker compose -f ./azure-service-bus/docker-compose.yaml down
	docker compose -f ./redis/docker-compose.yaml down
	@echo "All services stopped."