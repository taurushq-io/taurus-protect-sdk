import { RulesContainer } from "../../../src/internal/proto/request_reply";
import {
  ruleMessageBase64FromJson,
  ruleMessageJsonFromBase64,
  rulesContainerBase64FromJson,
  rulesContainerJsonFromBase64,
} from "../../../src/mappers/rules-container-json";

describe("rules container JSON bridge", () => {
  it("round-trips a rules container through canonical JSON", () => {
    const src = RulesContainer.fromJSON({ minimumDistinctUserSignatures: 3 });
    const encoded = Buffer.from(RulesContainer.encode(src).finish()).toString(
      "base64"
    );

    const json = rulesContainerJsonFromBase64(encoded);
    expect(json).toContain("minimumDistinctUserSignatures");

    const reEncoded = rulesContainerBase64FromJson(json);
    expect(reEncoded).toBe(encoded);
  });

  // The JSON bridge is the path the MCP governance tools drive: decode a container to
  // JSON, let a caller edit it, re-encode, submit — and approvers sign the re-encoded
  // bytes. An edited container arrives with its map keys in arbitrary order, so encoding
  // MUST be order-independent: ts-proto emits map entries in the object's own key order.
  // The input below carries the five `properties` keys in reverse-sorted order on purpose;
  // feeding back the canonical (already sorted) JSON would pass even with a
  // non-deterministic encoder and prove nothing. Expected bytes are the same vector
  // asserted for the typed encoder in rules-container-roundtrip.test.ts, and the same
  // reversed input is asserted in all four SDKs.
  it("encodes order-independently through the JSON bridge and matches the other SDKs byte for byte", () => {
    const expected =
      "CmgKAnUxEgNQRU0aAQMiEAoGa0FscGhhEgZrQWxwaGEiEAoGa0JyYXZvEgZrQnJhdm8iFAoIa0NoYXJsaWUSCGtDaGFybGllIhAKBmtEZWx0YRIGa0RlbHRhIg4KBWtFY2hvEgVrRWNob0oQCgZrQWxwaGESBmtBbHBoYUoQCgZrQnJhdm8SBmtCcmF2b0oUCghrQ2hhcmxpZRIIa0NoYXJsaWVKEAoGa0RlbHRhEgZrRGVsdGFKDgoFa0VjaG8SBWtFY2hv";
    const reversedKeyJson =
      '{"users":[{"id":"u1","publicKey":"PEM","roles":["SUPERADMIN"],"properties":' +
      '{"kEcho":"a0VjaG8=","kDelta":"a0RlbHRh","kCharlie":"a0NoYXJsaWU=",' +
      '"kBravo":"a0JyYXZv","kAlpha":"a0FscGhh"}}],"properties":' +
      '{"kEcho":"a0VjaG8=","kDelta":"a0RlbHRh","kCharlie":"a0NoYXJsaWU=",' +
      '"kBravo":"a0JyYXZv","kAlpha":"a0FscGhh"}}';

    for (let i = 0; i < 20; i++) {
      expect(rulesContainerBase64FromJson(reversedKeyJson)).toBe(expected);
    }
  });

  it("rejects invalid JSON when encoding a container", () => {
    expect(() => rulesContainerBase64FromJson("not json")).toThrow();
  });

  it("round-trips a rule message stably", () => {
    const encoded = ruleMessageBase64FromJson(
      "RuleFiatAmountRange",
      '{"minAmount":"1000"}'
    );
    const decoded = ruleMessageJsonFromBase64("RuleFiatAmountRange", encoded);
    expect(decoded).toBeDefined();
    expect(decoded).toContain("1000");
    expect(ruleMessageBase64FromJson("RuleFiatAmountRange", decoded!)).toBe(
      encoded
    );
  });

  it("returns undefined for empty rule message input", () => {
    expect(ruleMessageJsonFromBase64("RuleSource", "")).toBeUndefined();
    expect(ruleMessageJsonFromBase64("RuleSource", "   ")).toBeUndefined();
  });

  it("throws for an unknown rule message type", () => {
    expect(() => ruleMessageBase64FromJson("NoSuchRuleMessage", "{}")).toThrow();
  });

  it("throws for an empty rule message type", () => {
    expect(() => ruleMessageBase64FromJson("   ", "{}")).toThrow();
  });
});
