package com.taurushq.sdk.protect.client.model;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;

class RequestStatusTest {

    @Test
    void valueOfLabel() {
        assertEquals(RequestStatus.CREATED, RequestStatus.valueOfLabel("created"));
        assertEquals(RequestStatus.CREATED, RequestStatus.valueOfLabel("CREATED"));
        assertEquals(RequestStatus.FAST_APPROVED, RequestStatus.valueOfLabel("FAST_APPROVED"));
    }

    @Test
    void valueOfLabel_throwsOnUnrecognized() {
        assertThrows(IllegalArgumentException.class, () -> RequestStatus.valueOfLabel("test"));
    }

    @Test
    void valueOfLabel_throwsOnEmpty() {
        assertThrows(IllegalArgumentException.class, () -> RequestStatus.valueOfLabel(""));
    }

    @Test
    void valueOfLabel_throwsOnNull() {
        assertThrows(IllegalArgumentException.class, () -> RequestStatus.valueOfLabel(null));
    }
}
