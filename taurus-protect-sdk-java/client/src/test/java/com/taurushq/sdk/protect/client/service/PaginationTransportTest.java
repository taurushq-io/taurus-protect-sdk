package com.taurushq.sdk.protect.client.service;

import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.mapper.ApiExceptionMapper;
import com.taurushq.sdk.protect.client.model.ActionResult;
import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.AssetAddressesResult;
import com.taurushq.sdk.protect.client.model.BalanceResult;
import com.taurushq.sdk.protect.client.model.CursorPage;
import com.taurushq.sdk.protect.client.model.FeePayerResult;
import com.taurushq.sdk.protect.client.model.GovernanceRulesHistoryResult;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.OffsetPagination;
import com.taurushq.sdk.protect.client.model.PageRequest;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.client.model.RequestResult;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAddressEnvelope;
import com.taurushq.sdk.protect.client.model.SignedWhitelistedAssetEnvelope;
import com.taurushq.sdk.protect.client.model.User;
import com.taurushq.sdk.protect.client.model.UserResult;
import com.taurushq.sdk.protect.client.model.WalletResult;
import com.taurushq.sdk.protect.client.model.WalletTokensResult;
import com.taurushq.sdk.protect.client.model.WebhookResult;
import com.taurushq.sdk.protect.client.model.WhitelistMetadata;
import com.taurushq.sdk.protect.client.model.WhitelistedAddressApproval;
import com.taurushq.sdk.protect.client.model.WhitelistedAddressListResult;
import com.taurushq.sdk.protect.client.model.WhitelistedAssetApproval;
import com.taurushq.sdk.protect.client.model.WhitelistedAssetResult;
import com.taurushq.sdk.protect.client.testutil.SignedWhitelist;
import com.taurushq.sdk.protect.client.testutil.StubTransport;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.security.KeyPairGenerator;
import java.security.PrivateKey;
import java.security.spec.ECGenParameterSpec;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Each list family end to end through the transport stub: what goes on the wire, and the
 * pagination value built from the reply. The shared vectors pin the builders and the
 * request shape; these pin the wiring between them on real replies.
 */
class PaginationTransportTest {

    private static PrivateKey approverKey;

    @BeforeAll
    static void keys() throws Exception {
        KeyPairGenerator generator = KeyPairGenerator.getInstance("EC");
        generator.initialize(new ECGenParameterSpec("secp256r1"));
        approverKey = generator.generateKeyPair().getPrivate();
    }

    private static String rows(final String field, final int count, final int firstId) {
        StringBuilder json = new StringBuilder("\"" + field + "\":[");
        for (int i = 0; i < count; i++) {
            json.append(i == 0 ? "" : ",").append("{\"id\":\"").append(firstId + i).append("\"}");
        }
        return json.append(']').toString();
    }

    // ---- offset lists ------------------------------------------------------------------

    @Test
    @DisplayName("wallets: the reply offset IS the next page's offset")
    void walletsReplyOffsetIsTheNextOffset() throws Exception {
        StubTransport stub = StubTransport.replying("{" + rows("result", 2, 1) + ",\"totalItems\":\"5\",\"offset\":\"2\"}");
        WalletResult page = new ServicesUnderTest(stub).wallets().getWallets(2, 0);

        assertEquals(Collections.singletonList(Arrays.asList("limit", "2")), stub.only().query());
        assertEquals(2, page.getWallets().size());
        assertEquals(new OffsetPagination(2, 0, 5, 2, true), page.getPagination());
    }

    @Test
    @DisplayName("wallets: the last page reports no more, with the reply offset as next")
    void walletsLastPage() throws Exception {
        StubTransport stub = StubTransport.replying("{" + rows("result", 1, 5) + ",\"totalItems\":\"5\",\"offset\":\"5\"}");
        WalletResult page = new ServicesUnderTest(stub).wallets().getWallets(2, 4);

        assertEquals("4", stub.only().param("offset"));
        assertEquals(new OffsetPagination(2, 4, 5, 5, false), page.getPagination());
    }

    @Test
    @DisplayName("an empty page {} carries a full pagination value, never null")
    void emptyOffsetPage() throws Exception {
        StubTransport stub = StubTransport.replying("{}");
        ServicesUnderTest services = new ServicesUnderTest(stub);

        WalletResult wallets = services.wallets().getWallets(0, 0);
        assertTrue(wallets.getWallets().isEmpty());
        assertEquals(new OffsetPagination(20, 0, 0, 0, false), wallets.getPagination());

        FeePayerResult feePayers = services.feePayers().getFeePayers();
        assertEquals(new OffsetPagination(20, 0, 0, 0, false), feePayers.getPagination());
    }

