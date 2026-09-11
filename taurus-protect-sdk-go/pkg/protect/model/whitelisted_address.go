package model

import (
	"fmt"
	"time"
)

// WhitelistedAddress represents a whitelisted address in the system.
type WhitelistedAddress struct {
	// ID is the unique identifier for the whitelisted address.
	ID string `json:"id"`
	// TenantID is the tenant identifier.
	TenantID string `json:"tenant_id"`
	// Address is the blockchain address string (extracted from metadata).
	Address string `json:"address"`
	// Label is the human-readable label for this address (extracted from metadata).
	Label string `json:"label"`
	// Memo is the optional memo field for Stellar or destination tag for Ripple.
	Memo string `json:"memo,omitempty"`
	// CustomerId is the customer ID for external reconciliation.
	CustomerId string `json:"customer_id,omitempty"`
	// Blockchain is the blockchain network (e.g., "ETH", "BTC").
	Blockchain string `json:"blockchain"`
	// Network is the network type (e.g., "mainnet", "testnet").
	Network string `json:"network"`
	// Status is the current status of the whitelisted address.
	Status string `json:"status"`
	// Action is the pending action type if any.
	Action string `json:"action,omitempty"`
	// AddressType is the type of address (individual, exchange, contract, etc.).
	AddressType string `json:"address_type,omitempty"`
	// ContractType is the smart contract type (e.g., "CMTA20", "ERC20") for contract addresses.
	ContractType string `json:"contract_type,omitempty"`
	// ExchangeAccountId is the exchange account ID when the address belongs to an exchange.
	ExchangeAccountId int64 `json:"exchange_account_id,omitempty"`
	// Rule is the governance rule applied to this address.
	Rule string `json:"rule,omitempty"`
	// RulesContainer is the serialized rules container.
	RulesContainer string `json:"rules_container,omitempty"`
	// RulesContainerHash is the hash of the rules container.
	RulesContainerHash string `json:"rules_container_hash,omitempty"`
	// RulesSignatures contains the super-admin signatures for the rules.
	RulesSignatures string `json:"rules_signatures,omitempty"`
	// VisibilityGroupID is the visibility group this address belongs to.
	VisibilityGroupID string `json:"visibility_group_id,omitempty"`
	// TnParticipantID is the Taurus Network participant ID.
	TnParticipantID string `json:"tn_participant_id,omitempty"`
	// Metadata contains additional metadata about the address.
	// Uses WhitelistedAssetMetadata which is shared between assets and addresses.
	Metadata *WhitelistedAssetMetadata `json:"metadata,omitempty"`
	// Scores contains the risk scores from various providers.
	Scores []Score `json:"scores,omitempty"`
	// Trails contains the audit trail of actions on this address.
	// Uses shared Trail type.
	Trails []Trail `json:"trails,omitempty"`
	// CreatedAt is the creation timestamp extracted from the "created" trail action.
	CreatedAt *time.Time `json:"created_at,omitempty"`
	// Approvers contains the approval requirements and status.
	// Uses shared Approvers type.
	Approvers *Approvers `json:"approvers,omitempty"`
	// Attributes are custom key-value attributes.
	Attributes []WhitelistedAddressAttribute `json:"attributes,omitempty"`
	// SignedAddress contains the cryptographic signature data.
	SignedAddress *SignedWhitelistedAddress `json:"signed_address,omitempty"`
	// LinkedInternalAddresses is the list of internal addresses linked to this whitelisted address.
	LinkedInternalAddresses []InternalAddress `json:"linked_internal_addresses,omitempty"`
	// LinkedWallets is the list of internal wallets that can send to this whitelisted address.
	LinkedWallets []InternalWallet `json:"linked_wallets,omitempty"`
}

// InternalAddress represents an internal address linked to a whitelisted address.
type InternalAddress struct {
	// ID is the internal address identifier.
	ID int64 `json:"id"`
	// Label is the human-readable label for the address.
	Label string `json:"label,omitempty"`
}

