/**
 * The two bounds every signed whitelist payload passes before it is parsed: a size
 * ceiling, and a refusal to accept the same member name twice in one object.
 *
 * They live together, and in a leaf module, for the reason Go's peers do: both parse
 * functions have to apply both, and a new caller that remembers one and forgets the
 * other is the failure mode. `strict-base64.ts` is the same arrangement — a
 * dependency-free leaf so both `helpers/` and `models/` can call it without an import
 * cycle.
 *
 * ## Why duplicate keys are a bypass, not a curiosity
 *
 * `JSON.parse` keeps the LAST of two duplicate keys and reports no error. On the
 * whitelist verification paths that is a verification bypass rather than a curiosity:
 * the legacy-hash tolerance strips a member the parser would still read, so a
 * response-controlling server can append `,"label":"Coinbase Prime custody"` immediately
 * before the closing brace of a genuinely signed payload. The strip recovers the signed
 * bytes exactly, every signature check passes, and the parse then returns the attacker's
 * value as verified. Parsing the MATCHED VARIANT (see `LegacyPayloadVariant`) closes that
 * shape; this closes the shapes the strip does not reach — and it is the reason the two
 * defences have to land together.
 *
 * This has to be a separate structural pass over the TEXT. A `JSON.parse` reviver sees
 * the object AFTER the duplicate has already been collapsed, so it cannot tell one
 * member from two.
 */

import { IntegrityError } from "../errors";

/**
 * Maximum size of a signed payload before it is parsed.
 *
 * The payload is hash-checked in step 1 but not AUTHENTICATED until step 5's signatures
 * verify, so anything parsed in between is still attacker-influenced. A generous ceiling
 * — real payloads are a few hundred bytes — that only stops a hostile response consuming
 * memory before verification can reject it.
 *
 * Re-exported from `whitelist-hash-helper.ts`, which is where callers have always
 * imported it from.
 */
export const MAX_PAYLOAD_BYTES = 1 << 20;

/**
 * Maximum size of an encoded governance rules container before it is decoded.
 *
 * `MAX_PAYLOAD_BYTES` bounded only the signed whitelist payload, so nothing bounded the
 * container itself — and the container is the document every HSM and `PRICEUPDATER`
 * public key is read from, decoded on the address, asset and price read paths. An
 * unbounded decode there is reachable by anyone able to shape an API response.
 *
 * Deliberately 4x the payload ceiling: a real container is kilobytes, so this cannot
 * reject a legitimate ruleset, and it is a memory/CPU bound rather than a schema rule.
 */
export const MAX_RULES_CONTAINER_BYTES = 4 << 20;

/**
 * Maximum size of a single rule-cell payload before the typed codec parses it.
 *
 * Needed on top of the container bound because `ruleCellFromBytes` is exported and so
 * reachable with arbitrary bytes. An over-cap cell degrades to a `RawCell` — it is
 * preserved verbatim, never dropped — which keeps the codec's "decode never fails"
 * contract while refusing to run the typed parse (notably the big-integer
 * construction) over an attacker-chosen length.
 *
 * The largest cell in the cross-SDK golden vectors is 29 bytes, so 1 MiB is not a
 * limit any real ruleset can reach.
 */
export const MAX_CELL_PAYLOAD_BYTES = 1 << 20;

/**
 * How deep the pre-pass will walk before refusing.
 *
 * The signed payload is three levels at most (object -> linkedInternalAddresses ->
 * object), so anything deeper is not a payload this SDK can be looking at. The scan is
 * iterative rather than recursive, so this is a sanity bound on server-chosen nesting
 * rather than a stack guard — but it fails closed either way.
 */
export const MAX_SIGNED_PAYLOAD_DEPTH = 32;

/**
 * Returns the index one past the closing quote of the string starting at `start`,
 * or -1 when the string is unterminated.
 */
function scanStringToken(text: string, start: number): number {
  let i = start + 1;
  while (i < text.length) {
    const c = text[i]!;
    if (c === "\\") {
      // Skip the escape and whatever it escapes. A `\uXXXX` escape's four hex digits
      // are ordinary characters to this scan, none of which can be a quote.
      i += 2;
      continue;
    }
    if (c === '"') {
      return i + 1;
    }
    i++;
  }
  return -1;
}

