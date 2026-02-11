package mapper

import (
	"encoding/base64"
	"encoding/json"
	"os"
	"strings"
	"testing"

	pb "github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/proto"
	"google.golang.org/protobuf/proto"
)

// testFixture holds the structure of the test fixture file.
type testFixture struct {
	RulesSignatures     string          `json:"rulesSignatures"`
	RulesContainerJSON  json.RawMessage `json:"rulesContainerJson"`
	Metadata            *struct {
		Hash            string `json:"hash"`
		PayloadAsString string `json:"payloadAsString"`
	} `json:"metadata"`
}

// loadTestFixture loads the test fixture from the testdata directory.
func loadTestFixture(t *testing.T) *testFixture {
	t.Helper()
	data, err := os.ReadFile("../testdata/whitelisted_address_raw_response.json")
	if err != nil {
		t.Fatalf("Failed to load test fixture: %v", err)
	}

	var fixture testFixture
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatalf("Failed to parse test fixture: %v", err)
	}
	return &fixture
}

// createTestRulesContainer creates a protobuf RulesContainer for testing.
// This matches the structure in the fixture's rulesContainerJson.
func createTestRulesContainer() *pb.RulesContainer {
	return &pb.RulesContainer{
		Users: []*pb.User{
			{
				Id:        "superadmin1@bank.com",
				PublicKey: "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEyWjh6d+PgOK3LqockShMcDMtAHImitWjoVSX/FzBAWvemeaeNnYDKzEXiDDgiq2tILFL1Chdkqofhp9EdBZOlQ==\n-----END PUBLIC KEY-----",
				Roles:     []pb.Role{pb.Role_SUPERADMIN},
			},
			{
				Id:        "superadmin2@bank.com",
				PublicKey: "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAELJhEUNLLHgI8LiWJaeJGpaBfdvgoYyKsjSFyTMxECR/E+1qpzDlNNug7hDPgBPpZ3Z+U8QWjaKB4Mrbj2/kImQ==\n-----END PUBLIC KEY-----",
				Roles:     []pb.Role{pb.Role_SUPERADMIN},
			},
			{
				Id:        "team1@bank.com",
				PublicKey: "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEM2NtzaFhm7xIR3OvWq5chW3/GEvWL+3uqoE6lEJ13eWbulxsP/5h36VCqYDIGN/0wDeWwLYdpu5HhSXWhxCsCA==\n-----END PUBLIC KEY-----",
				Roles:     []pb.Role{pb.Role_TPUSER, pb.Role_REQUESTAPPROVER},
			},
			{
				Id:        "hsmslot@bank.com",
				PublicKey: "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEY9zGugzNLIfpZuaUrzywEh/8ZdtX4IIuIpDHLvJ36glFjfxxSZdOG6yHKFFlQh1GX3OCFZxHe+xeOGBJHBgraA==\n-----END PUBLIC KEY-----",
				Roles:     []pb.Role{pb.Role_HSMSLOT},
			},
		},
		Groups: []*pb.Group{
			{
				Id:      "team1",
				UserIds: []string{"team1@bank.com"},
			},
			{
				Id:      "superadmins",
				UserIds: []string{"superadmin1@bank.com", "superadmin2@bank.com"},
			},
		},
		MinimumDistinctUserSignatures:  0,
		MinimumDistinctGroupSignatures: 0,
		AddressWhitelistingRules: []*pb.RulesContainer_AddressWhitelistingRules{
			{
				Currency: "ALGO",
				Network:  "mainnet",
				ParallelThresholds: []*pb.SequentialThresholds{
					{
						Thresholds: []*pb.GroupThreshold{
							{
								GroupId:           "team1",
								MinimumSignatures: 1,
							},
						},
					},
				},
			},
		},
		ContractAddressWhitelistingRules: []*pb.RulesContainer_ContractAddressWhitelistingRules{},
		EnforcedRulesHash:                "",
		Timestamp:                        1706194800,
	}
}

// encodeRulesContainerToBase64 encodes a protobuf RulesContainer to base64.
func encodeRulesContainerToBase64(t *testing.T, container *pb.RulesContainer) string {
	t.Helper()
	data, err := proto.Marshal(container)
	if err != nil {
		t.Fatalf("Failed to marshal protobuf: %v", err)
	}
	return base64.StdEncoding.EncodeToString(data)
}

// =============================================================================
// Category A: Rules Container Decoding Tests (9 tests)
// =============================================================================

