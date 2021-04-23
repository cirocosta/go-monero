build:
	go build -v ./cmd/monero

debug:
	cd ./cmd/monero && dlv debug -- crawl

test-debug:
	cd ./pkg/levin && dlv test

test:
	go test -v ./pkg/...
