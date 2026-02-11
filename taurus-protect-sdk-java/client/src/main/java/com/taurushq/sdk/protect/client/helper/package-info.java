/**
 * Utility and helper classes for the Taurus Protect SDK.
 * <p>
 * This package contains helper classes that provide common functionality
 * used throughout the SDK, including validation, cryptographic operations,
 * and signature verification.
 * <p>
 * <h2>Available Helpers</h2>
 * <ul>
 *   <li>{@link com.taurushq.sdk.protect.client.helper.ValidationHelper} - Input validation utilities</li>
 *   <li>{@link com.taurushq.sdk.protect.client.helper.AddressSignatureVerifier} - Address signature verification</li>
 * </ul>
 * <p>
 * <h2>Validation Example</h2>
 * <pre>{@code
 * public void createWallet(String name, long balance) {
 *     ValidationHelper.requireNotBlank(name, "wallet name");
 *     ValidationHelper.requirePositive(balance, "initial balance");
 *     // proceed with creation
 * }
 * }</pre>
 *
 * @see com.taurushq.sdk.protect.client.helper.ValidationHelper
 */
package com.taurushq.sdk.protect.client.helper;
