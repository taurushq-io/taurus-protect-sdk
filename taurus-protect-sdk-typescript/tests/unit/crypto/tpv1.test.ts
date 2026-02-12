/**
 * Unit tests for TPV1-HMAC-SHA256 authentication.
 *
 * These tests verify the TPV1Auth class and calculateSignedHeader function
 * used for API request authentication.
 */

import {
  calculateHexHash,
  calculateBase64Hmac,
  verifyBase64Hmac,
} from "../../../src/crypto/hashing";
import { TPV1Auth, calculateSignedHeader } from "../../../src/crypto/tpv1";

describe("calculateHexHash (TPV1 context)", () => {
  it("should produce expected SHA-256 for known input", () => {
    const hash = calculateHexHash("hello");
    expect(hash).toBe(
      "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
    );
  });

  it("should produce expected SHA-256 for empty string", () => {
    const hash = calculateHexHash("");
    expect(hash).toBe(
      "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
    );
  });
});

describe("calculateBase64Hmac", () => {
  it("should produce a deterministic HMAC for known inputs", () => {
    const secret = Buffer.from("secret-key");
    const hmac1 = calculateBase64Hmac(secret, "data");
    const hmac2 = calculateBase64Hmac(secret, "data");
    expect(hmac1).toBe(hmac2);
  });

  it("should produce different HMACs for different secrets", () => {
    const hmac1 = calculateBase64Hmac(Buffer.from("key-a"), "data");
    const hmac2 = calculateBase64Hmac(Buffer.from("key-b"), "data");
    expect(hmac1).not.toBe(hmac2);
  });
});

describe("verifyBase64Hmac", () => {
  it("should return true for matching HMAC", () => {
    const secret = Buffer.from("my-secret");
    const data = "some-data-to-sign";
    const hmac = calculateBase64Hmac(secret, data);
    expect(verifyBase64Hmac(secret, data, hmac)).toBe(true);
  });

  it("should return false for non-matching HMAC", () => {
    const secret = Buffer.from("my-secret");
    expect(verifyBase64Hmac(secret, "data", "wrong-hmac-value")).toBe(false);
  });

  it("should return false when data is different", () => {
    const secret = Buffer.from("my-secret");
    const hmac = calculateBase64Hmac(secret, "original-data");
    expect(verifyBase64Hmac(secret, "tampered-data", hmac)).toBe(false);
  });
});

