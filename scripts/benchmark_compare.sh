#!/usr/bin/env bash
set -euo pipefail

PR_OUT=${PR_OUT:-bench_pr.txt}
MAIN_OUT=${MAIN_OUT:-bench_main.txt}
COUNT=${COUNT:-10}

echo "Running benchmarks for current branch..." >&2
go test -run=^$ -bench=. -benchmem -count="$COUNT" ./... | tee "$PR_OUT" >/dev/null

TMP=""
cleanup() {
  if [[ -n "$TMP" ]]; then
    git worktree remove --force "$TMP/main" >/dev/null 2>&1 || true
    rm -rf "$TMP"
  fi
}
trap cleanup EXIT

if git remote get-url origin >/dev/null 2>&1; then
  git fetch origin main >/dev/null
  TMP=$(mktemp -d)
  git worktree add --detach "$TMP/main" origin/main >/dev/null
  pushd "$TMP/main" >/dev/null
  go test -run=^$ -bench=. -benchmem -count="$COUNT" ./... > "$OLDPWD/$MAIN_OUT"
  popd >/dev/null
else
  cp "$PR_OUT" "$MAIN_OUT"
fi

if ! command -v benchstat >/dev/null 2>&1; then
  echo "Installing benchstat..." >&2
  if ! GOSUMDB=off go install golang.org/x/perf/cmd/benchstat@latest; then
    echo "Failed to install benchstat" >&2
    exit 1
  fi
fi

echo "## Benchmark Results"
echo
diff_csv=$(benchstat -format csv "$MAIN_OUT" "$PR_OUT" 2>/dev/null)
python3 - "$diff_csv" <<'PYTHON'
import csv, sys, collections, io

data = collections.defaultdict(dict)
reader = csv.reader(io.StringIO(sys.argv[1]))
current = None
for row in reader:
    if not row or row[0].startswith(('goos:', 'goarch:', 'pkg:', 'cpu:')):
        continue
    if row[0] == '' and len(row) > 1 and row[1] in ('sec/op', 'B/op', 'allocs/op'):
        current = row[1]
        continue
    if row[0] == '' or row[0] == 'geomean' or len(row) < 6:
        continue
    bench = row[0]
    try:
        base = float(row[1])
        pr = float(row[3])
    except ValueError:
        continue
    delta = row[5]
    d = data[bench]
    if current == 'sec/op':
        d['time_base'] = base * 1e9
        d['time_pr'] = pr * 1e9
        d['time_delta'] = delta
    elif current == 'B/op':
        d['bytes_base'] = base
        d['bytes_pr'] = pr
        d['bytes_delta'] = delta
    elif current == 'allocs/op':
        d['allocs_base'] = base
        d['allocs_pr'] = pr
        d['allocs_delta'] = delta

headers = ['Benchmark', 'base ns/op', 'PR ns/op', 'Δ', 'base B/op', 'PR B/op', 'Δ', 'base allocs/op', 'PR allocs/op', 'Δ']
print('| ' + ' | '.join(headers) + ' |')
print('|' + '---|' * len(headers))
for bench in sorted(data):
    v = data[bench]
    print("| {} | {:.2f} | {:.2f} | {} | {:.2f} | {:.2f} | {} | {:.2f} | {:.2f} | {} |".format(
        bench,
        v.get('time_base', float('nan')),
        v.get('time_pr', float('nan')),
        v.get('time_delta', ''),
        v.get('bytes_base', float('nan')),
        v.get('bytes_pr', float('nan')),
        v.get('bytes_delta', ''),
        v.get('allocs_base', float('nan')),
        v.get('allocs_pr', float('nan')),
        v.get('allocs_delta', ''),
    ))
PYTHON
