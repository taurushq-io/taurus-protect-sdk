package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Exchange;
import com.taurushq.sdk.protect.client.model.ExchangeCounterparty;
import com.taurushq.sdk.protect.client.model.ExchangeWithdrawalFee;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordExchange;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordExchangeCounterparty;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetExchangeWithdrawalFeeReply;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting exchange OpenAPI DTOs to client model objects.
 * <p>
 * This mapper handles the conversion of exchange data from OpenAPI generated models
 * to the SDK's clean domain models.
 */
@Mapper(uses = CurrencyMapper.class)
public interface ExchangeMapper {

    /**
     * Singleton instance of the mapper.
     */
    ExchangeMapper INSTANCE = Mappers.getMapper(ExchangeMapper.class);

    /**
     * Maps an exchange from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    @Mapping(target = "hasWLA", source = "hasWLA")
    Exchange fromDTO(TgvalidatordExchange dto);

    /**
     * Maps a list of exchanges from DTOs.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<Exchange> fromDTOList(List<TgvalidatordExchange> dtos);

    /**
     * Maps an exchange counterparty from DTO.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    ExchangeCounterparty fromCounterpartyDTO(TgvalidatordExchangeCounterparty dto);

    /**
     * Maps a list of exchange counterparties from DTOs.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<ExchangeCounterparty> fromCounterpartyDTOList(List<TgvalidatordExchangeCounterparty> dtos);

    /**
     * Maps a withdrawal fee from the reply.
     *
     * @param reply the OpenAPI reply
     * @return the domain model
     */
    @Mapping(target = "fee", source = "result")
    ExchangeWithdrawalFee fromWithdrawalFeeReply(TgvalidatordGetExchangeWithdrawalFeeReply reply);
}
