.PHONY: run
run_main_server:
	go run cmd/server/main.go

.PHONY: run
run_file_server:
	go run cmd/fileserver/main.go

.PHONY: test
test:
	go test -coverprofile=coverage1.out -coverpkg=./... -cover ./...
	cat coverage1.out | grep -v mock > cover.out
	go tool cover -func cover.out && go tool cover -html cover.out
