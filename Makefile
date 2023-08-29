grpc-server: 
	go run server/cmd/main.go

grpc-client: 
	go run client/cmd/main.go

proto:
	protoc --go_out=. --go-grpc_out=. ./proto/*.proto

# docker
build:
	docker-compose build

run:
	docker-compose up -d

run-server:
	docker-compose run app

run-client:
	docker-compose run client
