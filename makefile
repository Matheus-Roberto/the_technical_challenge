generate:
	@protoc --go_out=. --go-grpc_opt=require_unimplemented_servers=false --go-grpc_out=. proto/*.proto

run:
	@echo "---- Running Server ----"
	@go run cmd\server\main.go

run_client:
	@echo "---- Running Client ----"
	@go run cmd\client\main.go

run_test:
	@echo "---- Running Client ----"
	@go run test\test.go