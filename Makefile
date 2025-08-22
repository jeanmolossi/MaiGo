.PHONY: all test bench clean coverage-diff

all: test bench

test:
	go test -coverprofile=.coverage.out -covermode=atomic ./...
	go tool cover -html=.coverage.out -o coverage.html

bench:
	go test -run=^$$ -count=1 -bench=. -benchmem ./... > benchmark.txt

clean:
	@rm -f .coverage.out coverage.html
	@rm -f benchmark.*

coverage-diff:
	./scripts/coverage_report.sh > coverage_report.md
