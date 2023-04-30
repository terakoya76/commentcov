lint:
	@if [ -z `which golangci-lint 2> /dev/null` ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.52.2; \
	fi
	gofmt -s -w .
	golangci-lint run --timeout 2m

test:
	go test -race -v ./...

proto-generate:
	protoc \
		--go_out=./proto --go_opt=paths=import --go_opt=module=github.com/commentcov/commentcov/proto --go_opt=Mcommentcov-proto/commentcov_plugin.proto=github.com/commentcov/commentcov/proto \
		--go-grpc_out=./proto --go-grpc_opt=paths=import --go-grpc_opt=module=github.com/commentcov/commentcov/proto --go-grpc_opt=Mcommentcov-proto/commentcov_plugin.proto=github.com/commentcov/commentcov/proto \
		commentcov-proto/commentcov_plugin.proto
