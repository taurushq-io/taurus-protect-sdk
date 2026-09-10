package com.taurushq.sdk.protect.client.service;

import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.WhitelistSignature;
import com.taurushq.sdk.protect.client.model.WhitelistUserSignature;
import com.taurushq.sdk.protect.client.model.rulescontainer.DecodedRulesContainer;
import com.taurushq.sdk.protect.client.model.rulescontainer.GroupThreshold;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleGroup;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleUser;
import com.taurushq.sdk.protect.client.testutil.SignedFixtures;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import org.bouncycastle.jce.provider.BouncyCastleProvider;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.security.KeyPairGenerator;
import java.security.PublicKey;
import java.security.Security;
import java.security.spec.ECGenParameterSpec;
import java.util.ArrayList;
import java.util.Base64;
import java.util.Collections;
import java.util.HashSet;
import java.util.List;
import java.util.Set;

import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.fail;

/**
 * Cross-SDK signed fixtures for verification step 5 — the per-group threshold.
 * <p>
 * A different threshold from the SuperAdmin one in
 * {@code helper.SignedFixturesTest} (per group, not tenant-wide) but the same counting
 * rule, so both are gated from one file: the four SDKs cannot drift on one without
 * drifting on the other.
 * <p>
 * It lives in this package because the walk is a package-private service method — Java
 * keeps whitelist verification in the services rather than {@code helper/}.
 */
class SignedFixturesGroupThresholdTest {

    private static JsonObject fixtures;

    @BeforeAll
    static void loadFixtures() throws IOException {
        if (Security.getProvider(BouncyCastleProvider.PROVIDER_NAME) == null) {
            Security.addProvider(new BouncyCastleProvider());
        }
        fixtures = SignedFixtures.load();
    }

    /** The walk needs no network; only a client to construct the service. */
    private static WhitelistedAddressService service() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        PublicKey unrelated = generator.generateKeyPair().getPublic();
        return new WhitelistedAddressService(new ApiClient(), new ApiExceptionMapper(),
                Collections.singletonList(unrelated), 1);
    }

    private static DecodedRulesContainer rulesFor(final JsonObject c) throws Exception {
        DecodedRulesContainer container = new DecodedRulesContainer();

        List<RuleUser> users = new ArrayList<>();
        for (JsonElement element : c.getAsJsonArray("users")) {
            JsonObject u = element.getAsJsonObject();
            RuleUser user = new RuleUser();
            user.setId(u.get("user_id").getAsString());
            user.setPublicKey(CryptoTPV1.decodePublicKey(u.get("public_key_pem").getAsString()));
            users.add(user);
        }
        container.setUsers(users);

        List<String> memberIds = new ArrayList<>();
        for (JsonElement id : c.getAsJsonArray("group_user_ids")) {
            memberIds.add(id.getAsString());
        }
        RuleGroup group = new RuleGroup();
        group.setId(c.get("group_id").getAsString());
        group.setUserIds(memberIds);
        container.setGroups(Collections.singletonList(group));

        return container;
    }

    private static List<WhitelistSignature> signaturesFor(final JsonObject c) {
        List<WhitelistSignature> signatures = new ArrayList<>();
        for (JsonElement element : c.getAsJsonArray("signatures")) {
            JsonObject s = element.getAsJsonObject();
            WhitelistSignature sig = new WhitelistSignature();
            for (JsonElement hash : s.getAsJsonArray("hashes")) {
                sig.getHashes().add(hash.getAsString());
            }
            WhitelistUserSignature userSig = new WhitelistUserSignature();
            userSig.setUserId(s.get("user_id").getAsString());
            userSig.setSignature(Base64.getDecoder().decode(s.get("signature").getAsString()));
            sig.setSignature(userSig);
            signatures.add(sig);
        }
        return signatures;
    }

    @Test
    void agreesWithEveryRecordedOutcome() throws Exception {
        WhitelistedAddressService svc = service();

        for (JsonElement element : fixtures.getAsJsonArray("group_threshold")) {
            JsonObject c = element.getAsJsonObject();
            String description = c.get("description").getAsString();

            GroupThreshold gt = new GroupThreshold();
            gt.setGroupId(c.get("group_id").getAsString());
            gt.setMinimumSignatures(c.get("minimum_signatures").getAsInt());

            boolean expectError = "error".equals(c.get("expect").getAsString());
            try {
                svc.verifyGroupThreshold(gt, rulesFor(c), signaturesFor(c),
                        c.get("metadata_hash").getAsString());
                if (expectError) {
                    fail(description + ": expected the threshold to fail, it passed");
                }
            } catch (IntegrityException | IllegalArgumentException e) {
                if (!expectError) {
                    throw new AssertionError(
                            description + ": expected the threshold to pass", e);
                }
            }
        }
    }

    /**
     * Without a case where the in-group entry COUNT meets the threshold but the
     * distinct-signer count does not, the group section would pass against an
     * implementation that counts entries.
     */
    @Test
    void containsACaseThatDistinguishesDistinctSignersFromEntries() {
        boolean discriminating = false;
        for (JsonElement element : fixtures.getAsJsonArray("group_threshold")) {
            JsonObject c = element.getAsJsonObject();
            int minSigs = c.get("minimum_signatures").getAsInt();
            if (!"error".equals(c.get("expect").getAsString()) || minSigs <= 0) {
                continue;
            }

            Set<String> members = new HashSet<>();
            for (JsonElement id : c.getAsJsonArray("group_user_ids")) {
                members.add(id.getAsString());
            }
            int inGroup = 0;
            for (JsonElement sigElement : c.getAsJsonArray("signatures")) {
                if (members.contains(sigElement.getAsJsonObject().get("user_id").getAsString())) {
                    inGroup++;
                }
            }
            if (inGroup >= minSigs) {
                discriminating = true;
                break;
            }
        }
        assertTrue(discriminating,
                "no group vector distinguishes distinct-signer counting from entry counting");
    }
}
