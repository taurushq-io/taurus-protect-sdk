package com.taurushq.sdk.protect.client.model.rulescontainer;

import java.math.BigInteger;
import java.util.List;

import org.apache.commons.lang3.builder.EqualsBuilder;
import org.apache.commons.lang3.builder.HashCodeBuilder;
import org.apache.commons.lang3.builder.ToStringBuilder;

/**
 * Typed transaction-rule cell (cross-SDK contract).
 *
 * <p>Each concrete subclass corresponds to a column type and one of its
 * cell-type enum values in the governance protobuf schema (e.g.
 * {@link FiatAmountRange} is the {@code RuleFiatAmountRange} cell of a
 * {@code RuleFiatAmount} column). The {@code *Any} cells are the protobuf zero
 * values, so their serialized form is the empty cell. Cells whose type/content
 * this SDK version cannot represent losslessly are preserved verbatim as
 * {@link RawCell}, so nothing is ever silently dropped.
 *
 * <p>The same set exists in the Go, Python and TypeScript SDKs, and the wire
 * encoding is pinned by the shared golden vectors at
 * {@code scripts/resources/governance-cell-vectors.json}. Java 8 has no sealed
 * types; the closed set is enforced by convention and the shared vectors.
 *
 * <p>Encode/decode against a column's byte cells lives in {@link RuleCellCodec}.
 */
public abstract class RuleCell {

    /** The model type name (matches the other SDKs' discriminator). */
    public abstract String kind();

    @Override
    public boolean equals(Object o) {
        return EqualsBuilder.reflectionEquals(this, o, false);
    }

    @Override
    public int hashCode() {
        return HashCodeBuilder.reflectionHashCode(this);
    }

    @Override
    public String toString() {
        return ToStringBuilder.reflectionToString(this);
    }

    /** Shared base for cells that carry no fields (the {@code *Any}/{@code *IsZero} arms). */
    private abstract static class NoField extends RuleCell {
        NoField() {
            super();
        }
    }

    // --- RuleFiatAmount ---
    public static final class FiatAmountAny extends NoField {
        @Override
        public String kind() { return "FiatAmountAny"; }
    }

    public static final class FiatAmountIsZero extends NoField {
        @Override
        public String kind() { return "FiatAmountIsZero"; }
    }

    public static final class FiatAmountRange extends RuleCell {
        public final String minAmount;
        public final String maxAmount;
        public FiatAmountRange(String minAmount, String maxAmount) { this.minAmount = minAmount; this.maxAmount = maxAmount; }
        @Override
        public String kind() { return "FiatAmountRange"; }
    }

    // --- RuleSource ---
    public static final class SourceAny extends NoField {
        @Override
        public String kind() { return "SourceAny"; }
    }

    public static final class SourceInternalWallet extends RuleCell {
        public final String path;
        public SourceInternalWallet(String path) { this.path = path; }
        @Override
        public String kind() { return "SourceInternalWallet"; }
    }

    public static final class SourceInternalAddress extends RuleCell {
        public final String address;
        public final String path;
        public SourceInternalAddress(String address, String path) { this.address = address; this.path = path; }
        @Override
        public String kind() { return "SourceInternalAddress"; }
    }

    public static final class SourceAnyExchange extends NoField {
        @Override
        public String kind() { return "SourceAnyExchange"; }
    }

    public static final class SourceExchange extends RuleCell {
        public final String label;
        public SourceExchange(String label) { this.label = label; }
        @Override
        public String kind() { return "SourceExchange"; }
    }

    public static final class SourceExternalAddress extends RuleCell {
        public final String address;
        public final String memo;
        public SourceExternalAddress(String address, String memo) { this.address = address; this.memo = memo; }
        @Override
        public String kind() { return "SourceExternalAddress"; }
    }

    // --- RuleDestination ---
    public static final class DestinationAny extends NoField {
        @Override
        public String kind() { return "DestinationAny"; }
    }

    public static final class DestinationInternalWallet extends RuleCell {
        public final String path;
        public DestinationInternalWallet(String path) { this.path = path; }
        @Override
        public String kind() { return "DestinationInternalWallet"; }
    }

    public static final class DestinationInternalAddress extends RuleCell {
        public final String address;
        public final String path;
        public DestinationInternalAddress(String address, String path) { this.address = address; this.path = path; }
        @Override
        public String kind() { return "DestinationInternalAddress"; }
    }

    public static final class DestinationExternalAddress extends RuleCell {
        public final String address;
        public final String memo;
        public DestinationExternalAddress(String address, String memo) { this.address = address; this.memo = memo; }
        @Override
        public String kind() { return "DestinationExternalAddress"; }
    }

    public static final class DestinationAnyExchange extends NoField {
        @Override
        public String kind() { return "DestinationAnyExchange"; }
    }

