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
 * Whether `code` is a base64 alphabet character: `A-Z`, `a-z`, `0-9`, `+`, `/`.
 */
function isBase64Char(code: number): boolean {
  return (
    (code >= 0x41 && code <= 0x5a) || // A-Z
    (code >= 0x61 && code <= 0x7a) || // a-z
    (code >= 0x30 && code <= 0x39) || // 0-9
    code === 0x2b || // +
    code === 0x2f // /
  );
}

/**
 * Whether `s` is well-formed base64: a multiple of four alphabet characters, of which the
 * last one or two may instead be `=` pad characters.
 *
 * **A linear scan, deliberately, and not the regex it replaces.** The obvious spelling —
 * `/^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$/` — recurses per
 * four-character group in V8, so a multi-megabyte input throws
 * `RangeError: Maximum call stack size exceeded` out of the decoder rather than returning
 * a decision. Measured: 1.6 MB passes in 27 ms, 5.6 MB throws. That fails *closed*, so it
 * is availability only rather than a verification bypass — but it is reachable from any
 * governance response, through the signature-verification decode and the memo key as well
 * as the container decode, and only the container decode is behind a size cap. Go
 * (`base64.StdEncoding`), Python (`validate=True`) and Java all validate without a
 * quantified-group regex; TypeScript was the lone outlier.
 *
 * The accepted set is **unchanged**, which matters: what a base64 string means has to
 * agree byte-for-byte across the four SDKs, and the strictness property is cross-SDK
 * gated. Length-multiple-of-four plus "every non-pad character is in the alphabet" is
 * exactly the language the regex described.
 */
function isWellFormedBase64(s: string): boolean {
  const len = s.length;
  if (len % 4 !== 0) {
    return false;
  }

  // At most two trailing '=', and only at the very end — a '=' anywhere else is rejected
  // by the alphabet scan below, so padding cannot appear mid-string.
  let pad = 0;
  if (len > 0 && s.charCodeAt(len - 1) === 0x3d) {
    pad = s.charCodeAt(len - 2) === 0x3d ? 2 : 1;
  }

  for (let i = 0; i < len - pad; i++) {
    if (!isBase64Char(s.charCodeAt(i))) {
      return false;
    }
  }
  return true;
}

/**
 * Decodes base64, rejecting any character outside the base64 alphabet.
 *
 * @param data - Base64 text. The empty string decodes to an empty buffer.
 * @returns The decoded bytes.
 * @throws Error If the input is not well-formed base64.
 */
export function strictBase64Decode(data: string): Buffer {
  const compact = (data ?? "").replace(/[ \t\n\r\v\f]+/g, "");
  if (!isWellFormedBase64(compact)) {
    throw new Error("invalid base64 encoding");
  }
  return Buffer.from(compact, "base64");
}
