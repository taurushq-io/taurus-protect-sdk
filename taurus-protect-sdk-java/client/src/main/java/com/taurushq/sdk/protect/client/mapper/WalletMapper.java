package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Wallet;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWallet;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWalletInfo;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

/**
 * MapStruct mapper for converting wallet DTOs to domain models.
 * <p>
 * This mapper transforms OpenAPI-generated wallet objects into the SDK's
 * domain model objects. It uses {@link BalanceMapper}, {@link AttributeMapper},
 * and {@link CurrencyMapper} for nested object conversion.
 *
 * @see Wallet
 * @see TgvalidatordWalletInfo
 * @see TgvalidatordWallet
 */
@Mapper(uses = {BalanceMapper.class, AttributeMapper.class, CurrencyMapper.class})
public interface WalletMapper {

    /**
     * Singleton instance of the mapper.
     */
    WalletMapper INSTANCE = Mappers.getMapper(WalletMapper.class);

    /**
     * Converts a wallet info DTO to a domain model.
     * <p>
     * This method is used for responses from the V2 API endpoints.
     *
     * @param walletInfo the OpenAPI wallet info DTO
     * @return the domain model wallet
     */
    @Mapping(source = "isOmnibus", target = "omnibus")
    Wallet fromDTO(TgvalidatordWalletInfo walletInfo);

    /**
     * Converts a wallet DTO to a domain model.
     * <p>
     * Note: The {@code network} and {@code visibilityGroupID} fields are not
     * populated from this DTO type and will be null in the result.
     *
     * @param wallet the OpenAPI wallet DTO
     * @return the domain model wallet
     */
    @Mapping(source = "isOmnibus", target = "omnibus")
    @Mapping(target = "network", ignore = true)
    @Mapping(target = "visibilityGroupID", ignore = true)
    Wallet fromDTO(TgvalidatordWallet wallet);

}


