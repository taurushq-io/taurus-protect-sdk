package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents a custom key-value attribute that can be attached to various entities.
 * <p>
 * Attributes provide a flexible way to store additional metadata on wallets, addresses,
 * requests, and other objects in the Taurus Protect system.
 * <p>
 * Example usage:
 * <pre>{@code
 * // Add an attribute to an address
 * client.getAddressService().createAddressAttribute(addressId, "department", "finance");
 *
 * // Read attributes from an address
 * Address address = client.getAddressService().getAddress(addressId);
 * for (Attribute attr : address.getAttributes()) {
 *     System.out.println(attr.getKey() + " = " + attr.getValue());
 * }
 * }</pre>
 *
 * @see Address
 * @see Wallet
 * @see Request
 */
public class Attribute {

    /**
     * The unique identifier of the attribute.
     */
    private long id;

    /**
     * The attribute key (name).
     */
    private String key;

    /**
     * The attribute value.
     */
    private String value;

    /**
     * The content type for file attributes (e.g., "application/pdf").
     */
    private String contentType;

    /**
     * The owner identifier of the attribute.
     */
    private String owner;

    /**
     * The type classification of the attribute.
     */
    private String type;

    /**
     * The sub-type classification of the attribute.
     */
    private String subType;

    /**
     * Indicates whether the attribute value is a file reference.
     */
    private boolean isFile;


    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets id.
     *
     * @return the id
     */
    public long getId() {
        return id;
    }

    /**
     * Sets id.
     *
     * @param id the id
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * Gets key.
     *
     * @return the key
     */
    public String getKey() {
        return key;
    }

    /**
     * Sets key.
     *
     * @param key the key
     */
    public void setKey(String key) {
        this.key = key;
    }

    /**
     * Gets value.
     *
     * @return the value
     */
    public String getValue() {
        return value;
    }

    /**
     * Sets value.
     *
     * @param value the value
     */
    public void setValue(String value) {
        this.value = value;
    }

    /**
     * Gets content type.
     *
     * @return the content type
     */
    public String getContentType() {
        return contentType;
    }

    /**
     * Sets content type.
     *
     * @param contentType the content type
     */
    public void setContentType(String contentType) {
        this.contentType = contentType;
    }

    /**
     * Gets owner.
     *
     * @return the owner
     */
    public String getOwner() {
        return owner;
    }

    /**
     * Sets owner.
     *
     * @param owner the owner
     */
    public void setOwner(String owner) {
        this.owner = owner;
    }

    /**
     * Gets type.
     *
     * @return the type
     */
    public String getType() {
        return type;
    }

    /**
     * Sets type.
     *
     * @param type the type
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Gets sub type.
     *
     * @return the sub type
     */
    public String getSubType() {
        return subType;
    }

    /**
     * Sets sub type.
     *
     * @param subType the sub type
     */
    public void setSubType(String subType) {
        this.subType = subType;
    }

    /**
     * Is file boolean.
     *
     * @return the boolean
     */
    public boolean isFile() {
        return isFile;
    }

    /**
     * Sets is file.
     *
     * @param file the file
     */
    public void setIsFile(boolean file) {
        isFile = file;
    }
}
