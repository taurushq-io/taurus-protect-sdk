package com.taurushq.sdk.protect.client.model.rulescontainer;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertSame;

/**
 * Tests for {@link DecodedRulesContainer#findAddressWhitelistingRules(String, String)}.
 */
class DecodedRulesContainerTest {

    private DecodedRulesContainer container;

    @BeforeEach
    void setUp() {
        container = new DecodedRulesContainer();
    }

    // --- Helper methods ---

    private AddressWhitelistingRules createRule(String currency, String network) {
        AddressWhitelistingRules rule = new AddressWhitelistingRules();
        rule.setCurrency(currency);
        rule.setNetwork(network);
        return rule;
    }

    private void setRules(AddressWhitelistingRules... rules) {
        container.setAddressWhitelistingRules(Arrays.asList(rules));
    }

    // --- Null/empty rules list tests ---

    @Test
    void findAddressWhitelistingRules_nullRulesList_returnsNull() {
        container.setAddressWhitelistingRules(null);
        assertNull(container.findAddressWhitelistingRules("ETH", "mainnet"));
    }

    @Test
    void findAddressWhitelistingRules_emptyRulesList_returnsNull() {
        container.setAddressWhitelistingRules(new ArrayList<>());
        assertNull(container.findAddressWhitelistingRules("ETH", "mainnet"));
    }

    // --- Exact match tests ---

    @Test
    void findAddressWhitelistingRules_exactMatch() {
        AddressWhitelistingRules ethMainnet = createRule("ETH", "mainnet");
        setRules(ethMainnet);

        assertSame(ethMainnet, container.findAddressWhitelistingRules("ETH", "mainnet"));
    }

    @Test
    void findAddressWhitelistingRules_exactMatch_casePreserved() {
        AddressWhitelistingRules ethMainnet = createRule("ETH", "mainnet");
        setRules(ethMainnet);

        // Different case should NOT match (case-sensitive)
        assertNull(container.findAddressWhitelistingRules("eth", "mainnet"));
        assertNull(container.findAddressWhitelistingRules("ETH", "MAINNET"));
    }

    @Test
    void findAddressWhitelistingRules_noMatchDifferentBlockchain() {
        AddressWhitelistingRules ethMainnet = createRule("ETH", "mainnet");
        setRules(ethMainnet);

        assertNull(container.findAddressWhitelistingRules("BTC", "mainnet"));
    }

    @Test
    void findAddressWhitelistingRules_noMatchDifferentNetwork() {
        AddressWhitelistingRules ethMainnet = createRule("ETH", "mainnet");
        setRules(ethMainnet);

        // No fallback rule, so null
        assertNull(container.findAddressWhitelistingRules("ETH", "testnet"));
    }

    // --- Blockchain-only match (wildcard network) tests ---

    @Test
    void findAddressWhitelistingRules_blockchainOnlyMatch_nullNetwork() {
        AddressWhitelistingRules ethAnyNetwork = createRule("ETH", null);
        setRules(ethAnyNetwork);

        assertSame(ethAnyNetwork, container.findAddressWhitelistingRules("ETH", "mainnet"));
        assertSame(ethAnyNetwork, container.findAddressWhitelistingRules("ETH", "testnet"));
        assertSame(ethAnyNetwork, container.findAddressWhitelistingRules("ETH", "goerli"));
    }

    @Test
    void findAddressWhitelistingRules_blockchainOnlyMatch_emptyNetwork() {
        AddressWhitelistingRules ethAnyNetwork = createRule("ETH", "");
        setRules(ethAnyNetwork);

        assertSame(ethAnyNetwork, container.findAddressWhitelistingRules("ETH", "mainnet"));
        assertSame(ethAnyNetwork, container.findAddressWhitelistingRules("ETH", "testnet"));
    }

    @Test
    void findAddressWhitelistingRules_blockchainOnlyMatch_doesNotMatchOtherBlockchain() {
        AddressWhitelistingRules ethAnyNetwork = createRule("ETH", null);
        setRules(ethAnyNetwork);

        assertNull(container.findAddressWhitelistingRules("BTC", "mainnet"));
    }

