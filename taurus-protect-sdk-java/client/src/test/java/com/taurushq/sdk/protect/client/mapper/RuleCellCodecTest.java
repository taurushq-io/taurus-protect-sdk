package com.taurushq.sdk.protect.client.mapper;

import static org.junit.jupiter.api.Assertions.assertArrayEquals;
import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.math.BigInteger;
import java.util.Collections;

import com.google.protobuf.ByteString;
import com.google.protobuf.UnknownFieldSet;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleCell;
import com.taurushq.sdk.protect.proto.v1.RequestReply;
import org.junit.jupiter.api.Test;

/**
 * Branch and edge coverage for {@link RuleCellCodec}: encode error paths, the
 * RawCell lossless fallbacks, the magnitude helpers, {@code cellFamily}
 * resolution, and blockchain numeric passthrough.
 *
 * <p>Complements the shared golden-vector oracle ({@link GovernanceCellVectorsTest}),
 * which pins the happy-path wire bytes cross-SDK, and the container round-trip test.
 */
class RuleCellCodecTest {

    // --- Taxonomy 1: encode error branches ---

    @Test
    void encodeIntoWrongColumnFamilyThrows() {
        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class,
                () -> RuleCellCodec.encode("RuleFiatAmount", new RuleCell.StringEqualValue("x")));
        assertTrue(ex.getMessage().contains("is not valid for column type"));
    }

    @Test
    void encodeNegativeUnsignedIntegerThrows() {
        RuleCell[] negatives = {
                new RuleCell.UIntegerGreaterValue(BigInteger.valueOf(-1)),
                new RuleCell.UIntegerGreaterIsEqual(BigInteger.valueOf(-5)),
        };
        for (RuleCell cell : negatives) {
            IllegalArgumentException ex = assertThrows(IllegalArgumentException.class,
                    () -> RuleCellCodec.encode("RuleUIntegerGreater", cell));
            assertTrue(ex.getMessage().contains("requires a non-negative value"), cell.kind());
        }
    }

    @Test
    void encodeUnsupportedKindHitsTerminal() {
        // A recognised family prefix (FiatAmount) with no matching switch arm reaches
        // the encoder's terminal default; the empty columnType skips the family check.
        RuleCell bogus = new RuleCell() {
            @Override
            public String kind() {
                return "FiatAmountBogus";
            }
        };
        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class,
                () -> RuleCellCodec.encode("", bogus));
        assertTrue(ex.getMessage().contains("unsupported rule cell kind"));
    }

    @Test
    void cellFamilyUnknownKindThrows() {
        RuleCell zzz = new RuleCell() {
            @Override
            public String kind() {
                return "Zzz";
            }
        };
        IllegalArgumentException ex = assertThrows(IllegalArgumentException.class,
                () -> RuleCellCodec.cellFamily(zzz));
        assertTrue(ex.getMessage().contains("unknown rule cell kind"));
    }

    // --- Taxonomy 2: RawCell fallback branches ---

    @Test
    void decodeUnknownColumnTypeReturnsRawVerbatim() {
        byte[] payload = {0x0a, 0x03, 'a', 'b', 'c'};
        RuleCell cell = RuleCellCodec.decode("SomeFutureColumn", payload);
        assertTrue(cell instanceof RuleCell.RawCell);
        assertEquals("SomeFutureColumn", ((RuleCell.RawCell) cell).columnType);
        assertArrayEquals(payload, RuleCellCodec.encode("SomeFutureColumn", cell));
    }

    @Test
    void decodeUnknownColumnTypeEmptyPayloadReturnsRaw() {
        // Java wraps even empty bytes of an unknown column in a RawCell (it does not
        // return null), so the cell survives a decode/encode cycle unchanged.
        RuleCell cell = RuleCellCodec.decode("SomeFutureColumn", new byte[0]);
        assertTrue(cell instanceof RuleCell.RawCell);
        assertArrayEquals(new byte[0], RuleCellCodec.encode("SomeFutureColumn", cell));
    }

    @Test
    void decodeKnownCellWithUnknownSubfieldFallsBackToRaw() {
        // A RuleFiatAmount carrying a protobuf field this SDK does not model cannot
        // round-trip through the typed FiatAmountRange, so the lossless guard keeps it
        // as a RawCell that re-encodes byte-for-byte.
        byte[] range = RequestReply.RuleFiatAmountRange.newBuilder()
                .setMinAmount("1000").setMaxAmount("50000").build().toByteArray();
        UnknownFieldSet uf = UnknownFieldSet.newBuilder()
                .addField(500, UnknownFieldSet.Field.newBuilder().addVarint(42).build())
                .build();
        byte[] withUnknown = RequestReply.RuleFiatAmount.newBuilder()
                .setType(RequestReply.RuleFiatAmount.RuleFiatAmountType.RuleFiatAmountRange)
                .setPayload(ByteString.copyFrom(range))
                .setUnknownFields(uf)
                .build().toByteArray();

        RuleCell cell = RuleCellCodec.decode("RuleFiatAmount", withUnknown);
        assertTrue(cell instanceof RuleCell.RawCell);
        assertArrayEquals(withUnknown, RuleCellCodec.encode("RuleFiatAmount", cell));
    }

    @Test
    void decodePayloadFreeCellCarryingPayloadFallsBackToRaw() {
        // A payload-free cell type (FiatAmountAny) that nonetheless carries a payload
        // does not round-trip through the typed value and is preserved as RawCell.
        byte[] anyWithPayload = RequestReply.RuleFiatAmount.newBuilder()
                .setType(RequestReply.RuleFiatAmount.RuleFiatAmountType.RuleFiatAmountAny)
                .setPayload(ByteString.copyFromUtf8("unexpected"))
                .build().toByteArray();

        RuleCell cell = RuleCellCodec.decode("RuleFiatAmount", anyWithPayload);
        assertTrue(cell instanceof RuleCell.RawCell);
        assertArrayEquals(anyWithPayload, RuleCellCodec.encode("RuleFiatAmount", cell));
    }

    // --- Taxonomy 3: magnitude helpers (package-private statics) ---

    @Test
    void magnitudeAndFromMagnitudeRoundTrip() {
        assertArrayEquals(new byte[] {0}, RuleCellCodec.magnitude(BigInteger.ZERO));
        // 128 == 0x80: BigInteger.toByteArray prepends a 0x00 sign byte; magnitude drops it.
        assertArrayEquals(new byte[] {(byte) 0x80}, RuleCellCodec.magnitude(BigInteger.valueOf(128)));
        // 256 == 0x0100: no spurious leading zero to drop.
        assertArrayEquals(new byte[] {0x01, 0x00}, RuleCellCodec.magnitude(BigInteger.valueOf(256)));

        assertEquals(BigInteger.ZERO, RuleCellCodec.fromMagnitude(new byte[] {0}));
        assertEquals(BigInteger.valueOf(128), RuleCellCodec.fromMagnitude(new byte[] {(byte) 0x80}));
        assertEquals(BigInteger.valueOf(256), RuleCellCodec.fromMagnitude(new byte[] {0x01, 0x00}));
    }

    @Test
    void unsignedIntegerZeroRoundTrips() {
        // The UIntegerGreater* zero path (only the signed-int zero is covered elsewhere).
        RuleCell zero = new RuleCell.UIntegerGreaterValue(BigInteger.ZERO);
        assertEquals(zero, RuleCellCodec.decode("RuleUIntegerGreater",
                RuleCellCodec.encode("RuleUIntegerGreater", zero)));
    }

    // --- Taxonomy 4: cellFamily resolution incl. the prefix-ordering guards ---

    @Test
    void cellFamilyResolvesEachColumn() {
        // StringArrayEqual* must resolve before StringEqual, and UIntegerGreater* before
        // IntegerGreater — pin the mapping for both pairs plus the RawCell echo.
        assertEquals("RuleStringArrayEqual",
                RuleCellCodec.cellFamily(new RuleCell.StringArrayEqualValue(Collections.singletonList("a"))));
        assertEquals("RuleStringEqual", RuleCellCodec.cellFamily(new RuleCell.StringEqualValue("a")));
        assertEquals("RuleUIntegerGreater", RuleCellCodec.cellFamily(new RuleCell.UIntegerGreaterValue(BigInteger.ONE)));
        assertEquals("RuleIntegerGreater", RuleCellCodec.cellFamily(new RuleCell.IntegerGreaterValue(BigInteger.ONE)));
        assertEquals("RuleFiatAmount", RuleCellCodec.cellFamily(new RuleCell.RawCell("RuleFiatAmount", new byte[0])));
    }

    // --- Taxonomy 5: blockchain numeric passthrough (cell-level) ---

    @Test
    void contractAddressBlockchainNumericPassthroughRoundTrips() {
        // A blockchain value unknown to this SDK survives as its numeric string form.
        RuleCell.DestinationContractAddress original =
                new RuleCell.DestinationContractAddress("0xabc", "Token", "TKN", "4242");
        RuleCell decoded = RuleCellCodec.decode("RuleDestination",
                RuleCellCodec.encode("RuleDestination", original));
        assertEquals(original, decoded);
        assertEquals("4242", ((RuleCell.DestinationContractAddress) decoded).blockchain);
    }

    @Test
    void contractAddressUnparseableBlockchainThrows() {
        // A non-numeric, non-enum blockchain name is a caller error surfaced by parseInt.
        assertThrows(NumberFormatException.class, () -> RuleCellCodec.encode("RuleDestination",
                new RuleCell.DestinationContractAddress("0xabc", "Token", "TKN", "NotAChain")));
    }
}
