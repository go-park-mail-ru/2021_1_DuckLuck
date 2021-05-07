.PHONY: build
build:
	go build -o bin/api_server -v ./cmd/api_server
	go build -o bin/session_service -v ./cmd/session_service
	go build -o bin/cart_service -v ./cmd/cart_service
	go build -o bin/auth_service -v ./cmd/auth_service

.PHONY: api_server
api_server:
	go build -o bin/api_server -v ./cmd/api_server

.PHONY: session_service
session_service:
	go build -o bin/session_service -v ./cmd/session_service

.PHONY: session_service
cart_service:
	go build -o bin/cart_service -v ./cmd/cart_service

.PHONY: auth_service
auth_service:
	go build -o bin/cart_service -v ./cmd/auth_service

.PHONY: start_local
start_local:
	make stop_local
	echo TAG=local > .env
	docker build -t duckluckmarket/api-server:local --target api-server .
	docker build -t duckluckmarket/api-db:local --target api-db .
	docker build -t duckluckmarket/auth-service:local --target auth-service .
	docker build -t duckluckmarket/session-service:local --target session-service .
	docker build -t duckluckmarket/cart-service:local --target cart-service .
	docker-compose up -d

.PHONY: stop_local
stop_local:
	docker-compose down

.PHONY: remove_containers
remove_containers:
	-docker stop $$(docker ps -aq)
	-docker rm $$(docker ps -aq)

.PHONY: armageddon
armageddon:
	-make remove_containers
	-docker network prune -f
	-docker volume rm $$(docker volume ls --filter dangling=true -q)
	-docker rmi $$(docker images -a -q) -f


.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -coverprofile=coverage1.out -coverpkg=./... -cover ./...
	cat coverage1.out | grep -v mock | grep -v proto | grep -v cmd > cover.out
	go tool cover -func cover.out && go tool cover -html cover.out

.PHONY: init_db
init_db:
	sudo -u postgres psql -f scripts/postgresql/init_api_db.sql
	sudo -u postgres psql -f scripts/postgresql/init_auth_db.sql

.DEFAULT_GOAL := build
