package com.taurushq.sdk.protect.client.model;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Represents a whitelisted external address in the Taurus Protect system.
 * <p>
 * A whitelisted address is an approved destination for outgoing transfers.
 * Before funds can be sent to an external address, it must be added to the
 * whitelist and approved according to the configured governance rules.
 * <p>
 * Whitelisted addresses support various address types including:
 * <ul>
 *   <li>Individual addresses (personal wallets)</li>
 *   <li>Exchange addresses (centralized exchange accounts)</li>
 *   <li>Contract addresses (smart contracts)</li>
 *   <li>Validators and stake pools (for staking operations)</li>
 * </ul>
 *
 * @see WhitelistTrail
 * @see SignedWhitelistedAddress
 * @see InternalAddress
 */
public class WhitelistedAddress {

    /**
     * The blockchain identifier (e.g., "BTC", "ETH", "XTZ").
     */
    private String blockchain;

    /**
     * The type of address (individual, exchange, contract, etc.).
     */
    private AddressType addressType;

    /**
     * The actual blockchain address string.
     */
    private String address;

    /**
     * Optional memo field for Stellar or destination tag for Ripple.
     */
    private String memo;

    /**
     * Human-readable label or description for this whitelisted address.
     */
    private String label;

    /**
     * Customer ID for external reconciliation with other systems.
     */
    private String customerId;

    /**
     * Exchange account ID when the address belongs to an exchange.
     */
    private long exchangeAccountId;

    /**
     * List of internal addresses linked to this whitelisted address.
     */
    private List<InternalAddress> linkedInternalAddresses;

    /**
     * Smart contract type (e.g., "CMTA20", "ERC20") for contract addresses.
     */
    private String contractType;

    /**
     * List of internal wallets that can send to this whitelisted address.
     */
    private List<InternalWallet> linkedWallets;

    /**
     * The blockchain network (e.g., "mainnet", "testnet").
     */
    private String network;

    /**
     * Taurus Network participant ID for network-level identification.
     */
    private String tnParticipantID;

    /**
     * The creation timestamp extracted from the "created" trail action.
     */
    private java.time.OffsetDateTime createdAt;

    /**
     * Custom key-value attributes.
     */
    private Map<String, Object> attributes;

    /**
     * Returns the blockchain identifier.
     *
     * @return the blockchain identifier (e.g., "BTC", "ETH")
     */
    public String getBlockchain() {
        return blockchain;
    }

    /**
     * Sets the blockchain identifier.
     *
     * @param blockchain the blockchain identifier to set
     */
    public void setBlockchain(String blockchain) {
        this.blockchain = blockchain;
    }

    /**
     * Returns the type of this whitelisted address.
     *
     * @return the address type
     */
    public AddressType getAddressType() {
        return addressType;
    }

    /**
     * Sets the type of this whitelisted address.
     *
     * @param addressType the address type to set
     */
    public void setAddressType(AddressType addressType) {
        this.addressType = addressType;
    }

    /**
     * Returns the blockchain address string.
     *
     * @return the blockchain address
     */
    public String getAddress() {
        return address;
    }

    /**
     * Sets the blockchain address string.
     *
     * @param address the blockchain address to set
     */
    public void setAddress(String address) {
        this.address = address;
    }

    /**
     * Returns the optional memo (for Stellar) or destination tag (for Ripple).
     *
     * @return the memo or tag, or {@code null} if not applicable
     */
    public String getMemo() {
        return memo;
    }

    /**
     * Sets the optional memo (for Stellar) or destination tag (for Ripple).
     *
     * @param memo the memo or tag to set
     */
    public void setMemo(String memo) {
        this.memo = memo;
    }

    /**
     * Returns the human-readable label for this address.
     *
     * @return the address label
     */
    public String getLabel() {
        return label;
    }

    /**
     * Sets the human-readable label for this address.
     *
     * @param label the address label to set
     */
    public void setLabel(String label) {
        this.label = label;
    }

    /**
     * Returns the customer ID for external reconciliation.
     *
     * @return the customer ID
     */
    public String getCustomerId() {
        return customerId;
    }

    /**
     * Sets the customer ID for external reconciliation.
     *
     * @param customerId the customer ID to set
     */
    public void setCustomerId(String customerId) {
        this.customerId = customerId;
    }

    /**
     * Returns the exchange account ID.
     *
     * @return the exchange account ID, or 0 if not an exchange address
     */
    public long getExchangeAccountId() {
        return exchangeAccountId;
    }

    /**
     * Sets the exchange account ID.
     *
     * @param exchangeAccountId the exchange account ID to set
     */
    public void setExchangeAccountId(long exchangeAccountId) {
        this.exchangeAccountId = exchangeAccountId;
    }

