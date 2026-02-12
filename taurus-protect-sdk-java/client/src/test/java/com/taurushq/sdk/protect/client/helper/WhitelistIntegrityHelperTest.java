package com.taurushq.sdk.protect.client.helper;

import com.taurushq.sdk.protect.client.mapper.WhitelistedAddressMapper;
import com.taurushq.sdk.protect.client.mapper.WhitelistedContractAddressMapper;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import com.taurushq.sdk.protect.client.model.WhitelistedContractAddress;
import com.taurushq.sdk.protect.proto.v1.Whitelist;
import org.junit.jupiter.api.Test;

import java.util.Base64;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Tests for WhitelistIntegrityHelper that verify envelope field validation.
 * Test vectors are taken from the Go whitelist_test.go file.
 */
class WhitelistIntegrityHelperTest {

    // Base64 envelope from Go test for ChainLink Token contract address
    // blockchain=ETH, name=ChainLink Token, symbol=LINK, decimals=18, contractAddress=0x514910771af9ca656af840dff83e8264ecf986ca
    private static final String CONTRACT_ENVELOPE =
            "CgNFVEgSD0NoYWluTGluayBUb2tlbhoETElOSyASKioweDUxNDkxMDc3MWFmOWNhNjU2YWY4NDBkZmY4M2U4MjY0ZWNmOTg2Y2E=";

    // Base64 envelope for a whitelisted address
    private static final String ADDRESS_ENVELOPE =
            "CgNFVEgQARoqMHhmNjMxY2U4OTNlZGI0NDBlNDkxODhhOTkxMjUwNTFkMDc5NjgxODY0KkBteSBzZWNvbmQgRVRIIGFkZHJlc3Mgb24gYmlzdGFtcCAoYWN0dWFsbHksIGFuIGludGVybmFsIGFkZHJlc3MpOAM=";

    /**
     * Helper to create a test contract address envelope.
     */
    private String createContractEnvelope(String blockchain, String name, String symbol,
                                           long decimals, String contractAddress) {
        Whitelist.WhitelistedContractAddress proto = Whitelist.WhitelistedContractAddress.newBuilder()
                .setBlockchain(blockchain)
                .setName(name)
                .setSymbol(symbol)
                .setDecimals(decimals)
                .setContractAddress(contractAddress)
                .build();
        return Base64.getEncoder().encodeToString(proto.toByteArray());
    }

    /**
     * Helper to create a test address envelope.
     */
    private String createAddressEnvelope(String blockchain, String address, String label,
                                          String memo, String customerId,
                                          Whitelist.WhitelistedAddress.AddressType addressType) {
        Whitelist.WhitelistedAddress.Builder builder = Whitelist.WhitelistedAddress.newBuilder()
                .setBlockchain(blockchain)
                .setAddress(address)
                .setLabel(label);
        if (memo != null) {
            builder.setMemo(memo);
        }
        if (customerId != null) {
            builder.setCustomerId(customerId);
        }
        if (addressType != null) {
            builder.setAddressType(addressType);
        }
        return Base64.getEncoder().encodeToString(builder.build().toByteArray());
    }

    // ========== Contract Address Tests ==========

    @Test
    void testVerifyWLContractAddressIntegrity_ValidMatch() {
        String envelope = createContractEnvelope("ETH", "ChainLink Token", "LINK", 18,
                "0x514910771af9ca656af840dff83e8264ecf986ca");

        // Create DB address that matches envelope
        WhitelistedContractAddress dbAddr = new WhitelistedContractAddress();
        dbAddr.setBlockchain("ETH");
        dbAddr.setName("ChainLink Token");
        dbAddr.setSymbol("LINK");
        dbAddr.setDecimals(18);
        dbAddr.setContractAddress("0x514910771af9ca656af840dff83e8264ecf986ca");

        // Should not throw
        assertDoesNotThrow(() ->
                WhitelistIntegrityHelper.verifyWLContractAddressIntegrity(dbAddr, envelope));
    }

