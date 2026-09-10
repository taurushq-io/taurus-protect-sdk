/**
 * Loader for the cross-SDK lossless/parity vectors.
 *
 * The scenarios in governance-lossless-vectors.json cannot be produced by this SDK's
 * encoder — they are deliberately non-canonical or schema-newer wire bytes — so they
 * are not in governance-cell-vectors.json, which the Go SDK generates. They used to
 * live as base64 literals hand-copied into all four SDKs' round-trip suites, where
 * nothing compared the copies and they could drift silently.
 */
import * as fs from 'fs';
import * as path from 'path';

const VECTORS_PATH = path.resolve(
  __dirname,
  '../../../../scripts/resources/governance-lossless-vectors.json'
);

/**
 * Asserted so a vector added to the shared file without being consumed here fails
 * loudly instead of being silently ignored by this SDK.
 */
const VECTOR_COUNT = 9;

export interface LosslessVector {
  readonly description: string;
  readonly scenario: string;
  readonly wire_base64: string;
}

export function loadLosslessVectors(): LosslessVector[] {
  // readFileSync throws when the file is missing, which is the point: a missing gate
  // must fail, never silently pass.
  const vectors = JSON.parse(
    fs.readFileSync(VECTORS_PATH, 'utf-8')
  ) as LosslessVector[];

  if (vectors.length !== VECTOR_COUNT) {
    throw new Error(
      `shared lossless vectors file has ${vectors.length} entries, expected ${VECTOR_COUNT} — update every SDK's suite in lockstep`
    );
  }
  return vectors;
}

export function vectorsFor(scenario: string): LosslessVector[] {
  const out = loadLosslessVectors().filter((v) => v.scenario === scenario);
  if (out.length === 0) {
    throw new Error(`no shared lossless vector for scenario "${scenario}"`);
  }
  return out;
}

export function vectorFor(scenario: string): string {
  const out = vectorsFor(scenario);
  if (out.length !== 1) {
    throw new Error(
      `scenario "${scenario}" has ${out.length} vectors, expected exactly 1`
    );
  }
  return out[0]!.wire_base64;
}
