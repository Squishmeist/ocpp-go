.PHONY: ocpp, azure-service-bus, azure

ocpp:
	go run -v ./cmd/ocpp/main.go

azure-service-bus:
	docker compose -f ./service/azure-service-bus/docker-compose.yaml up -d

azure:
	go run -v ./cmd/azure-service-bus/main.go