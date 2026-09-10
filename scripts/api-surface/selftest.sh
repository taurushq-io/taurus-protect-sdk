#!/usr/bin/env bash
# Proves the differ actually detects a gap. A differ that always prints "OK" is worse
# than no differ, because the report reads as evidence of parity.
set -uo pipefail

HERE="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$HERE/../.." && pwd)"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

mkdir -p "$TMP/scripts/resources" "$TMP/scripts/api-surface"
cp "$HERE"/testdata/known-gap/*.json "$TMP/scripts/resources/"
cp "$HERE/diff.py" "$HERE/aliases.json" "$TMP/scripts/api-surface/"

output="$(python3 "$TMP/scripts/api-surface/diff.py" 2>&1)"
status=$?

if [ "$status" -eq 0 ]; then
  echo "SELFTEST FAILED: differ reported success on a surface set with a known gap"
  echo "$output"
  exit 1
fi
if ! grep -q "missing from typescript" <<<"$output"; then
  echo "SELFTEST FAILED: differ did not name the missing SDK"
  echo "$output"
  exit 1
fi

echo "differ selftest: OK (gap detected and attributed)"
