package com.taurushq.sdk.protect.client.service;

import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.taurushq.sdk.protect.client.mapper.ApiResponseCursorMapper;
import com.taurushq.sdk.protect.client.model.CursorPage;
import com.taurushq.sdk.protect.client.model.IntegrityException;
import com.taurushq.sdk.protect.client.model.OffsetPagination;
import com.taurushq.sdk.protect.client.model.OffsetRule;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.client.testutil.SharedVectors;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetWalletTokensReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import org.junit.jupiter.api.DynamicTest;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestFactory;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Locale;
import java.util.Map;
import java.util.Set;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assertions.fail;

/**
 * The shared pagination vectors ({@code scripts/resources/pagination-vectors.json}): the
 * same reply must give the same pagination value in all four SDKs.
 * <p>
 * The builders are called directly with the vector's request and reply. The cursor family's
 * reply cursor and the token family's reply go through the GENERATED Gson models first
 * ({@code fromJson}), so a token crosses the same {@code byte[]} boundary it does on a real
 * call.
 */
class PaginationVectorsTest {

    private static final String FILE = "pagination-vectors.json";

    /**
     * Paged operations this SDK does not call, with the reason. Everything else in the file
     * must be in {@link PagedOperation}, and vice versa.
     */
    static final Map<String, String> NOT_WRAPPED = new HashMap<>();

    static {
        String absent = "no Java method; methods absent from one SDK are out of scope (TODOS.md)";
        NOT_WRAPPED.put("ExchangeService_GetExchanges", absent);
        NOT_WRAPPED.put("StatisticsService_GetAggregatedTagStats", absent);
        NOT_WRAPPED.put("StatisticsService_GetPortfolioStatisticsHistory", absent);
        NOT_WRAPPED.put("TaurusNetworkService_GetSharedAssets", absent);
        NOT_WRAPPED.put("TaurusNetworkService_GetLendingAgreementsForApproval", absent);
        NOT_WRAPPED.put("TaurusNetworkService_GetPledgeActions", absent);
        NOT_WRAPPED.put("TaurusNetworkService_GetPledgeActionsForApproval", absent);
        NOT_WRAPPED.put("TaurusNetworkService_GetSettlementsForApproval", absent);
        NOT_WRAPPED.put("PriceService_ExportPricesHistory", absent);
    }

    private static JsonObject doc() {
        return SharedVectors.load(FILE,
                "operations", 53, "offset", 26, "cursor", 12, "page_size", 14, "offset_input", 4);
    }

    @Test
    void constantsMatchTheContract() {
        JsonObject constants = doc().getAsJsonObject("constants");
        assertEquals(constants.get("default_page_size").getAsInt(), Pagination.DEFAULT_PAGE_SIZE);
        assertEquals(constants.get("max_page_size").getAsInt(), Pagination.MAX_PAGE_SIZE);
        assertEquals(constants.get("price_history_max").getAsInt(), Pagination.MAX_PRICE_HISTORY_LIMIT);
    }

    @TestFactory
    List<DynamicTest> operationsUseTheirRule() {
        JsonObject ops = doc().getAsJsonObject("operations");
        List<DynamicTest> tests = new ArrayList<>();
        Map<String, PagedOperation> byId = new HashMap<>();
        for (PagedOperation op : PagedOperation.values()) {
            byId.put(op.operationId(), op);
        }
        Set<String> seen = new HashSet<>();
        for (Map.Entry<String, JsonElement> e : ops.entrySet()) {
            String id = e.getKey();
            JsonObject spec = e.getValue().getAsJsonObject();
            seen.add(id);
            tests.add(DynamicTest.dynamicTest(id, () -> {
                if (NOT_WRAPPED.containsKey(id)) {
                    assertTrue(!byId.containsKey(id), id + " is listed as not wrapped but the SDK calls it");
                    return;
                }
                PagedOperation op = byId.get(id);
                if (op == null) {
                    fail(id + " is a paged operation the SDK neither calls nor lists as not wrapped");
                }
                assertEquals(expectedKind(spec), describe(op), id);
                assertEquals(spec.get("total").getAsBoolean(), op.hasTotal(), id + " total");
                if (spec.has("size")) {
                    Integer max = "price_history".equals(spec.get("size").getAsString())
                            ? Integer.valueOf(Pagination.MAX_PRICE_HISTORY_LIMIT) : null;
                    assertEquals(max, op.maxSize(), id + " size maximum");
                } else {
                    assertEquals(Integer.valueOf(Pagination.MAX_PAGE_SIZE), op.maxSize(), id + " size maximum");
                }
            }));
        }
        tests.add(DynamicTest.dynamicTest("no SDK operation is missing from the file", () -> {
            for (PagedOperation op : PagedOperation.values()) {
                assertTrue(seen.contains(op.operationId()), op.operationId() + " is not in " + FILE);
            }
        }));
        return tests;
    }

    private static String expectedKind(final JsonObject spec) {
        String rule = spec.get("rule").getAsString();
        switch (rule) {
            case "cursor":
                return "CURSOR";
            case "token":
                return "TOKEN";
            case "limit_only":
                return "LIMIT_ONLY";
            default:
                return "OFFSET/" + rule.toUpperCase(Locale.ROOT);
        }
    }

    private static String describe(final PagedOperation op) {
        return op.kind() == PagedOperation.Kind.OFFSET ? "OFFSET/" + op.offsetRule().name() : op.kind().name();
    }

