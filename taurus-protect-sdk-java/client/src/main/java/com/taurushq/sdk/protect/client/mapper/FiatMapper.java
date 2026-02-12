package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;
import com.taurushq.sdk.protect.client.model.FiatProvider;
import com.taurushq.sdk.protect.client.model.FiatProviderAccount;
import com.taurushq.sdk.protect.client.model.FiatProviderAccountResult;
import com.taurushq.sdk.protect.client.model.FiatProviderCounterpartyAccount;
import com.taurushq.sdk.protect.client.model.FiatProviderCounterpartyAccountResult;
import com.taurushq.sdk.protect.client.model.FiatProviderOperation;
import com.taurushq.sdk.protect.client.model.FiatProviderOperationResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFiatProvider;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFiatProviderAccount;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFiatProviderCounterpartyAccount;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFiatProviderOperation;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderCounterpartyAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFiatProviderOperationsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting fiat provider DTOs to domain models.
 *
 * @see FiatProvider
 * @see FiatProviderAccount
 */
@Mapper
public interface FiatMapper {

    /**
     * Singleton instance of the mapper.
     */
    FiatMapper INSTANCE = Mappers.getMapper(FiatMapper.class);

    /**
     * Converts a fiat provider DTO to a domain model.
     *
     * @param dto the OpenAPI fiat provider DTO
     * @return the domain model fiat provider
     */
    FiatProvider fromProviderDTO(TgvalidatordFiatProvider dto);

    /**
     * Converts a list of fiat provider DTOs to domain models.
     *
     * @param dtos the list of OpenAPI fiat provider DTOs
     * @return the list of domain model fiat providers
     */
    List<FiatProvider> fromProviderDTOList(List<TgvalidatordFiatProvider> dtos);

    /**
     * Converts a fiat provider account DTO to a domain model.
     *
     * @param dto the OpenAPI fiat provider account DTO
     * @return the domain model fiat provider account
     */
    FiatProviderAccount fromAccountDTO(TgvalidatordFiatProviderAccount dto);

    /**
     * Converts a list of fiat provider account DTOs to domain models.
     *
     * @param dtos the list of OpenAPI fiat provider account DTOs
     * @return the list of domain model fiat provider accounts
     */
    List<FiatProviderAccount> fromAccountDTOList(List<TgvalidatordFiatProviderAccount> dtos);

    /**
     * Converts a get fiat provider accounts reply to a result.
     *
     * @param reply the OpenAPI reply
     * @return the fiat provider account result with pagination
     */
    @Mapping(target = "accounts", source = "result")
    @Mapping(target = "cursor", source = "cursor")
    FiatProviderAccountResult fromAccountsReply(TgvalidatordGetFiatProviderAccountsReply reply);

    /**
     * Converts a fiat provider counterparty account DTO to a domain model.
     *
     * @param dto the OpenAPI fiat provider counterparty account DTO
     * @return the domain model fiat provider counterparty account
     */
    @Mapping(target = "counterpartyId", source = "counterpartyID")
    FiatProviderCounterpartyAccount fromCounterpartyAccountDTO(TgvalidatordFiatProviderCounterpartyAccount dto);

    /**
     * Converts a list of fiat provider counterparty account DTOs to domain models.
     *
     * @param dtos the list of OpenAPI fiat provider counterparty account DTOs
     * @return the list of domain model fiat provider counterparty accounts
     */
    List<FiatProviderCounterpartyAccount> fromCounterpartyAccountDTOList(
            List<TgvalidatordFiatProviderCounterpartyAccount> dtos);

    /**
     * Converts a get fiat provider counterparty accounts reply to a result.
     *
     * @param reply the OpenAPI reply
     * @return the fiat provider counterparty account result with pagination
     */
    @Mapping(target = "accounts", source = "result")
    @Mapping(target = "cursor", source = "cursor")
    FiatProviderCounterpartyAccountResult fromCounterpartyAccountsReply(
            TgvalidatordGetFiatProviderCounterpartyAccountsReply reply);

    /**
     * Converts a fiat provider operation DTO to a domain model.
     *
     * @param dto the OpenAPI fiat provider operation DTO
     * @return the domain model fiat provider operation
     */
    FiatProviderOperation fromOperationDTO(TgvalidatordFiatProviderOperation dto);

    /**
     * Converts a list of fiat provider operation DTOs to domain models.
     *
     * @param dtos the list of OpenAPI fiat provider operation DTOs
     * @return the list of domain model fiat provider operations
     */
    List<FiatProviderOperation> fromOperationDTOList(List<TgvalidatordFiatProviderOperation> dtos);

    /**
     * Converts a get fiat provider operations reply to a result.
     *
     * @param reply the OpenAPI reply
     * @return the fiat provider operation result with pagination
     */
    @Mapping(target = "operations", source = "result")
    @Mapping(target = "cursor", source = "cursor")
    FiatProviderOperationResult fromOperationsReply(TgvalidatordGetFiatProviderOperationsReply reply);

    /**
     * Converts a response cursor DTO to a domain model.
     *
     * @param cursor the OpenAPI response cursor
     * @return the domain model cursor
     */
    ApiResponseCursor fromCursor(TgvalidatordResponseCursor cursor);
}
