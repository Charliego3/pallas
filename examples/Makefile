PROTOS=$(shell find . \( -path ./google \) -prune -type f -o -name "*.proto")
GOPATH=$(shell go env GOPATH)

.PHONY: gen
gen:
	protoc --go_out=./protos --go_opt=module=github.com/charliego3/pallas/examples/protos \
	--go-grpc_out=./protos --go-grpc_opt=module=github.com/charliego3/pallas/examples/protos \
	--go-pallas-http_out=./protos --go-pallas-http_opt=module=github.com/charliego3/pallas/examples/protos \
	$(PROTOS)
