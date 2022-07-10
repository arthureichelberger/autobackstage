.DEFAULT_GOAL := help
PLATFORM=linux-amd64
GRC=$(shell which grc)

test:
	ENVIRONMENT=test $(GRC) go test -v -p=1 -count=1 -race ./... -timeout 2m

build:
	@CGO_ENABLED=0 GOOS=linux go build -o deployment/bin/autobackstage cmd/main.go

build-mac:
	go build -o deployment/bin/autobackstage cmd/main.go