// InternalWallet represents an internal wallet linked to a whitelisted address.
type InternalWallet struct {
	// ID is the internal wallet identifier.
	ID int64 `json:"id"`
	// Path is the wallet path.
	Path string `json:"path,omitempty"`
	// Label is the human-readable label for the wallet.
	Label string `json:"label,omitempty"`
}

// WhitelistedAddressAttribute represents a custom attribute on a whitelisted address.
type WhitelistedAddressAttribute struct {
	// ID is the attribute identifier.
	ID string `json:"id"`
	// Key is the attribute name.
	Key string `json:"key"`
	// Value is the attribute value.
	Value string `json:"value"`
	// ContentType is the MIME type of the value.
	ContentType string `json:"content_type,omitempty"`
	// Owner is the user who created the attribute.
	Owner string `json:"owner,omitempty"`
	// Type is the attribute type.
	Type string `json:"type,omitempty"`
	// Subtype is the attribute subtype.
	Subtype string `json:"subtype,omitempty"`
	// IsFile indicates if the value is a file reference.
	IsFile bool `json:"is_file"`
}

// SignedWhitelistedAddress contains the cryptographic signature data for a whitelisted address.
type SignedWhitelistedAddress struct {
	// Payload is the base64-encoded signed payload.
	Payload string `json:"payload,omitempty"`
	// Signatures contains the list of signatures.
	// Uses shared WhitelistSignature type.
	Signatures []WhitelistSignature `json:"signatures,omitempty"`
}

// Score represents a risk score from a scoring provider.
type Score struct {
	// ID is the score identifier.
	ID string `json:"id"`
	// Provider is the scoring provider name.
	Provider string `json:"provider"`
	// Type is the type of score.
	Type string `json:"type"`
	// Score is the score value.
	Score string `json:"score"`
	// UpdateDate is when the score was last updated as a string.
	UpdateDate string `json:"update_date"`
}

// ListWhitelistedAddressesOptions contains options for listing whitelisted addresses.
type ListWhitelistedAddressesOptions struct {
	// Limit is the maximum number of addresses to return.
	Limit int64
	// Offset is the number of addresses to skip.
	Offset int64
	// Blockchain filters by blockchain.
	Blockchain string
	// Network filters by network.
	Network string
	// Query searches address names or values.
	Query string
	// AddressType filters by address type.
	AddressType string
	// IDs filters by specific whitelisted address IDs.
	IDs []string
	// Addresses filters by specific address values.
	Addresses []string
	// IncludeForApproval includes addresses pending approval.
	IncludeForApproval bool
	// AllowedForWalletID / AllowedForAddressID keep only the destinations that
	// wallet or address is permitted to send to.
	AllowedForWalletID  string
	AllowedForAddressID string
	// TNParticipantID filters to one Taurus-NETWORK counterparty.
	TNParticipantID string
	// TagIDs, ContractTypes and ExchangeAccountIDs are list filters. The endpoint
	// also has `actions`, absent from the generated client's stale spec.
	TagIDs             []string
	ContractTypes      []string
	ExchangeAccountIDs []string
}

