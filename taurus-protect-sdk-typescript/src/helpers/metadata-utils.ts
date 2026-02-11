import type { RequestMetadataAmount } from "../models/request";

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

function getPayloadValue(payloadAsString: string, key: string): Record<string, unknown> | undefined {
  for (const entry of parsePayloadEntries(payloadAsString)) {
    if (entry.key === key) {
      return entry.value;
    }
  }
  return undefined;
}

/** Extract source address from verified metadata payload. */
export function getSourceAddress(payloadAsString: string): string | undefined {
  const value = getPayloadValue(payloadAsString, "source");
  if (value) {
    const payload = value.payload as Record<string, unknown> | undefined;
    if (payload) {
      return payload.address as string | undefined;
    }
  }
  return undefined;
}

/** Extract destination address from verified metadata payload. */
export function getDestinationAddress(payloadAsString: string): string | undefined {
  const value = getPayloadValue(payloadAsString, "destination");
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

/** Extract amount information from verified metadata payload. */
export function getAmount(payloadAsString: string): RequestMetadataAmount | undefined {
  const value = getPayloadValue(payloadAsString, "amount");
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
