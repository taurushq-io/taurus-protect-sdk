/**
 * Cross-SDK decode tolerance.
 *
 * Every SDK must behave the same way when the server sends a field or an enum value the
 * client does not know, and when a caller passes an enum value the SDK does not know: the
 * field is kept in `additionalProperties` and written back, the enum value stays the raw
 * string, and a caller's value is sent verbatim. The shared vector file is consumed by the
 * Go/Java/Python/TypeScript suites alike, so one SDK drifting is caught for all four.
 */
import * as fs from 'fs';
import * as path from 'path';
import { MultiFactorSignatureApi } from '../../src/internal/openapi/apis';
import {
  TgvalidatordGetBalancesReplyFromJSON,
  TgvalidatordGetBalancesReplyToJSON,
} from '../../src/internal/openapi/models/TgvalidatordGetBalancesReply';
import {
  TgvalidatordGetMultiFactorSignatureEntitiesInfoReplyFromJSON,
  TgvalidatordGetMultiFactorSignatureEntitiesInfoReplyToJSON,
} from '../../src/internal/openapi/models/TgvalidatordGetMultiFactorSignatureEntitiesInfoReply';
import {
  instanceOfTgvalidatordMultiFactorSignaturesEntityType,
  TgvalidatordMultiFactorSignaturesEntityTypeFromJSON,
} from '../../src/internal/openapi/models/TgvalidatordMultiFactorSignaturesEntityType';
import {
  TgvalidatordStakeAccountFromJSON,
  TgvalidatordStakeAccountToJSON,
} from '../../src/internal/openapi/models/TgvalidatordStakeAccount';
import {
  instanceOfTgvalidatordStakeAccountType,
  TgvalidatordStakeAccountTypeFromJSON,
} from '../../src/internal/openapi/models/TgvalidatordStakeAccountType';
import {
  instanceOfTgvalidatordTokenType,
  TgvalidatordTokenTypeFromJSON,
} from '../../src/internal/openapi/models/TgvalidatordTokenType';
import type { AssetAddressTypeV2 } from '../../src/models/asset';
import type { MultiFactorSignatureEntityType } from '../../src/models/multi-factor-signature';
import { MultiFactorSignatureService } from '../../src/services/multi-factor-signature-service';
import { pagedServices, RecordingTransport } from './pagination/harness';

const VECTORS_PATH = path.resolve(
  __dirname,
  '../../../scripts/resources/decode-tolerance-vectors.json'
);

interface ModelVector {
  name: string;
  kind: 'model';
  model: string;
  wire: Record<string, unknown>;
  expected: { lossless: boolean; rootAdditionalProperties: Record<string, unknown> };
}

interface EnumVector {
  name: string;
  kind: 'enum';
  enum: string;
  wire: string;
  expected: { value: string; known: boolean };
}

interface ServiceVector {
  name: string;
  kind: 'service';
  operation: string;
  request: Record<string, unknown>;
  reply: unknown;
  expected: Record<string, unknown>;
}

type Vector = ModelVector | EnumVector | ServiceVector;

const VECTORS: Vector[] = JSON.parse(fs.readFileSync(VECTORS_PATH, 'utf-8')) as Vector[];

const MODELS: Record<string, { fromJSON: (json: any) => any; toJSON: (value: any) => any }> = {
  tgvalidatordGetBalancesReply: {
    fromJSON: TgvalidatordGetBalancesReplyFromJSON,
    toJSON: TgvalidatordGetBalancesReplyToJSON,
  },
  tgvalidatordStakeAccount: {
    fromJSON: TgvalidatordStakeAccountFromJSON,
    toJSON: TgvalidatordStakeAccountToJSON,
  },
  tgvalidatordGetMultiFactorSignatureEntitiesInfoReply: {
    fromJSON: TgvalidatordGetMultiFactorSignatureEntitiesInfoReplyFromJSON,
    toJSON: TgvalidatordGetMultiFactorSignatureEntitiesInfoReplyToJSON,
  },
};

