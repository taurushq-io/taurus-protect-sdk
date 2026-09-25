package com.taurushq.sdk.protect.client.service;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.client.testutil.SharedVectors;
import com.taurushq.sdk.protect.client.testutil.StubTransport;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.TestFactory;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.Collections;
import java.util.HashMap;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.TreeSet;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.fail;

/**
 * The shared list-request vectors ({@code scripts/resources/list-request-vectors.json}): the
 * same canonical options must put the same request on the wire in all four SDKs.
 * <p>
 * Each vector runs the real service method through the {@link StubTransport} with a
 * {@code {}} reply, so the generated client builds the request. The adapter table maps every
 * operation to its Java method and every canonical option to its parameter: an operation or
 * option the table does not map fails, except the operations this SDK does not wrap,
 * listed by name in {@link PaginationVectorsTest}, and options the Java signature cannot
 * take at all, which are acceptable only where the vector expects the call to be refused.
 */
class ListRequestVectorsTest {

    private static final String FILE = "list-request-vectors.json";

    /**
     * One vector's call.
     */
    @FunctionalInterface
    interface Call {
        void invoke(ServicesUnderTest s, Options o) throws Exception;
    }

    private static final class Adapter {
        private final Set<String> options;
        private final Set<String> unexpressible;
        private final Call call;

        Adapter(final Set<String> options, final Set<String> unexpressible, final Call call) {
            this.options = options;
            this.unexpressible = unexpressible;
            this.call = call;
        }
    }

    private static final Map<String, Adapter> ADAPTERS = new LinkedHashMap<>();

    private static final List<String> OFFSET = Arrays.asList("limit", "offset");
    private static final List<String> CURSOR = Arrays.asList("page_size", "cursor");

    private static void map(final String op, final List<String> options, final Call call) {
        ADAPTERS.put(op, new Adapter(new HashSet<>(options), Collections.<String>emptySet(), call));
    }

    private static List<String> with(final List<String> base, final String... more) {
        List<String> all = new ArrayList<>(base);
        all.addAll(Arrays.asList(more));
        return all;
    }

