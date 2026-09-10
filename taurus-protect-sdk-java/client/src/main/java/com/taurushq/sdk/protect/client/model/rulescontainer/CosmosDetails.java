package com.taurushq.sdk.protect.client.model.rulescontainer;

import java.util.List;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents Cosmos-specific scoping for a transaction rule.
 * <p>
 * A node rather than a flattened list on {@link TransactionRuleDetails}, so protobuf
 * fields from a newer schema can be preserved on this sub-message and re-emitted on
 * encode — matching the Go, Python and TypeScript SDKs.
 *
 * @see TransactionRuleDetails
 */
public class CosmosDetails extends RulesNode {

    /**
     * The Cosmos method signatures this rule applies to.
     */
    private List<String> methodSignatures;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the method signatures.
     *
     * @return the method signatures
     */
    public List<String> getMethodSignatures() {
        return methodSignatures;
    }

    /**
     * Sets the method signatures.
     *
     * @param methodSignatures the method signatures
     */
    public void setMethodSignatures(List<String> methodSignatures) {
        this.methodSignatures = methodSignatures;
    }
}
