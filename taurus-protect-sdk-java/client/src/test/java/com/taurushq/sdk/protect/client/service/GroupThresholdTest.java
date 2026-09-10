package com.taurushq.sdk.protect.client.service;

import com.google.gson.Gson;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.WhitelistUserSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.Security;
import java.security.spec.ECGenParameterSpec;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Base64;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;

/**
 * Per-group step-5 threshold: distinct signers, not signature entries.
 * <p>
 * The entries arrive in the server-supplied userSignatures blob, so counting them let a
 * duplicated entry from one group member satisfy an N-of-M group — promoting an
 * under-approved whitelist entry to approved.
 */
class GroupThresholdTest {

    private static final Gson GSON = new Gson();
    private static final String METADATA_HASH = "abc123";

    @BeforeAll
    static void setUpProvider() {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
    }

    private static KeyPair generateKeyPair() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        return generator.generateKeyPair();
    }

    /** A genuinely valid approval entry, as it arrives on the wire. */
    private static WhitelistSignature signedEntry(String userId, PrivateKey key) throws Exception {
        List<String> hashes = Collections.singletonList(METADATA_HASH);
        byte[] toSign = GSON.toJson(hashes).getBytes(StandardCharsets.UTF_8);

        WhitelistSignature sig = new WhitelistSignature();
        sig.getHashes().addAll(hashes);
        WhitelistUserSignature userSig = new WhitelistUserSignature();
        userSig.setUserId(userId);
        userSig.setSignature(
                Base64.getDecoder().decode(CryptoTPV1.calculateBase64Signature(key, toSign)));
        sig.setSignature(userSig);
        return sig;
    }

    private static DecodedRulesContainer container(List<String> userIds, List<PublicKey> keys) {
        DecodedRulesContainer container = new DecodedRulesContainer();
        List<RuleUser> users = new ArrayList<>();
        for (int i = 0; i < userIds.size(); i++) {
            RuleUser user = new RuleUser();
            user.setId(userIds.get(i));
            user.setPublicKey(keys.get(i));
            users.add(user);
        }
        container.setUsers(users);

        RuleGroup group = new RuleGroup();
        group.setId("approvers");
        group.setUserIds(new ArrayList<>(userIds));
        container.setGroups(Collections.singletonList(group));
        return container;
    }

    private static GroupThreshold threshold(int minSigs) {
        GroupThreshold gt = new GroupThreshold();
        gt.setGroupId("approvers");
        gt.setMinimumSignatures(minSigs);
        return gt;
    }

    /** The service needs no network for the walk; only a client to construct. */
    private static WhitelistedAddressService service(PublicKey superAdminKey) {
        return new WhitelistedAddressService(new ApiClient(), new ApiExceptionMapper(),
                Collections.singletonList(superAdminKey), 1);
    }

    @Test
    @DisplayName("one member's entry duplicated does not satisfy a 2-of-N group")
    void duplicateEntryFromOneMemberDoesNotMeetTwoOfN() throws Exception {
        KeyPair u1 = generateKeyPair();
        KeyPair u2 = generateKeyPair();
        DecodedRulesContainer rules = container(
                Arrays.asList("user1@bank.com", "user2@bank.com"),
                Arrays.asList(u1.getPublic(), u2.getPublic()));

        WhitelistSignature entry = signedEntry("user1@bank.com", u1.getPrivate());
        List<WhitelistSignature> signatures = Arrays.asList(entry, entry);

        WhitelistedAddressService svc = service(u1.getPublic());
        assertThrows(IntegrityException.class, () ->
                svc.verifyGroupThreshold(threshold(2), rules, signatures, METADATA_HASH));
    }

    @Test
    @DisplayName("two re-signed entries from one member do not satisfy a 2-of-N group")
    void reSignedEntryFromOneMemberDoesNotMeetTwoOfN() throws Exception {
        KeyPair u1 = generateKeyPair();
        KeyPair u2 = generateKeyPair();
        DecodedRulesContainer rules = container(
                Arrays.asList("user1@bank.com", "user2@bank.com"),
                Arrays.asList(u1.getPublic(), u2.getPublic()));

        WhitelistSignature first = signedEntry("user1@bank.com", u1.getPrivate());
        WhitelistSignature second = signedEntry("user1@bank.com", u1.getPrivate());
        // ECDSA is randomized, so the same key yields two different signatures.
        // Compared by content: assertNotEquals on arrays compares references.
        assertFalse(Arrays.equals(first.getSignature().getSignature(),
                second.getSignature().getSignature()));

        List<WhitelistSignature> signatures = Arrays.asList(first, second);
        WhitelistedAddressService svc = service(u1.getPublic());
        assertThrows(IntegrityException.class, () ->
                svc.verifyGroupThreshold(threshold(2), rules, signatures, METADATA_HASH));
    }

    @Test
    @DisplayName("two distinct members satisfy a 2-of-N group")
    void twoDistinctMembersMeetTwoOfN() throws Exception {
        KeyPair u1 = generateKeyPair();
        KeyPair u2 = generateKeyPair();
        DecodedRulesContainer rules = container(
                Arrays.asList("user1@bank.com", "user2@bank.com"),
                Arrays.asList(u1.getPublic(), u2.getPublic()));

        List<WhitelistSignature> signatures = Arrays.asList(
                signedEntry("user1@bank.com", u1.getPrivate()),
                signedEntry("user2@bank.com", u2.getPrivate()));

        WhitelistedAddressService svc = service(u1.getPublic());
        assertDoesNotThrow(() ->
                svc.verifyGroupThreshold(threshold(2), rules, signatures, METADATA_HASH));
    }

    @Test
    @DisplayName("one member satisfies a 1-of-N group")
    void oneMemberMeetsOneOfN() throws Exception {
        KeyPair u1 = generateKeyPair();
        KeyPair u2 = generateKeyPair();
        DecodedRulesContainer rules = container(
                Arrays.asList("user1@bank.com", "user2@bank.com"),
                Arrays.asList(u1.getPublic(), u2.getPublic()));

        List<WhitelistSignature> signatures =
                Collections.singletonList(signedEntry("user1@bank.com", u1.getPrivate()));

        WhitelistedAddressService svc = service(u1.getPublic());
        assertDoesNotThrow(() ->
                svc.verifyGroupThreshold(threshold(1), rules, signatures, METADATA_HASH));
    }

    @Test
    @DisplayName("two user IDs sharing one key count as one signer")
    void twoUserIdsSharingOneKeyCountOnce() throws Exception {
        KeyPair shared = generateKeyPair();
        DecodedRulesContainer rules = container(
                Arrays.asList("user1@bank.com", "user2@bank.com"),
                Arrays.asList(shared.getPublic(), shared.getPublic()));

        List<WhitelistSignature> signatures = Arrays.asList(
                signedEntry("user1@bank.com", shared.getPrivate()),
                signedEntry("user2@bank.com", shared.getPrivate()));

        WhitelistedAddressService svc = service(shared.getPublic());
        assertThrows(IntegrityException.class, () ->
                svc.verifyGroupThreshold(threshold(2), rules, signatures, METADATA_HASH));
    }

    @Test
    @DisplayName("a valid signer outside the group contributes nothing")
    void signerOutsideGroupContributesNothing() throws Exception {
        KeyPair u1 = generateKeyPair();
        KeyPair u2 = generateKeyPair();
        KeyPair outsider = generateKeyPair();

        DecodedRulesContainer rules = container(
                Arrays.asList("user1@bank.com", "user2@bank.com"),
                Arrays.asList(u1.getPublic(), u2.getPublic()));
        RuleUser extra = new RuleUser();
        extra.setId("outsider@bank.com");
        extra.setPublicKey(outsider.getPublic());
        rules.getUsers().add(extra);

        List<WhitelistSignature> signatures = Arrays.asList(
                signedEntry("user1@bank.com", u1.getPrivate()),
                signedEntry("outsider@bank.com", outsider.getPrivate()));

        WhitelistedAddressService svc = service(u1.getPublic());
        assertThrows(IntegrityException.class, () ->
                svc.verifyGroupThreshold(threshold(2), rules, signatures, METADATA_HASH));
    }
}
