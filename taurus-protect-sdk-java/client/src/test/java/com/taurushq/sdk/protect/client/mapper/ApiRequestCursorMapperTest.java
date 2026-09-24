package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;
import com.taurushq.sdk.protect.client.model.Pagination;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordRequestCursor;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class ApiRequestCursorMapperTest {

    @Test
    void toDTO_withCompleteData_mapsAllFields() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setCurrentPage("page-abc");
        cursor.setPageRequest(PageRequest.NEXT);
        cursor.setPageSize(50);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertNotNull(result);
        assertEquals("page-abc", result.getCurrentPage());
        assertEquals("NEXT", result.getPageRequest());
        assertEquals("50", result.getPageSize());
    }

    @Test
    void toDTO_withNullPageRequest_handlesGracefully() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setCurrentPage("page-1");
        cursor.setPageRequest(null);
        cursor.setPageSize(25);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertNotNull(result);
        assertEquals("page-1", result.getCurrentPage());
        assertNull(result.getPageRequest());
        assertEquals("25", result.getPageSize());
    }

    @Test
    void toDTO_withNullCursor_returnsNull() {
        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(null);

        // Then
        assertNull(result);
    }

    // ==================== PageRequest.FIRST Tests ====================

    @Test
    void toDTO_withFirstPageRequest_mapsToString() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setPageRequest(PageRequest.FIRST);
        cursor.setPageSize(10);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertEquals("FIRST", result.getPageRequest());
    }

    @Test
    void toDTO_withFirstPageRequest_withCurrentPage_mapsAllFields() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setCurrentPage("start-page");
        cursor.setPageRequest(PageRequest.FIRST);
        cursor.setPageSize(100);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertNotNull(result);
        assertEquals("start-page", result.getCurrentPage());
        assertEquals("FIRST", result.getPageRequest());
        assertEquals("100", result.getPageSize());
    }

    // ==================== PageRequest.PREVIOUS Tests ====================

    @Test
    void toDTO_withPreviousPageRequest_mapsToString() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setPageRequest(PageRequest.PREVIOUS);
        cursor.setPageSize(10);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertEquals("PREVIOUS", result.getPageRequest());
    }

    @Test
    void toDTO_withPreviousPageRequest_withCurrentPage_mapsAllFields() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setCurrentPage("page-5");
        cursor.setPageRequest(PageRequest.PREVIOUS);
        cursor.setPageSize(20);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertNotNull(result);
        assertEquals("page-5", result.getCurrentPage());
        assertEquals("PREVIOUS", result.getPageRequest());
        assertEquals("20", result.getPageSize());
    }

    // ==================== PageRequest.NEXT Tests ====================

    @Test
    void toDTO_withNextPageRequest_mapsToString() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setPageRequest(PageRequest.NEXT);
        cursor.setPageSize(10);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertEquals("NEXT", result.getPageRequest());
    }

    @Test
    void toDTO_withNextPageRequest_withCurrentPage_mapsAllFields() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setCurrentPage("page-3");
        cursor.setPageRequest(PageRequest.NEXT);
        cursor.setPageSize(30);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertNotNull(result);
        assertEquals("page-3", result.getCurrentPage());
        assertEquals("NEXT", result.getPageRequest());
        assertEquals("30", result.getPageSize());
    }

    // ==================== PageRequest.LAST Tests ====================

    @Test
    void toDTO_withLastPageRequest_mapsToString() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setPageRequest(PageRequest.LAST);
        cursor.setPageSize(10);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertEquals("LAST", result.getPageRequest());
    }

    @Test
    void toDTO_withLastPageRequest_withCurrentPage_mapsAllFields() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setCurrentPage("end-page");
        cursor.setPageRequest(PageRequest.LAST);
        cursor.setPageSize(50);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertNotNull(result);
        assertEquals("end-page", result.getCurrentPage());
        assertEquals("LAST", result.getPageRequest());
        assertEquals("50", result.getPageSize());
    }

    // ==================== All PageRequest Values Test ====================

    @Test
    void toDTO_allPageRequestValues_mapCorrectly() {
        for (PageRequest pageRequest : PageRequest.values()) {
            // Given
            ApiRequestCursor cursor = new ApiRequestCursor();
            cursor.setPageRequest(pageRequest);
            cursor.setPageSize(10);

            // When
            TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

            // Then
            assertEquals(pageRequest.name(), result.getPageRequest(),
                    "PageRequest." + pageRequest + " should map to \"" + pageRequest.name() + "\"");
        }
    }

    // ==================== Page Size Edge Cases ====================
    // An explicit page size must be between 1 and Pagination.MAX_PAGE_SIZE: the cursor used
    // to accept anything and send it, including 0, negatives and Long.MAX_VALUE.

    @Test
    void setPageSize_rejectsZero() {
        ApiRequestCursor cursor = new ApiRequestCursor();
        assertThrows(IllegalArgumentException.class, () -> cursor.setPageSize(0));
    }

    @Test
    void setPageSize_rejectsNegative() {
        ApiRequestCursor cursor = new ApiRequestCursor();
        assertThrows(IllegalArgumentException.class, () -> cursor.setPageSize(-1));
        assertThrows(IllegalArgumentException.class, () -> cursor.setPageSize(Long.MIN_VALUE));
    }

    @Test
    void setPageSize_rejectsAboveTheMaximum() {
        ApiRequestCursor cursor = new ApiRequestCursor();
        IllegalArgumentException e = assertThrows(IllegalArgumentException.class, () -> cursor.setPageSize(101));
        assertTrue(e.getMessage().contains("pageSize") && e.getMessage().contains("100"), e.getMessage());
        assertThrows(IllegalArgumentException.class, () -> cursor.setPageSize(Long.MAX_VALUE));
    }

    @Test
    void constructors_rejectAnOutOfRangePageSize() {
        assertThrows(IllegalArgumentException.class, () -> new ApiRequestCursor(PageRequest.FIRST, 0));
        assertThrows(IllegalArgumentException.class, () -> new ApiRequestCursor(PageRequest.FIRST, 101));
        assertThrows(IllegalArgumentException.class, () -> new ApiRequestCursor("page", PageRequest.NEXT, -1));
    }

    @Test
    void toDTO_atTheMaximumPageSize_mapsCorrectly() {
        ApiRequestCursor cursor = new ApiRequestCursor(PageRequest.FIRST, Pagination.MAX_PAGE_SIZE);

        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        assertEquals("100", result.getPageSize());
    }

    // ==================== Current Page Edge Cases ====================

    @Test
    void toDTO_withNullCurrentPage_mapsGracefully() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setCurrentPage(null);
        cursor.setPageRequest(PageRequest.FIRST);
        cursor.setPageSize(10);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertNotNull(result);
        assertNull(result.getCurrentPage());
        assertEquals("FIRST", result.getPageRequest());
    }

    @Test
    void toDTO_withEmptyCurrentPage_mapsEmptyString() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setCurrentPage("");
        cursor.setPageRequest(PageRequest.NEXT);
        cursor.setPageSize(25);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertNotNull(result);
        assertEquals("", result.getCurrentPage());
    }
}
