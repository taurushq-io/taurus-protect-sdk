/**
 * Price signature verification.
 *
 * Rate and decimals feed amount conversion, so an unverified price is a wrong number a
 * caller acts on.
 */

import type { KeyObject } from "crypto";

import { verifySignature, decodePublicKeyPem } from "../crypto";
import { IntegrityError } from "../errors";
import type { DecodedRulesContainer } from "../models/governance-rules";
import type { Price } from "../models/price";

/**
 * Returns the exact bytes a price signature covers.
 *
 * This is validatord's CurrencyPrice JSON projection: the five fields it serialises, in
 * its field order, with no spaces. The key order IS the signed byte sequence — do not
 * reorder or add fields.
 */
export function priceSignedBytes(price: Price): Buffer {
  const canonical =
    `{"blockchain":${JSON.stringify(price.blockchain ?? "")},` +
    `"currencyFrom":${JSON.stringify(price.currencyFrom ?? "")},` +
    `"currencyTo":${JSON.stringify(price.currencyTo ?? "")},` +
    `"decimals":${JSON.stringify(price.decimals ?? "")},` +
    `"rate":${JSON.stringify(price.rate ?? "")}}`;
  return Buffer.from(canonical, "utf-8");
}

/** Public keys of every user carrying the PRICEUPDATER role. */
function priceUpdaterKeys(rulesContainer: DecodedRulesContainer): KeyObject[] {
  const keys: KeyObject[] = [];
  for (const user of rulesContainer.users ?? []) {
    if (!user.roles?.includes("PRICEUPDATER") || !user.publicKeyPem) {
      continue;
    }
    try {
      keys.push(decodePublicKeyPem(user.publicKeyPem));
    } catch {
      // A key this SDK cannot decode cannot verify anything.
    }
  }
  return keys;
}

/**
 * Checks a price against the PRICEUPDATER keys in a verified rules container.
 *
 * The container decides whether prices must be signed at all. It is SuperAdmin-verified,
 * so "this tenant has no price signer" is trustworthy, whereas "this price carries no
 * signatures" is not — which is why a stripped signatures array on a tenant that DOES
 * have a PRICEUPDATER is an error rather than a skip.
 *
 * @throws {@link IntegrityError} If no PRICEUPDATER signature verifies
 */
export function verifyPrice(
  price: Price,
  rulesContainer: DecodedRulesContainer
): void {
  if (!price) {
    throw new IntegrityError("price cannot be null or undefined");
  }
  if (!rulesContainer) {
    throw new IntegrityError(
      "rules container required for price signature verification"
    );
  }

  const keys = priceUpdaterKeys(rulesContainer);
  if (keys.length === 0) {
    // This tenant does not sign prices.
    return;
  }

  const signatures = price.signatures ?? [];
  if (signatures.length === 0) {
    throw new IntegrityError(
      `price ${price.currencyFrom}/${price.currencyTo} carries no signatures but the ` +
        `rules container configures a PRICEUPDATER`
    );
  }

  const data = priceSignedBytes(price);
  for (const sig of signatures) {
    if (!sig.signature) {
      continue;
    }
    for (const key of keys) {
      try {
        if (verifySignature(key, data, sig.signature)) {
          return;
        }
      } catch {
        // A malformed signature is simply not a match.
      }
    }
  }

  throw new IntegrityError(
    `no PRICEUPDATER signature verifies for price ` +
      `${price.currencyFrom}/${price.currencyTo} ` +
      `(${signatures.length} signature(s) offered)`
  );
}

/** Verifies each price, throwing on the first failure. */
export function verifyPrices(
  prices: Price[],
  rulesContainer: DecodedRulesContainer
): void {
  for (const price of prices ?? []) {
    verifyPrice(price, rulesContainer);
  }
}
