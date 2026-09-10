package model

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"strings"
	"testing"
)

func mustKey(t *testing.T) *ecdsa.PublicKey {
	t.Helper()
	k, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	return &k.PublicKey
}

func TestFindUserByID(t *testing.T) {
	c := &DecodedRulesContainer{Users: []*RuleUser{{ID: "u1"}, {ID: "u2"}}}
	if got := c.FindUserByID("u2"); got == nil || got.ID != "u2" {
		t.Fatalf("expected u2, got %#v", got)
	}
	if got := c.FindUserByID("nope"); got != nil {
		t.Fatalf("expected nil for missing id, got %#v", got)
	}
	if got := c.FindUserByID(""); got != nil {
		t.Fatalf("expected nil for empty id, got %#v", got)
	}
	if got := (&DecodedRulesContainer{}).FindUserByID("u1"); got != nil {
		t.Fatalf("expected nil for nil users, got %#v", got)
	}
}

func TestFindGroupByID(t *testing.T) {
	c := &DecodedRulesContainer{Groups: []*RuleGroup{{ID: "g1"}}}
	if got := c.FindGroupByID("g1"); got == nil || got.ID != "g1" {
		t.Fatalf("expected g1, got %#v", got)
	}
	if got := c.FindGroupByID("nope"); got != nil {
		t.Fatalf("expected nil, got %#v", got)
	}
	if got := c.FindGroupByID(""); got != nil {
		t.Fatalf("expected nil for empty id, got %#v", got)
	}
}

func TestRuleUserHasRoleAndGroupContainsUser(t *testing.T) {
	u := &RuleUser{Roles: []string{"SUPERADMIN", "HSMSLOT"}}
	if !u.HasRole("HSMSLOT") {
		t.Fatal("expected HasRole(HSMSLOT) true")
	}
	if u.HasRole("NOPE") {
		t.Fatal("expected HasRole(NOPE) false")
	}
	g := &RuleGroup{UserIDs: []string{"u1", "u2"}}
	if !g.ContainsUser("u2") {
		t.Fatal("expected ContainsUser(u2) true")
	}
	if g.ContainsUser("u3") {
		t.Fatal("expected ContainsUser(u3) false")
	}
}

func TestFindAddressWhitelistingRules_Tiers(t *testing.T) {
	exact := &AddressWhitelistingRules{Currency: "ETH", Network: "mainnet"}
	blockchainOnly := &AddressWhitelistingRules{Currency: "ETH", Network: ""} // wildcard network
	global := &AddressWhitelistingRules{Currency: "Any", Network: "mainnet"}  // wildcard currency
	c := &DecodedRulesContainer{AddressWhitelistingRules: []*AddressWhitelistingRules{exact, blockchainOnly, global}}

	if got := c.FindAddressWhitelistingRules("ETH", "mainnet"); got != exact {
		t.Fatal("expected exact match")
	}
	if got := c.FindAddressWhitelistingRules("ETH", "testnet"); got != blockchainOnly {
		t.Fatal("expected blockchain-only (wildcard network) match")
	}
	if got := c.FindAddressWhitelistingRules("BTC", "mainnet"); got != global {
		t.Fatal("expected global-default (wildcard currency) match")
	}

	noGlobal := &DecodedRulesContainer{AddressWhitelistingRules: []*AddressWhitelistingRules{exact}}
	if got := noGlobal.FindAddressWhitelistingRules("SOL", "x"); got != nil {
		t.Fatalf("expected nil when nothing matches, got %#v", got)
	}
	if got := (&DecodedRulesContainer{}).FindAddressWhitelistingRules("ETH", "mainnet"); got != nil {
		t.Fatal("expected nil for nil slice")
	}
}

