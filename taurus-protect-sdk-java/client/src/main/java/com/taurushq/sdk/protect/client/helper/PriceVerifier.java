package com.taurushq.sdk.protect.client.helper;

import com.google.gson.Gson;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.Price;
import com.taurushq.sdk.protect.client.model.PriceSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;

import java.nio.charset.StandardCharsets;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.List;

/**
 * Price signature verification against the PRICEUPDATER role.
 * <p>
 * Rate and decimals feed amount conversion, so an unverified price is a wrong number a
 * caller acts on.
 */
public final class PriceVerifier {

    private static final Gson GSON = new Gson();
    private static final String PRICE_UPDATER_ROLE = "PRICEUPDATER";

    private PriceVerifier() {
    }

    /**
     * Returns the exact bytes a price signature covers.
     * <p>
     * This is validatord's CurrencyPrice JSON projection: the five fields it serialises,
     * in its field order, with no spaces. The key order IS the signed byte sequence — do
     * not reorder or add fields.
     *
     * @param price the price
     * @return the canonical signed bytes
     */
    public static byte[] priceSignedBytes(final Price price) {
        if (price == null) {
            throw new IllegalArgumentException("price cannot be null");
        }
        String canonical = "{\"blockchain\":" + GSON.toJson(nullToEmpty(price.getBlockchain()))
                + ",\"currencyFrom\":" + GSON.toJson(nullToEmpty(price.getCurrencyFrom()))
                + ",\"currencyTo\":" + GSON.toJson(nullToEmpty(price.getCurrencyTo()))
                + ",\"decimals\":" + GSON.toJson(nullToEmpty(price.getDecimals()))
                + ",\"rate\":" + GSON.toJson(nullToEmpty(price.getRate()))
                + "}";
        return canonical.getBytes(StandardCharsets.UTF_8);
    }

    /**
     * Checks a price against the PRICEUPDATER keys in a verified rules container.
     * <p>
     * The container decides whether prices must be signed at all. It is
     * SuperAdmin-verified, so "this tenant has no price signer" is trustworthy, whereas
     * "this price carries no signatures" is not — which is why a stripped signatures list
     * on a tenant that DOES have a PRICEUPDATER is an error rather than a skip.
     *
     * @param price          the price to verify
     * @param rulesContainer the verified rules container
     * @throws IntegrityException if no PRICEUPDATER signature verifies
     */
    public static void verifyPrice(final Price price, final DecodedRulesContainer rulesContainer) {
        if (price == null) {
            throw new IntegrityException("price cannot be null");
        }
        if (rulesContainer == null) {
            throw new IntegrityException(
                    "rules container required for price signature verification");
        }

        List<PublicKey> keys = priceUpdaterKeys(rulesContainer);
        if (keys.isEmpty()) {
            // This tenant does not sign prices.
            return;
        }

        List<PriceSignature> signatures = price.getSignatures();
        if (signatures == null || signatures.isEmpty()) {
            throw new IntegrityException(String.format(
                    "price %s/%s carries no signatures but the rules container configures a %s",
                    price.getCurrencyFrom(), price.getCurrencyTo(), PRICE_UPDATER_ROLE));
        }

        byte[] data = priceSignedBytes(price);
        for (PriceSignature sig : signatures) {
            if (sig == null || sig.getSignature() == null || sig.getSignature().isEmpty()) {
                continue;
            }
            if (SignatureVerifier.isValidSignature(data, sig.getSignature(), keys)) {
                return;
            }
        }

        throw new IntegrityException(String.format(
                "no %s signature verifies for price %s/%s (%d signature(s) offered)",
                PRICE_UPDATER_ROLE, price.getCurrencyFrom(), price.getCurrencyTo(),
                signatures.size()));
    }

    /**
     * Verifies each price, throwing on the first failure.
     *
     * @param prices         the prices to verify
     * @param rulesContainer the verified rules container
     * @throws IntegrityException if any price fails
     */
    public static void verifyPrices(final List<Price> prices,
                                    final DecodedRulesContainer rulesContainer) {
        if (prices == null) {
            return;
        }
        for (Price price : prices) {
            if (price != null) {
                verifyPrice(price, rulesContainer);
            }
        }
    }

    /**
     * Public keys of every user carrying the PRICEUPDATER role.
     */
    private static List<PublicKey> priceUpdaterKeys(final DecodedRulesContainer rulesContainer) {
        List<PublicKey> keys = new ArrayList<>();
        if (rulesContainer.getUsers() == null) {
            return keys;
        }
        for (RuleUser user : rulesContainer.getUsers()) {
            if (user == null || user.getPublicKey() == null || user.getRoles() == null) {
                continue;
            }
            if (user.getRoles().contains(PRICE_UPDATER_ROLE)) {
                keys.add(user.getPublicKey());
            }
        }
        return keys;
    }

    private static String nullToEmpty(final String value) {
        return value == null ? "" : value;
    }
}
