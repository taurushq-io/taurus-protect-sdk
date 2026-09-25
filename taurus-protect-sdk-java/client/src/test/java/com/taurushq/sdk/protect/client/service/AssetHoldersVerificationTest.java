package com.taurushq.sdk.protect.client.service;

import com.google.gson.JsonArray;
import com.google.gson.JsonObject;
import com.taurushq.sdk.protect.client.cache.RulesContainerCache;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ApiException;
import com.taurushq.sdk.protect.client.model.AssetAddressV2;
import com.taurushq.sdk.protect.client.model.AssetAddressV2Result;
import com.taurushq.sdk.protect.client.model.ContainerIntegrityException;
import com.taurushq.sdk.protect.client.model.ExcludedWhitelistedAddress;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.testutil.SignedWhitelist;
import com.taurushq.sdk.protect.client.testutil.StubTransport;
import com.taurushq.sdk.protect.openapi.ApiClient;
import com.taurushq.sdk.protect.openapi.auth.CryptoTPV1;
import com.taurushq.sdk.protect.proto.v1.RequestReply;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.nio.charset.StandardCharsets;
import java.security.KeyPair;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Base64;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The v2 asset-holders rows carry no signature, so {@code queryAssetAddresses} returns an
 * INTERNAL or WHITELISTED address only after its verified counterpart confirms it.
 * <p>
 * Everything runs through {@link StubTransport}: the holders page, the rules container the
 * managed-address reader verifies against, and the two verified readers, each answered by
 * the ids the request asks for. The keys are fresh per test, and every fixture row is
 * genuinely signed, so a kept row passed the real HSM or six-step verification.
 */
class AssetHoldersVerificationTest {

    private static final String HOLDERS = "/api/rest/v2/assets/a1/addresses/query";
    private static final String MANAGED = "/api/rest/v1/addresses";
    private static final String WHITELIST = "/api/rest/v1/whitelists/addresses";
    private static final String RULES = "/api/rest/v1/rules";
    private static final String NEXT = "bmV4dA==";

    private static final String INTERNAL = "ADDRESS_TYPE_V2_INTERNAL";
    private static final String WHITELISTED = "ADDRESS_TYPE_V2_WHITELISTED";
    private static final String EXTERNAL = "ADDRESS_TYPE_V2_EXTERNAL";

    private SignedWhitelist signedWhitelist;
    private KeyPair superAdmin;
    private KeyPair hsm;

    /** The holders page; it always claims a next page. */
    private final List<JsonObject> holders = new ArrayList<>();
    /** What the managed-address reader returns, by id. */
    private final Map<String, JsonObject> managed = new HashMap<>();
    /** What the whitelisted-address reader returns, by id. */
    private final Map<String, JsonObject> whitelisted = new HashMap<>();
    private boolean noHsmKey;
    private int managedStatus = 200;
    private StubTransport stub;

    @BeforeEach
    void setUp() throws Exception {
        signedWhitelist = new SignedWhitelist();
        superAdmin = signedWhitelist.superAdmin();
        hsm = SignedWhitelist.newKeyPair();
    }

    private static String pem(final PublicKey key) throws Exception {
        return CryptoTPV1.encodePublicKey(key);
    }

    private static String b64(final byte[] bytes) {
        return Base64.getEncoder().encodeToString(bytes);
    }

    /** The governance rules the managed-address reader verifies against: the HSMSLOT key. */
    private String rulesReply() {
        try {
            RequestReply.RulesContainer.Builder container = RequestReply.RulesContainer.newBuilder()
                    .addUsers(RequestReply.User.newBuilder().setId("superadmin@bank.com")
                            .setPublicKey(pem(superAdmin.getPublic())).addRoles(RequestReply.Role.SUPERADMIN));
            if (!noHsmKey) {
                container.addUsers(RequestReply.User.newBuilder().setId("hsm")
                        .setPublicKey(pem(hsm.getPublic())).addRoles(RequestReply.Role.HSMSLOT));
            }
            byte[] raw = container.build().toByteArray();
            JsonObject signature = new JsonObject();
            signature.addProperty("userId", "superadmin@bank.com");
            signature.addProperty("signature", CryptoTPV1.calculateBase64Signature(superAdmin.getPrivate(), raw));
            JsonArray signatures = new JsonArray();
            signatures.add(signature);
            JsonObject rules = new JsonObject();
            rules.addProperty("rulesContainer", b64(raw));
            rules.add("rulesSignatures", signatures);
            JsonObject reply = new JsonObject();
            reply.add("result", rules);
            return reply.toString();
        } catch (Exception e) {
            throw new IllegalStateException(e);
        }
    }

