build:
	go build -v ./cmd/monero

debug:
	cd ./cmd/test && dlv debug

test-debug:
	cd ./pkg/levin && dlv test

test:
	go test -v ./pkg/...
