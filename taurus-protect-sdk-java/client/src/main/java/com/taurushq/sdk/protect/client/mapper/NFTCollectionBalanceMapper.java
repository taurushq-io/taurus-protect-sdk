package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.NFTCollectionBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordNFTCollectionBalance;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * The interface NFT collection balance mapper.
 */
@Mapper(uses = {CurrencyMapper.class, BalanceMapper.class})
public interface NFTCollectionBalanceMapper {
    /**
     * The constant INSTANCE.
     */
    NFTCollectionBalanceMapper INSTANCE = Mappers.getMapper(NFTCollectionBalanceMapper.class);

    /**
     * From dto nft collection balance.
     *
     * @param nftCollectionBalance the nft collection balance
     * @return the nft collection balance
     */
    NFTCollectionBalance fromDTO(TgvalidatordNFTCollectionBalance nftCollectionBalance);

    /**
     * From dto list.
     *
     * @param nftCollectionBalances the nft collection balances
     * @return the list
     */
    List<NFTCollectionBalance> fromDTO(List<TgvalidatordNFTCollectionBalance> nftCollectionBalances);
}
