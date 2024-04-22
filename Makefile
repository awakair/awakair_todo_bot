SERVER_PACKAGE_PATH := cmd/TodoServiceServer/main.go
SERVER_BINARY_NAME := todo_service_server
SERVER_EXECUTABLE_PATH := /tmp/${SERVER_BINARY_NAME}
COVERAGE_FILE_PATH := /tmp/cover.out
COVERAGE_HTML_FILE_PATH = /tmp/cover.html
UNAME_S := $(shell uname -s)

ifeq (${UNAME_S}, Darwin)
	DEFAULT_OPEN := "open"
else
	DEFAULT_OPEN := "xdg-open"
endif

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## gen_proto_server: generate server's code from todo-service.proto specs
.PHONY: gen_proto_server
gen_proto_server:
	protoc \
	-I ./third_party/protovalidate/proto/protovalidate \
	-I ./api/todo-service \
	--go_out=api/todo-service \
	--go_opt=paths=source_relative \
  --go-grpc_out=api/todo-service \
	--go-grpc_opt=paths=source_relative \
  api/todo-service/todo-service.proto

## build_server: compile server's sources
.PHONY: build_server
build_server: gen_proto_server
	go build -o ${SERVER_EXECUTABLE_PATH} ${SERVER_PACKAGE_PATH}

## test_server: test server code
.PHONY: test_server
test_server:
	go test -cover -coverprofile ${COVERAGE_FILE_PATH} internal/TodoServiceServer/server.go internal/TodoServiceServer/server_test.go

## cover_server: measure coverage of server code and open it
.PHONE: cover_server
cover_server: test_server
	go tool cover -html ${COVERAGE_FILE_PATH} -o ${COVERAGE_HTML_FILE_PATH}
	${DEFAULT_OPEN} ${COVERAGE_HTML_FILE_PATH}

## run_server: compile and run server's sources
.PHONY: run_server
run_server: build_server
	${SERVER_EXECUTABLE_PATH}
