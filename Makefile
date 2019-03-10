all:
	go build -o procstat
stdout:
	go build -o procstat -tags stdout
fmt:
	go fmt
clean:
	go clean
	rm ./procstat
