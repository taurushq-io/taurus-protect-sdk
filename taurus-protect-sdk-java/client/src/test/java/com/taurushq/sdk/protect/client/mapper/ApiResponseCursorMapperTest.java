package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class ApiResponseCursorMapperTest {

    @Test
    void fromDTO_withCompleteData_mapsAllFields() {
        // Given
        TgvalidatordResponseCursor dto = new TgvalidatordResponseCursor();
        dto.setCurrentPage("page-token-123");
        dto.setHasPrevious(true);
        dto.setHasNext(true);

        // When
        ApiResponseCursor result = ApiResponseCursorMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals("page-token-123", result.getCurrentPage());
        assertTrue(result.hasPrevious());
        assertTrue(result.hasNext());
    }

    @Test
    void fromDTO_withNullFields_handlesGracefully() {
        // Given
        TgvalidatordResponseCursor dto = new TgvalidatordResponseCursor();
        // All fields left null

        // When
        ApiResponseCursor result = ApiResponseCursorMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertNull(result.getCurrentPage());
        assertFalse(result.hasPrevious());
        assertFalse(result.hasNext());
    }

    @Test
    void fromDTO_withAllFalse_mapsCorrectly() {
        // Given
        TgvalidatordResponseCursor dto = new TgvalidatordResponseCursor();
        dto.setCurrentPage("first-page");
        dto.setHasPrevious(false);
        dto.setHasNext(false);

        // When
        ApiResponseCursor result = ApiResponseCursorMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals("first-page", result.getCurrentPage());
        assertFalse(result.hasPrevious());
        assertFalse(result.hasNext());
    }

    @Test
    void fromDTO_withAllTrue_mapsCorrectly() {
        // Given
        TgvalidatordResponseCursor dto = new TgvalidatordResponseCursor();
        dto.setCurrentPage("middle-page");
        dto.setHasPrevious(true);
        dto.setHasNext(true);

        // When
        ApiResponseCursor result = ApiResponseCursorMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals("middle-page", result.getCurrentPage());
        assertTrue(result.hasPrevious());
        assertTrue(result.hasNext());
    }

    @Test
    void fromDTO_withNullDto_returnsNull() {
        // Given
        TgvalidatordResponseCursor dto = null;

        // When
        ApiResponseCursor result = ApiResponseCursorMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNull(result);
    }

    @Test
    void fromDTO_withEmptyCurrentPage_mapsEmptyString() {
        // Given
        TgvalidatordResponseCursor dto = new TgvalidatordResponseCursor();
        dto.setCurrentPage("");
        dto.setHasPrevious(false);
        dto.setHasNext(true);

        // When
        ApiResponseCursor result = ApiResponseCursorMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertEquals("", result.getCurrentPage());
    }

    @Test
    void fromDTO_withOnlyHasPrevious_mapsPartially() {
        // Given
        TgvalidatordResponseCursor dto = new TgvalidatordResponseCursor();
        dto.setHasPrevious(true);

        // When
        ApiResponseCursor result = ApiResponseCursorMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertNull(result.getCurrentPage());
        assertTrue(result.hasPrevious());
        assertFalse(result.hasNext());
    }

    @Test
    void fromDTO_withOnlyHasNext_mapsPartially() {
        // Given
        TgvalidatordResponseCursor dto = new TgvalidatordResponseCursor();
        dto.setHasNext(true);

        // When
        ApiResponseCursor result = ApiResponseCursorMapper.INSTANCE.fromDTO(dto);

        // Then
        assertNotNull(result);
        assertNull(result.getCurrentPage());
        assertFalse(result.hasPrevious());
        assertTrue(result.hasNext());
    }
}
