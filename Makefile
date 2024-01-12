VERSION := 1.0.0

install:
	go install -ldflags="-X main.Version=v${VERSION}" ./...

format:
	go install golang.org/x/tools/cmd/goimports@latest
	goimports -w .

test:
	go test ./...
