package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

import java.security.PublicKey;
import java.util.List;

/**
 * Represents the decoded rules container with all governance rules.
 * <p>
 * The rules container defines the complete governance structure for a Taurus Protect tenant,
 * including users, groups, signature thresholds, and rules for transactions and address
 * whitelisting. This class is populated by decoding the signed rules container from the API.
 * <p>
 * Key components:
 * <ul>
 *   <li><b>Users</b> - Individual users with their public keys and roles</li>
 *   <li><b>Groups</b> - Collections of users for group-based approvals</li>
 *   <li><b>Transaction Rules</b> - Rules governing transaction approvals</li>
 *   <li><b>Whitelisting Rules</b> - Rules for address and contract whitelisting</li>
 * </ul>
 *
 * @see RuleUser
 * @see RuleGroup
 * @see TransactionRules
 * @see AddressWhitelistingRules
 */
public class DecodedRulesContainer {

    private static final String HSMSLOT_ROLE = "HSMSLOT";
    private static final String ANY_WILDCARD = "Any";

    /**
     * List of users defined in the governance rules.
     */
    private List<RuleUser> users;

    /**
     * List of groups defined in the governance rules.
     */
    private List<RuleGroup> groups;

    /**
     * Minimum number of distinct user signatures required for rules container updates.
     */
    private int minimumDistinctUserSignatures;

    /**
     * Minimum number of distinct group signatures required for rules container updates.
     */
    private int minimumDistinctGroupSignatures;

    /**
     * Transaction approval rules organized by key (blockchain/action type).
     */
    private List<TransactionRules> transactionRules;

    /**
     * Address whitelisting rules organized by blockchain and network.
     */
    private List<AddressWhitelistingRules> addressWhitelistingRules;

    /**
     * Contract address whitelisting rules organized by blockchain and network.
     */
    private List<ContractAddressWhitelistingRules> contractAddressWhitelistingRules;

    /**
     * SHA-256 hash of the enforced rules for integrity verification.
     */
    private String enforcedRulesHash;

    /**
     * Unix timestamp when the rules container was created or updated.
     */
    private long timestamp;

    /**
     * Minimum number of commitment signatures required from HSM engines.
     */
    private int minimumCommitmentSignatures;

    /**
     * List of HSM engine identities (serial numbers) authorized for this tenant.
     */
    private List<String> engineIdentities;

    /**
     * Cached HSM public key (lazily resolved from users with HSMSLOT role).
     */
    private PublicKey hsmPublicKey;

    /**
     * Flag indicating whether the HSM public key has been resolved.
     */
    private boolean hsmPublicKeyResolved = false;

    /**
     * Lock object for thread-safe lazy initialization of HSM public key.
     */
    private final Object hsmKeyLock = new Object();

    private static boolean matches(String ruleValue, String searchValue) {
        if (ruleValue == null && searchValue == null) {
            return true;
        }
        if (ruleValue == null || searchValue == null) {
            return false;
        }
        return ruleValue.equals(searchValue);
    }

    /**
     * Checks if a value represents a wildcard (null, empty, or "Any").
     *
     * @param value the value to check
     * @return true if the value is a wildcard
     */
    private static boolean isWildcard(String value) {
        return value == null || value.isEmpty() || ANY_WILDCARD.equalsIgnoreCase(value);
    }

    /**
     * Checks if a rule is a global default (wildcard blockchain/currency).
     *
     * @param rule the rule to check
     * @return true if the rule has a wildcard currency
     */
    private static boolean isGlobalDefaultRule(AddressWhitelistingRules rule) {
        return isWildcard(rule.getCurrency());
    }

    /**
     * Checks if a rule has a wildcard network (matches any network).
     *
     * @param rule the rule to check
     * @return true if the rule has a wildcard network
     */
    private static boolean hasWildcardNetwork(AddressWhitelistingRules rule) {
        return isWildcard(rule.getNetwork());
    }

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the users in the rules.
     *
     * @return the users
     */
    public List<RuleUser> getUsers() {
        return users;
    }

    /**
     * Sets the users in the rules.
     *
     * @param users the users
     */
    public void setUsers(List<RuleUser> users) {
        this.users = users;
    }

    /**
     * Gets the groups in the rules.
     *
     * @return the groups
     */
    public List<RuleGroup> getGroups() {
        return groups;
    }

    /**
     * Sets the groups in the rules.
     *
     * @param groups the groups
     */
    public void setGroups(List<RuleGroup> groups) {
        this.groups = groups;
    }

    /**
     * Gets the minimum distinct user signatures required.
     *
     * @return the minimum distinct user signatures
     */
    public int getMinimumDistinctUserSignatures() {
        return minimumDistinctUserSignatures;
    }

