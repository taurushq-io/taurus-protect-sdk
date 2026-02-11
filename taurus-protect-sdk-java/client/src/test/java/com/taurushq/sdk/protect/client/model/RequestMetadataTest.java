package com.taurushq.sdk.protect.client.model;

import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;

class RequestMetadataTest {

    private final String payload = "[{\"key\":\"request_id\",\"type\":\"String\",\"value\":\"6286376\",\"column\":\"\"},{\"key\":\"rules_key\",\"type\":\"String\",\"value\":\"XLM\",\"column\":\"\"},{\"key\":\"currency\",\"type\":\"String\",\"value\":\"XLM\",\"column\":\"\"},{\"key\":\"currency_name\",\"type\":\"String\",\"value\":\"Stellar - XLM\",\"column\":\"\"},{\"key\":\"currency_id\",\"type\":\"String\",\"value\":\"1a31194be3de1a85216300c2c0640af1d9a87a2f600a2e8927b91f78740e10ff\",\"column\":\"\"},{\"key\":\"source\",\"type\":\"Source\",\"value\":{\"type\":\"SourceInternalAddress\",\"payload\":{\"id\":\"4\",\"address\":\"GBLNAHS75FMDPFPBZSEH2VC5ZAYVRETDVQDUT44M4T4FXWMJTZ6I2I2B\",\"label\":\"nostro XLM\",\"path\":\"m/44'/148'/0'\"}},\"column\":\"source\"},{\"key\":\"destination\",\"type\":\"Destination\",\"value\":{\"type\":\"DestinationInternalAddress\",\"payload\":{\"id\":\"6\",\"address\":\"GDWIH4JIZEDCHIRNQUQAOPPBGVUS23KCVH6JH7NM2D55LIC66QVY4MK6\",\"label\":\"stellar 1\",\"path\":\"m/44'/148'/1'\"}},\"column\":\"destination\"},{\"key\":\"amount\",\"type\":\"Amount\",\"value\":{\"valueFrom\":\"2\",\"valueTo\":\"0.0000\",\"rate\":\"0.08386473412745667\",\"decimals\":\"7\",\"currencyFrom\":\"XLM\",\"currencyTo\":\"CHF\"},\"column\":\"amount\"},{\"key\":\"total_fiat_amount\",\"type\":\"String\",\"value\":\"0.0000\",\"column\":\"\"},{\"key\":\"fee_limits\",\"type\":\"BigIntArray\",\"value\":[\"1500000\"],\"column\":\"\"},{\"key\":\"fee_paid_by_receiver\",\"type\":\"String\",\"value\":\"false\",\"column\":\"\"},{\"key\":\"use_unconfirmed_funds\",\"type\":\"String\",\"value\":\"false\",\"column\":\"\"},{\"key\":\"is_cancel_request\",\"type\":\"String\",\"value\":\"false\",\"column\":\"\"},{\"key\":\"transaction_ids\",\"type\":\"StringArray\",\"value\":[\"4:7671e0c5-623d-482a-9abd-ed90a262f6f6\"],\"column\":\"\"}]";

    @Test
    void getHash() {
        RequestMetadata metadata = new RequestMetadata();
        metadata.setHash("bcc15c43250c399221c421afb3e5848c1f02f3b6c51a86006282deaea2982077");
        metadata.setPayloadAsString(payload);
        assertEquals("bcc15c43250c399221c421afb3e5848c1f02f3b6c51a86006282deaea2982077", metadata.getHash());
        assertEquals("bcc15c43250c399221c421afb3e5848c1f02f3b6c51a86006282deaea2982077", CryptoTPV1.calculateHexHash(metadata.getPayloadAsString()));
    }

    @Test
    void getPayloadAsString() {
        RequestMetadata metadata = new RequestMetadata();
        metadata.setPayloadAsString(payload);
        assertEquals(payload, metadata.getPayloadAsString());
    }

    @Test
    void getRequestId() throws RequestMetadataException {
        RequestMetadata metadata = new RequestMetadata();
        metadata.setPayloadAsString(payload);
        assertEquals(6286376, metadata.getRequestId());
    }

    @Test
    void getCurrency() throws RequestMetadataException {
        RequestMetadata metadata = new RequestMetadata();
        metadata.setPayloadAsString(payload);
        assertEquals("XLM", metadata.getCurrency());
    }

    @Test
    void getRulesKey() throws RequestMetadataException {
        RequestMetadata metadata = new RequestMetadata();
        metadata.setPayloadAsString(payload);
        assertEquals("XLM", metadata.getRulesKey());
    }

    @Test
    void getSourceAddress() throws RequestMetadataException {
        RequestMetadata metadata = new RequestMetadata();
        metadata.setPayloadAsString(payload);
        assertEquals("GBLNAHS75FMDPFPBZSEH2VC5ZAYVRETDVQDUT44M4T4FXWMJTZ6I2I2B", metadata.getSourceAddress());
    }

    @Test
    void getDestinationAddress() throws RequestMetadataException {
        RequestMetadata metadata = new RequestMetadata();
        metadata.setPayloadAsString(payload);
        assertEquals("GDWIH4JIZEDCHIRNQUQAOPPBGVUS23KCVH6JH7NM2D55LIC66QVY4MK6", metadata.getDestinationAddress());
    }

    @Test
    void getAmount() throws RequestMetadataException {
        RequestMetadata metadata = new RequestMetadata();
        metadata.setPayloadAsString(payload);
        assertEquals("XLM", metadata.getAmount().getCurrencyFrom());
        assertEquals("CHF", metadata.getAmount().getCurrencyTo());
        assertEquals(7, metadata.getAmount().getDecimals());
        assertEquals("0.08386473412745667", metadata.getAmount().getRate());
        assertEquals("2", metadata.getAmount().getValueFrom());
        assertEquals("0.0000", metadata.getAmount().getValueTo());
    }
}