    private static JsonObject holder(final String address, final String type, final String addressId,
                                     final String whitelistedId) {
        JsonObject row = new JsonObject();
        row.addProperty("address", address);
        row.addProperty("balance", "5");
        row.addProperty("kycStatus", "KYC_STATUS_V2_APPROVED");
        if (type != null) {
            row.addProperty("addressType", type);
        }
        if (addressId != null) {
            row.addProperty("addressID", addressId);
        }
        if (whitelistedId != null) {
            row.addProperty("whitelistedAddressID", whitelistedId);
        }
        return row;
    }

    /** Makes the managed-address reader return id with an HSM signature over signedAddress. */
    private void manage(final String id, final String address, final String signedAddress) throws Exception {
        JsonObject row = new JsonObject();
        row.addProperty("id", id);
        row.addProperty("address", address);
        row.addProperty("signature", CryptoTPV1.calculateBase64Signature(hsm.getPrivate(),
                signedAddress.getBytes(StandardCharsets.UTF_8)));
        row.addProperty("status", "confirmed");
        managed.put(id, row);
    }

    /**
     * Makes the whitelisted-address reader return a signed envelope for id, whose address is
     * ADDR{id}; a tampered one carries a payload its hash does not cover.
     */
    private void whitelist(final String id, final boolean tampered) throws Exception {
        whitelisted.put(id, signedWhitelist.row(id, tampered));
    }

    private String holdersReply() {
        JsonArray result = new JsonArray();
        holders.forEach(result::add);
        JsonObject cursor = new JsonObject();
        cursor.addProperty("currentPage", NEXT);
        cursor.addProperty("hasNext", true);
        JsonObject reply = new JsonObject();
        reply.add("result", result);
        reply.add("cursor", cursor);
        return reply.toString();
    }

    private AssetAddressV2Result query() throws Exception {
        stub = StubTransport.replying()
                .route(HOLDERS, r -> holdersReply())
                .route(RULES, r -> rulesReply())
                .route(MANAGED, managedStatus, r -> managedStatus == 200
                        ? SignedWhitelist.page(managed, r.values("addressIds"))
                        : "{\"code\":13,\"message\":\"unavailable\"}")
                .route(WHITELIST, r -> SignedWhitelist.page(whitelisted, r.values("ids")));
        ApiClient client = stub.client();
        ApiExceptionMapper mapper = new ApiExceptionMapper();
        List<PublicKey> superAdminKeys = signedWhitelist.superAdminKeys();
        RulesContainerCache cache = new RulesContainerCache(
                new GovernanceRuleService(client, mapper, superAdminKeys, 1));
        AssetService assets = new AssetService(client, mapper, cache, new AddressService(client, mapper, cache),
                new WhitelistedAddressService(client, mapper, superAdminKeys, 1));
        return assets.queryAssetAddresses("a1", null, null, null, null);
    }

    /** The kept rows as address:verified. */
    private static List<String> kept(final AssetAddressV2Result result) {
        List<String> kept = new ArrayList<>();
        for (AssetAddressV2 address : result.getAddresses()) {
            kept.add(address.getAddress() + ":" + address.isVerified());
        }
        return kept;
    }

    private static List<String> excludedIds(final AssetAddressV2Result result) {
        List<String> ids = new ArrayList<>();
        for (ExcludedWhitelistedAddress excluded : result.getExcludedUnverified()) {
            ids.add(excluded.getId());
        }
        return ids;
    }

