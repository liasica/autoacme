.PHONY: build

export GO111MODULE=on
export CGO_ENABLED=0

TAG_NAME := $(shell git tag -l --contains HEAD)
SHA := $(shell git rev-parse HEAD)
VERSION := $(if $(TAG_NAME),$(TAG_NAME),$(SHA))

BIN_OUTPUT := build/aliacme
MAIN_DIRECTORY := ./cmd/aliacme

build:
	@echo Version: $(VERSION)
	go build -trimpath -ldflags '-X "main.version=${VERSION}"' -o ${BIN_OUTPUT} ${MAIN_DIRECTORY}

