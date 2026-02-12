package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.ADAStakePoolInfo;
import com.taurushq.sdk.protect.client.model.ApiResponseCursor;
import com.taurushq.sdk.protect.client.model.ETHValidatorInfo;
import com.taurushq.sdk.protect.client.model.FTMValidatorInfo;
import com.taurushq.sdk.protect.client.model.ICPNeuronInfo;
import com.taurushq.sdk.protect.client.model.NEARValidatorInfo;
import com.taurushq.sdk.protect.client.model.SolanaStakeAccount;
import com.taurushq.sdk.protect.client.model.StakeAccount;
import com.taurushq.sdk.protect.client.model.StakeAccountResult;
import com.taurushq.sdk.protect.client.model.XTZStakingRewards;
import com.taurushq.sdk.protect.openapi.model.GetICPNeuronInfoReplyNeuronState;
import com.taurushq.sdk.protect.openapi.model.SolanaStakeAccountState;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordETHValidatorInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetADAStakePoolInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetFTMValidatorInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetICPNeuronInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetNEARValidatorInfoReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetStakeAccountsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordGetXTZAddressStakingRewardsReply;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordResponseCursor;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSolanaStakeAccount;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordStakeAccount;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordStakeAccountType;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.Named;
import org.mapstruct.factory.Mappers;

import java.util.List;

/**
 * MapStruct mapper for converting staking-related OpenAPI DTOs to client model objects.
 * <p>
 * This mapper handles the conversion of validator info, stake account, and staking
 * rewards data from OpenAPI generated models to the SDK's clean domain models.
 */
@Mapper
public interface StakingMapper {

    /**
     * Singleton instance of the mapper.
     */
    StakingMapper INSTANCE = Mappers.getMapper(StakingMapper.class);

    /**
     * Maps an ADA stake pool info reply to the domain model.
     *
     * @param reply the OpenAPI reply
     * @return the domain model
     */
    ADAStakePoolInfo fromADAStakePoolInfoReply(TgvalidatordGetADAStakePoolInfoReply reply);

    /**
     * Maps an ETH validator info DTO to the domain model.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    @Mapping(target = "addressId", source = "addressID")
    ETHValidatorInfo fromETHValidatorInfo(TgvalidatordETHValidatorInfo dto);

    /**
     * Maps a list of ETH validator info DTOs to domain models.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<ETHValidatorInfo> fromETHValidatorInfoList(List<TgvalidatordETHValidatorInfo> dtos);

    /**
     * Maps an FTM validator info reply to the domain model.
     *
     * @param reply the OpenAPI reply
     * @return the domain model
     */
    @Mapping(target = "validatorId", source = "validatorID")
    @Mapping(target = "active", source = "isActive")
    FTMValidatorInfo fromFTMValidatorInfoReply(TgvalidatordGetFTMValidatorInfoReply reply);

    /**
     * Maps an ICP neuron info reply to the domain model.
     *
     * @param reply the OpenAPI reply
     * @return the domain model
     */
    @Mapping(target = "neuronState", source = "neuronState", qualifiedByName = "toNeuronStateString")
    @Mapping(target = "stakeE8s", source = "stakeE8S")
    ICPNeuronInfo fromICPNeuronInfoReply(TgvalidatordGetICPNeuronInfoReply reply);

    /**
     * Converts a neuron state enum to a string.
     *
     * @param state the neuron state enum
     * @return the string representation
     */
    @Named("toNeuronStateString")
    default String toNeuronStateString(GetICPNeuronInfoReplyNeuronState state) {
        return state != null ? state.getValue() : null;
    }

    /**
     * Maps a NEAR validator info reply to the domain model.
     *
     * @param reply the OpenAPI reply
     * @return the domain model
     */
    @Mapping(target = "stakingPaused", source = "isStakingPaused")
    NEARValidatorInfo fromNEARValidatorInfoReply(TgvalidatordGetNEARValidatorInfoReply reply);

    /**
     * Maps a stake accounts reply to the domain model result.
     *
     * @param reply the OpenAPI reply
     * @return the domain model result
     */
    @Mapping(target = "stakeAccounts", source = "stakeAccounts")
    @Mapping(target = "cursor", source = "cursor")
    StakeAccountResult fromStakeAccountsReply(TgvalidatordGetStakeAccountsReply reply);

    /**
     * Maps a stake account DTO to the domain model.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    @Mapping(target = "accountType", source = "accountType", qualifiedByName = "toAccountTypeString")
    @Mapping(target = "solanaStakeAccount", source = "solanaStakeAccount")
    StakeAccount fromStakeAccount(TgvalidatordStakeAccount dto);

    /**
     * Maps a list of stake account DTOs to domain models.
     *
     * @param dtos the list of OpenAPI DTOs
     * @return the list of domain models
     */
    List<StakeAccount> fromStakeAccountList(List<TgvalidatordStakeAccount> dtos);

    /**
     * Converts an account type enum to a string.
     *
     * @param type the account type enum
     * @return the string representation
     */
    @Named("toAccountTypeString")
    default String toAccountTypeString(TgvalidatordStakeAccountType type) {
        return type != null ? type.getValue() : null;
    }

    /**
     * Maps a Solana stake account DTO to the domain model.
     *
     * @param dto the OpenAPI DTO
     * @return the domain model
     */
    @Mapping(target = "state", source = "state", qualifiedByName = "toSolanaStateString")
    SolanaStakeAccount fromSolanaStakeAccount(TgvalidatordSolanaStakeAccount dto);

    /**
     * Converts a Solana stake account state enum to a string.
     *
     * @param state the state enum
     * @return the string representation
     */
    @Named("toSolanaStateString")
    default String toSolanaStateString(SolanaStakeAccountState state) {
        return state != null ? state.getValue() : null;
    }

    /**
     * Maps an XTZ staking rewards reply to the domain model.
     *
     * @param reply the OpenAPI reply
     * @return the domain model
     */
    XTZStakingRewards fromXTZStakingRewardsReply(TgvalidatordGetXTZAddressStakingRewardsReply reply);

    /**
     * Maps a response cursor DTO to the domain model.
     *
     * @param cursor the OpenAPI cursor
     * @return the domain model cursor
     */
    ApiResponseCursor fromCursor(TgvalidatordResponseCursor cursor);
}
