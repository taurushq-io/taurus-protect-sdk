"""Loader for the cross-SDK lossless/parity vectors.

The scenarios in governance-lossless-vectors.json cannot be produced by this SDK's
encoder — they are deliberately non-canonical or schema-newer wire bytes — so they are
not in governance-cell-vectors.json, which the Go SDK generates. They used to live as
base64 literals hand-copied into all four SDKs' round-trip suites, where nothing
compared the copies and they could drift silently.
"""

from __future__ import annotations

import json
from pathlib import Path
from typing import Dict, List

VECTORS_PATH = (
    Path(__file__).resolve().parents[4]
    / "scripts"
    / "resources"
    / "governance-lossless-vectors.json"
)

# Asserted so a vector added to the shared file without being consumed here fails
# loudly instead of being silently ignored by this SDK.
VECTOR_COUNT = 9


def load_lossless_vectors() -> List[Dict[str, str]]:
    vectors = json.loads(VECTORS_PATH.read_text())
    assert len(vectors) == VECTOR_COUNT, (
        f"shared lossless vectors file has {len(vectors)} entries, expected "
        f"{VECTOR_COUNT} — update every SDK's suite in lockstep"
    )
    return vectors


def vectors_for(scenario: str) -> List[Dict[str, str]]:
    out = [v for v in load_lossless_vectors() if v["scenario"] == scenario]
    assert out, f"no shared lossless vector for scenario {scenario!r}"
    return out


def vector_for(scenario: str) -> str:
    out = vectors_for(scenario)
    assert len(out) == 1, f"scenario {scenario!r} has {len(out)} vectors, expected 1"
    return out[0]["wire_base64"]
