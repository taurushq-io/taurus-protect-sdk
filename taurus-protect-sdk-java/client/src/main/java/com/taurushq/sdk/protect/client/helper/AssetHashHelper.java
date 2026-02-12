package com.taurushq.sdk.protect.client.helper;

import com.google.common.base.Strings;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.model.WhitelistException;
import com.taurushq.sdk.protect.client.model.WhitelistedAsset;

/**
 * Helper class for whitelisted asset (contract address) parsing operations.
 */
public final class AssetHashHelper {

    private AssetHashHelper() {
        // Prevent instantiation
    }

    /**
     * Parses a WhitelistedAsset from the verified JSON payload string.
     * The JSON format contains the signed fields (blockchain, symbol, contractAddress, etc.)
     * that have been cryptographically verified.
     *
     * @param json the verified JSON payload string
     * @return the WhitelistedAsset populated from JSON fields
     * @throws WhitelistException if parsing fails
     */
    public static WhitelistedAsset parseWhitelistedAssetFromJson(String json) throws WhitelistException {
        if (Strings.isNullOrEmpty(json)) {
            throw new WhitelistException("JSON payload cannot be null or empty");
        }
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
