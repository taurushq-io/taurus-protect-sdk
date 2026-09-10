/**
 * Price signature verification against the PRICEUPDATER role.
 *
 * Rate and decimals feed amount conversion, so an unverified price is a wrong number a
 * caller acts on. Whether prices must be signed is decided by the SuperAdmin-verified
 * rules container.
 */

import * as crypto from "crypto";

import { signData, encodePublicKeyPem } from "../../../src/crypto";
import { IntegrityError } from "../../../src/errors";
import {
  priceSignedBytes,
  verifyPrice,
  verifyPrices,
} from "../../../src/helpers/price-verifier";
import type { DecodedRulesContainer } from "../../../src/models/governance-rules";
import type { Price } from "../../../src/models/price";

function keyPair(): { privateKey: crypto.KeyObject; publicKey: crypto.KeyObject } {
  return crypto.generateKeyPairSync("ec", { namedCurve: "prime256v1" });
}

function price(overrides: Partial<Price> = {}): Price {
  return {
    blockchain: "ETH",
    currencyFrom: "ETH",
    currencyTo: "USD",
    decimals: "18",
    rate: "2500.00",
    ...overrides,
  } as Price;
}

function container(roles: string[], publicKey: crypto.KeyObject): DecodedRulesContainer {
  return {
    users: [
      { id: "price@bank.com", publicKeyPem: encodePublicKeyPem(publicKey), roles },
    ],
    groups: [],
  } as unknown as DecodedRulesContainer;
}

function signed(p: Price, privateKey: crypto.KeyObject): Price {
  return {
    ...p,
    signatures: [
      { userId: "price@bank.com", signature: signData(privateKey, priceSignedBytes(p)) },
    ],
  } as Price;
}

describe("price signature verification", () => {
  // The canonical form is validatord's CurrencyPrice JSON projection. Pinned by value: a
  // reordered or extended object silently changes every signature this SDK accepts.
  it("builds the canonical projection", () => {
    expect(priceSignedBytes(price()).toString("utf-8")).toBe(
      '{"blockchain":"ETH","currencyFrom":"ETH","currencyTo":"USD","decimals":"18","rate":"2500.00"}'
    );
  });

  it("accepts a PRICEUPDATER signature", () => {
    const { privateKey, publicKey } = keyPair();
    expect(() =>
      verifyPrice(signed(price(), privateKey), container(["PRICEUPDATER"], publicKey))
    ).not.toThrow();
  });

  it("rejects a signature from a key without the role", () => {
    const signer = keyPair();
    const updater = keyPair();
    expect(() =>
      verifyPrice(
        signed(price(), signer.privateKey),
        container(["PRICEUPDATER"], updater.publicKey)
      )
    ).toThrow(IntegrityError);
  });

  it("rejects a tampered rate", () => {
    const { privateKey, publicKey } = keyPair();
    const original = signed(price(), privateKey);
    // The rate is what a caller converts with.
    const tampered = { ...original, rate: "1.00" } as Price;

    expect(() =>
      verifyPrice(tampered, container(["PRICEUPDATER"], publicKey))
    ).toThrow(IntegrityError);
  });

  it("rejects stripped signatures when a PRICEUPDATER exists", () => {
    const { publicKey } = keyPair();
    expect(() => verifyPrice(price(), container(["PRICEUPDATER"], publicKey))).toThrow(
      IntegrityError
    );
  });

  it("passes through when no PRICEUPDATER is configured", () => {
    const { publicKey } = keyPair();
    expect(() =>
      verifyPrice(price(), container(["REQUESTAPPROVER"], publicKey))
    ).not.toThrow();
  });

  it("requires a rules container", () => {
    expect(() =>
      verifyPrice(price(), undefined as unknown as DecodedRulesContainer)
    ).toThrow(IntegrityError);
  });

  it("reports the first failure across a batch", () => {
    const { privateKey, publicKey } = keyPair();
    const good = signed(price(), privateKey);
    const bad = price({ currencyTo: "EUR" });

    expect(() =>
      verifyPrices([good, bad], container(["PRICEUPDATER"], publicKey))
    ).toThrow(IntegrityError);
  });
});
