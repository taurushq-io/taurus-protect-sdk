package com.taurushq.sdk.protect.client.model;

import com.google.gson.JsonElement;
import com.google.gson.JsonParser;
import org.apache.commons.lang3.builder.ToStringBuilder;

import java.util.Map;

/**
 * Represents the metadata associated with a transaction request.
 * <p>
 * Request metadata contains detailed information about the transaction including
 * source and destination addresses, amount details, currency information, and
 * business rules applied. The metadata is signed and hashed for integrity verification.
 * <p>
 * The metadata payload is stored as a JSON array of key-value pairs, which can be
 * accessed through convenience methods like {@link #getSourceAddress()},
 * {@link #getDestinationAddress()}, and {@link #getAmount()}.
 *
 * <h2>SECURITY DESIGN</h2>
 * <p>
 * The API returns metadata with two representations of the same data:
 * <ul>
 *   <li>{@code payload} - Raw JSON object/array (UNVERIFIED)</li>
 *   <li>{@code payloadAsString} - JSON string that is cryptographically hashed (VERIFIED)</li>
 * </ul>
 * <p>
 * The security model works as follows:
 * <ol>
 *   <li>The server computes: {@code metadata.hash = SHA256(payloadAsString)}</li>
 *   <li>The hash is signed by governance rules (SuperAdmin keys)</li>
 *   <li>Clients verify: {@code computed_hash(payloadAsString) == metadata.hash}</li>
 * </ol>
 * <p>
 * <strong>ATTACK VECTOR (if using raw payload):</strong>
 * An attacker intercepting API responses could:
 * <ol>
 *   <li>Modify the payload object (e.g., change destination address)</li>
 *   <li>Leave payloadAsString unchanged (hash still verifies)</li>
 *   <li>Client extracts data from modified payload → SECURITY BYPASS</li>
 * </ol>
 * <p>
 * <strong>SOLUTION:</strong>
 * All data extraction methods in this class use {@code payloadAsJson}, which is
 * parsed from {@code payloadAsString} via {@link #setPayloadAsString(String)}.
 * This ensures:
 * <ul>
 *   <li>All extracted data comes from the cryptographically verified source</li>
 *   <li>Any tampering with the raw payload object is ignored</li>
 *   <li>The integrity chain: payloadAsString → hash → signature is preserved</li>
 * </ul>
 *
 * @see RequestMetadataAmount
 * @see RequestMetadataException
 * @see Request
 */
public class RequestMetadata {

    /**
     * Cryptographic hash of the metadata payload for integrity verification.
     */
    private String hash;

    // SECURITY: payload field intentionally omitted - use payloadAsString only.
    // The raw payload object could be tampered with by an attacker while
    // payloadAsString remains unchanged (hash still verifies). By not having
    // this field, we enforce that all data extraction uses the verified source.
    // See getSourceAddress(), getDestinationAddress(), etc. for secure data extraction.

    /**
     * String representation of the payload for parsing.
     */
    private String payloadAsString;

    /**
     * Parsed JSON representation of the payload for data extraction.
     */
    private JsonElement payloadAsJson;


    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the cryptographic hash of the metadata payload.
     * <p>
     * This hash is used for integrity verification of the metadata.
     *
     * @return the metadata hash
     */
    public String getHash() {
        return hash;
    }

    /**
     * Sets the cryptographic hash of the metadata payload.
     *
     * @param hash the metadata hash to set
     */
    public void setHash(String hash) {
        this.hash = hash;
    }

    // SECURITY: getPayload() and setPayload() intentionally removed.
    // Use payloadAsJson (set via setPayloadAsString) for data extraction.

    /**
     * Returns the payload as a JSON string.
     *
     * @return the payload as a JSON string
     */
    public String getPayloadAsString() {
        return payloadAsString;
    }

    /**
     * Sets the payload from a JSON string.
     * <p>
     * This method also parses the string into a JSON element for
     * data extraction via convenience methods.
     * <p>
     * <strong>SECURITY NOTE:</strong> This setter is critical for security.
     * The {@code payloadAsJson} field is parsed from {@code payloadAsString},
     * which is the cryptographically verified source. All data extraction
     * methods ({@link #getSourceAddress()}, {@link #getDestinationAddress()},
     * {@link #getAmount()}, etc.) use {@code payloadAsJson}, NOT the raw
     * {@code payload} field. This ensures the integrity verification chain
     * is preserved.
     *
     * @param payloadAsString the payload as a JSON string
     */
    public void setPayloadAsString(String payloadAsString) {
        this.payloadAsString = payloadAsString;
        this.payloadAsJson = JsonParser.parseString(payloadAsString);
    }

    /**
     * Extracts the request ID from the metadata payload.
     *
     * @return the request ID associated with this metadata
     * @throws RequestMetadataException if the request ID is not found in the metadata
     */
    public long getRequestId() throws RequestMetadataException {
        return this.getMDLong("request_id");
    }

    /**
     * Extracts the currency code from the metadata payload.
     *
     * @return the currency code (e.g., "BTC", "ETH")
     * @throws RequestMetadataException if the currency is not found in the metadata
     */
    public String getCurrency() throws RequestMetadataException {
        return this.getMDString("currency");
    }

    /**
     * Extracts the business rules key from the metadata payload.
     * <p>
     * The rules key identifies which set of business rules were applied
     * during request validation.
     *
     * @return the rules key identifier
     * @throws RequestMetadataException if the rules key is not found in the metadata
     */
    public String getRulesKey() throws RequestMetadataException {
        return this.getMDString("rules_key");
    }

