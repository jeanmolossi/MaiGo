.PHONY: all test clean

all: test

test:
	go test -coverprofile=.coverage.out -covermode=atomic ./...
	go tool cover -html=.coverage.out -o coverage.html

clean:
	@rm -f .coverage.out coverage.html
