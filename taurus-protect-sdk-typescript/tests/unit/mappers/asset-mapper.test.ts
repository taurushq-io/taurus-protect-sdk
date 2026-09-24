/**
 * Unit tests for the v2 asset registry mappers.
 *
 * The rows are built with the generated FromJSON functions, so the wire key renames the
 * generator applies (`configuration` -> `_configuration`, `import` -> `_import`) are
 * exercised exactly as a reply delivers them.
 */

import {
  assetAddressV2FromDto,
  assetAddressesV2FromDto,
  assetOperationV2FromDto,
  assetOperationsV2FromDto,
  assetV2FromDto,
  assetsV2FromDto,
} from '../../../src/mappers/asset';
import { TgvalidatordAssetOperationV2FromJSON } from '../../../src/internal/openapi/models/TgvalidatordAssetOperationV2';
import { TgvalidatordAssetResourceV2FromJSON } from '../../../src/internal/openapi/models/TgvalidatordAssetResourceV2';

describe('assetV2FromDto', () => {
  it('should map the registry row, including the Canton instrument configuration', () => {
    const dto = TgvalidatordAssetResourceV2FromJSON({
      id: 'asset-1',
      tenantID: '7',
      version: '3',
      createdAt: '2026-01-01T00:00:00Z',
      label: 'Token',
      assetType: 'CANTON_NATIVE_TOKEN',
      status: 'ASSET_VIEW_STATUS_V2_ACTIVE',
      blockchain: 'CANTON',
      network: 'mainnet',
      currencyID: 'c-1',
      name: 'Token',
      symbol: 'TKN',
      decimals: '10',
      contractAddress: 'addr',
      attributes: [{ key: 'k', value: 'v' }],
      blockchainAsset: {
        cantonNativeTokenAsset: {
          instrumentID: 'inst-1',
          configuration: { cid: 'cid-1', requireCredentials: true, paused: false, operator: 'op' },
        },
      },
    });

    expect(assetV2FromDto(dto)).toEqual({
      id: 'asset-1',
      tenantId: '7',
      version: '3',
      createdAt: new Date('2026-01-01T00:00:00Z'),
      updatedAt: undefined,
      label: 'Token',
      assetType: 'CANTON_NATIVE_TOKEN',
      status: 'ASSET_VIEW_STATUS_V2_ACTIVE',
      blockchain: 'CANTON',
      network: 'mainnet',
      currencyId: 'c-1',
      name: 'Token',
      symbol: 'TKN',
      decimals: '10',
      contractAddress: 'addr',
      attributes: [{ key: 'k', value: 'v' }],
      cantonInstrumentId: 'inst-1',
      cantonConfiguration: {
        cid: 'cid-1',
        requireCredentials: true,
        paused: false,
        operator: 'op',
      },
    });
  });

  it('should return undefined for null input and [] for an absent array', () => {
    expect(assetV2FromDto(null)).toBeUndefined();
    expect(assetsV2FromDto(undefined)).toEqual([]);
  });
});

describe('assetAddressV2FromDto', () => {
  it('should map a holder row as unverified', () => {
    expect(
      assetAddressV2FromDto({
        address: '0xabc',
        kycStatus: 'KYC_STATUS_V2_APPROVED',
        balance: '5',
        addressType: 'ADDRESS_TYPE_V2_INTERNAL',
        addressID: '42',
        whitelistedAddressID: undefined,
      })
    ).toEqual({
      address: '0xabc',
      kycStatus: 'KYC_STATUS_V2_APPROVED',
      balance: '5',
      addressType: 'ADDRESS_TYPE_V2_INTERNAL',
      addressId: '42',
      whitelistedAddressId: undefined,
      verified: false,
    });
    expect(assetAddressesV2FromDto(null)).toEqual([]);
  });
});

describe('assetOperationV2FromDto', () => {
  it('should map the operation and the details of its type', () => {
    const dto = TgvalidatordAssetOperationV2FromJSON({
      id: 'op-1',
      assetID: 'asset-1',
      type: 'ASSET_OPERATION_TYPE_V2_IMPORT',
      status: 'ASSET_OPERATION_STATUS_V2_PENDING',
      initiatedByAddressID: '9',
      blockingReason: 'ASSET_OPERATION_BLOCKING_REASON_V2_MANUAL_PAUSE',
      import: {
        blockchain: 'CANTON',
        network: 'mainnet',
        label: 'Imported',
        price: '1',
        decimals: '6',
        address: 'addr',
      },
    });

    const op = assetOperationV2FromDto(dto);
    expect(op).toMatchObject({
      id: 'op-1',
      assetId: 'asset-1',
      type: 'ASSET_OPERATION_TYPE_V2_IMPORT',
      status: 'ASSET_OPERATION_STATUS_V2_PENDING',
      initiatedByAddressId: '9',
      import: {
        blockchain: 'CANTON',
        network: 'mainnet',
        label: 'Imported',
        price: '1',
        decimals: '6',
        address: 'addr',
      },
    });
    expect(op?.mint).toBeUndefined();
  });

  it('should map mint, burn and account-level targets', () => {
    const op = assetOperationV2FromDto(
      TgvalidatordAssetOperationV2FromJSON({
        mint: { destination: { addressID: '1' }, amount: '10', nftMetadata: ['m'] },
        burn: { destination: { whitelistedAddressID: '2' }, amount: '3', nftTokenIDs: ['t'] },
        setKyc: { target: { addressID: '4' }, status: 'KYC_STATUS_V2_REVOKED' },
      })
    );
    expect(op?.mint).toEqual({
      destination: { addressId: '1', whitelistedAddressId: undefined },
      amount: '10',
      nftMetadata: ['m'],
    });
    expect(op?.burn).toEqual({
      destination: { addressId: undefined, whitelistedAddressId: '2' },
      amount: '3',
      nftTokenIds: ['t'],
    });
    expect(op?.setKyc).toEqual({
      target: { addressId: '4', whitelistedAddressId: undefined },
      status: 'KYC_STATUS_V2_REVOKED',
    });
    expect(assetOperationsV2FromDto(undefined)).toEqual([]);
  });
});
