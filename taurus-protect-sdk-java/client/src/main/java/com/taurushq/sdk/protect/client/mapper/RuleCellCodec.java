package com.taurushq.sdk.protect.client.mapper;

import java.math.BigInteger;
import java.nio.charset.StandardCharsets;
import java.util.Arrays;

import com.google.protobuf.ByteString;
import com.google.protobuf.InvalidProtocolBufferException;
import com.taurushq.sdk.protect.client.model.rulescontainer.RuleCell;
import com.taurushq.sdk.protect.proto.v1.RequestReply;

/**
 * Codec between the typed {@link RuleCell} union and transaction-rule cell wire bytes.
 *
 * <p>Mirrors the Go, Python and TypeScript SDKs. A cell is a wrapper message
 * selected by the column's type, {@code {type: <cell-type enum>, payload:
 * <bytes>}}. Message-family payloads nest a serialized sub-message;
 * {@code RuleStringEqual}/{@code RuleBytesEqual} carry raw scalar bytes; the
 * integer families carry big-endian magnitude bytes with the sign expressed by
 * the enum arm. The {@code *Any} cells are the protobuf zero values, so their
 * serialized form is the empty cell.
 *
 * <p>Decode never fails: any cell that does not round-trip byte-identically
 * through the typed layer (unknown column type, unknown cell type, unknown
 * protobuf sub-fields) is returned as a {@link RuleCell.RawCell} so nothing is
 * silently dropped. The byte encoding is pinned across all four SDKs by the
 * shared golden vectors at {@code scripts/resources/governance-cell-vectors.json}.
 */
public final class RuleCellCodec {

    private RuleCellCodec() {
    }

    /** Column type constants (match the proto message names). */
    public static final String COL_FIAT_AMOUNT = "RuleFiatAmount";
    public static final String COL_SOURCE = "RuleSource";
    public static final String COL_DESTINATION = "RuleDestination";
    public static final String COL_WHITELISTED_CONTRACT = "RuleWhitelistedContract";
    public static final String COL_STRING_EQUAL = "RuleStringEqual";
    public static final String COL_BYTES_EQUAL = "RuleBytesEqual";
    public static final String COL_STRING_ARRAY_EQUAL = "RuleStringArrayEqual";
    public static final String COL_INTEGER_GREATER = "RuleIntegerGreater";
    public static final String COL_UINTEGER_GREATER = "RuleUIntegerGreater";

    /**
     * Returns the column type a typed cell belongs to.
     *
     * @param cell the typed cell
     * @return the column-type string
     */
    public static String cellFamily(RuleCell cell) {
        if (cell instanceof RuleCell.RawCell) {
            return ((RuleCell.RawCell) cell).columnType;
        }
        String k = cell.kind();
        if (k.startsWith("FiatAmount")) {
            return COL_FIAT_AMOUNT;
        }
        if (k.startsWith("Source")) {
            return COL_SOURCE;
        }
        if (k.startsWith("Destination")) {
            return COL_DESTINATION;
        }
        if (k.startsWith("WhitelistedContract")) {
            return COL_WHITELISTED_CONTRACT;
        }
        if (k.startsWith("StringArrayEqual")) {
            return COL_STRING_ARRAY_EQUAL;
        }
        if (k.startsWith("StringEqual")) {
            return COL_STRING_EQUAL;
        }
        if (k.startsWith("BytesEqual")) {
            return COL_BYTES_EQUAL;
        }
        if (k.startsWith("UIntegerGreater")) {
            return COL_UINTEGER_GREATER;
        }
        if (k.startsWith("IntegerGreater")) {
            return COL_INTEGER_GREATER;
        }
        throw new IllegalArgumentException("unknown rule cell kind " + k);
    }

