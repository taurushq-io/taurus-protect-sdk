package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/cache"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// decodeToleranceVectorsRelPath points at the vectors every SDK test suite consumes, so the four
// SDKs treat data they do not know the same way: unknown fields are kept, unknown enum values keep
// their raw string, and caller-supplied enum values are sent verbatim.
const decodeToleranceVectorsRelPath = "../../../../scripts/resources/decode-tolerance-vectors.json"

type decodeToleranceVector struct {
	Name      string          `json:"name"`
	Kind      string          `json:"kind"`
	Model     string          `json:"model"`
	Enum      string          `json:"enum"`
	Operation string          `json:"operation"`
	Wire      json.RawMessage `json:"wire"`
	Request   json.RawMessage `json:"request"`
	Reply     json.RawMessage `json:"reply"`
	Expected  json.RawMessage `json:"expected"`
}

func decodeToleranceVectors(t *testing.T, kind string) []decodeToleranceVector {
	t.Helper()
	path := filepath.Clean(decodeToleranceVectorsRelPath)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot read shared vectors file %s: %v", path, err)
	}
	var all []decodeToleranceVector
	if err := json.Unmarshal(raw, &all); err != nil {
		t.Fatalf("cannot parse shared vectors file: %v", err)
	}
	var out []decodeToleranceVector
	for _, v := range all {
		if v.Kind == kind {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		t.Fatalf("shared vectors file has no %q vectors", kind)
	}
	return out
}

func mustDecodeJSON(t *testing.T, raw json.RawMessage, into any) {
	t.Helper()
	if err := json.Unmarshal(raw, into); err != nil {
		t.Fatalf("cannot parse vector field %s: %v", raw, err)
	}
}

var decodeToleranceModels = map[string]func() any{
	"tgvalidatordGetBalancesReply":                         func() any { return &openapi.TgvalidatordGetBalancesReply{} },
	"tgvalidatordStakeAccount":                             func() any { return &openapi.TgvalidatordStakeAccount{} },
	"tgvalidatordGetMultiFactorSignatureEntitiesInfoReply": func() any { return &openapi.TgvalidatordGetMultiFactorSignatureEntitiesInfoReply{} },
}

type generatedEnum interface {
	~string
	IsValid() bool
}

// decodeEnum reports the decoded value and whether it is one of the generated constants.
func decodeEnum[E generatedEnum](raw []byte) (string, bool, error) {
	var v E
	if err := json.Unmarshal(raw, &v); err != nil {
		return "", false, err
	}
	return string(v), v.IsValid(), nil
}

var decodeToleranceEnums = map[string]func([]byte) (string, bool, error){
	"tgvalidatordTokenType":                       decodeEnum[openapi.TgvalidatordTokenType],
	"tgvalidatordMultiFactorSignaturesEntityType": decodeEnum[openapi.TgvalidatordMultiFactorSignaturesEntityType],
	"tgvalidatordStakeAccountType":                decodeEnum[openapi.TgvalidatordStakeAccountType],
}

func TestDecodeToleranceModelVectors(t *testing.T) {
	for _, v := range decodeToleranceVectors(t, "model") {
		t.Run(v.Name, func(t *testing.T) {
			newModel, ok := decodeToleranceModels[v.Model]
			if !ok {
				t.Fatalf("no generated model registered for %s", v.Model)
			}
			var expected struct {
				Lossless                 bool           `json:"lossless"`
				RootAdditionalProperties map[string]any `json:"rootAdditionalProperties"`
			}
			mustDecodeJSON(t, v.Expected, &expected)

			decoded := newModel()
			if err := json.Unmarshal(v.Wire, decoded); err != nil {
				t.Fatalf("decoding %s failed: %v", v.Model, err)
			}
			extras := reflect.ValueOf(decoded).Elem().FieldByName("AdditionalProperties")
			if !extras.IsValid() {
				t.Fatalf("%s has no AdditionalProperties: unknown fields are not preserved", v.Model)
			}
			if got, _ := extras.Interface().(map[string]any); !reflect.DeepEqual(got, expected.RootAdditionalProperties) {
				t.Fatalf("root AdditionalProperties = %v, want %v", got, expected.RootAdditionalProperties)
			}
			if !expected.Lossless {
				return
			}
			out, err := json.Marshal(decoded)
			if err != nil {
				t.Fatalf("re-serializing %s failed: %v", v.Model, err)
			}
			var got, want any
			mustDecodeJSON(t, out, &got)
			mustDecodeJSON(t, v.Wire, &want)
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("re-serialization is lossy:\n got %s\nwant %s", out, v.Wire)
			}
		})
	}
}

