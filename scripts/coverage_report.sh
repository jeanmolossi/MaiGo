#!/usr/bin/env bash
set -euo pipefail

MODULE=$(go list -m)
MIN_COVERAGE=${MIN_COVERAGE:-0}

cov_branch() {
  local profile=$1
  go test -coverprofile="$profile" ./... >&2
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

if git remote get-url origin >/dev/null 2>&1; then
  git fetch origin main >/dev/null
  TMP=$(mktemp -d)
  trap 'rm -rf "$TMP" main.out.func coverage.out.func' EXIT
  git worktree add --detach "$TMP/main" origin/main >/dev/null
  pushd "$TMP/main" >/dev/null
  cov_branch main.out
  popd >/dev/null
  mv "$TMP/main/main.out.func" main.out.func
  git worktree remove --force "$TMP/main" >/dev/null
else
  echo "origin remote not found" >&2
  exit 1
fi

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
