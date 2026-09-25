package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * An entity registered with a fiat provider.
 *
 * @see com.taurushq.sdk.protect.client.service.FiatService#listFiatProviderEntities
 */
public class FiatProviderEntity {

    private String id;
    private String provider;
    private String label;
    private String accountIdentifier;
    private String name;
    private String details;
    private OffsetDateTime creationDate;
    private OffsetDateTime updateDate;

    /**
     * Gets the entity id.
     *
     * @return the entity id
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the entity id.
     *
     * @param id the entity id
     */
    public void setId(final String id) {
        this.id = id;
    }

    /**
     * Gets the provider.
     *
     * @return the provider
     */
    public String getProvider() {
        return provider;
    }

    /**
     * Sets the provider.
     *
     * @param provider the provider
     */
    public void setProvider(final String provider) {
        this.provider = provider;
    }

    /**
     * Gets the label.
     *
     * @return the label
     */
    public String getLabel() {
        return label;
    }

    /**
     * Sets the label.
     *
     * @param label the label
     */
    public void setLabel(final String label) {
        this.label = label;
    }

    /**
     * Gets the account identifier.
     *
     * @return the account identifier
     */
    public String getAccountIdentifier() {
        return accountIdentifier;
    }

    /**
     * Sets the account identifier.
     *
     * @param accountIdentifier the account identifier
     */
    public void setAccountIdentifier(final String accountIdentifier) {
        this.accountIdentifier = accountIdentifier;
    }

    /**
     * Gets the name.
     *
     * @return the name
     */
    public String getName() {
        return name;
    }

    /**
     * Sets the name.
     *
     * @param name the name
     */
    public void setName(final String name) {
        this.name = name;
    }

    /**
     * Gets the details.
     *
     * @return the details
     */
    public String getDetails() {
        return details;
    }

    /**
     * Sets the details.
     *
     * @param details the details
     */
    public void setDetails(final String details) {
        this.details = details;
    }

    /**
     * Gets the creation date.
     *
     * @return the creation date
     */
    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    /**
     * Sets the creation date.
     *
     * @param creationDate the creation date
     */
    public void setCreationDate(final OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    /**
     * Gets the last update date.
     *
     * @return the last update date
     */
    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    /**
     * Sets the last update date.
     *
     * @param updateDate the last update date
     */
    public void setUpdateDate(final OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }
}
