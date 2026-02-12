package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Address;
import com.taurushq.sdk.protect.client.model.Asset;
import com.taurushq.sdk.protect.client.model.Wallet;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddress;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAsset;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordWalletInfo;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting asset-related DTOs to domain models.
 * <p>
 * This mapper transforms OpenAPI-generated asset objects into the SDK's
 * domain model objects. It also provides list conversion methods for
 * addresses and wallets used by the AssetService.
 *
 * @see Asset
 * @see Address
 * @see Wallet
 */
@Mapper(uses = {CurrencyMapper.class, AddressMapper.class, WalletMapper.class})
public interface AssetMapper {

    /**
     * Singleton instance of the mapper.
     */
    AssetMapper INSTANCE = Mappers.getMapper(AssetMapper.class);

    /**
     * Converts an asset DTO to a domain model.
     *
     * @param asset the OpenAPI asset DTO
     * @return the domain model asset
     */
    Asset fromDTO(TgvalidatordAsset asset);

    /**
     * Converts a list of address DTOs to domain models.
     *
     * @param addresses the list of OpenAPI address DTOs
     * @return the list of domain model addresses
     */
    List<Address> fromAddressDTOList(List<TgvalidatordAddress> addresses);

    /**
     * Converts a list of wallet info DTOs to domain models.
     *
     * @param wallets the list of OpenAPI wallet info DTOs
     * @return the list of domain model wallets
     */
    List<Wallet> fromWalletInfoDTOList(List<TgvalidatordWalletInfo> wallets);
}