package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.BalanceHistoryPoint;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBalanceHistoryPoint;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * The interface Balance history point mapper.
 */
@Mapper(uses = {BalanceMapper.class})
public interface BalanceHistoryPointMapper {
    /**
     * The constant INSTANCE.
     */
    BalanceHistoryPointMapper INSTANCE = Mappers.getMapper(BalanceHistoryPointMapper.class);

    /**
     * From dto balance history point.
     *
     * @param balanceHistoryPoint the balance history point
     * @return the balance history point
     */
    BalanceHistoryPoint fromDTO(TgvalidatordBalanceHistoryPoint balanceHistoryPoint);

    /**
     * From dto list.
     *
     * @param balanceHistoryPoints the balance history points
     * @return the list
     */
    List<BalanceHistoryPoint> fromDTO(List<TgvalidatordBalanceHistoryPoint> balanceHistoryPoints);

}