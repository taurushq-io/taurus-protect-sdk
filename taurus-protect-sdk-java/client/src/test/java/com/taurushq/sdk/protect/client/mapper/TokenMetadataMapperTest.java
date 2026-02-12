package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.TokenMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordERCTokenMetadata;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordFATokenMetadata;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

class TokenMetadataMapperTest {

    @Test
    void fromERCDTO_mapsAllFields() {
        TgvalidatordERCTokenMetadata dto = new TgvalidatordERCTokenMetadata();
        dto.setName("USD Coin");
        dto.setDescription("A stablecoin pegged to the US dollar");
        dto.setDecimals("6");
        dto.setDataType("image/png");
        dto.setBase64Data("iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==");
        dto.setUri("https://example.com/token/usdc.json");

        TokenMetadata metadata = TokenMetadataMapper.INSTANCE.fromERCDTO(dto);

        assertEquals("USD Coin", metadata.getName());
        assertEquals("A stablecoin pegged to the US dollar", metadata.getDescription());
        assertEquals("6", metadata.getDecimals());
        assertEquals("image/png", metadata.getDataType());
        assertNotNull(metadata.getBase64Data());
        assertEquals("https://example.com/token/usdc.json", metadata.getUri());
    }

    @Test
    void fromERCDTO_mapsERC20Token() {
        TgvalidatordERCTokenMetadata dto = new TgvalidatordERCTokenMetadata();
        dto.setName("Wrapped Ether");
        dto.setDecimals("18");

        TokenMetadata metadata = TokenMetadataMapper.INSTANCE.fromERCDTO(dto);

        assertEquals("Wrapped Ether", metadata.getName());
        assertEquals("18", metadata.getDecimals());
        assertNull(metadata.getDescription());
        assertNull(metadata.getDataType());
        assertNull(metadata.getBase64Data());
        assertNull(metadata.getUri());
    }

    @Test
    void fromERCDTO_mapsNFTToken() {
        TgvalidatordERCTokenMetadata dto = new TgvalidatordERCTokenMetadata();
        dto.setName("Bored Ape #1234");
        dto.setDescription("A unique NFT from the BAYC collection");
        dto.setDecimals("0");
        dto.setDataType("image/jpeg");
        dto.setBase64Data("base64encodedimage");
        dto.setUri("ipfs://QmXXX/1234.json");

        TokenMetadata metadata = TokenMetadataMapper.INSTANCE.fromERCDTO(dto);

        assertEquals("Bored Ape #1234", metadata.getName());
        assertEquals("A unique NFT from the BAYC collection", metadata.getDescription());
        assertEquals("0", metadata.getDecimals());
        assertEquals("image/jpeg", metadata.getDataType());
        assertEquals("base64encodedimage", metadata.getBase64Data());
        assertEquals("ipfs://QmXXX/1234.json", metadata.getUri());
    }

    @Test
    void fromERCDTO_handlesNullFields() {
        TgvalidatordERCTokenMetadata dto = new TgvalidatordERCTokenMetadata();

        TokenMetadata metadata = TokenMetadataMapper.INSTANCE.fromERCDTO(dto);

        assertNull(metadata.getName());
        assertNull(metadata.getDescription());
        assertNull(metadata.getDecimals());
        assertNull(metadata.getDataType());
        assertNull(metadata.getBase64Data());
        assertNull(metadata.getUri());
    }

    @Test
    void fromERCDTO_handlesNullDto() {
        TokenMetadata metadata = TokenMetadataMapper.INSTANCE.fromERCDTO(null);
        assertNull(metadata);
    }

    @Test
    void fromFADTO_mapsAllFields() {
        TgvalidatordFATokenMetadata dto = new TgvalidatordFATokenMetadata();
        dto.setName("Tezos BTC");
        dto.setSymbol("tzBTC");
        dto.setDecimals("8");
        dto.setDataType("image/png");
        dto.setBase64Data("base64data");
        dto.setUri("https://example.com/metadata.json");

        TokenMetadata metadata = TokenMetadataMapper.INSTANCE.fromFADTO(dto);

        assertEquals("Tezos BTC", metadata.getName());
        assertEquals("tzBTC", metadata.getSymbol());
        assertEquals("8", metadata.getDecimals());
        assertEquals("image/png", metadata.getDataType());
        assertEquals("base64data", metadata.getBase64Data());
        assertEquals("https://example.com/metadata.json", metadata.getUri());
    }

    @Test
    void fromFADTO_handlesNullFields() {
        TgvalidatordFATokenMetadata dto = new TgvalidatordFATokenMetadata();

        TokenMetadata metadata = TokenMetadataMapper.INSTANCE.fromFADTO(dto);

        assertNull(metadata.getName());
        assertNull(metadata.getSymbol());
        assertNull(metadata.getDecimals());
        assertNull(metadata.getDataType());
        assertNull(metadata.getBase64Data());
        assertNull(metadata.getUri());
    }

    @Test
    void fromFADTO_handlesNullDto() {
        TokenMetadata metadata = TokenMetadataMapper.INSTANCE.fromFADTO(null);
        assertNull(metadata);
    }

    @Test
    void fromFADTO_mapsFA2Token() {
        TgvalidatordFATokenMetadata dto = new TgvalidatordFATokenMetadata();
        dto.setName("Objkt NFT");
        dto.setSymbol("OBJKT");
        dto.setDecimals("0");

        TokenMetadata metadata = TokenMetadataMapper.INSTANCE.fromFADTO(dto);

        assertEquals("Objkt NFT", metadata.getName());
        assertEquals("OBJKT", metadata.getSymbol());
        assertEquals("0", metadata.getDecimals());
    }
}