const ENUMS: Record<string, { fromJSON: (json: any) => unknown; isKnown: (value: any) => boolean }> =
  {
    tgvalidatordTokenType: {
      fromJSON: TgvalidatordTokenTypeFromJSON,
      isKnown: instanceOfTgvalidatordTokenType,
    },
    tgvalidatordMultiFactorSignaturesEntityType: {
      fromJSON: TgvalidatordMultiFactorSignaturesEntityTypeFromJSON,
      isKnown: instanceOfTgvalidatordMultiFactorSignaturesEntityType,
    },
    tgvalidatordStakeAccountType: {
      fromJSON: TgvalidatordStakeAccountTypeFromJSON,
      isKnown: instanceOfTgvalidatordStakeAccountType,
    },
  };

function mfsService(transport: RecordingTransport): MultiFactorSignatureService {
  return new MultiFactorSignatureService(new MultiFactorSignatureApi(transport.configuration));
}

const SERVICES: Record<string, (vector: ServiceVector) => Promise<void>> = {
  getMultiFactorSignatureInfo: async (vector) => {
    const transport = new RecordingTransport(vector.reply);

    const info = await mfsService(transport).get(vector.request.id as string);

    expect(info.entityType).toBe(vector.expected.entityType);
  },

  createMultiFactorSignatures: async (vector) => {
    const transport = new RecordingTransport(vector.reply);

    await mfsService(transport).create({
      entityType: vector.request.entityType as MultiFactorSignatureEntityType,
      entityIds: vector.request.entityIds as string[],
    });

    expect(transport.single().body).toEqual(expect.objectContaining(vector.expected.requestBody));
  },

  queryAssetAddresses: async (vector) => {
    const transport = new RecordingTransport(vector.reply);

    const page = await pagedServices(transport).assets.queryAssetAddresses(
      vector.request.assetId as string,
      { addressType: vector.request.addressType as AssetAddressTypeV2 }
    );

    expect(transport.single().body).toEqual(expect.objectContaining(vector.expected.requestBody));
    expect(
      page.items.map((row) => ({
        address: row.address,
        addressType: row.addressType,
        verified: row.verified,
      }))
    ).toEqual(vector.expected.rows);
    expect(page.excludedUnverified).toHaveLength(vector.expected.excludedCount as number);
  },
};

function byKind<K extends Vector['kind']>(kind: K): Array<Extract<Vector, { kind: K }>> {
  return VECTORS.filter((v): v is Extract<Vector, { kind: K }> => v.kind === kind);
}

describe('decode tolerance vectors', () => {
  it('loads the shared vectors file and knows every vector', () => {
    expect(VECTORS.length).toBeGreaterThan(0);
    for (const vector of VECTORS) {
      const known =
        vector.kind === 'model'
          ? vector.model in MODELS
          : vector.kind === 'enum'
            ? vector.enum in ENUMS
            : vector.operation in SERVICES;
      expect(`${vector.name}: ${known ? 'consumed' : 'NOT consumed'}`).toBe(
        `${vector.name}: consumed`
      );
    }
  });

  it.each(byKind('model').map((v) => [v.name, v] as const))(
    'model %s keeps unknown fields and re-serializes losslessly',
    (_name, vector) => {
      const codec = MODELS[vector.model];

      const decoded = codec.fromJSON(vector.wire);

      expect(decoded.additionalProperties).toEqual(vector.expected.rootAdditionalProperties);
      if (vector.expected.lossless) {
        expect(JSON.parse(JSON.stringify(codec.toJSON(decoded)))).toEqual(vector.wire);
      }
    }
  );

  it.each(byKind('enum').map((v) => [v.name, v] as const))(
    'enum %s keeps the raw value',
    (_name, vector) => {
      const codec = ENUMS[vector.enum];

      const decoded = codec.fromJSON(vector.wire);

      expect(decoded).toBe(vector.expected.value);
      expect(codec.isKnown(decoded)).toBe(vector.expected.known);
    }
  );

  it.each(byKind('service').map((v) => [v.name, v] as const))(
    'service %s',
    async (_name, vector) => {
      await SERVICES[vector.operation](vector);
    }
  );
});
