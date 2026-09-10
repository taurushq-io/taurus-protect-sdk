package model

import "math/big"

// RuleCell is a typed transaction-rule cell. Each concrete type corresponds to a
// column type and one of its cell-type enum values in the governance protobuf
// schema (e.g. FiatAmountRange is the RuleFiatAmountRange cell of a
// RuleFiatAmount column).
//
// The *Any cell types are cells like any other. Being the protobuf zero
// values, their serialized form happens to be the empty cell. A nil RuleCell
// is treated as an empty cell at encode; decode yields nil only for cells
// under columns that have no cell types at all (RuleAny or unknown columns).
//
// Cells decoded from a container built with a newer schema than this SDK
// (unknown cell type, or unknown protobuf fields inside the cell) are
// represented as RawCell and re-encoded verbatim, so no data is ever silently
// dropped. See RawCell.
type RuleCell interface {
	isRuleCell()

	// Kind is the cell's discriminator — the variant name, e.g. "SourceInternalWallet"
	// or "RawCell". Java kind(), Python .kind and TS .kind all expose the same string,
	// so a caller can branch or log without a Go type switch and without a per-language
	// mapping table.
	Kind() string
}

// RawCell preserves a cell whose type or content is unknown to this SDK
// version. Its bytes are re-emitted verbatim on encode.
type RawCell struct {
	// ColumnType is the column type string the cell was decoded under
	// (may itself be a numeric passthrough for unknown column types).
	ColumnType string
	// Payload is the exact wire content of the cell. Named to match Java, Python
	// and TS, which all call it payload.
	Payload []byte
}

// --- RuleFiatAmount column ---

// FiatAmountAny matches any fiat amount.
type FiatAmountAny struct{}

// FiatAmountIsZero matches a zero fiat amount.
type FiatAmountIsZero struct{}

// FiatAmountRange matches a fiat amount within [MinAmount, MaxAmount].
type FiatAmountRange struct {
	// MinAmount is the inclusive lower bound as a decimal string.
	MinAmount string
	// MaxAmount is the inclusive upper bound as a decimal string.
	MaxAmount string
}

// --- RuleSource column ---

// SourceAny matches any transaction source.
type SourceAny struct{}

// SourceInternalWallet matches an internal wallet source by derivation path.
type SourceInternalWallet struct {
	Path string
}

// SourceInternalAddress matches an internal address source.
type SourceInternalAddress struct {
	Address string
	Path    string
}

// SourceAnyExchange matches any exchange source.
type SourceAnyExchange struct{}

// SourceExchange matches a specific exchange source by label.
type SourceExchange struct {
	Label string
}

// SourceExternalAddress matches an external address source.
type SourceExternalAddress struct {
	Address string
	Memo    string
}

// --- RuleDestination column ---

// DestinationAny matches any transaction destination.
type DestinationAny struct{}

// DestinationInternalWallet matches an internal wallet destination by derivation path.
type DestinationInternalWallet struct {
	Path string
}

// DestinationInternalAddress matches an internal address destination.
type DestinationInternalAddress struct {
	Address string
	Path    string
}

// DestinationExternalAddress matches an external (whitelisted) address destination.
type DestinationExternalAddress struct {
	Address string
	Memo    string
}

// DestinationAnyExchange matches any exchange destination.
type DestinationAnyExchange struct{}

// DestinationExchange matches a specific exchange destination.
type DestinationExchange struct {
	Label string
	Memo  string
}

// DestinationContractAddress matches a specific contract address destination.
type DestinationContractAddress struct {
	Address string
	Name    string
	Symbol  string
	// Blockchain is the blockchain enum value name (e.g. "ETH").
	Blockchain string
}

// DestinationAnyExternalAddress matches any external address destination.
type DestinationAnyExternalAddress struct{}

// DestinationAnyContractAddress matches any contract address destination.
type DestinationAnyContractAddress struct{}

