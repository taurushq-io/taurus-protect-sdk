/**
 * Bounds on the governance decoders: decode cost must be LINEAR in the payload, and
 * neither a cell payload nor a whole rules container may be unbounded.
 *
 * Both properties are reachable by anyone who can shape an API response. The rules
 * container is the document every HSM and PRICEUPDATER key is read from, so its decode
 * runs on the address, asset and price paths — a decoder that is quadratic in the
 * payload turns one response into a CPU stall for the whole process, and in a daemon
 * running one client per tenant that is every tenant at once.
 *
 * `magnitudeToBigint` is module-private, so it is exercised through the public codec
 * API, as this SDK's CLAUDE.md requires.
 */

import {
  cellFamily,
  ruleCellFromBytes,
  ruleCellToBytes,
} from "../../../src/mappers/rule-cell-codec";
import {
  MAX_CELL_PAYLOAD_BYTES,
  MAX_RULES_CONTAINER_BYTES,
} from "../../../src/helpers/signed-payload-guard";
import { rulesContainerFromBase64 } from "../../../src/mappers/governance-rules";
import { rulesContainerToBase64 } from "../../../src/mappers/protobuf-rules-container-encode";
import { createEmptyRulesContainer } from "../../../src/models/governance-rules";
import { tryDecodeProtobufRulesContainer } from "../../../src/mappers/protobuf-rules-container";
import { IntegrityError } from "../../../src/errors";
import type { RuleCell } from "../../../src/models/rule-cell";

describe("rule cell decode is linear in the payload size", () => {
  it("decodes a large but legal unsigned-integer magnitude well inside a time budget", () => {
    // 512 KiB of magnitude: legal (under the per-cell cap) and big enough that a
    // quadratic decode is unmistakable. Measured on this container: the shift-per-byte
    // loop takes ~23 s here, the one-pass construction ~30 ms.
    const width = 512 * 1024;
    expect(width).toBeLessThan(MAX_CELL_PAYLOAD_BYTES);

    const magnitude = new Uint8Array(width);
    magnitude[0] = 0x01; // no leading zero, so the magnitude round-trips exactly
    for (let i = 1; i < width; i++) {
      magnitude[i] = (i * 7 + 1) & 0xff;
    }
    const expected = BigInt(
      "0x" + Buffer.from(magnitude).toString("hex")
    );

    const cell: RuleCell = { kind: "UIntegerGreaterValue", value: expected };
    const bytes = ruleCellToBytes("RuleUIntegerGreater", cell);

    const started = Date.now();
    const decoded = ruleCellFromBytes("RuleUIntegerGreater", bytes);
    const elapsedMs = Date.now() - started;

    // The value must still be exact — a faster wrong answer is not the goal.
    expect(decoded.kind).toBe("UIntegerGreaterValue");
    expect((decoded as { value: bigint }).value).toBe(expected);

    // Generous: ~100x the linear cost, so a loaded machine cannot make this flake,
    // while the quadratic implementation is ~8x over it.
    expect([`elapsedMs<=3000`, elapsedMs <= 3000]).toEqual([
      `elapsedMs<=3000`,
      true,
    ]);
  }, 120000);
});