// TestDecodeRulesContainerFromBase64_Success tests successful decoding of a rules container.
func TestDecodeRulesContainerFromBase64_Success(t *testing.T) {
	container := createTestRulesContainer()
	base64Data := encodeRulesContainerToBase64(t, container)

	decoded, err := RulesContainerFromBase64(base64Data)
	if err != nil {
		t.Fatalf("RulesContainerFromBase64() error = %v", err)
	}

	if decoded == nil {
		t.Fatal("RulesContainerFromBase64() returned nil")
	}

	// Verify basic fields
	if decoded.Timestamp != 1706194800 {
		t.Errorf("Timestamp = %d, want %d", decoded.Timestamp, 1706194800)
	}
	if decoded.MinimumDistinctUserSignatures != 0 {
		t.Errorf("MinimumDistinctUserSignatures = %d, want 0", decoded.MinimumDistinctUserSignatures)
	}
	if decoded.MinimumDistinctGroupSignatures != 0 {
		t.Errorf("MinimumDistinctGroupSignatures = %d, want 0", decoded.MinimumDistinctGroupSignatures)
	}
}

// TestDecodedContainerHasUsers verifies that the decoded container has the expected users.
func TestDecodedContainerHasUsers(t *testing.T) {
	container := createTestRulesContainer()
	base64Data := encodeRulesContainerToBase64(t, container)

	decoded, err := RulesContainerFromBase64(base64Data)
	if err != nil {
		t.Fatalf("RulesContainerFromBase64() error = %v", err)
	}

	// Verify 4 users in decoded container
	if len(decoded.Users) != 4 {
		t.Errorf("Users count = %d, want 4", len(decoded.Users))
	}

	// Verify user IDs
	expectedUserIDs := []string{
		"superadmin1@bank.com",
		"superadmin2@bank.com",
		"team1@bank.com",
		"hsmslot@bank.com",
	}
	for i, user := range decoded.Users {
		if i < len(expectedUserIDs) && user.ID != expectedUserIDs[i] {
			t.Errorf("Users[%d].ID = %q, want %q", i, user.ID, expectedUserIDs[i])
		}
	}
}

// TestDecodedContainerHasGroups verifies that the decoded container has the expected groups.
func TestDecodedContainerHasGroups(t *testing.T) {
	container := createTestRulesContainer()
	base64Data := encodeRulesContainerToBase64(t, container)

	decoded, err := RulesContainerFromBase64(base64Data)
	if err != nil {
		t.Fatalf("RulesContainerFromBase64() error = %v", err)
	}

	// Verify 2 groups in decoded container
	if len(decoded.Groups) != 2 {
		t.Errorf("Groups count = %d, want 2", len(decoded.Groups))
	}

	// Verify group IDs
	expectedGroupIDs := []string{"team1", "superadmins"}
	for i, group := range decoded.Groups {
		if i < len(expectedGroupIDs) && group.ID != expectedGroupIDs[i] {
			t.Errorf("Groups[%d].ID = %q, want %q", i, group.ID, expectedGroupIDs[i])
		}
	}

	// Verify superadmins group has 2 users
	for _, group := range decoded.Groups {
		if group.ID == "superadmins" {
			if len(group.UserIDs) != 2 {
				t.Errorf("superadmins group UserIDs count = %d, want 2", len(group.UserIDs))
			}
			break
		}
	}
}

// TestDecodedUsersHavePemKeys verifies that each user has a public_key_pem field.
func TestDecodedUsersHavePemKeys(t *testing.T) {
	container := createTestRulesContainer()
	base64Data := encodeRulesContainerToBase64(t, container)

	decoded, err := RulesContainerFromBase64(base64Data)
	if err != nil {
		t.Fatalf("RulesContainerFromBase64() error = %v", err)
	}

	pemHeader := "-----BEGIN PUBLIC KEY-----"
	for i, user := range decoded.Users {
		if user.PublicKeyPEM == "" {
			t.Errorf("Users[%d].PublicKeyPEM is empty for user %q", i, user.ID)
			continue
		}
		// Verify PEM format
		if len(user.PublicKeyPEM) < 50 {
			t.Errorf("Users[%d].PublicKeyPEM is too short: %d chars", i, len(user.PublicKeyPEM))
		}
		// Verify PEM header is present (may have leading whitespace/newline)
		if !containsSubstr(user.PublicKeyPEM, pemHeader) {
			t.Errorf("Users[%d].PublicKeyPEM missing PEM header: %q", i, user.PublicKeyPEM[:50])
		}
	}
}

