package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.WhitelistedContractAddress;
import com.taurushq.sdk.protect.proto.v1.Whitelist;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;

class WhitelistedContractAddressMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        Whitelist.WhitelistedContractAddress proto = Whitelist.WhitelistedContractAddress.newBuilder()
                .setBlockchain("ETH")
                .setNetwork("mainnet")
                .setName("USD Coin")
                .setSymbol("USDC")
                .setDecimals(6)
                .setContractAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
                .setTokenId("0")
                .setKind("erc20")
                .build();

        WhitelistedContractAddress result = WhitelistedContractAddressMapper.INSTANCE.fromDTO(proto);

        assertNotNull(result);
        assertEquals("ETH", result.getBlockchain());
        assertEquals("mainnet", result.getNetwork());
        assertEquals("USD Coin", result.getName());
        assertEquals("USDC", result.getSymbol());
        assertEquals(6, result.getDecimals());
        assertEquals("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48", result.getContractAddress());
        assertEquals("0", result.getTokenId());
        assertEquals("erc20", result.getKind());
    }

    @Test
    void fromDTO_handlesEmptyProto() {
        Whitelist.WhitelistedContractAddress proto =
                Whitelist.WhitelistedContractAddress.getDefaultInstance();

        WhitelistedContractAddress result = WhitelistedContractAddressMapper.INSTANCE.fromDTO(proto);

        assertNotNull(result);
        assertEquals("", result.getBlockchain());
        assertEquals("", result.getNetwork());
        assertEquals("", result.getName());
        assertEquals("", result.getSymbol());
        assertEquals(0, result.getDecimals());
    }

    @Test
    void fromDTO_handlesNullDto() {
        WhitelistedContractAddress result = WhitelistedContractAddressMapper.INSTANCE.fromDTO(null);
        assertNull(result);
    }

    @Test
    void fromBytes_decodesProtobuf() throws Exception {
        Whitelist.WhitelistedContractAddress proto = Whitelist.WhitelistedContractAddress.newBuilder()
                .setBlockchain("MATIC")
                .setNetwork("mainnet")
                .setName("Wrapped ETH")
                .setSymbol("WETH")
                .setDecimals(18)
                .setContractAddress("0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619")
                .build();

        byte[] bytes = proto.toByteArray();
        WhitelistedContractAddress result = WhitelistedContractAddressMapper.INSTANCE.fromBytes(bytes);

        assertNotNull(result);
        assertEquals("MATIC", result.getBlockchain());
        assertEquals("Wrapped ETH", result.getName());
        assertEquals("WETH", result.getSymbol());
        assertEquals(18, result.getDecimals());
    }

    @Test
    void fromBase64String_decodesBase64() throws Exception {
        Whitelist.WhitelistedContractAddress proto = Whitelist.WhitelistedContractAddress.newBuilder()
                .setBlockchain("XTZ")
                .setNetwork("mainnet")
                .setName("Kolibri USD")
                .setSymbol("kUSD")
                .setDecimals(18)
                .setKind("fa2")
                .build();

        String base64 = java.util.Base64.getEncoder().encodeToString(proto.toByteArray());
        WhitelistedContractAddress result = WhitelistedContractAddressMapper.INSTANCE.fromBase64String(base64);

        assertNotNull(result);
        assertEquals("XTZ", result.getBlockchain());
        assertEquals("Kolibri USD", result.getName());
        assertEquals("kUSD", result.getSymbol());
        assertEquals("fa2", result.getKind());
    }
}
