/**
 * Domain model classes for the Taurus Protect SDK.
 * <p>
 * This package contains the data model classes that represent entities in the
 * Taurus Protect system, including wallets, addresses, transactions, and requests.
 * <p>
 * <h2>Core Entities</h2>
 * <ul>
 *   <li>{@link com.taurushq.sdk.protect.client.model.Wallet} - Blockchain wallets</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.Address} - Blockchain addresses</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.Transaction} - Blockchain transactions</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.Request} - Transaction requests</li>
 * </ul>
 * <p>
 * <h2>Supporting Classes</h2>
 * <ul>
 *   <li>{@link com.taurushq.sdk.protect.client.model.Balance} - Balance information</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.Currency} - Currency metadata</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.Score} - Compliance scores</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.Attribute} - Custom attributes</li>
 * </ul>
 * <p>
 * <h2>Request Builders</h2>
 * <ul>
 *   <li>{@link com.taurushq.sdk.protect.client.model.CreateWalletRequest} - Wallet creation</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.CreateAddressRequest} - Address creation</li>
 * </ul>
 * <p>
 * <h2>Exceptions</h2>
 * <ul>
 *   <li>{@link com.taurushq.sdk.protect.client.model.ApiException} - Base API exception</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.ValidationException} - Input validation errors</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.AuthenticationException} - Authentication failures</li>
 *   <li>{@link com.taurushq.sdk.protect.client.model.IntegrityException} - Data integrity failures</li>
 * </ul>
 *
 * @see com.taurushq.sdk.protect.client.service
 */
package com.taurushq.sdk.protect.client.model;
