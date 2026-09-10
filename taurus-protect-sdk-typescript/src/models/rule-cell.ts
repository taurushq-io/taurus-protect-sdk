/**
 * Typed transaction-rule cell union.
 *
 * Each variant corresponds to a column type and one of its cell-type enum
 * values in the governance protobuf schema (e.g. `FiatAmountRange` is the
 * `RuleFiatAmountRange` cell of a `RuleFiatAmount` column). The `kind`
 * discriminator carries the model type name.
 *
 * The `*Any` variants are cells like any other; being the protobuf zero
 * values, their serialized form is the empty cell. Cells whose type or content
 * this SDK version cannot fully represent (unknown cell type, unknown protobuf
 * fields inside the cell) are preserved verbatim as `RawCell`, so no data is
 * ever silently dropped.
 *
 * This union is a cross-SDK contract: the same set exists in the Go, Java, and
 * Python SDKs, and the wire encoding is pinned by the shared golden vectors at
 * `scripts/resources/governance-cell-vectors.json`.
 */
export type RuleCell =
  // RuleFiatAmount column
  | { readonly kind: 'FiatAmountAny' }
  | { readonly kind: 'FiatAmountIsZero' }
  | { readonly kind: 'FiatAmountRange'; readonly minAmount: string; readonly maxAmount: string }
  // RuleSource column
  | { readonly kind: 'SourceAny' }
  | { readonly kind: 'SourceInternalWallet'; readonly path: string }
  | { readonly kind: 'SourceInternalAddress'; readonly address: string; readonly path: string }
  | { readonly kind: 'SourceAnyExchange' }
  | { readonly kind: 'SourceExchange'; readonly label: string }
  | { readonly kind: 'SourceExternalAddress'; readonly address: string; readonly memo: string }
  // RuleDestination column
  | { readonly kind: 'DestinationAny' }
  | { readonly kind: 'DestinationInternalWallet'; readonly path: string }
  | { readonly kind: 'DestinationInternalAddress'; readonly address: string; readonly path: string }
  | { readonly kind: 'DestinationExternalAddress'; readonly address: string; readonly memo: string }
  | { readonly kind: 'DestinationAnyExchange' }
  | { readonly kind: 'DestinationExchange'; readonly label: string; readonly memo: string }
  | {
      readonly kind: 'DestinationContractAddress';
      readonly address: string;
      readonly name: string;
      readonly symbol: string;
      readonly blockchain: string;
    }
  | { readonly kind: 'DestinationAnyExternalAddress' }
  | { readonly kind: 'DestinationAnyContractAddress' }
  // RuleStringEqual column (payload is the raw UTF-8 string bytes)
  | { readonly kind: 'StringEqualAny' }
  | { readonly kind: 'StringEqualEmpty' }
  | { readonly kind: 'StringEqualValue'; readonly value: string }
  // RuleBytesEqual column (payload is the raw bytes)
  | { readonly kind: 'BytesEqualAny' }
  | { readonly kind: 'BytesEqualEmpty' }
  | { readonly kind: 'BytesEqualValue'; readonly value: Uint8Array }
  // RuleStringArrayEqual column
  | { readonly kind: 'StringArrayEqualAny' }
  | { readonly kind: 'StringArrayEqualEmpty' }
  | { readonly kind: 'StringArrayEqualValue'; readonly values: string[] }
  // RuleIntegerGreater column (sign selects the Value/NegValue wire arm)
  | { readonly kind: 'IntegerGreaterAny' }
  | { readonly kind: 'IntegerGreaterValue'; readonly value: bigint }
  // RuleUIntegerGreater column
  | { readonly kind: 'UIntegerGreaterAny' }
  | { readonly kind: 'UIntegerGreaterIsZero' }
  | { readonly kind: 'UIntegerGreaterValue'; readonly value: bigint }
  | { readonly kind: 'UIntegerGreaterIsEqual'; readonly value: bigint }
  // RuleWhitelistedContract column
  | { readonly kind: 'WhitelistedContractAny' }
  | {
      readonly kind: 'WhitelistedContractAddress';
      readonly address: string;
      readonly name: string;
      readonly symbol: string;
      readonly blockchain: string;
    }
  // Fallback: a cell whose type/content is unknown to this SDK version.
  // Re-emitted verbatim on encode.
  | { readonly kind: 'RawCell'; readonly columnType: string; readonly payload: Uint8Array };
