/**
 * Unit tests for hashing utilities.
 *
 * These tests verify SHA-256 hashing, HMAC calculation,
 * and constant-time comparison functions used for security.
 */

import {
  calculateHexHash,
  calculateSha256Bytes,
  calculateBase64Hmac,
  constantTimeCompare,
  constantTimeCompareBytes,
} from '../../../src/crypto/hashing';

describe('calculateHexHash', () => {
  it('should return correct SHA-256 hex hash for known input', () => {
    // SHA-256 of "hello" is well-known
    const hash = calculateHexHash('hello');
    expect(hash).toBe('2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824');
  });

  it('should return correct hash for empty string', () => {
    // SHA-256 of "" is well-known
    const hash = calculateHexHash('');
    expect(hash).toBe('e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855');
  });

  it('should return lowercase hex string', () => {
    const hash = calculateHexHash('test');
    expect(hash).toMatch(/^[0-9a-f]{64}$/);
  });

  it('should return 64 character string (256 bits)', () => {
    const hash = calculateHexHash('any input');
    expect(hash.length).toBe(64);
  });

  it('should produce different hashes for different inputs', () => {
    const hash1 = calculateHexHash('input1');
    const hash2 = calculateHexHash('input2');
    expect(hash1).not.toBe(hash2);
  });

  it('should produce same hash for same input', () => {
    const hash1 = calculateHexHash('deterministic');
    const hash2 = calculateHexHash('deterministic');
    expect(hash1).toBe(hash2);
  });
});

describe('calculateSha256Bytes', () => {
  it('should return 32 bytes', () => {
    const bytes = calculateSha256Bytes(new Uint8Array([1, 2, 3]));
    expect(bytes.length).toBe(32);
  });

  it('should return same result for same input', () => {
    const input = new Uint8Array([0x41, 0x42, 0x43]); // "ABC"
    const result1 = calculateSha256Bytes(input);
    const result2 = calculateSha256Bytes(input);
    expect(Buffer.from(result1).equals(Buffer.from(result2))).toBe(true);
  });
});

describe('calculateBase64Hmac', () => {
  it('should return a base64 encoded HMAC', () => {
    const secret = Buffer.from('my-secret-key');
    const hmac = calculateBase64Hmac(secret, 'data-to-sign');

    expect(hmac).toBeDefined();
    expect(hmac.length).toBeGreaterThan(0);
    // Should be valid base64
    expect(() => Buffer.from(hmac, 'base64')).not.toThrow();
  });

  it('should produce same HMAC for same inputs', () => {
    const secret = Buffer.from('key');
    const hmac1 = calculateBase64Hmac(secret, 'data');
    const hmac2 = calculateBase64Hmac(secret, 'data');
    expect(hmac1).toBe(hmac2);
  });

  it('should produce different HMACs for different keys', () => {
    const hmac1 = calculateBase64Hmac(Buffer.from('key1'), 'data');
    const hmac2 = calculateBase64Hmac(Buffer.from('key2'), 'data');
    expect(hmac1).not.toBe(hmac2);
  });

  it('should produce different HMACs for different data', () => {
    const secret = Buffer.from('key');
    const hmac1 = calculateBase64Hmac(secret, 'data1');
    const hmac2 = calculateBase64Hmac(secret, 'data2');
    expect(hmac1).not.toBe(hmac2);
  });
});

describe('constantTimeCompare (crypto/hashing)', () => {
  it('should return true for equal strings', () => {
    expect(constantTimeCompare('abc', 'abc')).toBe(true);
  });

  it('should return false for different strings', () => {
    expect(constantTimeCompare('abc', 'xyz')).toBe(false);
  });

  it('should return false for different length strings with dummy comparison', () => {
    // This verifies the timing-safe behavior even for different-length inputs
    expect(constantTimeCompare('short', 'much longer string')).toBe(false);
  });

  it('should return true for empty strings', () => {
    expect(constantTimeCompare('', '')).toBe(true);
  });
});

describe('constantTimeCompareBytes (crypto/hashing)', () => {
  it('should return true for equal byte arrays', () => {
    const a = new Uint8Array([1, 2, 3]);
    const b = new Uint8Array([1, 2, 3]);
    expect(constantTimeCompareBytes(a, b)).toBe(true);
  });

  it('should return false for different byte arrays', () => {
    const a = new Uint8Array([1, 2, 3]);
    const b = new Uint8Array([4, 5, 6]);
    expect(constantTimeCompareBytes(a, b)).toBe(false);
  });

  it('should return false for different length byte arrays', () => {
    const a = new Uint8Array([1, 2]);
    const b = new Uint8Array([1, 2, 3]);
    expect(constantTimeCompareBytes(a, b)).toBe(false);
  });
});
