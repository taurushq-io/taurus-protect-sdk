/**
 * Contract whitelisting mapper tests.
 *
 * Only the attribute mappers remain. The envelope and model mappers these tests used to
 * cover built a WhitelistedContract from the DTO with no verification — and parsed the
 * signed payload as fact — so their assertions pinned the bypass behaving consistently.
 * Verified reads are covered by the whitelisted-asset suites.
 */

import {
  whitelistedContractAttributeFromDto,
  whitelistedContractAttributesFromDto,
} from '../../../src/mappers/contract-whitelist';

describe('whitelistedContractAttributeFromDto', () => {
  it('maps every field', () => {
    const attr = whitelistedContractAttributeFromDto({
      id: 'attr-1',
      key: 'risk',
      value: 'low',
      contentType: 'text/plain',
      type: 'custom',
      subtype: 'note',
      isfile: false,
    });

    expect(attr).toEqual({
      id: 'attr-1',
      key: 'risk',
      value: 'low',
      contentType: 'text/plain',
      type: 'custom',
      subType: 'note',
      isFile: false,
    });
  });

  it('accepts snake_case and camelCase alternatives', () => {
    const attr = whitelistedContractAttributeFromDto({
      id: '1',
      key: 'k',
      value: 'v',
      content_type: 'application/json',
      sub_type: 'meta',
      is_file: true,
    });

    expect(attr?.contentType).toBe('application/json');
    expect(attr?.subType).toBe('meta');
    expect(attr?.isFile).toBe(true);
  });

  it('returns undefined for a non-object', () => {
    expect(whitelistedContractAttributeFromDto(undefined)).toBeUndefined();
    expect(whitelistedContractAttributeFromDto(null)).toBeUndefined();
    expect(whitelistedContractAttributeFromDto('nope')).toBeUndefined();
  });
});

describe('whitelistedContractAttributesFromDto', () => {
  it('maps a list', () => {
    const attrs = whitelistedContractAttributesFromDto([
      { id: '1', key: 'a', value: '1' },
      { id: '2', key: 'b', value: '2' },
    ]);

    expect(attrs).toHaveLength(2);
    expect(attrs[0]?.key).toBe('a');
    expect(attrs[1]?.key).toBe('b');
  });

  it('returns an empty list for null or undefined', () => {
    expect(whitelistedContractAttributesFromDto(null)).toEqual([]);
    expect(whitelistedContractAttributesFromDto(undefined)).toEqual([]);
  });
});
