.PHONY: generate integration-test unit-test test install
.DEFAULT_GOAL := test

generate:
	go generate ./...
	
integration-test: generate
	go test -tags=integration ./...

unit-test: generate
	go test ./...

test: unit-test integration-test

install: generate
	go install ./cmd/cpx