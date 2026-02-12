package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.UserDevicePairing;
import com.taurushq.sdk.protect.client.model.UserDevicePairingInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateUserDevicePairingReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordUserDevicePairingInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordUserDevicePairingInfoStatus;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

class UserDeviceMapperTest {

    @Test
    void fromCreateDTO_mapsPairingId() {
        TgvalidatordCreateUserDevicePairingReply dto = new TgvalidatordCreateUserDevicePairingReply();
        dto.setPairingID("pairing-abc-123");

        UserDevicePairing pairing = UserDeviceMapper.INSTANCE.fromCreateDTO(dto);

        assertNotNull(pairing);
        assertEquals("pairing-abc-123", pairing.getPairingId());
    }

    @Test
    void fromCreateDTO_handlesNullDto() {
        UserDevicePairing pairing = UserDeviceMapper.INSTANCE.fromCreateDTO(null);
        assertNull(pairing);
    }

    @Test
    void fromInfoDTO_mapsAllFields() {
        TgvalidatordUserDevicePairingInfo dto = new TgvalidatordUserDevicePairingInfo();
        dto.setPairingID("pairing-xyz");
        dto.setApiKey("api-key-secret");
        dto.setStatus(TgvalidatordUserDevicePairingInfoStatus.APPROVED);

        UserDevicePairingInfo info = UserDeviceMapper.INSTANCE.fromInfoDTO(dto);

        assertNotNull(info);
        assertEquals("pairing-xyz", info.getPairingId());
        assertEquals("api-key-secret", info.getApiKey());
        assertEquals("APPROVED", info.getStatus());
    }

    @Test
    void fromInfoDTO_mapsStatusWaiting() {
        TgvalidatordUserDevicePairingInfo dto = new TgvalidatordUserDevicePairingInfo();
        dto.setPairingID("test-pairing");
        dto.setStatus(TgvalidatordUserDevicePairingInfoStatus.WAITING);

        UserDevicePairingInfo info = UserDeviceMapper.INSTANCE.fromInfoDTO(dto);

        assertEquals("WAITING", info.getStatus());
    }

    @Test
    void fromInfoDTO_mapsStatusPairing() {
        TgvalidatordUserDevicePairingInfo dto = new TgvalidatordUserDevicePairingInfo();
        dto.setPairingID("test-pairing");
        dto.setStatus(TgvalidatordUserDevicePairingInfoStatus.PAIRING);

        UserDevicePairingInfo info = UserDeviceMapper.INSTANCE.fromInfoDTO(dto);

        assertEquals("PAIRING", info.getStatus());
    }

    @Test
    void fromInfoDTO_handlesNullApiKey() {
        TgvalidatordUserDevicePairingInfo dto = new TgvalidatordUserDevicePairingInfo();
        dto.setPairingID("pairing-no-key");
        dto.setStatus(TgvalidatordUserDevicePairingInfoStatus.WAITING);

        UserDevicePairingInfo info = UserDeviceMapper.INSTANCE.fromInfoDTO(dto);

        assertNotNull(info);
        assertEquals("pairing-no-key", info.getPairingId());
        assertNull(info.getApiKey());
    }

    @Test
    void fromInfoDTO_handlesNullDto() {
        UserDevicePairingInfo info = UserDeviceMapper.INSTANCE.fromInfoDTO(null);
        assertNull(info);
    }
}