describe("cell payload size cap", () => {
  /**
   * A GENUINELY DECODABLE string cell whose encoded length is exactly `total`.
   *
   * The fixture has to be one the codec would otherwise decode to a typed cell,
   * or the test passes with or without the cap: raw filler bytes are not a valid
   * `RuleStringEqual` wrapper and already degrade to a `RawCell`. That is the
   * "a test can assert the vulnerable behaviour" trap in this repo's CLAUDE.md.
   *
   * A string cell is used rather than an integer one so this test stays fast — the
   * property under test is the cap, not the arithmetic.
   */
  function stringCellOfExactLength(total: number): { bytes: Uint8Array; value: string } {
    let n = total;
    let length = 0;
    for (let i = 0; i < 8; i++) {
      const value = "A".repeat(n);
      const bytes = ruleCellToBytes("RuleStringEqual", {
        kind: "StringEqualValue",
        value,
      });
      length = bytes.length;
      if (length === total) return { bytes, value };
      n += total - length;
    }
    throw new Error(
      `test setup: could not build a string cell of exactly ${total} bytes (got ${length})`
    );
  }

  it("decodes a string cell at exactly the cap", () => {
    // Off-by-one guard: the check is `>`, so exactly MAX_CELL_PAYLOAD_BYTES passes.
    // This is also what proves the over-cap case below is refused BY THE CAP and not
    // because the fixture was undecodable.
    const { bytes, value } = stringCellOfExactLength(MAX_CELL_PAYLOAD_BYTES);
    expect(bytes.length).toBe(MAX_CELL_PAYLOAD_BYTES);

    const decoded = ruleCellFromBytes("RuleStringEqual", bytes);

    expect(decoded.kind).toBe("StringEqualValue");
    expect((decoded as { value: string }).value).toBe(value);
  });

  it("refuses to decode one byte over the cap, preserving it verbatim as a RawCell", () => {
    const { bytes } = stringCellOfExactLength(MAX_CELL_PAYLOAD_BYTES + 1);
    expect(bytes.length).toBe(MAX_CELL_PAYLOAD_BYTES + 1);

    const decoded = ruleCellFromBytes("RuleStringEqual", bytes);

    // Would be 'StringEqualValue' without the cap — that is the red/green line.
    expect(decoded.kind).toBe("RawCell");
    // Lossless: an over-cap cell is not silently dropped, it is passed through.
    // Compared with a native memcmp — jest's `toEqual` deep-walks a 1 MiB array
    // element by element and adds ~11 s to the suite for no extra assurance.
    const payload = (decoded as { payload: Uint8Array }).payload;
    expect(payload.length).toBe(bytes.length);
    expect(Buffer.compare(Buffer.from(payload), Buffer.from(bytes))).toBe(0);
    // A RawCell still reports the column family it came from.
    expect(cellFamily(decoded)).toBe("RuleStringEqual");
  });
});

describe("rules container size cap", () => {
  /**
   * A WIRE-VALID container that is merely too large.
   *
   * This matters: garbage bytes already throw `IntegrityError` from the
   * "not valid protobuf or JSON" branch, so a malformed fixture would make this test
   * pass whether the cap exists or not.
   */
  function oversizedButValidContainerBase64(): string {
    const container = createEmptyRulesContainer();
    return rulesContainerToBase64({
      ...container,
      users: [
        {
          id: "1",
          name: undefined,
          publicKeyPem: "A".repeat(MAX_RULES_CONTAINER_BYTES + 1024),
          roles: [],
          properties: {},
          unknownFields: undefined,
        },
      ],
    });
  }

  it("refuses an over-cap container, and says so", () => {
    const base64 = oversizedButValidContainerBase64();

    expect(() => rulesContainerFromBase64(base64)).toThrow(IntegrityError);
    // Anchored on the size message so this cannot pass via the generic
    // "not valid protobuf or JSON" branch.
    expect(() => rulesContainerFromBase64(base64)).toThrow(/exceeds/);
  });

  it("proves the oversized fixture is otherwise wire-valid", () => {
    // Without this guard the test above could be green against a decoder with no cap
    // at all — the trap this repo has hit twice.
    const container = createEmptyRulesContainer();
    const small = rulesContainerToBase64({
      ...container,
      users: [
        {
          id: "1",
          name: undefined,
          publicKeyPem: "A".repeat(1024),
          roles: [],
          properties: {},
          unknownFields: undefined,
        },
      ],
    });
    const decoded = rulesContainerFromBase64(small);
    expect(decoded.users).toHaveLength(1);
    expect(decoded.users[0].publicKeyPem).toBe("A".repeat(1024));
  });

  it("refuses an over-cap container at the protobuf entry point too", () => {
    // tryDecodeProtobufRulesContainer is exported, so it is reachable without going
    // through the base64 wrapper's cap.
    const oversized = new Uint8Array(MAX_RULES_CONTAINER_BYTES + 1);
    expect(tryDecodeProtobufRulesContainer(oversized)).toBeUndefined();
  });

  it("still decodes an empty container", () => {
    const decoded = rulesContainerFromBase64(
      rulesContainerToBase64(createEmptyRulesContainer())
    );
    expect(decoded.users).toEqual([]);
  });
});
