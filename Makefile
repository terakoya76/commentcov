lint:
	@if [ -z `which golangci-lint 2> /dev/null` ]; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.46.2; \
	fi
	gofmt -s -w .
	golangci-lint run --timeout 2m

test:
	go test -race -v ./...

proto-generate:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/commentcov.proto
