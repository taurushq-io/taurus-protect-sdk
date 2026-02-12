package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Currency;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

/**
 * The interface Currency mapper.
 */
@Mapper
public interface CurrencyMapper {
    /**
     * The constant INSTANCE.
     */
    CurrencyMapper INSTANCE = Mappers.getMapper(CurrencyMapper.class);

    /**
     * From dto currency.
     *
     * @param currency the currency
     * @return the currency
     */
    @Mapping(source = "tokenID", target = "tokenId")
    @Mapping(source = "isToken", target = "token")
    @Mapping(source = "isERC20", target = "ERC20")
    @Mapping(source = "isUTXOBased", target = "UTXOBased")
    @Mapping(source = "isAccountBased", target = "accountBased")
    @Mapping(source = "isFiat", target = "fiat")
    @Mapping(source = "isFA12", target = "FA12")
    @Mapping(source = "isFA20", target = "FA20")
    @Mapping(source = "isNFT", target = "NFT")
    Currency fromDTO(TgvalidatordCurrency currency);

}


