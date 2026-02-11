/**
 * MapStruct mappers for converting between OpenAPI DTOs and domain models.
 * <p>
 * This package contains mapper interfaces that use MapStruct to automatically
 * generate code for converting between the auto-generated OpenAPI client DTOs
 * and the SDK's domain model classes.
 * <p>
 * <h2>Available Mappers</h2>
 * <ul>
 *   <li>{@link com.taurushq.sdk.protect.client.mapper.WalletMapper} - Wallet conversion</li>
 *   <li>{@link com.taurushq.sdk.protect.client.mapper.AddressMapper} - Address conversion</li>
 *   <li>{@link com.taurushq.sdk.protect.client.mapper.TransactionMapper} - Transaction conversion</li>
 *   <li>{@link com.taurushq.sdk.protect.client.mapper.RequestMapper} - Request conversion</li>
 *   <li>{@link com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper} - Exception conversion</li>
 * </ul>
 * <p>
 * <h2>Implementation Note</h2>
 * Mapper implementations are generated at compile time by MapStruct and can be
 * found in {@code target/generated-sources/annotations}. Each mapper interface
 * provides a singleton INSTANCE field for convenient access.
 * <p>
 * <h2>Example Usage</h2>
 * <pre>{@code
 * // Convert a single DTO
 * Wallet wallet = WalletMapper.INSTANCE.fromDTO(walletInfoDto);
 *
 * // Convert a list of DTOs
 * List<Transaction> transactions = TransactionMapper.INSTANCE.fromDTO(transactionDtos);
 * }</pre>
 *
 * @see org.mapstruct.Mapper
 */
package com.taurushq.sdk.protect.client.mapper;
