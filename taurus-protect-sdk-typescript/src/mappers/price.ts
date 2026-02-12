/**
 * Price mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { ConversionResult, Price, PriceHistoryPoint } from '../models/price';
import { currencyFromDto } from './currency';
import { safeDate, safeMap, safeString } from './base';

/**
 * Maps a price DTO to a Price domain model.
 */
export function priceFromDto(dto: unknown): Price | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    blockchain: safeString(d.blockchain),
    currencyFrom: safeString(d.currencyFrom ?? d.currency_from),
    currencyTo: safeString(d.currencyTo ?? d.currency_to),
    decimals: safeString(d.decimals),
    rate: safeString(d.rate),
    changePercent24Hour: safeString(d.changePercent24Hour ?? d.change_percent_24_hour),
    source: safeString(d.source),
    // Map creationDate -> createdAt, updateDate -> updatedAt (as in Java mapper)
    createdAt: safeDate(d.creationDate ?? d.creation_date ?? d.createdAt ?? d.created_at),
    updatedAt: safeDate(d.updateDate ?? d.update_date ?? d.updatedAt ?? d.updated_at),
    currencyFromInfo: currencyFromDto(d.currencyFromInfo ?? d.currency_from_info),
    currencyToInfo: currencyFromDto(d.currencyToInfo ?? d.currency_to_info),
  };
}

/**
 * Maps an array of price DTOs to Price domain models.
 */
export function pricesFromDto(dtos: unknown[] | null | undefined): Price[] {
  return safeMap(dtos, priceFromDto);
}

/**
 * Maps a price history point DTO to a PriceHistoryPoint domain model.
 */
export function priceHistoryPointFromDto(dto: unknown): PriceHistoryPoint | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    periodStartDate: safeDate(d.periodStartDate ?? d.period_start_date),
    blockchain: safeString(d.blockchain),
    currencyFrom: safeString(d.currencyFrom ?? d.currency_from),
    currencyTo: safeString(d.currencyTo ?? d.currency_to),
    high: safeString(d.high),
    low: safeString(d.low),
    open: safeString(d.open),
    close: safeString(d.close),
    volumeFrom: safeString(d.volumeFrom ?? d.volume_from),
    volumeTo: safeString(d.volumeTo ?? d.volume_to),
    changePercent: safeString(d.changePercent ?? d.change_percent),
    currencyFromInfo: currencyFromDto(d.currencyFromInfo ?? d.currency_from_info),
    currencyToInfo: currencyFromDto(d.currencyToInfo ?? d.currency_to_info),
  };
}

/**
 * Maps an array of price history point DTOs to PriceHistoryPoint domain models.
 */
export function priceHistoryPointsFromDto(dtos: unknown[] | null | undefined): PriceHistoryPoint[] {
  return safeMap(dtos, priceHistoryPointFromDto);
}

/**
 * Maps a conversion value DTO to a ConversionResult domain model.
 */
export function conversionResultFromDto(dto: unknown): ConversionResult | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    symbol: safeString(d.symbol),
    value: safeString(d.value),
    mainUnitValue: safeString(d.mainUnitValue ?? d.main_unit_value),
    currencyInfo: currencyFromDto(d.currencyInfo ?? d.currency_info),
  };
}

/**
 * Maps an array of conversion value DTOs to ConversionResult domain models.
 */
export function conversionResultsFromDto(dtos: unknown[] | null | undefined): ConversionResult[] {
  return safeMap(dtos, conversionResultFromDto);
}
