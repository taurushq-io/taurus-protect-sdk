package com.taurushq.sdk.protect.client.helper;

import com.google.common.base.Strings;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAsset;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;

/**
 * Helper class for whitelisted asset (contract address) parsing operations.
 */
public final class AssetHashHelper {

    /** Matches {@code ,"isNFT":true|false} — the member with its LEADING comma. */
    private static final Pattern IS_NFT_LEADING_COMMA_PATTERN =
            Pattern.compile(",\"isNFT\":(true|false)");

    /** Matches {@code "isNFT":true|false,} — the member with its TRAILING comma. */
    private static final Pattern IS_NFT_TRAILING_COMMA_PATTERN =
            Pattern.compile("\"isNFT\":(true|false),");

    /** Matches {@code ,"kindType":"..."} — the member with its LEADING comma. */
    private static final Pattern KIND_TYPE_LEADING_COMMA_PATTERN =
            Pattern.compile(",\"kindType\":\"[^\"]*\"");

    /** Matches {@code "kindType":"...",} — the member with its TRAILING comma. */
    private static final Pattern KIND_TYPE_TRAILING_COMMA_PATTERN =
            Pattern.compile("\"kindType\":\"[^\"]*\",");

    private AssetHashHelper() {
        // Prevent instantiation
    }

    /**
     * Computes the backward-compatible payload rewrites for an ASSET, each paired with its
     * hash. This is the asset peer of
     * {@link WhitelistHashHelper#computeLegacyPayloadVariants(String)}.
     *
     * <p>Strategies, in order:
     * <ol>
     *   <li>remove {@code isNFT} (signed before {@code isNFT} existed)</li>
     *   <li>remove {@code kindType} (signed before {@code kindType} existed)</li>
     *   <li>remove both</li>
     * </ol>
     *
     * <p>No asset identity field is currently injectable through these strips — the
     * stripped members ({@code isNFT}, {@code kindType}) are not read by
     * {@link #parseWhitelistedAssetFromJson(String)}, so for a variant to hash-match the
     * inserted text has to be exactly what the regexes remove. The payload is carried
     * anyway, so the two flows stay symmetric and a future schema change that makes a
     * stripped field readable does not silently reopen the address defect here.
     *
     * @param payloadAsString the payload exactly as delivered, may be null
     * @return the variants to try, in strategy order, deduplicated by hash; never null
     */
    public static List<LegacyPayloadVariant> computeAssetLegacyPayloadVariants(
            final String payloadAsString) {
        if (payloadAsString == null) {
            return new ArrayList<>();
        }

        Map<String, LegacyPayloadVariant> unique = new LinkedHashMap<>();

        // Strategy 1: remove isNFT only.
        String withoutIsNFT = stripIsNFT(payloadAsString);
        WhitelistHashHelper.addVariant(unique, payloadAsString, withoutIsNFT);

        // Strategy 2: remove kindType only.
        String withoutKindType = stripKindType(payloadAsString);
        WhitelistHashHelper.addVariant(unique, payloadAsString, withoutKindType);

        // Strategy 3: remove both. isNFT first, then kindType — the fixed order is part
        // of the wire contract the cross-SDK legacy-hash vectors pin.
        String withoutBoth = stripKindType(stripIsNFT(payloadAsString));
        WhitelistHashHelper.addVariant(unique, payloadAsString, withoutBoth);

        return new ArrayList<>(unique.values());
    }

    /**
     * Computes alternative ASSET hashes for backward compatibility, discarding the payload
     * each one came from.
     *
     * <p>Verification uses {@link #computeAssetLegacyPayloadVariants(String)}: step 6
     * needs the payload, not only its hash. This projection remains for the cross-SDK
     * vector oracle, which asserts hashes and counts.
     *
     * @param payloadAsString the payload exactly as delivered, may be null
     * @return the legacy hashes to try, in strategy order; never null
     */
    public static List<String> computeAssetLegacyHashes(final String payloadAsString) {
        return WhitelistHashHelper.projectHashes(
                computeAssetLegacyPayloadVariants(payloadAsString));
    }

    private static String stripIsNFT(final String payload) {
        String stripped = IS_NFT_LEADING_COMMA_PATTERN.matcher(payload).replaceAll("");
        return IS_NFT_TRAILING_COMMA_PATTERN.matcher(stripped).replaceAll("");
    }

    private static String stripKindType(final String payload) {
        String stripped = KIND_TYPE_LEADING_COMMA_PATTERN.matcher(payload).replaceAll("");
        return KIND_TYPE_TRAILING_COMMA_PATTERN.matcher(stripped).replaceAll("");
    }

    /**
     * Parses a WhitelistedAsset from the verified JSON payload string.
     * The JSON format contains the signed fields (blockchain, symbol, contractAddress, etc.)
     * that have been cryptographically verified.
     *
     * <p>Callers must pass the payload the matched signature COVERED — see
     * {@link LegacyPayloadVariant} — not whatever the response delivered.
     *
     * @param json the verified JSON payload string
     * @return the WhitelistedAsset populated from JSON fields
     * @throws WhitelistException if parsing fails
     */
    public static WhitelistedAsset parseWhitelistedAssetFromJson(String json) throws WhitelistException {
        if (Strings.isNullOrEmpty(json)) {
            throw new WhitelistException("JSON payload cannot be null or empty");
        }
        // Bounded and duplicate-checked here rather than at the call sites, for the same
        // reason as the address parser: this is public, and a new caller must not be able
        // to forget either bound. Gson is last-duplicate-wins, which on a signed payload
        // means the parser can be steered to a value no signature covered.
        if (json.length() > WhitelistHashHelper.MAX_PAYLOAD_BYTES) {
            throw new WhitelistException(String.format(
                    "cannot parse whitelisted asset: payload exceeds %d bytes",
                    WhitelistHashHelper.MAX_PAYLOAD_BYTES));
        }
        WhitelistHashHelper.rejectDuplicateObjectKeys(json);
        try {
            JsonObject obj = JsonParser.parseString(json).getAsJsonObject();

            WhitelistedAsset asset = new WhitelistedAsset();
            asset.setBlockchain(getStringOrNull(obj, "blockchain"));
            asset.setSymbol(getStringOrNull(obj, "symbol"));
            asset.setContractAddress(getStringOrNull(obj, "contractAddress"));
            asset.setName(getStringOrNull(obj, "name"));
            asset.setNetwork(getStringOrNull(obj, "network"));
            asset.setTokenId(getStringOrNull(obj, "tokenId"));

            // Parse decimals (integer)
            if (obj.has("decimals") && !obj.get("decimals").isJsonNull()) {
                asset.setDecimals(obj.get("decimals").getAsInt());
            }

            // Parse kind enum
            String kindStr = getStringOrNull(obj, "kind");
            asset.setKind(WhitelistedAsset.AssetKind.fromString(kindStr));

            return asset;
        } catch (Exception e) {
            throw new WhitelistException("Failed to parse WhitelistedAsset from JSON: " + e.getMessage(), e);
        }
    }

    private static String getStringOrNull(JsonObject obj, String key) {
        if (!obj.has(key) || obj.get(key).isJsonNull()) {
            return null;
        }
        String value = obj.get(key).getAsString();
        return value.isEmpty() ? null : value;
    }
}
