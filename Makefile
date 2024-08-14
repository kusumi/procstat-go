bin1:
	go build
bin2:
	go build -tags stdout
clean:
	go clean
fmt:
	go fmt
lint1:
	golangci-lint run
lint2:
	golangci-lint run --build-tags stdout

xxx1:	fmt lint1
xxx2:	fmt lint2
