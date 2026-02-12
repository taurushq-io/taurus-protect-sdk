package com.taurushq.sdk.protect.client.model;

import java.time.OffsetDateTime;

/**
 * Represents a tag in the Taurus Protect system.
 * <p>
 * Tags are used to label and categorize various resources such as
 * wallets, addresses, and transactions for easier organization and filtering.
 *
 * @see TagService
 */
public class Tag {

    private String id;
    private String value;
    private OffsetDateTime creationDate;
    private String color;

    /**
     * Gets the unique identifier of the tag.
     *
     * @return the tag ID
     */
    public String getId() {
        return id;
    }

    /**
     * Sets the tag ID.
     *
     * @param id the ID to set
     */
    public void setId(String id) {
        this.id = id;
    }

    /**
     * Gets the tag value (the text/label of the tag).
     *
     * @return the tag value
     */
    public String getValue() {
        return value;
    }

    /**
     * Sets the tag value.
     *
     * @param value the value to set
     */
    public void setValue(String value) {
        this.value = value;
    }

    /**
     * Gets the creation date of the tag.
     *
     * @return the creation date
     */
    public OffsetDateTime getCreationDate() {
        return creationDate;
    }

    /**
     * Sets the creation date.
     *
     * @param creationDate the creation date to set
     */
    public void setCreationDate(OffsetDateTime creationDate) {
        this.creationDate = creationDate;
    }

    /**
     * Gets the color associated with the tag (for UI display).
     *
     * @return the color (e.g., hex code)
     */
    public String getColor() {
        return color;
    }

    /**
     * Sets the tag color.
     *
     * @param color the color to set
     */
    public void setColor(String color) {
        this.color = color;
    }
}
