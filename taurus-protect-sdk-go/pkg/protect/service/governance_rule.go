package service

import (
	"context"
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash"
	"sort"
	"strconv"
	"sync"

	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/internal/openapi"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/crypto"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/helper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/mapper"
	"github.com/taurushq-io/taurus-protect-sdk/taurus-protect-sdk-go/pkg/protect/model"
)

// maxVerifiedRulesets bounds the verification memo.
const maxVerifiedRulesets = 256

// GovernanceRuleService provides governance rules management operations.
type GovernanceRuleService struct {
	api                *openapi.GovernanceRulesAPIService
	errMapper          *ErrorMapper
	superAdminKeys     []*ecdsa.PublicKey
	minValidSignatures int

	// verified memoises the rulesets whose SuperAdmin signatures already checked out,
	// keyed by rulesetVerificationKey. See verifiedRuleset.
	verifiedMu sync.RWMutex
	verified   map[string]struct{}
}

// GovernanceRuleServiceConfig holds configuration for signature verification.
type GovernanceRuleServiceConfig struct {
	// SuperAdminKeys are the public keys used to verify governance rules signatures.
	SuperAdminKeys []*ecdsa.PublicKey
	// MinValidSignatures is the minimum number of valid SuperAdmin signatures required.
	MinValidSignatures int
}

// NewGovernanceRuleServiceWithVerification creates a new GovernanceRuleService with
// signature verification enabled. Every read verifies the SuperAdmin signatures before
// returning; a service built with no keys fails closed rather than returning unverified
// rules.
//
// There is deliberately no keyless constructor. The previous
// NewGovernanceRuleService ("without verification") is what made the
// `if len(superAdminKeys) > 0` skip in GetDecodedRulesContainer look necessary, and that
// skip returned an unverified container — the document the HSM key that authenticates
// every address is read from — with no error at all.
func NewGovernanceRuleServiceWithVerification(
	client *openapi.APIClient,
	config *GovernanceRuleServiceConfig,
) *GovernanceRuleService {
	svc := &GovernanceRuleService{
		api:       client.GovernanceRulesAPI,
		errMapper: NewErrorMapper(),
		verified:  make(map[string]struct{}),
	}

	if config != nil {
		svc.superAdminKeys = config.SuperAdminKeys
		svc.minValidSignatures = config.MinValidSignatures
	}

	return svc
}

// verifiedRuleset verifies a ruleset's SuperAdmin signatures and returns it, memoising
// the outcome so the same document is not re-verified on every read.
//
//	GetRules / GetRulesByID / GetRulesHistory ─▶ verifiedRuleset ─┬─ memo hit ─▶ return
//	                                                              └─ miss ─▶ ECDSA ─▶ memoise
//
// The container itself is already cached for the address/asset/price paths by
// cache.RulesContainerCache, which fetches through GetDecodedRulesContainer. This memo
// covers only the three reads that return the RAW ruleset, which that cache does not
// serve. Do not add a third cache of the same bytes.
//
// Only SUCCESSES are memoised: a failure must resurface its error on every call. The key
// covers the container AND its signature set, so re-signing the same container is a miss.
func (s *GovernanceRuleService) verifiedRuleset(rules *model.GovernanceRuleset) (*model.GovernanceRuleset, error) {
	if rules == nil {
		return nil, nil
	}

	key, memoisable := rulesetVerificationKey(rules)

	if memoisable {
		s.verifiedMu.RLock()
		_, hit := s.verified[key]
		s.verifiedMu.RUnlock()
		if hit {
			return rules, nil
		}
	}

	if _, err := s.VerifyGovernanceRules(rules); err != nil {
		return nil, err
	}

	// An undecodable container has no stable identity to memoise on; verification above
	// has already had the final say.
	if !memoisable {
		return rules, nil
	}

	s.verifiedMu.Lock()
	// Bounded so a long-running client walking a large history cannot grow it without
	// limit. Governance containers change rarely, so a plain reset beats LRU bookkeeping.
	if len(s.verified) >= maxVerifiedRulesets {
		s.verified = make(map[string]struct{}, maxVerifiedRulesets)
	}
	s.verified[key] = struct{}{}
	s.verifiedMu.Unlock()

	return rules, nil
}

