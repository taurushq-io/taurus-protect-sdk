package helper

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"testing"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

func samplePrice() *model.Price {
	return &model.Price{
		Blockchain:   "ETH",
		CurrencyFrom: "ETH",
		CurrencyTo:   "USD",
		Decimals:     "18",
		Rate:         "2500.00",
	}
}

func containerWith(roles []string, key *ecdsa.PublicKey) *model.DecodedRulesContainer {
	return &model.DecodedRulesContainer{
		Users: []*model.RuleUser{{ID: "price@bank.com", PublicKey: key, Roles: roles}},
	}
}

// The canonical form is validatord's CurrencyPrice JSON projection. Pinned by value: a
// reordered or extended struct silently changes every signature this SDK will accept.
func TestPriceSignedBytes_IsTheCanonicalProjection(t *testing.T) {
	data, err := PriceSignedBytes(samplePrice())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := `{"blockchain":"ETH","currencyFrom":"ETH","currencyTo":"USD","decimals":"18","rate":"2500.00"}`
	if string(data) != want {
		t.Errorf("canonical form drifted:\n got %s\nwant %s", data, want)
	}
}

func TestVerifyPrice_AcceptsAPriceUpdaterSignature(t *testing.T) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	price := samplePrice()
	data, _ := PriceSignedBytes(price)
	signature, err := crypto.SignData(key, data)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	price.Signatures = []model.PriceSignature{{UserID: "price@bank.com", Signature: signature}}

	container := containerWith([]string{"PRICEUPDATER"}, &key.PublicKey)
	if err := VerifyPrice(price, container); err != nil {
		t.Errorf("a genuine PRICEUPDATER signature must verify, got %v", err)
	}
}

func TestVerifyPrice_RejectsASignatureFromANonPriceUpdater(t *testing.T) {
	signer, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	updater, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	price := samplePrice()
	data, _ := PriceSignedBytes(price)
	signature, _ := crypto.SignData(signer, data)
	price.Signatures = []model.PriceSignature{{UserID: "someone@bank.com", Signature: signature}}

	container := containerWith([]string{"PRICEUPDATER"}, &updater.PublicKey)
	if err := VerifyPrice(price, container); err == nil {
		t.Error("a signature by a key without the PRICEUPDATER role must not verify")
	}
}

func TestVerifyPrice_RejectsATamperedRate(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	price := samplePrice()
	data, _ := PriceSignedBytes(price)
	signature, _ := crypto.SignData(key, data)
	price.Signatures = []model.PriceSignature{{UserID: "price@bank.com", Signature: signature}}

	// The rate is what a caller converts with.
	price.Rate = "1.00"

	container := containerWith([]string{"PRICEUPDATER"}, &key.PublicKey)
	if err := VerifyPrice(price, container); err == nil {
		t.Error("a rate altered after signing must not verify")
	}
}

// Decision 13A: the verified container decides. A stripped signatures array on a tenant
// that DOES configure a PRICEUPDATER is the fail-open case, so it must be an error.
func TestVerifyPrice_StrippedSignaturesAreRejectedWhenAPriceUpdaterExists(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	price := samplePrice() // no signatures

	container := containerWith([]string{"PRICEUPDATER"}, &key.PublicKey)
	err := VerifyPrice(price, container)
	if err == nil {
		t.Fatal("stripping the signatures must not bypass verification")
	}
	var integrityErr *model.IntegrityError
	if !errors.As(err, &integrityErr) {
		t.Errorf("expected an IntegrityError, got %T", err)
	}
}

// A tenant with no price signer legitimately serves unsigned prices, and the container
// saying so is trustworthy because it is SuperAdmin-verified.
func TestVerifyPrice_PassesThroughWhenNoPriceUpdaterIsConfigured(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	container := containerWith([]string{"REQUESTAPPROVER"}, &key.PublicKey)

	if err := VerifyPrice(samplePrice(), container); err != nil {
		t.Errorf("a tenant without a PRICEUPDATER must not be forced to sign, got %v", err)
	}
}

func TestVerifyPrice_RequiresARulesContainer(t *testing.T) {
	if err := VerifyPrice(samplePrice(), nil); err == nil {
		t.Error("expected an error without a rules container")
	}
}

func TestPriceUpdaterKeys_SelectsOnlyThatRole(t *testing.T) {
	updater, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	other, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	container := &model.DecodedRulesContainer{
		Users: []*model.RuleUser{
			{ID: "a", PublicKey: &other.PublicKey, Roles: []string{"HSMSLOT"}},
			{ID: "b", PublicKey: &updater.PublicKey, Roles: []string{"PRICEUPDATER"}},
			{ID: "c", PublicKey: nil, Roles: []string{"PRICEUPDATER"}},
		},
	}

	keys := container.PriceUpdaterKeys()
	if len(keys) != 1 || keys[0] != &updater.PublicKey {
		t.Errorf("expected exactly the one usable PRICEUPDATER key, got %d", len(keys))
	}
}

func TestVerifyPrices_ReportsTheFirstFailure(t *testing.T) {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	good := samplePrice()
	data, _ := PriceSignedBytes(good)
	signature, _ := crypto.SignData(key, data)
	good.Signatures = []model.PriceSignature{{Signature: signature}}

	bad := samplePrice()
	bad.CurrencyTo = "EUR"

	container := containerWith([]string{"PRICEUPDATER"}, &key.PublicKey)
	if err := VerifyPrices([]*model.Price{good, bad}, container); err == nil {
		t.Error("one unverifiable price must fail the batch")
	}
}
