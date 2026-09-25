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

/**
 * What a by-id re-read verified, and which rows failed their own check.
 *
 * For a caller that completes unsigned rows through a verified reader and must drop a
 * bad row rather than fail its page (the v2 asset-holders list).
 */
export interface VerifiedLookup<T> {
  /** The rows that verified, by id. */
  readonly verified: ReadonlyMap<string, T>;
  /** The rows that failed their own check, by id, with the reason. */
  readonly failed: ReadonlyMap<string, string>;
}

/**
 * Verifies each row of a by-id re-read, absorbing only per-row integrity failures.
 *
 * Anything {@link rethrowIfNotRowLevel} re-throws aborts the read. An id returned more
 * than once is reported as failed: a by-id read has one row per id, so a second one is
 * a server choosing which copy the caller keeps.
 *
 * @param rows - The reply rows
 * @param idOf - The row's id
 * @param verify - The verified reader's per-row check; undefined means "no row"
 */
export async function verifyRowsById<D, T>(
  rows: readonly D[],
  idOf: (row: D) => string,
  verify: (row: D) => T | undefined | Promise<T | undefined>
): Promise<VerifiedLookup<T>> {
  const verified = new Map<string, T>();
  const failed = new Map<string, string>();
  const seen = new Set<string>();
  for (const row of rows) {
    const id = idOf(row);
    if (seen.has(id)) {
      verified.delete(id);
      failed.set(id, `row ${id} was returned more than once`);
      continue;
    }
    seen.add(id);
    try {
      const value = await verify(row);
      if (value !== undefined) {
        verified.set(id, value);
      }
    } catch (error: unknown) {
      rethrowIfNotRowLevel(error);
      failed.set(id, error.message);
    }
  }
  return { verified, failed };
}
