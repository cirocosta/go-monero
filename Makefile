install:
	go install -v ./cmd/monero

build:
	go build -v ./cmd/monero

test:
	go test -v ./pkg/...

lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint run --config=.golangci.yaml

.images.lock.yaml: .images.yaml
	kbld -f $< --lock-output $@
.PHONY: .images.lock.yaml