    /**
     * Encodes a typed cell to the wrapped protobuf bytes for a line cell.
     *
     * @param columnType the column type the cell is placed in (may be empty to skip the check)
     * @param cell       the typed cell
     * @return the serialized cell bytes
     */
    public static byte[] encode(String columnType, RuleCell cell) {
        if (cell instanceof RuleCell.RawCell) {
            return ((RuleCell.RawCell) cell).payload;
        }
        String family = cellFamily(cell);
        if (columnType != null && !columnType.isEmpty() && !columnType.equals(family)) {
            throw new IllegalArgumentException("cell " + cell.kind() + " is not valid for column type " + columnType);
        }
        switch (cell.kind()) {
            // RuleFiatAmount
            case "FiatAmountAny":
                return fiat(RequestReply.RuleFiatAmount.RuleFiatAmountType.RuleFiatAmountAny, null);
            case "FiatAmountIsZero":
                return fiat(RequestReply.RuleFiatAmount.RuleFiatAmountType.RuleFiatAmountIsZero, null);
            case "FiatAmountRange": {
                RuleCell.FiatAmountRange c = (RuleCell.FiatAmountRange) cell;
                byte[] inner = RequestReply.RuleFiatAmountRange.newBuilder()
                        .setMinAmount(c.minAmount).setMaxAmount(c.maxAmount).build().toByteArray();
                return fiat(RequestReply.RuleFiatAmount.RuleFiatAmountType.RuleFiatAmountRange, inner);
            }

            // RuleSource
            case "SourceAny":
                return source(RequestReply.RuleSource.RuleSourceType.RuleSourceAny, null);
            case "SourceAnyExchange":
                return source(RequestReply.RuleSource.RuleSourceType.RuleSourceAnyExchange, null);
            case "SourceInternalWallet": {
                RuleCell.SourceInternalWallet c = (RuleCell.SourceInternalWallet) cell;
                byte[] inner = RequestReply.RuleSourceInternalWallet.newBuilder().setPath(c.path).build().toByteArray();
                return source(RequestReply.RuleSource.RuleSourceType.RuleSourceInternalWallet, inner);
            }
            case "SourceInternalAddress": {
                RuleCell.SourceInternalAddress c = (RuleCell.SourceInternalAddress) cell;
                byte[] inner = RequestReply.RuleSourceInternalAddress.newBuilder()
                        .setAddress(c.address).setPath(c.path).build().toByteArray();
                return source(RequestReply.RuleSource.RuleSourceType.RuleSourceInternalAddress, inner);
            }
            case "SourceExchange": {
                RuleCell.SourceExchange c = (RuleCell.SourceExchange) cell;
                byte[] inner = RequestReply.RuleSourceExchange.newBuilder().setLabel(c.label).build().toByteArray();
                return source(RequestReply.RuleSource.RuleSourceType.RuleSourceExchange, inner);
            }
            case "SourceExternalAddress": {
                RuleCell.SourceExternalAddress c = (RuleCell.SourceExternalAddress) cell;
                byte[] inner = RequestReply.RuleSourceExternalAddress.newBuilder()
                        .setAddress(c.address).setMemo(c.memo).build().toByteArray();
                return source(RequestReply.RuleSource.RuleSourceType.RuleSourceExternalAddress, inner);
            }

            // RuleDestination
            case "DestinationAny":
                return dest(RequestReply.RuleDestination.RuleDestinationType.RuleDestinationAny, null);
            case "DestinationAnyExchange":
                return dest(RequestReply.RuleDestination.RuleDestinationType.RuleDestinationAnyExchange, null);
            case "DestinationAnyExternalAddress":
                return dest(RequestReply.RuleDestination.RuleDestinationType.RuleDestinationAnyExternalAddress, null);
            case "DestinationAnyContractAddress":
                return dest(RequestReply.RuleDestination.RuleDestinationType.RuleDestinationAnyContractAddress, null);
            case "DestinationInternalWallet": {
                RuleCell.DestinationInternalWallet c = (RuleCell.DestinationInternalWallet) cell;
                byte[] inner = RequestReply.RuleDestinationInternalWallet.newBuilder().setPath(c.path).build().toByteArray();
                return dest(RequestReply.RuleDestination.RuleDestinationType.RuleDestinationInternalWallet, inner);
            }
            case "DestinationInternalAddress": {
                RuleCell.DestinationInternalAddress c = (RuleCell.DestinationInternalAddress) cell;
                byte[] inner = RequestReply.RuleDestinationInternalAddress.newBuilder()
                        .setAddress(c.address).setPath(c.path).build().toByteArray();
                return dest(RequestReply.RuleDestination.RuleDestinationType.RuleDestinationInternalAddress, inner);
            }
            case "DestinationExternalAddress": {
                RuleCell.DestinationExternalAddress c = (RuleCell.DestinationExternalAddress) cell;
                byte[] inner = RequestReply.RuleDestinationExternalAddress.newBuilder()
                        .setAddress(c.address).setMemo(c.memo).build().toByteArray();
                return dest(RequestReply.RuleDestination.RuleDestinationType.RuleDestinationExternalAddress, inner);
            }
            case "DestinationExchange": {
                RuleCell.DestinationExchange c = (RuleCell.DestinationExchange) cell;
                byte[] inner = RequestReply.RuleDestinationExchange.newBuilder()
                        .setLabel(c.label).setMemo(c.memo).build().toByteArray();
                return dest(RequestReply.RuleDestination.RuleDestinationType.RuleDestinationExchange, inner);
            }
            case "DestinationContractAddress": {
                RuleCell.DestinationContractAddress c = (RuleCell.DestinationContractAddress) cell;
                return dest(RequestReply.RuleDestination.RuleDestinationType.RuleDestinationContractAddress,
                        contractAddress(c.address, c.name, c.symbol, c.blockchain));
            }

            // RuleWhitelistedContract
            case "WhitelistedContractAny":
                return whitelisted(RequestReply.RuleWhitelistedContract.RuleWhitelistedContractType.RuleWhitelistedContractAny, null);
            case "WhitelistedContractAddress": {
                RuleCell.WhitelistedContractAddress c = (RuleCell.WhitelistedContractAddress) cell;
                return whitelisted(
                        RequestReply.RuleWhitelistedContract.RuleWhitelistedContractType.RuleWhitelistedContract_RuleDestinationContractAddress,
                        contractAddress(c.address, c.name, c.symbol, c.blockchain));
            }

            // RuleStringEqual (payload = raw UTF-8 string bytes)
            case "StringEqualAny":
                return stringEqual(RequestReply.RuleStringEqual.RuleStringEqualType.RuleStringEqualAny, null);
            case "StringEqualEmpty":
                return stringEqual(RequestReply.RuleStringEqual.RuleStringEqualType.RuleStringEqualEmpty, null);
            case "StringEqualValue": {
                RuleCell.StringEqualValue c = (RuleCell.StringEqualValue) cell;
                return stringEqual(RequestReply.RuleStringEqual.RuleStringEqualType.RuleStringEqualValue,
                        c.value.getBytes(StandardCharsets.UTF_8));
            }

            // RuleBytesEqual (payload = raw bytes)
            case "BytesEqualAny":
                return bytesEqual(RequestReply.RuleBytesEqual.RuleBytesEqualType.RuleBytesEqualAny, null);
            case "BytesEqualEmpty":
                return bytesEqual(RequestReply.RuleBytesEqual.RuleBytesEqualType.RuleBytesEqualEmpty, null);
            case "BytesEqualValue": {
                RuleCell.BytesEqualValue c = (RuleCell.BytesEqualValue) cell;
                return bytesEqual(RequestReply.RuleBytesEqual.RuleBytesEqualType.RuleBytesEqualValue, c.value);
            }

            // RuleStringArrayEqual
            case "StringArrayEqualAny":
                return stringArray(RequestReply.RuleStringArrayEqual.RuleStringArrayEqualType.RuleStringArrayEqualAny, null);
            case "StringArrayEqualEmpty":
                return stringArray(RequestReply.RuleStringArrayEqual.RuleStringArrayEqualType.RuleStringArrayEqualEmpty, null);
            case "StringArrayEqualValue": {
                RuleCell.StringArrayEqualValue c = (RuleCell.StringArrayEqualValue) cell;
                byte[] inner = RequestReply.RuleStringArrayEqualValue.newBuilder()
                        .addAllValues(c.values).build().toByteArray();
                return stringArray(RequestReply.RuleStringArrayEqual.RuleStringArrayEqualType.RuleStringArrayEqualValue, inner);
            }

            default:
                // The integer families live in their own method to keep this dispatch
                // within the project's method-length limit.
                return encodeIntegerCell(cell);
        }
    }