    static {
        map("ActionService_GetActions", OFFSET, (s, o) -> s.actions().getActions(o.limit(), o.offset(), null));
        map("AssetServiceV2_ListAssetOperationsV2", CURSOR, (s, o) -> s.assets().listAssetOperations(
                o.path("assetID"), null, null, o.pageSize(), o.cursor()));
        map("AssetServiceV2_QueryAssetAddressesV2", CURSOR, (s, o) -> s.assets().queryAssetAddresses(
                o.path("assetID"), null, null, o.pageSize(), o.cursor()));
        map("AssetServiceV2_QueryAssetsV2", CURSOR, (s, o) -> s.assets().queryAssets(
                null, null, null, null, null, null, o.pageSize(), o.cursor()));
        map("AuditService_GetAuditTrails", CURSOR, (s, o) -> s.audit().getAuditTrails(
                null, null, null, null, null, o.pageSize(), o.cursor()));
        map("ChangeService_GetChanges", CURSOR, (s, o) -> s.changes().getChanges(null, null, o.pageSize(), o.cursor()));
        map("ChangeService_GetChangesForApproval", CURSOR,
                (s, o) -> s.changes().getChangesForApproval(o.pageSize(), o.cursor()));
        map("EarnService_GetRewards", CURSOR, (s, o) -> s.earn().listRewards(null, o.pageSize(), o.cursor()));
        map("FeePayerService_GetFeePayers", OFFSET,
                (s, o) -> s.feePayers().getFeePayers(o.limit(), o.offset(), null, null, null));
        map("FiatProviderService_GetFiatProviderAccounts", with(CURSOR, "provider", "label"),
                (s, o) -> s.fiat().getFiatProviderAccounts(o.str("provider"), o.str("label"), null, null,
                        o.pageSize(), o.cursor()));
        map("FiatProviderService_GetFiatProviderCounterpartyAccounts", with(CURSOR, "provider", "label"),
                (s, o) -> s.fiat().getFiatProviderCounterpartyAccounts(o.str("provider"), o.str("label"), null,
                        null, o.pageSize(), o.cursor()));
        map("FiatProviderService_GetFiatProviderEntities", CURSOR,
                (s, o) -> s.fiat().listFiatProviderEntities(null, null, null, o.pageSize(), o.cursor()));
        map("FiatProviderService_GetFiatProviderOperations", CURSOR,
                (s, o) -> s.fiat().getFiatProviderOperations(null, null, null, o.pageSize(), o.cursor()));
        map("PriceService_GetPricesHistory", Collections.singletonList("limit"),
                (s, o) -> s.prices().getPriceHistory(o.path("base"), o.path("quote"), o.limit()));
        map("PriceService_QueryPricesV2", with(CURSOR, "from_currency_id", "to_currency_ids"),
                (s, o) -> s.prices().getPrices(o.str("from_currency_id"), o.list("to_currency_ids"), null, null,
                        o.pageSize(), o.cursor()));
        // The approval queue has no status or external-request-id filter in this SDK's signature
        // (validatord's approval paginator drops the latter), so both vectors that must be
        // refused are refused by the compiler.
        ADAPTERS.put("RequestService_GetRequestsForApprovalV2", new Adapter(new HashSet<>(CURSOR),
                new HashSet<>(Arrays.asList("statuses", "external_request_ids")),
                (s, o) -> s.requests().getRequestsForApproval(o.pageSize(), o.cursor())));
        map("RequestService_GetRequestsV2", CURSOR,
                (s, o) -> s.requests().getRequests(null, null, null, null, o.pageSize(), o.cursor()));
        map("RuleService_GetBusinessRulesV2", CURSOR,
                (s, o) -> s.businessRules().getBusinessRules(o.pageSize(), o.cursor()));
        map("RuleService_GetRulesHistory", CURSOR, (s, o) -> s.governance().getRulesHistory(o.pageSize(), o.cursor()));
        map("StakingService_GetStakeAccounts", CURSOR,
                (s, o) -> s.staking().getStakeAccounts(null, null, null, o.pageSize(), o.cursor()));
        map("TaurusNetworkService_GetLendingAgreements", CURSOR,
                (s, o) -> s.lending().getLendingAgreements(null, o.pageSize(), o.cursor()));
        map("TaurusNetworkService_GetLendingOffers", CURSOR,
                (s, o) -> s.lending().getLendingOffers(null, null, null, null, o.pageSize(), o.cursor()));
        map("TaurusNetworkService_GetPledges", CURSOR,
                (s, o) -> s.pledges().list(null, null, null, null, null, o.pageSize(), o.cursor()));
        map("TaurusNetworkService_GetPledgesWithdrawals", CURSOR,
                (s, o) -> s.pledges().listWithdrawals(null, null, null, o.pageSize(), o.cursor()));
        map("TaurusNetworkService_GetSettlements", CURSOR,
                (s, o) -> s.settlements().getSettlements(null, null, null, o.pageSize(), o.cursor()));
        // Seven filters leave no room for positional page arguments: the page is one cursor.
        map("TaurusNetworkService_GetSharedAddresses", CURSOR, (s, o) -> s.sharing().listSharedAddresses(
                null, null, null, null, null, null, null, Pagination.page(o.pageSize(), o.cursor())));
        map("TransactionService_ExportTransactions", Collections.singletonList("limit"),
                (s, o) -> s.transactions().exportTransactions(null, null, null, null, o.limit()));
        map("TransactionService_GetTransactions", OFFSET,
                (s, o) -> s.transactions().getTransactions(null, null, null, null, o.limit(), o.offset()));
        map("UserService_GetGroups", OFFSET, (s, o) -> s.groups().getGroups(o.limit(), o.offset(), null, null, null));
        map("UserService_GetUsers", OFFSET, (s, o) -> s.users().getUsers(o.limit(), o.offset()));
        map("WalletService_GetAddresses", with(OFFSET, "exclude_disabled"),
                (s, o) -> s.addresses().getAddresses((Long) null, o.limit(), o.offset(), o.excludeDisabled()));
        map("WalletService_GetAssetAddresses", with(CURSOR, "currency"),
                (s, o) -> s.assets().getAssetAddresses(o.str("currency"), null, null, o.pageSize(), o.cursor()));
        map("WalletService_GetAssetWallets", with(CURSOR, "currency"),
                (s, o) -> s.assets().getAssetWallets(o.str("currency"), o.pageSize(), o.cursor()));
        map("WalletService_GetBalances", CURSOR, (s, o) -> s.balances().getBalances(null, o.pageSize(), o.cursor()));
        map("WalletService_GetNFTCollectionBalances", CURSOR,
                (s, o) -> s.balances().getNFTCollectionBalances(null, null, o.pageSize(), o.cursor()));
        map("WalletService_GetReservations", CURSOR,
                (s, o) -> s.reservations().getReservations(null, null, null, null, o.pageSize(), o.cursor()));
        map("WalletService_GetWalletTokens", CURSOR, (s, o) -> s.wallets().getWalletTokens(
                Long.parseLong(o.path("id")), o.pageSize(), o.cursor()));
        map("WalletService_GetWalletsV2", with(OFFSET, "exclude_disabled"),
                (s, o) -> s.wallets().getWallets(o.limit(), o.offset(), o.excludeDisabled()));
        map("WebhookService_GetWebhookCalls", CURSOR,
                (s, o) -> s.webhookCalls().getWebhookCalls(null, null, null, null, o.pageSize(), o.cursor()));
        map("WebhookService_GetWebhooks", CURSOR,
                (s, o) -> s.webhooks().getWebhooks(null, null, o.pageSize(), o.cursor()));
        map("WhitelistService_GetWhitelistedAddresses", OFFSET,
                (s, o) -> s.whitelistedAddresses().getWhitelistedAddresses(o.limit(), o.offset()));
        map("WhitelistService_GetWhitelistedAddressesForApproval", OFFSET,
                (s, o) -> s.whitelistedAddresses().getWhitelistedAddressesForApproval(o.limit(), o.offset(), null, null));
        map("WhitelistService_GetWhitelistedContracts", OFFSET,
                (s, o) -> s.whitelistedAssets().getWhitelistedAssets(o.limit(), o.offset()));
        map("WhitelistService_GetWhitelistedContractsForApproval", OFFSET,
                (s, o) -> s.whitelistedAssets().getWhitelistedAssetsForApproval(o.limit(), o.offset(), null));
    }

