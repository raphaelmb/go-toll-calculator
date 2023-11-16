obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver

distance_calculator:
	@go build -o bin/distance_calculator ./distance_calculator
	@./bin/distance_calculator

aggregator:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator

gateway:
	@go build -o bin/gateway gateway/main.go
	@./bin/gateway

proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

.PHONY: obu receiver distance_calculator aggregator proto gateway