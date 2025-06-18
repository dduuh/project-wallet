build:
	docker exec -it kafka-wallet \
	bash -c "/opt/bitnami/kafka/bin/kafka-topics.sh --create --topic users --bootstrap-server localhost:9094 --partitions 1 --replication-factor 1"

run-consumer:
	go run cmd/users-consumer/main.go

run-producer:
	go run cmd/users-producer/main.go

run-server:
	go run cmd/server/main.go

target-lint:
	golangci-lint run -c .golangci.yml