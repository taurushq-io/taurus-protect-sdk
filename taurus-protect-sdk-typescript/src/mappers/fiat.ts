/**
 * Fiat provider mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { TgvalidatordFiatProviderEntity } from '../internal/openapi/models/TgvalidatordFiatProviderEntity';
import type {
  FiatProvider,
  FiatProviderAccount,
  FiatProviderCounterpartyAccount,
  FiatProviderEntity,
  FiatProviderOperation,
} from '../models/fiat';
import { currencyFromDto } from './currency';
import { safeDate, safeMap, safeString } from './base';

/**
 * Maps a fiat provider DTO to a FiatProvider domain model.
 */
export function fiatProviderFromDto(dto: unknown): FiatProvider | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    provider: safeString(d.provider),
    label: safeString(d.label),
    baseCurrencyValuation: safeString(d.baseCurrencyValuation),
  };
}

/**
 * Maps an array of fiat provider DTOs to FiatProvider domain models.
 */
export function fiatProvidersFromDto(
  dtos: unknown[] | null | undefined
): FiatProvider[] {
  return safeMap(dtos, fiatProviderFromDto);
}

/**
 * Maps a fiat provider account DTO to a FiatProviderAccount domain model.
 */
export function fiatProviderAccountFromDto(
  dto: unknown
): FiatProviderAccount | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    provider: safeString(d.provider),
    label: safeString(d.label),
    accountType: safeString(d.accountType),
    accountIdentifier: safeString(d.accountIdentifier),
    accountName: safeString(d.accountName),
    totalBalance: safeString(d.totalBalance),
    currencyId: safeString(d.currencyID ?? d.currencyId),
    currencyInfo: d.currencyInfo ? currencyFromDto(d.currencyInfo) : undefined,
    baseCurrencyValuation: safeString(d.baseCurrencyValuation),
    creationDate: safeDate(d.creationDate),
    updateDate: safeDate(d.updateDate),
  };
}

/**
 * Maps an array of fiat provider account DTOs to FiatProviderAccount domain models.
 */
export function fiatProviderAccountsFromDto(
  dtos: unknown[] | null | undefined
): FiatProviderAccount[] {
  return safeMap(dtos, fiatProviderAccountFromDto);
}

/**
 * Maps a fiat provider counterparty account DTO to a FiatProviderCounterpartyAccount domain model.
 */
export function fiatProviderCounterpartyAccountFromDto(
  dto: unknown
): FiatProviderCounterpartyAccount | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    provider: safeString(d.provider),
    label: safeString(d.label),
    accountType: safeString(d.accountType),
    accountIdentifier: safeString(d.accountIdentifier),
    accountName: safeString(d.accountName),
    counterpartyId: safeString(d.counterpartyID ?? d.counterpartyId),
    counterpartyName: safeString(d.counterpartyName),
    currencyId: safeString(d.currencyID ?? d.currencyId),
    currencyInfo: d.currencyInfo ? currencyFromDto(d.currencyInfo) : undefined,
    creationDate: safeDate(d.creationDate),
    updateDate: safeDate(d.updateDate),
  };
}

/**
 * Maps an array of fiat provider counterparty account DTOs to domain models.
 */
export function fiatProviderCounterpartyAccountsFromDto(
  dtos: unknown[] | null | undefined
): FiatProviderCounterpartyAccount[] {
  return safeMap(dtos, fiatProviderCounterpartyAccountFromDto);
}

/**
 * Maps a fiat provider operation DTO to a FiatProviderOperation domain model.
 */
export function fiatProviderOperationFromDto(
  dto: unknown
): FiatProviderOperation | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }

  const d = dto as Record<string, unknown>;
  return {
    id: safeString(d.id),
    provider: safeString(d.provider),
    label: safeString(d.label),
    operationType: safeString(d.operationType),
    operationIdentifier: safeString(d.operationIdentifier),
    operationDirection: safeString(d.operationDirection),
    status: safeString(d.status),
    amount: safeString(d.amount),
    currencyId: safeString(d.currencyID ?? d.currencyId),
    currencyInfo: d.currencyInfo ? currencyFromDto(d.currencyInfo) : undefined,
    fromAccountId: safeString(d.fromAccountID ?? d.fromAccountId),
    toAccountId: safeString(d.toAccountID ?? d.toAccountId),
    fromDetails: safeString(d.fromDetails),
    toDetails: safeString(d.toDetails),
    comment: safeString(d.comment),
    operationDetails: safeString(d.operationDetails),
    creationDate: safeDate(d.creationDate),
    updateDate: safeDate(d.updateDate),
  };
}

/**
 * Maps an array of fiat provider operation DTOs to domain models.
 */
export function fiatProviderOperationsFromDto(
  dtos: unknown[] | null | undefined
): FiatProviderOperation[] {
  return safeMap(dtos, fiatProviderOperationFromDto);
}

/**
 * Maps a fiat provider entity row to a FiatProviderEntity domain model.
 */
export function fiatProviderEntityFromDto(
  dto: TgvalidatordFiatProviderEntity | null | undefined
): FiatProviderEntity | undefined {
  if (!dto || typeof dto !== 'object') {
    return undefined;
  }
  return {
    id: dto.id,
    provider: dto.provider,
    label: dto.label,
    accountIdentifier: dto.accountIdentifier,
    name: dto.name,
    details: dto.details,
    creationDate: dto.creationDate,
    updateDate: dto.updateDate,
  };
}

/**
 * Maps an array of fiat provider entity rows to FiatProviderEntity domain models.
 */
export function fiatProviderEntitiesFromDto(
  dtos: TgvalidatordFiatProviderEntity[] | null | undefined
): FiatProviderEntity[] {
  return safeMap(dtos, fiatProviderEntityFromDto);
}
