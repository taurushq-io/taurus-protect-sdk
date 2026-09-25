package com.taurushq.sdk.protect.client.testutil;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

/**
 * The shared-vector loader must hard-fail: a skipped gate reports parity that is not there.
 */
class SharedVectorsTest {

    @Test
    void aMissingFileFailsLoudly() {
        IllegalStateException e = assertThrows(IllegalStateException.class,
                () -> SharedVectors.load("no-such-vectors.json"));
        assertTrue(e.getMessage().contains("no-such-vectors.json"), e.getMessage());
    }

    @Test
    void aCountMismatchFailsLoudly() {
        IllegalStateException e = assertThrows(IllegalStateException.class,
                () -> SharedVectors.load("pagination-vectors.json",
                        "operations", 53, "offset", 25, "cursor", 12, "page_size", 14, "offset_input", 4));
        assertTrue(e.getMessage().contains("offset"), e.getMessage());
    }
}
