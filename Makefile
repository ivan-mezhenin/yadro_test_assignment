SERVER_CONFIG = server/config/config.yaml
RESOLV_FILE  = /tmp/test_resolv.conf
BACKUP_FILE  = /tmp/test_resolv.conf.bak

PROTO_FILE = proto/dns.proto

.PHONY: build test run-server run-client proto clean docker-build docker-up

build:
	go build ./...

test:
	go test ./server/tests/... -v

run-server:
	go run ./server/cmd/server --config $(SERVER_CONFIG) --resolv $(RESOLV_FILE) --backup $(BACKUP_FILE)

run-client:
	go run ./client/cmd/client $(CMD)
	
proto:
	protoc --go_out=. --go_opt=module=dns-manager --go-grpc_out=. --go-grpc_opt=module=dns-manager $(PROTO_FILE)

clean:
	rm -f $(RESOLV_FILE) $(BACKUP_FILE)