// rulesetVerificationKey identifies the exact document verification cleared: the DECODED
// container bytes — the same bytes VerifyGovernanceRules checks the signatures over — plus
// every signature over them. A container whose signature set changed is a cache MISS
// rather than a stale hit.
//
// Every field is LENGTH-PREFIXED, and that is load-bearing. A plain concatenation is not
// injective: with sha256(container || sig || ...) the boundary between the container and
// the signature list is not committed, so a response-controlling attacker can shift bytes
// across it and make a MODIFIED container collide with a genuine one's key — inheriting
// its "already verified" status and skipping ECDSA entirely. Interleaving a 0x00 separator
// does NOT fix this; only length prefixes do.
//
// userId is deliberately excluded: verification never reads it (see
// helper.VerifyGovernanceRulesSignatures, which consumes only the signature), so a
// server-controlled field that verification ignores must not be able to influence a
// verification-skip decision.
//
// Keying on the decoded bytes rather than the base64 text also removes base64 leniency as
// a lever, and guarantees the key cannot identify something other than what was verified.
// Returns ok=false when the container does not decode, in which case there is nothing
// worth memoising and the caller must fall through to verification.
func rulesetVerificationKey(rules *model.GovernanceRuleset) (string, bool) {
	data, err := base64.StdEncoding.DecodeString(rules.RulesContainer)
	if err != nil {
		return "", false
	}

	h := sha256.New()
	writeLengthPrefixed(h, data)

	// Signature order is server-controlled; sort so an identical set keys identically.
	sigs := make([]string, 0, len(rules.Signatures))
	for _, sig := range rules.Signatures {
		sigs = append(sigs, sig.Signature)
	}
	sort.Strings(sigs)

	writeUint64(h, uint64(len(sigs)))
	for _, sig := range sigs {
		writeLengthPrefixed(h, []byte(sig))
	}
	return string(h.Sum(nil)), true
}

// writeLengthPrefixed writes an 8-byte big-endian length followed by the bytes, so a
// sequence of fields cannot be re-partitioned into a different sequence that encodes
// identically.
func writeLengthPrefixed(h hash.Hash, b []byte) {
	writeUint64(h, uint64(len(b)))
	h.Write(b)
}

func writeUint64(h hash.Hash, n uint64) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], n)
	h.Write(buf[:])
}

// SuperAdminKeys returns the configured SuperAdmin public keys.
func (s *GovernanceRuleService) SuperAdminKeys() []*ecdsa.PublicKey {
	return s.superAdminKeys
}

// MinValidSignatures returns the minimum number of valid signatures required.
func (s *GovernanceRuleService) MinValidSignatures() int {
	return s.minValidSignatures
}

