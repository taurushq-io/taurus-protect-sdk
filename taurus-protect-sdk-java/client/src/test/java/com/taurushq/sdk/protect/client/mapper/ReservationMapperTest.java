package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Reservation;
import com.taurushq.sdk.protect.client.model.ReservationUtxo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordReservation;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordUTXO;
import org.junit.jupiter.api.Test;

import java.time.OffsetDateTime;
import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class ReservationMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        OffsetDateTime creationDate = OffsetDateTime.now();

        TgvalidatordCurrency currencyDto = new TgvalidatordCurrency();
        currencyDto.setId("BTC");
        currencyDto.setName("Bitcoin");

        TgvalidatordReservation dto = new TgvalidatordReservation();
        dto.setId("res-123");
        dto.setAmount("1.5");
        dto.setCreationDate(creationDate);
        dto.setKind("utxo");
        dto.setComment("Test reservation");
        dto.setAddressid("addr-456");
        dto.setAddress("bc1qtest...");
        dto.setCurrencyInfo(currencyDto);
        dto.setResourceId("tx-789");
        dto.setResourceType("transaction");

        Reservation reservation = ReservationMapper.INSTANCE.fromDTO(dto);

        assertEquals("res-123", reservation.getId());
        assertEquals("1.5", reservation.getAmount());
        assertEquals(creationDate, reservation.getCreationDate());
        assertEquals("utxo", reservation.getKind());
        assertEquals("Test reservation", reservation.getComment());
        assertEquals("addr-456", reservation.getAddressId());
        assertEquals("bc1qtest...", reservation.getAddress());
        assertNotNull(reservation.getCurrencyInfo());
        assertEquals("BTC", reservation.getCurrencyInfo().getId());
        assertEquals("tx-789", reservation.getResourceId());
        assertEquals("transaction", reservation.getResourceType());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordReservation dto = new TgvalidatordReservation();
        dto.setId("minimal-res");

        Reservation reservation = ReservationMapper.INSTANCE.fromDTO(dto);

        assertEquals("minimal-res", reservation.getId());
        assertNull(reservation.getAmount());
        assertNull(reservation.getCreationDate());
        assertNull(reservation.getKind());
        assertNull(reservation.getComment());
        assertNull(reservation.getAddressId());
        assertNull(reservation.getCurrencyInfo());
    }

    @Test
    void fromDTO_handlesNullDto() {
        Reservation reservation = ReservationMapper.INSTANCE.fromDTO(null);
        assertNull(reservation);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordReservation dto1 = new TgvalidatordReservation();
        dto1.setId("res-1");
        dto1.setAmount("1.0");

        TgvalidatordReservation dto2 = new TgvalidatordReservation();
        dto2.setId("res-2");
        dto2.setAmount("2.0");

        List<Reservation> reservations = ReservationMapper.INSTANCE.fromDTOList(Arrays.asList(dto1, dto2));

        assertNotNull(reservations);
        assertEquals(2, reservations.size());
        assertEquals("res-1", reservations.get(0).getId());
        assertEquals("1.0", reservations.get(0).getAmount());
        assertEquals("res-2", reservations.get(1).getId());
        assertEquals("2.0", reservations.get(1).getAmount());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<Reservation> reservations = ReservationMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(reservations);
        assertTrue(reservations.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<Reservation> reservations = ReservationMapper.INSTANCE.fromDTOList(null);
        assertNull(reservations);
    }

    @Test
    void fromUtxoDTO_mapsFields() {
        TgvalidatordUTXO dto = new TgvalidatordUTXO();
        dto.setId("utxo-123");
        dto.setHash("abc123hash");
        dto.setValue("1000000");
        dto.setScript("script-data");
        dto.setBlockHeight("700000");
        dto.setReservedByRequestId("req-456");
        dto.setReservationId("res-789");

        ReservationUtxo utxo = ReservationMapper.INSTANCE.fromUtxoDTO(dto);

        assertNotNull(utxo);
        assertEquals("utxo-123", utxo.getId());
        assertEquals("abc123hash", utxo.getHash());
        assertEquals("1000000", utxo.getValue());
        assertEquals("script-data", utxo.getScript());
        assertEquals("700000", utxo.getBlockHeight());
        assertEquals("req-456", utxo.getReservedByRequestId());
        assertEquals("res-789", utxo.getReservationId());
    }

    @Test
    void fromUtxoDTO_handlesNullDto() {
        ReservationUtxo utxo = ReservationMapper.INSTANCE.fromUtxoDTO(null);
        assertNull(utxo);
    }

    @Test
    void fromDTO_mapsAddressIdCorrectly() {
        TgvalidatordReservation dto = new TgvalidatordReservation();
        dto.setId("res-test");
        dto.setAddressid("internal-addr-123");

        Reservation reservation = ReservationMapper.INSTANCE.fromDTO(dto);

        assertEquals("internal-addr-123", reservation.getAddressId());
    }
}
