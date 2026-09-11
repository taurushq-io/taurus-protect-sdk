package com.taurushq.sdk.protect.client.helper;

import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAddress;
import com.taurushq.sdk.protect.client.model.WhitelistedAsset;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * A duplicate key in a signed payload is refused, in all three places the SDK parses one.
 *
 * <p><b>Why this is a security gate and not JSON pedantry.</b> Gson keeps the LAST of two
 * duplicate keys. The legacy-hash tolerance accepts three backward-compatible rewrites of
 * the delivered payload, and those strips are not injective: a server can append a
 * duplicate {@code ,"label":"X"} immediately before the closing brace, the strip recovers
 * the genuinely signed bytes, every signature check from step 1 to step 5 passes — and a
 * parser reading the DELIVERED text hands the appended value back as verified. Parsing the
 * MATCHED variant closes the shapes the strip reaches; this closes the rest, including the
 * third parse in {@link WhitelistHashHelper#resolveRuleKey(String, String, String)} where a
 * duplicated {@code currency} or {@code network} would re-point rule selection at a weaker
 * tier.
 *
 * <p>It must be a separate structural pass: no Gson setting rejects duplicates, and by the
 * time a {@code JsonObject} exists they have already been collapsed.
 *
 * <p>The exception is the CHECKED {@link WhitelistException}, not the unchecked
 * {@code IntegrityException}: one bad row is a row-level failure, which is the class the
 * whitelist row loops catch and exclude on. An unchecked throw would escape those loops and
 * turn one unparseable row into an aborted listing.
 */
class DuplicateJsonKeyRejectionTest {

    private static final String GOOD_ADDRESS =
            "{\"address\":\"0xREAL\",\"label\":\"treasury\",\"currency\":\"ETH\"}";
    private static final String GOOD_ASSET =
            "{\"name\":\"USDC\",\"symbol\":\"USDC\",\"contractAddress\":\"0xREAL\","
                    + "\"blockchain\":\"ETH\"}";

    @Test
    @DisplayName("a clean payload still parses — the check is not a blanket refusal")
    void aCleanPayloadStillParses() throws Exception {
        WhitelistedAddress address =
                WhitelistHashHelper.parseWhitelistedAddressFromJson(GOOD_ADDRESS);
        assertEquals("0xREAL", address.getAddress());
        assertEquals("treasury", address.getLabel());

        WhitelistedAsset asset = AssetHashHelper.parseWhitelistedAssetFromJson(GOOD_ASSET);
        assertEquals("USDC", asset.getSymbol());
    }

    @Test
    @DisplayName("an appended duplicate label is refused, not last-wins")
    void appendedDuplicateLabelIsRefused() {
        String injected =
                "{\"address\":\"0xREAL\",\"label\":\"treasury\",\"label\":\"ATTACKER\"}";

        WhitelistException e = assertThrows(WhitelistException.class,
                () -> WhitelistHashHelper.parseWhitelistedAddressFromJson(injected));
        assertTrue(e.getMessage().contains("duplicate key"), e.getMessage());
        assertTrue(e.getMessage().contains("label"), e.getMessage());
    }

    @Test
    @DisplayName("a duplicate address is refused — the field that names the destination")
    void duplicateAddressIsRefused() {
        String injected = "{\"address\":\"0xREAL\",\"address\":\"0xATTACKER\"}";

        assertThrows(WhitelistException.class,
                () -> WhitelistHashHelper.parseWhitelistedAddressFromJson(injected));
    }

    @Test
    @DisplayName("a duplicate NESTED key is refused too")
    void duplicateNestedKeyIsRefused() {
        // The strip removes inner linkedInternalAddresses labels wholesale, so an
        // injection inside one of those objects is reachable on any strategy-2-era row.
        String injected = "{\"address\":\"0xREAL\",\"linkedInternalAddresses\":"
                + "[{\"id\":\"1\",\"label\":\"ops\",\"label\":\"ATTACKER\"}]}";

        assertThrows(WhitelistException.class,
                () -> WhitelistHashHelper.parseWhitelistedAddressFromJson(injected));
    }

    @Test
    @DisplayName("siblings in DIFFERENT objects may share a key")
    void siblingObjectsMaySharaAKey() throws Exception {
        // Per-object key sets, not one global set: two array elements each carrying
        // "label" is ordinary data, and rejecting it would refuse real payloads.
        String legitimate = "{\"address\":\"0xREAL\",\"linkedWallets\":"
                + "[{\"label\":\"a\"},{\"label\":\"b\"}]}";

        WhitelistedAddress address =
                WhitelistHashHelper.parseWhitelistedAddressFromJson(legitimate);
        assertEquals("0xREAL", address.getAddress());
    }

    @Test
    @DisplayName("a duplicate contractAddress is refused on the asset side")
    void duplicateContractAddressIsRefused() {
        String injected = "{\"name\":\"USDC\",\"contractAddress\":\"0xREAL\","
                + "\"contractAddress\":\"0xATTACKER\"}";

        WhitelistException e = assertThrows(WhitelistException.class,
                () -> AssetHashHelper.parseWhitelistedAssetFromJson(injected));
        assertTrue(e.getMessage().contains("duplicate key"), e.getMessage());
    }

    @Test
    @DisplayName("resolveRuleKey refuses a duplicated currency — the rule-tier selector")
    void resolveRuleKeyRefusesADuplicatedCurrency() {
        // This is the third parse, and the one with the widest consequence: currency and
        // network decide WHICH governance rule judges the row, so last-wins here could
        // steer an address to a tier with a weaker quorum.
        String injected = "{\"address\":\"0xREAL\",\"currency\":\"ETH\",\"currency\":\"XTZ\"}";

        WhitelistException e = assertThrows(WhitelistException.class,
                () -> WhitelistHashHelper.resolveRuleKey(injected, "ETH", "mainnet"));
        assertTrue(e.getMessage().contains("duplicate key"), e.getMessage());
    }

    @Test
    @DisplayName("trailing content after the document is refused")
    void trailingContentIsRefused() {
        // Lenient JSON absorbs trailing garbage, which would mean the document a
        // signature covers is not the document that was parsed.
        String trailing = GOOD_ADDRESS + "{\"address\":\"0xATTACKER\"}";

        WhitelistException e = assertThrows(WhitelistException.class,
                () -> WhitelistHashHelper.parseWhitelistedAddressFromJson(trailing));
        assertTrue(e.getMessage().contains("trailing content")
                || e.getMessage().contains("not well-formed"), e.getMessage());
    }

    @Test
    @DisplayName("excessive nesting is refused rather than recursed into")
    void excessiveNestingIsRefused() {
        StringBuilder deep = new StringBuilder("{\"a\":");
        for (int i = 0; i < 200; i++) {
            deep.append("[");
        }
        deep.append("1");
        for (int i = 0; i < 200; i++) {
            deep.append("]");
        }
        deep.append("}");

        assertThrows(WhitelistException.class,
                () -> WhitelistHashHelper.parseWhitelistedAddressFromJson(deep.toString()));
    }
}
