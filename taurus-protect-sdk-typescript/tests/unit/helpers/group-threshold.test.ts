/**
 * Per-group step-5 threshold: distinct signers, not signature entries.
 *
 * The entries arrive in the server-supplied userSignatures blob, so counting them let a
 * duplicated entry from one group member satisfy an N-of-M group — promoting an
 * under-approved whitelist entry to approved.
 */

import * as crypto from "crypto";

import { signData, encodePublicKeyPem } from "../../../src/crypto";
import { WhitelistedAddressVerifier } from "../../../src/helpers/whitelisted-address-verifier";
import type {
  DecodedRulesContainer,
  GroupThreshold,
} from "../../../src/models/governance-rules";
import type { WhitelistSignatureEntry } from "../../../src/models/whitelisted-address";

const METADATA_HASH = "abc123";

function generateP256KeyPair(): {
  privateKey: crypto.KeyObject;
  publicKey: crypto.KeyObject;
} {
  return crypto.generateKeyPairSync("ec", { namedCurve: "prime256v1" });
}

function keyToPem(publicKey: crypto.KeyObject): string {
  return encodePublicKeyPem(publicKey);
}

/** A genuinely valid approval entry, as it arrives on the wire. */
function signedEntry(
  userId: string,
  privateKey: crypto.KeyObject
): WhitelistSignatureEntry {
  const hashes = [METADATA_HASH];
  const hashesJson = Buffer.from(JSON.stringify(hashes), "utf-8");
  return {
    hashes,
    userSignature: { userId, signature: signData(privateKey, hashesJson) },
  } as WhitelistSignatureEntry;
}

function container(userPems: Record<string, string>): DecodedRulesContainer {
  return {
    users: Object.entries(userPems).map(([id, publicKeyPem]) => ({
      id,
      publicKeyPem,
      roles: ["USER"],
    })),
    groups: [{ id: "approvers", userIds: Object.keys(userPems) }],
  } as unknown as DecodedRulesContainer;
}

/** Drives the group walk directly; returns the failure message, or null on success. */
function walk(
  rules: DecodedRulesContainer,
  signatures: WhitelistSignatureEntry[],
  minSigs: number
): string | null {
  const verifier = new WhitelistedAddressVerifier({
    superAdminKeysPem: [keyToPem(generateP256KeyPair().publicKey)],
    minValidSignatures: 1,
  });
  const threshold = {
    groupId: "approvers",
    minimumSignatures: minSigs,
    threshold: 0,
  } as GroupThreshold;

  const precomputed = new Map<number, Buffer>();
  signatures.forEach((sig, idx) => {
    precomputed.set(idx, Buffer.from(JSON.stringify(sig.hashes), "utf-8"));
  });

  // verifyGroupThreshold is private; the walk is what is under test.
  return (
    verifier as unknown as {
      verifyGroupThreshold: (
        t: GroupThreshold,
        c: DecodedRulesContainer,
        s: WhitelistSignatureEntry[],
        h: string,
        p: Map<number, Buffer>,
        k: Map<string, crypto.KeyObject>
      ) => string | null;
    }
  ).verifyGroupThreshold(
    threshold,
    rules,
    signatures,
    METADATA_HASH,
    precomputed,
    new Map()
  );
}

describe("per-group threshold counts distinct signers", () => {
  it("rejects one member's entry duplicated against a 2-of-N group", () => {
    const u1 = generateP256KeyPair();
    const u2 = generateP256KeyPair();
    const rules = container({
      "user1@bank.com": keyToPem(u1.publicKey),
      "user2@bank.com": keyToPem(u2.publicKey),
    });

    const entry = signedEntry("user1@bank.com", u1.privateKey);
    expect(walk(rules, [entry, entry], 2)).not.toBeNull();
  });

  it("rejects two re-signed entries from one member against a 2-of-N group", () => {
    const u1 = generateP256KeyPair();
    const u2 = generateP256KeyPair();
    const rules = container({
      "user1@bank.com": keyToPem(u1.publicKey),
      "user2@bank.com": keyToPem(u2.publicKey),
    });

    const signatures = [
      signedEntry("user1@bank.com", u1.privateKey),
      signedEntry("user1@bank.com", u1.privateKey),
    ];
    // ECDSA is randomized, so the same key yields two different signatures.
    expect(signatures[0]!.userSignature!.signature).not.toBe(
      signatures[1]!.userSignature!.signature
    );
    expect(walk(rules, signatures, 2)).not.toBeNull();
  });

  it("accepts two distinct members against a 2-of-N group", () => {
    const u1 = generateP256KeyPair();
    const u2 = generateP256KeyPair();
    const rules = container({
      "user1@bank.com": keyToPem(u1.publicKey),
      "user2@bank.com": keyToPem(u2.publicKey),
    });

    const signatures = [
      signedEntry("user1@bank.com", u1.privateKey),
      signedEntry("user2@bank.com", u2.privateKey),
    ];
    expect(walk(rules, signatures, 2)).toBeNull();
  });

  it("accepts one member against a 1-of-N group", () => {
    const u1 = generateP256KeyPair();
    const u2 = generateP256KeyPair();
    const rules = container({
      "user1@bank.com": keyToPem(u1.publicKey),
      "user2@bank.com": keyToPem(u2.publicKey),
    });

    expect(walk(rules, [signedEntry("user1@bank.com", u1.privateKey)], 1)).toBeNull();
  });

  // One compromised secret must not satisfy 2-of-N, whatever it is called.
  it("counts two user IDs sharing one key as one signer", () => {
    const shared = generateP256KeyPair();
    const pem = keyToPem(shared.publicKey);
    const rules = container({ "user1@bank.com": pem, "user2@bank.com": pem });

    const signatures = [
      signedEntry("user1@bank.com", shared.privateKey),
      signedEntry("user2@bank.com", shared.privateKey),
    ];
    expect(walk(rules, signatures, 2)).not.toBeNull();
  });

  it("ignores a valid signer outside the group", () => {
    const u1 = generateP256KeyPair();
    const u2 = generateP256KeyPair();
    const outsider = generateP256KeyPair();

    const rules = container({
      "user1@bank.com": keyToPem(u1.publicKey),
      "user2@bank.com": keyToPem(u2.publicKey),
    });
    rules.users.push({
      id: "outsider@bank.com",
      publicKeyPem: keyToPem(outsider.publicKey),
      roles: ["USER"],
    } as unknown as DecodedRulesContainer["users"][number]);

    const signatures = [
      signedEntry("user1@bank.com", u1.privateKey),
      signedEntry("outsider@bank.com", outsider.privateKey),
    ];
    expect(walk(rules, signatures, 2)).not.toBeNull();
  });
});
