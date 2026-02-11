package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.AddressInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddressInfo;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

/**
 * MapStruct mapper for converting address DTOs to domain models.
 * <p>
 * This mapper transforms OpenAPI-generated address objects into the SDK's
 * domain model objects. It uses {@link WalletMapper}, {@link BalanceMapper},
 * {@link ScoreMapper}, {@link AttributeMapper}, and {@link CurrencyMapper}
 * for nested object conversion.
 *
 * @see Address
 * @see AddressInfo
 * @see TgvalidatordAddress
 */
@Mapper(uses = {WalletMapper.class, BalanceMapper.class, ScoreMapper.class, AttributeMapper.class, ScoreMapper.class, CurrencyMapper.class})
public interface AddressMapper {

    /**
     * Singleton instance of the mapper.
     */
    AddressMapper INSTANCE = Mappers.getMapper(AddressMapper.class);

    /**
     * Converts an address DTO to a domain model.
     *
     * @param address the OpenAPI address DTO
     * @return the domain model address
     */
    Address fromDTO(TgvalidatordAddress address);

    /**
     * Converts an address info DTO to a domain model.
     * <p>
     * AddressInfo is a simplified representation of an address used in
     * transaction source/destination lists.
     *
     * @param addressInfo the OpenAPI address info DTO
     * @return the domain model address info
     */
    AddressInfo fromDTO(TgvalidatordAddressInfo addressInfo);
}