    private static JsonObject doc() {
        return SharedVectors.load(FILE, "methods", 53, "vectors", 278);
    }

    @TestFactory
    List<DynamicTest> everyVectorPutsTheExpectedRequestOnTheWire() {
        JsonObject doc = doc();
        JsonObject methods = doc.getAsJsonObject("methods");
        List<DynamicTest> tests = new ArrayList<>();
        for (JsonElement element : doc.getAsJsonArray("vectors")) {
            JsonObject v = element.getAsJsonObject();
            String op = v.get("operation").getAsString();
            tests.add(DynamicTest.dynamicTest(op + ": " + v.get("description").getAsString(),
                    () -> run(methods.getAsJsonObject(op), v)));
        }
        return tests;
    }

    private static void run(final JsonObject method, final JsonObject v) throws Exception {
        String op = v.get("operation").getAsString();
        Adapter adapter = ADAPTERS.get(op);
        if (adapter == null) {
            assertTrue(PaginationVectorsTest.NOT_WRAPPED.containsKey(op),
                    op + " has no adapter and is not listed as not wrapped");
            return;
        }
        Set<String> given = v.getAsJsonObject("options").keySet();
        Set<String> unmapped = new TreeSet<>(given);
        unmapped.removeAll(adapter.options);
        unmapped.removeAll(adapter.unexpressible);
        assertTrue(unmapped.isEmpty(), op + ": options with no adapter mapping " + unmapped);
        boolean refused = v.has("expect_error");
        for (String option : given) {
            if (adapter.unexpressible.contains(option)) {
                assertTrue(refused, op + ": option " + option + " cannot be expressed in Java but the vector "
                        + "expects it to reach the wire");
                return;
            }
        }

        StubTransport stub = StubTransport.replying("{}");
        ServicesUnderTest services = new ServicesUnderTest(stub);
        Options options = new Options(v);
        if (refused) {
            assertThrows(IllegalArgumentException.class, () -> adapter.call.invoke(services, options));
            assertEquals(0, stub.requests().size(), op + ": a refused call must send nothing");
        } else {
            adapter.call.invoke(services, options);
            StubTransport.Recorded request = stub.only();
            assertEquals(method.get("method").getAsString(), request.method());
            assertEquals(path(method.get("path").getAsString(), v.getAsJsonObject("path_params")), request.path());
            JsonObject expect = v.getAsJsonObject("expect");
            assertEquals(expectedQuery(expect), request.query(), op + " query");
            if (expect.has("body")) {
                assertEquals(expect.get("body"), JsonParser.parseString(request.body()), op + " body");
            } else {
                assertTrue(request.body() == null || request.body().isEmpty(), op + " sent a body " + request.body());
            }
        }
        assertEquals(given, options.consumed, op + ": options read by the adapter");
    }