    private static String onlyReason(final AssetAddressV2Result result) {
        assertEquals(1, result.getExcludedUnverified().size(), result.getExcludedUnverified().toString());
        return result.getExcludedUnverified().get(0).getReason();
    }

    // --- INTERNAL ---------------------------------------------------------------------

    @Test
    @DisplayName("an INTERNAL row the verified managed address confirms is kept, verified, with its address")
    void internalConfirmed() throws Exception {
        holders.add(holder("party::8", INTERNAL, "8", null));
        manage("8", "party::8", "party::8");

        AssetAddressV2Result result = query();

        assertEquals(Collections.singletonList("party::8:true"), kept(result));
        assertTrue(result.getExcludedUnverified().isEmpty());
        assertEquals(Collections.singletonList("8"), stub.requestsTo(MANAGED).get(0).values("addressIds"));
        assertEquals("1", stub.requestsTo(MANAGED).get(0).param("limit"));
        // Exclusions and re-reads never move the holders cursor.
        assertTrue(result.getPage().hasMore());
        assertEquals(NEXT, result.getPage().getNextCursor());
    }

    @Test
    @DisplayName("an INTERNAL row whose address differs from the verified managed address is excluded")
    void internalDifferingAddress() throws Exception {
        holders.add(holder("party::7", EXTERNAL, null, null));
        holders.add(holder("party::8", INTERNAL, "8", null));
        manage("8", "party::other", "party::other");

        AssetAddressV2Result result = query();

        assertEquals(Collections.singletonList("party::7:false"), kept(result));
        assertEquals(Collections.singletonList("8"), excludedIds(result));
        assertTrue(onlyReason(result).contains("differs from the verified managed address 8"), onlyReason(result));
        assertTrue(result.getPage().hasMore());
        assertEquals(NEXT, result.getPage().getNextCursor());
    }

    @Test
    @DisplayName("an INTERNAL row the verified managed-address read does not return is excluded")
    void internalNotReturned() throws Exception {
        holders.add(holder("party::7", EXTERNAL, null, null));
        holders.add(holder("party::9", INTERNAL, "9", null));

        AssetAddressV2Result result = query();

        assertEquals(Collections.singletonList("party::7:false"), kept(result));
        assertEquals(Collections.singletonList("9"), excludedIds(result));
        assertTrue(onlyReason(result).contains("did not return addressID 9"), onlyReason(result));
    }

    @Test
    @DisplayName("an INTERNAL row without an addressID is excluded under its address, with no re-read")
    void internalMissingId() throws Exception {
        holders.add(holder("party::7", EXTERNAL, null, null));
        holders.add(holder("party::5", INTERNAL, null, null));
        holders.add(holder("party::6", INTERNAL, "0", null));

        AssetAddressV2Result result = query();

        assertEquals(Collections.singletonList("party::7:false"), kept(result));
        assertEquals(Arrays.asList("party::5", "party::6"), excludedIds(result));
        assertTrue(result.getExcludedUnverified().get(0).getReason().contains("no addressID"));
        assertTrue(stub.requestsTo(MANAGED).isEmpty(), "no id to re-read, so no managed-address request");
        assertTrue(stub.requestsTo(RULES).isEmpty());
    }

    @Test
    @DisplayName("an INTERNAL row whose HSM signature fails in the verified reader is excluded")
    void internalBadHsmSignature() throws Exception {
        holders.add(holder("party::8", INTERNAL, "8", null));
        holders.add(holder("party::12", INTERNAL, "12", null));
        manage("8", "party::8", "party::8");
        manage("12", "party::12", "party::forged");

        AssetAddressV2Result result = query();

        assertEquals(Collections.singletonList("party::8:true"), kept(result));
        assertEquals(Collections.singletonList("12"), excludedIds(result));
        assertTrue(onlyReason(result).contains("managed address did not verify"), onlyReason(result));
        assertTrue(onlyReason(result).contains("signature verification failed"), onlyReason(result));
    }

