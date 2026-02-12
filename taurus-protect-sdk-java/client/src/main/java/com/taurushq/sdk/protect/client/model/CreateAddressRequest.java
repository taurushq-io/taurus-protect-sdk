package com.taurushq.sdk.protect.client.model;

import static com.google.common.base.Preconditions.checkArgument;
import static com.google.common.base.Strings.isNullOrEmpty;

/**
 * Request object for creating a new address within a wallet.
 * <p>
 * Use the builder pattern for a cleaner API:
 * <pre>{@code
 * CreateAddressRequest request = CreateAddressRequest.builder()
 *     .walletId(123)
 *     .label("Deposit Address")
 *     .comment("Customer deposit address")
 *     .customerId("customer-456")
 *     .build();
 *
 * Address address = client.getAddressService().createAddress(request);
 * }</pre>
 */
public final class CreateAddressRequest {

    /**
     * The ID of the wallet where the address will be created.
     */
    private final long walletId;

    /**
     * The human-readable label for the address.
     */
    private final String label;

    /**
     * A comment or description for the address.
     */
    private final String comment;

    /**
     * An optional customer identifier associated with the address.
     */
    private final String customerId;

    private CreateAddressRequest(Builder builder) {
        this.walletId = builder.walletId;
        this.label = builder.label;
        this.comment = builder.comment;
        this.customerId = builder.customerId;
    }

    /**
     * Creates a new builder for CreateAddressRequest.
     *
     * @return a new builder
     */
    public static Builder builder() {
        return new Builder();
    }

    /**
     * Gets the wallet ID where the address will be created.
     *
     * @return the wallet ID
     */
    public long getWalletId() {
        return walletId;
    }

    /**
     * Gets the address label.
     *
     * @return the label
     */
    public String getLabel() {
        return label;
    }

    /**
     * Gets the address comment.
     *
     * @return the comment
     */
    public String getComment() {
        return comment;
    }

    /**
     * Gets the customer ID associated with this address.
     *
     * @return the customer ID, or empty string if not set
     */
    public String getCustomerId() {
        return customerId;
    }

    /**
     * Builder for constructing {@link CreateAddressRequest} instances.
     * <p>
     * Required fields:
     * <ul>
     *   <li>{@link #walletId(long)} - the wallet ID</li>
     *   <li>{@link #label(String)} - the address label</li>
     *   <li>{@link #comment(String)} - a comment for the address</li>
     * </ul>
     */
    public static final class Builder {

        /**
         * The wallet ID.
         */
        private long walletId;

        /**
         * The address label.
         */
        private String label;

        /**
         * The address comment.
         */
        private String comment;

        /**
         * The optional customer ID.
         */
        private String customerId = "";

        private Builder() {
        }

        /**
         * Sets the wallet ID where the address will be created (required).
         *
         * @param walletId the wallet ID
         * @return this builder
         */
        public Builder walletId(long walletId) {
            this.walletId = walletId;
            return this;
        }

        /**
         * Sets the address label (required).
         * <p>
         * The label is a human-readable identifier for the address.
         *
         * @param label the label
         * @return this builder
         */
        public Builder label(String label) {
            this.label = label;
            return this;
        }

        /**
         * Sets a comment for the address (required).
         *
         * @param comment the comment
         * @return this builder
         */
        public Builder comment(String comment) {
            this.comment = comment;
            return this;
        }

        /**
         * Sets an optional customer ID for the address.
         *
         * @param customerId the customer ID
         * @return this builder
         */
        public Builder customerId(String customerId) {
            this.customerId = customerId != null ? customerId : "";
            return this;
        }

        /**
         * Builds the CreateAddressRequest.
         *
         * @return the request object
         * @throws IllegalArgumentException if required fields are missing
         */
        public CreateAddressRequest build() {
            checkArgument(walletId > 0, "walletId must be positive");
            checkArgument(!isNullOrEmpty(label), "label is required");
            checkArgument(!isNullOrEmpty(comment), "comment is required");
            return new CreateAddressRequest(this);
        }
    }
}
