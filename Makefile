.PHONY: test lint run

shelf:
	go build -o shelf

test:
	source env.sh && go test -v

lint:
	golangci-lint run ./...

run: shelf
	source env.sh && ./shelf