    /**
     * Encodes the {@code RuleIntegerGreater} / {@code RuleUIntegerGreater} families,
     * whose payload is a big-endian magnitude with the sign selecting the enum arm.
     *
     * @param cell the typed cell
     * @return the serialized cell bytes
     */
    private static byte[] encodeIntegerCell(RuleCell cell) {
        switch (cell.kind()) {
            // RuleIntegerGreater (magnitude payload, sign selects the arm)
            case "IntegerGreaterAny":
                return integer(RequestReply.RuleIntegerGreater.RuleIntegerGreaterType.RuleIntegerGreaterAny, null);
            case "IntegerGreaterValue": {
                RuleCell.IntegerGreaterValue c = (RuleCell.IntegerGreaterValue) cell;
                RequestReply.RuleIntegerGreater.RuleIntegerGreaterType arm = c.value.signum() < 0
                        ? RequestReply.RuleIntegerGreater.RuleIntegerGreaterType.RuleIntegerGreaterNegValue
                        : RequestReply.RuleIntegerGreater.RuleIntegerGreaterType.RuleIntegerGreaterValue;
                return integer(arm, magnitude(c.value));
            }

            // RuleUIntegerGreater
            case "UIntegerGreaterAny":
                return uinteger(RequestReply.RuleUIntegerGreater.RuleUIntegerGreaterType.RuleUIntegerGreaterAny, null);
            case "UIntegerGreaterIsZero":
                return uinteger(RequestReply.RuleUIntegerGreater.RuleUIntegerGreaterType.RuleUIntegerGreaterIsZero, null);
            case "UIntegerGreaterValue": {
                RuleCell.UIntegerGreaterValue c = (RuleCell.UIntegerGreaterValue) cell;
                requireNonNegative(c.value, "UIntegerGreaterValue");
                return uinteger(RequestReply.RuleUIntegerGreater.RuleUIntegerGreaterType.RuleUIntegerGreaterValue, magnitude(c.value));
            }
            case "UIntegerGreaterIsEqual": {
                RuleCell.UIntegerGreaterIsEqual c = (RuleCell.UIntegerGreaterIsEqual) cell;
                requireNonNegative(c.value, "UIntegerGreaterIsEqual");
                return uinteger(RequestReply.RuleUIntegerGreater.RuleUIntegerGreaterType.RuleUIntegerGreaterIsEqual, magnitude(c.value));
            }

            default:
                throw new IllegalArgumentException("unsupported rule cell kind " + cell.kind());
        }
    }

