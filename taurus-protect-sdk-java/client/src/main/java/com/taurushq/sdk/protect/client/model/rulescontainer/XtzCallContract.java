package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents Tezos (XTZ) smart contract call configuration.
 * <p>
 * This class defines which contract types and methods are allowed for a transaction
 * rule on the Tezos blockchain.
 *
 * @see TransactionRuleDetails
 */
public class XtzCallContract {

    /**
     * The contract type (e.g., "FA12", "FA20" for Tezos token standards).
     */
    private String contractType;

    /**
     * The entrypoint/method signature or null for any entrypoint.
     */
    private String methodSignature;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the contract type.
     *
     * @return the contract type
     */
    public String getContractType() {
        return contractType;
    }

    /**
     * Sets the contract type.
     *
     * @param contractType the contract type
     */
    public void setContractType(String contractType) {
        this.contractType = contractType;
    }

    /**
     * Gets the method signature.
     *
     * @return the method signature
     */
    public String getMethodSignature() {
        return methodSignature;
    }

    /**
     * Sets the method signature.
     *
     * @param methodSignature the method signature
     */
    public void setMethodSignature(String methodSignature) {
        this.methodSignature = methodSignature;
    }
}
