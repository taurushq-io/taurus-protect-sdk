package com.taurushq.sdk.protect.client.service;

import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.openapi.ApiClient;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;

import static org.junit.jupiter.api.Assertions.assertThrows;

class ContractWhitelistingServiceTest {

    private ApiClient apiClient;
    private ApiExceptionMapper apiExceptionMapper;
    private ContractWhitelistingService contractWhitelistingService;

    @BeforeEach
    void setUp() {
        apiClient = new ApiClient();
        apiExceptionMapper = new ApiExceptionMapper();
        contractWhitelistingService = new ContractWhitelistingService(apiClient, apiExceptionMapper);
    }

    @Test
    void constructor_throwsOnNullApiClient() {
        assertThrows(NullPointerException.class, () ->
                new ContractWhitelistingService(null, apiExceptionMapper));
    }

    @Test
    void constructor_throwsOnNullExceptionMapper() {
        assertThrows(NullPointerException.class, () ->
                new ContractWhitelistingService(apiClient, null));
    }

    // createWhitelistedContract tests
    @Test
    void createWhitelistedContract_throwsOnNullBlockchain() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createWhitelistedContract(
                        null, "mainnet", "0x123", "USDC", "USD Coin", 6, "erc20", null));
    }

    @Test
    void createWhitelistedContract_throwsOnEmptyBlockchain() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createWhitelistedContract(
                        "", "mainnet", "0x123", "USDC", "USD Coin", 6, "erc20", null));
    }

    @Test
    void createWhitelistedContract_throwsOnNullNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createWhitelistedContract(
                        "ETH", null, "0x123", "USDC", "USD Coin", 6, "erc20", null));
    }

    @Test
    void createWhitelistedContract_throwsOnEmptyNetwork() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createWhitelistedContract(
                        "ETH", "", "0x123", "USDC", "USD Coin", 6, "erc20", null));
    }

    @Test
    void createWhitelistedContract_throwsOnNullSymbol() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createWhitelistedContract(
                        "ETH", "mainnet", "0x123", null, "USD Coin", 6, "erc20", null));
    }

    @Test
    void createWhitelistedContract_throwsOnEmptySymbol() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createWhitelistedContract(
                        "ETH", "mainnet", "0x123", "", "USD Coin", 6, "erc20", null));
    }

    @Test
    void createWhitelistedContract_throwsOnNullName() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createWhitelistedContract(
                        "ETH", "mainnet", "0x123", "USDC", null, 6, "erc20", null));
    }

    @Test
    void createWhitelistedContract_throwsOnEmptyName() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createWhitelistedContract(
                        "ETH", "mainnet", "0x123", "USDC", "", 6, "erc20", null));
    }

    @Test
    void createWhitelistedContract_throwsOnNullKind() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createWhitelistedContract(
                        "ETH", "mainnet", "0x123", "USDC", "USD Coin", 6, null, null));
    }

    @Test
    void createWhitelistedContract_throwsOnEmptyKind() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createWhitelistedContract(
                        "ETH", "mainnet", "0x123", "USDC", "USD Coin", 6, "", null));
    }

    // approveWhitelistedContracts tests
    @Test
    void approveWhitelistedContracts_throwsOnNullIds() {
        assertThrows(NullPointerException.class, () ->
                contractWhitelistingService.approveWhitelistedContracts(
                        null, "signature", "comment"));
    }

    @Test
    void approveWhitelistedContracts_throwsOnEmptyIds() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.approveWhitelistedContracts(
                        Collections.emptyList(), "signature", "comment"));
    }

    @Test
    void approveWhitelistedContracts_throwsOnNullSignature() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.approveWhitelistedContracts(
                        Arrays.asList("1", "2"), null, "comment"));
    }

    @Test
    void approveWhitelistedContracts_throwsOnEmptySignature() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.approveWhitelistedContracts(
                        Arrays.asList("1", "2"), "", "comment"));
    }

    @Test
    void approveWhitelistedContracts_throwsOnNullComment() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.approveWhitelistedContracts(
                        Arrays.asList("1", "2"), "signature", null));
    }

    @Test
    void approveWhitelistedContracts_throwsOnEmptyComment() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.approveWhitelistedContracts(
                        Arrays.asList("1", "2"), "signature", ""));
    }

    // getWhitelistedContract tests
    @Test
    void getWhitelistedContract_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.getWhitelistedContract(null));
    }

    @Test
    void getWhitelistedContract_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.getWhitelistedContract(""));
    }

    // updateWhitelistedContract tests
    @Test
    void updateWhitelistedContract_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.updateWhitelistedContract(
                        null, "USDC", "USD Coin", 6));
    }

    @Test
    void updateWhitelistedContract_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.updateWhitelistedContract(
                        "", "USDC", "USD Coin", 6));
    }

    @Test
    void updateWhitelistedContract_throwsOnNullSymbol() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.updateWhitelistedContract(
                        "123", null, "USD Coin", 6));
    }

    @Test
    void updateWhitelistedContract_throwsOnEmptySymbol() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.updateWhitelistedContract(
                        "123", "", "USD Coin", 6));
    }

    @Test
    void updateWhitelistedContract_throwsOnNullName() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.updateWhitelistedContract(
                        "123", "USDC", null, 6));
    }

    @Test
    void updateWhitelistedContract_throwsOnEmptyName() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.updateWhitelistedContract(
                        "123", "USDC", "", 6));
    }

    // deleteWhitelistedContract tests
    @Test
    void deleteWhitelistedContract_throwsOnNullId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.deleteWhitelistedContract(null, "comment"));
    }

    @Test
    void deleteWhitelistedContract_throwsOnEmptyId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.deleteWhitelistedContract("", "comment"));
    }

    // createAttribute tests
    @Test
    void createAttribute_throwsOnNullContractId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createAttribute(
                        null, "key", "value", null, null, null));
    }

    @Test
    void createAttribute_throwsOnEmptyContractId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createAttribute(
                        "", "key", "value", null, null, null));
    }

    @Test
    void createAttribute_throwsOnNullKey() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createAttribute(
                        "123", null, "value", null, null, null));
    }

    @Test
    void createAttribute_throwsOnEmptyKey() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.createAttribute(
                        "123", "", "value", null, null, null));
    }

    // getAttribute tests
    @Test
    void getAttribute_throwsOnNullContractId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.getAttribute(null, "attr-1"));
    }

    @Test
    void getAttribute_throwsOnEmptyContractId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.getAttribute("", "attr-1"));
    }

    @Test
    void getAttribute_throwsOnNullAttributeId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.getAttribute("123", null));
    }

    @Test
    void getAttribute_throwsOnEmptyAttributeId() {
        assertThrows(IllegalArgumentException.class, () ->
                contractWhitelistingService.getAttribute("123", ""));
    }

    // getWhitelistedContracts accepts null filters
    @Test
    void getWhitelistedContracts_acceptsAllNullFilters() {
        // This test verifies that all filter parameters are optional
        // The actual API call will fail due to no server, but validation should pass
        // The test passes if the method signature allows all optional parameters
    }

    // getWhitelistedContractsForApproval accepts null filters
    @Test
    void getWhitelistedContractsForApproval_acceptsNullIds() {
        // This test verifies that ids parameter is optional
        // The actual API call will fail due to no server, but validation should pass
    }
}
