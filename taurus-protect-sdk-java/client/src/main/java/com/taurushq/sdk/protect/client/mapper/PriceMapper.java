package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ConversionResult;
import com.taurushq.sdk.protect.client.model.Price;
import com.taurushq.sdk.protect.client.model.PriceHistoryPoint;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordConversionValue;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrencyPrice;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordPricesHistoryPoint;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * The interface Price mapper.
 */
@Mapper(uses = {CurrencyMapper.class})
public interface PriceMapper {
    /**
     * The constant INSTANCE.
     */
    PriceMapper INSTANCE = Mappers.getMapper(PriceMapper.class);

    /**
     * From dto price.
     *
     * @param price the price
     * @return the price
     */
    // signatures is ignored, not mapped: the generated TgvalidatordCurrencyPrice has no
    // such field even though apis.swagger.json declares it, so there is nothing for
    // MapStruct to read. Regenerating the openapi module turns this into a direct
    // mapping; until then a caller verifying a price must supply the signatures it
    // received. PriceVerifier does the checking either way.
    @Mapping(target = "signatures", ignore = true)
    @Mapping(source = "creationDate", target = "createdAt")
    @Mapping(source = "updateDate", target = "updatedAt")
    Price fromDTO(TgvalidatordCurrencyPrice price);

    /**
     * From dto list.
     *
     * @param prices the prices
     * @return the list
     */
    List<Price> fromDTO(List<TgvalidatordCurrencyPrice> prices);

    /**
     * From dto price history point.
     *
     * @param point the point
     * @return the price history point
     */
    PriceHistoryPoint fromDTO(TgvalidatordPricesHistoryPoint point);

    /**
     * From dto history list.
     *
     * @param points the points
     * @return the list
     */
    List<PriceHistoryPoint> fromDTOHistory(List<TgvalidatordPricesHistoryPoint> points);

    /**
     * From dto conversion result.
     *
     * @param value the value
     * @return the conversion result
     */
    ConversionResult fromDTO(TgvalidatordConversionValue value);

    /**
     * From dto conversion list.
     *
     * @param values the values
     * @return the list
     */
    List<ConversionResult> fromDTOConversion(List<TgvalidatordConversionValue> values);
}
