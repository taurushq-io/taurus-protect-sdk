/**
 * A compile-time marker that a value has passed verification.
 *
 * The type system could not previously tell an envelope straight off the wire from
 * one the verifier had cleared — both were `SignedWhitelistedAssetEnvelope` — so a
 * mapper accepted either. That is exactly how `WhitelistedAssetService.get()` came
 * to return attacker-controllable payload data labelled "verified".
 *
 *	         wire ──▶ SignedWhitelistedAssetEnvelope
 *	                            │
 *	                     verify()│  (the ONLY producer)
 *	                            ▼
 *	              Verified<SignedWhitelistedAssetEnvelope>
 *	                            │
 *	                            ▼
 *	                  mapper / accessor  ← accepts nothing else
 *
 * `attestVerified` is not exported from the package barrel, so a caller outside
 * these helpers cannot mint one without an explicit `as` cast.
 *
 * Erases completely: `Verified<T>` is `T` at runtime, so this costs nothing.
 *
 * What it proves and what it does not: that verification RAN, not that it ran
 * against the right keys or thresholds. It is a structural guard against forgetting
 * to call the verifier — the bug class above — not a cryptographic one.
 */

declare const verifiedBrand: unique symbol;

/** `T`, carrying proof that the verifier produced it. */
export type Verified<T> = T & { readonly [verifiedBrand]: true };

/**
 * Mints a {@link Verified} value. **Call this only after verification has passed** —
 * it is the single point where the marker is applied, so it is also the single point
 * to audit.
 *
 * @param value - the value the caller has just verified
 * @returns the same value, branded
 */
export function attestVerified<T>(value: T): Verified<T> {
  return value as Verified<T>;
}
