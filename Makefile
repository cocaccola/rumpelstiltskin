SHELL=/bin/bash
.ONESHELL:

linux: main.go
	@./scripts/build.sh amd64

linux-arm64: main.go
	@./scripts/build.sh arm64

local: main.go
	@go build -o rumpelstiltskin main.go

image: main.go Dockerfile
	@./scripts/build-docker.sh

clean:
	@rm rumpelstiltskin

.PHONY: linux linux-arm64 local image clean