    /**
     * Decodes a line cell. Falls back to {@link RuleCell.RawCell} if it does not
     * round-trip byte-identically.
     *
     * @param columnType the column type the cell is placed in
     * @param data       the serialized cell bytes
     * @return the typed cell, or a {@code RawCell} when it cannot be represented losslessly
     */
    public static RuleCell decode(String columnType, byte[] data) {
        RuleCell typed = decodeTyped(columnType, data);
        // Lossless guard: preserve verbatim if the typed value does not re-encode exactly.
        if (typed == null || !Arrays.equals(encode(columnType, typed), data)) {
            return new RuleCell.RawCell(columnType, data);
        }
        return typed;
    }

    private static RuleCell decodeTyped(String columnType, byte[] data) {
        try {
            if (COL_FIAT_AMOUNT.equals(columnType)) {
                RequestReply.RuleFiatAmount w = RequestReply.RuleFiatAmount.parseFrom(data);
                switch (w.getType()) {
                    case RuleFiatAmountAny:
                        return new RuleCell.FiatAmountAny();
                    case RuleFiatAmountIsZero:
                        return new RuleCell.FiatAmountIsZero();
                    case RuleFiatAmountRange: {
                        RequestReply.RuleFiatAmountRange i = RequestReply.RuleFiatAmountRange.parseFrom(w.getPayload());
                        return new RuleCell.FiatAmountRange(i.getMinAmount(), i.getMaxAmount());
                    }
                    default:
                        return null;
                }
            }
            if (COL_SOURCE.equals(columnType)) {
                RequestReply.RuleSource w = RequestReply.RuleSource.parseFrom(data);
                switch (w.getType()) {
                    case RuleSourceAny:
                        return new RuleCell.SourceAny();
                    case RuleSourceAnyExchange:
                        return new RuleCell.SourceAnyExchange();
                    case RuleSourceInternalWallet: {
                        RequestReply.RuleSourceInternalWallet i = RequestReply.RuleSourceInternalWallet.parseFrom(w.getPayload());
                        return new RuleCell.SourceInternalWallet(i.getPath());
                    }
                    case RuleSourceInternalAddress: {
                        RequestReply.RuleSourceInternalAddress i = RequestReply.RuleSourceInternalAddress.parseFrom(w.getPayload());
                        return new RuleCell.SourceInternalAddress(i.getAddress(), i.getPath());
                    }
                    case RuleSourceExchange: {
                        RequestReply.RuleSourceExchange i = RequestReply.RuleSourceExchange.parseFrom(w.getPayload());
                        return new RuleCell.SourceExchange(i.getLabel());
                    }
                    case RuleSourceExternalAddress: {
                        RequestReply.RuleSourceExternalAddress i = RequestReply.RuleSourceExternalAddress.parseFrom(w.getPayload());
                        return new RuleCell.SourceExternalAddress(i.getAddress(), i.getMemo());
                    }
                    default:
                        return null;
                }
            }
            if (COL_DESTINATION.equals(columnType)) {
                RequestReply.RuleDestination w = RequestReply.RuleDestination.parseFrom(data);
                switch (w.getType()) {
                    case RuleDestinationAny:
                        return new RuleCell.DestinationAny();
                    case RuleDestinationAnyExchange:
                        return new RuleCell.DestinationAnyExchange();
                    case RuleDestinationAnyExternalAddress:
                        return new RuleCell.DestinationAnyExternalAddress();
                    case RuleDestinationAnyContractAddress:
                        return new RuleCell.DestinationAnyContractAddress();
                    case RuleDestinationInternalWallet: {
                        RequestReply.RuleDestinationInternalWallet i = RequestReply.RuleDestinationInternalWallet.parseFrom(w.getPayload());
                        return new RuleCell.DestinationInternalWallet(i.getPath());
                    }
                    case RuleDestinationInternalAddress: {
                        RequestReply.RuleDestinationInternalAddress i = RequestReply.RuleDestinationInternalAddress.parseFrom(w.getPayload());
                        return new RuleCell.DestinationInternalAddress(i.getAddress(), i.getPath());
                    }
                    case RuleDestinationExternalAddress: {
                        RequestReply.RuleDestinationExternalAddress i = RequestReply.RuleDestinationExternalAddress.parseFrom(w.getPayload());
                        return new RuleCell.DestinationExternalAddress(i.getAddress(), i.getMemo());
                    }
                    case RuleDestinationExchange: {
                        RequestReply.RuleDestinationExchange i = RequestReply.RuleDestinationExchange.parseFrom(w.getPayload());
                        return new RuleCell.DestinationExchange(i.getLabel(), i.getMemo());
                    }
                    case RuleDestinationContractAddress: {
                        RequestReply.RuleDestinationContractAddress i = RequestReply.RuleDestinationContractAddress.parseFrom(w.getPayload());
                        return new RuleCell.DestinationContractAddress(i.getAddress(), i.getName(), i.getSymbol(), blockchainName(i.getBlockchainValue()));
                    }
                    default:
                        return null;
                }
            }
            if (COL_WHITELISTED_CONTRACT.equals(columnType)) {
                RequestReply.RuleWhitelistedContract w = RequestReply.RuleWhitelistedContract.parseFrom(data);
                switch (w.getType()) {
                    case RuleWhitelistedContractAny:
                        return new RuleCell.WhitelistedContractAny();
                    case RuleWhitelistedContract_RuleDestinationContractAddress: {
                        RequestReply.RuleDestinationContractAddress i = RequestReply.RuleDestinationContractAddress.parseFrom(w.getPayload());
                        return new RuleCell.WhitelistedContractAddress(i.getAddress(), i.getName(), i.getSymbol(), blockchainName(i.getBlockchainValue()));
                    }
                    default:
                        return null;
                }
            }
            if (COL_STRING_EQUAL.equals(columnType)) {
                RequestReply.RuleStringEqual w = RequestReply.RuleStringEqual.parseFrom(data);
                switch (w.getType()) {
                    case RuleStringEqualAny:
                        return new RuleCell.StringEqualAny();
                    case RuleStringEqualEmpty:
                        return new RuleCell.StringEqualEmpty();
                    case RuleStringEqualValue:
                        return new RuleCell.StringEqualValue(new String(w.getPayload().toByteArray(), StandardCharsets.UTF_8));
                    default:
                        return null;
                }
            }
            if (COL_BYTES_EQUAL.equals(columnType)) {
                RequestReply.RuleBytesEqual w = RequestReply.RuleBytesEqual.parseFrom(data);
                switch (w.getType()) {
                    case RuleBytesEqualAny:
                        return new RuleCell.BytesEqualAny();
                    case RuleBytesEqualEmpty:
                        return new RuleCell.BytesEqualEmpty();
                    case RuleBytesEqualValue:
                        return new RuleCell.BytesEqualValue(w.getPayload().toByteArray());
                    default:
                        return null;
                }
            }
            if (COL_STRING_ARRAY_EQUAL.equals(columnType)) {
                RequestReply.RuleStringArrayEqual w = RequestReply.RuleStringArrayEqual.parseFrom(data);
                switch (w.getType()) {
                    case RuleStringArrayEqualAny:
                        return new RuleCell.StringArrayEqualAny();
                    case RuleStringArrayEqualEmpty:
                        return new RuleCell.StringArrayEqualEmpty();
                    case RuleStringArrayEqualValue: {
                        RequestReply.RuleStringArrayEqualValue i = RequestReply.RuleStringArrayEqualValue.parseFrom(w.getPayload());
                        return new RuleCell.StringArrayEqualValue(i.getValuesList());
                    }
                    default:
                        return null;
                }
            }
            // The integer families live in their own method to keep this dispatch within
            // the project's method-length limit.
            return decodeIntegerCell(columnType, data);
        } catch (InvalidProtocolBufferException e) {
            return null;
        }
    }

