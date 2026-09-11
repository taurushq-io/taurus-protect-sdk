package com.taurushq.sdk.protect.client.model;

import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The verified envelopes used to expose a PUBLIC setter that took the asset/address AND
 * flipped the same {@code isInitialized} gate the getter checks. A caller could construct
 * an envelope, inject a fabricated value, and read it back with no exception — the
 * "verified" marker returned attacker-chosen data. Go's {@code helper.VerifiedAsset} is
 * forgeable but USELESS (unexported fields), and these now have the same property:
 * {@code markVerified} derives the value from the envelope's own signed payload.
 */
class VerifiedEnvelopeDerivationTest {

    private static final String ASSET_PAYLOAD =
            "{\"name\":\"USDC\",\"symbol\":\"USDC\",\"contractAddress\":\"0xREAL\","
                    + "\"decimals\":6,\"blockchain\":\"ETH\"}";

    @Test
    @DisplayName("the asset comes from the envelope's own payload, not from a caller")
    void assetIsDerivedFromItsOwnPayload() throws Exception {
        SignedWhitelistedAssetEnvelope envelope = new SignedWhitelistedAssetEnvelope();
        WhitelistMetadata metadata = new WhitelistMetadata();
        metadata.setPayloadAsString(ASSET_PAYLOAD);
        envelope.setMetadata(metadata);

        envelope.markVerified(new DecodedRulesContainer(), ASSET_PAYLOAD);

        // There is no seam to inject "0xATTACKER" through: markVerified takes no asset,
        // only the payload a counted signature covered.
        assertEquals("0xREAL", envelope.getWhitelistedAsset().getContractAddress());
        assertEquals("USDC", envelope.getWhitelistedAsset().getSymbol());
    }

    @Test
    @DisplayName("an envelope with no signed payload cannot be marked verified")
    void payloadLessEnvelopeCannotBeMarked() {
        SignedWhitelistedAssetEnvelope envelope = new SignedWhitelistedAssetEnvelope();
        envelope.setMetadata(new WhitelistMetadata());

        WhitelistException e = assertThrows(WhitelistException.class,
                () -> envelope.markVerified(new DecodedRulesContainer(), ASSET_PAYLOAD));
        assertTrue(e.getMessage().contains("no signed payload"), e.getMessage());

        // And the gate stayed shut, so the getter still refuses.
        assertThrows(IllegalStateException.class, envelope::getWhitelistedAsset);
    }

    @Test
    @DisplayName("the getter refuses an envelope nothing marked")
    void unmarkedEnvelopeGetterRefuses() {
        assertThrows(IllegalStateException.class,
                () -> new SignedWhitelistedAssetEnvelope().getWhitelistedAsset());
    }

    @Test
    @DisplayName("the address envelope derives from its own payload too")
    void addressIsDerivedFromItsOwnPayload() throws Exception {
        SignedWhitelistedAddressEnvelope envelope = new SignedWhitelistedAddressEnvelope();
        WhitelistMetadata metadata = new WhitelistMetadata();
        String payload = "{\"address\":\"0xREAL\",\"label\":\"treasury\"}";
        metadata.setPayloadAsString(payload);
        envelope.setMetadata(metadata);

        envelope.markVerified(new DecodedRulesContainer(), payload);

        assertEquals("0xREAL", envelope.getWhitelistedAddress().getAddress());
    }

    @Test
    @DisplayName("an address envelope with no signed payload cannot be marked verified")
    void payloadLessAddressEnvelopeCannotBeMarked() {
        SignedWhitelistedAddressEnvelope envelope = new SignedWhitelistedAddressEnvelope();
        envelope.setMetadata(new WhitelistMetadata());

        assertThrows(WhitelistException.class,
                () -> envelope.markVerified(new DecodedRulesContainer(), "{}"));
        assertThrows(IllegalStateException.class, envelope::getWhitelistedAddress);
    }

    @Test
    @DisplayName("the derivation uses the VERIFIED payload, not the delivered one")
    void derivationUsesTheVerifiedPayloadNotTheDeliveredOne() throws Exception {
        // What a server can do on a legacy row: append a duplicate label immediately
        // before the closing brace. The legacy strip recovers the genuinely signed bytes,
        // so every signature check passes -- and Gson keeps the LAST of two duplicate
        // keys, so parsing the DELIVERED text would return the attacker's value as
        // verified. Step 4 hands over what it matched; that is what must be parsed.
        String signed = "{\"address\":\"0xREAL\",\"label\":\"treasury\"}";
        String delivered =
                "{\"address\":\"0xREAL\",\"label\":\"treasury\",\"label\":\"ATTACKER\"}";

        SignedWhitelistedAddressEnvelope envelope = new SignedWhitelistedAddressEnvelope();
        WhitelistMetadata metadata = new WhitelistMetadata();
        metadata.setPayloadAsString(delivered);
        envelope.setMetadata(metadata);

        envelope.markVerified(new DecodedRulesContainer(), signed);

        assertEquals("treasury", envelope.getWhitelistedAddress().getLabel(),
                "the appended label must not reach the caller as verified");
    }

    @Test
    @DisplayName("an absent verified payload is a hard failure, never a fallback")
    void anAbsentVerifiedPayloadDoesNotFallBackToMetadata() {
        // Defaulting to metadata.payloadAsString here would reintroduce the bug above,
        // which is why it is refused rather than defaulted.
        SignedWhitelistedAddressEnvelope envelope = new SignedWhitelistedAddressEnvelope();
        WhitelistMetadata metadata = new WhitelistMetadata();
        metadata.setPayloadAsString("{\"address\":\"0xREAL\"}");
        envelope.setMetadata(metadata);

        WhitelistException e = assertThrows(WhitelistException.class,
                () -> envelope.markVerified(new DecodedRulesContainer(), null));
        assertTrue(e.getMessage().contains("no verified payload"), e.getMessage());
        assertThrows(IllegalStateException.class, envelope::getWhitelistedAddress);
    }
}
