install:
	go install -v ./cmd/monero

build:
	go build -v ./cmd/monero

test:
	go test -v ./pkg/...
