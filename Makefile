build:
	go build -v ./cmd/monero

debug:
	cd ./cmd/monero && dlv debug -- crawl --geo-ip-db=../../db.mmdb

test-debug:
	cd ./pkg/levin && dlv test

test:
	go test -v ./pkg/...


crawl-total-per-country:
	cat ./nodes.csv | awk -F ',' '{print $$3}' | sort | uniq -c | sort

crawl-total:
	cat ./nodes.csv | wc -l

crawl-reachable:
	cat ./nodes.csv | grep -v 'dial' | grep -v 'net' | grep -v 'reset' | wc -l

crawl-reachable-per-country:
	cat ./nodes.csv | grep -v 'dial' | grep -v 'net' | grep -v 'reset' | awk -F ',' '{print $$3}' | sort | uniq -c | sort
