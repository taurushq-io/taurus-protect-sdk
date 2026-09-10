/**
 * Payload accessors must refuse metadata that verification has not cleared.
 *
 * These helpers used to take the raw `payloadAsString`, which carries no record of
 * whether anything checked it — so they served data from unverified metadata and
 * the caller had no way to tell.
 */

import {
  getSourceAddress,
  getDestinationAddress,
  getAmount,
} from "../../../src/helpers/metadata-utils";
import { RequestMetadataError, UnverifiedMetadataError } from "../../../src/errors";
import type { RequestMetadata } from "../../../src/models/request";

const PAYLOAD = JSON.stringify([
  { key: "currency", value: "XLM" },
  { key: "source", value: { payload: { address: "SRC" } } },
  { key: "destination", value: { payload: { address: "DST" } } },
  { key: "amount", value: { valueFrom: "2", valueTo: "0.1", rate: "0.05", decimals: 7 } },
]);

const unverified: RequestMetadata = { hash: "deadbeef", payloadAsString: PAYLOAD };
const verified: RequestMetadata = { ...unverified, hashVerified: true };

describe("metadata accessors", () => {
  describe("unverified metadata", () => {
    it("refuses to yield the source address", () => {
      expect(() => getSourceAddress(unverified)).toThrow(UnverifiedMetadataError);
    });

    it("refuses to yield the destination address", () => {
      expect(() => getDestinationAddress(unverified)).toThrow(UnverifiedMetadataError);
    });

    it("refuses to yield the amount", () => {
      expect(() => getAmount(unverified)).toThrow(UnverifiedMetadataError);
    });

    it("is a RequestMetadataError, so existing catch blocks keep working", () => {
      expect(new UnverifiedMetadataError("x")).toBeInstanceOf(RequestMetadataError);
    });
  });

  describe("verified metadata", () => {
    it("serves every field", () => {
      expect(getSourceAddress(verified)).toBe("SRC");
      expect(getDestinationAddress(verified)).toBe("DST");
      expect(getAmount(verified)?.valueFrom).toBe("2");
    });

    // "not in the payload" and "not verified" must stay distinguishable.
    it("returns undefined for a field the payload does not carry", () => {
      const sparse: RequestMetadata = {
        hash: "deadbeef",
        payloadAsString: JSON.stringify([{ key: "currency", value: "XLM" }]),
        hashVerified: true,
      };
      expect(getSourceAddress(sparse)).toBeUndefined();
    });
  });

  // A request in an early status has nothing to verify and nothing to read.
  it("treats metadata without a payload as empty, not as an error", () => {
    expect(getSourceAddress(undefined)).toBeUndefined();
    expect(getSourceAddress({ hash: "", payloadAsString: undefined })).toBeUndefined();
  });
});