func TestDecodeToleranceEnumVectors(t *testing.T) {
	for _, v := range decodeToleranceVectors(t, "enum") {
		t.Run(v.Name, func(t *testing.T) {
			decode, ok := decodeToleranceEnums[v.Enum]
			if !ok {
				t.Fatalf("no generated enum registered for %s", v.Enum)
			}
			var expected struct {
				Value string `json:"value"`
				Known bool   `json:"known"`
			}
			mustDecodeJSON(t, v.Expected, &expected)

			value, known, err := decode(v.Wire)
			if err != nil {
				t.Fatalf("decoding %s %s failed: %v", v.Enum, v.Wire, err)
			}
			if value != expected.Value || known != expected.Known {
				t.Fatalf("decoded (%q, known=%v), want (%q, known=%v)", value, known, expected.Value, expected.Known)
			}
		})
	}
}

// toleranceServer answers reply on path and fails the test on any other request, so an unknown
// value cannot silently trigger a call the vector does not describe.
type toleranceServer struct {
	mu       sync.Mutex
	requests []recordedRequest
}

func newToleranceClient(t *testing.T, path string, reply json.RawMessage) (*openapi.APIClient, *toleranceServer) {
	t.Helper()
	srv := &toleranceServer{}
	client := walkClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		srv.mu.Lock()
		srv.requests = append(srv.requests, recordedRequest{method: r.Method, path: r.URL.Path, query: r.URL.RawQuery, body: body})
		srv.mu.Unlock()
		if r.URL.Path != path {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		_, _ = w.Write(reply)
	})
	return client, srv
}

// assertBodyContains checks that the single recorded request body carries every expected field.
func (s *toleranceServer) assertBodyContains(t *testing.T, expected map[string]any) {
	t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.requests) != 1 {
		t.Fatalf("got %d requests, want 1", len(s.requests))
	}
	var body map[string]any
	if err := json.Unmarshal(s.requests[0].body, &body); err != nil {
		t.Fatalf("request body %q is not a JSON object: %v", s.requests[0].body, err)
	}
	for key, want := range expected {
		if got := body[key]; !reflect.DeepEqual(got, want) {
			t.Fatalf("request body %s = %v, want %v (body %s)", key, got, want, s.requests[0].body)
		}
	}
}