    /**
     * Sets the minimum distinct user signatures required.
     *
     * @param minimumDistinctUserSignatures the minimum distinct user signatures
     */
    public void setMinimumDistinctUserSignatures(int minimumDistinctUserSignatures) {
        this.minimumDistinctUserSignatures = minimumDistinctUserSignatures;
    }

    /**
     * Gets the minimum distinct group signatures required.
     *
     * @return the minimum distinct group signatures
     */
    public int getMinimumDistinctGroupSignatures() {
        return minimumDistinctGroupSignatures;
    }

    /**
     * Sets the minimum distinct group signatures required.
     *
     * @param minimumDistinctGroupSignatures the minimum distinct group signatures
     */
    public void setMinimumDistinctGroupSignatures(int minimumDistinctGroupSignatures) {
        this.minimumDistinctGroupSignatures = minimumDistinctGroupSignatures;
    }

    /**
     * Gets the transaction rules.
     *
     * @return the transaction rules
     */
    public List<TransactionRules> getTransactionRules() {
        return transactionRules;
    }

    /**
     * Sets the transaction rules.
     *
     * @param transactionRules the transaction rules
     */
    public void setTransactionRules(List<TransactionRules> transactionRules) {
        this.transactionRules = transactionRules;
    }

    /**
     * Gets the address whitelisting rules.
     *
     * @return the address whitelisting rules
     */
    public List<AddressWhitelistingRules> getAddressWhitelistingRules() {
        return addressWhitelistingRules;
    }

    /**
     * Sets the address whitelisting rules.
     *
     * @param addressWhitelistingRules the address whitelisting rules
     */
    public void setAddressWhitelistingRules(List<AddressWhitelistingRules> addressWhitelistingRules) {
        this.addressWhitelistingRules = addressWhitelistingRules;
    }

    /**
     * Gets the contract address whitelisting rules.
     *
     * @return the contract address whitelisting rules
     */
    public List<ContractAddressWhitelistingRules> getContractAddressWhitelistingRules() {
        return contractAddressWhitelistingRules;
    }

    /**
     * Sets the contract address whitelisting rules.
     *
     * @param contractAddressWhitelistingRules the contract address whitelisting rules
     */
    public void setContractAddressWhitelistingRules(List<ContractAddressWhitelistingRules> contractAddressWhitelistingRules) {
        this.contractAddressWhitelistingRules = contractAddressWhitelistingRules;
    }

    /**
     * Gets the enforced rules hash.
     *
     * @return the enforced rules hash
     */
    public String getEnforcedRulesHash() {
        return enforcedRulesHash;
    }

    /**
     * Sets the enforced rules hash.
     *
     * @param enforcedRulesHash the enforced rules hash
     */
    public void setEnforcedRulesHash(String enforcedRulesHash) {
        this.enforcedRulesHash = enforcedRulesHash;
    }

    /**
     * Gets the timestamp.
     *
     * @return the timestamp
     */
    public long getTimestamp() {
        return timestamp;
    }

    /**
     * Sets the timestamp.
     *
     * @param timestamp the timestamp
     */
    public void setTimestamp(long timestamp) {
        this.timestamp = timestamp;
    }

    /**
     * Gets the minimum commitment signatures.
     *
     * @return the minimum commitment signatures
     */
    public int getMinimumCommitmentSignatures() {
        return minimumCommitmentSignatures;
    }

    /**
     * Sets the minimum commitment signatures.
     *
     * @param minimumCommitmentSignatures the minimum commitment signatures
     */
    public void setMinimumCommitmentSignatures(int minimumCommitmentSignatures) {
        this.minimumCommitmentSignatures = minimumCommitmentSignatures;
    }

    /**
     * Gets the engine identities (HSM serial numbers).
     *
     * @return the engine identities
     */
    public List<String> getEngineIdentities() {
        return engineIdentities;
    }

    /**
     * Sets the engine identities (HSM serial numbers).
     *
     * @param engineIdentities the engine identities
     */
    public void setEngineIdentities(List<String> engineIdentities) {
        this.engineIdentities = engineIdentities;
    }

