package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.FeePayer;
import com.taurushq.sdk.protect.client.model.FeePayerEth;
import com.taurushq.sdk.protect.client.model.FeePayerEthLocal;
import com.taurushq.sdk.protect.client.model.FeePayerEthRemote;
import com.taurushq.sdk.protect.client.model.FeePayerInfo;
import com.taurushq.sdk.protect.openapi.model.ETHLocal;
import com.taurushq.sdk.protect.openapi.model.ETHRemote;
import com.taurushq.sdk.protect.openapi.model.FeePayerETH;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFeePayer;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFeePayerEnvelope;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting fee payer-related DTOs to domain models.
 */
@Mapper
public interface FeePayerMapper {

    FeePayerMapper INSTANCE = Mappers.getMapper(FeePayerMapper.class);

    /**
     * Converts a FeePayerEnvelope DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    @Mapping(source = "feePayer", target = "feePayerInfo")
    FeePayer fromDTO(TgvalidatordFeePayerEnvelope dto);

    /**
     * Converts a list of FeePayerEnvelope DTOs to domain models.
     *
     * @param dtos the DTOs to convert
     * @return the domain models
     */
    List<FeePayer> fromDTOList(List<TgvalidatordFeePayerEnvelope> dtos);

    /**
     * Converts a FeePayer DTO to a FeePayerInfo domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    FeePayerInfo fromFeePayerDTO(TgvalidatordFeePayer dto);

    /**
     * Converts a FeePayerETH DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    FeePayerEth fromEthDTO(FeePayerETH dto);

    /**
     * Converts an ETHLocal DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    FeePayerEthLocal fromEthLocalDTO(ETHLocal dto);

    /**
     * Converts an ETHRemote DTO to a domain model.
     *
     * @param dto the DTO to convert
     * @return the domain model
     */
    FeePayerEthRemote fromEthRemoteDTO(ETHRemote dto);
}