// WhitelistedAddressEnvelope wraps a whitelisted address with its verification data.
// This is returned by GetWhitelistedAddressEnvelope after the 6-step verification flow.
type WhitelistedAddressEnvelope struct {
	// ID is the unique identifier for the whitelisted address.
	ID string `json:"id"`
	// Blockchain is the blockchain network (e.g., "ETH", "BTC").
	Blockchain string `json:"blockchain"`
	// Network is the network type (e.g., "mainnet", "testnet").
	Network string `json:"network"`
	// Metadata contains the address metadata.
	Metadata *WhitelistedAssetMetadata `json:"metadata,omitempty"`
	// SignedAddress contains the cryptographic signature data.
	SignedAddress *SignedWhitelistedAddress `json:"signed_address,omitempty"`
	// RulesContainer is the base64-encoded rules container.
	RulesContainer string `json:"rules_container,omitempty"`
	// RulesSignatures is the base64-encoded rules signatures.
	RulesSignatures string `json:"rules_signatures,omitempty"`
	// LinkedInternalAddresses is the list of internal addresses linked to this whitelisted address.
	LinkedInternalAddresses []InternalAddress `json:"linked_internal_addresses,omitempty"`
	// LinkedWallets is the list of internal wallets that can send to this whitelisted address.
	LinkedWallets []InternalWallet `json:"linked_wallets,omitempty"`

	// verifiedWhitelistedAddress is the verified whitelisted address parsed from the payload.
	verifiedWhitelistedAddress *WhitelistedAddress
	// verifiedRulesContainer is the verified and decoded rules container.
	verifiedRulesContainer *DecodedRulesContainer
}

// WhitelistedAddress returns the verified whitelisted address.
// This is populated after successful verification.
func (e *WhitelistedAddressEnvelope) WhitelistedAddress() *WhitelistedAddress {
	return e.verifiedWhitelistedAddress
}

// DecodedRulesContainer returns the verified and decoded rules container.
// This is populated after successful verification.
func (e *WhitelistedAddressEnvelope) DecodedRulesContainer() *DecodedRulesContainer {
	return e.verifiedRulesContainer
}

// SetVerified sets the verified address and rules container.
// This is called internally after successful verification.
func (e *WhitelistedAddressEnvelope) SetVerified(addr *WhitelistedAddress, rules *DecodedRulesContainer) {
	e.verifiedWhitelistedAddress = addr
	e.verifiedRulesContainer = rules
}

// ExcludedWhitelistedAddress names a row dropped from a list because it failed
// integrity verification, and why.
type ExcludedWhitelistedAddress struct {
	// ID is the address ID, or "" when the row carried no usable ID.
	ID string
	// Reason is why the row failed verification.
	Reason string
}

// WhitelistedAddressResult contains the result of a paginated whitelisted-address
// list query.
//
// Verification is lenient by design: a row that cannot be verified is excluded rather
// than failing the whole call, because one bad row used to deny access to every good
// one — and listing is how an operator finds the bad row. Excluding stays fail-closed,
// since an omitted destination cannot be selected.
//
// The omission is reported rather than silent: a shortened list must never be mistaken
// for a complete one.
// ListWhitelistedAddressesForApprovalOptions filters the approval queue.
//
// The approval queue is a different ENDPOINT, not a status filter on the general list:
// it is scoped to what the calling user may act on, which no filter reproduces.
type ListWhitelistedAddressesForApprovalOptions struct {
	// Limit is the maximum number of rows to return.
	Limit int64
	// Offset is the number of rows to skip.
	Offset int64
	// IDs filters by specific whitelisted address IDs.
	IDs []string
	// Blockchain filters by blockchain symbol.
	Blockchain string
	// Network filters by network.
	Network string
	// AddressType filters by address type.
	AddressType string
	// Query matches customer id, address, blockchain, label, memo and address type.
	Query string
	// IncludeAlreadySignedByUser includes rows the calling user has already signed,
	// which is how an approver tells "waiting for me" from "waiting for someone else".
	IncludeAlreadySignedByUser bool
}

type WhitelistedAddressResult struct {
	// Addresses is the list of verified whitelisted addresses in the current page.
	Addresses []*WhitelistedAddress
	// Pagination carries the page window, with TotalItems already reduced by the
	// number of excluded rows so HasMore stays honest.
	Pagination *Pagination
	// ExcludedUnverified names the rows dropped from Addresses, with the reason.
	ExcludedUnverified []ExcludedWhitelistedAddress
}

