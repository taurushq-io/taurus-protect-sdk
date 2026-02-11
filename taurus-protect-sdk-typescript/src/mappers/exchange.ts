/**
 * Exchange mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { Exchange, ExchangeCounterparty, ExchangeWithdrawalFee } from '../models/exchange';
import { currencyFromDto } from './currency';
import { safeBool, safeDate, safeMap, safeString } from './base';

/**
 * Maps an exchange DTO to an Exchange domain model.
 */
export function exchangeFromDto(dto: unknown): Exchange | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    exchange: safeString(d.exchange),
    account: safeString(d.account),
    currency: safeString(d.currency),
    type: safeString(d.type),
    totalBalance: safeString(d.totalBalance ?? d.total_balance),
    status: safeString(d.status),
    container: safeString(d.container),
    label: safeString(d.label),
    displayLabel: safeString(d.displayLabel ?? d.display_label),
    baseCurrencyValuation: safeString(d.baseCurrencyValuation ?? d.base_currency_valuation),
    hasWLA: safeBool(d.hasWLA ?? d.has_wla),
    currencyInfo: currencyFromDto(d.currencyInfo ?? d.currency_info),
    creationDate: safeDate(d.creationDate ?? d.creation_date),
    updateDate: safeDate(d.updateDate ?? d.update_date),
  };
}

/**
 * Maps an array of exchange DTOs to Exchange domain models.
 */
export function exchangesFromDto(dtos: unknown[] | null | undefined): Exchange[] {
  return safeMap(dtos, exchangeFromDto);
}

/**
 * Maps an exchange counterparty DTO to an ExchangeCounterparty domain model.
 */
export function exchangeCounterpartyFromDto(dto: unknown): ExchangeCounterparty | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    name: safeString(d.name),
    baseCurrencyValuation: safeString(d.baseCurrencyValuation ?? d.base_currency_valuation),
  };
}

/**
 * Maps an array of exchange counterparty DTOs to ExchangeCounterparty domain models.
 */
export function exchangeCounterpartiesFromDto(
  dtos: unknown[] | null | undefined
): ExchangeCounterparty[] {
  return safeMap(dtos, exchangeCounterpartyFromDto);
}

/**
 * Maps a withdrawal fee response to an ExchangeWithdrawalFee domain model.
 */
export function exchangeWithdrawalFeeFromDto(dto: unknown): ExchangeWithdrawalFee | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  // The API returns the fee as 'result' in the reply
  const fee = safeString(d.result ?? d.fee);

  if (fee === undefined) {
    return undefined;
  }

  return {
    fee,
  };
}
