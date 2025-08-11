
.PHONY: all test bench clean unit-test tests

all: test bench

test:
	go test -coverprofile=.coverage.out -covermode=atomic ./...
	go tool cover -html=.coverage.out -o coverage.html

bench:
	go test -run=^$$ -count=1 -bench=. -benchmem ./... > benchmark.txt

clean:
	@rm -f .coverage.out coverage.html
	@rm -f benchmark.*

unit-test tests:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