    @Test
    @DisplayName("actions: the reply total reaches the result")
    void actionsCarryTheTotal() throws Exception {
        StubTransport stub = StubTransport.replying("{\"result\":[],\"totalItems\":\"7\"}");
        ActionResult actions = new ServicesUnderTest(stub).actions().getActions(5, 5, null);

        assertEquals(7L, actions.getPagination().getTotalItems());
        assertEquals(5L, actions.getPagination().getNextOffset());
        assertFalse(actions.getPagination().hasMore(), "an empty page makes no progress and ends the walk");
    }

    @Test
    @DisplayName("users: a synthetic row beyond the limit does not push the next offset past it")
    void usersSyntheticRow() throws Exception {
        StubTransport stub = StubTransport.replying("{" + rows("result", 3, 1) + ",\"totalItems\":\"10\"}");
        UserResult page = new ServicesUnderTest(stub).users().getUsers(2, 0);

        assertEquals(3, page.getUsers().size());
        assertEquals(new OffsetPagination(2, 0, 10, 2, true), page.getPagination());
    }

    @Test
    @DisplayName("getUsersByEmail walks every page in pages of 100")
    void usersByEmailWalksThePages() throws Exception {
        StubTransport stub = StubTransport.replying(
                "{" + rows("result", 100, 1) + ",\"totalItems\":\"150\"}",
                "{" + rows("result", 50, 101) + ",\"totalItems\":\"150\"}");
        List<User> users = new ServicesUnderTest(stub).users()
                .getUsersByEmail(Collections.singletonList("someone@example.invalid"));

        assertEquals(150, users.size());
        assertEquals(2, stub.requests().size());
        assertEquals("100", stub.requests().get(0).param("limit"));
        assertNull(stub.requests().get(0).param("offset"));
        assertEquals("100", stub.requests().get(1).param("limit"));
        assertEquals("100", stub.requests().get(1).param("offset"));
    }

    // ---- cursor lists ------------------------------------------------------------------

    @Test
    @DisplayName("a continuation sends the cursor URL-encoded, NEXT and the page size")
    void cursorContinuation() throws Exception {
        String cursor = "eyJhIjoiKy8ifQ+/cQ==";
        StubTransport stub = StubTransport.replying("{\"cursor\":{\"currentPage\":\"bmV4dA==\",\"hasNext\":true}}");
        RequestResult page = new ServicesUnderTest(stub).requests().getRequestsForApproval(7, cursor);

        StubTransport.Recorded request = stub.only();
        assertEquals(cursor, request.param("cursor.currentPage"));
        assertTrue(request.rawQuery().contains("cursor.currentPage=eyJhIjoiKy8ifQ%2B%2FcQ%3D%3D"), request.rawQuery());
        assertEquals("NEXT", request.param("cursor.pageRequest"));
        assertEquals("7", request.param("cursor.pageSize"));
        assertEquals(new CursorPage(7, "bmV4dA==", true, null), page.getPage());
        assertTrue(page.hasNext());
    }

    @Test
    @DisplayName("a v2 last page: currentPage present, hasNext omitted, so no next cursor")
    void cursorLastPage() throws Exception {
        StubTransport stub = StubTransport.replying("{\"result\":[],\"cursor\":{\"currentPage\":\"abc\"}}");
        RequestResult page = new ServicesUnderTest(stub).requests().getRequests(null, null, null, null, 20, "prev");

        assertEquals(new CursorPage(20, "", false, null), page.getPage());
        assertNull(page.nextCursor(20), "no request cursor past the last page");
    }

    @Test
    @DisplayName("webhooks: no next cursor on the last page (it used to ignore hasNext)")
    void webhookNextCursorHonoursHasNext() throws Exception {
        StubTransport stub = StubTransport.replying("{\"webhooks\":[],\"cursor\":{\"currentPage\":\"abc\",\"hasNext\":false}}");
        WebhookResult page = new ServicesUnderTest(stub).webhooks().getWebhooks(null, null, 20, null);

        assertNull(page.nextCursor(20));
        assertFalse(page.getPage().hasMore());
        assertEquals("", page.getPage().getNextCursor());
    }

    @Test
    @DisplayName("hasNext without a currentPage is a malformed reply")
    void cursorWithoutToken() {
        StubTransport stub = StubTransport.replying("{\"cursor\":{\"hasNext\":true}}");
        assertThrows(IntegrityException.class,
                () -> new ServicesUnderTest(stub).changes().getChanges(null, null, 20, null));
    }