    /**
     * Decodes the {@code RuleIntegerGreater} / {@code RuleUIntegerGreater} families,
     * whose payload is a big-endian magnitude with the enum arm carrying the sign.
     *
     * @param columnType the column type the cell is aligned with
     * @param data       the serialized cell
     * @return the typed cell, or null when the column type is not an integer family
     * @throws InvalidProtocolBufferException if the cell is not valid protobuf
     */
    private static RuleCell decodeIntegerCell(String columnType, byte[] data)
            throws InvalidProtocolBufferException {
        if (COL_INTEGER_GREATER.equals(columnType)) {
            RequestReply.RuleIntegerGreater w = RequestReply.RuleIntegerGreater.parseFrom(data);
            switch (w.getType()) {
                case RuleIntegerGreaterAny:
                    return new RuleCell.IntegerGreaterAny();
                case RuleIntegerGreaterValue:
                    return new RuleCell.IntegerGreaterValue(fromMagnitude(w.getPayload().toByteArray()));
                case RuleIntegerGreaterNegValue:
                    return new RuleCell.IntegerGreaterValue(fromMagnitude(w.getPayload().toByteArray()).negate());
                default:
                    return null;
            }
        }
        if (COL_UINTEGER_GREATER.equals(columnType)) {
            RequestReply.RuleUIntegerGreater w = RequestReply.RuleUIntegerGreater.parseFrom(data);
            switch (w.getType()) {
                case RuleUIntegerGreaterAny:
                    return new RuleCell.UIntegerGreaterAny();
                case RuleUIntegerGreaterIsZero:
                    return new RuleCell.UIntegerGreaterIsZero();
                case RuleUIntegerGreaterValue:
                    return new RuleCell.UIntegerGreaterValue(fromMagnitude(w.getPayload().toByteArray()));
                case RuleUIntegerGreaterIsEqual:
                    return new RuleCell.UIntegerGreaterIsEqual(fromMagnitude(w.getPayload().toByteArray()));
                default:
                    return null;
            }
        }
        // Unknown column type (incl. RuleAny): no cell family to type it.
        return null;
    }

