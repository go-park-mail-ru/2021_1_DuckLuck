.PHONY: build
build:
	go build -o bin/api_server -v ./cmd/api_server
	go build -o bin/session_service -v ./cmd/session_service
	go build -o bin/cart_service -v ./cmd/cart_service
	go build -o bin/auth_service -v ./cmd/auth_service

.PHONY: local
local:
	echo TAG=local > .env
	docker build -t duckluckmarket/api-server:local --target api-server .
	docker build -t duckluckmarket/api-db:local --target api-db .
	docker build -t duckluckmarket/auth-service:local --target auth-service .
	docker build -t duckluckmarket/session-service:local --target session-service .
	docker build -t duckluckmarket/cart-service:local --target cart-service .
	docker-compose down
	docker-compose up -d

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -coverprofile=coverage1.out -coverpkg=./... -cover ./...
	cat coverage1.out | grep -v mock > cover.out
	go tool cover -func cover.out && go tool cover -html cover.out

.PHONY: init_db
init_db:
	sudo -u postgres psql -f configs/init.sql

.DEFAULT_GOAL := build