// containsSubstr checks if s contains substr.
func containsSubstr(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TestDecodedUsersHaveRoles verifies that users have correct roles array.
func TestDecodedUsersHaveRoles(t *testing.T) {
	container := createTestRulesContainer()
	base64Data := encodeRulesContainerToBase64(t, container)

	decoded, err := RulesContainerFromBase64(base64Data)
	if err != nil {
		t.Fatalf("RulesContainerFromBase64() error = %v", err)
	}

	// Map of expected roles for each user
	expectedRoles := map[string][]string{
		"superadmin1@bank.com": {"SUPERADMIN"},
		"superadmin2@bank.com": {"SUPERADMIN"},
		"team1@bank.com":       {"TPUSER", "REQUESTAPPROVER"},
		"hsmslot@bank.com":     {"HSMSLOT"},
	}

	for _, user := range decoded.Users {
		expected, ok := expectedRoles[user.ID]
		if !ok {
			continue
		}
		if len(user.Roles) != len(expected) {
			t.Errorf("User %q has %d roles, want %d", user.ID, len(user.Roles), len(expected))
			continue
		}
		for i, role := range expected {
			if user.Roles[i] != role {
				t.Errorf("User %q role[%d] = %q, want %q", user.ID, i, user.Roles[i], role)
			}
		}
	}
}

// TestFindSuperadminUsers finds users with SUPERADMIN role (should find 2).
func TestFindSuperadminUsers(t *testing.T) {
	container := createTestRulesContainer()
	base64Data := encodeRulesContainerToBase64(t, container)

	decoded, err := RulesContainerFromBase64(base64Data)
	if err != nil {
		t.Fatalf("RulesContainerFromBase64() error = %v", err)
	}

	superadminCount := 0
	for _, user := range decoded.Users {
		if user.HasRole("SUPERADMIN") {
			superadminCount++
		}
	}

	if superadminCount != 2 {
		t.Errorf("SUPERADMIN users count = %d, want 2", superadminCount)
	}
}

// TestFindHsmslotUser finds user with HSMSLOT role (should find 1).
func TestFindHsmslotUser(t *testing.T) {
	container := createTestRulesContainer()
	base64Data := encodeRulesContainerToBase64(t, container)

	decoded, err := RulesContainerFromBase64(base64Data)
	if err != nil {
		t.Fatalf("RulesContainerFromBase64() error = %v", err)
	}

	hsmslotCount := 0
	var hsmslotUser string
	for _, user := range decoded.Users {
		if user.HasRole("HSMSLOT") {
			hsmslotCount++
			hsmslotUser = user.ID
		}
	}

	if hsmslotCount != 1 {
		t.Errorf("HSMSLOT users count = %d, want 1", hsmslotCount)
	}
	if hsmslotUser != "hsmslot@bank.com" {
		t.Errorf("HSMSLOT user ID = %q, want %q", hsmslotUser, "hsmslot@bank.com")
	}

	// Also verify GetHsmPublicKey works
	hsmKey := decoded.GetHsmPublicKey()
	if hsmKey == nil {
		t.Error("GetHsmPublicKey() returned nil, want non-nil")
	}
}

// TestDecodedHasAddressWhitelistingRules verifies 1 ALGO/mainnet rule.
func TestDecodedHasAddressWhitelistingRules(t *testing.T) {
	container := createTestRulesContainer()
	base64Data := encodeRulesContainerToBase64(t, container)

	decoded, err := RulesContainerFromBase64(base64Data)
	if err != nil {
		t.Fatalf("RulesContainerFromBase64() error = %v", err)
	}

	// Verify 1 address whitelisting rule
	if len(decoded.AddressWhitelistingRules) != 1 {
		t.Errorf("AddressWhitelistingRules count = %d, want 1", len(decoded.AddressWhitelistingRules))
	}

	if len(decoded.AddressWhitelistingRules) > 0 {
		rule := decoded.AddressWhitelistingRules[0]

		// Verify ALGO/mainnet
		if rule.Currency != "ALGO" {
			t.Errorf("AddressWhitelistingRules[0].Currency = %q, want %q", rule.Currency, "ALGO")
		}
		if rule.Network != "mainnet" {
			t.Errorf("AddressWhitelistingRules[0].Network = %q, want %q", rule.Network, "mainnet")
		}

		// Verify parallel thresholds
		if len(rule.ParallelThresholds) != 1 {
			t.Errorf("ParallelThresholds count = %d, want 1", len(rule.ParallelThresholds))
		} else {
			pt := rule.ParallelThresholds[0]
			if len(pt.Thresholds) != 1 {
				t.Errorf("Thresholds count = %d, want 1", len(pt.Thresholds))
			} else {
				if pt.Thresholds[0].GroupID != "team1" {
					t.Errorf("GroupID = %q, want %q", pt.Thresholds[0].GroupID, "team1")
				}
				if pt.Thresholds[0].MinimumSignatures != 1 {
					t.Errorf("MinimumSignatures = %d, want 1", pt.Thresholds[0].MinimumSignatures)
				}
			}
		}
	}

	// Verify FindAddressWhitelistingRules method
	found := decoded.FindAddressWhitelistingRules("ALGO", "mainnet")
	if found == nil {
		t.Error("FindAddressWhitelistingRules(ALGO, mainnet) returned nil")
	} else if found.Currency != "ALGO" || found.Network != "mainnet" {
		t.Errorf("FindAddressWhitelistingRules returned wrong rule: %s/%s", found.Currency, found.Network)
	}
}

// TestInvalidBase64RaisesError tests that invalid base64 input raises an error.
func TestInvalidBase64RaisesError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid base64 characters",
			input: "not-valid-base64!!!",
		},
		{
			name:  "truncated base64",
			input: "YWJj", // truncated/incomplete
		},
		{
			name:  "valid base64 but invalid protobuf",
			input: base64.StdEncoding.EncodeToString([]byte("not a protobuf message")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := RulesContainerFromBase64(tt.input)
			if err == nil {
				t.Error("RulesContainerFromBase64() expected error, got nil")
			}
		})
	}
}

