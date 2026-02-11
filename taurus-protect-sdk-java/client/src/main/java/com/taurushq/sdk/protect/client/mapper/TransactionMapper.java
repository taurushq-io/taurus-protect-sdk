package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.Transaction;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTransaction;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.List;


/**
 * MapStruct mapper for converting transaction DTOs to domain models.
 * <p>
 * This mapper transforms OpenAPI-generated transaction objects into the SDK's
 * domain model objects. It uses {@link AddressMapper} and {@link CurrencyMapper}
 * for nested object conversion.
 *
 * @see Transaction
 * @see TgvalidatordTransaction
 */
@Mapper(uses = {AddressMapper.class, CurrencyMapper.class})
public interface TransactionMapper {

    /**
     * Singleton instance of the mapper.
     */
    TransactionMapper INSTANCE = Mappers.getMapper(TransactionMapper.class);

    /**
     * Converts a single transaction DTO to a domain model.
     *
     * @param transaction the OpenAPI transaction DTO
     * @return the domain model transaction
     */
    Transaction fromDTO(TgvalidatordTransaction transaction);

    /**
     * Converts a list of transaction DTOs to domain models.
     *
     * @param transactions the list of OpenAPI transaction DTOs
     * @return the list of domain model transactions
     */
    List<Transaction> fromDTO(List<TgvalidatordTransaction> transactions);

}


