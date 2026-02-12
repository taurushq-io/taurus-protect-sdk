package helper

import (
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
)

// vectorsFile holds the parsed test vectors from the shared JSON file.
type vectorsFile struct {
	Vectors struct {
		HexHash              []hexHashVector              `json:"hex_hash"`
		HMACSHA256           []hmacVector                 `json:"hmac_sha256"`
		ConstantTimeCompare  []constantTimeCompareVector  `json:"constant_time_compare"`
		LegacyHashAddress    []legacyHashAddressVector    `json:"legacy_hash_address"`
		LegacyHashAsset      []legacyHashAssetVector      `json:"legacy_hash_asset"`
	} `json:"vectors"`
}

type hexHashVector struct {
	Description string `json:"description"`
	Input       string `json:"input"`
	Expected    string `json:"expected"`
}

type hmacVector struct {
	Description    string `json:"description"`
	KeyHex         string `json:"key_hex"`
	Data           string `json:"data"`
	ExpectedBase64 string `json:"expected_base64"`
}

type constantTimeCompareVector struct {
	Description string `json:"description"`
	A           string `json:"a"`
	B           string `json:"b"`
	Expected    bool   `json:"expected"`
}

type legacyHashAddressVector struct {
	Description                  string `json:"description"`
	Payload                      string `json:"payload"`
	OriginalHash                 string `json:"original_hash"`
	ExpectedWithoutContractType  string `json:"expected_without_contract_type"`
	ExpectedWithoutLabels        string `json:"expected_without_labels"`
	ExpectedWithoutBoth          string `json:"expected_without_both"`
	ExpectedLegacyCount          int    `json:"expected_legacy_count"`
}

type legacyHashAssetVector struct {
	Description              string `json:"description"`
	Payload                  string `json:"payload"`
	OriginalHash             string `json:"original_hash"`
	ExpectedWithoutIsNFT     string `json:"expected_without_is_nft"`
	ExpectedWithoutKindType  string `json:"expected_without_kind_type"`
	ExpectedWithoutBoth      string `json:"expected_without_both"`
	ExpectedLegacyCount      int    `json:"expected_legacy_count"`
}

func loadVectors(t *testing.T) *vectorsFile {
	t.Helper()

	// Get directory of this test file
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Failed to get caller info")
	}

	// Navigate to repo root: helper/ -> protect/ -> pkg/ -> sdk-go/ -> repo/
	repoRoot := filepath.Join(filepath.Dir(filename), "..", "..", "..", "..")
	vectorsPath := filepath.Join(repoRoot, "docs", "test-vectors", "crypto-test-vectors.json")

	data, err := os.ReadFile(vectorsPath)
	if err != nil {
		t.Fatalf("Failed to read vectors file %s: %v", vectorsPath, err)
	}

	var vf vectorsFile
	if err := json.Unmarshal(data, &vf); err != nil {
		t.Fatalf("Failed to parse vectors JSON: %v", err)
	}
	return &vf
}

func TestCrossSdkHexHash(t *testing.T) {
	vf := loadVectors(t)

	for _, vec := range vf.Vectors.HexHash {
		t.Run(vec.Description, func(t *testing.T) {
			result := crypto.CalculateHexHash(vec.Input)
			if result != vec.Expected {
				t.Errorf("SHA-256 mismatch: got %s, want %s", result, vec.Expected)
			}
		})
	}
}

func TestCrossSdkHMACSHA256(t *testing.T) {
	vf := loadVectors(t)

	for _, vec := range vf.Vectors.HMACSHA256 {
		t.Run(vec.Description, func(t *testing.T) {
			key, err := hex.DecodeString(vec.KeyHex)
			if err != nil {
				t.Fatalf("Failed to decode key hex: %v", err)
			}
			result := crypto.CalculateBase64HMAC(key, vec.Data)
			if result != vec.ExpectedBase64 {
				t.Errorf("HMAC-SHA256 mismatch: got %s, want %s", result, vec.ExpectedBase64)
			}
		})
	}
}

func TestCrossSdkConstantTimeCompare(t *testing.T) {
	vf := loadVectors(t)

	for _, vec := range vf.Vectors.ConstantTimeCompare {
		t.Run(vec.Description, func(t *testing.T) {
			result := ConstantTimeCompare(vec.A, vec.B)
			if result != vec.Expected {
				t.Errorf("Constant-time compare mismatch: got %v, want %v", result, vec.Expected)
			}
		})
	}
}

func TestCrossSdkLegacyAddressHash(t *testing.T) {
	vf := loadVectors(t)

	for _, vec := range vf.Vectors.LegacyHashAddress {
		t.Run(vec.Description+" - original hash", func(t *testing.T) {
			result := crypto.CalculateHexHash(vec.Payload)
			if result != vec.OriginalHash {
				t.Errorf("Original hash mismatch: got %s, want %s", result, vec.OriginalHash)
			}
		})

		t.Run(vec.Description+" - legacy hashes", func(t *testing.T) {
			legacyHashes := ComputeLegacyHashes(vec.Payload)

			if len(legacyHashes) != vec.ExpectedLegacyCount {
				t.Errorf("Legacy hash count mismatch: got %d, want %d", len(legacyHashes), vec.ExpectedLegacyCount)
			}

			if vec.ExpectedLegacyCount > 0 {
				assertContains(t, legacyHashes, vec.ExpectedWithoutContractType, "without_contract_type")
				assertContains(t, legacyHashes, vec.ExpectedWithoutLabels, "without_labels")
				assertContains(t, legacyHashes, vec.ExpectedWithoutBoth, "without_both")
			}
		})
	}
}

func TestCrossSdkLegacyAssetHash(t *testing.T) {
	vf := loadVectors(t)

	for _, vec := range vf.Vectors.LegacyHashAsset {
		t.Run(vec.Description+" - original hash", func(t *testing.T) {
			result := crypto.CalculateHexHash(vec.Payload)
			if result != vec.OriginalHash {
				t.Errorf("Original hash mismatch: got %s, want %s", result, vec.OriginalHash)
			}
		})

		t.Run(vec.Description+" - legacy hashes", func(t *testing.T) {
			legacyHashes := ComputeAssetLegacyHashes(vec.Payload)

			if len(legacyHashes) != vec.ExpectedLegacyCount {
				t.Errorf("Legacy hash count mismatch: got %d, want %d", len(legacyHashes), vec.ExpectedLegacyCount)
			}

			if vec.ExpectedLegacyCount > 0 {
				assertContains(t, legacyHashes, vec.ExpectedWithoutIsNFT, "without_is_nft")
				assertContains(t, legacyHashes, vec.ExpectedWithoutKindType, "without_kind_type")
				assertContains(t, legacyHashes, vec.ExpectedWithoutBoth, "without_both")
			}
		})
	}
}

func assertContains(t *testing.T, hashes []string, expected, name string) {
	t.Helper()
	for _, h := range hashes {
		if h == expected {
			return
		}
	}
	t.Errorf("Missing %s hash %s in legacy hashes %v", name, expected, hashes)
}
