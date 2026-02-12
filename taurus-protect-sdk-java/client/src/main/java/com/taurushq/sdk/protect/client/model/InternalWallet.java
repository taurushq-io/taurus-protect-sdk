package com.taurushq.sdk.protect.client.model;

import com.google.gson.annotations.SerializedName;

/**
 * Represents an internal wallet linked to a whitelisted address.
 * <p>
 * Internal wallets are wallets managed within Taurus Protect that are
 * authorized to send funds to specific whitelisted addresses.
 *
 * @see WhitelistedAddress
 * @see Wallet
 */
public class InternalWallet {

    /**
     * Unique identifier for the internal wallet.
     */
    @SerializedName("id")
    private long id;

    /**
     * Human-readable name for the wallet.
     */
    @SerializedName("name")
    private String name;

    /**
     * Hierarchical path or location of the wallet in the organization structure.
     */
    @SerializedName("path")
    private String path;

    /**
     * Returns the unique identifier for this internal wallet.
     *
     * @return the wallet ID
     */
    public long getId() {
        return id;
    }

    /**
     * Sets the unique identifier for this internal wallet.
     *
     * @param id the wallet ID to set
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * Returns the human-readable name for this wallet.
     *
     * @return the wallet name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the human-readable name for this wallet.
     *
     * @param name the wallet name to set
     */
    public void setName(String name) {
        this.name = name;
    }

    /**
     * Returns the hierarchical path of this wallet.
     *
     * @return the wallet path
     */
    public String getPath() {
        return path;
    }

    /**
     * Sets the hierarchical path of this wallet.
     *
     * @param path the wallet path to set
     */
    public void setPath(String path) {
        this.path = path;
    }
}