// --- RuleStringEqual column (payload is the raw UTF-8 string bytes) ---

// StringEqualAny matches any string value.
type StringEqualAny struct{}

// StringEqualEmpty matches an empty string value.
type StringEqualEmpty struct{}

// StringEqualValue matches a specific string value.
type StringEqualValue struct {
	Value string
}

// --- RuleBytesEqual column (payload is the raw bytes) ---

// BytesEqualAny matches any bytes value.
type BytesEqualAny struct{}

// BytesEqualEmpty matches an empty bytes value.
type BytesEqualEmpty struct{}

// BytesEqualValue matches a specific bytes value.
type BytesEqualValue struct {
	Value []byte
}

// --- RuleStringArrayEqual column ---

// StringArrayEqualAny matches any string array.
type StringArrayEqualAny struct{}

// StringArrayEqualEmpty matches an empty string array.
type StringArrayEqualEmpty struct{}

// StringArrayEqualValue matches a specific string array.
type StringArrayEqualValue struct {
	Values []string
}

// --- RuleIntegerGreater column ---

// IntegerGreaterAny matches any signed integer.
type IntegerGreaterAny struct{}

// IntegerGreaterValue matches signed integers greater than Value.
//
// The wire format carries the magnitude with the sign expressed by the cell
// type enum (RuleIntegerGreaterValue vs RuleIntegerGreaterNegValue); the codec
// selects the arm from Value's sign, so callers just set the signed value.
type IntegerGreaterValue struct {
	Value *big.Int
}

// --- RuleUIntegerGreater column ---

// UIntegerGreaterAny matches any unsigned integer.
type UIntegerGreaterAny struct{}

// UIntegerGreaterIsZero matches a zero unsigned integer.
type UIntegerGreaterIsZero struct{}

// UIntegerGreaterValue matches unsigned integers greater than Value (must be non-negative).
type UIntegerGreaterValue struct {
	Value *big.Int
}

// UIntegerGreaterIsEqual matches unsigned integers equal to Value (must be non-negative).
type UIntegerGreaterIsEqual struct {
	Value *big.Int
}

// --- RuleWhitelistedContract column ---

// WhitelistedContractAny matches any whitelisted contract.
type WhitelistedContractAny struct{}

// WhitelistedContractAddress matches a specific whitelisted contract address
// (the RuleWhitelistedContract_RuleDestinationContractAddress cell type).
type WhitelistedContractAddress struct {
	Address string
	Name    string
	Symbol  string
	// Blockchain is the blockchain enum value name (e.g. "ETH").
	Blockchain string
}

func (RawCell) isRuleCell()                       {}
func (FiatAmountAny) isRuleCell()                 {}
func (FiatAmountIsZero) isRuleCell()              {}
func (FiatAmountRange) isRuleCell()               {}
func (SourceAny) isRuleCell()                     {}
func (SourceInternalWallet) isRuleCell()          {}
func (SourceInternalAddress) isRuleCell()         {}
func (SourceAnyExchange) isRuleCell()             {}
func (SourceExchange) isRuleCell()                {}
func (SourceExternalAddress) isRuleCell()         {}
func (DestinationAny) isRuleCell()                {}
func (DestinationInternalWallet) isRuleCell()     {}
func (DestinationInternalAddress) isRuleCell()    {}
func (DestinationExternalAddress) isRuleCell()    {}
func (DestinationAnyExchange) isRuleCell()        {}
func (DestinationExchange) isRuleCell()           {}
func (DestinationContractAddress) isRuleCell()    {}
func (DestinationAnyExternalAddress) isRuleCell() {}
func (DestinationAnyContractAddress) isRuleCell() {}
func (StringEqualAny) isRuleCell()                {}
func (StringEqualEmpty) isRuleCell()              {}
func (StringEqualValue) isRuleCell()              {}
func (BytesEqualAny) isRuleCell()                 {}
func (BytesEqualEmpty) isRuleCell()               {}
func (BytesEqualValue) isRuleCell()               {}
func (StringArrayEqualAny) isRuleCell()           {}
func (StringArrayEqualEmpty) isRuleCell()         {}
func (StringArrayEqualValue) isRuleCell()         {}
func (IntegerGreaterAny) isRuleCell()             {}
func (IntegerGreaterValue) isRuleCell()           {}
func (UIntegerGreaterAny) isRuleCell()            {}
func (UIntegerGreaterIsZero) isRuleCell()         {}
func (UIntegerGreaterValue) isRuleCell()          {}
func (UIntegerGreaterIsEqual) isRuleCell()        {}
func (WhitelistedContractAny) isRuleCell()        {}
func (WhitelistedContractAddress) isRuleCell()    {}

