/**
 * Unit tests for constant-time comparison utilities.
 *
 * These tests verify that the constant-time comparison functions
 * behave correctly for security-critical comparisons.
 */

import {
  constantTimeCompare,
  constantTimeCompareBytes,
} from '../../../src/helpers/constant-time';

describe('constantTimeCompare', () => {
  it('should return true for equal strings', () => {
    expect(constantTimeCompare('hello', 'hello')).toBe(true);
    expect(constantTimeCompare('abc123', 'abc123')).toBe(true);
  });

  it('should return false for different strings', () => {
    expect(constantTimeCompare('hello', 'world')).toBe(false);
    expect(constantTimeCompare('abc', 'xyz')).toBe(false);
  });

  it('should return false for strings of different length', () => {
    expect(constantTimeCompare('short', 'longer string')).toBe(false);
    expect(constantTimeCompare('a', 'ab')).toBe(false);
  });

  it('should return true for empty strings', () => {
    expect(constantTimeCompare('', '')).toBe(true);
  });

  it('should return false when one string is empty', () => {
    expect(constantTimeCompare('', 'notempty')).toBe(false);
    expect(constantTimeCompare('notempty', '')).toBe(false);
  });

  it('should handle hex hash strings correctly', () => {
    const hash1 = 'a'.repeat(64);
    const hash2 = 'a'.repeat(64);
    const hash3 = 'b'.repeat(64);

    expect(constantTimeCompare(hash1, hash2)).toBe(true);
    expect(constantTimeCompare(hash1, hash3)).toBe(false);
  });

  it('should be case-sensitive', () => {
    expect(constantTimeCompare('Hello', 'hello')).toBe(false);
    expect(constantTimeCompare('ABC', 'abc')).toBe(false);
  });
});

describe('constantTimeCompareBytes', () => {
  it('should return true for equal byte arrays', () => {
    const a = new Uint8Array([1, 2, 3, 4]);
    const b = new Uint8Array([1, 2, 3, 4]);
    expect(constantTimeCompareBytes(a, b)).toBe(true);
  });

  it('should return false for different byte arrays', () => {
    const a = new Uint8Array([1, 2, 3, 4]);
    const b = new Uint8Array([5, 6, 7, 8]);
    expect(constantTimeCompareBytes(a, b)).toBe(false);
  });

  it('should return false for different length byte arrays', () => {
    const a = new Uint8Array([1, 2, 3]);
    const b = new Uint8Array([1, 2, 3, 4]);
    expect(constantTimeCompareBytes(a, b)).toBe(false);
  });

  it('should return true for empty byte arrays', () => {
    const a = new Uint8Array([]);
    const b = new Uint8Array([]);
    expect(constantTimeCompareBytes(a, b)).toBe(true);
  });
});
