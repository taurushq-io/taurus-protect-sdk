/**
 * Unit tests for the Credentials authentication mechanisms.
 */

import type { ProtectClientConfig } from "../../src/client";
import { ProtectClient } from "../../src/client";
import { Credentials } from "../../src/credentials";
import { ConfigurationError } from "../../src/errors";

// A valid P-256 public key for testing (test key with no production value)
const TEST_SUPER_ADMIN_KEY_PEM = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEM2NtzaFhm7xIR3OvWq5chW3/GEvW
L+3uqoE6lEJ13eWbulxsP/5h36VCqYDIGN/0wDeWwLYdpu5HhSXWhxCsCA==
-----END PUBLIC KEY-----`;

const HOST = "https://protect.example.com";
const HEX_SECRET = "aabbccdd11223344556677889900aabbccddeeff";

describe("Credentials", () => {
  describe("factory validation", () => {
    it("apiKey rejects empty key or secret", () => {
      expect(() => Credentials.apiKey("", HEX_SECRET)).toThrow("apiKey is required");
      expect(() => Credentials.apiKey("k", "")).toThrow("apiSecret is required");
    });

    it("apiKey rejects a non-hex secret", () => {
      expect(() => Credentials.apiKey("k", "not-hex")).toThrow(ConfigurationError);
    });

    it("bearerToken rejects an empty token", () => {
      expect(() => Credentials.bearerToken("")).toThrow("bearerToken is required");
    });
  });

  describe("ProtectClient.create with credentials", () => {
    it("builds a bearer client with SuperAdmin keys", () => {
      const client = ProtectClient.create({
        host: HOST,
        credentials: Credentials.bearerTokenProvider(() => "tok"),
        superAdminKeysPem: [TEST_SUPER_ADMIN_KEY_PEM],
      });
      expect(client.isClosed).toBe(false);
      client.close();
    });

    it("builds an api-key client via Credentials", () => {
      const client = ProtectClient.create({
        host: HOST,
        credentials: Credentials.apiKey("k", HEX_SECRET),
        superAdminKeysPem: [TEST_SUPER_ADMIN_KEY_PEM],
      });
      expect(client.isClosed).toBe(false);
      client.close();
    });

    // close() must release the credential material for a client built the SUPPORTED way.
    // The secret is captured inside the auth middleware closure, and a property walk
    // cannot see closure scope — so the observable invariant is that the client no longer
    // references the closure at all. Overwriting config.apiSecret does not satisfy this:
    // under `credentials:` that field is undefined, so before the fix close() released
    // nothing and the auth closure stayed live on apiConfiguration.
    it("releases the auth closure on close, so no credential material stays reachable", () => {
      const client = ProtectClient.create({
        host: HOST,
        credentials: Credentials.apiKey("k", HEX_SECRET),
        superAdminKeysPem: [TEST_SUPER_ADMIN_KEY_PEM],
      });
      const middleware = (
        client as unknown as { apiConfiguration: { middleware: unknown[] } }
      ).apiConfiguration.middleware;

      expect(middleware.length).toBe(1);

      client.close();

      expect(middleware.length).toBe(0);
    });

    it("releases the auth closure on close under bearer auth too", () => {
      const client = ProtectClient.create({
        host: HOST,
        credentials: Credentials.bearerToken("tok"),
        superAdminKeysPem: [TEST_SUPER_ADMIN_KEY_PEM],
      });
      const middleware = (
        client as unknown as { apiConfiguration: { middleware: unknown[] } }
      ).apiConfiguration.middleware;

      client.close();

      expect(middleware.length).toBe(0);
    });

    it("requires SuperAdmin keys even under bearer auth", () => {
      // Omitting superAdminKeysPem is a compile error — the field is required, so the
      // mandatory-keys invariant is caught in the consumer's build, not in production.
      // An empty array is the reachable typed case; untyped JS callers reach the guard too.
      expect(() =>
        ProtectClient.create({
          host: HOST,
          credentials: Credentials.bearerToken("tok"),
          superAdminKeysPem: [],
        })
      ).toThrow("superAdminKeysPem is required");

      expect(() =>
        ProtectClient.create({
          host: HOST,
          credentials: Credentials.bearerToken("tok"),
        } as unknown as ProtectClientConfig)
      ).toThrow("superAdminKeysPem is required");
    });
  });
});
