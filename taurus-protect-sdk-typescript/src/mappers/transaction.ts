/**
 * Transaction mapper functions for converting OpenAPI DTOs to domain models.
 */

import type { TgvalidatordAddressInfo } from "../internal/openapi/models/TgvalidatordAddressInfo";
import type { TgvalidatordCurrency } from "../internal/openapi/models/TgvalidatordCurrency";
import type { TgvalidatordTransaction } from "../internal/openapi/models/TgvalidatordTransaction";
import type { TgvalidatordTransactionAttribute } from "../internal/openapi/models/TgvalidatordTransactionAttribute";
import type {
  Transaction,
  TransactionAddressInfo,
  TransactionAttribute,
  TransactionCurrencyInfo,
} from "../models/transaction";
import { TransactionStatus } from "../models/transaction";
import {
  safeBoolDefault,
  safeDate,
  safeInt,
  safeMap,
  safeString,
  safeStringDefault,
} from "./base";

/**
 * Maps a TgvalidatordAddressInfo DTO to a TransactionAddressInfo domain model.
 *
 * @param dto - The address info DTO from the OpenAPI response
 * @returns The TransactionAddressInfo domain model, or undefined if dto is null/undefined
 */
export function transactionAddressInfoFromDto(
  dto: TgvalidatordAddressInfo | null | undefined
): TransactionAddressInfo | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    address: safeString(dto.address),
    label: safeString(dto.label),
    container: safeString(dto.container),
    customerId: safeString(dto.customerId),
    amount: safeString(dto.amount),
    amountMainUnit: safeString(dto.amountMainUnit),
    type: safeString(dto.type),
    idx: safeString(dto.idx),
    internalAddressId: safeString(dto.internalAddressId),
    whitelistedAddressId: safeString(dto.whitelistedAddressId),
  };
}

/**
 * Maps a TgvalidatordCurrency DTO to a TransactionCurrencyInfo domain model.
 *
 * @param dto - The currency DTO from the OpenAPI response
 * @returns The TransactionCurrencyInfo domain model, or undefined if dto is null/undefined
 */
export function transactionCurrencyInfoFromDto(
  dto: TgvalidatordCurrency | null | undefined
): TransactionCurrencyInfo | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    id: safeString(dto.id),
    symbol: safeString(dto.symbol),
    name: safeString(dto.name),
    decimals: safeInt(dto.decimals),
    blockchain: safeString(dto.blockchain),
    network: safeString(dto.network),
  };
}

/**
 * Maps a TgvalidatordTransactionAttribute DTO to a TransactionAttribute domain model.
 *
 * @param dto - The attribute DTO from the OpenAPI response
 * @returns The TransactionAttribute domain model, or undefined if dto is null/undefined
 */
export function transactionAttributeFromDto(
  dto: TgvalidatordTransactionAttribute | null | undefined
): TransactionAttribute | undefined {
  if (!dto) {
    return undefined;
  }
  return {
    key: safeString(dto.key),
    value: safeString(dto.value),
  };
}

/**
 * Parses a transaction status string to the enum value.
 */
function parseTransactionStatus(
  status: string | undefined
): TransactionStatus | string | undefined {
  if (!status) {
    return undefined;
  }
  // Return enum value if it matches, otherwise return the raw string
  const upperStatus = status.toUpperCase();
  if (Object.values(TransactionStatus).includes(upperStatus as TransactionStatus)) {
    return upperStatus as TransactionStatus;
  }
  return status;
}

/**
 * Maps a TgvalidatordTransaction DTO to a Transaction domain model.
 *
 * @param dto - The transaction DTO from the OpenAPI response
 * @returns The Transaction domain model, or undefined if dto is null/undefined
 */
export function transactionFromDto(
  dto: TgvalidatordTransaction | null | undefined
): Transaction | undefined {
  if (!dto) {
    return undefined;
  }

  return {
    id: safeStringDefault(dto.id, ""),
    direction: safeString(dto.direction),
    currency: safeString(dto.currency),
    currencyInfo: transactionCurrencyInfoFromDto(dto.currencyInfo),
    blockchain: safeString(dto.blockchain),
    network: safeString(dto.network),
    hash: safeString(dto.hash),
    block: safeString(dto.block),
    confirmationBlock: safeString(dto.confirmationBlock),
    amount: safeString(dto.amount),
    amountMainUnit: safeString(dto.amountMainUnit),
    fee: safeString(dto.fee),
    feeMainUnit: safeString(dto.feeMainUnit),
    type: safeString(dto.type),
    status: parseTransactionStatus(dto.status),
    isConfirmed: safeBoolDefault(dto.isConfirmed, false),
    sources: safeMap(dto.sources, transactionAddressInfoFromDto),
    destinations: safeMap(dto.destinations, transactionAddressInfoFromDto),
    transactionId: safeString(dto.transactionId),
    uniqueId: safeString(dto.uniqueId),
    requestId: safeString(dto.requestId),
    requestVisible: dto.requestVisible,
    receptionDate: safeDate(dto.receptionDate),
    confirmationDate: safeDate(dto.confirmationDate),
    attributes: safeMap(dto.attributes, transactionAttributeFromDto),
    forkNumber: safeString(dto.forkNumber),
  };
}

/**
 * Maps an array of TgvalidatordTransaction DTOs to Transaction domain models.
 *
 * @param dtos - The array of transaction DTOs
 * @returns Array of Transaction domain models (undefined entries filtered out)
 */
export function transactionsFromDto(
  dtos: TgvalidatordTransaction[] | null | undefined
): Transaction[] {
  return safeMap(dtos, transactionFromDto);
}
