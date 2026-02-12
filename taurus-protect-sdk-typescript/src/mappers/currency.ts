/**
 * Currency mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { Currency } from '../models/currency';
import { safeBool, safeBoolDefault, safeInt, safeMap, safeString } from './base';

/**
 * Maps a currency DTO to a Currency domain model.
 */
export function currencyFromDto(dto: unknown): Currency | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id ?? d.currencyId ?? d.currency_id),
    symbol: safeString(d.symbol),
    name: safeString(d.name),
    displayName: safeString(d.displayName ?? d.display_name),
    type: safeString(d.type),
    blockchain: safeString(d.blockchain),
    network: safeString(d.network),
    decimals: safeInt(d.decimals ?? d.decimal),
    coinTypeIndex: safeString(d.coinTypeIndex ?? d.coin_type_index),
    contractAddress: safeString(d.contractAddress ?? d.contract_address ?? d.tokenContractAddress),
    tokenId: safeString(d.tokenId ?? d.token_id),
    wlcaId: safeInt(d.wlcaId ?? d.wlca_id),
    logo: safeString(d.logo),
    isToken: safeBoolDefault(d.isToken ?? d.is_token, false),
    isERC20: safeBoolDefault(d.isERC20 ?? d.is_erc20, false),
    isFA12: safeBoolDefault(d.isFA12 ?? d.is_fa12, false),
    isFA20: safeBoolDefault(d.isFA20 ?? d.is_fa20, false),
    isNFT: safeBoolDefault(d.isNFT ?? d.is_nft, false),
    isUTXOBased: safeBoolDefault(d.isUTXOBased ?? d.is_utxo_based, false),
    isAccountBased: safeBoolDefault(d.isAccountBased ?? d.is_account_based, false),
    isFiat: safeBoolDefault(d.isFiat ?? d.is_fiat, false),
    hasStaking: safeBoolDefault(d.hasStaking ?? d.has_staking, false),
    enabled: safeBoolDefault(d.enabled, true),
    isNative: safeBool(d.isNative ?? d.is_native ?? d.native),
    isDisabled: safeBool(d.isDisabled ?? d.is_disabled ?? d.disabled),
    logoUrl: safeString(d.logoUrl ?? d.logo_url ?? d.logo),
    price: safeString(d.price),
    priceCurrency: safeString(d.priceCurrency ?? d.price_currency ?? d.baseCurrency),
  };
}

/**
 * Maps an array of currency DTOs to Currency domain models.
 */
export function currenciesFromDto(dtos: unknown[] | null | undefined): Currency[] {
  return safeMap(dtos, currencyFromDto);
}
