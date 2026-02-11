package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Represents EVM (Ethereum Virtual Machine) smart contract call configuration.
 * <p>
 * This class defines which contract types and methods are allowed for a transaction
 * rule on EVM-compatible chains (Ethereum, Polygon, BSC, etc.).
 *
 * @see TransactionRuleDetails
 */
public class EvmCallContract {

    /**
     * The contract type (e.g., "GENERIC", "ERC20", "CMTAT", "CMTA20").
     */
    private String contractType;

    /**
     * The method signature (e.g., "transfer(address,uint256)") or null for any method.
     */
    private String methodSignature;

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /**
     * Gets the contract type.
     *
     * @return the contract type (e.g., GENERIC, ERC20, CMTAT, CMTA20)
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