// WhitelistedAddressApproval is the set of rows an approver reviewed, carrying the metadata hash
// each one had AT REVIEW TIME. It is the content pin the approval path signs against.
//
// Why this type exists rather than a plain []string of ids. The approval API accepts only ids:
// the SDK re-reads them and signs whatever the server returns under those ids. Nothing bound the
// approver's intent to the bytes signed, so a response-controlling server could answer the
// id-filtered re-read with a different row — one whose existing signatures already satisfy the
// container it presents — and harvest a genuine approver signature over content the approver
// never saw. This is the same shape GovernanceRuleService.ApproveRulesProposal was hardened
// against with its mandatory expectedContainerHash.
//
// `pinned` is unexported on purpose, which is the idiom helper.VerifiedAsset already uses here:
// Go permits model.WhitelistedAddressApproval{} from another package, so the value is forgeable
// but USELESS — a hand-built one pins nothing and IsEmpty reports true, which the approval
// refuses. Only a verified read can mint a usable one.
type WhitelistedAddressApproval struct {
	// pinned maps row id -> the metadata hash that row carried when it was reviewed.
	pinned map[string]string
}

// IDs returns the pinned row ids in no particular order. The approval path sorts them itself,
// because the endpoint requires ascending order and the signed array must not depend on the
// order the caller happened to select in.
func (a *WhitelistedAddressApproval) IDs() []string {
	if a == nil {
		return nil
	}
	ids := make([]string, 0, len(a.pinned))
	for id := range a.pinned {
		ids = append(ids, id)
	}
	return ids
}

// PinnedHash returns the reviewed metadata hash for id, and whether it was pinned at all.
func (a *WhitelistedAddressApproval) PinnedHash(id string) (string, bool) {
	if a == nil {
		return "", false
	}
	hash, ok := a.pinned[id]
	return hash, ok
}

// IsEmpty reports whether this selection pins nothing — true for a nil or hand-built value.
func (a *WhitelistedAddressApproval) IsEmpty() bool {
	return a == nil || len(a.pinned) == 0
}

// Select pins the given ids from this verified read.
//
// An id that is not in the result is an error rather than a silent omission: it means the caller
// is trying to approve something this read did not return — either it was excluded as
// unverifiable, or it was never on the page — and approving fewer rows than asked for would tell
// the approver they approved more than they did.
func (r *WhitelistedAddressResult) Select(ids ...string) (*WhitelistedAddressApproval, error) {
	if r == nil {
		return nil, fmt.Errorf("cannot select from a nil result")
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("cannot select an empty set of ids")
	}

	byID := make(map[string]*WhitelistedAddress, len(r.Addresses))
	for _, addr := range r.Addresses {
		if addr != nil {
			byID[addr.ID] = addr
		}
	}

	pinned := make(map[string]string, len(ids))
	for _, id := range ids {
		addr, ok := byID[id]
		if !ok {
			return nil, fmt.Errorf("whitelisted address %s is not in this verified read: it was "+
				"either excluded as unverifiable or not on this page", id)
		}
		if addr.Metadata == nil || addr.Metadata.Hash == "" {
			return nil, &IntegrityError{
				Message: fmt.Sprintf("whitelisted address %s carries no metadata hash, so there is "+
					"nothing to pin the approval to", id),
			}
		}
		pinned[id] = addr.Metadata.Hash
	}
	return &WhitelistedAddressApproval{pinned: pinned}, nil
}

// SelectAll pins every row this verified read returned.
//
// Note it pins what SURVIVED verification, not what the server sent: rows in
// ExcludedUnverified are not included, so a caller who wants to know about them must read
// that field. Approving is all-or-nothing over what is pinned here.
func (r *WhitelistedAddressResult) SelectAll() (*WhitelistedAddressApproval, error) {
	if r == nil {
		return nil, fmt.Errorf("cannot select from a nil result")
	}
	ids := make([]string, 0, len(r.Addresses))
	for _, addr := range r.Addresses {
		if addr != nil {
			ids = append(ids, addr.ID)
		}
	}
	if len(ids) == 0 {
		return nil, fmt.Errorf("this read returned no verified addresses to approve")
	}
	return r.Select(ids...)
}