/**
 * Decodes a raw string token (quotes included) into the member name it denotes.
 *
 * Escapes are decoded because `"a"` and `"a"` name the SAME member, and a
 * comparison on the raw text would let the second slip past the duplicate check.
 */
function decodeMemberName(rawToken: string): string | undefined {
  try {
    const parsed: unknown = JSON.parse(rawToken);
    return typeof parsed === "string" ? parsed : undefined;
  } catch {
    // Not a decodable string: leave it to the JSON.parse of the whole document rather
    // than inventing a second error message for the same malformed input.
    return undefined;
  }
}

/**
 * Throws when any object in `json` carries the same member name twice.
 *
 * Siblings sharing a key are fine — `[{"label":"a"},{"label":"b"}]` is two objects, each
 * with one `label` — so the key set is per open object, not per document.
 *
 * Malformed input is NOT this function's business: it scans structurally and returns
 * quietly on anything it cannot tokenise, so the caller's `JSON.parse` produces the
 * single authoritative parse error.
 *
 * @param json - the JSON document, as text
 * @throws {@link IntegrityError} on a duplicate member name, or on nesting deeper than
 *   {@link MAX_SIGNED_PAYLOAD_DEPTH}
 */
export function rejectDuplicateObjectKeys(json: string): void {
  // One frame per open container: a Set of member names for an object, null for an
  // array (where there are no member names to collide).
  const frames: Array<Set<string> | null> = [];

  // Whether the next string token is a MEMBER NAME rather than a value. True only
  // directly after `{`, or after a `,` whose enclosing container is an object — which is
  // exactly what distinguishes a key from a value without a full recursive descent.
  let expectMemberName = false;

  let i = 0;
  while (i < json.length) {
    const ch = json[i]!;

    if (ch === " " || ch === "\t" || ch === "\n" || ch === "\r") {
      i++;
      continue;
    }

    if (ch === "{" || ch === "[") {
      frames.push(ch === "{" ? new Set<string>() : null);
      if (frames.length > MAX_SIGNED_PAYLOAD_DEPTH) {
        throw new IntegrityError(
          `signed payload nests deeper than ${MAX_SIGNED_PAYLOAD_DEPTH} levels`
        );
      }
      expectMemberName = ch === "{";
      i++;
      continue;
    }

    if (ch === "}" || ch === "]") {
      frames.pop();
      expectMemberName = false;
      i++;
      continue;
    }

    if (ch === ",") {
      const enclosing = frames[frames.length - 1];
      expectMemberName = enclosing !== undefined && enclosing !== null;
      i++;
      continue;
    }

    if (ch === ":") {
      expectMemberName = false;
      i++;
      continue;
    }

    if (ch === '"') {
      const end = scanStringToken(json, i);
      if (end < 0) {
        return; // Unterminated string: JSON.parse reports it.
      }
      if (expectMemberName) {
        const frame = frames[frames.length - 1];
        const name = decodeMemberName(json.slice(i, end));
        if (name === undefined) {
          return;
        }
        if (frame) {
          if (frame.has(name)) {
            throw new IntegrityError(
              `signed payload carries duplicate key "${name}"; refusing to choose ` +
                `between two values for one field`
            );
          }
          frame.add(name);
        }
        expectMemberName = false;
      }
      i = end;
      continue;
    }

    // A scalar character (number, `true`, `false`, `null`): nothing to track. Validity
    // is JSON.parse's job.
    i++;
  }
}

/**
 * Applies both bounds to a signed payload, in the order a parse must apply them: size
 * first (so an oversized document is never walked), then duplicate member names.
 *
 * Called from inside `parseWhitelistedAddressFromJson` and
 * `parseWhitelistedAssetFromJson` rather than at their call sites, so a new caller
 * cannot forget either guard.
 *
 * @param json - the signed payload, as text
 * @param entity - what is being parsed, for the error message
 *   (e.g. `"whitelisted address"`)
 * @throws {@link IntegrityError} when the payload is oversized, nests too deeply, or
 *   carries a duplicate member name
 */
export function guardSignedPayload(json: string, entity: string): void {
  if (json.length > MAX_PAYLOAD_BYTES) {
    throw new IntegrityError(
      `cannot parse ${entity}: payload exceeds ${MAX_PAYLOAD_BYTES} bytes`
    );
  }
  rejectDuplicateObjectKeys(json);
}
