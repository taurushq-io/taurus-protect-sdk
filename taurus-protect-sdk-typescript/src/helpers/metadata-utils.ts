import type { RequestMetadata, RequestMetadataAmount } from "../models/request";
import { UnverifiedMetadataError } from "../errors";

interface PayloadEntry {
  key: string;
  value: Record<string, unknown>;
}

function parsePayloadEntries(payloadAsString: string): PayloadEntry[] {
  try {
    const parsed = JSON.parse(payloadAsString);
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

/**
 * Returns the payload string, or throws if verification has not cleared it.
 *
 * These helpers used to take the raw `payloadAsString`, which carries no record of
 * whether anything checked it — so they served data from unverified metadata and a
 * caller could not tell. They take the metadata object now, as the Go, Java and
 * Python peers do, and read `hashVerified` from it.
 *
 * Metadata with no payload is not an error: a request in an early status has
 * nothing to verify and nothing to read.
 */
function verifiedPayload(metadata: RequestMetadata | undefined): string | undefined {
  if (!metadata?.payloadAsString) {
    return undefined;
  }
  if (!metadata.hashVerified) {
    throw new UnverifiedMetadataError(
      "request metadata payload has not been verified"
    );
  }
  return metadata.payloadAsString;
}

function getPayloadValue(
  metadata: RequestMetadata | undefined,
  key: string
): Record<string, unknown> | undefined {
  const payloadAsString = verifiedPayload(metadata);
  if (!payloadAsString) {
    return undefined;
  }
  for (const entry of parsePayloadEntries(payloadAsString)) {
    if (entry.key === key) {
      return entry.value;
    }
  }
  return undefined;
}

/**
 * Extract the source address from VERIFIED metadata.
 *
 * @throws {@link UnverifiedMetadataError} if the metadata carries an unverified payload
 */
export function getSourceAddress(metadata: RequestMetadata | undefined): string | undefined {
  const value = getPayloadValue(metadata, "source");
  if (value) {
    const payload = value.payload as Record<string, unknown> | undefined;
    if (payload) {
      return payload.address as string | undefined;
    }
  }
  return undefined;
}

/**
 * Extract the destination address from VERIFIED metadata.
 *
 * @throws {@link UnverifiedMetadataError} if the metadata carries an unverified payload
 */
export function getDestinationAddress(metadata: RequestMetadata | undefined): string | undefined {
  const value = getPayloadValue(metadata, "destination");
  if (value) {
    const payload = value.payload as Record<string, unknown> | undefined;
    if (payload) {
      return payload.address as string | undefined;
    }
  }
  return undefined;
}

/**
 * Safely converts a JSON value to a string representation.
 * Handles strings, numbers, null, and undefined gracefully.
 */
function jsonValueToString(value: unknown): string {
  if (value === null || value === undefined) return '';
  if (typeof value === 'string') return value;
  if (typeof value === 'number') return String(value);
  return String(value);
}

/**
 * Extract amount information from VERIFIED metadata.
 *
 * @throws {@link UnverifiedMetadataError} if the metadata carries an unverified payload
 */
export function getAmount(metadata: RequestMetadata | undefined): RequestMetadataAmount | undefined {
  const value = getPayloadValue(metadata, "amount");
  if (!value) {
    return undefined;
  }
  return {
    valueFrom: jsonValueToString(value.valueFrom),
    valueTo: jsonValueToString(value.valueTo),
    rate: jsonValueToString(value.rate),
    decimals: Number(value.decimals),
    currencyFrom: String(value.currencyFrom ?? ""),
    currencyTo: String(value.currencyTo ?? ""),
  };
}
