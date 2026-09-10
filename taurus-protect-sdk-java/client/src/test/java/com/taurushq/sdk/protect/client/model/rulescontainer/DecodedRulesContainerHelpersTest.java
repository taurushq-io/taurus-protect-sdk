package com.taurushq.sdk.protect.client.model.rulescontainer;

import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertSame;

import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.spec.ECGenParameterSpec;
import java.util.Arrays;
import java.util.Collections;

import org.junit.jupiter.api.Test;

/**
 * Coverage for the {@link DecodedRulesContainer} lookup helpers beyond
 * findAddressWhitelistingRules (covered by {@link DecodedRulesContainerTest}):
 * the contract-address three-tier lookup, user/group id lookups, and the cached
 * HSM public key resolution.
 */
class DecodedRulesContainerHelpersTest {

    private static PublicKey p256Key() throws Exception {
        KeyPairGenerator kpg = KeyPairGenerator.getInstance("EC");
        kpg.initialize(new ECGenParameterSpec("secp256r1"));
        return kpg.generateKeyPair().getPublic();
    }

    private static ContractAddressWhitelistingRules contractRule(String blockchain, String network) {
        ContractAddressWhitelistingRules rule = new ContractAddressWhitelistingRules();
        rule.setBlockchain(blockchain);
        rule.setNetwork(network);
        return rule;
    }

    private static RuleUser user(String id, PublicKey key, String... roles) {
        RuleUser u = new RuleUser();
        u.setId(id);
        u.setRoles(Arrays.asList(roles));
        if (key != null) {
            u.setPublicKey(key);
        }
        return u;
    }

    // --- findContractAddressWhitelistingRules: exact > blockchain-only > global ---

    @Test
    void findContractAddressWhitelistingRules_tiers() {
        ContractAddressWhitelistingRules exact = contractRule("ETH", "mainnet");
        ContractAddressWhitelistingRules blockchainOnly = contractRule("ETH", "Any");
        ContractAddressWhitelistingRules global = contractRule("", "mainnet");
        DecodedRulesContainer c = new DecodedRulesContainer();
        c.setContractAddressWhitelistingRules(Arrays.asList(exact, blockchainOnly, global));

        assertSame(exact, c.findContractAddressWhitelistingRules("ETH", "mainnet"));
        assertSame(blockchainOnly, c.findContractAddressWhitelistingRules("ETH", "testnet"));
        assertSame(global, c.findContractAddressWhitelistingRules("BTC", "mainnet"));
    }

    @Test
    void findContractAddressWhitelistingRules_nullOrNoMatchReturnsNull() {
        assertNull(new DecodedRulesContainer().findContractAddressWhitelistingRules("ETH", "mainnet"));

        DecodedRulesContainer noMatch = new DecodedRulesContainer();
        noMatch.setContractAddressWhitelistingRules(Collections.singletonList(contractRule("ETH", "mainnet")));
        assertNull(noMatch.findContractAddressWhitelistingRules("SOL", "devnet"));
    }

    // --- findUserById / findGroupById ---

    @Test
    void findUserById() {
        DecodedRulesContainer c = new DecodedRulesContainer();
        RuleUser u1 = user("u1", null);
        RuleUser u2 = user("u2", null);
        c.setUsers(Arrays.asList(u1, u2));

        assertSame(u2, c.findUserById("u2"));
        assertNull(c.findUserById("nope"));
        assertNull(c.findUserById(null));
        assertNull(new DecodedRulesContainer().findUserById("u1"));
    }

    @Test
    void findGroupById() {
        DecodedRulesContainer c = new DecodedRulesContainer();
        RuleGroup g = new RuleGroup();
        g.setId("g1");
        c.setGroups(Collections.singletonList(g));

        assertSame(g, c.findGroupById("g1"));
        assertNull(c.findGroupById("nope"));
        assertNull(c.findGroupById(null));
        assertNull(new DecodedRulesContainer().findGroupById("g1"));
    }

    // --- getHsmPublicKey (cached; first HSMSLOT user with a parsed key) ---

    @Test
    void getHsmPublicKey_returnsAndCachesHsmSlotUsersKey() throws Exception {
        PublicKey key = p256Key();
        DecodedRulesContainer c = new DecodedRulesContainer();
        c.setUsers(Arrays.asList(user("auth", null, "SUPERADMIN"), user("hsm", key, "HSMSLOT")));

        assertSame(key, c.getHsmPublicKey());

        // Cached: mutating users afterwards does not change the already-resolved key.
        c.setUsers(Collections.singletonList(user("hsm2", p256Key(), "HSMSLOT")));
        assertSame(key, c.getHsmPublicKey());
    }

    @Test
    void getHsmPublicKey_nullWhenNoHsmSlotOrNoKey() throws Exception {
        DecodedRulesContainer noHsm = new DecodedRulesContainer();
        noHsm.setUsers(Collections.singletonList(user("u", p256Key(), "SUPERADMIN")));
        assertNull(noHsm.getHsmPublicKey());

        DecodedRulesContainer nullKey = new DecodedRulesContainer();
        nullKey.setUsers(Collections.singletonList(user("hsm", null, "HSMSLOT")));
        assertNull(nullKey.getHsmPublicKey());
    }
}
