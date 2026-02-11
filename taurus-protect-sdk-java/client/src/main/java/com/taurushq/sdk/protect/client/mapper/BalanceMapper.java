package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Balance;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBalance;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

/**
 * The interface Balance mapper.
 */
@Mapper
public interface BalanceMapper {
    /**
     * The constant INSTANCE.
     */
    BalanceMapper INSTANCE = Mappers.getMapper(BalanceMapper.class);

    /**
     * From dto balance.
     *
     * @param balance the balance
     * @return the balance
     */
    Balance fromDTO(TgvalidatordBalance balance);

}