// Kind returns each variant's discriminator string; see RuleCell.Kind.
func (RawCell) Kind() string                       { return "RawCell" }
func (FiatAmountAny) Kind() string                 { return "FiatAmountAny" }
func (FiatAmountIsZero) Kind() string              { return "FiatAmountIsZero" }
func (FiatAmountRange) Kind() string               { return "FiatAmountRange" }
func (SourceAny) Kind() string                     { return "SourceAny" }
func (SourceInternalWallet) Kind() string          { return "SourceInternalWallet" }
func (SourceInternalAddress) Kind() string         { return "SourceInternalAddress" }
func (SourceAnyExchange) Kind() string             { return "SourceAnyExchange" }
func (SourceExchange) Kind() string                { return "SourceExchange" }
func (SourceExternalAddress) Kind() string         { return "SourceExternalAddress" }
func (DestinationAny) Kind() string                { return "DestinationAny" }
func (DestinationInternalWallet) Kind() string     { return "DestinationInternalWallet" }
func (DestinationInternalAddress) Kind() string    { return "DestinationInternalAddress" }
func (DestinationExternalAddress) Kind() string    { return "DestinationExternalAddress" }
func (DestinationAnyExchange) Kind() string        { return "DestinationAnyExchange" }
func (DestinationExchange) Kind() string           { return "DestinationExchange" }
func (DestinationContractAddress) Kind() string    { return "DestinationContractAddress" }
func (DestinationAnyExternalAddress) Kind() string { return "DestinationAnyExternalAddress" }
func (DestinationAnyContractAddress) Kind() string { return "DestinationAnyContractAddress" }
func (StringEqualAny) Kind() string                { return "StringEqualAny" }
func (StringEqualEmpty) Kind() string              { return "StringEqualEmpty" }
func (StringEqualValue) Kind() string              { return "StringEqualValue" }
func (BytesEqualAny) Kind() string                 { return "BytesEqualAny" }
func (BytesEqualEmpty) Kind() string               { return "BytesEqualEmpty" }
func (BytesEqualValue) Kind() string               { return "BytesEqualValue" }
func (StringArrayEqualAny) Kind() string           { return "StringArrayEqualAny" }
func (StringArrayEqualEmpty) Kind() string         { return "StringArrayEqualEmpty" }
func (StringArrayEqualValue) Kind() string         { return "StringArrayEqualValue" }
func (IntegerGreaterAny) Kind() string             { return "IntegerGreaterAny" }
func (IntegerGreaterValue) Kind() string           { return "IntegerGreaterValue" }
func (UIntegerGreaterAny) Kind() string            { return "UIntegerGreaterAny" }
func (UIntegerGreaterIsZero) Kind() string         { return "UIntegerGreaterIsZero" }
func (UIntegerGreaterValue) Kind() string          { return "UIntegerGreaterValue" }
func (UIntegerGreaterIsEqual) Kind() string        { return "UIntegerGreaterIsEqual" }
func (WhitelistedContractAny) Kind() string        { return "WhitelistedContractAny" }
func (WhitelistedContractAddress) Kind() string    { return "WhitelistedContractAddress" }
