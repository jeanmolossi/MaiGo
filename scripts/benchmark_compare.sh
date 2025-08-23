#!/usr/bin/env bash
set -euo pipefail

PR_OUT=${PR_OUT:-bench_pr.txt}
MAIN_OUT=${MAIN_OUT:-bench_main.txt}

echo "Running benchmarks for current branch..." >&2
go test -run=^$ -bench=. -benchmem ./... | tee "$PR_OUT" >/dev/null

TMP=""
trap '[[ -n "$TMP" ]] && git worktree remove --force "$TMP/main" >/dev/null 2>&1 && rm -rf "$TMP"' EXIT

if git remote get-url origin >/dev/null 2>&1; then
  git fetch origin main >/dev/null
  TMP=$(mktemp -d)
  git worktree add --detach "$TMP/main" origin/main >/dev/null
  pushd "$TMP/main" >/dev/null
  go test -run=^$ -bench=. -benchmem ./... > "$OLDPWD/$MAIN_OUT"
  popd >/dev/null
else
  cp "$PR_OUT" "$MAIN_OUT"
fi

# install benchstat for comparison
if ! command -v benchstat >/dev/null 2>&1; then
  echo "Installing benchstat..." >&2
  if ! go install golang.org/x/perf/cmd/benchstat@latest; then
    echo "Failed to install benchstat" >&2
    exit 1
  fi
fi

echo "## Benchmark Results"
echo
echo '```'
benchstat -format text "$MAIN_OUT" "$PR_OUT"
echo '```'