    // --- Global default tests ---

    @Test
    void findAddressWhitelistingRules_globalDefault_nullCurrency() {
        AddressWhitelistingRules globalDefault = createRule(null, null);
        setRules(globalDefault);

        assertSame(globalDefault, container.findAddressWhitelistingRules("ETH", "mainnet"));
        assertSame(globalDefault, container.findAddressWhitelistingRules("BTC", "testnet"));
        assertSame(globalDefault, container.findAddressWhitelistingRules("SOL", "devnet"));
    }

    @Test
    void findAddressWhitelistingRules_globalDefault_emptyCurrency() {
        AddressWhitelistingRules globalDefault = createRule("", null);
        setRules(globalDefault);

        assertSame(globalDefault, container.findAddressWhitelistingRules("ETH", "mainnet"));
        assertSame(globalDefault, container.findAddressWhitelistingRules("BTC", "testnet"));
    }

    @Test
    void findAddressWhitelistingRules_globalDefault_anyCurrency() {
        AddressWhitelistingRules globalDefault = createRule("Any", null);
        setRules(globalDefault);

        assertSame(globalDefault, container.findAddressWhitelistingRules("ETH", "mainnet"));
        assertSame(globalDefault, container.findAddressWhitelistingRules("BTC", "testnet"));
        assertSame(globalDefault, container.findAddressWhitelistingRules("SOL", "mainnet"));
    }

    @Test
    void findAddressWhitelistingRules_globalDefault_anyCurrencyCaseInsensitive() {
        AddressWhitelistingRules globalDefaultLower = createRule("any", null);
        AddressWhitelistingRules globalDefaultUpper = createRule("ANY", null);

        container.setAddressWhitelistingRules(Arrays.asList(globalDefaultLower));
        assertSame(globalDefaultLower, container.findAddressWhitelistingRules("ETH", "mainnet"));

        container.setAddressWhitelistingRules(Arrays.asList(globalDefaultUpper));
        assertSame(globalDefaultUpper, container.findAddressWhitelistingRules("ETH", "mainnet"));
    }

    // --- Priority tests ---

    @Test
    void findAddressWhitelistingRules_priority_exactOverBlockchainOnly() {
        AddressWhitelistingRules ethMainnet = createRule("ETH", "mainnet");
        AddressWhitelistingRules ethAnyNetwork = createRule("ETH", null);

        // Order shouldn't matter - exact match should always win
        setRules(ethAnyNetwork, ethMainnet);
        assertSame(ethMainnet, container.findAddressWhitelistingRules("ETH", "mainnet"));

        setRules(ethMainnet, ethAnyNetwork);
        assertSame(ethMainnet, container.findAddressWhitelistingRules("ETH", "mainnet"));
    }

    @Test
    void findAddressWhitelistingRules_priority_blockchainOnlyOverGlobal() {
        AddressWhitelistingRules ethAnyNetwork = createRule("ETH", null);
        AddressWhitelistingRules globalDefault = createRule(null, null);

        // Order shouldn't matter - blockchain-only should win over global
        setRules(globalDefault, ethAnyNetwork);
        assertSame(ethAnyNetwork, container.findAddressWhitelistingRules("ETH", "mainnet"));

        setRules(ethAnyNetwork, globalDefault);
        assertSame(ethAnyNetwork, container.findAddressWhitelistingRules("ETH", "mainnet"));
    }

    @Test
    void findAddressWhitelistingRules_priority_exactOverGlobal() {
        AddressWhitelistingRules ethMainnet = createRule("ETH", "mainnet");
        AddressWhitelistingRules globalDefault = createRule("Any", null);

        // Order shouldn't matter - exact match should win over global
        setRules(globalDefault, ethMainnet);
        assertSame(ethMainnet, container.findAddressWhitelistingRules("ETH", "mainnet"));

        setRules(ethMainnet, globalDefault);
        assertSame(ethMainnet, container.findAddressWhitelistingRules("ETH", "mainnet"));
    }

