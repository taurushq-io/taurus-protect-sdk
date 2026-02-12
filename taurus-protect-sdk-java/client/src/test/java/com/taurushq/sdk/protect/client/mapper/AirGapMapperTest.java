package com.taurushq.sdk.protect.client.mapper;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertNotNull;

class AirGapMapperTest {

    @Test
    void instance_isNotNull() {
        assertNotNull(AirGapMapper.INSTANCE);
    }
}