    /**
     * Returns the list of internal addresses linked to this whitelisted address.
     *
     * @return the list of linked internal addresses
     */
    public List<InternalAddress> getLinkedInternalAddresses() {
        return linkedInternalAddresses;
    }

    /**
     * Sets the list of internal addresses linked to this whitelisted address.
     *
     * @param linkedInternalAddresses the list of linked internal addresses to set
     */
    public void setLinkedInternalAddresses(List<InternalAddress> linkedInternalAddresses) {
        this.linkedInternalAddresses = linkedInternalAddresses;
    }

    /**
     * Returns the smart contract type.
     *
     * @return the contract type (e.g., "CMTA20", "ERC20"), or {@code null} if not a contract
     */
    public String getContractType() {
        return contractType;
    }

    /**
     * Sets the smart contract type.
     *
     * @param contractType the contract type to set
     */
    public void setContractType(String contractType) {
        this.contractType = contractType;
    }

    /**
     * Returns the list of internal wallets that can send to this address.
     *
     * @return the list of linked wallets
     */
    public List<InternalWallet> getLinkedWallets() {
        return linkedWallets;
    }

    /**
     * Sets the list of internal wallets that can send to this address.
     *
     * @param linkedWallets the list of linked wallets to set
     */
    public void setLinkedWallets(List<InternalWallet> linkedWallets) {
        this.linkedWallets = linkedWallets;
    }

    /**
     * Returns the blockchain network.
     *
     * @return the network (e.g., "mainnet", "testnet")
     */
    public String getNetwork() {
        return network;
    }

    /**
     * Sets the blockchain network.
     *
     * @param network the network to set
     */
    public void setNetwork(String network) {
        this.network = network;
    }

    /**
     * Returns the Taurus Network participant ID.
     *
     * @return the participant ID
     */
    public String getTnParticipantID() {
        return tnParticipantID;
    }

    /**
     * Sets the Taurus Network participant ID.
     *
     * @param tnParticipantID the participant ID to set
     */
    public void setTnParticipantID(String tnParticipantID) {
        this.tnParticipantID = tnParticipantID;
    }

    /**
     * Returns the creation timestamp.
     *
     * @return the creation timestamp, or {@code null} if not available
     */
    public java.time.OffsetDateTime getCreatedAt() {
        return createdAt;
    }

    /**
     * Sets the creation timestamp.
     *
     * @param createdAt the creation timestamp to set
     */
    public void setCreatedAt(java.time.OffsetDateTime createdAt) {
        this.createdAt = createdAt;
    }

    /**
     * Returns the custom key-value attributes.
     *
     * @return the attributes map, or {@code null} if not available
     */
    public Map<String, Object> getAttributes() {
        return attributes;
    }

    /**
     * Sets the custom key-value attributes.
     *
     * @param attributes the attributes map to set
     */
    public void setAttributes(Map<String, Object> attributes) {
        this.attributes = attributes;
    }

    /**
     * Enumeration of supported whitelisted address types.
     * <p>
     * The address type determines the expected behavior and validation
     * rules for the whitelisted address.
     */
    public enum AddressType {
        /**
         * Individual personal wallet address.
         */
        individual(0),
        /**
         * Centralized exchange account address.
         */
        exchange(1),
        /**
         * Tezos baker address for delegation.
         */
        baker(2),
        /**
         * Smart contract address.
         */
        contract(3),
        /**
         * Cardano stake pool address for delegation.
         */
        stakepool(4),
        /**
         * Proof-of-stake validator address.
         */
        validator(5),
        /**
         * Blockchain node address.
         */
        node(6),
        /**
         * Fiat on/off ramp provider address.
         */
        fiatprovider(7),
        /**
         * Unknown or unrecognized address type.
         */
        UNRECOGNIZED(-1);

        /**
         * Map from numeric value to enum constant for efficient lookup.
         */
        private static final Map map = new HashMap<>();

        static {
            for (AddressType pageType : values()) {
                map.put(pageType.value, pageType);
            }
        }

        /**
         * The numeric value representing this address type.
         */
        private final int value;

        /**
         * Constructs an address type with the specified numeric value.
         *
         * @param value the numeric value for this type
         */
        AddressType(int value) {
            this.value = value;
        }

        /**
         * Returns the address type for the specified numeric value.
         *
         * @param value the numeric value to look up
         * @return the corresponding address type, or {@code null} if not found
         */
        public static AddressType valueOf(int value) {
            return (AddressType) map.get(value);
        }

        /**
         * Returns the numeric value of this address type.
         *
         * @return the numeric value
         */
        public int getValue() {
            return value;
        }
    }
}
