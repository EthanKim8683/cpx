-include .env
export

.PHONY: generate install

generate:
	go generate ./...

install: generate
	go install ./cmd/cpx