    @Test
    void testVerifyWLContractAddressIntegrity_InvalidBlockchain() {
        String envelope = createContractEnvelope("ETH", "ChainLink Token", "LINK", 18,
                "0x514910771af9ca656af840dff83e8264ecf986ca");

        // Create DB address with different blockchain
        WhitelistedContractAddress dbAddr = new WhitelistedContractAddress();
        dbAddr.setBlockchain("BTC"); // Different from envelope (ETH)
        dbAddr.setName("ChainLink Token");
        dbAddr.setSymbol("LINK");
        dbAddr.setDecimals(18);
        dbAddr.setContractAddress("0x514910771af9ca656af840dff83e8264ecf986ca");

        WhitelistException ex = assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLContractAddressIntegrity(dbAddr, envelope));
        assertTrue(ex.getMessage().contains("field Blockchain"));
    }

    @Test
    void testVerifyWLContractAddressIntegrity_InvalidName() {
        String envelope = createContractEnvelope("ETH", "ChainLink Token", "LINK", 18,
                "0x514910771af9ca656af840dff83e8264ecf986ca");

        WhitelistedContractAddress dbAddr = new WhitelistedContractAddress();
        dbAddr.setBlockchain("ETH");
        dbAddr.setName("DifferentName"); // Different from envelope
        dbAddr.setSymbol("LINK");
        dbAddr.setDecimals(18);
        dbAddr.setContractAddress("0x514910771af9ca656af840dff83e8264ecf986ca");

        WhitelistException ex = assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLContractAddressIntegrity(dbAddr, envelope));
        assertTrue(ex.getMessage().contains("field Name"));
    }

    @Test
    void testVerifyWLContractAddressIntegrity_InvalidSymbol() {
        String envelope = createContractEnvelope("ETH", "ChainLink Token", "LINK", 18,
                "0x514910771af9ca656af840dff83e8264ecf986ca");

        WhitelistedContractAddress dbAddr = new WhitelistedContractAddress();
        dbAddr.setBlockchain("ETH");
        dbAddr.setName("ChainLink Token");
        dbAddr.setSymbol("XXX"); // Different from envelope
        dbAddr.setDecimals(18);
        dbAddr.setContractAddress("0x514910771af9ca656af840dff83e8264ecf986ca");

        WhitelistException ex = assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLContractAddressIntegrity(dbAddr, envelope));
        assertTrue(ex.getMessage().contains("field Symbol"));
    }

    @Test
    void testVerifyWLContractAddressIntegrity_InvalidDecimals() {
        String envelope = createContractEnvelope("ETH", "ChainLink Token", "LINK", 18,
                "0x514910771af9ca656af840dff83e8264ecf986ca");

        WhitelistedContractAddress dbAddr = new WhitelistedContractAddress();
        dbAddr.setBlockchain("ETH");
        dbAddr.setName("ChainLink Token");
        dbAddr.setSymbol("LINK");
        dbAddr.setDecimals(8); // Different from envelope
        dbAddr.setContractAddress("0x514910771af9ca656af840dff83e8264ecf986ca");

        WhitelistException ex = assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLContractAddressIntegrity(dbAddr, envelope));
        assertTrue(ex.getMessage().contains("field Decimals"));
    }

    @Test
    void testVerifyWLContractAddressIntegrity_InvalidContractAddress() {
        String envelope = createContractEnvelope("ETH", "ChainLink Token", "LINK", 18,
                "0x514910771af9ca656af840dff83e8264ecf986ca");

        WhitelistedContractAddress dbAddr = new WhitelistedContractAddress();
        dbAddr.setBlockchain("ETH");
        dbAddr.setName("ChainLink Token");
        dbAddr.setSymbol("LINK");
        dbAddr.setDecimals(18);
        dbAddr.setContractAddress("0x0000000000000000000000000000000000000000"); // Different

        WhitelistException ex = assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLContractAddressIntegrity(dbAddr, envelope));
        assertTrue(ex.getMessage().contains("field ContractAddress"));
    }

    @Test
    void testVerifyWLContractAddressIntegrity_NullDbAddress() {
        String envelope = createContractEnvelope("ETH", "Test", "TST", 18, "0x1234");
        assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLContractAddressIntegrity(null, envelope));
    }

    @Test
    void testVerifyWLContractAddressIntegrity_NullEnvelope() {
        WhitelistedContractAddress dbAddr = new WhitelistedContractAddress();
        assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLContractAddressIntegrity(dbAddr, null));
    }

    @Test
    void testVerifyWLContractAddressIntegrity_EmptyEnvelope() {
        WhitelistedContractAddress dbAddr = new WhitelistedContractAddress();
        assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLContractAddressIntegrity(dbAddr, ""));
    }

    // ========== Address Tests ==========

    @Test
    void testVerifyWLAddressIntegrity_ValidMatch() throws Exception {
        WhitelistedAddress envelopeAddr =
                WhitelistedAddressMapper.INSTANCE.fromBase64String(ADDRESS_ENVELOPE);

        // Create DB address that matches envelope
        WhitelistedAddress dbAddr = new WhitelistedAddress();
        dbAddr.setBlockchain(envelopeAddr.getBlockchain());
        dbAddr.setAddress(envelopeAddr.getAddress());
        dbAddr.setLabel(envelopeAddr.getLabel());
        dbAddr.setMemo(envelopeAddr.getMemo());
        dbAddr.setCustomerId(envelopeAddr.getCustomerId());
        dbAddr.setAddressType(envelopeAddr.getAddressType());

        // Should not throw
        assertDoesNotThrow(() ->
                WhitelistIntegrityHelper.verifyWLAddressIntegrity(dbAddr, ADDRESS_ENVELOPE));
    }

    @Test
    void testVerifyWLAddressIntegrity_InvalidBlockchain() throws Exception {
        WhitelistedAddress envelopeAddr =
                WhitelistedAddressMapper.INSTANCE.fromBase64String(ADDRESS_ENVELOPE);

        WhitelistedAddress dbAddr = new WhitelistedAddress();
        dbAddr.setBlockchain("BTC"); // Different from envelope (ETH)
        dbAddr.setAddress(envelopeAddr.getAddress());
        dbAddr.setLabel(envelopeAddr.getLabel());
        dbAddr.setMemo(envelopeAddr.getMemo());
        dbAddr.setCustomerId(envelopeAddr.getCustomerId());
        dbAddr.setAddressType(envelopeAddr.getAddressType());

        WhitelistException ex = assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLAddressIntegrity(dbAddr, ADDRESS_ENVELOPE));
        assertTrue(ex.getMessage().contains("field Blockchain"));
    }

    @Test
    void testVerifyWLAddressIntegrity_InvalidAddress() throws Exception {
        WhitelistedAddress envelopeAddr =
                WhitelistedAddressMapper.INSTANCE.fromBase64String(ADDRESS_ENVELOPE);

        WhitelistedAddress dbAddr = new WhitelistedAddress();
        dbAddr.setBlockchain(envelopeAddr.getBlockchain());
        dbAddr.setAddress("0x0000000000000000000000000000000000000000"); // Different
        dbAddr.setLabel(envelopeAddr.getLabel());
        dbAddr.setMemo(envelopeAddr.getMemo());
        dbAddr.setCustomerId(envelopeAddr.getCustomerId());
        dbAddr.setAddressType(envelopeAddr.getAddressType());

        WhitelistException ex = assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLAddressIntegrity(dbAddr, ADDRESS_ENVELOPE));
        assertTrue(ex.getMessage().contains("field Address"));
    }

    @Test
    void testVerifyWLAddressIntegrity_InvalidLabel() throws Exception {
        WhitelistedAddress envelopeAddr =
                WhitelistedAddressMapper.INSTANCE.fromBase64String(ADDRESS_ENVELOPE);

        WhitelistedAddress dbAddr = new WhitelistedAddress();
        dbAddr.setBlockchain(envelopeAddr.getBlockchain());
        dbAddr.setAddress(envelopeAddr.getAddress());
        dbAddr.setLabel("Different Label"); // Different
        dbAddr.setMemo(envelopeAddr.getMemo());
        dbAddr.setCustomerId(envelopeAddr.getCustomerId());
        dbAddr.setAddressType(envelopeAddr.getAddressType());

        WhitelistException ex = assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLAddressIntegrity(dbAddr, ADDRESS_ENVELOPE));
        assertTrue(ex.getMessage().contains("field Label"));
    }

    @Test
    void testVerifyWLAddressIntegrity_NullDbAddress() {
        assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLAddressIntegrity(null, ADDRESS_ENVELOPE));
    }

    @Test
    void testVerifyWLAddressIntegrity_NullEnvelope() {
        WhitelistedAddress dbAddr = new WhitelistedAddress();
        assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLAddressIntegrity(dbAddr, null));
    }

    @Test
    void testVerifyWLAddressIntegrity_EmptyEnvelope() {
        WhitelistedAddress dbAddr = new WhitelistedAddress();
        assertThrows(WhitelistException.class, () ->
                WhitelistIntegrityHelper.verifyWLAddressIntegrity(dbAddr, ""));
    }
}
