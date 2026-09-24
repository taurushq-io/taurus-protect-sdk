package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.EarnReward;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordEarnReward;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordEarnRewardToken;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMerklTokenReward;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.ArrayList;
import java.util.List;

/**
 * Maps earn rewards to models, flattening the Merkl token reward.
 */
@Mapper
public interface EarnMapper {

    /**
     * Singleton instance of the mapper.
     */
    EarnMapper INSTANCE = Mappers.getMapper(EarnMapper.class);

    /**
     * Maps a reward.
     *
     * @param dto the DTO, may be null
     * @return the model, or null
     */
    default EarnReward fromDTO(final TgvalidatordEarnReward dto) {
        if (dto == null) {
            return null;
        }
        EarnReward reward = new EarnReward();
        reward.setId(dto.getId());
        reward.setRecipientAddressId(dto.getRecipientAddressId());
        reward.setRecipientAddress(dto.getRecipientAddress());
        reward.setRewardType(dto.getRewardType() == null ? null : dto.getRewardType().getValue());
        TgvalidatordMerklTokenReward merkl = dto.getMerklTokenReward();
        if (merkl != null) {
            reward.setAmount(merkl.getAmount());
            reward.setClaimed(merkl.getClaimed());
            reward.setPending(merkl.getPending());
            TgvalidatordEarnRewardToken token = merkl.getToken();
            if (token != null) {
                reward.setTokenAddress(token.getAddress());
                reward.setTokenSymbol(token.getSymbol());
                reward.setTokenAssetId(token.getAssetId());
            }
        }
        return reward;
    }

    /**
     * Maps a page of rewards.
     *
     * @param dtos the DTOs, may be null
     * @return the models, never null
     */
    default List<EarnReward> fromDTOList(final List<TgvalidatordEarnReward> dtos) {
        List<EarnReward> out = new ArrayList<>();
        if (dtos != null) {
            for (TgvalidatordEarnReward dto : dtos) {
                out.add(fromDTO(dto));
            }
        }
        return out;
    }
}
