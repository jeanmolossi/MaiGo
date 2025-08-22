#!/usr/bin/env bash
set -euo pipefail

MODULE=$(go list -m)
MIN_COVERAGE=${MIN_COVERAGE:-0}

go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out > current.func

# Fetch main and run tests in a temporary worktree
if git remote get-url origin >/dev/null 2>&1; then
    git fetch origin main >/dev/null
    TMPDIR=$(mktemp -d)
    trap 'rm -rf "$TMPDIR" current.func main.func' EXIT

    git worktree add --detach "$TMPDIR/main" origin/main >/dev/null
    pushd "$TMPDIR/main" >/dev/null
    go test -coverprofile=coverage.out ./...
    go tool cover -func=coverage.out > "$OLDPWD/main.func"
    popd >/dev/null
    git worktree remove --force "$TMPDIR/main" >/dev/null
else
    echo "origin remote not found" >&2
    exit 1
fi

python3 - "$MODULE" "$MIN_COVERAGE" <<'PY'
import os, sys, json
module = sys.argv[1]
min_cov = float(sys.argv[2])

def parse(path):
    pkg_cov = {}
    total = 0.0
    with open(path) as f:
        for line in f:
            line = line.rstrip()
            if line.startswith('total:'):
                total = float(line.split()[-1].strip('%'))
                continue
            parts = line.split('\t')
            if len(parts) < 3:
                continue
            file_part, _, cov_part = parts[0], parts[1], parts[2]
            if not file_part.startswith(module + '/'):
                continue
            rel = file_part[len(module)+1:]
            rel = rel.split(':')[0]
            pkg = os.path.dirname(rel)
            cov = float(cov_part.strip('%'))
            pkg_cov.setdefault(pkg, []).append(cov)
    for pkg in pkg_cov:
        covs = pkg_cov[pkg]
        pkg_cov[pkg] = sum(covs)/len(covs)
    return pkg_cov, total

curr_pkg, curr_total = parse('current.func')
main_pkg, main_total = parse('main.func')
all_pkgs = sorted(set(curr_pkg) | set(main_pkg))
print('| Package | main (%) | PR (%) | \u0394 (%) | Min (%) |')
print('| --- | --- | --- | --- | --- |')
for pkg in all_pkgs:
    mc = main_pkg.get(pkg, 0.0)
    cc = curr_pkg.get(pkg, 0.0)
    diff = cc - mc
    print(f'| {pkg} | {mc:.1f} | {cc:.1f} | {diff:+.1f} | {min_cov:.1f} |')
print(f'| **Total** | {main_total:.1f} | {curr_total:.1f} | {curr_total - main_total:+.1f} | {min_cov:.1f} |')
PY
