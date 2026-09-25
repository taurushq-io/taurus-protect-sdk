package com.taurushq.sdk.protect.client.model;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertSame;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Tests for {@link WhitelistedAssetResult} and the next-offset rule of its endpoint: skipped
 * rows keep their SQL slot, so the next page starts at offset + limit.
 */
class WhitelistedAssetResultTest {

    private static OffsetPagination page(final int limit, final long offset, final int rows, final String total) {
        return OffsetPagination.of(OffsetRule.PLUS_LIMIT, limit, offset, rows, 0, total, null);
    }

    @Test
    @DisplayName("hasMore is true while rows remain beyond the page")
    void hasMoreWhenRowsRemain() {
        assertTrue(page(100, 0, 100, "250").hasMore());
        assertTrue(page(100, 100, 100, "250").hasMore());
    }

    @Test
    @DisplayName("a short page is not the end: skipped rows keep their slot")
    void shortPageIsNotTheEnd() {
        OffsetPagination p = page(100, 0, 97, "250");
        assertTrue(p.hasMore());
        assertEquals(100L, p.getNextOffset());
    }

    @Test
    @DisplayName("hasMore is false on the last page")
    void hasMoreOnLastPage() {
        assertFalse(page(100, 200, 50, "250").hasMore());
        assertFalse(page(100, 250, 0, "250").hasMore());
    }

    @Test
    @DisplayName("offset plus page size cannot wrap into a false hasMore")
    void hasMoreDoesNotOverflow() {
        // The next offset is computed in long arithmetic: an int sum here would wrap to a
        // negative and promise a page that does not exist.
        assertFalse(page(1, Integer.MAX_VALUE, 0, "250").hasMore());
        assertFalse(page(100, Integer.MAX_VALUE, 0, "250").hasMore());
    }

    @Test
    @DisplayName("assets and pagination round-trip")
    void accessorsRoundTrip() {
        List<SignedWhitelistedAssetEnvelope> assets = new ArrayList<>();
        assets.add(new SignedWhitelistedAssetEnvelope());
        OffsetPagination pagination = page(20, 0, 1, "1");

        WhitelistedAssetResult result = new WhitelistedAssetResult();
        result.setAssets(assets);
        result.setPagination(pagination);

        assertEquals(1, result.getAssets().size());
        assertSame(pagination, result.getPagination());
        assertEquals(1L, result.getPagination().getTotalItems());
    }

    @Test
    @DisplayName("an empty page reports no more results")
    void emptyPageHasNoMore() {
        OffsetPagination p = page(50, 0, 0, null);
        assertFalse(p.hasMore());
        assertEquals(0L, p.getTotalItems());
    }
}
