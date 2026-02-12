package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Reservation;
import com.taurushq.sdk.protect.client.model.ReservationUtxo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordReservation;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordUTXO;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting reservation-related DTOs to domain models.
 */
@Mapper(uses = {CurrencyMapper.class})
public interface ReservationMapper {

    ReservationMapper INSTANCE = Mappers.getMapper(ReservationMapper.class);

    /**
     * Converts a Reservation DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    @Mapping(source = "addressid", target = "addressId")
    Reservation fromDTO(TgvalidatordReservation dto);

    /**
     * Converts a list of Reservation DTOs to domain models.
     *
     * @param dtos the DTOs to convert
     * @return the domain models
     */
    List<Reservation> fromDTOList(List<TgvalidatordReservation> dtos);

    /**
     * Converts a UTXO DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    ReservationUtxo fromUtxoDTO(TgvalidatordUTXO dto);
}