    /**
     * Extracts the source blockchain address from the metadata payload.
     * <p>
     * This is the address from which funds are being sent.
     * <p>
     * <strong>SECURITY NOTE:</strong> This method extracts from {@code payloadAsJson},
     * which is parsed from the cryptographically verified {@code payloadAsString}.
     * See class Javadoc for security rationale.
     *
     * @return the source blockchain address
     * @throws RequestMetadataException if the source address is not found in the metadata
     */
    public String getSourceAddress() throws RequestMetadataException {
        for (JsonElement bm : this.payloadAsJson.getAsJsonArray()) {
            Map<String, JsonElement> map = bm.getAsJsonObject().asMap();
            if ("source".equals(map.get("key").getAsString())) {
                com.google.gson.JsonObject value = map.get("value").getAsJsonObject();
                com.google.gson.JsonObject payload = value.getAsJsonObject("payload");
                // Some source types (e.g., SourceInternalWallet) don't have an address field
                if (payload != null && payload.has("address")) {
                    return payload.get("address").getAsString();
                }
                throw new RequestMetadataException("source address not available for source type: " + value.get("type"));
            }
        }
        throw new RequestMetadataException("source address not found in metadata");
    }

    /**
     * Extracts the destination blockchain address from the metadata payload.
     * <p>
     * This is the address to which funds are being sent.
     * <p>
     * <strong>SECURITY NOTE:</strong> This method extracts from {@code payloadAsJson},
     * which is parsed from the cryptographically verified {@code payloadAsString}.
     * See class Javadoc for security rationale.
     *
     * @return the destination blockchain address
     * @throws RequestMetadataException if the destination address is not found in the metadata
     */
    public String getDestinationAddress() throws RequestMetadataException {
        for (JsonElement bm : this.payloadAsJson.getAsJsonArray()) {
            Map<String, JsonElement> map = bm.getAsJsonObject().asMap();
            if ("destination".equals(map.get("key").getAsString())) {
                com.google.gson.JsonObject value = map.get("value").getAsJsonObject();
                com.google.gson.JsonObject payload = value.getAsJsonObject("payload");
                // Some destination types may not have an address field
                if (payload != null && payload.has("address")) {
                    return payload.get("address").getAsString();
                }
                throw new RequestMetadataException("destination address not available for destination type: " + value.get("type"));
            }
        }
        throw new RequestMetadataException("destination address not found in metadata");
    }

    /**
     * Extracts the amount details from the metadata payload.
     * <p>
     * The amount includes the value in both source and target currencies,
     * conversion rate, and decimal precision.
     * <p>
     * <strong>SECURITY NOTE:</strong> This method extracts from {@code payloadAsJson},
     * which is parsed from the cryptographically verified {@code payloadAsString}.
     * See class Javadoc for security rationale.
     *
     * @return the amount details for this transaction
     * @throws RequestMetadataException if the amount is not found in the metadata
     */
    public RequestMetadataAmount getAmount() throws RequestMetadataException {
        for (JsonElement bm : this.payloadAsJson.getAsJsonArray()) {
            Map<String, JsonElement> map = bm.getAsJsonObject().asMap();
            if ("amount".equals(map.get("key").getAsString())) {

                com.google.gson.JsonObject value = map.get("value").getAsJsonObject();
                RequestMetadataAmount amount = new RequestMetadataAmount();
                amount.setValueFrom(jsonValueToString(value.get("valueFrom")));
                amount.setValueTo(jsonValueToString(value.get("valueTo")));
                amount.setRate(jsonValueToString(value.get("rate")));
                amount.setDecimals(value.get("decimals").getAsInt());
                amount.setCurrencyFrom(value.get("currencyFrom").getAsString());
                amount.setCurrencyTo(value.get("currencyTo").getAsString());

                return amount;
            }
        }
        throw new RequestMetadataException("amount not found in metadata");
    }

    /**
     * Converts a JSON value to its string representation, preserving arbitrary precision.
     * Handles both string and numeric JSON values for backward compatibility.
     *
     * @param element the JSON element to convert
     * @return the string representation, or null if the element is null
     */
    private static String jsonValueToString(JsonElement element) {
        if (element == null || element.isJsonNull()) {
            return null;
        }
        if (element.getAsJsonPrimitive().isString()) {
            return element.getAsString();
        }
        // For numeric values, use BigDecimal for lossless conversion
        return new java.math.BigDecimal(element.getAsString()).toPlainString();
    }


    /**
     * Retrieves a string value from the metadata payload by key.
     *
     * @param key the key to look up in the metadata
     * @return the string value for the specified key
     * @throws RequestMetadataException if the key is not found in the metadata
     */
    private String getMDString(String key) throws RequestMetadataException {
        for (JsonElement bm : this.payloadAsJson.getAsJsonArray()) {
            Map<String, JsonElement> map = bm.getAsJsonObject().asMap();
            if (map.get("key").getAsString().equals(key)) {
                return map.get("value").getAsString();
            }
        }
        throw new RequestMetadataException("key '" + key + "' not found in metadata");
    }

    /**
     * Retrieves a long value from the metadata payload by key.
     *
     * @param key the key to look up in the metadata
     * @return the long value for the specified key
     * @throws RequestMetadataException if the key is not found in the metadata
     */
    private long getMDLong(String key) throws RequestMetadataException {
        for (JsonElement bm : this.payloadAsJson.getAsJsonArray()) {
            Map<String, JsonElement> map = bm.getAsJsonObject().asMap();
            if (map.get("key").getAsString().equals(key)) {
                return map.get("value").getAsLong();
            }
        }
        throw new RequestMetadataException("key '" + key + "' not found in metadata");
    }
}
