install:
	go install ./...

format:
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .

test:
	go test ./...