    // --- wrapper builders ---

    private static byte[] fiat(RequestReply.RuleFiatAmount.RuleFiatAmountType type, byte[] payload) {
        RequestReply.RuleFiatAmount.Builder b = RequestReply.RuleFiatAmount.newBuilder().setType(type);
        if (payload != null) {
            b.setPayload(ByteString.copyFrom(payload));
        }
        return b.build().toByteArray();
    }

    private static byte[] source(RequestReply.RuleSource.RuleSourceType type, byte[] payload) {
        RequestReply.RuleSource.Builder b = RequestReply.RuleSource.newBuilder().setType(type);
        if (payload != null) {
            b.setPayload(ByteString.copyFrom(payload));
        }
        return b.build().toByteArray();
    }

    private static byte[] dest(RequestReply.RuleDestination.RuleDestinationType type, byte[] payload) {
        RequestReply.RuleDestination.Builder b = RequestReply.RuleDestination.newBuilder().setType(type);
        if (payload != null) {
            b.setPayload(ByteString.copyFrom(payload));
        }
        return b.build().toByteArray();
    }

    private static byte[] whitelisted(RequestReply.RuleWhitelistedContract.RuleWhitelistedContractType type, byte[] payload) {
        RequestReply.RuleWhitelistedContract.Builder b = RequestReply.RuleWhitelistedContract.newBuilder().setType(type);
        if (payload != null) {
            b.setPayload(ByteString.copyFrom(payload));
        }
        return b.build().toByteArray();
    }

