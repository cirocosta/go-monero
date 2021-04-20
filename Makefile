build:
	go build -v ./cmd/monero

test-debug:
	cd ./pkg/levin && dlv test

test:
	go test -v ./pkg/...