// =============================================================================
// Category B: User Signatures Decoding Tests (5 tests)
// =============================================================================

// TestDecodeUserSignaturesFromBase64_Success tests successful decoding using rulesSignatures from fixture.
func TestDecodeUserSignaturesFromBase64_Success(t *testing.T) {
	fixture := loadTestFixture(t)

	signatures, err := UserSignaturesFromBase64(fixture.RulesSignatures)
	if err != nil {
		t.Fatalf("UserSignaturesFromBase64() error = %v", err)
	}

	if signatures == nil {
		t.Fatal("UserSignaturesFromBase64() returned nil")
	}
}

// TestSignaturesContainUserIds verifies that signatures contain superadmin1@bank.com and superadmin2@bank.com.
func TestSignaturesContainUserIds(t *testing.T) {
	fixture := loadTestFixture(t)

	signatures, err := UserSignaturesFromBase64(fixture.RulesSignatures)
	if err != nil {
		t.Fatalf("UserSignaturesFromBase64() error = %v", err)
	}

	expectedUserIDs := map[string]bool{
		"superadmin1@bank.com": false,
		"superadmin2@bank.com": false,
	}

	for _, sig := range signatures {
		if _, exists := expectedUserIDs[sig.UserID]; exists {
			expectedUserIDs[sig.UserID] = true
		}
	}

	for userID, found := range expectedUserIDs {
		if !found {
			t.Errorf("Expected user ID %q not found in signatures", userID)
		}
	}
}

// TestSignaturesContainSignatureBytes verifies that signature field is not empty.
func TestSignaturesContainSignatureBytes(t *testing.T) {
	fixture := loadTestFixture(t)

	signatures, err := UserSignaturesFromBase64(fixture.RulesSignatures)
	if err != nil {
		t.Fatalf("UserSignaturesFromBase64() error = %v", err)
	}

	for i, sig := range signatures {
		if sig.Signature == "" {
			t.Errorf("signatures[%d].Signature is empty for user %q", i, sig.UserID)
		}
		// Verify signature is base64 encoded (from the mapper)
		_, err := base64.StdEncoding.DecodeString(sig.Signature)
		if err != nil {
			t.Errorf("signatures[%d].Signature is not valid base64: %v", i, err)
		}
	}
}

// TestSignaturesCountMatchesExpected verifies that there are exactly 2 signatures.
func TestSignaturesCountMatchesExpected(t *testing.T) {
	fixture := loadTestFixture(t)

	signatures, err := UserSignaturesFromBase64(fixture.RulesSignatures)
	if err != nil {
		t.Fatalf("UserSignaturesFromBase64() error = %v", err)
	}

	if len(signatures) != 2 {
		t.Errorf("Signatures count = %d, want 2", len(signatures))
	}
}

// TestInvalidSignaturesBase64RaisesError tests that invalid input raises an error.
func TestInvalidSignaturesBase64RaisesError(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "invalid base64 characters",
			input: "not-valid-base64!!!",
		},
		{
			name:  "truncated base64",
			input: "YWJj", // truncated/incomplete
		},
		{
			name:  "valid base64 but invalid protobuf",
			input: base64.StdEncoding.EncodeToString([]byte("not a protobuf message")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := UserSignaturesFromBase64(tt.input)
			if err == nil {
				t.Error("UserSignaturesFromBase64() expected error, got nil")
			}
		})
	}
}
