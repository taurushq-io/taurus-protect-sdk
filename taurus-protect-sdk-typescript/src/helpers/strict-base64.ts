/**
 * Strict base64 decoding for security-critical paths.
 *
 * `Buffer.from(s, "base64")` is lenient: it silently DISCARDS characters outside the
 * base64 alphabet and absorbs the rest. That makes the decoding non-injective — two
 * different strings can decode to the same bytes, and, worse, a string carrying embedded
 * separators decodes to a *valid container with attacker-chosen bytes appended*. Protobuf
 * treats concatenation as merge, so appended bytes ADD entries to repeated fields such as
 * `users`, which is how an attacker-controlled HSM or PRICEUPDATER key would reach the
 * governance trust root.
 *
 * Every governance path that turns a base64 string into bytes must agree byte-for-byte on
 * what that string means: the verification memo key, the container decode, and the
 * signature verification that decides what a signature actually covers. Go uses strict
 * `base64.StdEncoding` throughout, and that is precisely why Go was not exploitable by
 * this route.
 *
 * Throws a plain `Error` and lets each caller wrap it in the appropriate SDK error type,
 * so this module stays a dependency-free leaf.
 */

/**
 * Well-formed base64: groups of four alphabet characters, optionally ending in one or two
 * `=` pad characters. ASCII whitespace is stripped before the test, since line-wrapped
 * base64 is legitimate and introduces no ambiguity about the decoded bytes.
 */
const BASE64 = /^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$/;

/**
 * Decodes base64, rejecting any character outside the base64 alphabet.
 *
 * @param data - Base64 text. The empty string decodes to an empty buffer.
 * @returns The decoded bytes.
 * @throws Error If the input is not well-formed base64.
 */
export function strictBase64Decode(data: string): Buffer {
  const compact = (data ?? "").replace(/[ \t\n\r\v\f]+/g, "");
  if (!BASE64.test(compact)) {
    throw new Error("invalid base64 encoding");
  }
  return Buffer.from(compact, "base64");
}
