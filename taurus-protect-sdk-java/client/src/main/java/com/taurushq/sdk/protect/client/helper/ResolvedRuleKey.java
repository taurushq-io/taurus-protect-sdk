package com.taurushq.sdk.protect.client.helper;

/**
 * The {@code (blockchain, network)} pair that selects the governance rules for an entity,
 * plus the one fact a caller cannot otherwise recover: whether the NETWORK came from the
 * signed payload or was taken from the unsigned response DTO.
 *
 * <p>Step 5 needs that distinction because the network selects which rule — and therefore
 * which group quorum — judges the row. When the network is unsigned, a single rule lookup
 * lets the server pick the quorum; see
 * {@link com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer#findAddressWhitelistingRuleCandidates(String)}
 * for what to do instead.
 *
 * <p>{@link WhitelistHashHelper#resolveRuleKey(String, String, String)} keeps its
 * two-element {@code String[]} shape because the shared {@code rule_key} vectors in
 * {@code scripts/resources/verification-behaviour-vectors.json} assert exactly that shape
 * across all four SDKs.
 */
public final class ResolvedRuleKey {

    private final String blockchain;
    private final String network;
    private final boolean networkFromPayload;

    /**
     * Constructs a resolved key.
     *
     * @param blockchain         the chain, always from the signed payload
     * @param network            the network, from the payload when it carries one
     * @param networkFromPayload whether the network above was signed
     */
    ResolvedRuleKey(final String blockchain, final String network,
                    final boolean networkFromPayload) {
        this.blockchain = blockchain;
        this.network = network;
        this.networkFromPayload = networkFromPayload;
    }

    /**
     * Returns the chain. Always sourced from the signed payload — an absent chain is an
     * error rather than a wildcard, because an empty value matches every rule tier.
     *
     * @return the blockchain
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Returns the network, which may have come from the unsigned DTO.
     *
     * @return the network, possibly null
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Returns whether {@link #getNetwork()} was carried by the signed payload.
     *
     * <p>False means governance had {@code includeNetworkInPayload} off for this rule, so
     * the value came from the response and the caller must enforce every tier that value
     * could have selected instead of trusting the one it named.
     *
     * @return true when the network was signed
     */
    public boolean isNetworkFromPayload() {
        return networkFromPayload;
    }

    @Override
    public String toString() {
        return "ResolvedRuleKey{blockchain=" + blockchain + ", network=" + network
                + ", networkFromPayload=" + networkFromPayload + "}";
    }
}
