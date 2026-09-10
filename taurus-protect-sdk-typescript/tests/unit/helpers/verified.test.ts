/**
 * The verified brand is a compile-time guard, so most of this file is type-level.
 *
 * `@ts-expect-error` is the assertion: ts-jest fails the suite if the marked line
 * compiles cleanly. So each one proves the unbranded value is REJECTED — if the
 * brand stopped working, these lines would compile and the suite would go red.
 */

import { attestVerified, type Verified } from "../../../src/helpers/verified";
import type { SignedWhitelistedAssetEnvelope } from "../../../src/models/whitelisted-asset";

function readsOnlyVerified(
  envelope: Verified<SignedWhitelistedAssetEnvelope>
): string {
  return envelope.metadata.hash;
}

const bare: SignedWhitelistedAssetEnvelope = {
  id: 1,
  metadata: { hash: "abc", payloadAsString: "{}" },
  rulesContainerBase64: "",
  rulesSignaturesBase64: "",
  signedContractAddress: { payload: undefined, signatures: [] },
  blockchain: "ETH",
  network: "mainnet",
};

describe("Verified brand", () => {
  it("rejects an envelope that has not been through the verifier", () => {
    // @ts-expect-error a bare envelope is not Verified<...>
    readsOnlyVerified(bare);

    // The call still runs at runtime — the brand erases — which is the point:
    // the guarantee is structural, not behavioural.
    expect(true).toBe(true);
  });

  it("accepts a value the verifier attested", () => {
    expect(readsOnlyVerified(attestVerified(bare))).toBe("abc");
  });

  it("erases at runtime: the branded value is the same object", () => {
    const branded = attestVerified(bare);
    expect(branded).toBe(bare);
  });

  it("does not let a plain object masquerade as verified", () => {
    // @ts-expect-error the brand is a unique symbol; no literal can supply it
    const forged: Verified<SignedWhitelistedAssetEnvelope> = { ...bare };
    expect(forged).toBeDefined();
  });
});
