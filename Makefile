.PHONY: all test bench clean coverage-diff benchmark-diff lint

all: test bench

test:
	go test -coverprofile=.coverage.out -covermode=atomic ./...
	go tool cover -html=.coverage.out -o coverage.html

bench:
	go test -run=^$$ -count=1 -bench=. -benchmem ./... > benchmark.txt

clean:
	@rm -f .coverage.out coverage.html coverage.out coverage_report.md
	@rm -f benchmark.*

coverage-diff:
	bash -eo pipefail scripts/coverage_report.sh | tee coverage_report.md

benchmark-diff:
	bash -eo pipefail scripts/benchmark_compare.sh | tee benchmark_report.md

lint:
	golangci-lint run
