#!/usr/bin/env bash
# Regenerates the shared API surface for every SDK, injects the generated method index
# into each docs/SERVICES.md, and runs the cross-SDK parity differ.
#
# Each SDK's ./build.sh docs calls this for its own language; run it with no argument to
# do all four plus the differ (which needs all four surfaces present).
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
RES="$ROOT/scripts/resources"
HERE="$ROOT/scripts/api-surface"
MODE="${2:-render}"   # render | check

surface_for() {
  local sdk="$1"
  echo "$RES/api-surface.$sdk.json"
}

extract() {
  local sdk="$1"
  case "$sdk" in
    go)
      (cd "$ROOT/taurus-protect-sdk-go" && go run "$HERE/extract-go.go" pkg/protect/service "$(surface_for go)")
      ;;
    python)
      local py="$ROOT/taurus-protect-sdk-python/.venv/bin/python"
      [ -x "$py" ] || py=python3
      (cd "$ROOT/taurus-protect-sdk-python" && "$py" "$HERE/extract-python.py" . "$(surface_for python)")
      ;;
    typescript)
      (cd "$ROOT/taurus-protect-sdk-typescript" && node "$HERE/extract-ts.mjs" . "$(surface_for typescript)")
      ;;
    java)
      # Needs the compiled classes; javap reads those, not the sources.
      (cd "$ROOT/taurus-protect-sdk-java" && python3 "$HERE/extract-java.py" . "$(surface_for java)")
      ;;
    *) echo "unknown sdk: $sdk" >&2; return 2 ;;
  esac
}

one() {
  local sdk="$1"
  extract "$sdk"
  python3 "$HERE/docs.py" "$MODE" "$sdk" "$(surface_for "$sdk")" \
    "$ROOT/taurus-protect-sdk-$sdk/docs/SERVICES.md"
}

if [ "${1:-all}" = "all" ]; then
  for sdk in go java python typescript; do one "$sdk"; done
  python3 "$HERE/diff.py"
else
  one "$1"
fi
