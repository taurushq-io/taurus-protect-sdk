package com.taurushq.sdk.protect.client.helper;

import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.mapper.WhitelistedAddressMapper;
import com.taurushq.sdk.protect.client.mapper.WhitelistedContractAddressMapper;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import com.taurushq.sdk.protect.client.model.WhitelistedContractAddress;

import java.util.Objects;

/**
 * Helper class for verifying whitelist integrity.
 * Validates that envelope fields match database fields.
 * Equivalent to Go's VerifyWLAddressesIntegrity and VerifyWLContractAddressesIntegrity functions.
 */
public final class WhitelistIntegrityHelper {

    private WhitelistIntegrityHelper() {
        // Prevent instantiation
    }

    /**
     * Verifies that the whitelisted address envelope fields match the database fields.
     *
     * @param dbAddress      the address from the database
     * @param envelopeBase64 the base64-encoded envelope
     * @throws WhitelistException if validation fails or fields don't match
     */
    public static void verifyWLAddressIntegrity(WhitelistedAddress dbAddress,
                                                 String envelopeBase64) throws WhitelistException {
        if (dbAddress == null) {
            throw new WhitelistException("database address cannot be null");
        }
        if (envelopeBase64 == null || envelopeBase64.isEmpty()) {
            throw new WhitelistException("envelope cannot be null or empty");
        }

        WhitelistedAddress envelopeAddress;
        try {
            envelopeAddress = WhitelistedAddressMapper.INSTANCE.fromBase64String(envelopeBase64);
        } catch (InvalidProtocolBufferException e) {
            throw new WhitelistException("failed to decode envelope", e);
        }

        // Validate fields match
        validateFieldMatch("Blockchain", dbAddress.getBlockchain(), envelopeAddress.getBlockchain());
        validateFieldMatch("Address", dbAddress.getAddress(), envelopeAddress.getAddress());
        validateFieldMatch("Label", dbAddress.getLabel(), envelopeAddress.getLabel());
        validateFieldMatch("Memo", dbAddress.getMemo(), envelopeAddress.getMemo());
        validateFieldMatch("CustomerId", dbAddress.getCustomerId(), envelopeAddress.getCustomerId());
        validateFieldMatch("AddressType",
                dbAddress.getAddressType() != null ? dbAddress.getAddressType().name() : null,
                envelopeAddress.getAddressType() != null ? envelopeAddress.getAddressType().name() : null);
    }

    /**
     * Verifies that the whitelisted contract address envelope fields match the database fields.
     *
     * @param dbContractAddress the contract address from the database
     * @param envelopeBase64    the base64-encoded envelope
     * @throws WhitelistException if validation fails or fields don't match
     */
    public static void verifyWLContractAddressIntegrity(WhitelistedContractAddress dbContractAddress,
                                                         String envelopeBase64) throws WhitelistException {
        if (dbContractAddress == null) {
            throw new WhitelistException("database contract address cannot be null");
        }
        if (envelopeBase64 == null || envelopeBase64.isEmpty()) {
            throw new WhitelistException("envelope cannot be null or empty");
        }

        WhitelistedContractAddress envelopeAddress;
        try {
            envelopeAddress = WhitelistedContractAddressMapper.INSTANCE.fromBase64String(envelopeBase64);
        } catch (InvalidProtocolBufferException e) {
            throw new WhitelistException("failed to decode envelope", e);
        }

        // Validate fields match
        validateFieldMatch("Blockchain", dbContractAddress.getBlockchain(), envelopeAddress.getBlockchain());
        validateFieldMatch("Name", dbContractAddress.getName(), envelopeAddress.getName());
        validateFieldMatch("Symbol", dbContractAddress.getSymbol(), envelopeAddress.getSymbol());
        validateFieldMatch("Decimals",
                String.valueOf(dbContractAddress.getDecimals()),
                String.valueOf(envelopeAddress.getDecimals()));
        validateFieldMatch("ContractAddress",
                dbContractAddress.getContractAddress(), envelopeAddress.getContractAddress());
    }

    /**
     * Validates that two field values match.
     *
     * @param fieldName     the name of the field being validated
     * @param dbValue       the value from the database
     * @param envelopeValue the value from the envelope
     * @throws WhitelistException if the values don't match
     */
    private static void validateFieldMatch(String fieldName, String dbValue, String envelopeValue)
            throws WhitelistException {
        if (!Objects.equals(dbValue, envelopeValue)) {
            throw new WhitelistException(String.format(
                    "invalid whitelist signature: field %s mismatch (db='%s', envelope='%s')",
                    fieldName, dbValue, envelopeValue));
        }
    }
}
