/**
 * Taurus Protect SDK for Java - High-level client library.
 * <p>
 * This package provides the main entry point for interacting with the Taurus Protect API.
 * The primary class is {@link com.taurushq.sdk.protect.client.ProtectClient}, which provides
 * access to all SDK services.
 * <p>
 * <h2>Getting Started</h2>
 * <pre>{@code
 * // Create a client using the builder
 * ProtectClient client = ProtectClient.builder()
 *     .host("https://api.protect.taurushq.com")
 *     .credentials(apiKey, apiSecret)
 *     .superAdminKeysPem(pemKeys)
 *     .minValidSignatures(2)
 *     .build();
 *
 * // Access services
 * WalletService walletService = client.getWalletService();
 * AddressService addressService = client.getAddressService();
 * RequestService requestService = client.getRequestService();
 * TransactionService transactionService = client.getTransactionService();
 * }</pre>
 * <p>
 * <h2>Package Structure</h2>
 * <ul>
 *   <li>{@link com.taurushq.sdk.protect.client.model} - Domain model classes</li>
 *   <li>{@link com.taurushq.sdk.protect.client.service} - Service classes for API operations</li>
 *   <li>{@link com.taurushq.sdk.protect.client.mapper} - MapStruct mappers for DTO conversion</li>
 *   <li>{@link com.taurushq.sdk.protect.client.helper} - Utility and helper classes</li>
 * </ul>
 *
 * @see com.taurushq.sdk.protect.client.ProtectClient
 * @see com.taurushq.sdk.protect.client.ProtectClientBuilder
 */
package com.taurushq.sdk.protect.client;