func TestFindContractAddressWhitelistingRules_Tiers(t *testing.T) {
	exact := &ContractAddressWhitelistingRules{Blockchain: "ETH", Network: "mainnet"}
	blockchainOnly := &ContractAddressWhitelistingRules{Blockchain: "ETH", Network: "Any"}
	global := &ContractAddressWhitelistingRules{Blockchain: "", Network: "mainnet"}
	c := &DecodedRulesContainer{ContractAddressWhitelistingRules: []*ContractAddressWhitelistingRules{exact, blockchainOnly, global}}

	if got := c.FindContractAddressWhitelistingRules("ETH", "mainnet"); got != exact {
		t.Fatal("expected exact match")
	}
	if got := c.FindContractAddressWhitelistingRules("ETH", "testnet"); got != blockchainOnly {
		t.Fatal("expected blockchain-only match")
	}
	if got := c.FindContractAddressWhitelistingRules("BTC", "mainnet"); got != global {
		t.Fatal("expected global-default match")
	}
	if got := (&DecodedRulesContainer{}).FindContractAddressWhitelistingRules("ETH", "mainnet"); got != nil {
		t.Fatal("expected nil for nil slice")
	}
}

func TestGetHsmPublicKey(t *testing.T) {
	key := mustKey(t)
	c := &DecodedRulesContainer{Users: []*RuleUser{
		{ID: "auth", Roles: []string{"SUPERADMIN"}},
		{ID: "hsm", Roles: []string{"HSMSLOT"}, PublicKey: key},
	}}
	if got := c.GetHsmPublicKey(); got != key {
		t.Fatal("expected the HSMSLOT user's key")
	}

	// sync.Once caches the first result even if Users later changes.
	c.Users = []*RuleUser{{ID: "hsm2", Roles: []string{"HSMSLOT"}, PublicKey: mustKey(t)}}
	if got := c.GetHsmPublicKey(); got != key {
		t.Fatal("expected the cached key after mutation")
	}

	// No HSMSLOT user -> nil.
	noHsm := &DecodedRulesContainer{Users: []*RuleUser{{ID: "u", Roles: []string{"SUPERADMIN"}}}}
	if got := noHsm.GetHsmPublicKey(); got != nil {
		t.Fatal("expected nil when no HSMSLOT user")
	}

	// HSMSLOT user with a nil parsed key -> nil.
	nilKey := &DecodedRulesContainer{Users: []*RuleUser{{ID: "hsm", Roles: []string{"HSMSLOT"}}}}
	if got := nilKey.GetHsmPublicKey(); got != nil {
		t.Fatal("expected nil when HSMSLOT user has no parsed key")
	}
}

// TestHasUnknownFields_PerNode isolates an unknown field (or raw cell/source) at
// each node in turn and asserts HasUnknownFields still reports true, so the
// upgrade-before-editing hint fires wherever a newer schema's data lands. The
// roundtrip test proves such data survives re-encode; this proves it is detected
// per node (HasUnknownFields short-circuits, so a single all-nodes fixture cannot
// exercise every detection branch). A clean container reports false.
func TestHasUnknownFields_PerNode(t *testing.T) {
	b := []byte{0x01}
	cases := map[string]*DecodedRulesContainer{
		"container":            {UnknownFields: b},
		"user":                 {Users: []*RuleUser{{UnknownFields: b}}},
		"group":                {Groups: []*RuleGroup{{UnknownFields: b}}},
		"transaction rules":    {TransactionRules: []*TransactionRules{{UnknownFields: b}}},
		"column":               {TransactionRules: []*TransactionRules{{Columns: []*RuleColumn{{UnknownFields: b}}}}},
		"line":                 {TransactionRules: []*TransactionRules{{Lines: []*RuleLine{{UnknownFields: b}}}}},
		"line thresholds":      {TransactionRules: []*TransactionRules{{Lines: []*RuleLine{{ParallelThresholds: []*SequentialThresholds{{UnknownFields: b}}}}}}},
		"line group threshold": {TransactionRules: []*TransactionRules{{Lines: []*RuleLine{{ParallelThresholds: []*SequentialThresholds{{Thresholds: []*GroupThreshold{{UnknownFields: b}}}}}}}}},
		"line raw cell":        {TransactionRules: []*TransactionRules{{Lines: []*RuleLine{{Cells: []RuleCell{RawCell{}}}}}}},
		"details":              {TransactionRules: []*TransactionRules{{Details: &TransactionRuleDetails{UnknownFields: b}}}},
		"awr":                  {AddressWhitelistingRules: []*AddressWhitelistingRules{{UnknownFields: b}}},
		"awr thresholds":       {AddressWhitelistingRules: []*AddressWhitelistingRules{{ParallelThresholds: []*SequentialThresholds{{UnknownFields: b}}}}},
		"awr line":             {AddressWhitelistingRules: []*AddressWhitelistingRules{{Lines: []*AddressWhitelistingLine{{UnknownFields: b}}}}},
		"awr line thresholds":  {AddressWhitelistingRules: []*AddressWhitelistingRules{{Lines: []*AddressWhitelistingLine{{ParallelThresholds: []*SequentialThresholds{{UnknownFields: b}}}}}}},
		"awr line raw source":  {AddressWhitelistingRules: []*AddressWhitelistingRules{{Lines: []*AddressWhitelistingLine{{Cells: []*RuleSource{{Raw: b}}}}}}},
		"cawr":                 {ContractAddressWhitelistingRules: []*ContractAddressWhitelistingRules{{UnknownFields: b}}},
		"cawr thresholds":      {ContractAddressWhitelistingRules: []*ContractAddressWhitelistingRules{{ParallelThresholds: []*SequentialThresholds{{UnknownFields: b}}}}},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if !c.HasUnknownFields() {
				t.Fatalf("%s: expected HasUnknownFields true", name)
			}
		})
	}

	clean := &DecodedRulesContainer{
		Users:                            []*RuleUser{{ID: "u"}},
		TransactionRules:                 []*TransactionRules{{Key: "k", Lines: []*RuleLine{{Cells: []RuleCell{}}}, Details: &TransactionRuleDetails{}}},
		AddressWhitelistingRules:         []*AddressWhitelistingRules{{Currency: "ETH", Lines: []*AddressWhitelistingLine{{Cells: []*RuleSource{{}}}}}},
		ContractAddressWhitelistingRules: []*ContractAddressWhitelistingRules{{Blockchain: "ETH"}},
	}
	if clean.HasUnknownFields() {
		t.Fatal("expected HasUnknownFields false for a clean container")
	}
}

