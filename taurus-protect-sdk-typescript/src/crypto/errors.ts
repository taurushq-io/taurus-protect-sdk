/**
 * Deciding which errors a verification loop may swallow.
 *
 * Four sites across this SDK carried the same copy-pasted predicate:
 *
 *	error.message.includes('signature') || error.message.includes('key') ||
 *	error.message.includes('Invalid')   || error.message.includes('decode') ||
 *	error.message.includes('ERR_OSSL')
 *
 * Two problems with that. `includes('key')` swallows ANY error whose message
 * happens to contain the word — a `TypeError: Cannot read properties of undefined
 * (reading 'key')` from a real bug was silently treated as "this signature did not
 * verify". And it is coupled to wording: a Node or OpenSSL upgrade that rephrases a
 * message changes behaviour with no test failure.
 *
 * Node attaches a stable `code` to crypto and OpenSSL errors, so that is what this
 * matches. Java catches `SignatureException | InvalidKeyException` and Python
 * catches `InvalidSignature, ValueError`; this is the TypeScript equivalent.
 */

/** Node error codes that mean "this input did not verify", not "the code is wrong". */
const VERIFICATION_ERROR_CODE_PREFIXES = ["ERR_OSSL", "ERR_CRYPTO"] as const;

const VERIFICATION_ERROR_CODES = new Set([
  "ERR_INVALID_ARG_TYPE",
  "ERR_INVALID_ARG_VALUE",
]);

/**
 * Reports whether an error means the supplied key or signature was unusable,
 * rather than that something in the SDK is broken.
 *
 * A caller may swallow these and move to the next key; anything else must be
 * rethrown, so a genuine defect surfaces instead of reading as a failed signature.
 *
 * @param error - the thrown value
 * @returns true when the error describes bad crypto input
 */
export function isCryptoVerificationError(error: unknown): boolean {
  if (!(error instanceof Error)) {
    return false;
  }

  const code = (error as NodeJS.ErrnoException).code;
  if (typeof code === "string") {
    if (VERIFICATION_ERROR_CODES.has(code)) {
      return true;
    }
    return VERIFICATION_ERROR_CODE_PREFIXES.some((prefix) =>
      code.startsWith(prefix)
    );
  }

  // Some paths throw plain Errors with no code — DER parsing and PEM decoding in
  // this SDK. Those messages are ours, written right here, so matching on them is
  // safe in a way that matching OpenSSL's wording is not. Anchored with
  // startsWith, never includes, so an unrelated error that happens to contain one
  // of these words is not swallowed.
  return (
    error.message.startsWith("Invalid DER signature") ||
    error.message.startsWith("DER signature truncated") ||
    error.message.startsWith("Failed to decode public key") ||
    error.message.startsWith("Failed to decode private key") ||
    error.message.startsWith("Only P-256")
  );
}
