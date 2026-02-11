package com.taurushq.sdk.protect.client.mapper;

import com.taurushq.sdk.protect.client.model.BlockchainInfo;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordBlockchainEntity;
import com.taurushq.sdk.protect.openapi.model.TgvalidatordCurrency;
import org.junit.jupiter.api.Test;

import java.util.Arrays;
import java.util.Collections;
import java.util.List;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertNotNull;
import static org.junit.jupiter.api.Assertions.assertNull;
import static org.junit.jupiter.api.Assertions.assertTrue;

class BlockchainMapperTest {

    @Test
    void fromDTO_mapsAllFields() {
        TgvalidatordCurrency currency = new TgvalidatordCurrency();
        currency.setId("ETH");
        currency.setName("Ethereum");

        TgvalidatordBlockchainEntity dto = new TgvalidatordBlockchainEntity();
        dto.setSymbol("ETH");
        dto.setName("Ethereum");
        dto.setNetwork("mainnet");
        dto.setChainId("1");
        dto.setConfirmations("12");
        dto.setBlockHeight("18500000");
        dto.setBlackholeAddress("0x0000000000000000000000000000000000000000");
        dto.setIsLayer2Chain(false);
        dto.setLayer1Network(null);
        dto.setBaseCurrency(currency);

        BlockchainInfo blockchain = BlockchainMapper.INSTANCE.fromDTO(dto);

        assertEquals("ETH", blockchain.getSymbol());
        assertEquals("Ethereum", blockchain.getName());
        assertEquals("mainnet", blockchain.getNetwork());
        assertEquals("1", blockchain.getChainId());
        assertEquals("12", blockchain.getConfirmations());
        assertEquals("18500000", blockchain.getBlockHeight());
        assertEquals("0x0000000000000000000000000000000000000000", blockchain.getBlackholeAddress());
        assertEquals(false, blockchain.getIsLayer2Chain());
        assertNull(blockchain.getLayer1Network());
        assertNotNull(blockchain.getBaseCurrency());
        assertEquals("ETH", blockchain.getBaseCurrency().getId());
    }

    @Test
    void fromDTO_mapsLayer2Chain() {
        TgvalidatordBlockchainEntity dto = new TgvalidatordBlockchainEntity();
        dto.setSymbol("MATIC");
        dto.setName("Polygon");
        dto.setNetwork("mainnet");
        dto.setChainId("137");
        dto.setIsLayer2Chain(true);
        dto.setLayer1Network("ETH");

        BlockchainInfo blockchain = BlockchainMapper.INSTANCE.fromDTO(dto);

        assertEquals("MATIC", blockchain.getSymbol());
        assertEquals("Polygon", blockchain.getName());
        assertEquals("137", blockchain.getChainId());
        assertEquals(true, blockchain.getIsLayer2Chain());
        assertEquals("ETH", blockchain.getLayer1Network());
    }

    @Test
    void fromDTO_handlesNullFields() {
        TgvalidatordBlockchainEntity dto = new TgvalidatordBlockchainEntity();
        dto.setSymbol("BTC");

        BlockchainInfo blockchain = BlockchainMapper.INSTANCE.fromDTO(dto);

        assertEquals("BTC", blockchain.getSymbol());
        assertNull(blockchain.getName());
        assertNull(blockchain.getNetwork());
        assertNull(blockchain.getChainId());
        assertNull(blockchain.getConfirmations());
        assertNull(blockchain.getBlockHeight());
        assertNull(blockchain.getBlackholeAddress());
        assertNull(blockchain.getIsLayer2Chain());
        assertNull(blockchain.getLayer1Network());
        assertNull(blockchain.getBaseCurrency());
    }

    @Test
    void fromDTO_handlesNullDto() {
        BlockchainInfo blockchain = BlockchainMapper.INSTANCE.fromDTO(null);
        assertNull(blockchain);
    }

    @Test
    void fromDTOList_mapsList() {
        TgvalidatordBlockchainEntity eth = new TgvalidatordBlockchainEntity();
        eth.setSymbol("ETH");
        eth.setName("Ethereum");
        eth.setNetwork("mainnet");

        TgvalidatordBlockchainEntity btc = new TgvalidatordBlockchainEntity();
        btc.setSymbol("BTC");
        btc.setName("Bitcoin");
        btc.setNetwork("mainnet");

        List<BlockchainInfo> blockchains = BlockchainMapper.INSTANCE.fromDTOList(Arrays.asList(eth, btc));

        assertNotNull(blockchains);
        assertEquals(2, blockchains.size());

        assertEquals("ETH", blockchains.get(0).getSymbol());
        assertEquals("Ethereum", blockchains.get(0).getName());

        assertEquals("BTC", blockchains.get(1).getSymbol());
        assertEquals("Bitcoin", blockchains.get(1).getName());
    }

    @Test
    void fromDTOList_handlesEmptyList() {
        List<BlockchainInfo> blockchains = BlockchainMapper.INSTANCE.fromDTOList(Collections.emptyList());
        assertNotNull(blockchains);
        assertTrue(blockchains.isEmpty());
    }

    @Test
    void fromDTOList_handlesNullList() {
        List<BlockchainInfo> blockchains = BlockchainMapper.INSTANCE.fromDTOList(null);
        assertNull(blockchains);
    }

    @Test
    void fromDTOList_mapsMultipleNetworks() {
        TgvalidatordBlockchainEntity ethMainnet = new TgvalidatordBlockchainEntity();
        ethMainnet.setSymbol("ETH");
        ethMainnet.setNetwork("mainnet");
        ethMainnet.setChainId("1");

        TgvalidatordBlockchainEntity ethGoerli = new TgvalidatordBlockchainEntity();
        ethGoerli.setSymbol("ETH");
        ethGoerli.setNetwork("goerli");
        ethGoerli.setChainId("5");

        TgvalidatordBlockchainEntity ethSepolia = new TgvalidatordBlockchainEntity();
        ethSepolia.setSymbol("ETH");
        ethSepolia.setNetwork("sepolia");
        ethSepolia.setChainId("11155111");

        List<BlockchainInfo> blockchains = BlockchainMapper.INSTANCE.fromDTOList(
                Arrays.asList(ethMainnet, ethGoerli, ethSepolia));

        assertEquals(3, blockchains.size());
        assertEquals("mainnet", blockchains.get(0).getNetwork());
        assertEquals("1", blockchains.get(0).getChainId());
        assertEquals("goerli", blockchains.get(1).getNetwork());
        assertEquals("5", blockchains.get(1).getChainId());
        assertEquals("sepolia", blockchains.get(2).getNetwork());
        assertEquals("11155111", blockchains.get(2).getChainId());
    }
}
