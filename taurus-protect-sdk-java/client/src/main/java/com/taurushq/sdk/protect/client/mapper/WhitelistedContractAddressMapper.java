package com.taurushq.sdk.protect.client.mapper;

import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.model.WhitelistedContractAddress;
import com.taurushq.sdk.protect.proto.v1.Whitelist;
import org.mapstruct.Mapper;
import org.mapstruct.Mapping;
import org.mapstruct.factory.Mappers;

import java.util.Base64;

/**
 * Mapper for WhitelistedContractAddress.
 */
@Mapper
public interface WhitelistedContractAddressMapper {

    /**
     * The constant INSTANCE.
     */
    WhitelistedContractAddressMapper INSTANCE = Mappers.getMapper(WhitelistedContractAddressMapper.class);

    /**
     * Maps protobuf WhitelistedContractAddress to model WhitelistedContractAddress.
     *
     * @param proto the protobuf object
     * @return the model object
     */
    @Mapping(target = "decimals", expression = "java((int) proto.getDecimals())")
    WhitelistedContractAddress fromDTO(Whitelist.WhitelistedContractAddress proto);

    /**
     * Decodes a WhitelistedContractAddress protobuf from raw bytes.
     *
     * @param data the raw protobuf bytes
     * @return the decoded WhitelistedContractAddress
     * @throws InvalidProtocolBufferException if the data cannot be parsed
     */
    default WhitelistedContractAddress fromBytes(byte[] data) throws InvalidProtocolBufferException {
        return fromDTO(Whitelist.WhitelistedContractAddress.parseFrom(data));
    }

    /**
     * Decodes a WhitelistedContractAddress protobuf from a base64-encoded string.
     *
     * @param base64 the base64-encoded protobuf data
     * @return the decoded WhitelistedContractAddress
     * @throws InvalidProtocolBufferException if the data cannot be parsed
     */
    default WhitelistedContractAddress fromBase64String(String base64) throws InvalidProtocolBufferException {
        byte[] data = Base64.getDecoder().decode(base64);
        return fromBytes(data);
    }
}
