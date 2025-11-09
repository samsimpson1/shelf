.PHONY: test

shelf:
	go build -o shelf

test:
	source env.sh && go test -v

make run: shelf
	source env.sh && ./shelf