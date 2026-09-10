package helper

import (
	"encoding/json"
	"fmt"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// canonicalPrice is the projection of a price that a PRICEUPDATER signs: the five fields
// validatord's CurrencyPrice serialises, in its field order. Go's encoding/json emits
// struct fields in declaration order, so this order IS the signed byte sequence — do not
// reorder or add fields.
type canonicalPrice struct {
	Blockchain   string `json:"blockchain"`
	CurrencyFrom string `json:"currencyFrom"`
	CurrencyTo   string `json:"currencyTo"`
	Decimals     string `json:"decimals"`
	Rate         string `json:"rate"`
}

// PriceSignedBytes returns the exact bytes a price signature covers.
func PriceSignedBytes(price *model.Price) ([]byte, error) {
	if price == nil {
		return nil, fmt.Errorf("price cannot be nil")
	}
	return json.Marshal(canonicalPrice{
		Blockchain:   price.Blockchain,
		CurrencyFrom: price.CurrencyFrom,
		CurrencyTo:   price.CurrencyTo,
		Decimals:     price.Decimals,
		Rate:         price.Rate,
	})
}

// VerifyPrice checks a price against the PRICEUPDATER keys in a verified rules container.
//
// The container decides whether prices must be signed at all. It is SuperAdmin-verified,
// so "this tenant has no price signer" is trustworthy, whereas "this price carries no
// signatures" is not — which is why a stripped signatures array on a tenant that DOES
// have a PRICEUPDATER is an error rather than a skip.
//
// Rate and decimals feed amount conversion, so an unverified price is a wrong number a
// caller would otherwise act on.
//
//	verified container ─▶ any PRICEUPDATER key? ──no──▶ PASS (tenant does not sign prices)
//	                              │yes
//	                              ▼
//	                      price.Signatures empty? ──yes──▶ IntegrityError
//	                              │no                      (only the container may excuse
//	                              ▼                         a price from being signed)
//	                      one entry valid over the
//	                      canonical 5-field JSON? ──no──▶ IntegrityError
//	                              │yes
//	                              ▼
//	                             PASS
func VerifyPrice(price *model.Price, rulesContainer *model.DecodedRulesContainer) error {
	if price == nil {
		return &model.IntegrityError{Message: "price cannot be nil"}
	}
	if rulesContainer == nil {
		return &model.IntegrityError{
			Message: "rules container required for price signature verification",
		}
	}

	keys := rulesContainer.PriceUpdaterKeys()
	if len(keys) == 0 {
		// This tenant does not sign prices.
		return nil
	}

	if len(price.Signatures) == 0 {
		return &model.IntegrityError{
			Message: fmt.Sprintf(
				"price %s/%s carries no signatures but the rules container configures a PRICEUPDATER",
				price.CurrencyFrom, price.CurrencyTo),
		}
	}

	data, err := PriceSignedBytes(price)
	if err != nil {
		return &model.IntegrityError{
			Message: fmt.Sprintf("cannot build the signed form of price %s/%s: %v",
				price.CurrencyFrom, price.CurrencyTo, err),
		}
	}

	for _, sig := range price.Signatures {
		if sig.Signature == "" {
			continue
		}
		if IsValidSignature(data, sig.Signature, keys) {
			return nil
		}
	}

	return &model.IntegrityError{
		Message: fmt.Sprintf(
			"no PRICEUPDATER signature verifies for price %s/%s (%d signature(s) offered)",
			price.CurrencyFrom, price.CurrencyTo, len(price.Signatures)),
	}
}

// VerifyPrices verifies each price, returning the first failure.
func VerifyPrices(prices []*model.Price, rulesContainer *model.DecodedRulesContainer) error {
	for _, price := range prices {
		if price == nil {
			continue
		}
		if err := VerifyPrice(price, rulesContainer); err != nil {
			return err
		}
	}
	return nil
}