    private static byte[] stringEqual(RequestReply.RuleStringEqual.RuleStringEqualType type, byte[] payload) {
        RequestReply.RuleStringEqual.Builder b = RequestReply.RuleStringEqual.newBuilder().setType(type);
        if (payload != null) {
            b.setPayload(ByteString.copyFrom(payload));
        }
        return b.build().toByteArray();
    }

    private static byte[] bytesEqual(RequestReply.RuleBytesEqual.RuleBytesEqualType type, byte[] payload) {
        RequestReply.RuleBytesEqual.Builder b = RequestReply.RuleBytesEqual.newBuilder().setType(type);
        if (payload != null) {
            b.setPayload(ByteString.copyFrom(payload));
        }
        return b.build().toByteArray();
    }

    private static byte[] stringArray(RequestReply.RuleStringArrayEqual.RuleStringArrayEqualType type, byte[] payload) {
        RequestReply.RuleStringArrayEqual.Builder b = RequestReply.RuleStringArrayEqual.newBuilder().setType(type);
        if (payload != null) {
            b.setPayload(ByteString.copyFrom(payload));
        }
        return b.build().toByteArray();
    }

    private static byte[] integer(RequestReply.RuleIntegerGreater.RuleIntegerGreaterType type, byte[] payload) {
        RequestReply.RuleIntegerGreater.Builder b = RequestReply.RuleIntegerGreater.newBuilder().setType(type);
        if (payload != null) {
            b.setPayload(ByteString.copyFrom(payload));
        }
        return b.build().toByteArray();
    }

    private static byte[] uinteger(RequestReply.RuleUIntegerGreater.RuleUIntegerGreaterType type, byte[] payload) {
        RequestReply.RuleUIntegerGreater.Builder b = RequestReply.RuleUIntegerGreater.newBuilder().setType(type);
        if (payload != null) {
            b.setPayload(ByteString.copyFrom(payload));
        }
        return b.build().toByteArray();
    }

    private static byte[] contractAddress(String address, String name, String symbol, String blockchain) {
        return RequestReply.RuleDestinationContractAddress.newBuilder()
                .setAddress(address).setName(name).setSymbol(symbol)
                .setBlockchainValue(blockchainToInt(blockchain))
                .build().toByteArray();
    }

    // --- scalar helpers (must match the other SDKs byte-for-byte) ---

    private static void requireNonNegative(BigInteger value, String kind) {
        if (value.signum() < 0) {
            throw new IllegalArgumentException(kind + " requires a non-negative value");
        }
    }

    /** Minimal big-endian magnitude bytes; zero encodes as a single 0x00 byte. */
    static byte[] magnitude(BigInteger value) {
        BigInteger a = value.abs();
        if (a.signum() == 0) {
            return new byte[] {0};
        }
        byte[] b = a.toByteArray();
        // BigInteger.toByteArray() prepends a 0x00 sign byte when the top bit is set; drop it.
        int off = (b.length > 1 && b[0] == 0) ? 1 : 0;
        return off == 0 ? b : Arrays.copyOfRange(b, off, b.length);
    }

    static BigInteger fromMagnitude(byte[] data) {
        return new BigInteger(1, data);
    }

    private static int blockchainToInt(String name) {
        try {
            return RequestReply.Blockchain.valueOf(name).getNumber();
        } catch (IllegalArgumentException e) {
            return Integer.parseInt(name); // numeric passthrough for values newer than this SDK
        }
    }

    private static String blockchainName(int value) {
        RequestReply.Blockchain b = RequestReply.Blockchain.forNumber(value);
        return b == null ? String.valueOf(value) : b.name();
    }
}
