#!/usr/bin/env bash
set -euo pipefail

PR_OUT=${PR_OUT:-bench_pr.txt}
MAIN_OUT=${MAIN_OUT:-bench_main.txt}

echo "Running benchmarks for current branch..." >&2
go test -bench=. -benchmem ./... | tee "$PR_OUT" >/dev/null

TMP=""
trap '[[ -n "$TMP" ]] && git worktree remove --force "$TMP/main" >/dev/null 2>&1 && rm -rf "$TMP"' EXIT

if git remote get-url origin >/dev/null 2>&1; then
  git fetch origin main >/dev/null
  TMP=$(mktemp -d)
  git worktree add --detach "$TMP/main" origin/main >/dev/null
  pushd "$TMP/main" >/dev/null
  go test -bench=. -benchmem ./... > "$OLDPWD/$MAIN_OUT"
  popd >/dev/null
else
  cp "$PR_OUT" "$MAIN_OUT"
fi

# install benchstat for comparison
if ! command -v benchstat >/dev/null 2>&1; then
  go install golang.org/x/perf/cmd/benchstat@latest
fi

echo "## Benchmark Results"
benchstat -markdown "$MAIN_OUT" "$PR_OUT"
