obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver

distance_calculator:
	@go build -o bin/distance_calculator ./distance_calculator
	@./bin/distance_calculator

.PHONY: obu receiver distance_calculator