    public static final class DestinationExchange extends RuleCell {
        public final String label;
        public final String memo;
        public DestinationExchange(String label, String memo) { this.label = label; this.memo = memo; }
        @Override
        public String kind() { return "DestinationExchange"; }
    }

    public static final class DestinationContractAddress extends RuleCell {
        public final String address;
        public final String name;
        public final String symbol;
        public final String blockchain;
        public DestinationContractAddress(String address, String name, String symbol, String blockchain) {
            this.address = address; this.name = name; this.symbol = symbol; this.blockchain = blockchain;
        }
        @Override
        public String kind() { return "DestinationContractAddress"; }
    }

    public static final class DestinationAnyExternalAddress extends NoField {
        @Override
        public String kind() { return "DestinationAnyExternalAddress"; }
    }

    public static final class DestinationAnyContractAddress extends NoField {
        @Override
        public String kind() { return "DestinationAnyContractAddress"; }
    }

    // --- RuleStringEqual (payload is raw UTF-8 string bytes) ---
    public static final class StringEqualAny extends NoField {
        @Override
        public String kind() { return "StringEqualAny"; }
    }

    public static final class StringEqualEmpty extends NoField {
        @Override
        public String kind() { return "StringEqualEmpty"; }
    }

    public static final class StringEqualValue extends RuleCell {
        public final String value;
        public StringEqualValue(String value) { this.value = value; }
        @Override
        public String kind() { return "StringEqualValue"; }
    }

    // --- RuleBytesEqual (payload is raw bytes) ---
    public static final class BytesEqualAny extends NoField {
        @Override
        public String kind() { return "BytesEqualAny"; }
    }

    public static final class BytesEqualEmpty extends NoField {
        @Override
        public String kind() { return "BytesEqualEmpty"; }
    }

    public static final class BytesEqualValue extends RuleCell {
        public final byte[] value;
        public BytesEqualValue(byte[] value) { this.value = value; }
        @Override
        public String kind() { return "BytesEqualValue"; }
    }

    // --- RuleStringArrayEqual ---
    public static final class StringArrayEqualAny extends NoField {
        @Override
        public String kind() { return "StringArrayEqualAny"; }
    }

    public static final class StringArrayEqualEmpty extends NoField {
        @Override
        public String kind() { return "StringArrayEqualEmpty"; }
    }

    public static final class StringArrayEqualValue extends RuleCell {
        public final List<String> values;
        public StringArrayEqualValue(List<String> values) { this.values = values; }
        @Override
        public String kind() { return "StringArrayEqualValue"; }
    }

    // --- RuleIntegerGreater (sign selects the Value/NegValue wire arm) ---
    public static final class IntegerGreaterAny extends NoField {
        @Override
        public String kind() { return "IntegerGreaterAny"; }
    }

    public static final class IntegerGreaterValue extends RuleCell {
        public final BigInteger value;
        public IntegerGreaterValue(BigInteger value) { this.value = value; }
        @Override
        public String kind() { return "IntegerGreaterValue"; }
    }

    // --- RuleUIntegerGreater ---
    public static final class UIntegerGreaterAny extends NoField {
        @Override
        public String kind() { return "UIntegerGreaterAny"; }
    }

    public static final class UIntegerGreaterIsZero extends NoField {
        @Override
        public String kind() { return "UIntegerGreaterIsZero"; }
    }

    public static final class UIntegerGreaterValue extends RuleCell {
        public final BigInteger value;
        public UIntegerGreaterValue(BigInteger value) { this.value = value; }
        @Override
        public String kind() { return "UIntegerGreaterValue"; }
    }

    public static final class UIntegerGreaterIsEqual extends RuleCell {
        public final BigInteger value;
        public UIntegerGreaterIsEqual(BigInteger value) { this.value = value; }
        @Override
        public String kind() { return "UIntegerGreaterIsEqual"; }
    }

    // --- RuleWhitelistedContract ---
    public static final class WhitelistedContractAny extends NoField {
        @Override
        public String kind() { return "WhitelistedContractAny"; }
    }

    public static final class WhitelistedContractAddress extends RuleCell {
        public final String address;
        public final String name;
        public final String symbol;
        public final String blockchain;
        public WhitelistedContractAddress(String address, String name, String symbol, String blockchain) {
            this.address = address; this.name = name; this.symbol = symbol; this.blockchain = blockchain;
        }
        @Override
        public String kind() { return "WhitelistedContractAddress"; }
    }

    // --- Fallback for cells unknown to this SDK version (preserved verbatim) ---
    public static final class RawCell extends RuleCell {
        public final String columnType;
        public final byte[] payload;
        public RawCell(String columnType, byte[] payload) { this.columnType = columnType; this.payload = payload; }
        @Override
        public String kind() { return "RawCell"; }
    }
}
