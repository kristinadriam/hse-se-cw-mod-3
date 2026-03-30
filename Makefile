.PHONY: build test pytest proto fmt

PYTHON := $(if $(wildcard .venv/bin/python),.venv/bin/python,python3)

build:
	go build -o bin/chat ./cmd/chat

test:
	go test ./...

pytest: build
	$(PYTHON) -m pytest

fmt:
	gofmt -w .

proto:
	protoc --go_out=. --go_opt=module=hse-se-cw-mod-3 \
		--go-grpc_out=. --go-grpc_opt=module=hse-se-cw-mod-3 \
		proto/chat/v1/chat.proto