    @Test
    @DisplayName("balances: requestCursor only, and the page carries the server total")
    void balancesCarryTheTotal() throws Exception {
        StubTransport stub = StubTransport.replying(
                "{\"balances\":[],\"total\":\"57\",\"cursor\":{\"currentPage\":\"cDI=\",\"hasNext\":true}}");
        BalanceResult page = new ServicesUnderTest(stub).balances().getBalances(null, 20, null);

        assertEquals(Collections.singletonList(Arrays.asList("requestCursor.pageSize", "20")), stub.only().query());
        assertEquals(new CursorPage(20, "cDI=", true, 57L), page.getPage());
    }

    @Test
    @DisplayName("asset addresses: page two is reachable through the body request cursor")
    void assetAddressesPageTwo() throws Exception {
        StubTransport stub = StubTransport.replying("{\"addresses\":[],\"totalItems\":\"41\"}");
        AssetAddressesResult page = new ServicesUnderTest(stub).assets()
                .getAssetAddresses("ETH", null, null, 20, "cDI=");

        JsonObject body = JsonParser.parseString(stub.only().body()).getAsJsonObject();
        JsonObject requestCursor = body.getAsJsonObject("requestCursor");
        assertEquals("cDI=", requestCursor.get("currentPage").getAsString());
        assertEquals("NEXT", requestCursor.get("pageRequest").getAsString());
        assertEquals("20", requestCursor.get("pageSize").getAsString());
        assertFalse(body.has("limit") || body.has("cursor"), "the legacy limit/cursor is never sent: " + body);
        assertEquals(new CursorPage(20, "", false, 41L), page.getPage());
    }

    @Test
    @DisplayName("reservations send pageRequest NEXT and a page size, not currentPage alone")
    void reservationsContinuation() throws Exception {
        StubTransport stub = StubTransport.replying("{}");
        new ServicesUnderTest(stub).reservations().getReservations(null, null, null, null, 20, "cDI=");

        assertEquals(Arrays.asList(Arrays.asList("cursor.currentPage", "cDI="), Arrays.asList("cursor.pageRequest", "NEXT"),
                Arrays.asList("cursor.pageSize", "20")), stub.only().query());
    }

    @Test
    @DisplayName("the low-level request cursor still reaches the wire as given")
    void lowLevelCursor() throws Exception {
        StubTransport stub = StubTransport.replying("{}");
        new ServicesUnderTest(stub).changes().getChanges(null, null,
                new ApiRequestCursor("tok", PageRequest.PREVIOUS, 30));

        assertEquals(Arrays.asList(Arrays.asList("cursor.currentPage", "tok"), Arrays.asList("cursor.pageRequest", "PREVIOUS"),
                Arrays.asList("cursor.pageSize", "30")), stub.only().query());
    }

    @Test
    @DisplayName("an unset low-level cursor sends the default page size")
    void unsetLowLevelCursor() throws Exception {
        StubTransport stub = StubTransport.replying("{}");
        new ServicesUnderTest(stub).changes().getChangesForApproval(new ApiRequestCursor());

        assertEquals(Collections.singletonList(Arrays.asList("cursor.pageSize", "20")), stub.only().query());
    }

    // ---- token lists -------------------------------------------------------------------

    @Test
    @DisplayName("wallet tokens: the byte token comes back as text and goes back escaped")
    void walletTokensToken() throws Exception {
        StubTransport stub = StubTransport.replying("{\"balances\":[],\"next\":\"ab+/cQ==\",\"total\":\"120\"}");
        ServicesUnderTest services = new ServicesUnderTest(stub);

        WalletTokensResult first = services.wallets().getWalletTokens(7, 20, null);
        assertEquals(new CursorPage(20, "ab+/cQ==", true, 120L), first.getPage());

        services.wallets().getWalletTokens(7, 20, first.getPage().getNextCursor());
        StubTransport.Recorded continuation = stub.requests().get(1);
        assertEquals("/api/rest/v1/wallets/7/tokens", continuation.path());
        assertTrue(continuation.rawQuery().contains("cursor=ab%2B%2FcQ%3D%3D"), continuation.rawQuery());
        assertEquals("20", continuation.param("limit"));
    }