    @Test
    @DisplayName("51 INTERNAL ids are re-read in exactly two requests: 50 then 1")
    void internalChunks() throws Exception {
        for (int i = 1; i <= 51; i++) {
            String id = String.valueOf(i);
            holders.add(holder("party::" + id, INTERNAL, id, null));
            manage(id, "party::" + id, "party::" + id);
        }

        AssetAddressV2Result result = query();

        assertEquals(51, result.getAddresses().size());
        assertTrue(result.getExcludedUnverified().isEmpty());
        List<StubTransport.Recorded> reads = stub.requestsTo(MANAGED);
        assertEquals(2, reads.size());
        assertEquals(50, reads.get(0).values("addressIds").size());
        assertEquals("50", reads.get(0).param("limit"));
        assertEquals(Collections.singletonList("51"), reads.get(1).values("addressIds"));
        assertEquals("1", reads.get(1).param("limit"));
        // One container fetch for the whole page.
        assertEquals(1, stub.requestsTo(RULES).size());
    }

    // --- WHITELISTED ------------------------------------------------------------------

    @Test
    @DisplayName("a WHITELISTED row the verified whitelisted address confirms is kept, verified, with its address")
    void whitelistedConfirmed() throws Exception {
        holders.add(holder("ADDR3", WHITELISTED, null, "3"));
        whitelist("3", false);

        AssetAddressV2Result result = query();

        assertEquals(Collections.singletonList("ADDR3:true"), kept(result));
        assertTrue(result.getExcludedUnverified().isEmpty());
        List<StubTransport.Recorded> reads = stub.requestsTo(WHITELIST);
        assertEquals(1, reads.size());
        StubTransport.Recorded read = reads.get(0);
        assertEquals(Collections.singletonList("3"), read.values("ids"));
        assertEquals("1", read.param("limit"));
        assertEquals("true", read.param("rulesContainerNormalized"));
        // The default list: a holder is an approved whitelisted address, not a pending one.
        assertNull(read.param("includeForApproval"));
        assertTrue(stub.requestsTo(MANAGED).isEmpty());
        assertTrue(result.getPage().hasMore());
        assertEquals(NEXT, result.getPage().getNextCursor());
    }

    @Test
    @DisplayName("a WHITELISTED row whose address differs from the verified whitelisted address is excluded")
    void whitelistedDifferingAddress() throws Exception {
        holders.add(holder("party::7", EXTERNAL, null, null));
        holders.add(holder("OTHER3", WHITELISTED, null, "3"));
        whitelist("3", false);

        AssetAddressV2Result result = query();

        assertEquals(Collections.singletonList("party::7:false"), kept(result));
        assertEquals(Collections.singletonList("3"), excludedIds(result));
        assertTrue(onlyReason(result).contains("differs from the verified whitelisted address 3"), onlyReason(result));
    }

    @Test
    @DisplayName("a WHITELISTED row the verified whitelisted-address read does not return is excluded")
    void whitelistedNotReturned() throws Exception {
        holders.add(holder("party::7", EXTERNAL, null, null));
        holders.add(holder("ADDR4", WHITELISTED, null, "4"));

        AssetAddressV2Result result = query();

        assertEquals(Collections.singletonList("party::7:false"), kept(result));
        assertEquals(Collections.singletonList("4"), excludedIds(result));
        assertTrue(onlyReason(result).contains("did not return whitelistedAddressID 4"), onlyReason(result));
    }

    @Test
    @DisplayName("a WHITELISTED row without a whitelistedAddressID is excluded under its address, with no re-read")
    void whitelistedMissingId() throws Exception {
        holders.add(holder("party::7", EXTERNAL, null, null));
        holders.add(holder("ADDR5", WHITELISTED, null, null));

        AssetAddressV2Result result = query();

        assertEquals(Collections.singletonList("party::7:false"), kept(result));
        assertEquals(Collections.singletonList("ADDR5"), excludedIds(result));
        assertTrue(onlyReason(result).contains("no whitelistedAddressID"), onlyReason(result));
        assertTrue(stub.requestsTo(WHITELIST).isEmpty());
    }

