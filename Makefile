.PHONY: run
run_main_server:
	go run cmd/server/main.go

.PHONY: run
run_file_server:
	go run cmd/fileserver/main.go
