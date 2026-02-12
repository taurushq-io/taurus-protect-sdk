package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.AssetBalance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAssetBalance;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * The interface Asset balance mapper.
 */
@Mapper(uses = {AssetMapper.class, BalanceMapper.class})
public interface AssetBalanceMapper {
    /**
     * The constant INSTANCE.
     */
    AssetBalanceMapper INSTANCE = Mappers.getMapper(AssetBalanceMapper.class);

    /**
     * From dto asset balance.
     *
     * @param assetBalance the asset balance
     * @return the asset balance
     */
    AssetBalance fromDTO(TgvalidatordAssetBalance assetBalance);

    /**
     * From dto list.
     *
     * @param assetBalances the asset balances
     * @return the list
     */
    List<AssetBalance> fromDTO(List<TgvalidatordAssetBalance> assetBalances);

}