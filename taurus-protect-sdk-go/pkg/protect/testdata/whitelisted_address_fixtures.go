// Package testdata contains test fixtures for whitelisted address unit tests.
//
// This file provides real API response data and legacy hash test cases for verifying
// the whitelisted address integrity verification functionality, including backward
// compatibility with addresses signed before schema changes.
package testdata

// =============================================================================
// Real API Response Data
// =============================================================================
// These values were captured from a live integration test to ensure the SDK
// correctly handles actual API responses.

// RealPayloadAsString is the payload from a real whitelisted address (ID: 36663, ALGO blockchain).
// This represents the current schema format with all fields including contractType and tnParticipantID.
const RealPayloadAsString = `{"currency":"ALGO","addressType":"individual","address":"P4QCJV2YYLAEULGLJQAW4XTU3EBOHWL5C46I5SPLH2H7AJEE367ZDACV5A","memo":"","label":"TN_Bank ACC Cockroach_WTRTest","customerId":"","exchangeAccountId":"","linkedInternalAddresses":[],"contractType":"","tnParticipantID":"84dc35e3-0af8-4b6b-be75-785f4b149d16"}`

// RealMetadataHash is the expected SHA-256 hash of RealPayloadAsString.
// This hash is stored in the envelope's metadata and verified against SuperAdmin signatures.
const RealMetadataHash = "830063cfa8c1dbd696d670fc8360e85fbc57c3ffa66d22358b9a7d6befabb2f0"

// =============================================================================
// Legacy Hash Test Case 1: Address 509 (contractType added after signing)
// =============================================================================
// This address was signed before the 'contractType' field was added to the schema.
// The original payload did NOT include contractType, but the current API returns it.
//
// Timeline:
// - 2020: Address signed with original schema (no contractType)
// - Later: 'contractType' field added to schema
// - Now: API returns payload with contractType:"" but signatures were made on original payload
//
// Verification must try legacy hash (without contractType) when current hash fails.

// Case1LegacyPayload is the original payload format WITHOUT contractType (2020 schema).
// This is what was actually signed by the HSM.
const Case1LegacyPayload = `{"currency":"ETH","addressType":"individual","address":"0x012566A179a935ACF1d81d4D237495DE933D12E6","memo":"","label":"CMTA20-KYC - 0x012566A179a935ACF1d81d4D237495DE933D12E6 (request 6826)","customerId":"","exchangeAccountId":"","linkedInternalAddresses":[]}`

// Case1CurrentPayload is the current payload format WITH contractType (current schema).
// This is what the API returns now, but it doesn't match the original signature.
const Case1CurrentPayload = `{"currency":"ETH","addressType":"individual","address":"0x012566A179a935ACF1d81d4D237495DE933D12E6","memo":"","label":"CMTA20-KYC - 0x012566A179a935ACF1d81d4D237495DE933D12E6 (request 6826)","customerId":"","exchangeAccountId":"","linkedInternalAddresses":[],"contractType":""}`

// Case1LegacyHash is the SHA-256 hash of Case1LegacyPayload.
// This hash matches what the HSM signed.
const Case1LegacyHash = "cda66e821ec26f2432a717feaa1ef49be39a7ad9e93b6b8fcdce606659e964df"

// Case1CurrentHash is the SHA-256 hash of Case1CurrentPayload.
// This hash does NOT match any signatures because the payload has been modified.
const Case1CurrentHash = "d95ae4359bea509c2542acf410649f1e361233da5e1ac7c7a198b6d6a2bbbe1f"

// =============================================================================
// Legacy Hash Test Case 2: Address 391 (contractType AND labels added after signing)
// =============================================================================
// This address was signed before TWO schema changes:
// 1. 'contractType' field was added to the root payload
// 2. 'label' field was added to linkedInternalAddresses objects
//
// The verification must try multiple legacy hash strategies:
// - Strategy 1: Remove both contractType AND labels from linkedInternalAddresses
// - Strategy 2: Remove only labels from linkedInternalAddresses (keep contractType)
//
// This handles addresses signed at different points during the schema evolution.

