/**
 * Cross-SDK role-extraction alignment.
 *
 * Every SDK must derive the same roles from the same server message, so the shared
 * vector file is consumed by the Go/Java/Python/TypeScript suites alike and a wording
 * change is caught in one place for all four.
 */
import * as fs from 'fs';
import * as path from 'path';
import { AuthorizationError, mapHttpError, parseRequiredRoles } from '../../src/errors';

const VECTORS_PATH = path.resolve(
  __dirname,
  '../../../scripts/resources/authorization-error-vectors.json'
);

interface AuthzVector {
  description: string;
  message: string;
  expected_roles: string[];
}

const VECTORS: AuthzVector[] = JSON.parse(fs.readFileSync(VECTORS_PATH, 'utf-8')) as AuthzVector[];

describe('authorization error vectors', () => {
  it('loads the shared vectors file', () => {
    expect(VECTORS.length).toBeGreaterThan(0);
  });

  it.each(VECTORS.map((v) => [v.description, v] as const))(
    'parses roles: %s',
    (_description, vector) => {
      expect(parseRequiredRoles(vector.message)).toEqual(vector.expected_roles);
    }
  );
});

describe('AuthorizationError.requiredRoles', () => {
  it('is populated from the message', () => {
    const err = new AuthorizationError("one of the 'admin - adminreadonly' role is required");

    expect(err.requiredRoles).toEqual(['admin', 'adminreadonly']);
  });

  it('is empty for a non-role denial', () => {
    expect(new AuthorizationError('This endpoint has been disabled').requiredRoles).toEqual([]);
  });

  it('is populated through mapHttpError', () => {
    const err = mapHttpError(403, "one of the 'tpuser' role is required");

    expect(err).toBeInstanceOf(AuthorizationError);
    expect((err as AuthorizationError).requiredRoles).toEqual(['tpuser']);
  });

  it('tolerates an undefined message', () => {
    expect(parseRequiredRoles(undefined)).toEqual([]);
  });
});