    @Test
    @DisplayName("a WHITELISTED row whose envelope fails the six-step verification is excluded")
    void whitelistedTampered() throws Exception {
        holders.add(holder("ADDR3", WHITELISTED, null, "3"));
        holders.add(holder("ADDR6", WHITELISTED, null, "6"));
        whitelist("3", false);
        whitelist("6", true);

        AssetAddressV2Result result = query();

        assertEquals(Collections.singletonList("ADDR3:true"), kept(result));
        assertEquals(Collections.singletonList("6"), excludedIds(result));
        assertTrue(onlyReason(result).contains("whitelisted address did not verify"), onlyReason(result));
    }

    @Test
    @DisplayName("101 WHITELISTED ids are re-read in exactly two requests: 100 then 1")
    void whitelistedChunks() throws Exception {
        for (int i = 1; i <= 101; i++) {
            String id = String.valueOf(i);
            holders.add(holder("ADDR" + id, WHITELISTED, null, id));
            whitelist(id, false);
        }

        AssetAddressV2Result result = query();

        assertEquals(101, result.getAddresses().size());
        List<StubTransport.Recorded> reads = stub.requestsTo(WHITELIST);
        assertEquals(2, reads.size());
        assertEquals(100, reads.get(0).values("ids").size());
        assertEquals(Collections.singletonList("101"), reads.get(1).values("ids"));
    }

    // --- other types, whole-page outcomes ---------------------------------------------

    @Test
    @DisplayName("EXTERNAL and untyped rows are kept unverified, with no extra request")
    void externalKeptUnverified() throws Exception {
        holders.add(holder("party::7", EXTERNAL, null, null));
        holders.add(holder("party::8", null, "8", "3"));

        AssetAddressV2Result result = query();

        assertEquals(Arrays.asList("party::7:false", "party::8:false"), kept(result));
        assertTrue(result.getExcludedUnverified().isEmpty());
        assertEquals(1, stub.requests().size(), "only the holders page is read");
        assertEquals(HOLDERS, stub.requests().get(0).path());
    }

    @Test
    @DisplayName("rows came back and none survived: the call fails rather than return an empty page")
    void everyRowExcludedThrows() throws Exception {
        holders.add(holder("party::8", INTERNAL, "8", null));
        holders.add(holder("ADDR4", WHITELISTED, null, "4"));
        manage("8", "party::8", "party::forged");

        IntegrityException e = assertThrows(IntegrityException.class, this::query);
        assertTrue(e.getMessage().contains("all 2 asset address(es) failed verification"), e.getMessage());
    }

    @Test
    @DisplayName("an empty page is empty, with no extra request")
    void emptyPage() throws Exception {
        AssetAddressV2Result result = query();

        assertTrue(result.getAddresses().isEmpty());
        assertTrue(result.getExcludedUnverified().isEmpty());
        assertEquals(1, stub.requests().size());
    }

    // --- the verified readers fail closed ---------------------------------------------

    @Test
    @DisplayName("a rules container without an HSMSLOT key fails the call instead of excluding every row")
    void noHsmKeyFailsTheCall() throws Exception {
        noHsmKey = true;
        holders.add(holder("party::7", EXTERNAL, null, null));
        holders.add(holder("party::8", INTERNAL, "8", null));
        manage("8", "party::8", "party::8");

        assertThrows(ContainerIntegrityException.class, this::query);
        assertTrue(stub.requestsTo(MANAGED).isEmpty(), "no address read against a container that cannot verify it");
    }

    @Test
    @DisplayName("a failing managed-address read fails the call")
    void managedReadErrorPropagates() throws Exception {
        managedStatus = 503;
        holders.add(holder("party::7", EXTERNAL, null, null));
        holders.add(holder("party::8", INTERNAL, "8", null));

        ApiException e = assertThrows(ApiException.class, this::query);
        assertEquals(503, e.getCode());
        assertFalse(stub.requestsTo(MANAGED).isEmpty());
    }
}
