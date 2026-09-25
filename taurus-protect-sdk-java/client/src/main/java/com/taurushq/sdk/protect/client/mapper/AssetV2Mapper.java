package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.AddressTargetV2;
import com.taurushq.sdk.protect.client.model.AssetAddressV2;
import com.taurushq.sdk.protect.client.model.AssetAttributeV2;
import com.taurushq.sdk.protect.client.model.AssetOperationV2;
import com.taurushq.sdk.protect.client.model.AssetV2;
import com.taurushq.sdk.protect.client.model.CantonNativeTokenV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAddressTargetV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAssetAddressV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAssetAttributeV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAssetOperationV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordAssetResourceV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBurnOperationDetailsV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCantonInstrumentConfigurationV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCantonNativeTokenAssetV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCantonNativeTokenParamsV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCreateOperationDetailsV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordImportOperationDetailsV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordMintOperationDetailsV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordSetKYCOperationDetailsV2;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordUpdateOperationDetailsV2;
import org.mapstruct.Mapper;
import org.mapstruct.factory.Mappers;

import java.util.ArrayList;
import java.util.List;

/**
 * Maps the v2 asset service DTOs (assets, asset addresses, asset operations) to models.
 * <p>
 * Hand-written: the DTOs name ids {@code *ID}, carry enums the models expose as their
 * wire values, and nest single-valued details the models flatten, so MapStruct's
 * by-name mapping would leave most fields silently null.
 */
@Mapper
public interface AssetV2Mapper {

    /**
     * Singleton instance of the mapper.
     */
    AssetV2Mapper INSTANCE = Mappers.getMapper(AssetV2Mapper.class);

    /**
     * Maps an asset.
     *
     * @param dto the DTO, may be null
     * @return the model, or null
     */
    default AssetV2 fromAssetDTO(final TgvalidatordAssetResourceV2 dto) {
        if (dto == null) {
            return null;
        }
        AssetV2 asset = new AssetV2();
        asset.setId(dto.getId());
        asset.setTenantId(dto.getTenantID());
        asset.setVersion(dto.getVersion());
        asset.setCreatedAt(dto.getCreatedAt());
        asset.setUpdatedAt(dto.getUpdatedAt());
        asset.setLabel(dto.getLabel());
        asset.setAssetType(dto.getAssetType());
        asset.setStatus(dto.getStatus() == null ? null : dto.getStatus().getValue());
        asset.setBlockchain(dto.getBlockchain());
        asset.setNetwork(dto.getNetwork());
        asset.setCurrencyId(dto.getCurrencyID());
        asset.setName(dto.getName());
        asset.setSymbol(dto.getSymbol());
        asset.setDecimals(dto.getDecimals());
        asset.setContractAddress(dto.getContractAddress());
        List<AssetAttributeV2> attributes = new ArrayList<>();
        if (dto.getAttributes() != null) {
            for (TgvalidatordAssetAttributeV2 a : dto.getAttributes()) {
                AssetAttributeV2 attribute = new AssetAttributeV2();
                attribute.setKey(a.getKey());
                attribute.setValue(a.getValue());
                attributes.add(attribute);
            }
        }
        asset.setAttributes(attributes);
        if (dto.getBlockchainAsset() != null) {
            asset.setCantonNativeToken(
                    fromCantonNativeTokenDTO(dto.getBlockchainAsset().getCantonNativeTokenAsset()));
        }
        return asset;
    }

    /**
     * Maps a page of assets.
     *
     * @param dtos the DTOs, may be null
     * @return the models, never null
     */
    default List<AssetV2> fromAssetDTOList(final List<TgvalidatordAssetResourceV2> dtos) {
        List<AssetV2> out = new ArrayList<>();
        if (dtos != null) {
            for (TgvalidatordAssetResourceV2 dto : dtos) {
                out.add(fromAssetDTO(dto));
            }
        }
        return out;
    }

    /**
     * Maps the Canton native-token details of an asset.
     *
     * @param dto the DTO, may be null
     * @return the model, or null
     */
    default CantonNativeTokenV2 fromCantonNativeTokenDTO(final TgvalidatordCantonNativeTokenAssetV2 dto) {
        if (dto == null) {
            return null;
        }
        CantonNativeTokenV2 token = new CantonNativeTokenV2();
        token.setInstrumentId(dto.getInstrumentID());
        TgvalidatordCantonInstrumentConfigurationV2 configuration = dto.getConfiguration();
        if (configuration != null) {
            token.setCid(configuration.getCid());
            token.setRequireCredentials(configuration.getRequireCredentials());
            token.setPaused(configuration.getPaused());
            token.setOperator(configuration.getOperator());
        }
        return token;
    }

    /**
     * Maps an asset address.
     *
     * @param dto the DTO, may be null
     * @return the model, or null
     */
    default AssetAddressV2 fromAssetAddressDTO(final TgvalidatordAssetAddressV2 dto) {
        if (dto == null) {
            return null;
        }
        AssetAddressV2 address = new AssetAddressV2();
        address.setAddress(dto.getAddress());
        address.setKycStatus(dto.getKycStatus() == null ? null : dto.getKycStatus().getValue());
        address.setBalance(dto.getBalance());
        address.setAddressType(dto.getAddressType() == null ? null : dto.getAddressType().getValue());
        address.setAddressId(dto.getAddressID());
        address.setWhitelistedAddressId(dto.getWhitelistedAddressID());
        return address;
    }

