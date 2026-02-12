/**
 * Service classes for interacting with the Taurus Protect API.
 * <p>
 * This package contains service classes that provide high-level operations
 * for managing wallets, addresses, transactions, and requests. Services are
 * accessed through the {@link com.taurushq.sdk.protect.client.ProtectClient}.
 * <p>
 * <h2>Available Services</h2>
 * <ul>
 *   <li>{@link com.taurushq.sdk.protect.client.service.WalletService} - Wallet management</li>
 *   <li>{@link com.taurushq.sdk.protect.client.service.AddressService} - Address management</li>
 *   <li>{@link com.taurushq.sdk.protect.client.service.TransactionService} - Transaction queries</li>
 *   <li>{@link com.taurushq.sdk.protect.client.service.RequestService} - Request lifecycle management</li>
 * </ul>
 * <p>
 * <h2>Error Handling</h2>
 * All service methods throw {@link com.taurushq.sdk.protect.client.model.ApiException}
 * or one of its subclasses on failure:
 * <ul>
 *   <li>{@link com.taurushq.sdk.protect.client.model.ValidationException} - Invalid input (400)</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.AuthenticationException} - Authentication failure (401)</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.AuthorizationException} - Insufficient permissions (403)</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.NotFoundException} - Resource not found (404)</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.RateLimitException} - Rate limit exceeded (429)</li>
 * </ul>
 * <p>
 * <h2>Example Usage</h2>
 * <pre>{@code
 * try {
 *     Wallet wallet = client.getWalletService().getWallet(walletId);
 * } catch (NotFoundException e) {
 *     System.out.println("Wallet not found: " + e.getMessage());
 * } catch (RateLimitException e) {
 *     // Wait and retry
 *     Thread.sleep(e.getSuggestedRetryDelayMs());
 * } catch (ApiException e) {
 *     System.out.println("API error: " + e.getMessage());
 * }
 * }</pre>
 *
 * @see com.taurushq.sdk.protect.client.ProtectClient
 * @see com.taurushq.sdk.protect.client.model.ApiException
 */
package com.taurushq.sdk.protect.client.service;