// GetRules retrieves the currently enforced governance rules.
func (s *GovernanceRuleService) GetRules(ctx context.Context) (*model.GovernanceRuleset, error) {
	resp, httpResp, err := s.api.RuleServiceGetRules(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return s.verifiedRuleset(mapper.GovernanceRulesetFromDTO(resp.Result))
}

// GetRulesByID retrieves a governance ruleset by its ID.
func (s *GovernanceRuleService) GetRulesByID(ctx context.Context, id string) (*model.GovernanceRuleset, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	resp, httpResp, err := s.api.RuleServiceGetRulesByID(ctx, id).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return s.verifiedRuleset(mapper.GovernanceRulesetFromDTO(resp.Result))
}

// GetRulesHistory retrieves the history of governance rules with pagination.
func (s *GovernanceRuleService) GetRulesHistory(ctx context.Context, opts *model.ListRulesHistoryOptions) (*model.GovernanceRulesHistoryResult, error) {
	req := s.api.RuleServiceGetRulesHistory(ctx)

	if opts != nil {
		if opts.Limit > 0 {
			req = req.Limit(fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Cursor != "" {
			req = req.Cursor(opts.Cursor)
		}
	}

	resp, httpResp, err := req.Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	// Every historical ruleset is a past ENFORCED document, so each is verified. The
	// memo makes this one ECDSA pass per distinct container rather than one per row.
	//
	//	entries ──▶ verifiedRuleset ──┬─ ok ──────▶ kept
	//	                              └─ error ───▶ excluded + named in the result
	//
	// LENIENT, unlike the whitelist lists, and it does NOT error when nothing survives:
	// a SuperAdmin key rotation makes every pre-rotation ruleset unverifiable, so
	// aborting the page would deny access to the entire audit trail from then on.
	kept := make([]*model.GovernanceRuleset, 0, len(resp.Result))
	var excluded []model.ExcludedRuleset
	for _, rules := range mapper.GovernanceRulesetsFromDTO(resp.Result) {
		if _, err := s.verifiedRuleset(rules); err != nil {
			excluded = append(excluded, model.ExcludedRuleset{
				CreatedAt: rules.CreatedAt,
				Reason:    err.Error(),
			})
			continue
		}
		kept = append(kept, rules)
	}

	result := &model.GovernanceRulesHistoryResult{
		Rules:              kept,
		ExcludedUnverified: excluded,
	}

	if resp.TotalItems != nil {
		if total, parseErr := strconv.ParseInt(*resp.TotalItems, 10, 64); parseErr == nil {
			// Reduced by the exclusions: the server counts rows it returned, the caller
			// receives only those that verified, so the server total would promise a
			// page that can never be fully read.
			result.TotalItems = total - int64(len(excluded))
		}
	}

	if resp.Cursor != nil {
		result.Cursor = *resp.Cursor
	}

	return result, nil
}

// GetRulesProposal retrieves the proposed governance rules.
// Requires SuperAdmin or SuperAdminReadOnly role.
func (s *GovernanceRuleService) GetRulesProposal(ctx context.Context) (*model.GovernanceRuleset, error) {
	resp, httpResp, err := s.api.RuleServiceGetRulesProposal(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	if resp.Result == nil {
		return nil, nil
	}

	return mapper.GovernanceRulesetFromDTO(resp.Result), nil
}

// GetPublicKeys retrieves the list of SuperAdmin public keys.
func (s *GovernanceRuleService) GetPublicKeys(ctx context.Context) ([]*model.SuperAdminPublicKey, error) {
	resp, httpResp, err := s.api.RuleServiceGetPublicKeys(ctx).Execute()
	if err != nil {
		return nil, s.errMapper.MapError(err, httpResp)
	}

	return mapper.SuperAdminPublicKeysFromDTO(resp.PublicKeys), nil
}

// GetDecodedRulesContainer decodes and verifies a GovernanceRuleset's rules container.
// If SuperAdmin keys are configured, it verifies the signatures before returning.
// Returns the decoded rules container or an error if verification fails.
func (s *GovernanceRuleService) GetDecodedRulesContainer(
	rules *model.GovernanceRuleset,
) (*model.DecodedRulesContainer, error) {
	if rules == nil {
		return nil, fmt.Errorf("governance rules cannot be nil")
	}

	if rules.RulesContainer == "" {
		return nil, fmt.Errorf("rules container is empty")
	}

	// Verification is mandatory, never conditional: this container carries the HSM key
	// that address verification trusts, so returning it unverified when no keys are
	// configured would hand back an unauthenticated trust root with no error.
	if _, err := s.verifiedRuleset(rules); err != nil {
		return nil, err
	}

	// Decode the rules container
	return mapper.RulesContainerFromBase64(rules.RulesContainer)
}

// UpdateRulesProposal submits a rules container as a governance proposal.
// Requires SuperAdmin role.
//
// The container is encoded to the wire format internally; the server-controlled
// enforcedRulesHash and timestamp fields are stripped (the server recomputes
// them and rejects submissions asserting stale values). The endpoint returns no
// body; callers that need the persisted proposal should call GetRulesProposal —
// note the server caches rules reads, so an immediate read-back may be stale.
func (s *GovernanceRuleService) UpdateRulesProposal(ctx context.Context, container *model.DecodedRulesContainer) error {
	if container == nil {
		return fmt.Errorf("rules container cannot be nil")
	}

	encoded, err := mapper.RulesContainerToBase64(container)
	if err != nil {
		return fmt.Errorf("failed to encode rules container: %w", err)
	}

	body := openapi.TgvalidatordUpdateRulesProposalRequest{RulesContainer: encoded}
	if _, httpResp, err := s.api.RuleServiceUpdateRulesProposal(ctx).Body(body).Execute(); err != nil {
		return s.errMapper.MapError(err, httpResp)
	}
	return nil
}

// ProposalContainerHash returns the canonical SHA-256 hex digest of a ruleset's decoded
// rules container. This is the value to pass as expectedContainerHash to
// ApproveRulesProposal, and it is what pins the approval to reviewed content.
//
// Not to be confused with the row-to-container label used by the whitelist list paths,
// which is validatord's convention over the base64 TEXT. Both are SHA-256 digests of the
// same document, which is what makes the mix-up silent — this one is client-side only and
// never reaches the wire.
//
// It digests the DECODED bytes, not the base64 text: base64 has multiple encodings of the
// same bytes, so hashing the text would let a re-encoded but byte-identical container read
// as a mismatch. Those decoded bytes are also exactly what gets signed.
func (s *GovernanceRuleService) ProposalContainerHash(rules *model.GovernanceRuleset) (string, error) {
	if rules == nil || rules.RulesContainer == "" {
		return "", fmt.Errorf("ruleset carries no rules container")
	}
	data, err := base64.StdEncoding.DecodeString(rules.RulesContainer)
	if err != nil {
		return "", fmt.Errorf("failed to decode rules container: %w", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

// DecodeProposalForReview decodes a PENDING rules proposal so a SuperAdmin can inspect it
// before approving.
//
// UNVERIFIED BY DESIGN. A pending proposal legitimately carries 0..N signatures — signing
// IS the approval step — so there is no threshold to check yet and GetDecodedRulesContainer
// (which verifies unconditionally) always fails on one. That is why this is a separate,
// explicitly-named entry point rather than a flag: the absence of verification has to be
// visible at the call site.
//
// Pair it with ProposalContainerHash and pass that digest to ApproveRulesProposal, so the
// bytes signed are the bytes reviewed.
func (s *GovernanceRuleService) DecodeProposalForReview(rules *model.GovernanceRuleset) (*model.DecodedRulesContainer, error) {
	if rules == nil || rules.RulesContainer == "" {
		return nil, fmt.Errorf("ruleset carries no rules container")
	}
	return mapper.RulesContainerFromBase64(rules.RulesContainer)
}

// ApproveRulesProposal signs the pending rules proposal with a SuperAdmin
// private key and submits the approval. Requires SuperAdmin role.
//
// expectedContainerHash PINS the content being approved: it must be the
// ProposalContainerHash of the proposal the caller actually reviewed. The proposal is
// re-fetched here, and if its container does not match that digest the call aborts
// WITHOUT signing.
//
// That pin is the security control. Without it a server able to shape responses could
// serve the benign proposal to the review call and a different container to the re-fetch
// inside this method, obtaining a GENUINE SuperAdmin signature over bytes of its choosing
// — and such a container then verifies clean everywhere, including in independent and
// air-gapped verifiers that never trusted that server. Nothing else binds the two calls:
// the wire format carries no proposal id, hash or version, and a pending proposal has no
// signatures to verify against, so re-verification cannot substitute for pinning.
//
// The signature is computed over the decoded bytes of the pending proposal's rules
// container (SHA-256 + P-256 ECDSA, base64 raw r||s).
//
// Review flow:
//
//	proposal, _ := svc.GetRulesProposal(ctx)
//	container, _ := svc.DecodeProposalForReview(proposal)  // inspect this
//	pin, _ := svc.ProposalContainerHash(proposal)
//	err := svc.ApproveRulesProposal(ctx, key, "lgtm", pin)
func (s *GovernanceRuleService) ApproveRulesProposal(
	ctx context.Context,
	privateKey *ecdsa.PrivateKey,
	comment string,
	expectedContainerHash string,
) error {
	if privateKey == nil {
		return fmt.Errorf("private key cannot be nil")
	}
	if expectedContainerHash == "" {
		return fmt.Errorf("expectedContainerHash is required: it pins the approval to the " +
			"container you reviewed (see ProposalContainerHash)")
	}

	proposal, err := s.GetRulesProposal(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch pending rules proposal: %w", err)
	}
	if proposal == nil || proposal.RulesContainer == "" {
		return fmt.Errorf("no pending rules proposal to approve")
	}

	actualHash, err := s.ProposalContainerHash(proposal)
	if err != nil {
		return fmt.Errorf("failed to hash pending rules container: %w", err)
	}
	if subtle.ConstantTimeCompare([]byte(actualHash), []byte(expectedContainerHash)) != 1 {
		return fmt.Errorf(
			"refusing to sign: the pending rules proposal changed since it was reviewed "+
				"(reviewed %s, pending %s)", expectedContainerHash, actualHash)
	}

	data, err := base64.StdEncoding.DecodeString(proposal.RulesContainer)
	if err != nil {
		return fmt.Errorf("failed to decode pending rules container: %w", err)
	}

	signature, err := crypto.SignData(privateKey, data)
	if err != nil {
		return fmt.Errorf("failed to sign rules container: %w", err)
	}

	body := openapi.TgvalidatordApproveRulesProposalRequest{Signature: signature, Comment: comment}
	if _, httpResp, err := s.api.RuleServiceApproveRulesProposal(ctx).Body(body).Execute(); err != nil {
		return s.errMapper.MapError(err, httpResp)
	}
	return nil
}

// RejectRulesProposal rejects the pending rules proposal with a comment.
// Requires SuperAdmin role.
func (s *GovernanceRuleService) RejectRulesProposal(ctx context.Context, comment string) error {
	body := openapi.TgvalidatordRejectRulesProposalRequest{Comment: comment}
	if _, httpResp, err := s.api.RuleServiceRejectRulesProposal(ctx).Body(body).Execute(); err != nil {
		return s.errMapper.MapError(err, httpResp)
	}
	return nil
}

// VerifyGovernanceRules verifies the SuperAdmin signatures on the governance rules and
// returns the verified rules, so a caller can chain verification into an assignment
// rather than relying on the error alone. Java and Python return the rules the same way.
func (s *GovernanceRuleService) VerifyGovernanceRules(rules *model.GovernanceRuleset) (*model.GovernanceRuleset, error) {
	if len(rules.Signatures) == 0 {
		return nil, &model.IntegrityError{Message: "no signatures provided for governance rules"}
	}

	// Decode the rules container data
	rulesData, err := base64.StdEncoding.DecodeString(rules.RulesContainer)
	if err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("failed to decode rules container: %v", err),
		}
	}

	// Convert signatures to the format expected by the verifier
	signatures := make([]*model.RuleUserSignature, len(rules.Signatures))
	for i := range rules.Signatures {
		signatures[i] = &model.RuleUserSignature{
			UserID:    rules.Signatures[i].UserID,
			Signature: rules.Signatures[i].Signature,
		}
	}

	// Verify signatures
	if err := helper.VerifyGovernanceRulesSignatures(
		rulesData,
		signatures,
		s.superAdminKeys,
		s.minValidSignatures,
	); err != nil {
		return nil, &model.IntegrityError{
			Message: fmt.Sprintf("governance rules signature verification failed: %v", err),
		}
	}

	return rules, nil
}
