/**
 * Which verification failures a list path may absorb into `excludedUnverified`.
 *
 * Shared by the whitelisted-address and whitelisted-asset services, which run the same
 * lenient-list shape over different entities. Extracted rather than copied: the two
 * copies of this classification are exactly how the list paths drifted from `get()` in
 * the first place, and a private copy per service is what the repo's "one seam" rule
 * exists to prevent.
 */

import { ContainerIntegrityError, IntegrityError, WhitelistError } from "../errors";

/**
 * Re-throws anything that is not a per-row integrity failure.
 *
 * Two classes must never be swallowed into `excludedUnverified`:
 *
 * - `ContainerIntegrityError` — this SDK cannot interpret the rules container, which
 *   invalidates every row judged against it, not just this one. Excluding them one by
 *   one would empty the whitelist and report success. Go re-throws via `errors.As`,
 *   Python via `except ContainerIntegrityError: raise`, and Java by the exception being
 *   unchecked; this SDK's blanket `catch (error: unknown)` was the lone outlier.
 * - anything that is not an integrity failure at all — a `TypeError` from a genuine
 *   defect was being reported to the caller as "this address failed verification",
 *   indistinguishable from a tampered row.
 *
 * `ContainerIntegrityError extends IntegrityError`, so the narrower check comes first.
 *
 * The caller may treat the error as an `Error` after this returns: both surviving
 * branches are `Error` subclasses.
 */
export function rethrowIfNotRowLevel(error: unknown): asserts error is Error {
  if (error instanceof ContainerIntegrityError) {
    throw error;
  }
  if (!(error instanceof IntegrityError) && !(error instanceof WhitelistError)) {
    throw error;
  }
}
