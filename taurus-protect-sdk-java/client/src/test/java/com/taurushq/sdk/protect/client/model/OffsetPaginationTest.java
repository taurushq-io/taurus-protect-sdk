package com.taurushq.sdk.protect.client.model;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * Properties of the offset builder the shared vectors do not isolate.
 */
class OffsetPaginationTest {

    @Test
    @DisplayName("hasMore uses the server total, not the total reduced by exclusions")
    void hasMoreUsesTheServerTotal() {
        // 20 served of 21, two withheld: the reduced total (19) is below the next offset (20),
        // but row 21 still exists on the server and must stay reachable.
        OffsetPagination p = OffsetPagination.of(OffsetRule.PLUS_SERVER_ROWS, 20, 0, 20, 2, "21", null);
        assertEquals(19L, p.getTotalItems());
        assertEquals(20L, p.getNextOffset());
        assertTrue(p.hasMore());
    }

    @Test
    @DisplayName("a page that makes no progress ends the walk")
    void noProgressEndsTheWalk() {
        OffsetPagination p = OffsetPagination.of(OffsetRule.REPLY_OFFSET, 20, 40, 0, 0, "100", "40");
        assertEquals(40L, p.getNextOffset());
        assertFalse(p.hasMore());
    }

    @Test
    @DisplayName("a malformed count is an integrity failure, not a zero")
    void malformedCount() {
        assertThrows(IntegrityException.class,
                () -> OffsetPagination.of(OffsetRule.PLUS_ROWS, 20, 0, 0, 0, "1e3", null));
    }
}