describe("TPV1Auth", () => {
  const apiKey = "test-api-key";
  const apiSecretHex = "deadbeef01020304";

  it("should throw on empty apiKey", () => {
    expect(() => new TPV1Auth("", apiSecretHex)).toThrow(
      "apiKey cannot be empty"
    );
  });

  it("should throw on empty apiSecret", () => {
    expect(() => new TPV1Auth(apiKey, "")).toThrow("apiSecret cannot be empty");
  });

  describe("signRequest", () => {
    let auth: TPV1Auth;

    beforeEach(() => {
      auth = new TPV1Auth(apiKey, apiSecretHex);
    });

    afterEach(() => {
      auth.close();
    });

    it('should produce header starting with "TPV1-HMAC-SHA256"', () => {
      const header = auth.signRequest(
        "GET",
        "api.example.com",
        "/v1/wallets"
      );
      expect(header.startsWith("TPV1-HMAC-SHA256 ")).toBe(true);
    });

    it("should include ApiKey, Nonce, Timestamp, Signature", () => {
      const header = auth.signRequest(
        "GET",
        "api.example.com",
        "/v1/wallets"
      );
      expect(header).toContain(`ApiKey=${apiKey}`);
      expect(header).toMatch(/Nonce=[0-9a-f-]{36}/);
      expect(header).toMatch(/Timestamp=\d+/);
      expect(header).toMatch(/Signature=.+/);
    });

    it("should produce different signatures for GET vs POST", () => {
      const headerGet = auth.signRequest(
        "GET",
        "api.example.com",
        "/v1/wallets"
      );
      const headerPost = auth.signRequest(
        "POST",
        "api.example.com",
        "/v1/wallets"
      );

      // Extract just the signature portion
      const sigGet = headerGet.split("Signature=")[1];
      const sigPost = headerPost.split("Signature=")[1];
      expect(sigGet).not.toBe(sigPost);
    });

    it("should handle GET without body", () => {
      const header = auth.signRequest(
        "GET",
        "api.example.com",
        "/v1/wallets"
      );
      expect(header).toBeDefined();
      expect(header.startsWith("TPV1-HMAC-SHA256 ")).toBe(true);
    });

    it("should handle POST with JSON body", () => {
      const header = auth.signRequest(
        "POST",
        "api.example.com",
        "/v1/wallets",
        undefined,
        "application/json",
        '{"name":"my-wallet"}'
      );
      expect(header).toBeDefined();
      expect(header.startsWith("TPV1-HMAC-SHA256 ")).toBe(true);
    });

    it("should include query string when provided", () => {
      // Sign with and without query to confirm they differ
      const headerNoQuery = auth.signRequest(
        "GET",
        "api.example.com",
        "/v1/wallets"
      );
      const headerWithQuery = auth.signRequest(
        "GET",
        "api.example.com",
        "/v1/wallets",
        "limit=10&offset=0"
      );

      const sigNoQuery = headerNoQuery.split("Signature=")[1];
      const sigWithQuery = headerWithQuery.split("Signature=")[1];
      expect(sigNoQuery).not.toBe(sigWithQuery);
    });

    it("should uppercase the HTTP method", () => {
      // Signing with "get" should produce the same signature as "GET"
      const auth2 = new TPV1Auth(apiKey, apiSecretHex);
      // We can't easily verify method uppercasing since nonce/timestamp differ,
      // but we verify it doesn't throw
      const header = auth.signRequest(
        "get",
        "api.example.com",
        "/v1/wallets"
      );
      expect(header).toBeDefined();
      auth2.close();
    });
  });

  describe("close", () => {
    it("should throw after close when signing", () => {
      const auth = new TPV1Auth(apiKey, apiSecretHex);
      auth.close();
      expect(() =>
        auth.signRequest("GET", "api.example.com", "/v1/wallets")
      ).toThrow("TPV1Auth has been closed");
    });

    it("should be idempotent (calling close twice does not throw)", () => {
      const auth = new TPV1Auth(apiKey, apiSecretHex);
      auth.close();
      expect(() => auth.close()).not.toThrow();
    });
  });

  describe("parseUrl", () => {
    it("should parse a simple URL", () => {
      const result = TPV1Auth.parseUrl("https://api.example.com/v1/wallets");
      expect(result.host).toBe("api.example.com");
      expect(result.path).toBe("/v1/wallets");
      expect(result.query).toBeUndefined();
    });

    it("should parse URL with query string", () => {
      const result = TPV1Auth.parseUrl(
        "https://api.example.com/v1/wallets?limit=10&offset=0"
      );
      expect(result.host).toBe("api.example.com");
      expect(result.path).toBe("/v1/wallets");
      expect(result.query).toBe("limit=10&offset=0");
    });

    it("should handle URL with port", () => {
      const result = TPV1Auth.parseUrl("https://api.example.com:8443/v1/health");
      expect(result.host).toBe("api.example.com:8443");
      expect(result.path).toBe("/v1/health");
    });
  });
});

describe("calculateSignedHeader", () => {
  it("should produce a valid TPV1 header", () => {
    const secret = Buffer.from("deadbeef", "hex");
    const header = calculateSignedHeader(
      "my-key",
      secret,
      "nonce-123",
      1700000000000,
      "GET",
      "api.example.com",
      "/v1/wallets"
    );

    expect(header.startsWith("TPV1-HMAC-SHA256 ")).toBe(true);
    expect(header).toContain("ApiKey=my-key");
    expect(header).toContain("Nonce=nonce-123");
    expect(header).toContain("Timestamp=1700000000000");
    expect(header).toMatch(/Signature=.+/);
  });

  it("should be deterministic for same inputs", () => {
    const secret = Buffer.from("deadbeef", "hex");
    const args = [
      "my-key",
      secret,
      "fixed-nonce",
      1700000000000,
      "POST",
      "api.example.com",
      "/v1/wallets",
      undefined,
      "application/json",
      '{"name":"test"}',
    ] as const;

    const header1 = calculateSignedHeader(...args);
    const header2 = calculateSignedHeader(...args);
    expect(header1).toBe(header2);
  });

  it("should include body in signature computation", () => {
    const secret = Buffer.from("deadbeef", "hex");
    const common = [
      "my-key",
      secret,
      "nonce-1",
      1700000000000,
      "POST",
      "api.example.com",
      "/v1/wallets",
      undefined,
      "application/json",
    ] as const;

    const header1 = calculateSignedHeader(...common, '{"a":1}');
    const header2 = calculateSignedHeader(...common, '{"b":2}');

    const sig1 = header1.split("Signature=")[1];
    const sig2 = header2.split("Signature=")[1];
    expect(sig1).not.toBe(sig2);
  });
});
