package mapper

import (
	"math/big"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// cellCase describes one typed transaction-rule cell together with its
// cross-SDK vector metadata. The table is the single source of truth for the
// codec round-trip test and the shared golden-vector file.
type cellCase struct {
	Description string
	ColumnType  string
	CellType    string
	Cell        model.RuleCell
	TypedJSON   string
}

// allCellCases covers every concrete cell type of the governance grammar.
// The *Any cell types are proto zero values: their wire form is the empty
// cell (wire_base64 "") and decode picks the *Any type from the column.
var allCellCases = []cellCase{
	// RuleFiatAmount
	{"fiat amount any", "RuleFiatAmount", "RuleFiatAmountAny", model.FiatAmountAny{}, `{}`},
	{"fiat amount is zero", "RuleFiatAmount", "RuleFiatAmountIsZero", model.FiatAmountIsZero{}, `{}`},
	{"fiat amount range", "RuleFiatAmount", "RuleFiatAmountRange",
		model.FiatAmountRange{MinAmount: "1000", MaxAmount: "50000"},
		`{"minAmount":"1000","maxAmount":"50000"}`},

	// RuleSource
	{"source any", "RuleSource", "RuleSourceAny", model.SourceAny{}, `{}`},
	{"source internal wallet", "RuleSource", "RuleSourceInternalWallet",
		model.SourceInternalWallet{Path: "m/44'/60'/0'"}, `{"path":"m/44'/60'/0'"}`},
	{"source internal address", "RuleSource", "RuleSourceInternalAddress",
		model.SourceInternalAddress{Address: "0xabc", Path: "m/44'/60'/0'/0/0"},
		`{"address":"0xabc","path":"m/44'/60'/0'/0/0"}`},
	{"source any exchange", "RuleSource", "RuleSourceAnyExchange", model.SourceAnyExchange{}, `{}`},
	{"source exchange", "RuleSource", "RuleSourceExchange",
		model.SourceExchange{Label: "kraken-main"}, `{"label":"kraken-main"}`},
	{"source external address", "RuleSource", "RuleSourceExternalAddress",
		model.SourceExternalAddress{Address: "0xdef", Memo: "memo-1"},
		`{"address":"0xdef","memo":"memo-1"}`},

	// RuleDestination
	{"destination any", "RuleDestination", "RuleDestinationAny", model.DestinationAny{}, `{}`},
	{"destination internal wallet", "RuleDestination", "RuleDestinationInternalWallet",
		model.DestinationInternalWallet{Path: "m/44'/60'/1'"}, `{"path":"m/44'/60'/1'"}`},
	{"destination internal address", "RuleDestination", "RuleDestinationInternalAddress",
		model.DestinationInternalAddress{Address: "0x111", Path: "m/44'/60'/1'/0/0"},
		`{"address":"0x111","path":"m/44'/60'/1'/0/0"}`},
	{"destination external address", "RuleDestination", "RuleDestinationExternalAddress",
		model.DestinationExternalAddress{Address: "0x222", Memo: "dest-memo"},
		`{"address":"0x222","memo":"dest-memo"}`},
	{"destination any exchange", "RuleDestination", "RuleDestinationAnyExchange", model.DestinationAnyExchange{}, `{}`},
	{"destination exchange", "RuleDestination", "RuleDestinationExchange",
		model.DestinationExchange{Label: "binance-desk", Memo: "x"},
		`{"label":"binance-desk","memo":"x"}`},
	{"destination contract address", "RuleDestination", "RuleDestinationContractAddress",
		model.DestinationContractAddress{Address: "0x333", Name: "USDC", Symbol: "USDC", Blockchain: "ETH"},
		`{"address":"0x333","name":"USDC","symbol":"USDC","blockchain":"ETH"}`},
	{"destination contract address (unknown blockchain)", "RuleDestination", "RuleDestinationContractAddress",
		model.DestinationContractAddress{Address: "0x1", Blockchain: "4242"},
		`{"address":"0x1","name":"","symbol":"","blockchain":"4242"}`},
	{"destination any external address", "RuleDestination", "RuleDestinationAnyExternalAddress",
		model.DestinationAnyExternalAddress{}, `{}`},
	{"destination any contract address", "RuleDestination", "RuleDestinationAnyContractAddress",
		model.DestinationAnyContractAddress{}, `{}`},

	// RuleStringEqual
	{"string equal any", "RuleStringEqual", "RuleStringEqualAny", model.StringEqualAny{}, `{}`},
	{"string equal empty", "RuleStringEqual", "RuleStringEqualEmpty", model.StringEqualEmpty{}, `{}`},
	{"string equal value", "RuleStringEqual", "RuleStringEqualValue",
		model.StringEqualValue{Value: "contract-id-42"}, `{"value":"contract-id-42"}`},

	// RuleBytesEqual
	{"bytes equal any", "RuleBytesEqual", "RuleBytesEqualAny", model.BytesEqualAny{}, `{}`},
	{"bytes equal empty", "RuleBytesEqual", "RuleBytesEqualEmpty", model.BytesEqualEmpty{}, `{}`},
	{"bytes equal value", "RuleBytesEqual", "RuleBytesEqualValue",
		model.BytesEqualValue{Value: []byte{0xde, 0xad, 0xbe, 0xef}},
		`{"value_base64":"3q2+7w=="}`},

	// RuleStringArrayEqual
	{"string array equal any", "RuleStringArrayEqual", "RuleStringArrayEqualAny", model.StringArrayEqualAny{}, `{}`},
	{"string array equal empty", "RuleStringArrayEqual", "RuleStringArrayEqualEmpty", model.StringArrayEqualEmpty{}, `{}`},
	{"string array equal value", "RuleStringArrayEqual", "RuleStringArrayEqualValue",
		model.StringArrayEqualValue{Values: []string{"a", "b", "c"}},
		`{"values":["a","b","c"]}`},

	// RuleIntegerGreater (sign selects the wire arm)
	{"integer greater any", "RuleIntegerGreater", "RuleIntegerGreaterAny", model.IntegerGreaterAny{}, `{}`},
	{"integer greater positive value", "RuleIntegerGreater", "RuleIntegerGreaterValue",
		model.IntegerGreaterValue{Value: big.NewInt(50)}, `{"value":"50"}`},
	{"integer greater negative value", "RuleIntegerGreater", "RuleIntegerGreaterNegValue",
		model.IntegerGreaterValue{Value: big.NewInt(-50)}, `{"value":"-50"}`},
	{"integer greater zero", "RuleIntegerGreater", "RuleIntegerGreaterValue",
		model.IntegerGreaterValue{Value: new(big.Int)}, `{"value":"0"}`},

	// RuleUIntegerGreater
	{"uinteger greater any", "RuleUIntegerGreater", "RuleUIntegerGreaterAny", model.UIntegerGreaterAny{}, `{}`},
	{"uinteger greater is zero", "RuleUIntegerGreater", "RuleUIntegerGreaterIsZero", model.UIntegerGreaterIsZero{}, `{}`},
	{"uinteger greater value", "RuleUIntegerGreater", "RuleUIntegerGreaterValue",
		model.UIntegerGreaterValue{Value: new(big.Int).SetUint64(18446744073709551615)},
		`{"value":"18446744073709551615"}`},
	{"uinteger greater is equal", "RuleUIntegerGreater", "RuleUIntegerGreaterIsEqual",
		model.UIntegerGreaterIsEqual{Value: big.NewInt(10)}, `{"value":"10"}`},

	// RuleWhitelistedContract
	{"whitelisted contract any", "RuleWhitelistedContract", "RuleWhitelistedContractAny",
		model.WhitelistedContractAny{}, `{}`},
	{"whitelisted contract address", "RuleWhitelistedContract", "RuleWhitelistedContract_RuleDestinationContractAddress",
		model.WhitelistedContractAddress{Address: "0x444", Name: "DAI", Symbol: "DAI", Blockchain: "ETH"},
		`{"address":"0x444","name":"DAI","symbol":"DAI","blockchain":"ETH"}`},
}