    @Test
    void findAddressWhitelistingRules_priority_allThreeTiers() {
        AddressWhitelistingRules ethMainnet = createRule("ETH", "mainnet");
        AddressWhitelistingRules ethAnyNetwork = createRule("ETH", null);
        AddressWhitelistingRules globalDefault = createRule(null, null);

        setRules(globalDefault, ethAnyNetwork, ethMainnet);

        // ETH/mainnet -> exact match
        assertSame(ethMainnet, container.findAddressWhitelistingRules("ETH", "mainnet"));

        // ETH/testnet -> blockchain-only match
        assertSame(ethAnyNetwork, container.findAddressWhitelistingRules("ETH", "testnet"));

        // BTC/mainnet -> global default
        assertSame(globalDefault, container.findAddressWhitelistingRules("BTC", "mainnet"));
    }

    // --- First match within tier tests ---

    @Test
    void findAddressWhitelistingRules_firstBlockchainOnlyMatchWins() {
        AddressWhitelistingRules ethNull = createRule("ETH", null);
        AddressWhitelistingRules ethEmpty = createRule("ETH", "");

        setRules(ethNull, ethEmpty);
        assertSame(ethNull, container.findAddressWhitelistingRules("ETH", "testnet"));

        setRules(ethEmpty, ethNull);
        assertSame(ethEmpty, container.findAddressWhitelistingRules("ETH", "testnet"));
    }

    @Test
    void findAddressWhitelistingRules_firstGlobalDefaultWins() {
        AddressWhitelistingRules globalNull = createRule(null, null);
        AddressWhitelistingRules globalAny = createRule("Any", null);

        setRules(globalNull, globalAny);
        assertSame(globalNull, container.findAddressWhitelistingRules("BTC", "mainnet"));

        setRules(globalAny, globalNull);
        assertSame(globalAny, container.findAddressWhitelistingRules("BTC", "mainnet"));
    }

    // --- Edge cases ---

    @Test
    void findAddressWhitelistingRules_nullSearchBlockchain() {
        AddressWhitelistingRules ethMainnet = createRule("ETH", "mainnet");
        AddressWhitelistingRules globalDefault = createRule(null, null);

        setRules(ethMainnet, globalDefault);

        // Null blockchain should match global default
        assertSame(globalDefault, container.findAddressWhitelistingRules(null, "mainnet"));
    }

    @Test
    void findAddressWhitelistingRules_nullSearchNetwork() {
        AddressWhitelistingRules ethNull = createRule("ETH", null);
        AddressWhitelistingRules ethMainnet = createRule("ETH", "mainnet");

        setRules(ethMainnet, ethNull);

        // Null network should not match "mainnet" but should match null
        assertSame(ethNull, container.findAddressWhitelistingRules("ETH", null));
    }

    @Test
    void findAddressWhitelistingRules_bothSearchValuesNull() {
        AddressWhitelistingRules globalDefault = createRule(null, null);
        setRules(globalDefault);

        assertSame(globalDefault, container.findAddressWhitelistingRules(null, null));
    }

    @Test
    void findAddressWhitelistingRules_multipleBlockchains() {
        AddressWhitelistingRules ethMainnet = createRule("ETH", "mainnet");
        AddressWhitelistingRules btcMainnet = createRule("BTC", "mainnet");
        AddressWhitelistingRules btcAny = createRule("BTC", null);
        AddressWhitelistingRules globalDefault = createRule("Any", null);

        setRules(ethMainnet, btcMainnet, btcAny, globalDefault);

        assertSame(ethMainnet, container.findAddressWhitelistingRules("ETH", "mainnet"));
        assertSame(btcMainnet, container.findAddressWhitelistingRules("BTC", "mainnet"));
        assertSame(btcAny, container.findAddressWhitelistingRules("BTC", "testnet"));
        assertSame(globalDefault, container.findAddressWhitelistingRules("SOL", "mainnet"));
        assertSame(globalDefault, container.findAddressWhitelistingRules("ETH", "testnet"));
    }
}