    /**
     * Finds AddressWhitelistingRules matching the given blockchain and network.
     * Uses a three-tier priority system:
     * <ol>
     *   <li>Exact match - both blockchain and network match exactly</li>
     *   <li>Blockchain-only match - blockchain matches, rule has wildcard network (null/empty)</li>
     *   <li>Global default - rule has wildcard blockchain (null/empty/"Any")</li>
     * </ol>
     *
     * @param blockchain the blockchain (currency) identifier
     * @param network    the network identifier
     * @return the matching AddressWhitelistingRules, or null if not found
     */
    public AddressWhitelistingRules findAddressWhitelistingRules(String blockchain, String network) {
        if (addressWhitelistingRules == null) {
            return null;
        }

        AddressWhitelistingRules blockchainOnlyMatch = null;
        AddressWhitelistingRules globalDefault = null;

        for (AddressWhitelistingRules rule : addressWhitelistingRules) {
            boolean ruleIsGlobalDefault = isGlobalDefaultRule(rule);
            boolean blockchainMatches = !ruleIsGlobalDefault && matches(rule.getCurrency(), blockchain);
            boolean networkMatches = matches(rule.getNetwork(), network);
            boolean ruleHasWildcardNetwork = hasWildcardNetwork(rule);

            // Priority 1: Exact match (blockchain + network)
            if (blockchainMatches && networkMatches) {
                return rule;
            }

            // Priority 2: Blockchain match with wildcard network
            if (blockchainMatches && ruleHasWildcardNetwork && blockchainOnlyMatch == null) {
                blockchainOnlyMatch = rule;
            }

            // Priority 3: Global default (wildcard blockchain)
            if (ruleIsGlobalDefault && globalDefault == null) {
                globalDefault = rule;
            }
        }

        // Return best match by priority
        if (blockchainOnlyMatch != null) {
            return blockchainOnlyMatch;
        }
        return globalDefault;
    }

    /**
     * Finds ContractAddressWhitelistingRules matching the given blockchain and network.
     * Uses a three-tier priority system:
     * <ol>
     *   <li>Exact match - both blockchain and network match exactly</li>
     *   <li>Blockchain-only match - blockchain matches, rule has wildcard network (null/empty)</li>
     *   <li>Global default - rule has wildcard blockchain (null/empty/"Any")</li>
     * </ol>
     *
     * @param blockchain the blockchain identifier
     * @param network    the network identifier
     * @return the matching ContractAddressWhitelistingRules, or null if not found
     */
    public ContractAddressWhitelistingRules findContractAddressWhitelistingRules(
            String blockchain, String network) {
        if (contractAddressWhitelistingRules == null) {
            return null;
        }

        ContractAddressWhitelistingRules blockchainOnlyMatch = null;
        ContractAddressWhitelistingRules globalDefault = null;

        for (ContractAddressWhitelistingRules rule : contractAddressWhitelistingRules) {
            boolean ruleIsGlobalDefault = isWildcard(rule.getBlockchain());
            boolean blockchainMatches = !ruleIsGlobalDefault && matches(rule.getBlockchain(), blockchain);
            boolean networkMatches = matches(rule.getNetwork(), network);
            boolean ruleHasWildcardNetwork = isWildcard(rule.getNetwork());

            // Priority 1: Exact match (blockchain + network)
            if (blockchainMatches && networkMatches) {
                return rule;
            }

            // Priority 2: Blockchain match with wildcard network
            if (blockchainMatches && ruleHasWildcardNetwork && blockchainOnlyMatch == null) {
                blockchainOnlyMatch = rule;
            }

            // Priority 3: Global default (wildcard blockchain)
            if (ruleIsGlobalDefault && globalDefault == null) {
                globalDefault = rule;
            }
        }

        // Return best match by priority
        if (blockchainOnlyMatch != null) {
            return blockchainOnlyMatch;
        }
        return globalDefault;
    }

    /**
     * Finds a RuleUser by ID.
     *
     * @param userId the user ID to find
     * @return the RuleUser, or null if not found
     */
    public RuleUser findUserById(String userId) {
        if (users == null || userId == null) {
            return null;
        }
        for (RuleUser user : users) {
            if (userId.equals(user.getId())) {
                return user;
            }
        }
        return null;
    }

    /**
     * Finds a RuleGroup by ID.
     *
     * @param groupId the group ID to find
     * @return the RuleGroup, or null if not found
     */
    public RuleGroup findGroupById(String groupId) {
        if (groups == null || groupId == null) {
            return null;
        }
        for (RuleGroup group : groups) {
            if (groupId.equals(group.getId())) {
                return group;
            }
        }
        return null;
    }

    /**
     * Gets the HSM slot public key (cached).
     * Finds the first user with the HSMSLOT role.
     * This method is thread-safe.
     *
     * @return the HSM public key, or null if no user with HSMSLOT role exists
     */
    public PublicKey getHsmPublicKey() {
        synchronized (hsmKeyLock) {
            if (!hsmPublicKeyResolved) {
                hsmPublicKey = findHsmPublicKey();
                hsmPublicKeyResolved = true;
            }
            return hsmPublicKey;
        }
    }

    private PublicKey findHsmPublicKey() {
        if (users == null || users.isEmpty()) {
            return null;
        }
        for (RuleUser user : users) {
            if (user.getRoles() != null && user.getRoles().contains(HSMSLOT_ROLE)) {
                if (user.getPublicKey() != null) {
                    return user.getPublicKey();
                }
            }
        }
        return null;
    }
}
