all:
	go build -o procstat-go
stdout:
	go build -o procstat-go -tags stdout
fmt:
	go fmt
clean:
	go clean

lint:
	golangci-lint run
lint-stdout:
	golangci-lint run --build-tags stdout

xxx1:	fmt lint
xxx2:	fmt lint-stdout
