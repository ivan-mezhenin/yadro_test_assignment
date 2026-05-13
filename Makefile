SERVER_CONFIG = server/config/config.yaml
RESOLV_FILE  = /tmp/test_resolv.conf
BACKUP_FILE  = /tmp/test_resolv.conf.bak

.PHONY: build test run-server run-client clean docker-build docker-up

build:
	go build ./...

test:
	go test ./server/tests/... -v

run-server:
	go run ./server/cmd/server --config $(SERVER_CONFIG) --resolv $(RESOLV_FILE) --backup $(BACKUP_FILE)

run-client:
	go run ./client/cmd/client $(CMD)
	
clean:
	rm -f $(RESOLV_FILE) $(BACKUP_FILE)
