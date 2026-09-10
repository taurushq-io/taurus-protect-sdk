package mapper

import (
	"testing"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// This file complements rules_container_roundtrip_test.go with per-node
// unknown-field preservation, enum passthrough for every encoded enum, the
// RuleSource variant codec, the Xtz details arm, and the encode error paths.

// TestRulesContainer_UnknownFieldsSurviveAllNodes extends the top-level /
// TransactionRules / User coverage to every remaining node type, so a newer
// schema's fields are preserved wherever they appear.
func TestRulesContainer_UnknownFieldsSurviveAllNodes(t *testing.T) {
	uf := func(b byte) protoreflect.RawFields { return protoreflect.RawFields{0xA0, 0x1F, b} } // field 500, varint b

	c := buildRichPbContainer(t)
	c.Groups[0].ProtoReflect().SetUnknown(uf(0x01))
	c.TransactionRules[0].Columns[0].ProtoReflect().SetUnknown(uf(0x02))
	c.TransactionRules[0].Lines[0].ProtoReflect().SetUnknown(uf(0x03))
	c.TransactionRules[0].Details.ProtoReflect().SetUnknown(uf(0x04))
	c.AddressWhitelistingRules[0].ProtoReflect().SetUnknown(uf(0x05))
	c.AddressWhitelistingRules[0].Lines[0].ProtoReflect().SetUnknown(uf(0x06))
	c.ContractAddressWhitelistingRules[0].ProtoReflect().SetUnknown(uf(0x07))
	c.TransactionRules[0].Lines[0].ParallelThresholds[0].ProtoReflect().SetUnknown(uf(0x08))
	c.TransactionRules[0].Lines[0].ParallelThresholds[0].Thresholds[0].ProtoReflect().SetUnknown(uf(0x09))

	decoded, err := RulesContainerFromBytes(mustMarshal(t, c))
	mustNoErr(t, err)
	if !decoded.HasUnknownFields() {
		t.Fatal("expected HasUnknownFields to report true")
	}

	var back pb.RulesContainer
	mustNoErr(t, proto.Unmarshal(mustMarshalContainer(t, decoded), &back))

	check := func(name string, got protoreflect.RawFields, want byte) {
		wantEqual(t, []byte(uf(want)), []byte(got))
		_ = name
	}
	check("group", back.GetGroups()[0].ProtoReflect().GetUnknown(), 0x01)
	check("column", back.GetTransactionRules()[0].GetColumns()[0].ProtoReflect().GetUnknown(), 0x02)
	check("line", back.GetTransactionRules()[0].GetLines()[0].ProtoReflect().GetUnknown(), 0x03)
	check("details", back.GetTransactionRules()[0].GetDetails().ProtoReflect().GetUnknown(), 0x04)
	check("awr", back.GetAddressWhitelistingRules()[0].ProtoReflect().GetUnknown(), 0x05)
	check("awl", back.GetAddressWhitelistingRules()[0].GetLines()[0].ProtoReflect().GetUnknown(), 0x06)
	check("cawr", back.GetContractAddressWhitelistingRules()[0].ProtoReflect().GetUnknown(), 0x07)
	check("thresholds", back.GetTransactionRules()[0].GetLines()[0].GetParallelThresholds()[0].ProtoReflect().GetUnknown(), 0x08)
	check("groupThreshold", back.GetTransactionRules()[0].GetLines()[0].GetParallelThresholds()[0].GetThresholds()[0].ProtoReflect().GetUnknown(), 0x09)
}

// TestRulesContainer_AllEnumsPassThroughNumerically covers numeric passthrough
// for every enum the encoder touches (previously only ColumnType).
func TestRulesContainer_AllEnumsPassThroughNumerically(t *testing.T) {
	c := buildRichPbContainer(t)
	c.Users[0].Roles = append(c.Users[0].Roles, pb.Role(991))
	c.TransactionRules[0].Details.Domain = pb.RulesContainer_TransactionRules_TransactionRuleDetails_RuleDomain(992)
	c.TransactionRules[0].Details.SubDomain = pb.RulesContainer_TransactionRules_TransactionRuleDetails_RuleSubDomain(993)
	c.ContractAddressWhitelistingRules[0].Blockchain = pb.Blockchain(994)

	decoded, err := RulesContainerFromBytes(mustMarshal(t, c))
	mustNoErr(t, err)
	roles := decoded.Users[0].Roles
	wantEqual(t, "991", roles[len(roles)-1])
	wantEqual(t, "992", decoded.TransactionRules[0].Details.Domain)
	wantEqual(t, "993", decoded.TransactionRules[0].Details.SubDomain)
	wantEqual(t, "994", decoded.ContractAddressWhitelistingRules[0].Blockchain)

	var back pb.RulesContainer
	mustNoErr(t, proto.Unmarshal(mustMarshalContainer(t, decoded), &back))
	bRoles := back.GetUsers()[0].GetRoles()
	wantEqual(t, pb.Role(991), bRoles[len(bRoles)-1])
	wantEqual(t, pb.RulesContainer_TransactionRules_TransactionRuleDetails_RuleDomain(992), back.GetTransactionRules()[0].GetDetails().GetDomain())
	wantEqual(t, pb.RulesContainer_TransactionRules_TransactionRuleDetails_RuleSubDomain(993), back.GetTransactionRules()[0].GetDetails().GetSubDomain())
	wantEqual(t, pb.Blockchain(994), back.GetContractAddressWhitelistingRules()[0].GetBlockchain())
}

func TestEnumNumber(t *testing.T) {
	n, err := enumNumber(pb.Role_value, "", "role")
	mustNoErr(t, err)
	wantEqual(t, int32(0), n) // empty -> zero value

	n, err = enumNumber(pb.Role_value, "991", "role")
	mustNoErr(t, err)
	wantEqual(t, int32(991), n) // decimal passthrough

	_, err = enumNumber(pb.Role_value, "NotARole", "role")
	wantErrContaining(t, err, `unknown role "NotARole"`)
}

func TestUint32FromInt(t *testing.T) {
	v, err := uint32FromInt(5, "x")
	mustNoErr(t, err)
	wantEqual(t, uint32(5), v)

	_, err = uint32FromInt(-1, "x")
	wantErrContaining(t, err, "out of range")
}

// TestRulesContainer_AuthoredUnknownEnumsError pins fail-loud for caller typos
// in the domain/sub-domain/role enums (ColumnType is covered elsewhere).
func TestRulesContainer_AuthoredUnknownEnumsError(t *testing.T) {
	cases := []struct {
		name      string
		container *model.DecodedRulesContainer
		substr    string
	}{
		{"role", &model.DecodedRulesContainer{Users: []*model.RuleUser{{ID: "u", Roles: []string{"NotARole"}}}}, "unknown user role"},
		{"domain", &model.DecodedRulesContainer{TransactionRules: []*model.TransactionRules{{Key: "k", Details: &model.TransactionRuleDetails{Domain: "NopeDomain"}}}}, "unknown rule domain"},
		{"subdomain", &model.DecodedRulesContainer{TransactionRules: []*model.TransactionRules{{Key: "k", Details: &model.TransactionRuleDetails{SubDomain: "NopeSub"}}}}, "unknown rule sub-domain"},
		{"blockchain", &model.DecodedRulesContainer{ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{{Blockchain: "NopeChain"}}}, "unknown blockchain"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := RulesContainerToBytes(tc.container)
			wantErrContaining(t, err, tc.substr)
		})
	}
}

