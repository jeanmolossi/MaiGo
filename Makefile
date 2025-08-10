.PHONY: all test clean

all: test bench

test:
	go test -coverprofile=.coverage.out -covermode=atomic ./...
	go tool cover -html=.coverage.out -o coverage.html

bench:
	go test -bench=. -benchmem ./... > benchmark.txt

clean:
	@rm -f .coverage.out coverage.html
	@rm -f benchmark.txt
