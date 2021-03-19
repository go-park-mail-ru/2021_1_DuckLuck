.PHONY: build
build:
	go build -o bin/server -v ./cmd/server
	go build -o bin/fileserver -v ./cmd/fileserver

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