// Every RuleCell variant must report a Kind() matching its Go type name, and the set
// must match the discriminators the other three SDKs expose (Java kind(), Python .kind,
// TS .kind). Go had no discriminator at all, forcing callers into a type switch and
// leaving nothing to compare against the peers.
func TestRuleCell_KindMatchesTheTypeNameForEveryVariant(t *testing.T) {
	cells := []RuleCell{
		RawCell{}, FiatAmountAny{}, FiatAmountIsZero{}, FiatAmountRange{},
		SourceAny{}, SourceInternalWallet{}, SourceInternalAddress{},
		SourceAnyExchange{}, SourceExchange{}, SourceExternalAddress{},
		DestinationAny{}, DestinationInternalWallet{}, DestinationInternalAddress{},
		DestinationExternalAddress{}, DestinationAnyExchange{}, DestinationExchange{},
		DestinationContractAddress{}, DestinationAnyExternalAddress{},
		DestinationAnyContractAddress{},
		StringEqualAny{}, StringEqualEmpty{}, StringEqualValue{},
		BytesEqualAny{}, BytesEqualEmpty{}, BytesEqualValue{},
		StringArrayEqualAny{}, StringArrayEqualEmpty{}, StringArrayEqualValue{},
		IntegerGreaterAny{}, IntegerGreaterValue{},
		UIntegerGreaterAny{}, UIntegerGreaterIsZero{}, UIntegerGreaterValue{},
		UIntegerGreaterIsEqual{},
		WhitelistedContractAny{}, WhitelistedContractAddress{},
	}

	if len(cells) != 36 {
		t.Fatalf("listed %d variants, expected all 36 — a new variant needs a Kind()", len(cells))
	}

	seen := make(map[string]struct{}, len(cells))
	for _, c := range cells {
		typeName := fmt.Sprintf("%T", c)
		typeName = typeName[strings.LastIndex(typeName, ".")+1:]

		if c.Kind() != typeName {
			t.Errorf("Kind() = %q, want %q", c.Kind(), typeName)
		}
		if _, dup := seen[c.Kind()]; dup {
			t.Errorf("duplicate Kind() %q", c.Kind())
		}
		seen[c.Kind()] = struct{}{}
	}
}
