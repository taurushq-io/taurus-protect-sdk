package com.taurushq.sdk.protect.client.model;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import java.util.ArrayList;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Tests for {@link WhitelistedAssetResult}.
 */
class WhitelistedAssetResultTest {

    @Test
    @DisplayName("hasMore is true while rows remain beyond the page")
    void hasMoreWhenRowsRemain() {
        WhitelistedAssetResult result = new WhitelistedAssetResult();
        result.setTotalItems(250L);

        assertTrue(result.hasMore(0, 100));
        assertTrue(result.hasMore(100, 100));
    }

    @Test
    @DisplayName("hasMore is false on the last page")
    void hasMoreOnLastPage() {
        WhitelistedAssetResult result = new WhitelistedAssetResult();
        result.setTotalItems(250L);

        assertFalse(result.hasMore(200, 100));
        assertFalse(result.hasMore(250, 100));
    }

    @Test
    @DisplayName("offset plus page size cannot wrap into a false hasMore")
    void hasMoreDoesNotOverflow() {
        // The unsafe form, (currentOffset + pageSize) < totalItems, overflows to a
        // negative here and promises a page that does not exist.
        WhitelistedAssetResult result = new WhitelistedAssetResult();
        result.setTotalItems(250L);

        assertFalse(result.hasMore(Integer.MAX_VALUE, 1));
        assertFalse(result.hasMore(Integer.MAX_VALUE, Integer.MAX_VALUE));
    }

    @Test
    @DisplayName("assets and totalItems round-trip")
    void accessorsRoundTrip() {
        List<SignedWhitelistedAssetEnvelope> assets = new ArrayList<>();
        assets.add(new SignedWhitelistedAssetEnvelope());

        WhitelistedAssetResult result = new WhitelistedAssetResult();
        result.setAssets(assets);
        result.setTotalItems(1L);

        assertEquals(1, result.getAssets().size());
        assertEquals(1L, result.getTotalItems());
    }

    @Test
    @DisplayName("an empty page reports no more results")
    void emptyPageHasNoMore() {
        WhitelistedAssetResult result = new WhitelistedAssetResult();
        result.setTotalItems(0L);

        assertFalse(result.hasMore(0, 50));
    }
}
