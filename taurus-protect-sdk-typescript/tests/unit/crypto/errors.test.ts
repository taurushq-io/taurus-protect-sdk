/**
 * The predicate that decides which errors a verification loop may swallow.
 *
 * The four copy-pasted predicates it replaces matched on `message.includes('key')`,
 * which swallowed ANY error containing the word — including a TypeError from a real
 * bug, silently reported as "this signature did not verify".
 */

import { isCryptoVerificationError } from "../../../src/crypto/errors";

describe("isCryptoVerificationError", () => {
  it("accepts OpenSSL and Node crypto error codes", () => {
    const ossl = Object.assign(new Error("bad decrypt"), { code: "ERR_OSSL_BAD_DECRYPT" });
    const crypto = Object.assign(new Error("nope"), { code: "ERR_CRYPTO_INVALID_KEY" });

    expect(isCryptoVerificationError(ossl)).toBe(true);
    expect(isCryptoVerificationError(crypto)).toBe(true);
  });

  it("accepts this SDK's own DER and PEM failures", () => {
    expect(isCryptoVerificationError(new Error("Invalid DER signature: too short"))).toBe(true);
    expect(isCryptoVerificationError(new Error("Failed to decode public key: unknown error"))).toBe(true);
  });

  // The regression this predicate exists to prevent.
  it("does NOT swallow a programming error that merely mentions a key", () => {
    const bug = new TypeError(
      "Cannot read properties of undefined (reading 'key')"
    );
    expect(isCryptoVerificationError(bug)).toBe(false);
  });

  it("does NOT swallow arbitrary errors containing 'signature' or 'Invalid'", () => {
    expect(isCryptoVerificationError(new Error("signature service unavailable"))).toBe(false);
    expect(isCryptoVerificationError(new Error("Invalid configuration supplied"))).toBe(false);
  });

  it("rejects non-Error values", () => {
    expect(isCryptoVerificationError("ERR_OSSL")).toBe(false);
    expect(isCryptoVerificationError(undefined)).toBe(false);
    expect(isCryptoVerificationError(null)).toBe(false);
  });
});
