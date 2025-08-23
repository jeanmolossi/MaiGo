#!/usr/bin/env bash
set -euo pipefail

MODULE=$(go list -m 2>/dev/null || true)
if [[ -z "$MODULE" ]]; then
  echo "Failed to determine module via 'go list -m'; ensure you are in a module or workspace" >&2
  exit 1
fi
MIN_COVERAGE=${MIN_COVERAGE:-0}

cov_branch() {
  local profile=$1
  go test -count=1 -covermode=atomic -coverprofile="$profile" ./... >&2
  go tool cover -func="$profile" | awk -v module="$MODULE" '
    $1 == "total:" {
      gsub("%", "", $NF)
      total = $NF
      next
    }
    index($1, module) == 1 {
      split($1, path, ":")
      rel = substr(path[1], length(module)+2)
      pkg = rel
      sub(/\/[^\/]+$/, "", pkg)
      if (pkg == "") pkg = "."
      cov = $NF
      gsub("%", "", cov)
      sum[pkg] += cov
      count[pkg]++
    }
    END {
      for (p in sum) printf "%s\t%.1f\n", p, sum[p]/count[p]
      printf "__total__\t%s\n", total
    }' > "$profile.func"
}

cov_branch coverage.out

TMP=""
trap 'rm -rf "${TMP:-}" main.out.func coverage.out.func' EXIT

BASE_BRANCH="${BASE_REF:-${GITHUB_BASE_REF:-}}"
if [[ -z "$BASE_BRANCH" ]]; then
  upstream=$(git rev-parse --abbrev-ref --symbolic-full-name @{u} 2>/dev/null || true)
  [[ -n "$upstream" ]] && BASE_BRANCH=${upstream#*/}
fi
if [[ -z "$BASE_BRANCH" ]]; then
  BASE_BRANCH=$(git symbolic-ref --short HEAD 2>/dev/null || true)
fi

if git remote get-url origin >/dev/null 2>&1 && [[ -n "$BASE_BRANCH" ]]; then
  git fetch origin "$BASE_BRANCH" >/dev/null
  TMP=$(mktemp -d)
  git worktree add --detach "$TMP/main" "origin/$BASE_BRANCH" >/dev/null
  pushd "$TMP/main" >/dev/null
  cov_branch main.out
  popd >/dev/null
  mv "$TMP/main/main.out.func" main.out.func
  git worktree remove --force "$TMP/main" >/dev/null
else
  cp coverage.out.func main.out.func
fi

echo "## Coverage Report"
awk -F'\t' -v min="$MIN_COVERAGE" '
NR==FNR { main[$1]=$2; next }
{ curr[$1]=$2 }
END {
  total_main = main["__total__"] + 0
  total_curr = curr["__total__"] + 0
  delete main["__total__"]; delete curr["__total__"]
  print "| Package | main (%) | PR (%) | Δ (%) | Min (%) |"
  print "| --- | --- | --- | --- | --- |"
  for (p in main) pkgs[p]=1
  for (p in curr) pkgs[p]=1
  n=asorti(pkgs, sorted)
  for (i=1; i<=n; i++) {
    p = sorted[i]
    mc = (p in main)?main[p]:0
    cc = (p in curr)?curr[p]:0
    diff = cc - mc
    printf "| %s | %.1f | %.1f | %+0.1f | %.1f |\n", p, mc, cc, diff, min
  }
  diff_total = total_curr - total_main
  printf "| **Total** | %.1f | %.1f | %+0.1f | %.1f |\n", total_main, total_curr, diff_total, min
}' main.out.func coverage.out.func
