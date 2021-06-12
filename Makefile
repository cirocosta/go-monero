build:
	go build -v ./cmd/monero

test:
	go test -v ./pkg/...
