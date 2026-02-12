package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiRequestCursor;
import com.taurushq.sdk.protect.client.model.PageRequest;
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

    @Test
    void toDTO_withZeroPageSize_mapsToZeroString() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setPageRequest(PageRequest.FIRST);
        cursor.setPageSize(0);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertEquals("0", result.getPageSize());
    }

    @Test
    void toDTO_withLargePageSize_mapsCorrectly() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setPageRequest(PageRequest.FIRST);
        cursor.setPageSize(Long.MAX_VALUE);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertEquals(String.valueOf(Long.MAX_VALUE), result.getPageSize());
    }

    @Test
    void toDTO_withNegativePageSize_mapsCorrectly() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setPageRequest(PageRequest.FIRST);
        cursor.setPageSize(-1);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertEquals("-1", result.getPageSize());
    }

    @Test
    void toDTO_withMinLongPageSize_mapsCorrectly() {
        // Given
        ApiRequestCursor cursor = new ApiRequestCursor();
        cursor.setPageRequest(PageRequest.NEXT);
        cursor.setPageSize(Long.MIN_VALUE);

        // When
        TgvalidatordRequestCursor result = ApiResponseCursorMapper.INSTANCE.toDTO(cursor);

        // Then
        assertEquals(String.valueOf(Long.MIN_VALUE), result.getPageSize());
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
