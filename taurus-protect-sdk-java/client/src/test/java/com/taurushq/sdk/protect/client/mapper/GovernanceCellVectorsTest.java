package com.taurushq.sdk.protect.client.mapper;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.io.IOException;
import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.Base64;
import java.util.LinkedHashMap;
import java.util.Map;

import com.google.gson.JsonArray;
import com.google.gson.JsonElement;
import com.google.gson.JsonObject;
import com.google.gson.JsonParser;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleCell;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.Test;

/**
 * Cross-SDK cell wire-format alignment.
 *
 * <p>Every entry in the shared golden-vector file (produced by the Go SDK) must
 * encode to exactly the recorded bytes and decode back to the same typed cell.
 * The same file is consumed by the Go/Python/TypeScript suites, pinning
 * byte-for-byte parity across all four SDKs.
 */
class GovernanceCellVectorsTest {

    private static final String RAW_DESCRIPTION = "raw cell (unknown cell type preserved verbatim)";

    private static JsonArray vectors;

    /**
     * Typed cell cases keyed by the shared vectors' {@code description}, mirroring
     * the Go SDK's allCellCases table. The vectors file is the byte oracle; the
     * expected column type is derived from the cell itself via
     * {@link RuleCellCodec#cellFamily}.
     */
    private static Map<String, RuleCell> cases() {
        Map<String, RuleCell> m = new LinkedHashMap<>();
        m.put("fiat amount any", new RuleCell.FiatAmountAny());
        m.put("fiat amount is zero", new RuleCell.FiatAmountIsZero());
        m.put("fiat amount range", new RuleCell.FiatAmountRange("1000", "50000"));
        m.put("source any", new RuleCell.SourceAny());
        m.put("source internal wallet", new RuleCell.SourceInternalWallet("m/44'/60'/0'"));
        m.put("source internal address", new RuleCell.SourceInternalAddress("0xabc", "m/44'/60'/0'/0/0"));
        m.put("source any exchange", new RuleCell.SourceAnyExchange());
        m.put("source exchange", new RuleCell.SourceExchange("kraken-main"));
        m.put("source external address", new RuleCell.SourceExternalAddress("0xdef", "memo-1"));
        m.put("destination any", new RuleCell.DestinationAny());
        m.put("destination internal wallet", new RuleCell.DestinationInternalWallet("m/44'/60'/1'"));
        m.put("destination internal address", new RuleCell.DestinationInternalAddress("0x111", "m/44'/60'/1'/0/0"));
        m.put("destination external address", new RuleCell.DestinationExternalAddress("0x222", "dest-memo"));
        m.put("destination any exchange", new RuleCell.DestinationAnyExchange());
        m.put("destination exchange", new RuleCell.DestinationExchange("binance-desk", "x"));
        m.put("destination contract address", new RuleCell.DestinationContractAddress("0x333", "USDC", "USDC", "ETH"));
        m.put("destination contract address (unknown blockchain)", new RuleCell.DestinationContractAddress("0x1", "", "", "4242"));
        m.put("destination any external address", new RuleCell.DestinationAnyExternalAddress());
        m.put("destination any contract address", new RuleCell.DestinationAnyContractAddress());
        m.put("string equal any", new RuleCell.StringEqualAny());
        m.put("string equal empty", new RuleCell.StringEqualEmpty());
        m.put("string equal value", new RuleCell.StringEqualValue("contract-id-42"));
        m.put("bytes equal any", new RuleCell.BytesEqualAny());
        m.put("bytes equal empty", new RuleCell.BytesEqualEmpty());
        m.put("bytes equal value", new RuleCell.BytesEqualValue(new byte[] {(byte) 0xDE, (byte) 0xAD, (byte) 0xBE, (byte) 0xEF}));
        m.put("string array equal any", new RuleCell.StringArrayEqualAny());
        m.put("string array equal empty", new RuleCell.StringArrayEqualEmpty());
        m.put("string array equal value", new RuleCell.StringArrayEqualValue(Arrays.asList("a", "b", "c")));
        m.put("integer greater any", new RuleCell.IntegerGreaterAny());
        m.put("integer greater positive value", new RuleCell.IntegerGreaterValue(BigInteger.valueOf(50)));
        m.put("integer greater negative value", new RuleCell.IntegerGreaterValue(BigInteger.valueOf(-50)));
        m.put("integer greater zero", new RuleCell.IntegerGreaterValue(BigInteger.ZERO));
        m.put("uinteger greater any", new RuleCell.UIntegerGreaterAny());
        m.put("uinteger greater is zero", new RuleCell.UIntegerGreaterIsZero());
        m.put("uinteger greater value", new RuleCell.UIntegerGreaterValue(new BigInteger("18446744073709551615")));
        m.put("uinteger greater is equal", new RuleCell.UIntegerGreaterIsEqual(BigInteger.valueOf(10)));
        m.put("whitelisted contract any", new RuleCell.WhitelistedContractAny());
        m.put("whitelisted contract address", new RuleCell.WhitelistedContractAddress("0x444", "DAI", "DAI", "ETH"));
        return m;
    }

    @BeforeAll
    static void loadVectors() throws IOException {
        String[] candidates = {
                "../../scripts/resources/governance-cell-vectors.json", // from client/
                "../scripts/resources/governance-cell-vectors.json",    // from sdk-java/
                "scripts/resources/governance-cell-vectors.json",       // from repo root
        };
        for (String candidate : candidates) {
            Path p = Paths.get(candidate).toAbsolutePath().normalize();
            if (Files.exists(p)) {
                String json = new String(Files.readAllBytes(p), StandardCharsets.UTF_8);
                vectors = JsonParser.parseString(json).getAsJsonArray();
                return;
            }
        }
        throw new IOException("Cannot find governance-cell-vectors.json. Working directory: " + System.getProperty("user.dir"));
    }

    @Test
    void vectorsCoverAllCases() {
        Map<String, RuleCell> cases = cases();
        for (String desc : cases.keySet()) {
            boolean found = false;
            for (JsonElement e : vectors) {
                if (desc.equals(e.getAsJsonObject().get("description").getAsString())) {
                    found = true;
                    break;
                }
            }
            assertTrue(found, "vector missing for " + desc);
        }
        assertEquals(cases.size() + 1, vectors.size(), "typed cases + one raw case");
    }

    @Test
    void cellVectorsEncodeAndDecode() {
        Map<String, RuleCell> cases = cases();
        for (JsonElement e : vectors) {
            JsonObject v = e.getAsJsonObject();
            String desc = v.get("description").getAsString();
            String colType = v.get("column_type").getAsString();
            byte[] wire = Base64.getDecoder().decode(v.get("wire_base64").getAsString());

            if (RAW_DESCRIPTION.equals(desc)) {
                RuleCell decoded = RuleCellCodec.decode(colType, wire);
                assertTrue(decoded instanceof RuleCell.RawCell, desc + ": expected RawCell, got " + decoded);
                assertArrayEquals(wire, RuleCellCodec.encode(colType, decoded), desc);
                continue;
            }

            RuleCell cell = cases.get(desc);
            assertTrue(cell != null, "no case for vector " + desc);
            assertEquals(colType, RuleCellCodec.cellFamily(cell), desc);
            // encode(typed) == recorded wire bytes (cross-SDK byte parity)
            assertArrayEquals(wire, RuleCellCodec.encode(colType, cell), desc);
            // decode(wire) == typed
            assertEquals(cell, RuleCellCodec.decode(colType, wire), desc);
        }
    }
}
