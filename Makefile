#!make
SHELL=/bin/sh
include .deploy/env/local.env
export $(shell sed 's/=.*//' ./.deploy/env/local.env)

RELEASE=$(shell git describe --always --tags)
BUILD_TIME?=$(shell date '+%F_%T%z')

OUT_DIR = ./bin
APP = grpc-captcha
APP_MAIN = ./cmd/$(APP)

.PHONY: default
default: all

.PHONY: all
all: service

.PHONY: deps
deps:
	@go mod tidy
	@go mod vendor

.PHONY: service
service: linux64

.PHONY: run
run:
	$(OUT_DIR)/$(APP)_linux_amd64

.PHONY: env
env: service
	$(OUT_DIR)/$(APP) --help | grep -o '$$[^ ]*'

.PHONY: test
test:
	@echo "Running autotests..."
	@go clean -testcache
	@go test -p 1 ./...

.PHONY: lint
lint:
	@echo "Running golangci-lint..."
	@golangci-lint run

.PHONY: image
HUB = gtsh77workshop
image:
	@docker image build -t $(HUB)/$(APP):$(RELEASE) -t $(HUB)/$(APP):latest -f Dockerfile.service .

.PHONY: protobuf
PROTO_SRC = ./pkg/proto/$(APP)
PROTO_DEST = ./pkg/proto/$(APP)
protobuf: 
	@protoc -I$(PROTO_SRC) --go_out=$(PROTO_DEST) --go_opt=paths=source_relative $(PROTO_SRC)/$(APP).proto 
	@protoc -I$(PROTO_SRC) --go-grpc_out=$(PROTO_DEST) --go-grpc_opt paths=source_relative $(PROTO_SRC)/$(APP).proto

.PHONY: release
release: clean linux64 linux32 windows64 windows32 freebsd64 freebsd32 darwin64
	@echo OK

.PHONY: linux64
linux64:
	@echo BUILDING $(OUT_DIR)/$(APP)_linux_amd64
	@GOOS=linux GOARCH=amd64 go build -ldflags "-s -w -X main.name=$(APP) -X main.version=$(RELEASE) -X main.compiledAt=$(BUILD_TIME)" -o $(OUT_DIR)/$(APP)_linux_amd64 $(APP_MAIN)
.PHONY: linux32
linux32:
	@echo BUILDING $(OUT_DIR)/$(APP)_linux_386
	@GOOS=linux GOARCH=386 go build -ldflags "-s -w -X main.name=$(APP) -X main.version=$(RELEASE) -X main.compiledAt=$(BUILD_TIME)" -o $(OUT_DIR)/$(APP)_linux_386 $(APP_MAIN)
.PHONY: windows64	
windows64:
	@echo BUILDING $(OUT_DIR)/$(APP)_windows_amd64
	@GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -X main.name=$(APP) -X main.version=$(RELEASE) -X main.compiledAt=$(BUILD_TIME)" -o $(OUT_DIR)/$(APP)_windows_amd64 $(APP_MAIN)
.PHONY: windows32
windows32:
	@echo BUILDING $(OUT_DIR)/$(APP)_windows_386
	@GOOS=windows GOARCH=386 go build -ldflags "-s -w -X main.name=$(APP) -X main.version=$(RELEASE) -X main.compiledAt=$(BUILD_TIME)" -o $(OUT_DIR)/$(APP)_windows_386 $(APP_MAIN)
.PHONY: freebsd64
freebsd64:
	@echo BUILDING $(OUT_DIR)/$(APP)_freebsd_amd64
	@GOOS=freebsd GOARCH=amd64 go build -ldflags "-s -w -X main.name=$(APP) -X main.version=$(RELEASE) -X main.compiledAt=$(BUILD_TIME)" -o $(OUT_DIR)/$(APP)_freebsd_amd64 $(APP_MAIN)
.PHONY: freebsd32
freebsd32:
	@echo BUILDING $(OUT_DIR)/$(APP)_freebsd_386
	@GOOS=freebsd GOARCH=386 go build -ldflags "-s -w -X main.name=$(APP) -X main.version=$(RELEASE) -X main.compiledAt=$(BUILD_TIME)" -o $(OUT_DIR)/$(APP)_freebsd_386 $(APP_MAIN)
.PHONY: darwin64
darwin64:
	@echo BUILDING $(OUT_DIR)/$(APP)_darwin_amd64
	@GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w -X main.name=$(APP) -X main.version=$(RELEASE) -X main.compiledAt=$(BUILD_TIME)" -o $(OUT_DIR)/$(APP)_darwin_amd64 $(APP_MAIN)

.PHONY: clean
clean: 
	@rm -rf $(OUT_DIR)

