/**
 * `BaseService.handleError` must not relabel an integrity failure as a retryable
 * server error.
 *
 * `WhitelistError`, `ConfigurationError` and `RequestMetadataError` all extend plain
 * `Error`, not `APIError` — so they fell to `handleError`'s final branch and came back
 * as `ServerError(500)`, whose `isRetryable()` is `true`. The consequence is not
 * cosmetic: the response that produced the failure is one the adversary controls, so
 * telling the caller to retry invites them to keep asking until an unverifiable answer
 * is accepted, and it hides a governance misconfiguration behind a "server error" that
 * looks transient. Both `getEnvelope` doc comments promise `@throws WhitelistError`,
 * which was unreachable through `execute`.
 *
 * The classification must match the sibling `rethrowIfNotRowLevel`
 * (`src/services/row-level-error.ts`), which already recognises `WhitelistError`.
 */

import {
  APIError,
  ConfigurationError,
  ContainerIntegrityError,
  IntegrityError,
  NotFoundError,
  RequestMetadataError,
  ServerError,
  UnverifiedMetadataError,
  ValidationError,
  WhitelistError,
} from "../../../src/errors";
import { BaseService } from "../../../src/services/base";

/**
 * Minimal concrete service: `execute` is protected, so a subclass is the only way to
 * drive the real error path without a stub API.
 */
class ProbeService extends BaseService {
  async run<T>(apiCall: () => Promise<T>): Promise<T> {
    return this.execute(apiCall);
  }
}

describe("BaseService error classification", () => {
  let service: ProbeService;

  beforeEach(() => {
    service = new ProbeService();
  });

  describe("integrity and configuration failures pass through untouched", () => {
    const passThrough: ReadonlyArray<[string, Error]> = [
      ["WhitelistError", new WhitelistError("group threshold not met")],
      [
        "ConfigurationError",
        new ConfigurationError("minValidSignatures must be positive"),
      ],
      [
        "RequestMetadataError",
        new RequestMetadataError("sourceAddress not in payload"),
      ],
      [
        "UnverifiedMetadataError",
        new UnverifiedMetadataError("metadata not verified"),
      ],
      ["IntegrityError", new IntegrityError("hash mismatch")],
      [
        "ContainerIntegrityError",
        new ContainerIntegrityError("cannot interpret rules container"),
      ],
      ["APIError", new NotFoundError("no such wallet")],
      ["ValidationError", new ValidationError("id must be positive")],
    ];

    it.each(passThrough)("%s reaches the caller as itself", async (label, thrown) => {
      // Tuple assertions: jest has no per-assertion message argument, so the label
      // has to travel with the value to stay identifiable in a failure.
      const caught = await service.run(() => Promise.reject(thrown)).then(
        () => undefined,
        (e: unknown) => e
      );

      expect([label, caught]).toEqual([label, thrown]);
    });

    it.each(passThrough)("%s is not turned into a ServerError", async (label, thrown) => {
      const caught = await service.run(() => Promise.reject(thrown)).then(
        () => undefined,
        (e: unknown) => e
      );

      // ServerError is what the final branch produces; NotFoundError/ValidationError
      // are APIErrors already and must not be re-wrapped either.
      expect([label, caught instanceof ServerError]).toEqual([label, false]);
    });

    it.each([
      ["WhitelistError", new WhitelistError("group threshold not met")],
      [
        "ConfigurationError",
        new ConfigurationError("minValidSignatures must be positive"),
      ],
      [
        "RequestMetadataError",
        new RequestMetadataError("sourceAddress not in payload"),
      ],
    ] as ReadonlyArray<[string, Error]>)(
      "%s is never advertised as retryable",
      async (label, thrown) => {
        const caught = await service.run(() => Promise.reject(thrown)).then(
          () => undefined,
          (e: unknown) => e
        );

        // The whole point: an attacker-controlled response must not come back with
        // "try me again". These types are not APIErrors, so they carry no
        // isRetryable() at all — which is the correct answer.
        const retryable =
          caught instanceof APIError ? caught.isRetryable() : false;
        expect([label, retryable]).toEqual([label, false]);
      }
    );
  });

  describe("genuine defects and unknown throws still become ServerError", () => {
    it("wraps a TypeError from a real defect", async () => {
      // Must NOT be reclassified as an integrity failure: a defect reported as
      // "this row failed verification" is indistinguishable from a tampered row.
      const caught = await service
        .run(() => Promise.reject(new TypeError("cannot read property 'x'")))
        .then(
          () => undefined,
          (e: unknown) => e
        );

      expect(caught).toBeInstanceOf(ServerError);
      expect((caught as ServerError).statusCode).toBe(500);
      expect((caught as ServerError).message).toMatch(/cannot read property 'x'/);
    });

    it("wraps a non-Error throw", async () => {
      const caught = await service.run(() => Promise.reject("a string")).then(
        () => undefined,
        (e: unknown) => e
      );

      expect(caught).toBeInstanceOf(ServerError);
      expect((caught as ServerError).message).toBe("Unknown error occurred");
    });

    it("returns the value when nothing throws", async () => {
      await expect(service.run(() => Promise.resolve(42))).resolves.toBe(42);
    });
  });
});