// TestRuleSource_VariantCodecRoundTrip round-trips the source variants beyond
// InternalWallet/Any (previously untested).
func TestRuleSource_VariantCodecRoundTrip(t *testing.T) {
	sources := []*model.RuleSource{
		{Type: model.RuleSourceTypeInternalAddress, InternalAddress: &model.RuleSourceInternalAddress{Address: "0xabc", Path: "m/44'/60'/0'/0/0"}},
		{Type: model.RuleSourceTypeExchange, Exchange: &model.RuleSourceExchange{Label: "kraken"}},
		{Type: model.RuleSourceTypeExternalAddress, ExternalAddress: &model.RuleSourceExternalAddress{Address: "0xdef", Memo: "m"}},
		{Type: model.RuleSourceTypeAnyExchange},
	}
	c := &model.DecodedRulesContainer{
		AddressWhitelistingRules: []*model.AddressWhitelistingRules{{
			Currency: "ETH",
			Lines:    []*model.AddressWhitelistingLine{{Cells: sources}},
		}},
	}
	encoded, err := RulesContainerToBytes(c)
	mustNoErr(t, err)
	decoded, err := RulesContainerFromBytes(encoded)
	mustNoErr(t, err)
	wantEqual(t, sources, decoded.AddressWhitelistingRules[0].Lines[0].Cells)
}