    /**
     * Maps a page of asset addresses.
     *
     * @param dtos the DTOs, may be null
     * @return the models, never null
     */
    default List<AssetAddressV2> fromAssetAddressDTOList(final List<TgvalidatordAssetAddressV2> dtos) {
        List<AssetAddressV2> out = new ArrayList<>();
        if (dtos != null) {
            for (TgvalidatordAssetAddressV2 dto : dtos) {
                out.add(fromAssetAddressDTO(dto));
            }
        }
        return out;
    }

    /**
     * Maps an asset operation, flattening its type-specific details.
     *
     * @param dto the DTO, may be null
     * @return the model, or null
     */
    default AssetOperationV2 fromAssetOperationDTO(final TgvalidatordAssetOperationV2 dto) {
        if (dto == null) {
            return null;
        }
        AssetOperationV2 op = new AssetOperationV2();
        op.setId(dto.getId());
        op.setAssetId(dto.getAssetID());
        op.setType(dto.getType() == null ? null : dto.getType().getValue());
        op.setStatus(dto.getStatus() == null ? null : dto.getStatus().getValue());
        op.setCreatedAt(dto.getCreatedAt());
        op.setUpdatedAt(dto.getUpdatedAt());
        op.setInitiatedByAddressId(dto.getInitiatedByAddressID());
        op.setFailureReason(dto.getFailureReason() == null ? null : dto.getFailureReason().getValue());
        op.setBlockingReason(dto.getBlockingReason() == null ? null : dto.getBlockingReason().getValue());
        applyCreate(op, dto.getCreate());
        applyUpdate(op, dto.getUpdate());
        applyImport(op, dto.getImport());
        TgvalidatordMintOperationDetailsV2 mint = dto.getMint();
        if (mint != null) {
            op.setAmount(mint.getAmount());
            op.setDestination(fromAddressTargetDTO(mint.getDestination()));
            op.setNftMetadata(mint.getNftMetadata());
        }
        TgvalidatordBurnOperationDetailsV2 burn = dto.getBurn();
        if (burn != null) {
            op.setAmount(burn.getAmount());
            op.setDestination(fromAddressTargetDTO(burn.getDestination()));
            op.setNftTokenIds(burn.getNftTokenIDs());
        }
        if (dto.getPauseAccount() != null) {
            op.setTarget(fromAddressTargetDTO(dto.getPauseAccount().getTarget()));
        }
        if (dto.getUnpauseAccount() != null) {
            op.setTarget(fromAddressTargetDTO(dto.getUnpauseAccount().getTarget()));
        }
        TgvalidatordSetKYCOperationDetailsV2 setKyc = dto.getSetKyc();
        if (setKyc != null) {
            op.setTarget(fromAddressTargetDTO(setKyc.getTarget()));
            op.setKycStatus(setKyc.getStatus() == null ? null : setKyc.getStatus().getValue());
        }
        return op;
    }

    /**
     * Maps a page of asset operations.
     *
     * @param dtos the DTOs, may be null
     * @return the models, never null
     */
    default List<AssetOperationV2> fromAssetOperationDTOList(final List<TgvalidatordAssetOperationV2> dtos) {
        List<AssetOperationV2> out = new ArrayList<>();
        if (dtos != null) {
            for (TgvalidatordAssetOperationV2 dto : dtos) {
                out.add(fromAssetOperationDTO(dto));
            }
        }
        return out;
    }

    /**
     * Maps the address an operation acts on.
     *
     * @param dto the DTO, may be null
     * @return the model, or null
     */
    default AddressTargetV2 fromAddressTargetDTO(final TgvalidatordAddressTargetV2 dto) {
        if (dto == null) {
            return null;
        }
        AddressTargetV2 target = new AddressTargetV2();
        target.setAddressId(dto.getAddressID());
        target.setWhitelistedAddressId(dto.getWhitelistedAddressID());
        return target;
    }

    /**
     * Copies the create details onto an operation.
     *
     * @param op     the operation
     * @param create the details, may be null
     */
    default void applyCreate(final AssetOperationV2 op, final TgvalidatordCreateOperationDetailsV2 create) {
        if (create == null) {
            return;
        }
        op.setLabel(create.getLabel());
        op.setPrice(create.getPrice());
        op.setDecimals(create.getDecimals());
        op.setBlockchain(create.getBlockchain());
        op.setNetwork(create.getNetwork());
        op.setAssetType(create.getAssetType());
        if (create.getParams() != null) {
            TgvalidatordCantonNativeTokenParamsV2 canton = create.getParams().getCantonNativeTokenParams();
            if (canton != null) {
                op.setInstrumentId(canton.getInstrumentID());
                op.setName(canton.getName());
                op.setSymbol(canton.getSymbol());
            }
        }
    }

    /**
     * Copies the update details onto an operation.
     *
     * @param op     the operation
     * @param update the details, may be null
     */
    default void applyUpdate(final AssetOperationV2 op, final TgvalidatordUpdateOperationDetailsV2 update) {
        if (update == null) {
            return;
        }
        op.setLabel(update.getLabel());
        op.setPrice(update.getPrice());
        if (update.getParams() != null && update.getParams().getCantonUpdateNativeTokenParams() != null) {
            op.setRequireCredentials(
                    update.getParams().getCantonUpdateNativeTokenParams().getRequireCredentials());
        }
    }

    /**
     * Copies the import details onto an operation.
     *
     * @param op         the operation
     * @param importDto  the details, may be null
     */
    default void applyImport(final AssetOperationV2 op, final TgvalidatordImportOperationDetailsV2 importDto) {
        if (importDto == null) {
            return;
        }
        op.setBlockchain(importDto.getBlockchain());
        op.setNetwork(importDto.getNetwork());
        op.setLabel(importDto.getLabel());
        op.setPrice(importDto.getPrice());
        op.setDecimals(importDto.getDecimals());
        op.setAddress(importDto.getAddress());
    }
}