// Case2CurrentPayload is the current payload format with contractType AND labels in linkedInternalAddresses.
// This is what the API returns now.
const Case2CurrentPayload = `{"currency":"ETH","addressType":"individual","address":"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380","memo":"","label":"20200324 test address 2","customerId":"1","exchangeAccountId":"","linkedInternalAddresses":[{"id":"10","address":"0x589ef3d7585f54f0539e24253050887c691c9bd8","label":"client 0 ETH "},{"id":"13","address":"0x669805f31178faf0dca39c8a5c49ecc531b5156e","label":"ETH internal client 02.02"},{"id":"20","address":"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d","label":"LBR 07.02"},{"id":"98","address":"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e","label":"ETH LBR internal client 26.02"},{"id":"25","address":"0x9bc28e6710f5bb2511372987f613a436618e28ad","label":"LBR IC 13.02"}],"contractType":""}`

// Case2LegacyPayload is the original payload format (oldest schema).
// This is what was signed: NO contractType AND NO labels in linkedInternalAddresses.
const Case2LegacyPayload = `{"currency":"ETH","addressType":"individual","address":"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380","memo":"","label":"20200324 test address 2","customerId":"1","exchangeAccountId":"","linkedInternalAddresses":[{"id":"10","address":"0x589ef3d7585f54f0539e24253050887c691c9bd8"},{"id":"13","address":"0x669805f31178faf0dca39c8a5c49ecc531b5156e"},{"id":"20","address":"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d"},{"id":"98","address":"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e"},{"id":"25","address":"0x9bc28e6710f5bb2511372987f613a436618e28ad"}]}`

// Case2LegacyHash is the SHA-256 hash of Case2LegacyPayload (oldest schema).
// This matches signatures on addresses signed before both schema changes.
const Case2LegacyHash = "88e4e456f7ca1fc4ca415c6c571f828c0eb047e9f15f36d547c103b2ea0def9b"

// Case2CurrentHash is the SHA-256 hash of Case2CurrentPayload.
// This does NOT match any signatures because the payload has been modified twice.
const Case2CurrentHash = "7d62d7f78ed55c716ea1278473d6cac5b60a31e1df941873118932822df32b03"

// =============================================================================
// Legacy Hash Strategy 2: Labels removed only (contractType kept)
// =============================================================================
// This handles addresses signed AFTER contractType was added but BEFORE labels
// were added to linkedInternalAddresses.
//
// Timeline:
// - 2020: Original schema (no contractType, no labels in linkedInternalAddresses)
// - Mid-2020: contractType field added
// - Later: labels added to linkedInternalAddresses
//
// Addresses signed between these two changes need Strategy 2.

// Strategy2LegacyPayload has contractType but NO labels in linkedInternalAddresses.
// Use this for addresses signed after contractType but before labels were added.
const Strategy2LegacyPayload = `{"currency":"ETH","addressType":"individual","address":"0x5c2697f5faf6faaeefa9f2fa1e5a18bb248a6380","memo":"","label":"20200324 test address 2","customerId":"1","exchangeAccountId":"","linkedInternalAddresses":[{"id":"10","address":"0x589ef3d7585f54f0539e24253050887c691c9bd8"},{"id":"13","address":"0x669805f31178faf0dca39c8a5c49ecc531b5156e"},{"id":"20","address":"0x6cf6ab78ebb80d7dde4ec11d7f139ea4d0210c3d"},{"id":"98","address":"0x2dc5b7f8f94cbb2a1d1306dda130325d7384296e"},{"id":"25","address":"0x9bc28e6710f5bb2511372987f613a436618e28ad"}],"contractType":""}`

// =============================================================================
// Test Helper Functions
// =============================================================================

// LegacyHashStrategies returns the ordered list of legacy payload transformations
// that should be tried when verifying whitelisted address integrity.
//
// The order matters - try most recent schema first, then progressively older:
// 1. Current payload (no transformation)
// 2. Remove labels from linkedInternalAddresses only (keep contractType)
// 3. Remove both contractType AND labels from linkedInternalAddresses
func LegacyHashStrategies() []string {
	return []string{
		"current",                  // No transformation
		"remove_labels",            // Strategy 2: Remove labels from linkedInternalAddresses
		"remove_contractType",      // Strategy 1a: Remove contractType only
		"remove_labels_and_contractType", // Strategy 1: Remove both
	}
}