    @Test
    @DisplayName("rules history: the cursor round-trips and exclusions reduce the total")
    void rulesHistoryExclusions() throws Exception {
        // One ruleset carries no SuperAdmin signature, so it is withheld and named; the
        // page total drops by one while the next cursor stays the server's.
        StubTransport stub = StubTransport.replying(
                "{\"result\":[{\"rulesContainer\":\"CgA=\"}],\"cursor\":\"YWJj\",\"totalItems\":\"3\"}");
        GovernanceRulesHistoryResult page = new ServicesUnderTest(stub).governance().getRulesHistory(20, "eHl6");

        assertEquals("eHl6", stub.only().param("cursor"));
        assertTrue(page.getRules().isEmpty());
        assertEquals(1, page.getExcludedUnverified().size());
        assertEquals(new CursorPage(20, "YWJj", true, 2L), page.getPage());
    }

    @Test
    @DisplayName("a token cursor that is not base64 is refused before any request")
    void badTokenCursor() {
        StubTransport stub = StubTransport.replying("{}");
        assertThrows(IllegalArgumentException.class,
                () -> new ServicesUnderTest(stub).governance().getRulesHistory(20, "%%%"));
        assertEquals(0, stub.requests().size());
    }

    // ---- page-size validation ----------------------------------------------------------

    @Test
    @DisplayName("page sizes: 0 and null mean the default, 101 and -1 are refused by name")
    void pageSizeBounds() {
        assertEquals(20, Pagination.page(null, null).getPageSize());
        assertEquals(20, Pagination.page(0, null).getPageSize());
        IllegalArgumentException above = assertThrows(IllegalArgumentException.class, () -> Pagination.page(101, null));
        assertTrue(above.getMessage().contains("pageSize") && above.getMessage().contains("100"), above.getMessage());
        assertThrows(IllegalArgumentException.class, () -> Pagination.page(-1, "x"));
        assertEquals(PageRequest.NEXT, Pagination.page(20, "x").getPageRequest());
        assertNull(Pagination.page(20, null).getPageRequest(), "a first page sends no page request");
    }

    // ---- approval re-reads -------------------------------------------------------------

    private static List<Long> ids(final int count) {
        List<Long> ids = new ArrayList<>();
        for (long id = 1; id <= count; id++) {
            ids.add(id);
        }
        return ids;
    }

    private static int count(final StubTransport.Recorded request, final String name) {
        int n = 0;
        for (List<String> pair : request.query()) {
            if (name.equals(pair.get(0))) {
                n++;
            }
        }
        return n;
    }

    @Test
    @DisplayName("an address approval re-reads its ids in chunks of 100")
    void addressApprovalChunks() {
        List<SignedWhitelistedAddressEnvelope> envelopes = new ArrayList<>();
        for (Long id : ids(150)) {
            SignedWhitelistedAddressEnvelope e = new SignedWhitelistedAddressEnvelope();
            e.setId(id);
            WhitelistMetadata metadata = new WhitelistMetadata();
            metadata.setHash("hash-" + id);
            e.setMetadata(metadata);
            envelopes.add(e);
        }
        WhitelistedAddressApproval selection = new WhitelistedAddressListResult(envelopes,
                Collections.emptyList(), new OffsetPagination(100, 0, 150, 100, true)).selectAll();

        StubTransport stub = StubTransport.replying("{}");
        // The re-read returns nothing, so the approval refuses to sign after reading.
        assertThrows(IntegrityException.class, () -> new ServicesUnderTest(stub).whitelistedAddresses()
                .approveWhitelistedAddresses(selection, approverKey, "reviewed"));

        assertEquals(2, stub.requests().size());
        assertEquals(100, count(stub.requests().get(0), "ids"));
        assertEquals("100", stub.requests().get(0).param("limit"));
        assertEquals(50, count(stub.requests().get(1), "ids"));
        assertEquals("50", stub.requests().get(1).param("limit"));
        assertEquals("true", stub.requests().get(1).param("includeForApproval"));
    }

    @Test
    @DisplayName("whitelisted addresses: a failing row is excluded, the rest verify, and the walk is not shifted")
    void whitelistedAddressesExclusionKeepsTheWalk() throws Exception {
        SignedWhitelist whitelist = new SignedWhitelist();
        StubTransport stub = StubTransport.replying(SignedWhitelist.page(Arrays.asList(
                whitelist.row("1", false), whitelist.row("2", true), whitelist.row("3", false)), "5"));

        WhitelistedAddressListResult page = new WhitelistedAddressService(stub.client(),
                new ApiExceptionMapper(), whitelist.superAdminKeys(), 1)
                .getWhitelistedAddresses(3, 0);

        List<String> kept = new ArrayList<>();
        for (SignedWhitelistedAddressEnvelope envelope : page.getEnvelopes()) {
            kept.add(envelope.getId() + "=" + envelope.getWhitelistedAddress().getAddress());
        }
        assertEquals(Arrays.asList("1=ADDR1", "3=ADDR3"), kept);
        assertEquals(1, page.getExcludedUnverified().size());
        assertEquals("2", page.getExcludedUnverified().get(0).getId());
        // The next page starts after every row the SERVER returned; only the total drops.
        assertEquals(new OffsetPagination(3, 0, 4, 3, true), page.getPagination());
    }

