package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ApiResponseCursor;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingAgreement;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingAgreementResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingOffer;
import com.taurushq.sdk.protect.client.model.taurusnetwork.LendingOfferResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.Participant;
import com.taurushq.sdk.protect.client.model.taurusnetwork.Pledge;
import com.taurushq.sdk.protect.client.model.taurusnetwork.PledgeResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.PledgeWithdrawal;
import com.taurushq.sdk.protect.client.model.taurusnetwork.PledgeWithdrawalResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.Settlement;
import com.taurushq.sdk.protect.client.model.taurusnetwork.SettlementResult;
import com.taurushq.sdk.protect.client.model.taurusnetwork.SharedAddress;
import com.taurushq.sdk.protect.client.model.taurusnetwork.SharedAddressResult;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetLendingAgreementsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetLendingOffersReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPledgesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetPledgesWithdrawalsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSettlementsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetSharedAddressesReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordLendingAgreement;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnLendingOffer;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnParticipant;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnPledge;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnPledgeWithdrawal;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnSettlement;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordTnSharedAddress;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting Taurus Network DTOs to domain models.
 */
@Mapper
public interface TaurusNetworkMapper {

    /**
     * Singleton instance of the mapper.
     */
    TaurusNetworkMapper INSTANCE = Mappers.getMapper(TaurusNetworkMapper.class);

    /**
     * Converts a participant DTO to a domain model.
     */
    Participant fromParticipantDTO(TgvalidatordTnParticipant dto);

    /**
     * Converts a list of participant DTOs to domain models.
     */
    List<Participant> fromParticipantDTOList(List<TgvalidatordTnParticipant> dtos);

    /**
     * Converts a pledge DTO to a domain model.
     */
    Pledge fromPledgeDTO(TgvalidatordTnPledge dto);

    /**
     * Converts a list of pledge DTOs to domain models.
     */
    List<Pledge> fromPledgeDTOList(List<TgvalidatordTnPledge> dtos);

    /**
     * Converts a get pledges reply to a result.
     */
    @Mapping(target = "pledges", source = "pledges")
    @Mapping(target = "cursor", source = "cursor")
    PledgeResult fromPledgesReply(TgvalidatordGetPledgesReply reply);

    /**
     * Converts a pledge withdrawal DTO to a domain model.
     */
    PledgeWithdrawal fromPledgeWithdrawalDTO(TgvalidatordTnPledgeWithdrawal dto);

    /**
     * Converts a list of pledge withdrawal DTOs to domain models.
     */
    List<PledgeWithdrawal> fromPledgeWithdrawalDTOList(List<TgvalidatordTnPledgeWithdrawal> dtos);

    /**
     * Converts a get pledge withdrawals reply to a result.
     */
    @Mapping(target = "withdrawals", source = "withdrawals")
    @Mapping(target = "cursor", source = "cursor")
    PledgeWithdrawalResult fromPledgeWithdrawalsReply(TgvalidatordGetPledgesWithdrawalsReply reply);

    /**
     * Converts a shared address DTO to a domain model.
     */
    SharedAddress fromSharedAddressDTO(TgvalidatordTnSharedAddress dto);

    /**
     * Converts a list of shared address DTOs to domain models.
     */
    List<SharedAddress> fromSharedAddressDTOList(List<TgvalidatordTnSharedAddress> dtos);

    /**
     * Converts a get shared addresses reply to a result.
     */
    @Mapping(target = "sharedAddresses", source = "sharedAddresses")
    @Mapping(target = "cursor", source = "cursor")
    SharedAddressResult fromSharedAddressesReply(TgvalidatordGetSharedAddressesReply reply);

    /**
     * Converts a settlement DTO to a domain model.
     */
    Settlement fromSettlementDTO(TgvalidatordTnSettlement dto);

    /**
     * Converts a list of settlement DTOs to domain models.
     */
    List<Settlement> fromSettlementDTOList(List<TgvalidatordTnSettlement> dtos);

    /**
     * Converts a get settlements reply to a result.
     */
    @Mapping(target = "settlements", source = "result")
    @Mapping(target = "cursor", source = "cursor")
    SettlementResult fromSettlementsReply(TgvalidatordGetSettlementsReply reply);

    /**
     * Converts a lending offer DTO to a domain model.
     */
    LendingOffer fromLendingOfferDTO(TgvalidatordTnLendingOffer dto);

    /**
     * Converts a list of lending offer DTOs to domain models.
     */
    List<LendingOffer> fromLendingOfferDTOList(List<TgvalidatordTnLendingOffer> dtos);

    /**
     * Converts a get lending offers reply to a result.
     */
    @Mapping(target = "offers", source = "lendingOffers")
    @Mapping(target = "cursor", source = "cursor")
    LendingOfferResult fromLendingOffersReply(TgvalidatordGetLendingOffersReply reply);

    /**
     * Converts a lending agreement DTO to a domain model.
     */
    LendingAgreement fromLendingAgreementDTO(TgvalidatordLendingAgreement dto);

    /**
     * Converts a list of lending agreement DTOs to domain models.
     */
    List<LendingAgreement> fromLendingAgreementDTOList(List<TgvalidatordLendingAgreement> dtos);

    /**
     * Converts a get lending agreements reply to a result.
     */
    @Mapping(target = "agreements", source = "lendingAgreements")
    @Mapping(target = "cursor", source = "cursor")
    LendingAgreementResult fromLendingAgreementsReply(TgvalidatordGetLendingAgreementsReply reply);

    /**
     * Converts a response cursor DTO to a domain model.
     */
    ApiResponseCursor fromCursor(TgvalidatordResponseCursor cursor);
}