var decodeToleranceOperations = map[string]func(t *testing.T, v decodeToleranceVector){
	"getMultiFactorSignatureInfo": func(t *testing.T, v decodeToleranceVector) {
		var request struct {
			ID string `json:"id"`
		}
		var expected struct {
			EntityType string `json:"entityType"`
		}
		mustDecodeJSON(t, v.Request, &request)
		mustDecodeJSON(t, v.Expected, &expected)

		client, _ := newToleranceClient(t, "/api/rest/v1/multifactor-signature/"+request.ID, v.Reply)
		info, err := NewMultiFactorSignatureService(client).GetMultiFactorSignatureInfo(context.Background(), request.ID)
		if err != nil {
			t.Fatalf("GetMultiFactorSignatureInfo failed: %v", err)
		}
		if got := string(info.EntityType); got != expected.EntityType {
			t.Fatalf("EntityType = %q, want %q", got, expected.EntityType)
		}
	},
	"createMultiFactorSignatures": func(t *testing.T, v decodeToleranceVector) {
		var request struct {
			EntityIDs  []string `json:"entityIds"`
			EntityType string   `json:"entityType"`
		}
		var expected struct {
			RequestBody map[string]any `json:"requestBody"`
		}
		mustDecodeJSON(t, v.Request, &request)
		mustDecodeJSON(t, v.Expected, &expected)

		client, srv := newToleranceClient(t, "/api/rest/v1/multifactor-signatures", v.Reply)
		_, err := NewMultiFactorSignatureService(client).CreateMultiFactorSignatures(
			context.Background(), request.EntityIDs, model.MultiFactorSignatureEntityType(request.EntityType))
		if err != nil {
			t.Fatalf("CreateMultiFactorSignatures failed: %v", err)
		}
		srv.assertBodyContains(t, expected.RequestBody)
	},
	"queryAssetAddresses": func(t *testing.T, v decodeToleranceVector) {
		var request struct {
			AssetID     string `json:"assetId"`
			AddressType string `json:"addressType"`
		}
		var expected struct {
			RequestBody map[string]any `json:"requestBody"`
			Rows        []struct {
				Address     string `json:"address"`
				AddressType string `json:"addressType"`
				Verified    bool   `json:"verified"`
			} `json:"rows"`
			ExcludedCount int `json:"excludedCount"`
		}
		mustDecodeJSON(t, v.Request, &request)
		mustDecodeJSON(t, v.Expected, &expected)

		client, srv := newToleranceClient(t, "/api/rest/v2/assets/"+request.AssetID+"/addresses/query", v.Reply)
		rules := cache.NewRulesContainerCache(time.Hour, func(context.Context) (*model.DecodedRulesContainer, error) {
			return &model.DecodedRulesContainer{}, nil
		})
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			t.Fatal(err)
		}
		whitelisted := NewWhitelistedAddressServiceWithVerification(client, &WhitelistedAddressServiceConfig{
			SuperAdminKeys: []*ecdsa.PublicKey{&key.PublicKey}, MinValidSignatures: 1})
		svc := NewAssetService(client, rules, NewAddressService(client, rules), whitelisted)

		result, err := svc.QueryAssetAddresses(context.Background(), request.AssetID,
			&model.QueryAssetAddressesOptions{AddressType: request.AddressType})
		if err != nil {
			t.Fatalf("QueryAssetAddresses failed: %v", err)
		}
		srv.assertBodyContains(t, expected.RequestBody)
		if len(result.Addresses) != len(expected.Rows) {
			t.Fatalf("got %d rows, want %d", len(result.Addresses), len(expected.Rows))
		}
		for i, want := range expected.Rows {
			got := result.Addresses[i]
			if got.Address != want.Address || got.AddressType != want.AddressType || got.Verified != want.Verified {
				t.Fatalf("row %d = {%q %q verified=%v}, want {%q %q verified=%v}", i,
					got.Address, got.AddressType, got.Verified, want.Address, want.AddressType, want.Verified)
			}
		}
		if len(result.ExcludedUnverified) != expected.ExcludedCount {
			t.Fatalf("ExcludedUnverified = %v, want %d entries", result.ExcludedUnverified, expected.ExcludedCount)
		}
	},
}

func TestDecodeToleranceServiceVectors(t *testing.T) {
	for _, v := range decodeToleranceVectors(t, "service") {
		t.Run(v.Name, func(t *testing.T) {
			run, ok := decodeToleranceOperations[v.Operation]
			if !ok {
				t.Fatalf("no runner registered for operation %s", v.Operation)
			}
			run(t, v)
		})
	}
}

// The generated decoder deletes "payload" from the extras map; without that, the raw payload the
// post-generation patch removes from TgvalidatordMetadata would come back through
// AdditionalProperties, where it could be read unhashed.
func TestMetadataPayloadStaysOutOfAdditionalProperties(t *testing.T) {
	var metadata openapi.TgvalidatordMetadata
	wire := `{"hash":"h","payload":[{"key":"value"}],"payloadAsString":"s","futureMetadataField":"f"}`
	if err := json.Unmarshal([]byte(wire), &metadata); err != nil {
		t.Fatalf("decoding metadata failed: %v", err)
	}
	extras := reflect.ValueOf(&metadata).Elem().FieldByName("AdditionalProperties")
	if !extras.IsValid() {
		t.Fatal("TgvalidatordMetadata has no AdditionalProperties: unknown fields are not preserved")
	}
	got, _ := extras.Interface().(map[string]any)
	if _, ok := got["payload"]; ok {
		t.Fatalf("payload leaked into AdditionalProperties: %v", got)
	}
	if got["futureMetadataField"] != "f" {
		t.Fatalf("AdditionalProperties = %v, want futureMetadataField kept", got)
	}
}
