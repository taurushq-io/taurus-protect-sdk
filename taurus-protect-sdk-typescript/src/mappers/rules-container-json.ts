/**
 * Canonical proto-JSON <-> base64 bridge for governance rules.
 *
 * Mirrors the Go SDK's rules_container_json.go: operates on the generated
 * ts-proto messages (encode/decode + fromJSON/toJSON), independent of the
 * DecodedRulesContainer convenience model, so callers can round-trip a rules
 * container or an individual rule-cell message through canonical protobuf JSON.
 */

import * as requestReply from "../internal/proto/request_reply";
import { strictBase64Decode } from "../helpers/strict-base64";

const { RulesContainer } = requestReply;

/** Minimal shape of a ts-proto message codec, for dynamic message lookup. */
interface ProtoCodec {
  encode(message: unknown): { finish(): Uint8Array };
  decode(input: Uint8Array): unknown;
  fromJSON(object: unknown): unknown;
  toJSON(message: unknown): unknown;
}

/** Decode a base64 protobuf RulesContainer into canonical proto JSON. */
export function rulesContainerJsonFromBase64(base64Data: string): string {
  const message = RulesContainer.decode(strictBase64Decode(base64Data));
  return JSON.stringify(RulesContainer.toJSON(message));
}

/**
 * Encode canonical proto JSON into a base64 protobuf RulesContainer.
 *
 * Deterministic `properties` ordering is enforced inside the generated encoders
 * (patched by `scripts/generate-proto.sh`), so this path needs no pre-sort of its
 * own — a caller-side sort could not deliver it anyway, since JS enumerates
 * integer-like object keys first regardless of insertion order.
 */
export function rulesContainerBase64FromJson(json: string): string {
  const message = RulesContainer.fromJSON(JSON.parse(json));
  return Buffer.from(RulesContainer.encode(message).finish()).toString("base64");
}

/**
 * Encode a governance rule protobuf message (resolved by name) into base64.
 *
 * messageType is a concrete rule message name from request_reply.proto such as
 * RuleSource, RuleFiatAmountRange, RulesContainer_Line, or
 * RulesContainer_TransactionRules.
 */
export function ruleMessageBase64FromJson(
  messageType: string,
  json: string
): string {
  const codec = governanceRuleMessage(messageType);
  const message = codec.fromJSON(JSON.parse(json));
  return Buffer.from(codec.encode(message).finish()).toString("base64");
}

/**
 * Decode a base64 protobuf rule message (by name) into canonical proto JSON.
 * Returns undefined for empty input (an empty cell means "match any").
 */
export function ruleMessageJsonFromBase64(
  messageType: string,
  base64Data: string
): string | undefined {
  if (!base64Data || base64Data.trim() === "") {
    return undefined;
  }
  const bytes = strictBase64Decode(base64Data);
  if (bytes.length === 0) {
    return undefined;
  }
  const codec = governanceRuleMessage(messageType);
  return JSON.stringify(codec.toJSON(codec.decode(bytes)));
}

/**
 * Resolve a rule message codec by its generated name. ts-proto has no reflective
 * registry, but exports each message as a named codec const, so the module
 * namespace serves the same role as Go's proto registry (nested types use the
 * same underscore-flattened names, e.g. RulesContainer_Line).
 */
function governanceRuleMessage(messageType: string): ProtoCodec {
  const name = (messageType ?? "").trim();
  if (!name) {
    throw new Error("messageType cannot be empty");
  }
  const codec = (requestReply as Record<string, unknown>)[name] as
    | ProtoCodec
    | undefined;
  if (
    !codec ||
    typeof codec.decode !== "function" ||
    typeof codec.fromJSON !== "function"
  ) {
    throw new Error(`unsupported governance rule message type "${messageType}"`);
  }
  return codec;
}