    @TestFactory
    List<DynamicTest> offsetVectors() {
        List<DynamicTest> tests = new ArrayList<>();
        for (JsonElement element : doc().getAsJsonArray("offset")) {
            JsonObject v = element.getAsJsonObject();
            tests.add(DynamicTest.dynamicTest(v.get("description").getAsString(), () -> {
                JsonObject request = v.getAsJsonObject("request");
                JsonObject reply = v.getAsJsonObject("reply");
                OffsetRule rule = OffsetRule.valueOf(v.get("rule").getAsString().toUpperCase(Locale.ROOT));
                int limit = request.get("limit").getAsInt();
                long offset = request.get("offset").getAsLong();
                int served = v.get("served_rows").getAsInt();
                int excluded = v.get("sdk_excluded").getAsInt();
                String total = text(reply, "totalItems");
                String replyOffset = text(reply, "offset");
                if (v.has("expect_error")) {
                    assertThrows(IntegrityException.class, () -> OffsetPagination.of(
                            rule, limit, offset, served, excluded, total, replyOffset));
                    return;
                }
                JsonObject expect = v.getAsJsonObject("expect");
                OffsetPagination p = OffsetPagination.of(rule, limit, offset, served, excluded, total, replyOffset);
                assertEquals(new OffsetPagination(expect.get("limit").getAsInt(), expect.get("offset").getAsLong(),
                        expect.get("total_items").getAsLong(), expect.get("next_offset").getAsLong(),
                        expect.get("has_more").getAsBoolean()), p);
            }));
        }
        return tests;
    }

    @TestFactory
    List<DynamicTest> cursorVectors() {
        List<DynamicTest> tests = new ArrayList<>();
        for (JsonElement element : doc().getAsJsonArray("cursor")) {
            JsonObject v = element.getAsJsonObject();
            tests.add(DynamicTest.dynamicTest(v.get("description").getAsString(), () -> {
                int pageSize = v.get("page_size").getAsInt();
                boolean hasTotal = v.get("has_total").getAsBoolean();
                String total = text(v, "reply_total");
                if (v.has("expect_error")) {
                    assertThrows(IntegrityException.class, () -> page(v, pageSize, total, hasTotal));
                    return;
                }
                JsonObject expect = v.getAsJsonObject("expect");
                Long expectedTotal = expect.get("total_items").isJsonNull()
                        ? null : expect.get("total_items").getAsLong();
                assertEquals(new CursorPage(expect.get("page_size").getAsInt(),
                        expect.get("next_cursor").getAsString(), expect.get("has_more").getAsBoolean(),
                        expectedTotal), page(v, pageSize, total, hasTotal));
            }));
        }
        return tests;
    }

    private static CursorPage page(final JsonObject v, final int pageSize, final String total,
                                   final boolean hasTotal) throws Exception {
        if ("token".equals(v.get("family").getAsString())) {
            // Through the generated byte[] field, as on a real wallet-tokens call.
            JsonObject reply = new JsonObject();
            if (!v.get("reply_token").isJsonNull()) {
                reply.add("next", v.get("reply_token"));
            }
            if (total != null) {
                reply.addProperty("total", total);
            }
            TgvalidatordGetWalletTokensReply dto = TgvalidatordGetWalletTokensReply.fromJson(reply.toString());
            return CursorPage.fromToken(pageSize, CursorRequest.tokenText(dto.getNext()), dto.getTotal(), hasTotal);
        }
        JsonElement cursor = v.get("reply_cursor");
        TgvalidatordResponseCursor dto = cursor.isJsonNull() ? null : TgvalidatordResponseCursor.fromJson(cursor.toString());
        return CursorPage.fromCursor(pageSize, ApiResponseCursorMapper.INSTANCE.fromDTO(dto), total, hasTotal);
    }

    @TestFactory
    List<DynamicTest> pageSizeVectors() {
        List<DynamicTest> tests = new ArrayList<>();
        for (JsonElement element : doc().getAsJsonArray("page_size")) {
            JsonObject v = element.getAsJsonObject();
            tests.add(DynamicTest.dynamicTest(v.get("description").getAsString(), () -> {
                String kind = v.get("kind").getAsString();
                Integer max = "page".equals(kind) ? Integer.valueOf(Pagination.MAX_PAGE_SIZE)
                        : "price_history".equals(kind) ? Integer.valueOf(Pagination.MAX_PRICE_HISTORY_LIMIT) : null;
                Integer input = v.get("input").isJsonNull() ? null : v.get("input").getAsInt();
                if (v.has("expect_error")) {
                    IllegalArgumentException e = assertThrows(IllegalArgumentException.class,
                            () -> Pagination.resolveSize("pageSize", input, max));
                    assertTrue(e.getMessage().startsWith("pageSize "), "the error names the option: " + e.getMessage());
                    return;
                }
                assertEquals(v.get("expect").getAsInt(), Pagination.resolveSize("pageSize", input, max));
            }));
        }
        return tests;
    }

    @TestFactory
    List<DynamicTest> offsetInputVectors() {
        List<DynamicTest> tests = new ArrayList<>();
        for (JsonElement element : doc().getAsJsonArray("offset_input")) {
            JsonObject v = element.getAsJsonObject();
            tests.add(DynamicTest.dynamicTest(v.get("description").getAsString(), () -> {
                // An unset offset is 0 in this SDK's signatures (a primitive long).
                long input = v.get("input").isJsonNull() ? 0L : v.get("input").getAsLong();
                if (v.has("expect_error")) {
                    IllegalArgumentException e = assertThrows(IllegalArgumentException.class,
                            () -> Pagination.resolveOffset("offset", input));
                    assertTrue(e.getMessage().startsWith("offset "), e.getMessage());
                    return;
                }
                assertEquals(v.get("expect").getAsLong(), Pagination.resolveOffset("offset", input));
            }));
        }
        return tests;
    }

    private static String text(final JsonObject o, final String key) {
        return o.has(key) && !o.get(key).isJsonNull() ? o.get(key).getAsString() : null;
    }
}
