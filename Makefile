#!make
include .deploy/env/local.env
export $(shell sed 's/=.*//' ./.deploy/env/local.env)

V := @

RELEASE=$(shell git describe --always --tags)
BUILD_TIME?=$(shell date '+%F_%T%z')

OUT_DIR := ./bin
APP := grpc-captcha

default: all

.PHONY: all
all:service

deps:
	$(V)go mod tidy
	$(V) go mod vendor

.PHONY: service
service: APP_OUT := $(OUT_DIR)/$(APP)
service: APP_MAIN := ./cmd/$(APP)
service:
	@echo BUILDING $(APP_OUT)
	$(V)go build -ldflags "-s -w -X main.name=${APP} -X main.version=${RELEASE} -X main.compiledAt=${BUILD_TIME}" -o $(APP_OUT) $(APP_MAIN)
	@echo DONE

run:
	$(OUT_DIR)/${APP}

env: | service
	$(OUT_DIR)/${APP} --help | grep -o '$$[^ ]*'

test:
	@echo "Running autotests..."
	$(V)go clean -testcache
	$(V)go test -p 1 ./...

lint:
	@echo "Running golangci-lint..."
	$(V)golangci-lint run

HUB := "gtsh77workshop"
.PHONY: image
image:
	docker image build -t ${HUB}/${APP}:${RELEASE} -t ${HUB}/${APP}:latest -f Dockerfile.service .

.PHONY: protobuf
protobuf: PROTO_SRC:= ./pkg/proto/$(APP)
protobuf: PROTO_DEST:= ./pkg/proto/$(APP)
protobuf:
	protoc -I$(PROTO_SRC) --go_out=$(PROTO_DEST) --go_opt=paths=source_relative $(PROTO_SRC)/$(APP).proto 
	protoc -I$(PROTO_SRC) --go-grpc_out=$(PROTO_DEST) --go-grpc_opt paths=source_relative $(PROTO_SRC)/$(APP).proto
