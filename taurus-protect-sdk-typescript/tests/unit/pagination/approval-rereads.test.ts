/**
 * Approval re-reads obey the page-size maximum: a batch larger than one page is re-read
 * in id-filtered pages of at most 100, never as one oversized page.
 */

import * as crypto from "crypto";
import { IntegrityError } from "../../../src/errors";
import {
  WhitelistedAddressApproval,
} from "../../../src/models/whitelisted-address";
import { WhitelistedAssetApproval } from "../../../src/models/whitelisted-asset";
import { RecordingTransport, pagedServices } from "./harness";

const { privateKey } = crypto.generateKeyPairSync("ec", { namedCurve: "P-256" });

function ids(n: number): string[] {
  return Array.from({ length: n }, (_v, i) => String(i + 1));
}

describe("approval re-reads page in chunks of 100", () => {
  it("whitelisted addresses: 150 ids are re-read as 100 + 50", async () => {
    const transport = new RecordingTransport({});
    const selection = new WhitelistedAddressApproval(
      new Map(ids(150).map((id) => [id, `h${id}`]))
    );
    // The empty replies verify nothing, so the completeness check refuses to sign.
    await expect(
      pagedServices(transport).whitelistedAddresses.approve(selection, privateKey, "ok")
    ).rejects.toBeInstanceOf(IntegrityError);

    expect(transport.requests).toHaveLength(2);
    const [first, second] = transport.requests;
    const idsOf = (q: Array<[string, string]>): string[] =>
      q.filter(([k]) => k === "ids").map(([, v]) => v);
    expect(idsOf(first.query)).toEqual(ids(100));
    expect(idsOf(second.query)).toEqual(ids(150).slice(100));
    expect(first.query).toContainEqual(["limit", "100"]);
    expect(second.query).toContainEqual(["limit", "50"]);
    for (const request of transport.requests) {
      expect(request.query).toContainEqual(["includeForApproval", "true"]);
      expect(request.query).toContainEqual(["rulesContainerNormalized", "true"]);
      expect(request.path).toBe("/api/rest/v1/whitelists/addresses");
    }
  });

  it("whitelisted assets: 150 ids are re-read as 100 + 50", async () => {
    const transport = new RecordingTransport({});
    const selection = new WhitelistedAssetApproval(
      new Map(ids(150).map((id) => [Number(id), `h${id}`]))
    );
    await expect(
      pagedServices(transport).whitelistedAssets.approve(selection, privateKey, "ok")
    ).rejects.toBeInstanceOf(IntegrityError);

    expect(transport.requests).toHaveLength(2);
    const [first, second] = transport.requests;
    const idsOf = (q: Array<[string, string]>): string[] =>
      q.filter(([k]) => k === "whitelistedContractAddressIds").map(([, v]) => v);
    expect(idsOf(first.query)).toEqual(ids(100));
    expect(idsOf(second.query)).toEqual(ids(150).slice(100));
    expect(first.query).toContainEqual(["limit", "100"]);
    expect(second.query).toContainEqual(["limit", "50"]);
  });

  it("a batch that fits one page is re-read once", async () => {
    const transport = new RecordingTransport({});
    const selection = new WhitelistedAddressApproval(new Map([["7", "h7"], ["3", "h3"]]));
    await expect(
      pagedServices(transport).whitelistedAddresses.approve(selection, privateKey, "ok")
    ).rejects.toBeInstanceOf(IntegrityError);
    expect(transport.requests).toHaveLength(1);
    expect(transport.requests[0].query.filter(([k]) => k === "ids")).toEqual([
      ["ids", "3"],
      ["ids", "7"],
    ]);
  });
});
