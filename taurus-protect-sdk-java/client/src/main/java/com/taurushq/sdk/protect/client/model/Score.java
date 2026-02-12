package com.taurushq.sdk.protect.client.model;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.time.OffsetDateTime;

/**
 * Represents a compliance score from a risk assessment provider.
 * <p>
 * Scores are provided by third-party compliance services such as Chainalysis,
 * Coinfirm, Elliptic, or Scorechain. These scores help assess the risk associated
 * with blockchain addresses for anti-money laundering (AML) and compliance purposes.
 * <p>
 * Different providers use different scoring methodologies and scales. The score
 * value should be interpreted according to the specific provider's documentation.
 * <p>
 * Example usage:
 * <pre>{@code
 * Address address = client.getAddressService().getAddress(addressId);
 * for (Score score : address.getScores()) {
 *     System.out.println("Provider: " + score.getProvider());
 *     System.out.println("Score: " + score.getScore());
 *     System.out.println("Type: " + score.getType());
 * }
 * }</pre>
 *
 * @see Address
 * @see WhitelistedAddress
 */
public class Score {

    /**
     * The unique identifier of the score record.
     */
    private long id;

    /**
     * The name of the risk assessment provider (e.g., "chainalysis", "elliptic").
     */
    private String provider;

    /**
     * The type of score (e.g., "in" for incoming, "out" for outgoing).
     */
    private String type;

    /**
     * The score value as assigned by the provider.
     */
    private String score;

    /**
     * The date and time when the score was last updated.
     */
    private OffsetDateTime updateDate;


    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Returns the unique identifier of the score record.
     *
     * @return the score ID
     */
    public long getId() {
        return id;
    }

    /**
     * Sets the unique identifier of the score record.
     *
     * @param id the score ID to set
     */
    public void setId(long id) {
        this.id = id;
    }

    /**
     * Returns the name of the risk assessment provider.
     * <p>
     * Common providers include "chainalysis", "coinfirm", "elliptic", and "scorechain".
     *
     * @return the provider name
     */
    public String getProvider() {
        return provider;
    }

    /**
     * Sets the name of the risk assessment provider.
     *
     * @param provider the provider name to set
     */
    public void setProvider(String provider) {
        this.provider = provider;
    }

    /**
     * Returns the type of score.
     * <p>
     * Common types include "in" (for incoming transaction risk) and "out" (for outgoing).
     *
     * @return the score type
     */
    public String getType() {
        return type;
    }

    /**
     * Sets the type of score.
     *
     * @param type the score type to set
     */
    public void setType(String type) {
        this.type = type;
    }

    /**
     * Returns the score value as assigned by the provider.
     * <p>
     * The interpretation of this value depends on the provider's scoring methodology.
     *
     * @return the score value
     */
    public String getScore() {
        return score;
    }

    /**
     * Sets the score value.
     *
     * @param score the score value to set
     */
    public void setScore(String score) {
        this.score = score;
    }

    /**
     * Returns the date and time when the score was last updated.
     *
     * @return the update date
     */
    public OffsetDateTime getUpdateDate() {
        return updateDate;
    }

    /**
     * Sets the date and time when the score was last updated.
     *
     * @param updateDate the update date to set
     */
    public void setUpdateDate(OffsetDateTime updateDate) {
        this.updateDate = updateDate;
    }
}