func TestRuleSource_UnknownTypeWithoutRawErrors(t *testing.T) {
	c := &model.DecodedRulesContainer{
		AddressWhitelistingRules: []*model.AddressWhitelistingRules{{
			Currency: "ETH",
			Lines:    []*model.AddressWhitelistingLine{{Cells: []*model.RuleSource{{Type: model.RuleSourceType(99)}}}},
		}},
	}
	_, err := RulesContainerToBytes(c)
	wantErrContaining(t, err, "unknown rule source type")
}

func TestRuleSource_MalformedPayloadSurvivesAsRaw(t *testing.T) {
	// A source claiming InternalWallet but carrying an undecodable payload must
	// be preserved verbatim rather than dropped.
	badSource := mustMarshal(t, &pb.RuleSource{
		Type:    pb.RuleSource_RuleSourceInternalWallet,
		Payload: []byte{0xff, 0xff, 0xff, 0xff},
	})
	c := buildRichPbContainer(t)
	c.AddressWhitelistingRules[0].Lines[0].Cells = [][]byte{badSource}

	decoded, err := RulesContainerFromBytes(mustMarshal(t, c))
	mustNoErr(t, err)
	wantEqual(t, badSource, decoded.AddressWhitelistingRules[0].Lines[0].Cells[0].Raw)

	var back pb.RulesContainer
	mustNoErr(t, proto.Unmarshal(mustMarshalContainer(t, decoded), &back))
	wantEqual(t, badSource, back.GetAddressWhitelistingRules()[0].GetLines()[0].GetCells()[0])
}

// TestRulesContainer_XtzDetailsRoundTrip covers the Xtz arm of the details
// encoder (the rich fixture exercises Evm/Cash/Cosmos).
func TestRulesContainer_XtzDetailsRoundTrip(t *testing.T) {
	c := &model.DecodedRulesContainer{
		TransactionRules: []*model.TransactionRules{{
			Key: "XTZ/call",
			Details: &model.TransactionRuleDetails{
				Domain:          "RuleDomainCallContract",
				XtzCallContract: &model.XtzCallContract{ContractType: "FA2", MethodSignature: "transfer"},
			},
		}},
	}
	encoded, err := RulesContainerToBytes(c)
	mustNoErr(t, err)
	decoded, err := RulesContainerFromBytes(encoded)
	mustNoErr(t, err)
	wantEqual(t, &model.XtzCallContract{ContractType: "FA2", MethodSignature: "transfer"}, decoded.TransactionRules[0].Details.XtzCallContract)
}

// TestRulesContainer_EncodeNilElementErrors covers the defensive nil-element
// guards in every encoder loop.
func TestRulesContainer_EncodeNilElementErrors(t *testing.T) {
	line := &model.RuleLine{ParallelThresholds: []*model.SequentialThresholds{nil}}
	cases := map[string]*model.DecodedRulesContainer{
		"nil user":             {Users: []*model.RuleUser{nil}},
		"nil transaction rule": {TransactionRules: []*model.TransactionRules{nil}},
		"nil column":           {TransactionRules: []*model.TransactionRules{{Key: "k", Columns: []*model.RuleColumn{nil}}}},
		"nil line":             {TransactionRules: []*model.TransactionRules{{Key: "k", Lines: []*model.RuleLine{nil}}}},
		"nil awr":              {AddressWhitelistingRules: []*model.AddressWhitelistingRules{nil}},
		"nil awr line":         {AddressWhitelistingRules: []*model.AddressWhitelistingRules{{Currency: "ETH", Lines: []*model.AddressWhitelistingLine{nil}}}},
		"nil cawr":             {ContractAddressWhitelistingRules: []*model.ContractAddressWhitelistingRules{nil}},
		"nil sequential thr":   {TransactionRules: []*model.TransactionRules{{Key: "k", Lines: []*model.RuleLine{line}}}},
		"minSigs out of range": {MinimumDistinctUserSignatures: -1},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if _, err := RulesContainerToBytes(c); err == nil {
				t.Fatalf("%s: expected encode error", name)
			}
		})
	}
}

func TestRulesContainer_EncodeNilGroupThresholdErrors(t *testing.T) {
	c := &model.DecodedRulesContainer{
		TransactionRules: []*model.TransactionRules{{
			Key:   "k",
			Lines: []*model.RuleLine{{ParallelThresholds: []*model.SequentialThresholds{{Thresholds: []*model.GroupThreshold{nil}}}}},
		}},
	}
	if _, err := RulesContainerToBytes(c); err == nil {
		t.Fatal("expected error for nil group threshold")
	}
}

func mustMarshalContainer(t *testing.T, c *model.DecodedRulesContainer) []byte {
	t.Helper()
	data, err := RulesContainerToBytes(c)
	mustNoErr(t, err)
	return data
}