    private static String path(final String template, final JsonObject params) {
        String path = template;
        for (Map.Entry<String, JsonElement> p : params.entrySet()) {
            path = path.replace("{" + p.getKey() + "}", p.getValue().getAsString());
        }
        return path;
    }

    private static List<List<String>> expectedQuery(final JsonObject expect) {
        List<List<String>> pairs = new ArrayList<>();
        if (expect.has("query")) {
            for (JsonElement pair : expect.getAsJsonArray("query")) {
                JsonArray a = pair.getAsJsonArray();
                pairs.add(Arrays.asList(a.get(0).getAsString(), a.get(1).getAsString()));
            }
        }
        return pairs;
    }

    @TestFactory
    List<DynamicTest> everyAdapterOptionIsExercised() {
        JsonObject doc = doc();
        Map<String, Set<String>> used = new HashMap<>();
        for (JsonElement element : doc.getAsJsonArray("vectors")) {
            JsonObject v = element.getAsJsonObject();
            used.computeIfAbsent(v.get("operation").getAsString(), k -> new HashSet<>())
                    .addAll(v.getAsJsonObject("options").keySet());
        }
        List<DynamicTest> tests = new ArrayList<>();
        for (Map.Entry<String, Adapter> e : ADAPTERS.entrySet()) {
            tests.add(DynamicTest.dynamicTest(e.getKey(), () -> {
                if (!doc.getAsJsonObject("methods").has(e.getKey())) {
                    fail(e.getKey() + " is adapted but is not in " + FILE);
                }
                Set<String> stale = new TreeSet<>(e.getValue().options);
                stale.removeAll(used.getOrDefault(e.getKey(), Collections.<String>emptySet()));
                assertTrue(stale.isEmpty(), e.getKey() + ": adapter options no vector uses " + stale);
            }));
        }
        return tests;
    }

    /**
     * A vector's options, typed, recording which ones the adapter read.
     */
    static final class Options {
        private final JsonObject options;
        private final JsonObject pathParams;
        private final Set<String> consumed = new HashSet<>();

        Options(final JsonObject v) {
            this.options = v.getAsJsonObject("options");
            this.pathParams = v.getAsJsonObject("path_params");
        }

        private JsonElement read(final String name) {
            if (!options.has(name) || options.get(name).isJsonNull()) {
                return null;
            }
            consumed.add(name);
            return options.get(name);
        }

        int limit() {
            JsonElement e = read("limit");
            return e == null ? 0 : e.getAsInt();
        }

        long offset() {
            JsonElement e = read("offset");
            return e == null ? 0L : e.getAsLong();
        }

        Integer pageSize() {
            JsonElement e = read("page_size");
            return e == null ? null : e.getAsInt();
        }

        String cursor() {
            return str("cursor");
        }

        Boolean excludeDisabled() {
            JsonElement e = read("exclude_disabled");
            return e == null ? null : e.getAsBoolean();
        }

        String str(final String name) {
            JsonElement e = read(name);
            return e == null ? null : e.getAsString();
        }

        List<String> list(final String name) {
            JsonElement e = read(name);
            if (e == null) {
                return null;
            }
            List<String> values = new ArrayList<>();
            for (JsonElement item : e.getAsJsonArray()) {
                values.add(item.getAsString());
            }
            return values;
        }

        String path(final String name) {
            if (!pathParams.has(name)) {
                throw new AssertionError("vector has no path parameter " + name);
            }
            return pathParams.get(name).getAsString();
        }
    }
}
