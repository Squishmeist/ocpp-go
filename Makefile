.PHONY: azure-service-bus ocpp send-message dev

azure-service-bus:
	docker compose -f ./service/azure-service-bus/docker-compose.yaml up -d

ocpp:
	go run -v ./cmd/ocpp/main.go

send-message:
	go run -v ./cmd/azure-service-bus/main.go $(ARGS)

dev:
	@echo "1. Run 'make azure-service-bus' to start the emulator"
	@echo "2. Run 'make ocpp' to start the OCPP listener"
	@echo "3. Run 'make send-message ARGS=heartbeatrequest' to send a message"
