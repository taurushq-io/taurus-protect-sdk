/**
 * Shared fixtures for whitelisted-asset verification tests.
 *
 * Builds a genuinely signed envelope — real P-256 keys, a real SuperAdmin
 * signature over the rules container, and a real user signature over the JSON
 * hashes array — so a test that passes here proves the verification path ran,
 * not that it was skipped.
 *
 *	generateP256KeyPair ──▶ payload ──▶ SHA-256 ──▶ metadataHash
 *	                          │                          │
 *	                          │                          ├──▶ user signs [hash]
 *	                          └──▶ rules container ──────┴──▶ SuperAdmin signs
 *
 * Consumed by both `helpers/whitelisted-asset-verifier.test.ts` (step-level
 * cases) and `services/whitelisted-asset-service.test.ts` (read-path cases).
 */

import * as crypto from "crypto";

import { calculateHexHash, signData, encodePublicKeyPem } from "../../../src/crypto";
import type {
  DecodedRulesContainer,
  RuleUserSignature,
} from "../../../src/models/governance-rules";
import type { SignedWhitelistedAssetEnvelope } from "../../../src/models/whitelisted-asset";

export function generateP256KeyPair(): {
  privateKey: crypto.KeyObject;
  publicKey: crypto.KeyObject;
} {
  const { privateKey, publicKey } = crypto.generateKeyPairSync("ec", {
    namedCurve: "P-256",
  });
  return { privateKey, publicKey };
}

export function keyToPem(publicKey: crypto.KeyObject): string {
  return encodePublicKeyPem(publicKey);
}

export function buildAssetPayload(
  overrides?: Partial<{
    blockchain: string;
    network: string;
    contractAddress: string;
    name: string;
    symbol: string;
    decimals: number;
    isNFT: boolean;
    kindType: string;
  }>
): Record<string, unknown> {
  return {
    blockchain: "ETH",
    network: "mainnet",
    contractAddress: "0xUSDC",
    name: "USDC",
    symbol: "USDC",
    decimals: 6,
    ...overrides,
  };
}

export function payloadToString(payload: Record<string, unknown>): string {
  return JSON.stringify(payload);
}

export interface AssetTestFixture {
  saPriv: crypto.KeyObject;
  saPub: crypto.KeyObject;
  userPriv: crypto.KeyObject;
  userPub: crypto.KeyObject;
  saPem: string;
  userPem: string;
  envelope: SignedWhitelistedAssetEnvelope;
  rulesContainerDecoder: (b64: string) => DecodedRulesContainer;
  userSignaturesDecoder: (b64: string) => RuleUserSignature[];
}

export function buildFullAssetFixture(overrides?: {
  blockchain?: string;
  network?: string;
  groupId?: string;
  userId?: string;
}): AssetTestFixture {
  const { privateKey: saPriv, publicKey: saPub } = generateP256KeyPair();
  const { privateKey: userPriv, publicKey: userPub } = generateP256KeyPair();
  const saPem = keyToPem(saPub);
  const userPem = keyToPem(userPub);

  const blockchain = overrides?.blockchain ?? "ETH";
  const network = overrides?.network ?? "mainnet";
  const groupId = overrides?.groupId ?? "approvers";
  const userId = overrides?.userId ?? "user1@bank.com";

  const payload = buildAssetPayload({ blockchain, network });
  const payloadStr = payloadToString(payload);
  const metadataHash = calculateHexHash(payloadStr);

  const rulesJson = JSON.stringify({
    users: [{ id: userId, publicKey: userPem, roles: ["USER"] }],
    groups: [{ id: groupId, userIds: [userId] }],
    contractAddressWhitelistingRules: [
      {
        blockchain,
        network,
        parallelThresholds: [{ groupId, minimumSignatures: 1 }],
      },
    ],
  });
  const rulesB64 = Buffer.from(rulesJson).toString("base64");
  const rulesData = Buffer.from(rulesB64, "base64");

  const saSig = signData(saPriv, rulesData);

  const hashes = [metadataHash];
  const hashesJson = JSON.stringify(hashes);
  const userSig = signData(userPriv, Buffer.from(hashesJson, "utf-8"));

  const envelope: SignedWhitelistedAssetEnvelope = {
    id: 1,
    metadata: {
      hash: metadataHash,
      payloadAsString: payloadStr,
    },
    rulesContainerBase64: rulesB64,
    rulesSignaturesBase64: Buffer.from("dummy").toString("base64"),
    signedContractAddress: {
      payload: undefined,
      signatures: [
        {
          userSignature: {
            userId,
            signature: userSig,
            comment: undefined,
          },
          hashes,
        },
      ],
    },
    blockchain,
    network,
  };

  const rulesContainerDecoder = (_b64: string): DecodedRulesContainer => ({
    users: [
      {
        id: userId,
        name: "User 1",
        publicKeyPem: userPem,
        roles: ["USER"],
      },
    ],
    groups: [{ id: groupId, name: "Approvers", userIds: [userId] }],
    addressWhitelistingRules: [],
    contractAddressWhitelistingRules: [
      {
        blockchain,
        network,
        parallelThresholds: [
          {
            thresholds: [{ groupId, minimumSignatures: 1, threshold: 0 }],
          },
        ],
      },
    ],
    transactionRules: [],
    minimumDistinctUserSignatures: 0,
    minimumDistinctGroupSignatures: 0,
    enforcedRulesHash: "",
    timestamp: 0,
    hsmSlotId: 0,
    minimumCommitmentSignatures: 0,
    engineIdentities: [],
  });

  const userSignaturesDecoder = (_b64: string): RuleUserSignature[] => [
    { userId: "sa@bank.com", signature: saSig },
  ];

  return {
    saPriv,
    saPub,
    userPriv,
    userPub,
    saPem,
    userPem,
    envelope,
    rulesContainerDecoder,
    userSignaturesDecoder,
  };
}

/**
 * Shapes a fixture envelope into the DTO the ContractWhitelisting API returns,
 * so service tests exercise the real DTO → envelope → verify path rather than
 * handing the service an already-built envelope.
 */
export function fixtureToDto(
  envelope: SignedWhitelistedAssetEnvelope,
  overrides?: { payloadAsString?: string; id?: string }
): Record<string, unknown> {
  return {
    id: overrides?.id ?? String(envelope.id),
    metadata: {
      hash: envelope.metadata.hash,
      payloadAsString:
        overrides?.payloadAsString ?? envelope.metadata.payloadAsString,
    },
    rulesContainer: envelope.rulesContainerBase64,
    rulesSignatures: envelope.rulesSignaturesBase64,
    signedContractAddress: {
      payload: envelope.signedContractAddress.payload,
      signatures: (envelope.signedContractAddress.signatures ?? []).map((s) => ({
        signature: s.userSignature
          ? {
              userId: s.userSignature.userId,
              signature: s.userSignature.signature,
              comment: s.userSignature.comment,
            }
          : undefined,
        hashes: s.hashes,
      })),
    },
    blockchain: envelope.blockchain,
    network: envelope.network,
  };
}