    @Test
    @DisplayName("an address approval whose re-read row fails verification names it and signs nothing")
    void addressApprovalNamesTheFailedRow() throws Exception {
        SignedWhitelist whitelist = new SignedWhitelist();
        SignedWhitelistedAddressEnvelope reviewed = new SignedWhitelistedAddressEnvelope();
        reviewed.setId(6L);
        WhitelistMetadata metadata = new WhitelistMetadata();
        metadata.setHash(SignedWhitelist.hashOf("6"));
        reviewed.setMetadata(metadata);
        WhitelistedAddressApproval selection = new WhitelistedAddressListResult(
                Collections.singletonList(reviewed), Collections.emptyList(),
                new OffsetPagination(1, 0, 1, 1, false)).selectAll();

        StubTransport stub = StubTransport.replying(
                SignedWhitelist.page(Collections.singletonList(whitelist.row("6", true)), "1"));
        IntegrityException e = assertThrows(IntegrityException.class, () -> new WhitelistedAddressService(
                stub.client(), new ApiExceptionMapper(),
                whitelist.superAdminKeys(), 1).approveWhitelistedAddresses(selection, approverKey, "reviewed"));

        assertTrue(e.getMessage().contains("address 6 failed verification"), e.getMessage());
        assertEquals(1, stub.requests().size(), "the approval itself is never sent");
    }

    @Test
    @DisplayName("an asset approval re-reads its ids in chunks of 100")
    void assetApprovalChunks() {
        List<SignedWhitelistedAssetEnvelope> envelopes = new ArrayList<>();
        for (Long id : ids(150)) {
            SignedWhitelistedAssetEnvelope e = new SignedWhitelistedAssetEnvelope();
            e.setId(id);
            WhitelistMetadata metadata = new WhitelistMetadata();
            metadata.setHash("hash-" + id);
            e.setMetadata(metadata);
            envelopes.add(e);
        }
        WhitelistedAssetResult read = new WhitelistedAssetResult();
        read.setAssets(envelopes);
        WhitelistedAssetApproval selection = read.selectAll();

        StubTransport stub = StubTransport.replying("{}");
        assertThrows(IntegrityException.class, () -> new ServicesUnderTest(stub).whitelistedAssets()
                .approveWhitelistedAssets(selection, approverKey, "reviewed"));

        // The contract list names its id filter whitelistedContractAddressIds on the wire.
        assertEquals(2, stub.requests().size());
        assertEquals(100, count(stub.requests().get(0), "whitelistedContractAddressIds"));
        assertEquals("100", stub.requests().get(0).param("limit"));
        assertEquals(50, count(stub.requests().get(1), "whitelistedContractAddressIds"));
        assertEquals("50", stub.requests().get(1).param("limit"));
    }

    // ---- endpoints ---------------------------------------------------------------------

    @Test
    @DisplayName("fees come from the v2 endpoint, in the v2 shape")
    void feesV2() throws Exception {
        StubTransport stub = StubTransport.replying(
                "{\"result\":[{\"currencyId\":\"ETH\",\"value\":\"21000\",\"denom\":\"wei\"}]}");
        List<com.taurushq.sdk.protect.client.model.Fee> fees = new ServicesUnderTest(stub).fees().getFees();

        assertEquals("/api/rest/v2/fees", stub.only().path());
        assertEquals("ETH", fees.get(0).getCurrencyId());
        assertEquals("wei", fees.get(0).getDenom());
    }

    @Test
    @DisplayName("ERC token metadata comes from the EVM endpoint")
    void tokenMetadataEvm() throws Exception {
        StubTransport stub = StubTransport.replying("{}");
        new ServicesUnderTest(stub).tokenMetadata().getEVMERCTokenMetadata("mainnet", "0xabc", "1", false, "ETH");

        assertEquals("/api/rest/v1/evm/mainnet/erc/contract/0xabc/token/1/metadata", stub.only().path());
    }
}
