/**
 * Cross-SDK enforced-in-rules flags.
 *
 * validatord computes enforcedInRules only on GetMe (when asked), GetUsers and GetGroups, and
 * omits a false bool from its JSON. Every SDK therefore reads a present flag as sent, an absent
 * one as false where the endpoint computes it and as absent where it does not, and never
 * defaults publicKeyEnforcedInRules. The shared vector file is consumed by all four suites.
 */
import * as fs from 'fs';
import * as path from 'path';
import type { Group, User } from '../../src/models/user';
import { pagedServices, RecordingTransport } from './pagination/harness';

const VECTORS_PATH = path.resolve(
  __dirname,
  '../../../scripts/resources/enforced-in-rules-vectors.json'
);
const VECTOR_COUNT = 7;

interface UserFlags {
  enforcedInRules: boolean | null;
  publicKeyEnforcedInRules: boolean | null;
  groupsEnforcedInRules: Array<boolean | null>;
}

interface GroupFlags {
  enforcedInRules: boolean | null;
  usersEnforcedInRules: Array<boolean | null>;
}

interface Vector {
  name: string;
  operation: string;
  request: { id?: string };
  reply: unknown;
  expected: { query?: Record<string, string>; users?: UserFlags[]; groups?: GroupFlags[] };
}

// Read inside each test: a throw at module or describe level drops the suite instead of failing it.
function loadVectors(): Vector[] {
  return JSON.parse(fs.readFileSync(VECTORS_PATH, 'utf-8')) as Vector[];
}

// The vectors spell an absent value as null; here it is undefined.
function orNull(value: boolean | undefined): boolean | null {
  return value ?? null;
}

function userFlags(user: User): UserFlags {
  return {
    enforcedInRules: orNull(user.enforcedInRules),
    publicKeyEnforcedInRules: orNull(user.publicKeyEnforcedInRules),
    groupsEnforcedInRules: (user.groups ?? []).map((group) => orNull(group.enforcedInRules)),
  };
}

function groupFlags(group: Group): GroupFlags {
  return {
    enforcedInRules: orNull(group.enforcedInRules),
    usersEnforcedInRules: (group.users ?? []).map((user) => orNull(user.enforcedInRules)),
  };
}

const OPERATIONS: Record<
  string,
  (transport: RecordingTransport, vector: Vector) => Promise<{ users?: User[]; groups?: Group[] }>
> = {
  getMe: async (transport) => ({ users: [await pagedServices(transport).users.getCurrentUser()] }),
  getUser: async (transport, vector) => ({
    users: [await pagedServices(transport).users.get(vector.request.id as string)],
  }),
  listUsers: async (transport) => ({ users: (await pagedServices(transport).users.list()).items }),
  listGroups: async (transport) => ({ groups: (await pagedServices(transport).groups.list()).items }),
};

describe('enforced-in-rules vectors', () => {
  it('loads every shared vector and knows each operation', () => {
    const vectors = loadVectors();

    expect(vectors).toHaveLength(VECTOR_COUNT);
    for (const vector of vectors) {
      expect(`${vector.name}: ${vector.operation in OPERATIONS ? 'consumed' : 'NOT consumed'}`).toBe(
        `${vector.name}: consumed`
      );
    }
  });

  it('reads the flags each endpoint computes and nothing else', async () => {
    for (const vector of loadVectors()) {
      const transport = new RecordingTransport(vector.reply);

      const result = await OPERATIONS[vector.operation](transport, vector);

      const sent = transport.single().query;
      for (const [key, value] of Object.entries(vector.expected.query ?? {})) {
        expect([vector.name, sent.filter(([k]) => k === key)]).toEqual([vector.name, [[key, value]]]);
      }
      expect(vector.expected.users ?? vector.expected.groups).toBeDefined();
      if (vector.expected.users) {
        expect([vector.name, result.users?.map(userFlags)]).toEqual([
          vector.name,
          vector.expected.users,
        ]);
      }
      if (vector.expected.groups) {
        expect([vector.name, result.groups?.map(groupFlags)]).toEqual([
          vector.name,
          vector.expected.groups,
        ]);
      }
    }
  });
});
