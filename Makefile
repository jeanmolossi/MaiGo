.PHONY: all test bench clean coverage-diff benchmark-diff lint install-go-lint

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

lint: install-go-lint
	go fmt ./...
	golangci-lint run ./... --fix

install-go-lint:
	@if ! command -v golangci-lint >/dev/null; then \
		read -p "Go's linter is not installed on your machine. Do you want to install it? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $$(go env GOPATH)/bin v2.3.0; \
			if ! command -v golangci-lint >/dev/null; then \
				echo "Go linter installation failed. Exiting..."; \
				exit 1; \
			fi; \
		fi; \
	fi
