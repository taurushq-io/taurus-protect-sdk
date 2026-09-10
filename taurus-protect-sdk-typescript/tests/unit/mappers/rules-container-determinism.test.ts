/**
 * Cross-SDK determinism of the `properties` map.
 *
 * `properties` is `map<string, bytes>` and appears at seven levels of the rules
 * container. Its encoded order is part of the bytes a SuperAdmin signs on a rules
 * proposal, so all four SDKs must agree on it.
 *
 *   Go    proto.MarshalOptions{Deterministic: true}  ─┐
 *   Java  useDeterministicSerialization               ├─ sort keys lexicographically
 *   Py    SerializeToString(deterministic=True)      ─┘
 *   TS    Object.entries(...)  ← enumerates INTEGER-LIKE KEYS FIRST, ascending
 *
 * So for keys {"10","2","a"} the other three write 10,2,a and stock ts-proto writes
 * 2,10,a. No amount of sorting in the caller fixes that: JS object key order is
 * spec-fixed, and re-inserting keys in sorted order is silently discarded for
 * numeric-like keys. `scripts/generate-proto.sh` therefore patches the generated
 * encoders to sort at the emit site; these tests fail if that patch is ever lost
 * to a regeneration.
 */

import { createEmptyRulesContainer } from "../../../src/models/governance-rules";
import type { DecodedRulesContainer } from "../../../src/models/governance-rules";
import { rulesContainerToBytes } from "../../../src/mappers/protobuf-rules-container-encode";

/** Recovers the on-wire order of the given keys from the encoded container. */
function wireOrder(container: DecodedRulesContainer, keys: string[]): string[] {
  const wire = Buffer.from(rulesContainerToBytes(container)).toString("latin1");
  return keys
    .map((k) => ({ k, at: wire.indexOf(k) }))
    .filter((e) => e.at >= 0)
    .sort((a, b) => a.at - b.at)
    .map((e) => e.k);
}

function withProperties(keysInInsertionOrder: string[]): DecodedRulesContainer {
  const properties: { [k: string]: Uint8Array } = {};
  for (const k of keysInInsertionOrder) {
    properties[k] = new Uint8Array([1]);
  }
  return { ...createEmptyRulesContainer(), properties };
}

describe("properties map determinism", () => {
  it("writes numeric-like keys in lexicographic order, matching the other three SDKs", () => {
    // Lexicographic: "10" < "2" because '1' < '2'. Stock ts-proto emits 2,10,a.
    expect(wireOrder(withProperties(["10", "2", "a"]), ["10", "2", "a"])).toEqual([
      "10",
      "2",
      "a",
    ]);
  });

  it("is independent of the order the object was built in", () => {
    const a = rulesContainerToBytes(withProperties(["a", "2", "10"]));
    const b = rulesContainerToBytes(withProperties(["10", "a", "2"]));
    const c = rulesContainerToBytes(withProperties(["2", "10", "a"]));
    expect(Buffer.from(a)).toEqual(Buffer.from(b));
    expect(Buffer.from(b)).toEqual(Buffer.from(c));
  });

  it("keeps purely non-numeric keys sorted too", () => {
    expect(
      wireOrder(withProperties(["zebra", "alpha", "money"]), [
        "zebra",
        "alpha",
        "money",
      ])
    ).toEqual(["alpha", "money", "zebra"]);
  });
});
